# Implementation Plan: Universal Power System Framework

This document outlines the step-by-step refactoring plan required to transition the `tge-go` codebase from a linear Cultivation-focused tree to the Universal DAG and Contextual Rules Engine defined in `ai/research/abstract_structure.md`.

## 1. Refactor `progression.PowerNode` (Replacing `Power`)
**Target File**: `internal/core/domain/progression/power.go`

Currently, `Power` is a strict parent-child tree (`Children []Power`). We must evolve this into a **Directed Acyclic Graph (DAG)** to support non-linear webs, grids, and mutually exclusive powers.

**Action Items**:
- Rename the `Power` struct to `PowerNode`.
- Add a unique `ID` (string) to reference nodes in a flat map instead of a deep tree.
- Replace `Children []Power` with relational string arrays: `Parents []string`, `Siblings []string`, and `MutuallyExclusive []string`.
- Add taxonomy fields: `Category string` (e.g., "Logia", "Mover") and `Tags []string` (e.g., ["verbal", "fire"]).
- Add stat fields: `BasePower float64` and `StatVector map[string]string`.
- Add constraint fields: `Drawbacks []string` and `MaterialReq map[string]int`.

## 2. Introduce `MechanicState` (Evolving `PowerState`)
**Target File**: `internal/core/domain/progression/power_system.go` & `character.go`

Currently, progression is tracked via `PowerState` (implemented by `CultivationState` and `SuperPowerState`), which only returns a `Kind()` and `Power() float64`. We need to track multi-dimensional progression mechanics.

**Action Items**:
- Create a new `MechanicState` struct (or `UniversalPowerState` that implements `PowerState`) containing:
  - `BasePower float64` (for Zenkai boosts)
  - `Tier int`
  - `IsAwakened bool` (MHA/One Piece flag)
  - `PermanentTraits []string` (Iron Bodies, Savantism)
  - `Vows []string` (JJK Binding Vows)
  - `Alignment float64` (Depravity, Sanity)
  - `EnergyPools map[string]int` (Mana, Shinsu)
  - `SpellSlots map[int]int` (D&D Vancian magic)
- Update `Character` in `internal/core/domain/character/character.go` to hold this `MechanicState`, either at the root level or nested per `PowerSystem`.

## 3. Flatten `PowerSystem` Structure
**Target File**: `internal/core/domain/progression/power_system.go`

Currently, `PowerSystem` holds a forest of roots: `Powers []Power`.

**Action Items**:
- Refactor `PowerSystem` to hold a flat map of nodes: `Nodes map[string]*PowerNode`.
- Ensure tree traversal methods (`Names()`, `findInForest()`) are updated to walk the DAG using the `Parents` and `Siblings` relations.

## 4. Implement the `CombatContext` Rules Engine
**Target File**: `internal/core/domain/progression/combat.go` (New File)

We must move away from static `Damage = Power - Defense` math.

**Action Items**:
- Create the `CombatContext` struct:
  - `Environment string` (e.g., "Ocean", "SolarEclipse")
  - `AmbientDensity float64` (Tower of God)
  - `LocalHumeLevel float64` (SCP Reality)
- Implement an `ExecuteAction(node PowerNode, target *Character, ctx CombatContext)` function.
- **Rules to implement inside `ExecuteAction`**:
  - *Environmental Lethality*: Check `Target.Resistance` vs `ctx.AmbientDensity`.
  - *Inventory Constraints*: Check `node.MaterialReq` against `Character.Inventory`.
  - *Status Constraints*: Fail if `node` has "verbal" tag and `Character` has "Silenced" status.
  - *Absolute Immunities*: Intangibility (Logia) vs specific counters (Armament Haki).
  - *Meta-Overrides*: If `Character` has "The_Anomaly" trait, allow direct modification of enemy structs.

## 5. Add Progression Hooks (Post-Simulation)
**Target File**: `internal/core/service/character_service.go`

**Action Items**:
- Add a `PostCombatUpdate` method to `CharacterService`.
- Implement **Zenkai Boosts**: If `HP < 5%`, permanently increment `MechanicState.BasePower`.
- Implement **Trauma Awakenings**: If `HP == 1`, small % chance to set `IsAwakened = true`.
- Implement **Truth Rebounds**: If an action fails due to missing materials, permanently decrement a physical stat.

## Next Steps
Once this plan is approved, we will begin by refactoring the `PowerNode` and `PowerSystem` domain models in `internal/core/domain/progression/`.
