# Power System: Jujutsu Kaisen (Cursed Energy)

## Overview
In Gege Akutami's *Jujutsu Kaisen*, the power system revolves around "Cursed Energy" (Juryoku)—spiritual power generated from human negative emotions. It is a highly lethal, physics-bending system governed by strict spatial, psychological, and karmic rules.

## Core Mechanics
1. **Cursed Energy & Techniques**
   Cursed Energy is the fuel; Cursed Techniques are the engine. Energy can be used broadly to enhance physical stats (reinforcement), but sorcerers channel it through their innate techniques (unique genetic abilities) to manifest specific phenomena (e.g., controlling shadows, manipulating space).

2. **Binding Vows (Keiyaku)**
   A foundational law of the universe. Vows are supernatural contracts made with oneself or others based on a high-risk, high-reward principle. 
   - **Self-Imposed**: By sacrificing something (e.g., revealing how your technique works to the enemy, or restricting your weapon's use to certain times), the user gains a massive multiplier to their technique's output.
   - **Breaking Vows**: Breaking a vow with oneself just removes the buff. Breaking a vow made with another person results in unpredictable, catastrophic karmic punishment.

3. **Domain Expansion**
   The pinnacle of jujutsu sorcery. A sorcerer externalizes their "innate domain" (mental landscape) into reality using a barrier. 
   - **Sure-Hit Effect**: Within the barrier, the user's technique is imbued into the space itself, guaranteeing that attacks cannot miss.
   - **Domain Clashes**: If two sorcerers activate their domains, the more "refined" domain overwrites the weaker one. 

4. **Heavenly Restriction**
   A type of binding vow forced upon a character at birth. A person might be born with absolutely zero cursed energy in exchange for god-like physical prowess and immunity to domain sure-hit effects.

## Simulation / Implementation Concepts for CLI
- **Vow Multipliers**: The `MechanicState` must support `BindingVows`—dynamic modifiers where a user can manually invoke a restriction in exchange for a temporary `PowerMultiplier`. 
- **Domain Override**: In combat simulations, `DomainExpansion` acts as a field effect. If `Char A DomainRefinement > Char B DomainRefinement`, Char A gains a 100% accuracy buff, while Char B's techniques are suppressed.
