# RPG Domain Class Diagram

This diagram maps out the core RPG entities that exist inside `internal/core/domain/rpg`.
They represent standard definitions of objects, skills, and mechanics within the game
world. Every identity-bearing entity is keyed by `Name`; `Stats` and `Inventory` are
components of a `Character` rather than standalone entities.

```mermaid
classDiagram
    class Item {
        +String Name
        +String Description
        +String Grade
    }
    class Ability {
        +String Name
        +String Description
        +String Grade
    }
    class Skill {
        +String Name
        +String Description
        +String Grade
    }
    class Class {
        +String Name
        +String Description
        +String Grade
    }
    class Profession {
        +String Name
        +String Description
        +String Grade
    }
    class Stats {
        +float64 STR
        +float64 AGI
        +float64 INT
        +float64 VIT
        +float64 DEX
        +float64 WIS
        +float64 CHA
        +float64 LUK
        +Add(Stats) Stats
    }
    class EquipmentSlot {
        <<enumeration>>
        Weapon
        Armor
        Accessory
    }
    class Equipment {
        +String Name
        +EquipmentSlot Slot
        +Stats Bonus
    }
    class EffectKind {
        <<enumeration>>
        Buff
        Debuff
        Status
    }
    class Effect {
        +String Name
        +EffectKind Kind
        +String Description
    }
    class ItemStack {
        +String Item
        +int Quantity
    }
    class Inventory {
        +List~ItemStack~ Items
        +Add(String, int) error
    }
    class Objective {
        +int Order
        +String Description
    }
    class Quest {
        +String Name
        +String Description
        +List~Objective~ Objectives
        +AddObjective(int, String) error
    }
    class Ingredient {
        +String Item
        +int Quantity
    }
    class Recipe {
        +String Name
        +String Output
        +List~Ingredient~ Inputs
        +AddInput(String, int) error
    }

    Quest *-- Objective : contains (unique by Order, sorted)
    Recipe *-- Ingredient : requires (unique by Item)
    Inventory *-- ItemStack : holds (merged by Item)
    Equipment --> Stats : yields bonus
    Equipment --> EquipmentSlot : worn in
    Effect --> EffectKind : classified as
```

Stats are `float64`, not integers, and `BaseStats()` starts a character at `0.65` across all
eight attributes (`defaults.yml` sets the same `0.65` for `character.stats`, and that config
value is what `CharacterService` actually passes — `BaseStats()` is only the domain-level
fallback). Only `Add` exists (there is no `Sub`); it is used for base + equipment bonus.
`NewStats` rejects a negative attribute.

`Inventory` is a slice of `ItemStack`, not a map: `Add` merges into an existing stack for
the same item and rejects a blank item or a non-positive quantity. `Quest.AddObjective`
keeps objectives sorted by `Order` and rejects a duplicate order; `Recipe.AddInput`
rejects a duplicate ingredient item.

`Grade` is a free-form string on Item/Ability/Skill/Class/Profession. `defaults.yml`
defines a Common → Mythic ladder with `power_multiplier`s, but nothing consumes it yet and
the seeder does not pass catalog grades through — see
[../decisions.md](../decisions.md) §9.
