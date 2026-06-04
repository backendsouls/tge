package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
	"sync"

	"github.com/pressly/goose/v3"
)

// migrationsFS holds the goose SQL migrations applied on open. They are embedded
// so the binary is self-contained — no migration files ship alongside it.
//
//go:embed migrations/*.sql
var migrationsFS embed.FS

// gooseInit configures goose's global state (dialect, embedded FS, silent
// logger) exactly once, regardless of how many repositories call open.
var gooseInit sync.Once

// open opens the database at dsn and brings its schema up to date by applying
// any pending goose migrations.
func open(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	var initErr error
	gooseInit.Do(func() {
		goose.SetBaseFS(migrationsFS)
		goose.SetLogger(goose.NopLogger()) // keep migration noise out of CLI output
		initErr = goose.SetDialect("sqlite3")
	})
	if initErr != nil {
		db.Close()
		return nil, fmt.Errorf("configure migrations: %w", initErr)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		db.Close()
		return nil, fmt.Errorf("apply migrations: %w", err)
	}
	return db, nil
}
