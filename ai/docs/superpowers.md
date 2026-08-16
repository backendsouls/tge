# SuperPower Power System

This document outlines the design for the new `SuperPower` power system, which categorizes various known superpowers into distinct categories and progression levels.

## Overview

The **SuperPower** system models abilities commonly found in comic books and superhero media. Unlike traditional Cultivation paths that require slow gathering of Qi or spiritual energy, superpowers are often innate, awakened through trauma, or gained through scientific accidents/genetic mutation. 

This system functions as a distinct `PowerSystem` within the domain, under its own
`powersystem.PowerSystemType` value: **`SuperPower`** (alongside `Cultivation`, `Magic`,
`Reiatsu` and `Gamer`).

## Super Power Categories (Kinds)

Powers are divided into five main branches, modelled as root `PowerNode`s in the system's
DAG with the specific powers as children (`EdgeParent`):

1. **Physical Enhancements**
   - *Super Strength:* Immense physical power.
   - *Super Speed:* Moving and thinking at superhuman velocities.
   - *Invulnerability:* High resistance to physical damage.
   - *Enhanced Senses:* Heightened hearing, vision, smell, etc.

2. **Mental / Psychic Abilities**
   - *Telepathy:* Reading and communicating through minds.
   - *Telekinesis:* Moving objects with the mind.
   - *Precognition:* Seeing into the future.
   - *Mind Control:* Imposing will over others.

3. **Energy Manipulation**
   - *Pyrokinesis:* Creation and control of fire.
   - *Electrokinesis:* Creation and control of electricity.
   - *Energy Projection:* Firing concussive or destructive energy blasts.
   - *Force Fields:* Creating defensive barriers of energy.

4. **Biological / Alteration**
   - *Regeneration:* Rapid healing from injuries.
   - *Shapeshifting:* Altering physical appearance and form.
   - *Invisibility:* Bending light to become unseen.
   - *Elasticity:* Stretching and contorting the body.

5. **Reality / Spacial Manipulation**
   - *Teleportation:* Instantaneous spatial displacement.
   - *Time Manipulation:* Slowing, accelerating, or rewinding time.
   - *Gravity Control:* Altering gravitational fields.

## Progression Levels

Instead of standard Cultivation "Realms", the progression of a specific superpower is measured in **Tiers** or **Levels**. As a character progresses, their `PowerState` (e.g., `SuperPowerState`) advances through these ranks, unlocking higher potential and lifespan modifiers.

The multipliers below are the ones `superpower.SuperPowerState.Power()` returns for tiers
0–5 (`NewSuperPowerState` rejects anything outside that range). The lifespan modifiers are
design intent only — `SuperPowerState` has no lifespan method.

0. **Level 0: Base / Dormant**
   - *Description:* The baseline level, often before awakening or when the power is completely suppressed.
   - *Power Multiplier:* 1.0x

1. **Level 1: Latent / Novice**
   - *Description:* The power has just awakened. Control is poor, and it often manifests only under extreme stress or emotion.
   - *Power Multiplier:* 5x
   - *Lifespan Modifier:* None (unless biological)

2. **Level 2: Enhanced / Basic**
   - *Description:* The user has conscious control over their power but lacks stamina and precision. Suitable for street-level encounters.
   - *Power Multiplier:* 10x

3. **Level 3: Advanced / Mastery**
   - *Description:* Full mastery over the base ability. The user can perform complex maneuvers and sustain their power for extended periods.
   - *Power Multiplier:* 20x
   - *Lifespan Modifier:* Minor increase (e.g., +50 years)

4. **Level 4: Expert / City-Level**
   - *Description:* The power reaches a devastating scale. The user can affect entire city blocks or overpower armies.
   - *Power Multiplier:* 50x
   - *Lifespan Modifier:* Moderate increase

5. **Level 5: Omega / Planetary**
   - *Description:* The absolute pinnacle of the power. Omega-level mutants or cosmic entities capable of affecting entire planets or bending fundamental laws of reality.
   - *Power Multiplier:* 100x+
   - *Lifespan Modifier:* Immortality or vast increase (e.g., +10,000 years)

## Integration status

**Shipped.**

1. `powersystem.SuperPower` exists as a `PowerSystemType`.
2. `superpower.SuperPowerState{Tier}` implements `power.PowerState`, with the tier→multiplier
   table above and bounds checking on construction.
3. `defaults.yml` seeds a `SuperPower` system whose roots are the five categories, each with
   its specific powers attached as children via `EdgeParent`.

**Not yet wired.**

- Nothing constructs a `SuperPowerState`. A character's shipped progression runs through
  `character.NodeProgress` (a per-node `Level` + `Progress` against a `100 × Level²` gate),
  which is system-agnostic and ignores the tier table entirely — so a superpower node
  currently progresses on the same curve as a cultivation node. Mapping tier ↔ node level,
  or routing `SuperPower` systems through `SuperPowerState`, is the open work.
- Lifespan modifiers have no representation anywhere.
- The categories seed as plain nodes with `BasePower` 0 and category `"General"`; per-power
  `BasePower`, `StatVector`, `MaterialReq` and `Drawbacks` are unset.

See [decisions.md](decisions.md) §5 for why the `PowerState` implementations are kept
staged rather than deleted or wired in.
