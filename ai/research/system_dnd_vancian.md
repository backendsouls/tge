# Power System: Dungeons & Dragons (Vancian Magic)

## Overview
Named after author Jack Vance, the traditional D&D magic system focuses on strict resource management, preparation, and scarcity rather than an endless pool of "mana."

## Core Mechanics
1. **Prepare and Forget (Vancian Core)**
   In its strictest form, a wizard must study a spellbook to "memorize" a specific spell. Casting the spell physically burns the knowledge from their mind. If they prepare two Fireballs, they can cast exactly two Fireballs that day.

2. **Spell Slots**
   Modern iterations use "Spell Slots" as ammunition. A character has a fixed number of slots across different power levels (e.g., four 1st-level slots, two 3rd-level slots). They can use these slots to fuel any spell they know of that level or lower. When slots run out, the caster is powerless until they undergo a "Long Rest" (8 hours of sleep).

3. **Spell Components**
   Spells require physical triggers to activate:
   - **Verbal (V)**: Spoken incantations. Countered by silencing magic or gagging.
   - **Somatic (S)**: Precise hand gestures. Countered by binding the caster's hands or wielding heavy armor/shields without proficiency.
   - **Material (M)**: Physical reagents (e.g., a bat guano for a fireball, or a 500gp diamond for resurrection). Some are consumed upon casting.

## Simulation / Implementation Concepts for CLI
- **Hard Resource Scarcity**: Instead of `Mana = 100`, the `Character` state utilizes a `map[int]int` for `SpellSlots` (e.g., `[Level1: 4, Level2: 2]`). Attempting an action deducts exactly 1 from the corresponding key.
- **Component Constraint Checks**: Before a `Cast` action resolves, the simulation validates constraints: If `Context.IsSilenced == true` and `Node.RequiresVerbal == true`, the action instantly fails. If `Inventory` lacks the specific `MaterialCost`, it fails.
