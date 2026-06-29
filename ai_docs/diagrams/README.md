# Diagrams

Visual companion to [../architecture.md](../architecture.md), [../domain.md](../domain.md)
and [../decisions.md](../decisions.md). Each file below covers one logical area of
the current implementation with a Mermaid diagram and a short explanation.

Diagrams render natively on GitHub and in most Markdown viewers.

| Diagram | What it shows |
|---|---|
| [architecture.md](architecture.md) | Hexagonal layering: CLI → driving ports → services → driven ports → SQLite. |
| [domain-model.md](domain-model.md) | Core domain entities and how they relate (Character, PowerSystem, Realm, …). |
| [character-power.md](character-power.md) | A character's power: definition-vs-instance split and the `PowerState` polymorphism. |
| [cosmology.md](cosmology.md) | The world containment hierarchy (Reality → … → Realm) and timelines. |
| [create-character-flow.md](create-character-flow.md) | Sequence: creating a name-only main character + default-world provisioning. |
| [cultivate-flow.md](cultivate-flow.md) | Sequence: `character cultivate` writing cultivation state, and `status` rendering it. |
| [persistence-migrations.md](persistence-migrations.md) | How a DB is opened and brought up to date with embedded goose migrations. |

> These describe the state of the code as of the cultivation-state / goose-migrations
> work. When a model or flow changes, update the matching diagram here.
