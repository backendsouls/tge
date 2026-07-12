# Power System: Tower of God (Shinsu)

## Overview
In the webtoon *Tower of God*, the world is entirely permeated by **Shinsu** (Divine Water). Shinsu acts as the air, water, and fundamental magic energy of the Tower. As characters climb to higher floors, the density of Shinsu increases exponentially, crushing those without the innate talent or resistance to withstand it.

## Core Mechanics
1. **Shinsu Manipulation**
   Shinsu can be molded into any element or force (fire, light, pure kinetic energy). Control is quantified by three specific metrics:
   - **Baang**: The number of independent units or spheres of Shinsu a user can control simultaneously.
   - **Myun**: The size or area of effect of the Shinsu.
   - **Soo**: The density or concentration of power packed into a single Baang.

2. **Irregulars vs. Regulars**
   - **Regulars**: Inhabitants of the Tower who are bound by its rules. They must form contracts with the Administrator of each floor to be granted permission to use Shinsu.
   - **Irregulars**: Those who open the doors of the Tower from the outside. They are unbound by the Tower's contracts and can manipulate Shinsu with absolute, terrifying freedom, essentially acting as demigods.

3. **Team Positions**
   Combat is highly structured into RPG-like roles:
   - **Fisherman**: Front-line melee and heavy assault.
   - **Wave Controller**: The "mage", manipulating Shinsu from the rear for AoE or support.
   - **Spear Bearer**: Long-range snipers.
   - **Light Bearer**: Tactical support, using floating "Lighthouses" to gather data, calculate variables, and project shields.
   - **Scout**: Vanguard reconnaissance.

## Simulation / Implementation Concepts for CLI
- **Environmental Density Pressure**: If a character's `ShinsuResistance` stat is lower than the `FloorDensity` variable of the current arena, they are immediately paralyzed or crushed, instantly ending the simulation.
- **Multicast Scaling (Baangs)**: Progression isn't just increasing a flat damage number. Leveling up should increase a `MaxBaangs` integer. An attack's total damage is calculated as `Damage = Soo (Density) * ActiveBaangs`.
