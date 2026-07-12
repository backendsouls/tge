# Power System: LitRPG (Literary Role-Playing Game)

## Overview
LitRPG blends traditional narrative storytelling with the explicit, quantifiable mechanics of video games. In these universes, reality is governed by "The System"—a ubiquitous, often mysterious framework that manages physics, magic, and biological progression through math and UI screens.

## Core Mechanics
1. **The System Interface**: Characters perceive the world through blue screens, notifications, and character sheets. Every action is quantified.
2. **Stats vs. Skills**:
   - **Attributes/Stats** (STR, AGI, INT, LUK): Represent raw physical and mental capacity.
   - **Skills** (e.g., Swordsmanship Lv. 5, Fireball Lv. 2): Represent specialized knowledge and techniques, improving through repetitive use (Grinding).
3. **Levels and XP**: Defeating enemies (often monsters in dungeons) or completing System-issued "Quests" yields Experience Points. Accumulating enough XP triggers a Level Up, fully healing the user and granting stat points to distribute.
4. **Classes and Jobs**: Characters select overarching paths (e.g., Necromancer, Spellsword) that dictate stat growth rates and skill acquisition pools.

## Common Tropes & Challenges
- **Min-Maxing & System Hacks**: Protagonists often discover loopholes in the System, allocating points into hyper-specialized, "broken" builds that exploit obscure mechanics.
- **Stat Creep**: A major narrative hurdle where numbers become so astronomically high they lose meaning. Authors often counter this by introducing "Tiers" (E to SSS) or "Evolutions" that reset numbers but massively increase baseline quality.
- **Dungeon Crawling**: The world is often dotted with instanced zones containing escalating threats, traps, and boss monsters designed specifically to test players and distribute "Loot".

## Simulation / Implementation Concepts for CLI
- **Extensible Character Sheets**: The `Character` struct requires dynamic `map[string]int` fields to track diverse attributes (STR, AGI, LUK) and a dedicated `Inventory` for loot.
- **XP / Leveling Hooks**: Instead of manually assigning tiers, implement an `AddXP()` method that automatically handles overflow, triggers level-ups, and allocates attribute points.
- **Skill Proficiency Grinding**: Skills (nodes) should have their own localized XP bars that increment upon use, simulating the LitRPG trope of practicing a spell repeatedly to rank it up.
