# A Character's Power: Definition vs. Instance

```mermaid
%%{init: {"theme":"dark","themeVariables":{"lineColor":"#c8d3f5"},"themeCSS":".edgePath path,.flowchart-link{stroke-width:2px} .messageLine0,.messageLine1{stroke-width:2px} .relation{stroke-width:2px} .actor{stroke-width:2px} .node rect,.node circle,.node polygon,.node path{stroke-width:2px} .cluster rect{stroke-width:2px}"}}%%
flowchart LR
    subgraph CH["Character"]
        direction TB
        DEF["PowerSystems []PowerSystem<br/>(DEFINITIONS — what it can grow in)"]
        INST["Power []PowerSystem<br/>(INSTANCES — what it has grown)"]
    end

    subgraph DEFTREE["Definition tree (shared, authored once)"]
        D1["PowerSystem{Kind: Cultivation}"]
        D1 --> DP1["Power 'Spirit' (State = nil)"]
        D1 --> DP2["Power 'Body' (State = nil)"]
    end

    subgraph INSTTREE["Instance tree (this character's progress)"]
        I1["PowerSystem{Kind: Cultivation}"]
        I1 --> IP1["Power 'Spirit'<br/>State = CultivationState{Realm, Level, Points}"]
    end

    DEF --> DEFTREE
    INST --> INSTTREE

    classDef def fill:#1f2a44,stroke:#7aa2f7,color:#e5e7eb,stroke-width:2px
    classDef inst fill:#16351f,stroke:#7ad08a,color:#e5e7eb,stroke-width:2px
    class D1,DP1,DP2 def
    class I1,IP1 inst
```

A `Character` carries **two** `[]progression.PowerSystem` fields that use the same type
for two different jobs. `PowerSystems` are the **definitions** — the complete power
trees the character has access to, authored once and shared; their `Power` nodes carry
no state (`State == nil`). `Power` is the character's **instances** — its own progressed
copies, where each node it has advanced carries a `PowerState`. Keeping both as
`PowerSystem` (rather than a cultivation-specific type) is deliberate: `SystemKind`
distinguishes families (`Cultivation`, `Magic`), and the node-level `PowerState`
interface lets each family bring its own progress shape. Today the only implementation
is `CultivationState` (Realm + Level + accumulated Points + Progress), so a cultivation
instance node reads "at *this* Realm and Level"; a future Magic system would attach a
different `PowerState` to the same structure. `PowerValue` (a rendered string) is a
separate, third thing — a summary number, not the structured state. See
[../decisions.md](../decisions.md) §1–§2 for the rationale.
