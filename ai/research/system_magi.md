# Power System: Magi (Rukh & Magoi)

## Overview
The power system in *Magi: The Labyrinth of Magic* is heavily tied to destiny and the spiritual flow of the world. It merges elemental magic with the summoning of ancient Djinn spirits from dungeons.

## Core Mechanics
1. **Rukh and Magoi**
   - **Rukh**: Spiritual entities (white birds) that represent the flow of fate. They are the source of all souls and magic.
   - **Magoi**: The raw energy produced by Rukh. Magicians use Magoi to issue commands to the Rukh to create elemental magic. Over-exhausting Magoi leads to severe physical hemorrhage and death.
   - **Falling into Depravity**: If a person curses their fate or succumbs to absolute despair, their Rukh turns black. They lose connection to natural magic but gain access to destructive "Black Magoi."

2. **Djinn and Metal Vessels**
   Djinn are massive magical beings that reside in Dungeons. When a warrior conquers a Dungeon, the Djinn inhabits a personal item of the conqueror (a sword, jewelry), turning it into a "Metal Vessel."

3. **Djinn Equip (Masou)**
   The ultimate combat form. The user completely merges with the Djinn inside their Metal Vessel, taking on a mythological appearance. This grants them absolute control over a specific element and allows them to cast "Extreme Magic"—a devastating, ultimate attack that consumes massive amounts of Magoi.

## Simulation / Implementation Concepts for CLI
- **Alignment Penalty (Depravity)**: If a character's `Morality` or `Hope` stat drops below 0, their `MagoiAlignment` switches to `Black`. This locks out healing nodes but doubles the damage of destructive nodes.
- **Weapon-Bound Equipment**: `DjinnEquip` is tied to the `Inventory`. The simulation checks if `Equipment.Slot1` contains a `MetalVessel`. Activating the equip applies a temporary multiplier to stats but applies a `MagoiDrain` debuff every turn.
