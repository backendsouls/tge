# Flow: Create a (name-only) Main Character

```mermaid
%%{init: {"theme":"dark","themeVariables":{"lineColor":"#c8d3f5"},"themeCSS":".edgePath path,.flowchart-link{stroke-width:2px} .messageLine0,.messageLine1{stroke-width:2px} .relation{stroke-width:2px} .actor{stroke-width:2px} .node rect,.node circle,.node polygon,.node path{stroke-width:2px} .cluster rect{stroke-width:2px}"}}%%
sequenceDiagram
    actor User
    participant CLI as cli.App
    participant CS as CharacterService
    participant DW as DefaultWorldService
    participant Repo as CharacterRepository (SQLite)

    User->>CLI: tge character create --name "Lin Feng"
    CLI->>CS: CreateCharacter(in{Name, Type=MainCharacter})
    note over CS: Type defaults to MainCharacter
    CS->>DW: EnsureDefaults(ctx)
    activate DW
    note over DW: idempotent — reuses existing world
    DW->>DW: save Human base species
    DW->>DW: create default PowerSystem ("Spirit")
    DW->>DW: create Universe + SaveSystems + SaveRealms
    DW->>DW: create Multiverse / Omniverse / Reality (+ timelines)
    DW-->>CS: DefaultWorld{Species, PowerSystem, …}
    deactivate DW
    CS->>CS: fill blanks (species=Human, systems=[Spirit], gender=Male)
    CS->>CS: CheckRole + NewMortalCharacter(cfg)
    CS->>Repo: Save(character)
    Repo-->>CS: ok
    CS-->>CLI: Character
    CLI-->>User: created MainCharacter "Lin Feng" (Male, Human)
```

Creating a character from **just a name** is the marquee flow. `CreateCharacter`
recognises the (defaulted) `MainCharacter` type and asks `DefaultWorldService` to
provision the default cosmology — Human base species, the default "Spirit" power
system, and the full `Reality → Omniverse → Multiverse → Universe → Realm` chain with
their timelines — **idempotently**, so repeated creations reuse one shared world.
`CreateCharacter` then fills any unset fields from that default world (species → Human,
systems → `[Spirit]`, gender → the species/global default), enforces the cross-character
role rules (`CheckRole`, e.g. a Hero needs a female main character), builds a fresh
**mortal** character via `NewMortalCharacter`, and persists it. The result is a mortal
with `Power` (instances) empty — it hasn't cultivated yet. Adding cultivation is the
next flow, [cultivate-flow.md](cultivate-flow.md).
