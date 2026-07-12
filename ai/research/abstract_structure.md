# Abstract Structure: Universal Power System Framework

## Objective
To extend the existing `tge-go` backend to support the vast array of mechanics observed across 20+ distinct power systems—ranging from linear hyper-scaling (Dragon Ball) to strict resource scarcity (D&D Vancian), conceptual reality bending (SCP/Matrix), and psychological horror progression (LotM).

## Core Architectural Enhancements

To accommodate this diversity, the rigid, vertical `PowerSystem` tree must evolve into a **Component-Based Directed Acyclic Graph (DAG)**, and the simulation engine must transition from static linear math to a contextual rules engine.

---

### 1. The Universal Power Node (The DAG)
Power can no longer be a strict parent-child tree. Systems like Path of Exile (Web), Mistborn (4x4 Grid), or One Piece (Singular Fruit) require a flexible node structure.

```go
type PowerNode struct {
    ID          string
    Name        string
    Category    string            // e.g., "Mover" (Worm), "Logia" (One Piece)
    Tags        []string          // e.g., ["fire", "projectile", "verbal_component"]
    
    // Graph Connectivity
    Parents     []string          // Directed prereqs
    Siblings    []string          // Horizontal/undirected connections (PoE style)
    MutuallyExclusive []string    // Mistborn: Can't have two base metals
    
    // Base Metrics
    BasePower   float64
    StatVector  map[string]string // JoJo parameters (DestructivePower: "A", Speed: "C") or SCP ACS classes
    
    // Constraints & Drawbacks
    Drawbacks   []Drawback        // Path of Exile keystones
    MaterialReq map[string]int    // Inventory checks for Fullmetal Alchemy (Equivalent Exchange)
}
```

---

### 2. The Unified Mechanic State (Character Level)
The character's progression state must track significantly more than just `Tier`. It must hold the permutations of Vows, Awakenings, and Alignments.

```go
type MechanicState struct {
    // Linear Scaling
    BasePower        float64         // Dragon Ball Zenkai base
    Tier             int             // Standard level
    
    // Evolutionary Flags
    IsAwakened       bool            // My Hero Academia / One Piece Awakenings
    PermanentTraits  []string        // Cradle Iron Bodies, Mistborn Savantism
    
    // Modifiers and Rules
    Vows             []Vow           // Jujutsu Kaisen / HxH conditional multipliers
    Alignment        float64         // Mental state: LotM Acting, Magi Rukh Depravity, Xianxia Sanity
    
    // Resource Pools
    EnergyPools      map[string]int  // Tracks Mana, Od, Shinsu, or Stamina
    SpellSlots       map[int]int     // D&D Vancian strict cast limits
}
```

---

### 3. The Simulation Engine (Context & Physics)
Combat math can no longer rely on `Damage = Power - Defense`. It must evaluate context, reality-bending, and absolute immunities.

```go
type CombatContext struct {
    Environment      string          // e.g., "Ocean", "SolarEclipse" (Avatar Bending buffs/debuffs)
    AmbientDensity   float64         // Tower of God Shinsu density checks
    LocalHumeLevel   float64         // SCP reality stability
    ActiveDomains    []Domain        // Jujutsu Kaisen Domain Expansions overwriting the arena
}

// Example Evaluation Engine
func (s MechanicState) ExecuteAction(node PowerNode, target *Character, ctx CombatContext) (Result, error) {
    // 1. Environmental Lethality Check (Tower of God)
    if s.Stats.Resistance < ctx.AmbientDensity {
        return InstantlyCrushed()
    }
    
    // 2. Resource & Component Constraints (D&D / Fullmetal / Xuanhuan)
    if node.HasTag("verbal") && target.HasStatus("Silenced") {
        return SpellFailed()
    }
    if !ConsumeInventory(node.MaterialReq) {
        return TruthRebound() // FMA consequence for failed equivalence
    }
    
    // 3. Absolute Hierarchy & Immunities (One Piece / Xianxia)
    if target.HasTrait("Logia_Intangibility") && !s.HasTrait("Armament_Haki") {
        return ZeroDamage()
    }
    
    // 4. Meta-Administrative Overrides (The Matrix)
    if s.HasTrait("The_Anomaly") {
        target.Stats.Damage = 0 // Directly editing the opponent's struct
    }
    
    // 5. Calculate Final Output (JJK Vows * DBZ Multipliers)
    power := calculateBase(node, s)
    power *= EvaluateVows(s.Vows, ctx)
    
    return ApplyDamage(power)
}
```

---

### 4. Broad Archetype Support Hooks
To fully accommodate the researched literature, the Go backend must support these macro-mechanics:
- **Zenkai Progression**: A hook that triggers post-simulation. If `Character.HP < 5%`, permanently increase `MechanicState.BasePower`.
- **Recoil & Feedback**: Quirk Singularity (MHA) and Qi Deviation (Wuxia). If node power output exceeds physical constitution thresholds, deal recoil damage or apply a permanent `Paralyzed` / `Madness` affliction.
- **Weapon Binding**: Djinn Equips (Magi) requiring `Inventory.EquippedWeapon` to match a specific `MetalVessel` ID to unlock Extreme Magic nodes.
- **AoE Pacification**: Song Energy (Macross) emitting fold waves that don't deal damage, but trigger `Willpower` checks to flip an enemy's hostile alignment to neutral.
- **Hume Reality Overrides**: If `Character.InternalHume > Context.LocalHumeLevel`, physics calculations are bypassed in favor of the character (SCP Reality Bending).

## Summary
By decomposing power into **Graphs (Structure)**, **MechanicStates (Conditionals)**, and an **Environmental Rules Engine (Context)**, the `tge-go` framework can dynamically simulate everything from throwing a simple Fireball to altering the server's gravity parameters as The One.
