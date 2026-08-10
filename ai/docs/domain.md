# Domain Models

The domain for `tge` centers around creating and managing elements of a Cultivation-style game or novel.

## Key Entities

1. **Character**: A being in the story, born as a "Mortal" with attributes like `Age` and `Lifespan`.
   - **Type**: `MainCharacter`, `SideCharacter`, `SupportCharacter`, `Hero`, `Heroine`.
   - **Gender**: `Male` or `Female` (a `Gender` enum).
   - **Species**: `Species []character.Species` — a character may be of one or more species (e.g. a hybrid gained "by any means"); a fresh character is created with a single one. The primary (first) species defines its base `PowerValue` and `Lifespan`. (Persistence currently stores only the primary species name; see [decisions.md](decisions.md).)
   - **Rules**: A `Hero` requires a `Female` `MainCharacter` to exist, and a `Heroine` requires a `Male` `MainCharacter`.
   - **Power, two facets** — see [decisions.md](decisions.md) for the definition-vs-instance rationale:
     - `PowerSystems []progression.PowerSystem` — the complete **tree definitions** the character has access to (its membership/catalogue). This is the renamed former `Systems` field.
     - `Power []progression.PowerSystem` — the character's own **progressed instances** of those systems: same `PowerSystem` type, but with per-node progress (`PowerState`) hung on the tree. Empty for a fresh mortal; grows as the character cultivates. (Not persisted yet.)
     - `PowerValue string` — the rendered numeric power value (computed later; the renamed former `Power` string).
   - **Name-only `MainCharacter`**: A `MainCharacter` can be created from just a name. It is born from the built-in **Human base** species (its initial-status template: baseline average `Power` 0.65, peak `Power` 1 (mortal), `Lifespan` 80), defaults to `Male`, and is attached to a default `PowerSystem` ("Mortal Path"). Creating one provisions the default cosmology in the background — `Reality` "The Box" → `Omniverse` "Origin Omniverse" → `Multiverse` "Origin Multiverse" → `Universe` "Origin Universe" → `Realm` "Mortal Realm" — idempotently, so repeated creations reuse the same single default world. Explicitly provided type/gender/species/systems are honoured.
   - **RPG attributes** (see the RPG context below): a character has an optional `Class` (`rpg.Class`) and `Profession` (`rpg.Profession`) — held as the entity values, validated against and resolved from existing entities when given — a `Stats` block (defaults to the base 5/5/5/… spread, persisted in a separate `character_stats` table), and an `Inventory` of items (added via `tge character give-item`, stacks merge by item).

2. **Species**: A biological classification defining base status values and a per-species default.
   - Sets the initial `Power` and `Lifespan` for Mortal characters of that species.
   - **DefaultGender** (optional, `Male`/`Female`): when a character of this species is created without a gender, it falls back to the species' default (then to the global default). The built-in **Human** base species defaults to `Male`. Several species ship in the seed catalog (Human, Demon, Spirit, Beastkin, Dragon, Fae), each with its own default gender; more can be added via `tge species add --default-gender ...` or the defaults YAML.

3. **PowerSystem**: A named tree (forest) of powers, and the unit shared by a character's two power facets (definition in `PowerSystems`, instance in `Power`).
   - Each system (e.g., "Universe A Cultivation") contains root `Power` elements, which can optionally have child sub-powers (e.g., "Body", "Soul", "Spirit").
   - **Kind** (`SystemKind`): the family the system belongs to — `Cultivation` (realm/level progression) or `Magic` (placeholder for a future, differently-progressing kind). `NewPowerSystem` defaults to `Cultivation`. The kind exists so a Magic system can progress by its own rules without forcing cultivation-specific shapes onto everything (see [decisions.md](decisions.md)). (Not persisted yet.)
   - **Power node state**: a `Power` node carries an optional `State PowerState` — `nil` in a system *definition*, set only in a character's *instance* (`Character.Power`). `PowerState` is an interface so each `SystemKind` brings its own progress type. `CultivationState` is the first implementation: the current `Realm`, the current `Level`, the `Points` accumulated toward `Level.BreakthroughPoints`, and the `Progress` (the `x` fed into the realm's power/lifespan formulas). It also exposes `Power()`/`Lifespan()`. (This replaced the standalone `Cultivation{Path, Realm, Progress}` struct — the "path" is now simply the tree node the state hangs under.)

4. **Reality** (also known as the **Box**): The outermost collection, grouping one or more `Omniverse`s together.
   - Member omniverses are referenced by name and must already exist.
   - An omniverse belongs to at most one reality.

