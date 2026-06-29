# Hexagonal Architecture

```mermaid
%%{init: {"theme":"dark","themeVariables":{"lineColor":"#c8d3f5"},"themeCSS":".edgePath path,.flowchart-link{stroke-width:2px} .messageLine0,.messageLine1{stroke-width:2px} .relation{stroke-width:2px} .actor{stroke-width:2px} .node rect,.node circle,.node polygon,.node path{stroke-width:2px} .cluster rect{stroke-width:2px}"}}%%
flowchart TB
    subgraph Root["cmd/tge — composition root"]
        MAIN["main + seed<br/>wires everything"]
    end

    subgraph Driving["Inbound / driving adapter"]
        CLI["cli.App<br/>parses commands, renders output"]
    end

    subgraph Core["internal/core — the hexagon"]
        direction TB
        DP["Driving ports<br/>CharacterService, RealmService,<br/>UniverseService, NovelService, …"]
        SVC["Service implementations<br/>CharacterService, DefaultWorldService, …"]
        DOMAIN["Domain<br/>character · progression · cosmology · rpg · novel"]
        DR["Driven ports<br/>CharacterRepository, RealmRepository, …"]
    end

    subgraph Driven["Outbound / driven adapter"]
        SQLITE["sqlite repositories"]
        DB[("SQLite file<br/>goose migrations")]
    end

    MAIN --> CLI
    MAIN -. constructs .-> SVC
    MAIN -. constructs .-> SQLITE
    CLI --> DP
    DP -. implemented by .-> SVC
    SVC --> DOMAIN
    SVC --> DR
    DR -. implemented by .-> SQLITE
    SQLITE --> DB
```

The application is **Ports and Adapters**: every dependency points *inward* toward
the domain. The CLI (inbound adapter) only ever calls **driving ports** — interfaces
like `port.CharacterService` — never a concrete service. Those ports are implemented
by the **service** layer, which orchestrates the pure **domain** and reaches storage
only through **driven ports** like `port.CharacterRepository`. The SQLite adapter
implements those driven ports. The one place that knows about concrete types is
`cmd/tge`, the **composition root**, which opens the database, constructs the SQLite
repositories, injects them into the services, and hands the services to the CLI. This
inversion is what lets the entire CLI be exercised with in-memory fakes (no process or
storage) and lets the domain stay ignorant of SQLite, goose, and flags.
