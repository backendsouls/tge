# Flow: Create a (name-only) Main Character

```mermaid
%%{init: {"theme":"dark","themeVariables":{"lineColor":"#c8d3f5"},"themeCSS":".edgePath path,.flowchart-link{stroke-width:2px} .messageLine0,.messageLine1{stroke-width:2px} .relation{stroke-width:2px} .actor{stroke-width:2px} .node rect,.node circle,.node polygon,.node path{stroke-width:2px} .cluster rect{stroke-width:2px}"}}%%
sequenceDiagram
    actor User
    participant CLI as cli.App
    participant CS as CharacterService
    participant DW as DefaultWorldService
    participant SP as SpeciesRepository (SQLite)
    participant Repo as CharacterRepository (file/JSON)

    User->>CLI: tge character create --name "Lin Feng"
    CLI->>CS: CreateCharacter(in{Name, Type=MainCharacter})
    note over CS: --type defaults to MainCharacter
    CS->>DW: EnsureDefaults(ctx)
    activate DW
    note over DW: idempotent — every step ignores "already exists"
    DW->>DW: save Human base species
    DW->>DW: create default PowerSystem (kind=Cultivation)
    DW->>DW: create Universe + SaveSystems + SaveRealms
    DW->>DW: create Multiverse / Omniverse / Reality
    DW->>DW: ensure a Timeline for all 5 locations
    DW-->>CS: DefaultWorld{Species, PowerSystem, …}
    deactivate DW
    CS->>Repo: MainCharacters(ctx)
    CS->>CS: CheckRole(type, mains)
    CS->>SP: FindByName(species)
    SP-->>CS: Species{Power, Lifespan, DefaultGender}
    CS->>CS: gender ← species default → global default
    loop until the name is free
        CS->>Repo: FindByName(candidate)
        CS->>CS: on hit, try "Lin Feng (1)", "(2)", …
    end
    CS->>CS: NewMortalCharacter(cfg) → CalculateTotalPower()
    CS->>Repo: Save(character) — one JSON aggregate
    Repo-->>CS: ok
    CS-->>CLI: Character
    CLI-->>User: [Ding! System binding complete. Welcome, Host Lin Feng!]<br/>created MainCharacter "Lin Feng" (Male, Human)
```

Creating a character from **just a name** is the marquee flow. `CreateCharacter`
recognises the (defaulted) `MainCharacter` type and asks `DefaultWorldService` to
provision the default cosmology — the Human base species, a default power system of kind
`Cultivation` (named by `world.power_system` in `defaults.yml`, shipped as `"Cultivation"`;
the in-code fallback is `"Mortal Path"`), and the full `Reality → Omniverse → Multiverse →
Universe → Realm` chain with a timeline for each — **idempotently**, so repeated creations
converge on one shared world.

`CreateCharacter` then fills the blanks: species → the default world's Human base, gender →
the species' `DefaultGender`, then the global default from `defaults.yml`. It resolves the
optional `--class` / `--profession` against their repositories (failing if they don't
exist), enforces the cross-character role rules (`CheckRole` — a `Hero` needs a female main
character, a `Heroine` a male one), and builds a fresh **mortal** via `NewMortalCharacter`
(age 16, `BaseStats()` at 0.65, `MechanicState` tier 0 / BasePower 1.0, one idle slot),
which computes `PowerValue` before returning.

Two behaviours worth noting:

- **Names never collide, they suffix.** Because the file adapter derives a filename from a
  slug of the name, `CreateCharacter` loops `FindByName` and appends `" (1)"`, `" (2)"`, …
  until the slug is free rather than returning "already exists".
- **The character is not joined to the default system.** `Systems` only contains whatever
  `--system` passed; the provisioned default power system is created but not attached.
  Use `tge character add-power` to join one.

The result is a mortal with no `UnlockedNodes` — it hasn't trained yet. Progression is the
next flow, [progression-flow.md](progression-flow.md).
