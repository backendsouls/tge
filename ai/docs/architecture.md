# Hexagonal Architecture

The `tge` (Cultivation Game CLI) is structured using **Hexagonal Architecture** (also known as Ports and Adapters).

## Project Structure

- `cmd/tge`: The **Composition Root**. It connects concrete adapters to the application services, initializing dependencies (the SQLite repositories, the flat-file repositories) and starting the application (CLI adapter). On startup it loads the default values (`internal/config`) and seeds the starter catalog (`seed.go`).
- `internal/config`: Loads default values from YAML. A baseline `defaults.yml` is **embedded** in the binary (`go:embed`); setting `TGE_DEFAULTS` to a file path overrides it (decoded on top of the embedded defaults, so only changed keys need to be specified). It provides two things: *fill-in defaults* (default cosmology names, the Human base species, a new character's gender/age/stats, the system-log template — consumed by `DefaultWorldService`, `CharacterService` and `logger`) and a *starter Catalog* of entity instances seeded into the store on startup when absent (idempotent — existing entities are left untouched).
- `internal/logger`: Two output channels. `Dev` writes developer diagnostics to **stderr** with a `[DEV]` prefix; `System` writes the in-fiction "system message" to **stdout** through a configurable template (`system_log.template`, default `[Ding! %s]`). The CLI uses `System` for the novel-facing beats of `character create` and `character train-node`. Note both are package-level and write to the process's real streams via `log`/`fmt.Printf` — they bypass the `out`/`err` writers injected into `cli.App`, so a test capturing `App` output will not see them.
- `internal/core`: The **Core** of the hexagon, containing business logic and interfaces, independent of outside concerns.
  - `domain`: Core business entities organized into bounded contexts — `character`, `cosmology`, `novel`, `powersystem`, `power`, `cultivation`, `superpower` and `rpg`. It handles intrinsic business rules (cross-character gender roles, DAG cycle rejection, novel volume uniqueness, realm level caps) and depends on nothing outside itself. The `rpg` context (Ability, Skill, Item, Effect, Equipment, Profession, Class, Quest, Recipe, plus the Stats/Inventory components of a character) provides reusable RPG building blocks.
  - `port`: The interfaces for driving and driven adapters.
    - *Driving Ports*: Use-case interfaces like `CharacterService` (inward facing, implemented by `service`).
    - *Driven Ports*: Repository interfaces like `CharacterRepository` (outward facing, implemented by `adapter/sqlite` and `adapter/file`).
  - `service`: Implements driving ports and orchestrates domain logic using driven ports. A service may span bounded contexts: e.g. `DefaultWorldService` provisions the default cosmology (Reality → Omniverse → Multiverse → Universe → Realm) plus a power system and the Human base species, which `CharacterService` uses to back a name-only main character. `IdleService` is an *internal* collaborator — it is not a port; `CharacterService` owns one and exposes idle use cases through its own interface.
- `internal/adapter`: Concrete implementations of ports.
  - `cli`: Driving adapter. Parses CLI commands with `flag` and formats output (`text/tabwriter` for tables). It depends only on driving-port interfaces, so the whole CLI can be exercised with fake services and in-memory buffers.
  - `sqlite`: Driven adapter. SQLite repositories (via `modernc.org/sqlite`, a pure-Go driver — no cgo) for realms, species, novels, timelines and the whole cosmology and RPG catalogue. Schema is managed with embedded goose migrations.
  - `file`: Driven adapter. Flat-file **JSON aggregate** repositories for the two entities whose shape is a graph rather than a table: `Character` (`data/characters/<slug>.json`) and `PowerSystem` (`data/power_systems/<slug>.json`). Writes go to a `.tmp` file and are `rename`d into place so a save is atomic. `Clean` soft-deletes by renaming to `.deleted`.

### Which adapter stores what

| Entity | Adapter | Storage |
|---|---|---|
| Character | `file` | `data/characters/*.json` |
| PowerSystem (+ nodes, edges) | `file` | `data/power_systems/*.json` |
| Realm + Levels | `sqlite` | `realms`, `realm_levels` |
| Reality / Omniverse / Multiverse / Universe | `sqlite` | one table per tier + a link table |
| Timeline + Events | `sqlite` | `timelines`, `timeline_events` |
| Species | `sqlite` | `species` |
| Novel / Volume / Chapter | `sqlite` | `novels`, `volumes`, `chapters` |
| RPG catalogue (ability, skill, item, effect, equipment, profession, class, quest, recipe) | `sqlite` | one table each (+ `quest_objectives`, `recipe_inputs`) |

The split is deliberate: entities that are naturally rows stayed in SQL; the two
aggregates that are really object graphs (a character with its unlocked nodes, idle
slots and mechanic state; a power system with its whole DAG) moved to serialized JSON
so the repository interface collapses to `Save`/`FindByName`/`List` with no join
mapping. See [decisions.md](decisions.md) §3.

## Database Migrations

Schema is managed with **[goose](https://github.com/pressly/goose) v3**, used as an embedded library (not a CLI):

- Migrations live in `internal/adapter/sqlite/migrations/*.sql` as goose SQL files (`-- +goose Up` / `-- +goose Down`), embedded into the binary via `//go:embed` — so the binary is self-contained and ships no loose migration files.
- Migrations are **one per entity**, numbered in dependency order: `00001_realms.sql` … `00019_novels.sql` create the tables, then `00020_polymorphic_powers.sql` and `00021_add_grade.sql` evolve existing ones with `ALTER TABLE`. New changes go in the next numbered file.
- `open(dsn)` in `schema.go` applies pending migrations on every database open. goose is configured once (dialect `sqlite3`, embedded FS, and a **`NopLogger`** so migration output never pollutes CLI stdout) and `goose.Up` is idempotent, so the many repository constructors that each call `open` are safe. Applied versions are tracked in goose's `goose_db_version` table.
- This replaced an earlier hand-rolled bootstrap (a slice of `CREATE TABLE IF NOT EXISTS` statements). That approach could only *add* tables; goose is what enables evolving existing ones. The `Up` blocks deliberately keep `IF NOT EXISTS`, so a database created by the old bootstrap adopts migrations cleanly.
- Migrations currently run only implicitly on open; there is no `tge migrate` CLI command yet.
- **Stale tables**: `00002_power_systems.sql` (`power_systems`, `powers`), `00004_characters.sql` (`characters`, `character_stats`, `character_systems`, `character_items`, `character_cultivations`) and `00020`'s `character_superpowers` are still created and migrated, but nothing reads or writes them since characters and power systems moved to the `file` adapter. They are dead schema awaiting a drop migration.

## Testing

- **Unit tests**: Standard Go tests located alongside the code in `internal/...`.
- **Functional tests**: Located in `test/functional`, using `testcontainers-go` to run the `tge` binary inside a Docker container, verifying end-to-end behaviour across multiple executions. They require Docker.
  - `cli_test.go` — multi-character flow, the Hero/Heroine role rule, novels, universes, "character requires an existing system", and stats rendering.
  - `character_ft_test.go` — `TestFunctional_CharacterInventoryAndState` (inventory + state round-trip).
  - `cycle_detection_ft_test.go` — the DAG rejects a parent edge that would close a cycle.
  - `mechanic_state_ft_test.go` — `MechanicState` survives a save/load cycle.
  - `progression_ft_test.go` — `TestFunctional_ProgressionDAG`, node unlocking against the graph.
