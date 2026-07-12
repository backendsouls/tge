# Power System: SCP Foundation (Anomaly Classification System & Humes)

## Overview
The SCP Foundation is a collaborative fiction project centered around a clandestine organization that secures and studies anomalous entities. Unlike traditional progression fantasy, power in the SCP universe is not "cultivated" but rather "contained" or "measured." The mechanics of power are defined by the threat an anomaly poses to baseline reality, humanity, and secrecy.

## Core Mechanics
1. **Anomaly Classification System (ACS)**
   Rather than a simple "power level," anomalies are categorized using a three-pillar system that assesses entirely different axes of threat:
   - **Containment Class (Object Class)**: Measures the *difficulty of containment*. 
     - *Safe*: Easily and safely contained.
     - *Euclid*: Unpredictable or requires significant resources to contain.
     - *Keter*: Exceedingly difficult to contain consistently.
     - *(Note: A button that destroys the universe but is locked in a box is "Safe", while a cat that randomly teleports is "Keter".)*
   - **Disruption Class**: Measures the anomaly's potential to break the "Veil" of secrecy and disrupt global normalcy (ranging from *Dark* [low impact] to *Amida* [universal disruption]).
   - **Risk Class**: Measures the localized lethality and danger posed to a single individual interacting with the anomaly (ranging from *Notice* to *Critical*).

2. **Hume Levels (Reality Mechanics)**
   The Foundation quantifies the stability of reality using a measurement called "Humes." 
   - **Baseline Reality** is exactly 1 Hume.
   - **Reality Benders (Type Greens)** often have a personal Hume level much higher than 1, while the environment around them drops below 1. This means they are more "real" than their surroundings, allowing their will to override local physics.
   - **Scranton Reality Anchors (SRAs)** are technological devices used by the Foundation to artificially enforce a baseline 1 Hume field, stripping reality benders of their power.

## Common Tropes
- **The Veil**: The absolute necessity of keeping the anomalous hidden from the mundane world.
- **Cognitohazards & Memetics**: Information that is dangerous simply to perceive or know (e.g., a symbol that causes instant death, or an idea that spreads like a virus).
- **Thaumiel Class**: Anomalies that are highly dangerous but are actively weaponized by the Foundation to contain other, worse anomalies.

## Simulation / Implementation Concepts for CLI
- **Multi-Axis Threat Vectors**: Similar to the Worm PRT system, an SCP node shouldn't have a single `Tier`. It should have an `ACS` struct: `{Containment: Keter, Disruption: Vlam, Risk: Warning}`.
- **Hume Physics Engine**: In a combat simulation, characters could have an `InternalHume` and an `EnvironmentalHume` aura. If `Char A InternalHume > Char B EnvironmentalHume`, Char A can bypass standard defenses (simulating reality bending).
- **Containment vs. Destruction**: Combat victory conditions can shift from "HP = 0" to "Target Contained," requiring the deployment of specific countermeasures (like Reality Anchors) to neutralize a target's Hume advantage.
