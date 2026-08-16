# Hexagonal Architecture

```mermaid
%%{init: {"theme":"dark","themeVariables":{"lineColor":"#c8d3f5"},"themeCSS":".edgePath path,.flowchart-link{stroke-width:2px} .messageLine0,.messageLine1{stroke-width:2px} .relation{stroke-width:2px} .actor{stroke-width:2px} .node rect,.node circle,.node polygon,.node path{stroke-width:2px} .cluster rect{stroke-width:2px}"}}%%
flowchart TB
    subgraph Root["cmd/tge — composition root"]
        MAIN["main + seed<br/>wires everything, seeds the catalog"]
    end

    subgraph Driving["Inbound / driving adapter"]
        CLI["cli.App<br/>parses commands, renders output"]
    end

    subgraph Core["internal/core — the hexagon"]
        direction TB
        DP["Driving ports<br/>CharacterService, RealmService,<br/>PowerSystemService, NovelService, …"]
        SVC["Service implementations<br/>CharacterService (+IdleService),<br/>DefaultWorldService, …"]
        DOMAIN["Domain<br/>character · powersystem · power · cultivation<br/>cosmology · novel · rpg · superpower"]
        DR["Driven ports<br/>CharacterRepository, PowerSystemRepository,<br/>RealmRepository, …"]
    end

    subgraph Driven["Outbound / driven adapters"]
        FILE["file repositories<br/>(aggregate JSON)"]
        SQLITE["sqlite repositories"]
        JSON[("data/characters/*.json<br/>data/power_systems/*.json")]
        DB[("SQLite file<br/>goose migrations")]
    end

    CFG["internal/config<br/>embedded defaults.yml"]
    LOG["internal/logger<br/>Dev → stderr · System → stdout"]

    MAIN --> CLI
    MAIN -. constructs .-> SVC
    MAIN -. constructs .-> SQLITE
    MAIN -. constructs .-> FILE
    CFG -. defaults + catalog .-> MAIN
    CLI -. novel beats .-> LOG
    CLI --> DP
    DP -. implemented by .-> SVC
    SVC --> DOMAIN
    SVC --> DR
    DR -. implemented by .-> SQLITE
    DR -. implemented by .-> FILE
    SQLITE --> DB
    FILE --> JSON
```

The application is **Ports and Adapters**: every dependency points *inward* toward
the domain. The CLI (inbound adapter) only ever calls **driving ports** — interfaces
like `port.CharacterService` — never a concrete service. Those ports are implemented
by the **service** layer, which orchestrates the pure **domain** and reaches storage
only through **driven ports** like `port.CharacterRepository`.

There are **two** driven adapter families behind those ports. `adapter/sqlite` backs the
row-shaped entities (realms and levels, the cosmology tiers, timelines, species, novels
and the whole RPG catalogue), opening the database through `open(dsn)` which applies
embedded goose migrations. `adapter/file` backs the two aggregate-shaped entities —
`Character` and `PowerSystem` — as single JSON documents, because a character's unlocked
nodes / idle slots / mechanic state and a system's whole DAG are object graphs, not rows.
Swapping either one is invisible to the core: that substitutability is the point of the
port. See [../decisions.md](../decisions.md) §3.

The one place that knows about concrete types is `cmd/tge`, the **composition root**,
which loads the embedded defaults, opens the database, constructs both repository
families, injects them into the services, seeds the starter catalog idempotently, and
hands the services to the CLI. This inversion is what lets the entire CLI be exercised
with in-memory fakes (no process or storage) and lets the domain stay ignorant of SQLite,
goose, JSON and flags.
