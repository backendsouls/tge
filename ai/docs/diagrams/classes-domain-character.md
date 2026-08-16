# Character Domain Class Diagram

This diagram shows the composition of a Character inside `internal/core/domain/character`.
It integrates with `rpg` (Class, Profession, Stats, Inventory) and `power`
(`MechanicState`), and references `powersystem` **by name/ID only** — a character never
embeds a power system.

```mermaid
classDiagram
    class Character {
        +String Name
        +CharacterType Type
        +Gender Gender
        +List~Species~ Species
        +List~String~ Systems
        +String PowerValue
        +MechanicState MechanicState
        +List~NodeProgress~ UnlockedNodes
        +Mortal Mortal
        +Class Class
        +Profession Profession
        +Stats Stats
        +Inventory Inventory
        +IdleState IdleState
        +int64 NovelTime
        +CalculateTotalPower() float64
        +AdvanceNode(String, String, float64) float64
        +CurrentEnergyPools(int64) Map~String, int~
    }

    class CharacterType {
        <<enumeration>>
        MainCharacter
        SideCharacter
        SupportCharacter
        Hero
        Heroine
    }

    class Gender {
        <<enumeration>>
        Male
        Female
    }

    class Mortal {
        +int Age
        +int Lifespan
    }

    class Species {
        +String Name
        +float64 Power
        +int Lifespan
        +Gender DefaultGender
    }

    class NodeProgress {
        +String System
        +String NodeID
        +int Level
        +float64 Progress
        +float64 BasePower
        +CalculateBreakthrough() float64
        +Advance(float64) float64
    }

    class IdleState {
        +List~IdleSlot~ Slots
        +int TotalSlots
    }

    class IdleSlot {
        +int64 StartTime
        +float64 Duration
        +String System
        +String Power
        +float64 Rate
    }

    class MechanicState {
        <<external: domain/power>>
    }
    class Stats {
        <<external: domain/rpg>>
    }
    class Inventory {
        <<external: domain/rpg>>
    }
    class Class {
        <<external: domain/rpg>>
    }
    class Profession {
        <<external: domain/rpg>>
    }

    Character --> CharacterType : typed as
    Character --> Gender : has
    Character *-- Mortal : mortal status
    Character "1" o-- "*" Species : hybrid-capable
    Character "1" o-- "*" NodeProgress : UnlockedNodes
    Character *-- IdleState : background activity
    IdleState "1" o-- "*" IdleSlot : Slots
    Character *-- MechanicState : cross-system state
    Character --> Stats : has
    Character --> Inventory : possesses
    Character --> Class : optional
    Character --> Profession : optional
```

`NewMortalCharacter` defaults age to `16` when unset, stats to `rpg.BaseStats()` (`0.65`
across the board), `MechanicState` to tier `0` / `BasePower 1.0`, and `IdleState.TotalSlots`
to `1`, then computes `PowerValue`. `CheckRole` enforces the cross-character rules: a
`Hero` requires a **female** main character to already exist, a `Heroine` a **male** one.

`CalculateTotalPower()` is
`(MechanicState.BasePower + Σ node.BasePower × node.Level) × Π species.Power`, floored at
`1.0` and rendered into `PowerValue`. `NodeProgress.CalculateBreakthrough()` is the gate
for the current level, `100 × Level²`; `Advance` pours points in, levels up when the gate
fills, and returns the unconsumed remainder.

The whole struct is serialized as one JSON document (`data/characters/<slug>.json`), so
every field here — including the multi-species list and the idle slots — round-trips.
