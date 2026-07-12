# Power System: JoJo's Bizarre Adventure (Stands)

## Overview
A "Stand" is the visual manifestation of a person's life energy and soul, a concept pioneered by Hirohiko Araki. This system completely bypasses linear "power levels" in favor of highly specialized, lateral rulesets where intellect and creative application beat raw strength.

## Core Mechanics
1. **The Stand Manifestation**
   Stands are invisible to non-Stand users. They typically appear as humanoid figures, but can be objects, swarms, or phenomena. Any damage dealt to a Stand is reflected on the user's body.

2. **The "One Rule" Exceptions**
   Stands are defined by extreme hyper-specialization. A Stand usually does *one* specific thing incredibly well (e.g., unzipping any surface, accelerating time, turning touched objects into bombs). Combat is entirely puzzle-based: figuring out the enemy's rule and exploiting its blind spot.

3. **Stand Stats (Parameters)**
   To quantify these ethereal powers, Stands are graded (A to E) across six parameters:
   - **Destructive Power**: Raw physical strength or lethal capacity.
   - **Speed**: Agility and attack frequency.
   - **Range**: How far the Stand can travel from the user. (Note: Higher range usually results in lower Destructive Power).
   - **Persistence/Stamina**: How long the ability can be maintained.
   - **Precision**: Accuracy and motor control.
   - **Developmental Potential**: The capacity to learn new applications or evolve.

4. **Evolution (Requiem)**
   A Stand can evolve (e.g., pierced by a Stand Arrow) into a "Requiem" Stand. This completely rewrites its ruleset to grant a god-like ability specifically tailored to grant the user their most desperate, immediate desire at the exact moment of evolution.

## Simulation / Implementation Concepts for CLI
- **Parameter Grading**: Instead of a flat `Tier` integer, `PowerNode` definitions should utilize a 6-axis stat block (`DestructivePower: "A", Speed: "C", Range: "E"`). The simulation engine must translate these letter grades into underlying numeric multipliers.
- **Inverse Scaling Restrictions**: Enforce design logic during node creation—if a user defines a node with an "A" in Range, the engine could dynamically cap Destructive Power at "C" to simulate the inverse scaling of Stands.
- **Puzzle-Based Combat Metrics**: In CLI simulations, replace linear HP depletion with a `UtilityMatchup` check. A Stand with high `Precision` but low `Power` might instantly defeat a high `Power` brute if the combat context favors tactical traps over head-on brawling.
