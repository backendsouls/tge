# Domain Models

The domain for `tge` centers around creating and managing elements of a Cultivation-style game or novel.

Bounded contexts live under `internal/core/domain/`: `character`, `cosmology`,
`novel`, `powersystem`, `power`, `cultivation`, `superpower` and `rpg`.

## Key Entities

1. **Character** (`domain/character`): A being in the story, born as a "Mortal" with a `Mortal{Age, Lifespan}` block.
   - **Type**: `MainCharacter`, `SideCharacter`, `SupportCharacter`, `Hero`, `Heroine`.
   - **Gender**: `Male` or `Female` (a `Gender` enum).
   - **Species**: `Species []character.Species` — a character may be of one or more species (e.g. a hybrid gained "by any means"); a fresh character is created with a single one. The species' `Lifespan` seeds `Mortal.Lifespan`, and **every** species' `Power` multiplies into the total power (see below).
   - **Rules**: A `Hero` requires a `Female` `MainCharacter` to exist, and a `Heroine` requires a `Male` `MainCharacter` (`CheckRole`).
   - **Name collisions**: `CharacterService.CreateCharacter` never rejects a duplicate name — it appends `" (1)"`, `" (2)"`, … until the name is free.
   - **Power** is held in three parts:
     - `Systems []string` — the **names** of the power systems the character has access to (its membership). Added at creation (`--system`, repeatable) or later via `tge character add-power`.
     - `MechanicState power.MechanicState` — the character-level, system-agnostic progression state (see §4).
     - `UnlockedNodes []NodeProgress` — per-node progress inside those systems (see §5).
     - `PowerValue string` — the **rendered** total, recomputed by `CalculateTotalPower()` on every mutation:
       `(MechanicState.BasePower + Σ node.BasePower × node.Level) × Π species.Power`, floored at `1.0`.
   - **Idle / novel time**: `NovelTime int64` is the character's own clock in seconds within the story, advanced only by `tge character pass-time`. `IdleState{Slots []IdleSlot, TotalSlots int}` holds the background activities running against that clock (see §5).
   - **RPG attributes** (see the RPG context below): an optional `Class` (`rpg.Class`) and `Profession` (`rpg.Profession`) — held as the entity values, validated against and resolved from existing entities when given — a `Stats` block (defaults to `rpg.BaseStats()`, `0.65` across the board) and an `Inventory` of items (added via `tge character give-item`; stacks merge by item).
   - **Persistence**: the whole `Character` is one JSON aggregate on disk (`data/characters/<slug>.json`), so *every* field above round-trips, including the multi-species list.

2. **Species** (`domain/character`): A biological classification defining base status values and a per-species default.
   - `Power` is a **multiplier** into `CalculateTotalPower()`; `Lifespan` seeds the character's `Mortal.Lifespan`. Because it multiplies, a species with `Power < 1` (Human, at `0.65`) makes a character *weaker* than its raw base — the reason a fresh mortal's total is floored at `1.0`.
   - **The Human base is defined twice, inconsistently.** `defaults.yml` has `character.species.power: 1` (the fill-in template passed to `DefaultWorldService`), while `catalog.species[Human].power: 0.65` — and `character.HumanBase()` in code also says `0.65`. The catalog wins: `seedCatalog` runs before any command, inserts Human at `0.65`, and `EnsureDefaults` then hits `ErrSpeciesExists` and silently keeps the seeded row. So the effective Human `Power` is **0.65** and the `character.species.power: 1` key is dead. Worth reconciling.
   - **DefaultGender** (optional, `Male`/`Female`): when a character of this species is created without a gender, it falls back to the species' default (then to the global default from `defaults.yml`). The built-in **Human** base species (`character.HumanBase()`) defaults to `Male`. Seven species ship in the seed catalog (Human, Demon, Spirit, Beastkin, Dragon, Fae, Elf), each with its own default gender; more can be added via `tge species create --default-gender …` or the defaults YAML.

