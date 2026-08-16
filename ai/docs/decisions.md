# Design Decisions

A running log of the non-obvious modelling decisions behind `tge` — the *why*
that the code alone doesn't explain. See [domain.md](domain.md) for the entities
themselves and [architecture.md](architecture.md) for the layering.

## 1. A power system is a DAG, not a tree

**Decision.** `powersystem.PowerSystem` holds `Nodes map[string]*PowerNode` — a flat
map keyed by node ID — and the shape comes from edges stored *on the nodes*
(`Parents`, `Siblings`, `MutuallyExclusive`). `AddEdge` with `EdgeParent` runs a DFS up
the parent chain and returns `ErrCyclicDependency` if the edge would close a cycle.

**Why.** The original model was a literal Go tree (`Power{Name, Children}`). Real power
systems are not trees:

- A technique can require **several** prerequisites at once (multi-parent) — a tree
  forces you to pick one parent and lie about the rest.
- Branches are frequently **mutually exclusive** (pick Kaioken *or* Super Saiyan), which
  has no representation in a parent/child tree at all.
- Traversal and "does this character have node X?" checks want O(1) lookup by ID, not a
  recursive walk.

The flat map + explicit edges gives all three, and the cycle check is what keeps it a
DAG rather than a general graph — progression has to have a direction.

**Consequence.** `PowerNode.ID` is a derived slug (lowercase, spaces → `_`), so IDs are
stable and predictable from names, and the CLI/seeder can reference nodes without a
lookup table. Rendering a system (`tge powersystem show`) has to reconstruct a display
tree from the parent edges and mark already-visited nodes `(see above)`, because a
multi-parent node genuinely appears in more than one place.

## 2. Progression state lives on the character, not on the system nodes

**Decision.** A character carries `MechanicState` (character-level) plus
`UnlockedNodes []NodeProgress` (per node). The power system a character references is
the **shared definition** and holds no per-character state at all.

**Why.** The earlier design gave `Character` two `[]PowerSystem` fields — definitions in
`PowerSystems` and progressed copies in `Power`, with per-node state hung on a node's
`State PowerState` interface. That copied the entire tree per character just to annotate
a few nodes, and made "which nodes has this character unlocked?" a full traversal of a
duplicated structure. Flipping it around, the character holds a short list of
`{System, NodeID, Level, Progress, BasePower}` records, and the system stays a single
shared, immutable definition.

`MechanicState` exists alongside it because a lot of progression is **not** per-node:
tier, awakening, alignment, energy pools, spell slots, vows, permanent traits. Those are
facts about the character across every system it participates in, so they belong on the
character once rather than being duplicated into each system's state type.

**Consequence.**

- `Character.Systems` degraded from `[]PowerSystem` to `[]string` — it is membership by
  name, and the definitions are loaded from the repository when needed.
- `PowerValue` is a pure derivation, recomputed by `CalculateTotalPower()` on every
  mutation: `(MechanicState.BasePower + Σ node.BasePower × node.Level) × Π species.Power`.
  Note that **every** species multiplies in, which is what makes hybrids compound.
- The breakthrough curve moved onto `NodeProgress` as `100 × Level²` — a single global
  curve rather than per-realm authored thresholds.

## 3. Characters and power systems persist as JSON aggregates, not SQL rows

**Decision.** `Character` and `PowerSystem` moved out of SQLite into `adapter/file`,
one JSON document per entity (`data/characters/<slug>.json`,
`data/power_systems/<slug>.json`). Everything else stayed in SQLite.

**Why.** Both are object graphs with open-ended nested collections — a character has
unlocked nodes, idle slots, energy pools, spell slots, an inventory, a stats block and a
multi-species list; a system has a node map where each node has three edge lists and two
maps. Modelling that relationally meant a table (and a load/save fan-out) per nested
collection, and every new field on `MechanicState` was a migration. It was also the
direct cause of the old "not persisted yet" gaps: extra species, system kind, nested
non-root nodes and node state were all silently dropped on save because no column
existed for them.

