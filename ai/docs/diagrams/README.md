# Diagrams

Visual companion to [../architecture.md](../architecture.md), [../domain.md](../domain.md)
and [../decisions.md](../decisions.md). Each file below covers one logical area of
the current implementation with a Mermaid diagram and a short explanation.

Diagrams render natively on GitHub and in most Markdown viewers.

## Overview

| Diagram | What it shows |
|---|---|
| [architecture.md](architecture.md) | Hexagonal layering: CLI → driving ports → services → driven ports → **both** the SQLite and flat-file adapters. |
| [domain-model.md](domain-model.md) | Core domain entities and how they relate (Character, MechanicState, PowerSystem DAG, Realm, …). |
| [character-power.md](character-power.md) | A character's power: shared system definition vs. the progress the character holds. |
| [cosmology.md](cosmology.md) | The world containment hierarchy (Reality → … → Realm) and timelines. |

## Class diagrams (per package)

| Diagram | Package |
|---|---|
| [classes-domain-character.md](classes-domain-character.md) | `domain/character` — Character, Species, NodeProgress, IdleState. |
| [classes-domain-progression.md](classes-domain-progression.md) | `domain/powersystem`, `domain/power`, `domain/cultivation` — the DAG, MechanicState, Realm/Level. |
| [classes-domain-cosmology.md](classes-domain-cosmology.md) | `domain/cosmology` — Reality → Universe, Location, Timeline. |
| [classes-domain-rpg.md](classes-domain-rpg.md) | `domain/rpg` — Item, Skill, Stats, Equipment, Quest, Recipe, … |
| [classes-core-ports.md](classes-core-ports.md) | `core/port` + `core/service` — the repeated service/repository shape, and a service spanning several ports. |
| [classes-adapters.md](classes-adapters.md) | `adapter/cli`, `adapter/sqlite`, `adapter/file`. |

## Flows

| Diagram | What it shows |
|---|---|
| [create-character-flow.md](create-character-flow.md) | Sequence: creating a name-only main character + default-world provisioning. |
| [progression-flow.md](progression-flow.md) | Sequence: `train-node` → `idle` → `pass-time` → `status`, and how idle gains commit. |
| [persistence-migrations.md](persistence-migrations.md) | How the SQLite DB is opened and brought up to date with embedded goose migrations. |

> These describe the state of the code as of the power-system-DAG / flat-file-aggregate
> work. When a model or flow changes, update the matching diagram here.