3. **PowerSystem** (`domain/powersystem`): A named **DAG** (directed acyclic graph) of `PowerNode`s — *not* a tree. Nodes are held in a flat `map[string]*PowerNode` keyed by node ID for O(1) lookup.
   - **PowerSystemType**: the family the system belongs to — `Cultivation`, `Magic`, `SuperPower`, `Reiatsu` or `Gamer`. `NewPowerSystem` defaults to `Cultivation` when the kind is blank and rejects an unknown one. The type drives presentation (`tge status` shows a `Gamer` character's `Stats`, and labels an un-progressed `Cultivation` node "Mortal" vs. "Unawakened" elsewhere).
   - **PowerNode**: `ID` (a slug derived from the name — lowercased, spaces → `_`), `Name`, `Category`, `Tags`, plus the graph edges (`Parents`, `Siblings`, `MutuallyExclusive`) and the mechanical payload: `BasePower`, `StatVector map[string]float64`, `MaterialReq map[string]int` and `Drawbacks []string`.
   - **Edges** (`AddEdge(nodeID, targetID, EdgeType)`): `parent`, `sibling` or `mutually_exclusive`. Adding a **parent** edge runs a DFS up the parent chain and rejects the edge with `ErrCyclicDependency` if it would close a cycle — this is what makes the graph a DAG. A `mutually_exclusive` edge is written **symmetrically** onto both nodes; a node cannot be mutually exclusive with itself.
   - **Multi-parent** is the point of the DAG: a node may require several prerequisites at once, which a tree cannot express.
   - **Persistence**: a system is one JSON aggregate on disk (`data/power_systems/<slug>.json`), so the whole graph — kind, nodes and all edges — round-trips.

4. **MechanicState** (`domain/power`): the character-level, system-agnostic progression tracker that replaced per-system state structs on the aggregate.
   - `Tier int` (non-negative), `BasePower float64`, `IsAwakened bool`, `Alignment float64` (bounded `-100..100` by `SetAlignment`), `EnergyPools map[string]int`, `SpellSlots map[int]int`, `PermanentTraits []string`, `Vows []string`.
   - A new character starts at tier `0` with `BasePower` `1.0`.
   - Energy pools are how idle activity banks unspent points: the pool name is `"<system>_<power>"`.

5. **NodeProgress & idle progression** (`domain/character`):
   - **NodeProgress**: `{System, NodeID, Level, Progress, BasePower}` — a character's standing at one node of one system. `CalculateBreakthrough()` is the gate for the current level: `100 × Level²`. `Advance(points)` pours points into `Progress`, levels up whenever the gate fills (resetting `Progress`), and returns any unconsumed remainder — which for an unlocked node is always `0`, since the loop runs until the points are spent.
     Two quirks follow from the arithmetic. `train-node` grants a node at `Level 1`/`Progress 100`, which is *exactly* the level-1 gate (`100 × 1²`), so the very next `Advance` call promotes it to Level 2 before spending a single point. And at `Level 0` the gate is `0`, so a level-0 node likewise jumps to Level 1 for free.
   - **IdleSlot**: `{StartTime, Duration, System, Power, Rate}` — a background activity started at a `NovelTime` (seconds), lasting `Duration` **hours** (`<= 0` means indefinite), generating `Rate` points/hour (currently hard-coded to `10.0` by `IdleService`).
   - **`TotalSlots`** caps concurrent activities (1 for a new character); assigning past the cap fails.
   - `IdleService.CommitOfflineGains` is the **write** path: it converts elapsed hours to points, pours them through `Character.AdvanceNode`, banks only the *leftover* into `MechanicState.EnergyPools`, drops finished slots, and re-stamps surviving ones to the current `NovelTime`.
   - `CurrentEnergyPools(now)` is a **read-only projection** used by `tge status`: base pools plus `elapsed × Rate` for each running slot. Note it does **not** agree with the write path — it adds the raw points straight to a pool, whereas a commit would have spent most of them on node progress. In practice the discrepancy never surfaces: every path that moves `NovelTime` (`pass-time`) commits first, and a commit re-stamps `StartTime = NovelTime`, so by the time `status` runs `now > slot.StartTime` is false and the slot term contributes `0`. The loop is effectively dead code.

6. **Reality** (also known as the **Box**) (`domain/cosmology`): The outermost collection, grouping one or more `Omniverse`s together.
   - Member omniverses are referenced by name and must already exist.
   - An omniverse belongs to at most one reality.

7. **Omniverse**: A collection that groups one or more `Multiverse`s together.
   - Member multiverses are referenced by name and must already exist.
   - A multiverse belongs to at most one omniverse.

8. **Multiverse**: A collection that groups one or more `Universe`s together.

9. **Universe**: A collection grouping multiple `PowerSystem`s together and containing in-universe `Realm`s (`cosmology.Location`).
   - Systems belong exclusively to one Universe.
   - Universes can be grouped under a Multiverse.
   - Realms here are **locations/bubbles** within a Universe — a `cosmology.Location`, entirely unrelated to the cultivation `Realm` below. They share only the word.
   - `AddRealms` exists on the domain, port and service, but **nothing calls the service method** — there is no CLI command, and `DefaultWorldService` writes the default realm through the repository's `SaveRealms` directly. It is unreachable code today.

10. **Realm (Cultivation Stage)** (`domain/cultivation`): A single stage of cultivation defining how power and lifespan grow.
    - Uses linear equations `ax + b` (`Multiplier * x + Adder`) for calculating current `Power` and `Lifespan`.
    - **Tier**: orders realms into a sequence (`1` = lowest, `0` = unordered). The seeder assigns tiers from catalog order.
    - **Levels**: A realm is subdivided into ordered `Level`s (the `Cultivation → Realm → Level` hierarchy), e.g. the "First Level" of the "Spirit Gathering" realm. Each `Level` has a positive `Number` (unique within the realm), a `Name`, and `BreakthroughPoints` to advance to the next level. Levels are kept sorted by number. Managed via `tge realm create-level --breakthrough …` / `tge realm show`.
    - **Per-tier level caps**: A realm caps how many levels a character may reach, and the **main character gets a higher cap** than a normal character — e.g. `MaxLevels` 9 for normal characters, `MainCharacterMaxLevels` 13 for the main character (`--max-levels 9 --max-levels-main 13`). `0` = unlimited; an unset main cap inherits the normal cap. `MaxLevelsFor(isMain)` returns the applicable cap. A realm may define levels up to the highest (main) cap; adding a level whose number exceeds that is rejected, so the realm holds at most that many levels.
    - **CultivationState**: the `Cultivation` kind's `power.PowerState` — `{Realm, Level, Points, Progress}`, with `Ready()` (gate full) and `AdvanceWithin(points)` (fills the gate, breaks through level by level, returns the leftover once the realm's last level is full so the caller can carry it into the next realm). `Power()`/`Lifespan()` evaluate the realm's formulas at `Progress`.
      **Status: not wired in.** The shipped progression path is `NodeProgress` (§5); `CultivationState` is complete and tested but no service or adapter constructs one yet. See [decisions.md](decisions.md) §5.
    - **CultivationPath** (`Body`, `Spirit`, `Soul`): modelled as an open interface rather than an enum so each path can accrue its own attributes later (OCP). Also not referenced by the shipped progression path.
    - **Default catalog**: the seeded cultivation is **9 realms** (Spirit Gathering → Spirit Sovereign), **each with 9 levels** (First … Ninth Level). In the defaults YAML a realm's `level_count` auto-generates that many ordinally-named levels (with `level_breakthrough_points` each) instead of listing them.
    - **Persistence**: realms and their levels live in SQLite (`realms`, `realm_levels`).

