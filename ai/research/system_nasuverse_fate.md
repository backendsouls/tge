# Power System: Nasuverse (Fate Series)

## Overview
The power system of the Nasuverse (Type-Moon) is built entirely on the concept of **Mystery**. Power is not just raw energy; it is intrinsically tied to age, concealment, and human belief. The older and more hidden a phenomenon is, the stronger it is. As modern science explains the world, it erodes Mystery, making modern Magecraft inherently weaker than that of the Age of Gods.

## Core Mechanics
1. **Magic Circuits**
   The "hardware" of a Magus. A pseudo-nervous system present at birth used to convert life force (Od) or atmospheric mana into usable magical energy. Because you cannot easily increase your circuit count, power is largely determined by genetics and centuries of selective breeding within Magus families.

2. **Magecraft vs. Magic**
   - **Magecraft**: The artificial reenactment of phenomena that *could* be achieved by modern science (e.g., creating fire). It acts like computer programming: circuits generate energy, which is sent to a pre-existing "Foundation" (system) engraved in the world to execute a spell.
   - **True Magic**: Miracles that are entirely impossible for modern humanity to replicate regardless of time or resources (e.g., time travel, materialization of the soul). 

3. **Heroic Spirits and Noble Phantasms**
   Servants are familiars made from the souls of mythological/historical heroes. 
   - **Noble Phantasms**: The crystallization of a hero's legend. A Noble Phantasm's power is dictated by the hero's fame and the sheer weight of their "Mystery." They range from weapons (Excalibur) to conceptual abilities (reversing causality to guarantee a spear pierces the heart before it is even thrown).

## Simulation / Implementation Concepts for CLI
- **Mystery Scaling**: In a combat simulation, `Age` or `MysteryLevel` acts as an absolute armor multiplier. A modern spell (`Mystery: Low`) physically cannot damage a mythological beast (`Mystery: High`) regardless of the raw damage stat.
- **Circuit Overload**: Pushing beyond a character's maximum `CircuitQuality` shouldn't just drain MP; it should cause severe physical backlash, draining HP to simulate burning out one's nervous system.
