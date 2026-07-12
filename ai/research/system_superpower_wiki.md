# Power System: Superpower Wiki (Fandom)

## Overview
The Superpower Wiki (powerlisting.fandom.com) acts as a massive crowdsourced encyclopedia of abilities. While it lacks a single canonical power system, its "Fanon" community has established widespread categorization methods to define, scale, and limit powers for worldbuilding.

## Core Mechanics
1. **Categorization by Function**: Abilities are grouped into archetypes rather than strict power levels. Common categories include:
   - **Elemental Manipulation**: Control over fire, water, earth, etc.
   - **Metaphysical/Conceptual**: Control over abstract concepts (e.g., Hierarchy Manipulation, Fate Manipulation).
   - **Physical/Enhancement**: Super strength, speed, durability.
   
2. **Threat/Power Tiers**: Fanon systems often utilize letter grades (F through SSS) or Greek letters (Alpha, Beta, Omega) to rank destructive capacity or reality-warping potential.
   - *Street Tier*: Can destroy a wall/building.
   - *Planetary Tier*: Can destroy a planet.
   - *Multiversal Tier*: Can alter or destroy multiple realities.

3. **Energy Constraints**: Powers often draw on universal resources such as Chi, Mana, or Cosmic Energy, which acts as a stamina pool.

4. **Hierarchy Manipulation**: A unique conceptual power frequently discussed, allowing a user to rewrite the rules, social rank, or existential laws of reality itself.

## Simulation / Implementation Concepts for CLI
- **Tag-based Categorization**: Assign tags to SuperPower nodes (e.g., `tags: ["elemental", "fire", "offensive"]`).
- **Destructive Capacity (DC) Tiers**: Add a `Scale` or `Tier` enum (Street, City, Continental, Planetary) to nodes to calculate collateral damage in simulations.
- **Resource Pools**: Implement a `Stamina` or `Mana` pool in the `Character` entity that drains when powers are invoked based on their Tier.
