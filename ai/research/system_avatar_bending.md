# Power System: Avatar (Bending Arts)

## Overview
In the *Avatar: The Last Airbender* universe, "Bending" is the ability to manipulate the four classical elements (Water, Earth, Fire, Air) by projecting one's internal life energy (Chi) outward. 

## Core Mechanics
1. **Chi and Martial Arts**
   Bending is not telekinesis; it is strictly tied to physical martial arts movements. Each element is based on a real-world martial art (e.g., Tai Chi for water, Hung Gar for earth) to dictate how the Chi flows from the body to the element. 

2. **Sub-Bending Disciplines**
   True mastery of an element allows practitioners to bend variants or sub-elements:
   - **Water**: Bloodbending (manipulating water inside a living body), Healing, Spiritbending.
   - **Earth**: Metalbending (manipulating unrefined earth impurities inside metal), Lavabending, Seismic Sense.
   - **Fire**: Lightning Generation (separating positive and negative energies), Combustionbending.
   - **Air**: Flight (releasing all earthly tethers), Spiritual Projection.

3. **Energybending**
   The ancient, original form of bending. Instead of manipulating elements, the user manipulates the Chi pathways within another person's body directly. This is used to permanently strip away or restore a person's bending abilities.

## Simulation / Implementation Concepts for CLI
- **Physical Action Coupling**: A `BendingState` must require the character to not be physically restrained. A `StatusEffect: Paralyzed` or `Chi-Blocked` must entirely disable access to `PowerNodes`.
- **Environmental Buffs/Debuffs**: Elemental systems are heavily reliant on the arena. In a combat simulation, `Context: Ocean` should grant a 3x multiplier to Water nodes, while `Context: Solar Eclipse` should set Fire nodes to 0.
