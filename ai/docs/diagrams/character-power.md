# A Character's Power: Shared Definition vs. Held Progress

```mermaid
%%{init: {"theme":"dark","themeVariables":{"lineColor":"#c8d3f5"},"themeCSS":".edgePath path,.flowchart-link{stroke-width:2px} .messageLine0,.messageLine1{stroke-width:2px} .relation{stroke-width:2px} .actor{stroke-width:2px} .node rect,.node circle,.node polygon,.node path{stroke-width:2px} .cluster rect{stroke-width:2px}"}}%%
flowchart LR
    subgraph CH["Character (data/characters/lin_feng.json)"]
        direction TB
        SYS["Systems []string<br/>['Cultivation']  — membership by NAME"]
        MS["MechanicState<br/>Tier · BasePower · IsAwakened<br/>Alignment · EnergyPools · SpellSlots"]
        UN["UnlockedNodes []NodeProgress<br/>{System, NodeID, Level, Progress, BasePower}"]
        PV["PowerValue string<br/>(derived, recomputed on every mutation)"]
    end

    subgraph DEF["Shared definition — illustrative graph, not the shipped catalog"]
        direction TB
        PS["PowerSystem{Name, PowerSystemType: Cultivation}<br/>Nodes map[string]*PowerNode"]
        N1["'spirit'<br/>BasePower, StatVector, Tags"]
        N2["'body'"]
        N3["'soul'"]
        N4["'golden_core'<br/>Parents: [spirit, body]"]
        N5["'demonic_core'<br/>MutuallyExclusive: [golden_core]"]
        PS --> N1
        PS --> N2
        PS --> N3
        N1 --> N4
        N2 --> N4
        N4 -.->|mutually exclusive| N5
    end

    SYS -->|resolved by name at use time| PS
    UN -->|NodeID| N1
    UN -->|NodeID| N4
    MS --> PV
    UN --> PV

    classDef held fill:#16351f,stroke:#7ad08a,color:#e5e7eb,stroke-width:2px
    classDef def fill:#1f2a44,stroke:#7aa2f7,color:#e5e7eb,stroke-width:2px
    class SYS,MS,UN,PV held
    class PS,N1,N2,N3,N4,N5 def
```

> The graph above is drawn to show what the DAG *can* express. The seeded `Cultivation`
> system is actually flat — three parentless roots (`spirit`, `body`, `soul`) and no edges
> at all; `golden_core`/`demonic_core` are invented. The only seeded system with real parent
> edges is `SuperPower`, whose five category roots each have children. Nothing in the
> shipped catalog uses `EdgeSibling` or `EdgeMutuallyExclusive`.

A `PowerSystem` is a **shared definition** and holds no per-character state whatsoever.
It is a DAG of `PowerNode`s in a flat `map[string]*PowerNode` keyed by a slug ID derived
from the node name, with shape carried by edges *on the nodes*: `Parents` (multi-parent —
`golden_core` needs both `spirit` and `body`), `Siblings`, and `MutuallyExclusive` (written
symmetrically onto both nodes). A parent edge that would close a cycle is rejected with
`ErrCyclicDependency`, which is what keeps it a DAG rather than a general graph.

The **character** holds three separate things pointing at that definition:

- `Systems []string` — membership, by name. The definition is loaded from the repository
  when it's actually needed.
- `MechanicState` — everything true about the character *across* systems: tier, awakening,
  alignment, energy pools, spell slots, vows, permanent traits. Stored once, not
  duplicated per system.
- `UnlockedNodes []NodeProgress` — a short list of `{System, NodeID, Level, Progress,
  BasePower}` records for the individual nodes it has actually unlocked.

`PowerValue` is not stored state but a rendered derivation, recomputed by
`CalculateTotalPower()` on every mutation:
`(MechanicState.BasePower + Σ node.BasePower × node.Level) × Π species.Power`, floored at
`1.0`. Note the product over **all** species — that's what makes hybrids compound.

This replaced an earlier design where a character carried two `[]PowerSystem` fields —
tree *definitions* and progressed *instance* copies with state hung on each node. That
duplicated the whole tree per character to annotate a few nodes. See
[../decisions.md](../decisions.md) §1–§2 for the full rationale.
