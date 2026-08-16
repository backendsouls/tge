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
        +int64 NovelTime
        +Mortal Mortal
        +Class Class
        +Profession Profession
        +Stats Stats
        +Inventory Inventory
        +CalculateTotalPower() float64
        +AdvanceNode(system, nodeID, points) float64
        +CurrentEnergyPools(now) Map
    }
    class Species {
        +string Name
        +float64 Power
        +int Lifespan
        +Gender DefaultGender
    }
    class MechanicState {
        +int Tier
        +float64 BasePower
        +bool IsAwakened
        +float64 Alignment
        +Map~string,int~ EnergyPools
        +Map~int,int~ SpellSlots
        +AddEnergyPool(name, max)
        +SetAlignment(v) error
    }
    class NodeProgress {
        +string System
        +string NodeID
        +int Level
        +float64 Progress
        +float64 BasePower
        +CalculateBreakthrough() float64
        +Advance(points) float64
    }
    class IdleState {
        +int TotalSlots
    }
    class IdleSlot {
        +int64 StartTime
        +float64 Duration
        +string System
        +string Power
        +float64 Rate
    }
    class PowerSystem {
        +string Name
        +PowerSystemType PowerSystemType
        +Map~string,PowerNode~ Nodes
        +AddNode(PowerNode) error
        +AddEdge(node, target, EdgeType) error
    }
    class PowerNode {
        +string ID
        +string Name
        +string Category
        +float64 BasePower
        +Map~string,float64~ StatVector
        +Map~string,int~ MaterialReq
    }
    class Realm {
        +string Name
        +int Tier
        +float64 PowerMultiplier
        +float64 PowerAdder
        +int MaxLevels
        +int MainCharacterMaxLevels
        +MaxLevelsFor(isMain) int
        +Power(x) float64
    }
    class Level {
        +int Number
        +string Name
        +int BreakthroughPoints
    }
    class PowerState {
        <<interface>>
        +Kind() PowerSystemType
        +Power() float64
    }
    class CultivationState {
        +int Points
        +float64 Progress
        +Ready() bool
        +AdvanceWithin(points) CultivationState
    }
    class SuperPowerState {
        +int Tier
    }

    Character "1" o-- "*" Species : Species (hybrid-capable)
    Character "1" *-- "1" MechanicState : MechanicState
    Character "1" o-- "*" NodeProgress : UnlockedNodes
    Character "1" *-- "1" IdleState : IdleState
    IdleState "1" o-- "*" IdleSlot : Slots
    Character ..> PowerSystem : Systems (by name)
    NodeProgress ..> PowerNode : NodeID (by ID)
    PowerSystem "1" o-- "*" PowerNode : Nodes (DAG)
    PowerNode "*" --> "*" PowerNode : Parents / Siblings / MutuallyExclusive
    Realm "1" o-- "*" Level : ordered sub-stages
    PowerState <|.. CultivationState : implements
    PowerState <|.. SuperPowerState : implements
    CultivationState --> "1" Realm : Realm
    CultivationState --> "1" Level : Level
```

The center of gravity is the **Character**. Its power is held in three parts: a
character-level `MechanicState` (tier, awakening, alignment, energy pools, spell slots —
the facts that are true across *every* system it participates in), a list of
`NodeProgress` records naming individual nodes it has unlocked, and `Systems []string`
— membership by name only. The power systems themselves are **shared definitions** that
carry no per-character state: a `PowerSystem` is a **DAG** of `PowerNode`s held in a flat
ID-keyed map, where shape comes from edges on the nodes (`Parents`, `Siblings`,
`MutuallyExclusive`) and a parent edge is rejected if it would close a cycle. Multi-parent
and mutual exclusion are exactly what a tree could not express.

`PowerValue` is a pure derivation — `(MechanicState.BasePower + Σ node.BasePower × Level)
× Π species.Power` — recomputed on every mutation, and **every** species multiplies in,
which is what makes hybrids compound. `IdleState` drives background progression against
`NovelTime`, the character's own in-story clock (see
[progression-flow.md](progression-flow.md)).

`Realm`/`Level` and the `PowerState` implementations (`CultivationState`,
`SuperPowerState`) are the **staged** progression model: authored by `tge realm`, but not
wired to any service — the shipped path is `NodeProgress`. `cultivation` is unit-tested;
`superpower` has no tests and no references at all. See
[../decisions.md](../decisions.md) §1–§2 and §5.
