# Domain Model

> Slice/collection fields are shown with `*` multiplicity on the relationships
> rather than `[]` in the boxes (Mermaid class members don't take Go's `[]T`).

```mermaid
%%{init: {"theme":"dark","themeVariables":{"lineColor":"#c8d3f5"},"themeCSS":".edgePath path,.flowchart-link{stroke-width:2px} .messageLine0,.messageLine1{stroke-width:2px} .relation{stroke-width:2px} .actor{stroke-width:2px} .node rect,.node circle,.node polygon,.node path{stroke-width:2px} .cluster rect{stroke-width:2px}"}}%%
classDiagram
    class Character {
        +string Name
        +CharacterType Type
        +Gender Gender
        +string PowerValue
        +Mortal Mortal
        +Class Class
        +Profession Profession
        +Stats Stats
        +Inventory Inventory
    }
    class Species {
        +string Name
        +float64 Power
        +int Lifespan
        +Gender DefaultGender
    }
    class PowerSystem {
        +string Name
        +SystemKind Kind
    }
    class Power {
        +string Name
    }
    class PowerState {
        <<interface>>
        +Kind() SystemKind
    }
    class CultivationState {
        +int Points
        +float64 Progress
        +Power() float64
        +Lifespan() float64
    }
    class Realm {
        +string Name
        +float64 PowerMultiplier
        +float64 PowerAdder
        +int BottleneckPoints
        +int MaxLevels
        +int MainCharacterMaxLevels
    }
    class Level {
        +int Number
        +string Name
        +int BreakthroughPoints
        +int BottleneckPoints
    }

    Character "1" o-- "*" Species : Species (hybrid-capable)
    Character "1" o-- "*" PowerSystem : PowerSystems (defs) + Power (instances)
    PowerSystem "1" o-- "*" Power : Powers tree
    Power "1" o-- "*" Power : Children
    Power "1" --> "0..1" PowerState : State (instances only)
    PowerState <|.. CultivationState : implements
    CultivationState --> "1" Realm : Realm
    CultivationState --> "1" Level : Level
    Realm "1" o-- "*" Level : ordered sub-stages
```

The center of gravity is the **Character** and the **progression** vocabulary it
draws on. A `PowerSystem` is a named tree of `Power` nodes; a `Realm` is a
cultivation stage whose `Power`/`Lifespan` follow linear `ax + b` formulas and which
is subdivided into ordered `Level`s carrying the breakthrough/bottleneck thresholds.
The key polymorphism is `PowerState`: a `Power` node's optional per-character progress,
implemented today by `CultivationState` (anchored to a `Realm` + `Level`) and, later,
by other system kinds (e.g. Magic) without changing `Power` or `PowerSystem`. RPG
building blocks (`Class`, `Profession`, `Stats`, `Inventory`) and the cosmology and
novel contexts hang off this core — see the other diagrams for those. Note the two
distinct roles the same `PowerSystem` type plays on a character (definitions in
`PowerSystems`, instances in `Power`), detailed next in
[character-power.md](character-power.md).
