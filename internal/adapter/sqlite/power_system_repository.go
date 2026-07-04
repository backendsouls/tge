package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"tge/internal/core/domain/progression"
	"tge/internal/core/port"

	sqlitedrv "modernc.org/sqlite"
	sqlitelib "modernc.org/sqlite/lib"
)

// PowerSystemRepository implements port.PowerSystemRepository over SQLite,
// storing each system's power tree as an adjacency list.
type PowerSystemRepository struct {
	db *sql.DB
}

// NewPowerSystemRepository opens (creating if needed) the database at dsn.
func NewPowerSystemRepository(dsn string) (*PowerSystemRepository, error) {
	db, err := open(dsn)
	if err != nil {
		return nil, err
	}
	return &PowerSystemRepository{db: db}, nil
}

// Close releases the underlying database handle.
func (r *PowerSystemRepository) Close() error {
	return r.db.Close()
}

// Create inserts an empty system, returning port.ErrPowerSystemExists if taken.
func (r *PowerSystemRepository) Create(ctx context.Context, name string) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO power_systems (name) VALUES (?)`, name)
	if err != nil {
		if serr, ok := errors.AsType[*sqlitedrv.Error](err); ok {
			switch serr.Code() {
			case sqlitelib.SQLITE_CONSTRAINT_PRIMARYKEY, sqlitelib.SQLITE_CONSTRAINT_UNIQUE:
				return fmt.Errorf("%w: %q", port.ErrPowerSystemExists, name)
			}
		}
		return fmt.Errorf("create power system: %w", err)
	}
	return nil
}

// FindByName loads a system and its power tree, or port.ErrPowerSystemNotFound.
func (r *PowerSystemRepository) FindByName(ctx context.Context, name string) (progression.PowerSystem, error) {
	var found string
	err := r.db.QueryRowContext(ctx, `SELECT name FROM power_systems WHERE name = ?`, name).Scan(&found)
	if errors.Is(err, sql.ErrNoRows) {
		return progression.PowerSystem{}, fmt.Errorf("%w: %q", port.ErrPowerSystemNotFound, name)
	}
	if err != nil {
		return progression.PowerSystem{}, fmt.Errorf("find power system: %w", err)
	}

	ps := progression.PowerSystem{Name: found}
	if err := r.loadPowers(ctx, &ps); err != nil {
		return progression.PowerSystem{}, err
	}
	return ps, nil
}

// List returns all systems with their power trees, ordered by name.
func (r *PowerSystemRepository) List(ctx context.Context) ([]progression.PowerSystem, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT name FROM power_systems ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list power systems: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var systems []progression.PowerSystem
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("scan power system: %w", err)
		}
		systems = append(systems, progression.PowerSystem{Name: name})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate power systems: %w", err)
	}
	for i := range systems {
		if err := r.loadPowers(ctx, &systems[i]); err != nil {
			return nil, err
		}
	}
	return systems, nil
}

// SavePowers replaces the stored powers of a system with the given tree. It runs
// in a transaction so a failed rebuild never leaves a partial tree.
func (r *PowerSystemRepository) SavePowers(ctx context.Context, system progression.PowerSystem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `DELETE FROM powers WHERE system = ?`, system.Name); err != nil {
		return fmt.Errorf("clear powers: %w", err)
	}

	const ins = `INSERT INTO powers (system, name, parent) VALUES (?, ?, ?)`
	var insert func(powers []progression.Power, parent string) error
	insert = func(powers []progression.Power, parent string) error {
		for _, p := range powers {
			var parentArg any
			if parent != "" {
				parentArg = parent
			}
			if _, err := tx.ExecContext(ctx, ins, system.Name, p.Name, parentArg); err != nil {
				return fmt.Errorf("insert power %q: %w", p.Name, err)
			}
			if err := insert(p.Children, p.Name); err != nil {
				return err
			}
		}
		return nil
	}
	if err := insert(system.Powers, ""); err != nil {
		return err
	}
	return tx.Commit()
}

// loadPowers reads a system's adjacency rows and rebuilds its tree in place.
func (r *PowerSystemRepository) loadPowers(ctx context.Context, ps *progression.PowerSystem) error {
	rows, err := r.db.QueryContext(ctx,
		`SELECT name, parent FROM powers WHERE system = ? ORDER BY rowid`, ps.Name)
	if err != nil {
		return fmt.Errorf("load powers: %w", err)
	}
	defer func() { _ = rows.Close() }()

	// node keeps mutable children while we link the adjacency list.
	type node struct {
		name     string
		children []*node
	}
	nodes := map[string]*node{}
	type rel struct{ name, parent string }
	var rels []rel // preserves insertion order

	for rows.Next() {
		var name string
		var parent sql.NullString
		if err := rows.Scan(&name, &parent); err != nil {
			return fmt.Errorf("scan power: %w", err)
		}
		nodes[name] = &node{name: name}
		rels = append(rels, rel{name: name, parent: parent.String})
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate powers: %w", err)
	}

	var roots []*node
	for _, rl := range rels {
		if rl.parent == "" {
			roots = append(roots, nodes[rl.name])
			continue
		}
		if p := nodes[rl.parent]; p != nil {
			p.children = append(p.children, nodes[rl.name])
		}
	}

	var convert func(n *node) progression.Power
	convert = func(n *node) progression.Power {
		p := progression.Power{Name: n.name}
		for _, c := range n.children {
			p.Children = append(p.Children, convert(c))
		}
		return p
	}
	ps.Powers = nil
	for _, root := range roots {
		ps.Powers = append(ps.Powers, convert(root))
	}
	return nil
}

var _ port.PowerSystemRepository = (*PowerSystemRepository)(nil)
