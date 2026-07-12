# Power System: Dragon Ball (Ki & Transformations)

## Overview
*Dragon Ball* established the gold standard for linear, hyper-scaling power systems in shonen anime. It relies on **Ki** (life force) and genetic evolutionary traits to continuously push characters past their absolute limits.

## Core Mechanics
1. **Ki (Life Force)**
   Ki is the physical and spiritual energy within all living beings. It can be utilized to:
   - Multiply raw physical strength, speed, and durability.
   - Project destructive energy blasts.
   - Achieve flight and sense the presence/power level of others across planets.

2. **Linear Power Scaling (Power Levels)**
   Unlike highly conceptual systems (like JoJo or Fate), Dragon Ball combat is largely dictated by pure math. If Character A's total Ki output is significantly higher than Character B's, Character A is virtually immune to Character B's attacks.

3. **Transformations (Multipliers)**
   The primary method of progression. Saiyans can unlock transformations (Super Saiyan 1, 2, 3, God, Blue) that act as massive, flat multipliers to their base power level (e.g., SSJ1 is a 50x multiplier). These forms rapidly drain stamina, creating a tactical timer on their usage.

4. **The Zenkai Boost**
   A biological cheat code exclusive to the Saiyan race. If a Saiyan is severely injured and brought to the brink of death, surviving and recovering will result in a permanent, massive evolutionary boost to their base power level.

## Simulation / Implementation Concepts for CLI
- **Stat Multiplier Toggles**: Instead of permanent stats, a character possesses a `BasePower` and an array of `Transformations`. Activating a transformation applies a massive multiplier (e.g., `* 50`) but initiates a heavy `StaminaDrain` per turn.
- **Zenkai Progression Hook**: In the progression engine, if a simulation ends with the character at `< 5% HP` but victorious (or surviving), applying a `Zenkai` multiplier to their permanent `BasePower` upon returning to the menu.
- **Absolute Defense Check**: Combat math should enforce linear supremacy: If `Attacker.Ki < (Defender.Ki * 0.2)`, `Damage = 0`.
