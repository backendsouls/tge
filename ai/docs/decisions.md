# Design Decisions

A running log of the non-obvious modelling decisions behind `tge` — the *why*
that the code alone doesn't explain. See [domain.md](domain.md) for the entities
themselves and [architecture.md](architecture.md) for the layering.

## 1. A character's power has two facets: definition vs. instance

**Decision.** A `Character` holds **two** `[]progression.PowerSystem` fields:

- `PowerSystems` — the complete tree **definitions** the character has access to
  (its membership/catalogue). Authored once, shared, lives in the cosmology.
- `Power` — the character's own **progressed instances** of those systems: the
  same `PowerSystem` type, but with per-node progress attached.

**Why.** `PowerSystems` answers *"in which systems can this character grow at
all?"*; `Power` answers *"how strong is it, and how did it get there?"*. Keeping
them as the same type (rather than a bespoke progress struct) means a new kind of
system — Magic vs. Cultivation — slots into `Power` with no shape change. This is
the classic definition-vs-instance / template-vs-state split.

**Consequence.** `PowerSystems` is the renamed former `Systems` field; the former
`Power string` (the rendered numeric value) was renamed to `PowerValue` to free
the `Power` name for the instance list.

## 2. Per-node progress is polymorphic (`SystemKind` + `PowerState`)

**Decision.** Progress is **not** typed as `Cultivation`. Instead:

- `PowerSystem` carries a `Kind` (`SystemKind`: `Cultivation`, `Magic`, …).
- A `Power` tree node carries an optional `State PowerState` — an **interface**,
  `nil` in a definition and set only in a character's instance.
- `CultivationState` is the first `PowerState`: `Realm`, current `Level`,
  `Points` accumulated toward `Level.BreakthroughPoints`, and `Progress`.

**Why.** If the state were typed `Cultivation`, a Magic system would be stuck — it
doesn't cultivate through realms/levels. The interface keeps `Character.Power`
general over the *kind* of system, so a future `MagicState` implements
`PowerState` and reuses the same `Power`/`PowerSystem` machinery unchanged. This
mirrors how `Path` was already modelled as an open type for Open/Closed reasons.

**Consequence.** The standalone `Cultivation{Path, Realm, Progress}` struct was
replaced by `CultivationState`. The "path" is no longer a field — it is simply the
tree node the state hangs under. `Path`/`Body`/`Spirit`/`Soul` in `path.go` are
now unused by progression and are a candidate for later removal.

## 3. Class, Profession, Species, Gender are entities/enums, not loose strings

**Decision.** On `Character`: `Class` is an `rpg.Class`, `Profession` an
`rpg.Profession`, `Species` a `[]character.Species`, and `Gender` the `Gender`
enum. The `CharacterService` resolves the class/profession entities from their
repositories (instead of discarding the lookup result) and passes them through.

**Why.** Carrying the resolved entity rather than a name keeps the character
self-describing and avoids re-lookups downstream. `Species` is a list because a
character may acquire additional species "by any means" (hybrids); creation still
starts from a single species, whose values seed `PowerValue`/`Lifespan`.

## 4. `Cultivate` — first vertical slice (IMPLEMENTED), growth rates (PENDING)

The goal that motivated the reshape. Built as a full vertical slice:

- **Driving port + service.** `CharacterService.Cultivate(port.CultivateInput)`
  **sets (upserts) a character's cultivation state at one power node**: `System`
  (defaults to the character's first power system), `Path` (defaults to the
  `System` name), anchored to a `Realm` + `Level` (plus `Points`/`Progress`). It
  saves, then returns the reloaded character.
- **CLI.** `tge character cultivate --name … [--system …] [--path …] [--realm …]
  --level N`. The CLI resolves the realm + level via `a.realms` and passes
  concrete names in, so `CharacterService` gained **no** new dependency.
- **Persistence.** A `character_cultivations` table (one row per
  `(character, system, path)`), upserted by `CharacterRepository.SaveCultivation`
  and rebuilt into `Character.Power` by `loadCultivations` on every load. This
  closes most of the old "Power not persisted" gap (see §5).
- **Status rendering.** `tge status` prints the `Power:` tree grouped by system
  `Kind`, each node showing its `CultivationState` Realm/Level.

**Still pending — the growth-rate mechanics.** The current `Cultivate` sets state
directly (you pass the target realm/level); it does **not** yet accumulate points
or break through. The agreed-but-unbuilt design:

- **Growth rate** (`Normal`, `Overpowered`, `Insane`, `Slow`, `None`): a
  multiplier on base points gained. `None` = 0 (cannot progress).
- **Rate source** (decided): **aptitude default with per-call override** — a
  character carries a default rate (a "talent"); a call may override it for a
  one-off boost (a pill/event).
- **Mechanics**: add `rate × base` to the node's `Points`; on reaching
  `Level.BreakthroughPoints`, advance to the next `Level` (respecting the realm's
  per-tier caps, higher for the main character).

## 5. Deferred / known gaps

- **Power persistence — mostly closed.** Cultivation state now round-trips via
  `character_cultivations` (§4). Still not persisted: `PowerSystem.Kind` (assumed
  `Cultivation` on load) and any **nested** (non-root) power nodes — only
  root-level `(system, path)` nodes are stored.
- **"First realm" ordering.** `tge character cultivate` defaults `--realm` to the
  first realm from `ListRealms`, which is ordered **alphabetically** (e.g. resolves
  to `Nascent Spirit`, not the tier-1 `Spirit Gathering`). `Realm` has no
  tier/sequence field, so "first" can't mean "lowest tier" yet — pass `--realm`
  explicitly, or add realm ordering (domain + schema + seed + `ORDER BY`).
- **Additional species.** Only the **primary** species name is persisted (single
  `species` column); extra species a character gains are dropped on save.
- **`path.go`.** Retained and still tested, but unused by `CultivationState`;
  revisit whether tree nodes fully subsume the `Path` concept.

## 6. Migrations on goose

Schema is managed with **goose v3** as an embedded library (migrations in
`internal/adapter/sqlite/migrations/*.sql`, applied on `open`). This replaced a
hand-rolled `CREATE TABLE IF NOT EXISTS` bootstrap that could only add tables, not
evolve existing ones. See [architecture.md](architecture.md#database-migrations)
for the full setup; column-level schema changes (e.g. renaming the `power` column
to match `PowerValue`, widening for multi-species) now go in new numbered
migration files.

## 7. GitOps and Semantic Versioning (Release Please)

**Decision.** Replaced manual Git tags with Google's `release-please-action`.

**Why.** Automates semantic versioning (SemVer) and changelog generation using Conventional Commits (`feat:`, `fix:`, etc.). It opens a "Release PR" accumulating all merged changes. Merging this PR instantly cuts a Git tag and a GitHub Release, which triggers GoReleaser to attach cross-platform binaries. This decouples merging features from publishing releases while eliminating manual version bookkeeping.