11. **Novel** (`domain/novel`): A story containing `Volumes` and `Chapters`, led by a single `MainCharacter`.
    - A `MainCharacter` can only be the lead of one `Novel`.
    - `Volumes` and `Chapters` have strict ordering and uniqueness constraints.

12. **Timeline** (`domain/cosmology`): An ordered sequence of `Event`s (each an `Order` + `Description`) owned by a **location**.
    - Every location — `Realm` (location), `Universe`, `Multiverse`, `Omniverse` and `Reality` (the Box) — owns exactly one Timeline. A location is referenced by a `port.LocationRef` (kind + name; a realm is also scoped by its `Universe`, since realm names are unique only within a universe).
    - A Timeline is provisioned automatically with a derived default name (`defaultTimelineName` = `"<location> Timeline"`) whenever a location is created. `CreateReality`/`CreateOmniverse`/`CreateMultiverse`/`CreateUniverse` each call `ensureTimeline`, as does the default-world provisioner for all five of its locations. `UniverseService.AddRealms` would do the same for a realm — but nothing calls it, so the only realm timeline that ever exists is the default world's "Mortal Realm".
    - `Event`s are unique by `Order` within a Timeline and are kept sorted ascending.
    - A Timeline is **not** the same clock as `Character.NovelTime`: timelines are authored narrative events, `NovelTime` is the simulation clock idle progression runs against.

13. **SuperPowerState** (`domain/superpower`): a second `power.PowerState` implementation — a `Tier` from 0 to 5 mapping to a fixed power (`1, 5, 10, 20, 50, 100`). Currently **unreferenced** by any service or adapter; kept as the worked example that `PowerState` generalises beyond cultivation. See [superpowers.md](superpowers.md).

