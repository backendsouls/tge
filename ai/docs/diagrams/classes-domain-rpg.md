# RPG Domain Class Diagram

This diagram maps out the core RPG entities that exist inside `internal/core/domain/rpg`. They represent standard definitions of objects, skills, and mechanics within the game world.

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
        +int STR
        +int AGI
        +int INT
        +int VIT
        +int DEX
        +int WIS
        +int CHA
        +int LUK
        +Add(Stats) Stats
        +Sub(Stats) Stats
    }
    class Equipment {
        +String Name
        +EquipmentSlot Slot
        +Stats Bonus
    }
    class Effect {
        +String Name
        +EffectKind Kind
        +String Description
    }
    class Inventory {
        +Map~String, int~ items
        +AddItem(String, int)
        +RemoveItem(String, int)
        +HasItem(String, int) bool
    }
    class Objective {
        +int Order
        +String Description
    }
    class Quest {
        +String Name
        +String Description
        +List~Objective~ Objectives
        +AddObjective(int, String)
    }
    class Ingredient {
        +String Item
        +int Quantity
    }
    class Recipe {
        +String Name
        +String Output
        +List~Ingredient~ Inputs
        +AddInput(String, int)
    }

    Quest *-- Objective : contains
    Recipe *-- Ingredient : requires
    Equipment --> Stats : yields bonus
```
