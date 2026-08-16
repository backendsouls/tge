# Progression Domain Class Diagram

Progression is split across three packages under `internal/core/domain`:

- `powersystem` — the shared DAG definition (`PowerSystem`, `PowerNode`).
- `power` — the character-level `MechanicState` and the `PowerState` interface.
- `cultivation` — the realm/level model and `CultivationState`.

The **shipped** progression path is `character.NodeProgress` against a `PowerNode`.
`cultivation` and `superpower` are staged but not wired to any service — see
[../decisions.md](../decisions.md) §5.

```mermaid
classDiagram
    class PowerSystem {
        +String Name
        +PowerSystemType PowerSystemType
        +Map~String, PowerNode~ Nodes
        +AddNode(PowerNode) error
        +AddEdge(nodeID, targetID, EdgeType) error
        +Names() List~String~
    }

    class PowerNode {
        +String ID
        +String Name
        +String Category
        +List~String~ Tags
        +List~String~ Parents
        +List~String~ Siblings
        +List~String~ MutuallyExclusive
        +float64 BasePower
        +Map~String, float64~ StatVector
        +Map~String, int~ MaterialReq
        +List~String~ Drawbacks
        +AddMutuallyExclusive(String) error
    }

    class PowerSystemType {
        <<enumeration>>
        Cultivation
        Magic
        SuperPower
        Reiatsu
        Gamer
    }

    class EdgeType {
        <<enumeration>>
        parent
        sibling
        mutually_exclusive
    }

    class MechanicState {
        +int Tier
        +float64 BasePower
        +bool IsAwakened
        +float64 Alignment
        +Map~String, int~ EnergyPools
        +Map~int, int~ SpellSlots
        +List~String~ PermanentTraits
        +List~String~ Vows
        +AddEnergyPool(String, int)
        +SetAlignment(float64) error
    }

    class PowerState {
        <<interface>>
        +Kind() PowerSystemType
        +Power() float64
    }

    class Realm {
        +String Name
        +int Tier
        +float64 PowerMultiplier
        +float64 PowerAdder
        +float64 LifespanMultiplier
        +float64 LifespanAdder
        +int MaxLevels
        +int MainCharacterMaxLevels
        +List~Level~ Levels
        +AddLevel(int, String, int) error
        +MaxLevelsFor(bool) int
        +Power(float64) float64
        +Lifespan(float64) float64
    }

    class Level {
        +int Number
        +String Name
        +int BreakthroughPoints
    }

    class CultivationState {
        +Realm Realm
        +Level Level
        +int Points
        +float64 Progress
        +Ready() bool
        +AdvanceWithin(points) CultivationState
    }

    class SuperPowerState {
        +int Tier
    }

    class CultivationPath {
        <<interface>>
        +Name() String
    }

    PowerSystem "1" o-- "*" PowerNode : Nodes (flat map)
    PowerNode "*" --> "*" PowerNode : Parents / Siblings / MutuallyExclusive
    PowerSystem --> PowerSystemType : kind
    PowerSystem ..> EdgeType : AddEdge
    Realm "1" *-- "*" Level : ordered sub-stages
    PowerState <|.. CultivationState : implements
    PowerState <|.. SuperPowerState : implements
    CultivationState --> Realm : anchored to
    CultivationState --> Level : current
```

`PowerSystem.AddEdge` with `EdgeParent` runs a DFS up the parent chain and returns
`ErrCyclicDependency` if the edge would close a cycle; an `EdgeMutuallyExclusive` edge is
written onto **both** nodes. `PowerNode.ID` is a slug of the name (lowercase, spaces → `_`).

`Realm` uses linear `ax + b` formulas for both power and lifespan, evaluated at a
progress value `x`. `MaxLevelsFor(isMain)` returns the main character's higher cap when
one is set, and `AddLevel` rejects a level number above the realm's *effective* cap (the
higher of the two), so a realm never defines more levels than the most privileged
character could reach. `CultivationState.AdvanceWithin` fills the current level's
breakthrough gate, breaks through level by level, and hands back any points left over once
the realm's last level is full so the caller can carry them into the next realm.

`CultivationPath` (`BodyCultivation`, `SpiritCultivation`, `SoulCultivation`) is an open
interface rather than an enum so each path can accrue its own attributes and items
without forcing changes on what references it (OCP). It is currently unreferenced.
