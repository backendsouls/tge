# Persistence & Migrations (goose)

```mermaid
%%{init: {"theme":"dark","themeVariables":{"lineColor":"#c8d3f5"},"themeCSS":".edgePath path,.flowchart-link{stroke-width:2px} .messageLine0,.messageLine1{stroke-width:2px} .relation{stroke-width:2px} .actor{stroke-width:2px} .node rect,.node circle,.node polygon,.node path{stroke-width:2px} .cluster rect{stroke-width:2px}"}}%%
flowchart TB
    subgraph Ctor["Any repository constructor"]
        NEW["NewCharacterRepository(dsn)<br/>NewRealmRepository(dsn) …"]
    end

    NEW --> OPEN["open(dsn)"]
    OPEN --> SQLOPEN["sql.Open(\"sqlite\", dsn)"]
    OPEN --> ONCE{"gooseInit<br/>(sync.Once)"}
    ONCE -->|first call only| CFG["SetBaseFS(embed) +<br/>SetDialect(sqlite3) +<br/>SetLogger(NopLogger)"]
    OPEN --> UP["goose.Up(db, \"migrations\")"]

    EMB[["//go:embed migrations/*.sql<br/>00001_init.sql, 00002_…"]]
    EMB -. embedded FS .-> UP

    UP --> CHK{"pending<br/>migrations?"}
    CHK -->|yes| APPLY["apply in order,<br/>record in goose_db_version"]
    CHK -->|no| NOOP["no-op"]
    APPLY --> DB[("SQLite file")]
    NOOP --> DB
```

Schema is managed with **goose v3 as an embedded library**, not a CLI. The migration
SQL files live in `internal/adapter/sqlite/migrations/*.sql` (goose `-- +goose Up` /
`-- +goose Down` format) and are compiled into the binary via `//go:embed`, so it ships
self-contained. Every repository constructor calls `open(dsn)`, which opens the database
and runs `goose.Up`; goose's global config (embedded FS, `sqlite3` dialect, and a
**`NopLogger`** so migration output never leaks into CLI stdout) is set exactly once via
a `sync.Once`. `goose.Up` is **idempotent** — it records applied versions in the
`goose_db_version` table, so the many constructors that each call `open` are safe and
subsequent opens are no-ops. The whole existing schema is consolidated into
`00001_init.sql`; its `Up` keeps `CREATE TABLE IF NOT EXISTS` so a database made by the
earlier hand-rolled bootstrap adopts migrations cleanly. Future schema changes —
including column evolution (`ALTER TABLE`) that the old bootstrap could not do — go into
new numbered files (`00002_*.sql`, …). See
[../architecture.md](../architecture.md#database-migrations).
