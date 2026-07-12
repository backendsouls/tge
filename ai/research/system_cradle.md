# Power System: Cradle (Sacred Arts)

## Overview
Will Wight's *Cradle* is the definitive Western execution of the "Progression Fantasy" genre. Set on a planet with immense gravity, the system revolves around "Sacred Arts"—absorbing and refining environmental vital aura into personal "Madra" to evolve the physical and spiritual body.

## Core Mechanics
1. **Aura and Madra**
   The world generates Vital Aura (fire, wind, sword, dream). Sacred Artists cycle this aura through their spirit, converting it into Madra (personal mana). 

2. **Paths and Techniques**
   A "Path" is a specialized martial art utilizing specific aspects of Madra. Most paths consists of four technique types:
   - *Striker*: Ranged attacks.
   - *Enforcer*: Internal reinforcement for speed/strength.
   - *Ruler*: Controlling external ambient aura.
   - *Forger*: Solidifying Madra into physical constructs.

3. **The Iron Body**
   The Foundation, Copper, and Iron stages are preparatory. Advancing to "Iron" requires reforging the physical body. A "Perfect Iron Body" is achieved by undergoing extreme, torturous rituals (e.g., drinking sacred beast venom to gain the Bloodforged Iron Body for hyper-regeneration). This permanent physical passive dictates how the artist fights forever.

4. **The Lord Realms & Soulfire**
   Progression shifts from pure energy accumulation to philosophical enlightenment. Advancing to Underlord, Overlord, and Archlord requires the artist to weave "Soulfire" (distilled aura of the world) and discover deep, personal revelations about their core motivations (e.g., "Why do I practice the Sacred Arts?").

5. **Sages, Heralds, and Monarchs**
   The pinnacle tiers. 
   - *Sages* connect to universal "Icons" (conceptual authorities like the Sword Icon) to bend reality.
   - *Heralds* physically merge with their own spiritual Remnant to achieve absolute physical supremacy.
   - *Monarchs* achieve both simultaneously.

## Simulation / Implementation Concepts for CLI
- **Path Specialization**: The `PowerSystem` tree should enforce a "Technique Type" constraint (Striker, Enforcer, Ruler, Forger).
- **Permanent Passives (Iron Body)**: The `MechanicState` must support `IronBody` traits—permanent passive buffs applied at a specific tier threshold (e.g., Tier 3) that provide flat multipliers to regeneration or defense regardless of current energy pools.
- **Revelation Gates**: Breaking through higher bottlenecks (Lord realms) shouldn't just cost `Points`; it should require fulfilling a narrative or contextual "Revelation" condition in the simulation.
