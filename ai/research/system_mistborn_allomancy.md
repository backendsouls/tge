# Power System: Allomancy (Mistborn)

## Overview
Allomancy is a "hard magic" system created by Brandon Sanderson in the *Mistborn* series. It is highly scientific, relying on the ingestion and "burning" of specific metals to yield precise, predictable physical or mental effects. 

## Core Mechanics
1. **The 16 Metals (Tech Tree)**: 
   The system is organized into a rigid periodic-table-like structure. There are four quadrants (Physical, Mental, Temporal, Enhancement), each containing two base metals and two alloys.
   - **Push vs. Pull**: Metals are paired. One pushes (e.g., Steel pushes on external metals), the other pulls (e.g., Iron pulls on external metals).
   - **Internal vs. External**: Metals affect either the user's own body/mind (Internal) or the outside world (External).

2. **Mistings vs. Mistborn**:
   - **Mistings**: Born with the genetic ability to burn exactly *one* metal (specialists).
   - **Mistborn**: Born with the genetic ability to burn *all* metals (generalists). 
   - *Rule*: You can never be born with the ability to burn exactly two or three metals. It is always one or all.

3. **Savantism**:
   If an Allomancer constantly burns their metal for an extended period, their body adapts, making them a "Savant." This dramatically increases their power but introduces severe physical/mental dependencies and drawbacks (e.g., a Tin savant has super-senses but can be blinded/deafened easily by bright lights or loud noises).

4. **God Metals**: Exceptionally rare metals (like Atium or Lerasium) that break the standard rules, offering reality-bending powers like precognition or rewriting spiritual DNA.

## Simulation / Implementation Concepts for CLI
- **Symmetric Graph Architecture**: Construct the `PowerSystem` tree as a highly symmetric, categorized grid (4x4) where nodes are tightly coupled as Push/Pull pairs.
- **Savant Progression State**: Add a secondary progression axis. Instead of just Tier 1-5, a user can attain `Savant` status on a node, granting a 3x multiplier but applying a persistent negative trait to the character.
- **Binary Access Control**: In the CLI, character creation rules could enforce the "One or All" mechanic for granting access to nodes within a specific system.