5. **Omniverse**: A collection that groups one or more `Multiverse`s together.
   - Member multiverses are referenced by name and must already exist.
   - A multiverse belongs to at most one omniverse.

6. **Multiverse**: A collection that groups one or more `Universe`s together.

7. **Universe**: A collection grouping multiple `PowerSystem`s together and containing in-universe `Realm`s (locations).
   - Systems belong exclusively to one Universe.
   - Universes can be grouped under a Multiverse.
   - Realms are locations/bubbles within a Universe.

8. **Realm (Cultivation Stage)**: A single stage of cultivation defining how power and lifespan grow.
   - Uses linear equations `ax + b` (`Multiplier * x + Adder`) for calculating current `Power` and `Lifespan`.
   - A realm has a realm-wide `BottleneckPoints` barrier but **no** breakthrough of its own — breakthrough is a per-`Level` concept.
   - **Levels**: A realm is subdivided into ordered `Level`s (the `Cultivation → Realm → Level` hierarchy), e.g. the "First Level" of the "Qi Condensation" realm. Each `Level` has a positive `Number` (unique within the realm), a `Name`, and **both** progression concepts: `BreakthroughPoints` to advance to the next level and `BottleneckPoints` lost when a breakthrough fails. Levels are kept sorted by number. Managed via `tge realm add-level --breakthrough … --bottleneck …` / `tge realm show`.
   - **Per-tier level caps**: A realm caps how many levels a character may reach, and the **main character gets a higher cap** than a normal character — e.g. `MaxLevels` 9 for normal characters, `MainCharacterMaxLevels` 13 for the main character (`--max-levels 9 --max-levels-main 13`). `0` = unlimited; an unset main cap inherits the normal cap. `MaxLevelsFor(isMain)` returns the applicable cap. A realm may define levels up to the highest (main) cap; adding a level whose number exceeds that is rejected, so the realm holds at most that many levels.
   - **Default cultivation**: The default cultivation (the default power system) is **"Spirit"**, seeded with **9 realms** (Spirit Gathering → Spirit Sovereign), **each with 9 levels** (First … Ninth Level). In the defaults YAML a realm's `level_count` auto-generates that many ordinally-named levels instead of listing them.

9. **Novel**: A story containing `Volumes` and `Chapters`, led by a single `MainCharacter`.
   - A `MainCharacter` can only be the lead of one `Novel`.
   - `Volumes` and `Chapters` have strict ordering and uniqueness constraints.

10. **Timeline**: An ordered sequence of `Event`s (each an `Order` + `Description`) owned by a **location**.
    - Every location — `Realm`, `Universe`, `Multiverse`, `Omniverse` and `Reality` (the Box) — owns exactly one Timeline. A location is referenced by a kind + name (a realm is also scoped by its `Universe`, since realm names are unique only within a universe).
    - A Timeline is provisioned automatically with a derived default name (`"<location> Timeline"`) whenever a location is created — both for the default cosmology and for user-created locations.
    - `Event`s are unique by `Order` within a Timeline and are kept sorted ascending.

## RPG Context (`internal/core/domain/rpg`)

A separate bounded context of role-playing-game building blocks. Each identity-bearing entity is keyed by `Name` and has a full vertical slice (domain → port → service → SQLite → CLI `add`/`list`/`show`).

11. **Stats**: A value block of eight attributes — `STR`, `AGI`, `INT`, `VIT`, `DEX`, `WIS`, `CHA`, `LUK` — all non-negative. `BaseStats()` is the default starting spread; `Add` combines blocks (e.g. base + equipment bonus). Stats are a component of `Character`, stored in the `character_stats` table.
12. **Ability**: An innate power (`Name`, `Description`).
13. **Skill**: A learned, trainable ability (`Name`, `Description`).
14. **Item**: A thing that can be held in an inventory (`Name`, `Description`).
15. **Effect**: A modifier or condition with a `Kind` (`Buff`, `Debuff`, `Status`).
16. **Equipment**: An equippable item with a `Slot` (`Weapon`, `Armor`, `Accessory`) and a `Stats` bonus.
17. **Profession**: A vocation (e.g. Blacksmith); a character may have one.
18. **Class**: A combat archetype (e.g. Warrior); a character may have one.
19. **Quest**: Content with a `Name`/`Description` and ordered `Objective`s (unique by `Order`, kept sorted).
20. **Recipe**: Crafting that turns input `Ingredient`s (item + positive quantity, unique per item) into an `Output` item.

**Inventory** is a `Character` component (not a standalone entity): a list of `ItemStack`s (item + quantity) that merge by item.
