# Power System: The Matrix (Simulation Mechanics)

## Overview
In *The Matrix*, the physical world is irrelevant. The power system is entirely based on computational processing speed, code manipulation, and realizing that the constraints of reality (gravity, speed, durability) are just parameters in an operating system.

## Core Mechanics
1. **The Battery/Processor Truth**
   Humans are plugged into a massive simulation. Their physical bodies provide energy/processing power, while their minds are trapped in a digital construct bound by programmed physics.

2. **The Anomaly (The One)**
   "The One" (Neo) is a systemic anomaly—a mathematical remainder in the Matrix's balancing equations. Because of his unique, high-bandwidth connection to the system, he can "see" the raw code of the simulation rather than the rendered textures.

3. **Code Manipulation & Bullet Time**
   Because the rules (gravity, velocity) are just code, an Anomaly who believes the rules are fake can rewrite them on the fly. 
   - **Bullet Time**: This is not actually "moving fast"; it is the Anomaly processing the spatiotemporal data of the simulation faster than the server's default tick rate. They perceive the bullet's trajectory code before the physical event renders, allowing them to casually step around it.

## Simulation / Implementation Concepts for CLI
- **Admin Privileges / Meta-Gaming**: A `SimulationState` where a character achieves `Enlightenment`. Instead of checking stats against the enemy, the character checks stats against the *Environment*. 
- **Parameter Override**: In CLI combat, an "Anomaly" character has a unique action: `OverrideCode`. They can literally spend their turn to target an enemy's `Damage` variable and reset it to `0`, or change the `Gravity` context to `false`, treating combat as a debugging exercise rather than a physical brawl.
