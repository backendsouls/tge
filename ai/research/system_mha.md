# Power System: My Hero Academia (Quirks)

## Overview
In *My Hero Academia*, power is entirely biological. "Quirks" are genetic mutations present in 80% of the population, functioning exactly like physical muscles or organs.

## Core Mechanics
1. **The Biological Limit**
   Because Quirks are biological, they have physical drawbacks. Using an explosion quirk strains the wrists; using a fire quirk overheats the body; using a zero-gravity quirk causes nausea. Stamina and physical conditioning are the primary limiters of power.

2. **Quirk Awakenings**
   A spontaneous, micro-evolutionary event. When a user experiences extreme, life-threatening trauma or emotional distress, their Quirk can "awaken." This is not just a power boost; the fundamental rules of the Quirk change (e.g., a touch-based disintegration quirk evolving to spread disintegration through the ground without direct contact).

3. **Quirk Singularity Doomsday Theory**
   A macro-evolutionary concept. With each generation, bloodlines mix, causing Quirks to become increasingly complex and destructive. The theory posits an eventual "Singularity" where Quirks will become so potent that human bodies will biologically be unable to withstand them, leading to self-destruction.

## Simulation / Implementation Concepts for CLI
- **Recoil Damage**: Every `QuirkNode` in the tech tree must have a `RecoilStat` property (e.g., `StaminaDrain`, `SelfDamage`). Using an overpowered node without leveling up the underlying `PhysicalConstitution` stat will result in the character killing themselves.
- **Trauma-Triggered Awakening**: A progression hook: If a character survives a simulation with exactly `1 HP`, there is a % chance to trigger an `Awakening`, permanently unlocking a hidden, overpowered branch of their tech tree.
