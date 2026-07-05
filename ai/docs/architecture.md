# Hexagonal Architecture

The `tge` (Cultivation Game CLI) is structured using **Hexagonal Architecture** (also known as Ports and Adapters).

## Project Structure

- `cmd/tge`: The **Composition Root**. It connects concrete adapters to the application services, initializing dependencies (e.g. SQLite DB and repositories) and starting the application (CLI adapter). On startup it loads the default values (`internal/config`) and seeds the starter catalog (`seed.go`).
- `internal/config`: Loads default values from YAML. A baseline `defaults.yml` is **embedded** in the binary (`go:embed`); setting `TGE_DEFAULTS` to a file path overrides it (decoded on top of the embedded defaults, so only changed keys need to be specified). It provides two things: *fill-in defaults* (default cosmology names, the Human base species, a new character's gender/age/stats — consumed by `DefaultWorldService` and `CharacterService`) and a *starter Catalog* of entity instances seeded into the store on startup when absent (idempotent — existing entities are left untouched).
- `internal/core`: The **Core** of the hexagon, containing business logic and interfaces, independent of outside concerns.
  - `domain`: Contains core business entities organized into bounded contexts (`cosmology`, `progression`, `character`, `novel`, `rpg`). It handles intrinsic business rules (e.g., enforcing cross-character gender roles, novel volume uniqueness). It depends on nothing outside itself. The `rpg` context (Ability, Skill, Item, Effect, Equipment, Profession, Class, Quest, Recipe, plus the Stats/Inventory components of a character) provides reusable RPG building blocks.
  - `port`: The interfaces for driving and driven adapters.
    - *Driving Ports*: Interfaces like `CharacterService` use-cases (inward facing, implemented by `service`).
    - *Driven Ports*: Repository interfaces like `CharacterRepository` (outward facing, implemented by `adapter/sqlite`).
  - `service`: Implements driving ports and orchestrates domain logic using driven ports. A service may span bounded contexts: e.g. `DefaultWorldService` provisions the default cosmology (Reality → Omniverse → Multiverse → Universe → Realm) plus a power system and the Human base species, which `CharacterService` uses to back a name-only main character.
- `internal/adapter`: Concrete implementations of ports.
  - `cli`: Driving adapter. The user interface logic parsing CLI commands and formatting output.
  - `sqlite`: Driven adapter. The SQLite repository implementations implementing `port` interfaces for persistence.
  - `memory`: Driven adapter. (Likely used for testing or in-memory runs).

## Database Migrations

Schema is managed with **[goose](https://github.com/pressly/goose) v3**, used as an embedded library (not a CLI):

- Migrations live in `internal/adapter/sqlite/migrations/*.sql` as goose SQL files (`-- +goose Up` / `-- +goose Down`), embedded into the binary via `//go:embed` — so the binary is self-contained and ships no loose migration files. The current schema is consolidated into `00001_init.sql`; subsequent changes (including column evolution via `ALTER TABLE`) go in new numbered files (`00002_*.sql`, …).
- `open(dsn)` in `schema.go` applies pending migrations on every database open. goose is configured once (dialect `sqlite3`, embedded FS, and a **`NopLogger`** so migration output never pollutes CLI stdout) and `goose.Up` is idempotent, so the many repository constructors that each call `open` are safe. Applied versions are tracked in goose's `goose_db_version` table.
- This replaced an earlier hand-rolled bootstrap (a slice of `CREATE TABLE IF NOT EXISTS` statements). That approach could only *add* tables; goose is what enables evolving existing ones. The `00001_init.sql` `Up` deliberately keeps `IF NOT EXISTS`, so a database created by the old bootstrap adopts migrations cleanly (the creates are no-ops and goose just records version 1).
- Migrations currently run only implicitly on open; there is no `tge migrate` CLI command yet.

## Testing

- **Unit tests**: Standard Go tests located in `internal/...` directories.
- **Functional tests**: Located in `test/functional`, using `testcontainers-go` to run the `tge` binary inside a Docker container, verifying end-to-end functionality including SQLite persistence across multiple executions.
