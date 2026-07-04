package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"tge/internal/core/domain/cosmology"
	"tge/internal/core/domain/progression"
	"tge/internal/core/port"

	sqlitedrv "modernc.org/sqlite"
	sqlitelib "modernc.org/sqlite/lib"
)

// UniverseRepository implements port.UniverseRepository over SQLite. A universe's
// member power systems are stored in universe_systems (system unique, so a system
// belongs to at most one universe).
type UniverseRepository struct {
	db *sql.DB
}

// NewUniverseRepository opens (creating if needed) the database at dsn.
func NewUniverseRepository(dsn string) (*UniverseRepository, error) {
	db, err := open(dsn)
	if err != nil {
		return nil, err
	}
	return &UniverseRepository{db: db}, nil
}

// Close releases the underlying database handle.
func (r *UniverseRepository) Close() error {
	return r.db.Close()
}

// Create inserts an empty universe, returning port.ErrUniverseExists if taken.
func (r *UniverseRepository) Create(ctx context.Context, name string) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO universes (name) VALUES (?)`, name)
	if err != nil {
		if serr, ok := errors.AsType[*sqlitedrv.Error](err); ok {
			switch serr.Code() {
			case sqlitelib.SQLITE_CONSTRAINT_PRIMARYKEY, sqlitelib.SQLITE_CONSTRAINT_UNIQUE:
				return fmt.Errorf("%w: %q", port.ErrUniverseExists, name)
			}
		}
		return fmt.Errorf("create universe: %w", err)
	}
	return nil
}

// FindByName loads a universe and its member systems, or port.ErrUniverseNotFound.
func (r *UniverseRepository) FindByName(ctx context.Context, name string) (cosmology.Universe, error) {
	var found string
	err := r.db.QueryRowContext(ctx, `SELECT name FROM universes WHERE name = ?`, name).Scan(&found)
	if errors.Is(err, sql.ErrNoRows) {
		return cosmology.Universe{}, fmt.Errorf("%w: %q", port.ErrUniverseNotFound, name)
	}
	if err != nil {
		return cosmology.Universe{}, fmt.Errorf("find universe: %w", err)
	}
	u := cosmology.Universe{Name: found}
	if err := r.loadSystems(ctx, &u); err != nil {
		return cosmology.Universe{}, err
	}
	if err := r.loadRealms(ctx, &u); err != nil {
		return cosmology.Universe{}, err
	}
	return u, nil
}

// List returns all universes with their member systems, ordered by name.
func (r *UniverseRepository) List(ctx context.Context) ([]cosmology.Universe, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT name FROM universes ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list universes: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var universes []cosmology.Universe
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("scan universe: %w", err)
		}
		universes = append(universes, cosmology.Universe{Name: name})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate universes: %w", err)
	}
	for i := range universes {
		if err := r.loadSystems(ctx, &universes[i]); err != nil {
			return nil, err
		}
		if err := r.loadRealms(ctx, &universes[i]); err != nil {
			return nil, err
		}
	}
	return universes, nil
}

// SaveSystems replaces a universe's membership in a transaction.
func (r *UniverseRepository) SaveSystems(ctx context.Context, u cosmology.Universe) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `DELETE FROM universe_systems WHERE universe = ?`, u.Name); err != nil {
		return fmt.Errorf("clear universe systems: %w", err)
	}
	for _, s := range u.Systems {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO universe_systems (universe, system) VALUES (?, ?)`, u.Name, s.Name); err != nil {
			return fmt.Errorf("insert universe system %q: %w", s.Name, err)
		}
	}
	return tx.Commit()
}

// SaveRealms replaces a universe's in-universe realms in a transaction.
func (r *UniverseRepository) SaveRealms(ctx context.Context, u cosmology.Universe) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `DELETE FROM universe_realms WHERE universe = ?`, u.Name); err != nil {
		return fmt.Errorf("clear universe realms: %w", err)
	}
	for _, realm := range u.Realms {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO universe_realms (universe, name) VALUES (?, ?)`, u.Name, realm.Name); err != nil {
			return fmt.Errorf("insert universe realm %q: %w", realm.Name, err)
		}
	}
	return tx.Commit()
}

// loadRealms populates a universe's realms (locations) ordered by insertion.
func (r *UniverseRepository) loadRealms(ctx context.Context, u *cosmology.Universe) error {
	rows, err := r.db.QueryContext(ctx,
		`SELECT name FROM universe_realms WHERE universe = ? ORDER BY rowid`, u.Name)
	if err != nil {
		return fmt.Errorf("load universe realms: %w", err)
	}
	defer func() { _ = rows.Close() }()

	u.Realms = nil
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("scan universe realm: %w", err)
		}
		u.Realms = append(u.Realms, cosmology.Location{Name: name})
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate universe realms: %w", err)
	}
	return nil
}

// FindBySystem returns the universe a power system belongs to, or
// port.ErrUniverseNotFound.
func (r *UniverseRepository) FindBySystem(ctx context.Context, system string) (cosmology.Universe, error) {
	var name string
	err := r.db.QueryRowContext(ctx,
		`SELECT universe FROM universe_systems WHERE system = ?`, system).Scan(&name)
	if errors.Is(err, sql.ErrNoRows) {
		return cosmology.Universe{}, port.ErrUniverseNotFound
	}
	if err != nil {
		return cosmology.Universe{}, fmt.Errorf("find universe by system: %w", err)
	}
	return cosmology.Universe{Name: name}, nil
}

// loadSystems populates a universe's member systems (names) ordered by insertion.
func (r *UniverseRepository) loadSystems(ctx context.Context, u *cosmology.Universe) error {
	rows, err := r.db.QueryContext(ctx,
		`SELECT system FROM universe_systems WHERE universe = ? ORDER BY rowid`, u.Name)
	if err != nil {
		return fmt.Errorf("load universe systems: %w", err)
	}
	defer func() { _ = rows.Close() }()

	u.Systems = nil
	for rows.Next() {
		var system string
		if err := rows.Scan(&system); err != nil {
			return fmt.Errorf("scan universe system: %w", err)
		}
		u.Systems = append(u.Systems, progression.PowerSystem{Name: system})
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate universe systems: %w", err)
	}
	return nil
}

var _ port.UniverseRepository = (*UniverseRepository)(nil)
