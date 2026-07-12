# Power System: PRT Classification (Worm)

## Overview
In the web serial *Worm* by Wildbow, the Parahuman Response Team (PRT) uses a tactical classification system to categorize superpowers. Crucially, this system does *not* measure raw power or utility; it measures **threat level** to first responders and dictates tactical engagement protocols.

## Core Mechanics
1. **The 12 Categories**:
   - **Mover**: Mobility/speed.
   - **Shaker**: Area-of-effect/environmental control.
   - **Brute**: Enhanced strength/durability.
   - **Breaker**: Shifts into a different physical/mental state.
   - **Master**: Minion generation/mind control.
   - **Tinker**: Advanced, futuristic technology crafting.
   - **Blaster**: Long-ranged offensive attacks.
   - **Thinker**: Enhanced perception, data gathering, or intelligence.
   - **Striker**: Touch-based or melee abilities.
   - **Changer**: Form alteration.
   - **Trump**: Manipulation or granting of other powers.
   - **Stranger**: Stealth, infiltration, or social concealment.

2. **Numerical Threat Ratings (1 - 12+)**:
   - **0-2**: Normal trained humans can handle.
   - **3-4**: Requires trained parahumans.
   - **5-7**: Requires a full parahuman squad.
   - **8-9**: Requires highly specialized teams; evacuation prioritized.
   - **10+**: Uncontainable; focus entirely on survival and mass evacuation.

3. **Hybridization**: Capes often possess multiple ratings (e.g., a Master 5 / Shaker 3) to reflect multi-faceted threats.

## Simulation / Implementation Concepts for CLI
- **Vector-Based Power States**: Instead of a single `Tier` integer, a node could grant a vector of PRT ratings (e.g., `Mover: 2, Brute: 4`).
- **Tactical Matchups**: In the CLI combat simulation, certain categories could rock-paper-scissors others (e.g., a high Trump rating counters high ratings in other categories; Strangers bypass Brute defenses).
- **Threat Aggregation**: A character's total "PowerValue" could be calculated as the highest rating across all categories, or a weighted sum of their PRT vector.
