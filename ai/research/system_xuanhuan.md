# Power System: Xuanhuan (Mysterious Fantasy)

## Overview
Xuanhuan is a highly fluid, modern Chinese web novel genre that translates to "Mysterious Fantasy." While it heavily utilizes the progression grind of *Xianxia* (cultivating energy to reach higher realms), it strips away the strict Daoist religious mythology. This makes it a "wild west" genre where authors freely mix Eastern cultivation with Western magic, sci-fi technology, and video game mechanics.

## Core Mechanics
1. **Hybrid Energy Systems**
   While characters may still cultivate *Qi*, Xuanhuan worlds often introduce Elemental Magic (Mana), Aura, Spiritual Power, or Divine Essences. It is common for a protagonist to dual-cultivate both a Western-style magic core and an Eastern-style physical Dantian.

2. **The Progression Grind (Realms)**
   Similar to Xianxia, progression is quantified in strict Realms (e.g., Bronze, Silver, Gold, or Spirit Apprentice, Spirit Master, Spirit Grandmaster). The mechanical gap between these realms is absolute—a lower realm cannot defeat a higher realm without an absurd external advantage.

3. **External Aids and Professions**
   Progression relies heavily on the economy of the world. Cultivators consume massive amounts of alchemical pills, spirit beast cores, and medicinal baths. Sub-professions like Alchemists, Forgers, and Array Masters often hold higher social status than pure fighters due to their role in the progression economy.

## Common Tropes
- **The Cheat / The System**: Xuanhuan is famous for protagonists who receive an unfair advantage early on. This might be a literal LitRPG "System" in their mind, an artifact that accelerates time for training, or an ancient god's soul granting them lost knowledge.
- **Multiversal Scope (Ascension)**: The narrative scope often escalates infinitely. Once the protagonist becomes the strongest in their world, they "ascend" to a higher dimensional realm where they are once again the weakest, starting the grind over.
- **Talent vs. Effort**: While Western fantasy often relies on chosen ones or innate fate, Xuanhuan strictly dictates that while innate talent speeds up cultivation, ruthless effort, grinding, and ruthlessness in securing resources are what ultimately determine survival.

## Simulation / Implementation Concepts for CLI
- **Hybrid Power Aggregate**: Xuanhuan requires characters to possess multiple `PowerSystem` trees simultaneously (e.g., `MageTree` + `WarriorTree`). The `syncPower` function must aggregate these multiplicatively.
- **Resource Economy Integration**: `Cultivate` actions should require `Inventory` checks. To simulate the Xuanhuan grind, breaking a bottleneck should explicitly consume `ItemStack` items like "Spirit Beast Core" or "Foundation Pill," tying the `Inventory` and `Progression` domains together.