Serializing the aggregate root collapses `CharacterRepository` to
`Save`/`FindByName`/`List`, and the persistence gap closes by construction: if it's a
field, it round-trips.

**Consequence.**

- Writes are atomic via write-to-`.tmp` + `rename`. `Clean` soft-deletes by renaming to
  `.deleted` rather than unlinking, so a mistaken `character clean` is recoverable.
- The filename is a slug of the name, so names differing only in case or spacing collide
  on disk. `CreateCharacter` sidesteps this by never rejecting a duplicate name — it
  appends `" (1)"`, `" (2)"`, … until the slug is free.
- `List` reads and decodes every file; fine at authoring scale, linear in character count.
- The `characters`, `character_*`, `power_systems` and `powers` tables still exist in the
  migrations but are dead (see [architecture.md](architecture.md#database-migrations)).

## 4. Idle progression runs on a per-character novel clock

**Decision.** `Character.NovelTime` (seconds) is the character's own clock, advanced only
by `tge character pass-time`. Idle activities are `IdleSlot`s stamped with a `StartTime`
on that clock, and gains are computed as elapsed hours × rate.

**Why.** Wall-clock time is wrong for an authoring tool: the author is writing a story
where months pass between chapters, and the CLI is invoked sporadically. Anchoring to an
explicit in-story clock makes progression *deterministic and replayable* — the same
sequence of commands always yields the same character — and lets the author skip forward
by narrative amounts (`--days`, `--hours`, `--minutes`) instead of waiting.

**Consequence.**

- Two different reads of "how much has idling produced", and they **disagree**.
  `IdleService.CommitOfflineGains` is the write path: it calls `AdvanceNode` and banks only
  the leftover into `MechanicState.EnergyPools` (pool name `"<system>_<power>"`), then
  expires finished slots. `CurrentEnergyPools(now)` is the non-mutating projection `tge
  status` uses, and it adds `elapsed × Rate` straight to the pool with no `AdvanceNode`
  step. The discrepancy is currently invisible only because the projection's slot term is
  dead: every path that moves `NovelTime` commits first, and a commit re-stamps
  `StartTime = NovelTime`, so `now > slot.StartTime` is never true at read time. Worth
  reconciling before either side changes.
- A slot with `Duration <= 0` is indefinite and is re-stamped to the current `NovelTime`
  on every commit, so it accrues continuously without unbounded catch-up.
- `TotalSlots` (1 for a new character) caps concurrency — the knob a future "unlock more
  idle slots" progression would turn.
- The generation rate is currently hard-coded to `10.0` points/hour in `IdleService`;
  making it a function of aptitude/technique is the obvious next step.

## 5. `PowerState`, `CultivationState` and `SuperPowerState` are staged, not wired

**Decision.** `power.PowerState` (`Kind() PowerSystemType` + `Power() float64`) and its
two implementations — `cultivation.CultivationState` (Realm + Level + Points + Progress,
with realm-aware `AdvanceWithin`) and `superpower.SuperPowerState` (tier 0–5) — are kept
and **not referenced by any service or adapter**. The `cultivation` package is well
covered by unit tests; `superpower` has **no test file at all** and is referenced by
nothing, tests included.

**Why.** They encode the realm/level progression rules the cultivation domain actually
wants (per-level breakthrough gates, per-realm caps that differ for the main character,
leftover points carried into the next realm) at a fidelity `NodeProgress`'s single
`100 × Level²` curve does not. Deleting them would throw away the modelling; wiring them
in now would mean two competing progression paths on the same character. They stay as the
target design for the next progression pass.

**Consequence.** `tge realm` / `tge realm create-level` author realms and levels that
nothing currently consumes at runtime — the seeded 9×9 Spirit ladder is data waiting for
its engine. Likewise `cultivation.CultivationPath` (`Body`/`Spirit`/`Soul`) is modelled as
an open interface rather than an enum (so each path can accrue attributes later) but is
not yet referenced. `power.Power` — the old tree node with a `State PowerState` field — is
the one piece with no remaining role and is a straightforward removal.

## 6. Class, Profession, Species, Gender are entities/enums, not loose strings

**Decision.** On `Character`: `Class` is an `rpg.Class`, `Profession` an
`rpg.Profession`, `Species` a `[]character.Species`, and `Gender` the `Gender`
enum. The `CharacterService` resolves the class/profession entities from their
repositories (instead of discarding the lookup result) and passes them through.

**Why.** Carrying the resolved entity rather than a name keeps the character
self-describing and avoids re-lookups downstream. `Species` is a list because a
character may acquire additional species "by any means" (hybrids); creation still
starts from a single species.

**Consequence.** Now that characters serialize as JSON aggregates (§3), the full species
list round-trips — this was the one gap the SQL schema could not express.

## 7. Migrations on goose

Schema is managed with **goose v3** as an embedded library (migrations in
`internal/adapter/sqlite/migrations/*.sql`, applied on `open`). This replaced a
hand-rolled `CREATE TABLE IF NOT EXISTS` bootstrap that could only add tables, not
evolve existing ones. See [architecture.md](architecture.md#database-migrations)
for the full setup; column-level schema changes (adding `kind` to `power_systems`, adding
`grade` across the RPG tables) go in new numbered migration files.

## 8. Two output channels: developer log vs. system log

**Decision.** `internal/logger` splits output. `Dev` goes to **stderr** with a `[DEV]`
prefix; `System` goes to **stdout** through a template (`[Ding! %s]`, overridable via
`system_log.template` in the defaults YAML).

**Why.** `tge` writes for two readers at once. The author wants the in-fiction beat —
*"[Ding! System binding complete. Welcome, Host Lin Feng!]"* — as copy they can lift
straight into the novel. The developer wants diagnostics. Putting them on different
streams with different formats means the narrative channel can be redirected, restyled
per-setting, or piped, without either polluting the other. `character create` and
`character train-node` print both, under explicit `--- Novel Log` / `--- Internal System
Log` headers.

## 9. Known gaps

- **Two progression models coexist.** `NodeProgress` (shipped) and `CultivationState`
  (staged, §5). Realms/levels are authored but unused at runtime.
- **`train-node` does not train.** It validates prerequisites (parents unlocked, no
  mutually-exclusive conflict, not already unlocked) and then grants the node outright at
  `Level 1` / `Progress 100`. It cannot be called twice on the same node, so there is no
  path to level 2+ except idle accrual through `AdvanceNode`.
- **Grades are inert.** `defaults.yml` defines a Common→Mythic ladder with power
  multipliers and tags catalog entries with grade names, but `seedCatalog` never passes
  `Grade` through to the create inputs, and nothing reads `power_multiplier`.
- **`powersystem add-power` is advertised but missing.** The sub-command appears in the
  help text and has no case in the dispatch switch; nodes and edges can only be created by
  the seeder.
- **No CLI for universe realms.** `AddRealms` exists on the domain, port and service, but
  only `DefaultWorldService` calls it.
- **Realm "first"/"next" ordering.** `Realm.Tier` exists and the seeder assigns it from
  catalog order, but nothing resolves "the next realm" through it yet.
- **`List` is O(files).** The flat-file repositories decode every document to list.

## 10. GitOps and Semantic Versioning (Release Please)

**Decision.** Replaced manual Git tags with Google's `release-please-action`.

**Why.** Automates semantic versioning (SemVer) and changelog generation using Conventional Commits (`feat:`, `fix:`, etc.). It opens a "Release PR" accumulating all merged changes. Merging this PR instantly cuts a Git tag and a GitHub Release, which triggers GoReleaser to attach cross-platform binaries. This decouples merging features from publishing releases while eliminating manual version bookkeeping.
