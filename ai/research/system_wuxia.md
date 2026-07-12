# Power System: Wuxia (Martial Heroes)

## Overview
Wuxia is a "low-fantasy" Chinese fiction genre centered around martial artists operating in the *Jianghu* (an underground society of sects, wanderers, and outlaws). Unlike its high-fantasy counterpart *Xianxia*, Wuxia characters are essentially human; they cannot fly or destroy mountains, but through intense training and life energy (Qi), they achieve superhuman athleticism and combat prowess.

## Core Mechanics
1. **Qi and Neigong (Internal Arts)**
   Characters cultivate internal energy (Qi) through meditation and breathing exercises (Neigong). This energy is stored in the *Dantian* and circulated through meridians. It is used to enhance strength, heal wounds, and project force (e.g., palm strikes that hit from a distance).

2. **Wai-gong (External Arts)**
   Physical techniques, weaponry, and body conditioning. In the Wuxia hierarchy, pure external martial arts are usually considered inferior to internal arts, though the ultimate goal is perfect harmony between the two.

3. **Specialized Techniques**
   - **Qinggong (Lightness Skill)**: Using Qi to reduce effective body weight, allowing characters to leap over houses, run across water, or glide on tree branches.
   - **Dianxue (Acupoint Striking)**: Striking specific meridian points on a body to paralyze, mute, kill, or even heal the target.

4. **Qi Deviation**
   A severe mechanical drawback. If a martial artist trains incorrectly, rushes their progress, or experiences extreme emotional trauma, their Qi can flow backward or become chaotic. This causes "Qi Deviation," resulting in crippled cultivation, physical paralysis, or homicidal madness.

## Common Tropes
- **The Jianghu**: The lawless, honor-bound "rivers and lakes" where martial artists operate outside the Emperor's laws.
- **The Xia Code**: A chivalric code mandating loyalty, righteousness, and protecting the weak.
- **Hidden Manuals**: Plots often revolve around the discovery of or fight over a lost, supreme martial arts manual.
- **Energy Transfer**: A classic trope where an older master forcefully transfers their lifetime of cultivated Qi into a young protagonist, often dying in the process but instantly leveling up the youth.

## Simulation / Implementation Concepts for CLI
- **Dual-Track Progression**: `WuxiaState` could require separate tracking for `InternalQi` (mana/stamina) and `ExternalConditioning` (HP/Defense).
- **Qi Deviation Risk**: When attempting a `Train` or `Breakthrough` action, a failed roll doesn't just result in 0 points gained; it triggers a `QiDeviation` affliction that temporarily halves the character's stats until cured by specific medicine.
