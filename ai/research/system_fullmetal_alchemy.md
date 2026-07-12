# Power System: Fullmetal Alchemist (Alchemy)

## Overview
Alchemy in *Fullmetal Alchemist* is treated as a rigorous, metaphysical science rather than mystical magic. It is strictly governed by the law of Equivalent Exchange, meaning mass and elemental composition cannot be created from nothing; they can only be broken down and reconstructed.

## Core Mechanics
1. **The Three Stages of Transmutation**
   - **Comprehension**: Understanding the molecular/atomic structure of the target material.
   - **Deconstruction**: Breaking down the physical structure using energy (drawn from tectonic shifts or ley lines).
   - **Reconstruction**: Reassembling the material into a new shape.

2. **Equivalent Exchange**
   "To obtain, something of equal value must be lost." You cannot turn 1kg of lead into 2kg of lead, nor can you turn lead into gold (different atomic structures).

3. **Human Transmutation and The Truth**
   A taboo act. Attempting to transmute a human fails because a human soul lacks a quantifiable material equivalent. 
   - **The Rebound**: The alchemist is dragged to the "Gate of Truth," where they are forcibly shown universal knowledge.
   - **The Toll**: In exchange for this unearned knowledge, "Truth" (a god-like entity) takes a physical toll from the alchemist—a limb, organs, or their sight.

4. **Philosopher's Stones**
   A mythical red stone that allows an alchemist to bypass the law of Equivalent Exchange entirely. However, the grim reality is that the stone's fuel source is condensed human souls.

## Simulation / Implementation Concepts for CLI
- **Material Dependency (Inventory Check)**: The simulation CLI must tie a `Transmute` action directly to `Inventory`. If `Node: Steel Spear` is used, the system must verify `Inventory` contains `Iron >= 5kg`.
- **Rebound Affliction**: Attempting a high-tier skill without sufficient resources shouldn't just fail; it should trigger a `Rebound` penalty, permanently degrading a random character stat (simulating the toll of Truth).
