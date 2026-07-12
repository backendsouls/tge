# Power System: Xianxia / Cultivation

## Overview
Xianxia ("Immortal Hero") is a genre of Chinese fantasy heavily influenced by Daoist philosophy and traditional medicine. Its progression systems revolve around "Cultivation"—the practice of absorbing environmental spiritual energy (Qi) to transcend mortal limits and achieve godhood/immortality.

## Core Mechanics
1. **The Realm Hierarchy**: Cultivation follows a strict, almost universal ladder of "Realms" (Major tiers) and "Stages/Levels" (Minor tiers).
   - *Qi Condensation (Gathering)*: Sensing and absorbing ambient energy.
   - *Foundation Establishment*: Compressing energy to form a stable base (dantian).
   - *Core Formation (Golden Core)*: Crystallizing energy into a dense, semi-permanent core.
   - *Nascent Soul*: Developing a spiritual avatar capable of surviving body death.
   - *Ascension / Godhood*: Shedding the mortal coil to ascend to higher dimensional planes.

2. **Bottlenecks and Breakthroughs**: Transitioning between levels is never linear. Cultivators hit "bottlenecks" where accumulating more Qi is useless without a philosophical epiphany, rare pill, or extreme physical stress.
   - **Heavenly Tribulations**: Breaking through major realms often triggers lightning strikes or karmic tests sent by the heavens to destroy those defying natural limits.

3. **Realm Suppression**: Power scaling is exponentially steep. A cultivator at the Foundation Establishment realm can often effortlessly crush dozens of Qi Condensation cultivators simply by releasing their aura ("suppression").

4. **Dao Comprehension**: At higher levels, raw energy takes a backseat to understanding fundamental universal laws ("The Dao" of Fire, Space, Time, Sword).

## Common Tropes
- **The Golden Finger**: A unique cheat (ancient artifact, spirit mentor, special bloodline) that lets the protagonist cultivate faster than peers.
- **Sect Politics**: Characters join massive martial sects to compete for scarce resources (spirit stones, pills).
- **Face-Slapping**: Humiliating arrogant antagonists who look down on the protagonist's initially low cultivation base.

## Simulation / Implementation Concepts for CLI
- **Non-Linear XP (Bottlenecks)**: Implement a mechanic where `Points` cap out at a bottleneck, requiring a separate `Breakthrough` action with a `Failure Rate` (Tribulation).
- **Aura Suppression Status**: In combat simulations, if `Char A Realm > Char B Realm`, automatically apply a flat percentage debuff to Char B's stats before calculations begin.