## The default world

Creating a `MainCharacter` (including from just a name) provisions a default cosmology
in the background, **idempotently**, so repeated creations reuse the same single world:

- `Reality` "The Box" → `Omniverse` "Origin Omniverse" → `Multiverse` "Origin Multiverse"
  → `Universe` "Origin Universe" → `Realm` (location) "Mortal Realm", each with its Timeline.
- A default `PowerSystem` of kind `Cultivation`, named by `world.power_system` in
  `defaults.yml` (shipped value: **"Cultivation"**). The built-in fallback in code, used
  only when configuration supplies nothing, is `"Mortal Path"`.
- The **Human base** species as the initial-status template.

Explicitly provided type/gender/species/systems are honoured. Note that the character is
*not* auto-joined to the default power system — `Systems` is only what `--system` passed.

## RPG Context (`internal/core/domain/rpg`)

A separate bounded context of role-playing-game building blocks. Each identity-bearing entity is keyed by `Name` and has a full vertical slice (domain → port → service → SQLite → CLI `create`/`list`/`show`).

14. **Stats**: A value block of eight `float64` attributes — `STR`, `AGI`, `INT`, `VIT`, `DEX`, `WIS`, `CHA`, `LUK` — all non-negative. `BaseStats()` is the default starting spread (`0.65` each); `Add` combines blocks (e.g. base + equipment bonus). Stats are a component of `Character`, serialized with it.
15. **Ability**: An innate power (`Name`, `Description`, `Grade`).
16. **Skill**: A learned, trainable ability (`Name`, `Description`, `Grade`).
17. **Item**: A thing that can be held in an inventory (`Name`, `Description`, `Grade`).
18. **Effect**: A modifier or condition with a `Kind` (`Buff`, `Debuff`, `Status`).
19. **Equipment**: An equippable item with a `Slot` (`Weapon`, `Armor`, `Accessory`) and a `Stats` bonus.
20. **Profession**: A vocation (e.g. Blacksmith), with a `Grade`; a character may have one.
21. **Class**: A combat archetype (e.g. Warrior), with a `Grade`; a character may have one.
22. **Quest**: Content with a `Name`/`Description` and ordered `Objective`s (unique by `Order`, kept sorted).
23. **Recipe**: Crafting that turns input `Ingredient`s (item + positive quantity, unique per item) into an `Output` item.

**Inventory** is a `Character` component (not a standalone entity): a list of `ItemStack`s (item + quantity) that merge by item.

**Grades** are free-form strings on the entity. `defaults.yml` defines a `rpg.grades`
ladder (Common → Mythic, each with a `level` and `power_multiplier`) and the catalog tags
items/abilities/skills with grade names, but nothing consumes either yet: the seeder does
not pass `Grade` through, and no power calculation reads the multiplier.

## CLI surface

| Command | Sub-commands |
|---|---|
| `character` | `create`, `add-power`, `list`, `clean`, `give-item`, `train-node`, `idle`, `pass-time` |
| `status` | *(flags only: `--name`, defaults to the main character)* |
| `species` | `create`, `list` |
| `realm` | `create`, `create-level`, `list`, `show` |
| `powersystem` | `create`, `list`, `show` |
| `universe` | `create`, `create-system`, `list`, `show` |
| `multiverse` | `create`, `create-universe`, `list`, `show` |
| `omniverse` | `create`, `create-multiverse`, `list`, `show` |
| `reality` | `create`, `create-omniverse`, `list`, `show` |
| `timeline` | `show`, `create-event` |
| `novel` | `create`, `create-volume`, `create-chapter`, `list`, `show` |
| `ability` / `skill` / `item` / `profession` / `class` | `create`, `list`, `show` |
| `effect` / `equipment` | `create`, `list`, `show` |
| `quest` | `create`, `create-objective`, `list`, `show` |
| `recipe` | `create`, `create-input`, `list`, `show` |

Every long flag gets an auto-derived single-letter shorthand (first free letter of the
flag name), and long flags must use `--`; `-foo` is rejected with a usage error.

Known gaps in the surface: `tge powersystem` advertises an `add-power` sub-command in its
help text that is **not** implemented (adding nodes/edges is only reachable through the
seeder), and there is no command to add a `Location` realm to a universe.
