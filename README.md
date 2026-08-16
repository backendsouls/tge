# tge — Cultivation Game CLI

`tge` is a command-line tool for authoring and managing the elements of a
**cultivation-style** game or web novel: characters, their power systems and
cultivation progress, the cosmology they live in, RPG building blocks, and the
novels that narrate them.

It is written in Go, persists to a local SQLite database, and is structured with
**Hexagonal Architecture** (Ports & Adapters) so the domain logic stays independent
of the CLI and the storage engine.

## Features

- **Characters** — create mortals (and name-only main characters that bootstrap a
  default world), with species, gender/role rules (Hero/Heroine), RPG class,
  profession, stats and inventory.
- **Progression** — power systems modeled as trees of powers; cultivation realms
  with linear `ax + b` power/lifespan formulas, subdivided into ordered levels with
  breakthrough thresholds; per-character **cultivation state**.
- **Cosmology** — a containment hierarchy of Reality → Omniverse → Multiverse →
  Universe → Realm, each owning a timeline of ordered events.
- **RPG toolkit** — stats, items, abilities, skills, effects, equipment, classes,
  professions, quests and recipes.
- **Novels** — stories with ordered volumes and chapters, each led by one main
  character.
- **Self-contained persistence** — SQLite with schema managed by embedded
  [goose](https://github.com/pressly/goose) migrations; a starter catalog and
  defaults are seeded on first run.

## Getting started

Requires Go 1.26+.

```bash
# Build
go build -o tge ./cmd/tge

# Create a main character from just a name (bootstraps the default world)
./tge character create --name "Lin Feng"

# Add cultivation and inspect the character
./tge character cultivate --name "Lin Feng" --realm "Spirit Gathering" --level 2
./tge status
```

Example `status` output:

```
Lin Feng  (Male, Human)
  power: 1
  Age 16/80
  Stats: STR 5, AGI 5, INT 5, VIT 5, DEX 5, WIS 5, CHA 5, LUK 5
  Systems: Spirit
  Power:
    - Cultivation:
      - Spirit:
        - Realm: Spirit Gathering
        - Level: Second Level
  Inventory: empty
```

Run `./tge <command>` with no arguments to see its sub-commands. Top-level commands
include: `character`, `status`, `species`, `realm`, `powersystem`, `universe`,
`multiverse`, `omniverse`, `reality`, `timeline`, `novel`, and the RPG entities
(`class`, `profession`, `item`, `ability`, `skill`, `effect`, `equipment`, `quest`,
`recipe`).

## Configuration

Both are optional:

- `TGE_DB` — path to the SQLite database file (default: `tge.db`). Migrations are
  applied automatically on open.
- `TGE_DEFAULTS` — path to a YAML file overriding the embedded defaults
  (`internal/config/defaults.yml`); only the keys you set are overridden.

## Project layout

```
cmd/tge/                 composition root (wires adapters to services) + catalog seeding
internal/
  core/
    domain/              pure business entities: character, progression, cosmology, rpg, novel
    port/                driving & driven port interfaces
    service/             use-case implementations
  adapter/
    cli/                 inbound CLI adapter
    sqlite/              outbound SQLite adapter + goose migrations
  config/                embedded defaults & starter catalog (YAML)
ai_docs/                 architecture, domain and design docs, plus Mermaid diagrams
test/functional/         end-to-end tests (testcontainers)
```

## Testing

```bash
go test ./...            # unit + service + adapter tests
```

The functional tests under `test/functional` build the `Dockerfile` image with
[testcontainers-go](https://golang.testcontainers.org/) and exercise the built
binary against a real SQLite file, so they require Docker.

## Documentation

Deeper design notes live in [`ai_docs/`](ai_docs/):

- [architecture.md](ai_docs/architecture.md) — hexagonal layering and migrations.
- [domain.md](ai_docs/domain.md) — the domain model and entities.
- [decisions.md](ai_docs/decisions.md) — key design decisions and trade-offs.
- [diagrams/](ai_docs/diagrams/) — Mermaid diagrams of the architecture, domain and flows.
