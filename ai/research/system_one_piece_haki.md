# Power System: One Piece (Devil Fruits & Haki)

## Overview
*One Piece* features a dual power system: **Devil Fruits** (supernatural abilities with severe drawbacks) and **Haki** (the physical manifestation of sheer willpower used to counter them). 

## Core Mechanics
1. **Devil Fruits**
   Consuming a Devil Fruit grants a unique power but permanently curses the user, making them sink like an anvil in seawater. 
   - **Paramecia**: Alters the user's body or environment (e.g., rubber body, creating strings).
   - **Zoan**: Allows transformation into an animal, mythical beast, or hybrid form, heavily boosting physical stats.
   - **Logia**: The user can create, control, and transform their entire body into a natural element (fire, light, magma), rendering them intangible to normal physical attacks.

2. **Devil Fruit Awakening**
   The pinnacle of fruit mastery. 
   - Awakened Paramecia can turn their surrounding environment into their element (e.g., turning the ground into rubber).
   - Awakened Zoans gain near-immortal regeneration and strength, but risk losing their mind to the "animal" spirit of the fruit.

3. **Haki (Willpower)**
   A dormant power in all living things, utilized in three forms to level the playing field against overpowered Devil Fruits:
   - **Observation Haki**: Precognition and spatial awareness.
   - **Armament Haki**: An invisible armor of willpower that can be used offensively to bypass Logia intangibility and strike their "true body".
   - **Conqueror's Haki**: A rare, innate ability to project one's dominating will to knock out weak-willed opponents. Advanced users coat their weapons in it for devastating strikes.

## Simulation / Implementation Concepts for CLI
- **Dual-System Rock-Paper-Scissors**: The simulation engine needs a `DefenseType` variable. A `Logia` node grants `Intangibility = true`. Standard physical attacks deal 0 damage unless the attacker simultaneously activates an `ArmamentHaki` node to override the `Intangibility` flag.
- **Environmental Instant-Kill**: If a character possesses a `DevilFruit` node, encountering `Context: Seawater` immediately applies a `Paralyzed` state and drains HP.
