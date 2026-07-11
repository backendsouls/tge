// Package sqlite is a driven adapter that persists realms in a SQLite database
// using the pure-Go modernc.org/sqlite driver (no CGO required).
package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"tge/internal/core/domain/cultivation"
	"tge/internal/core/port"

	sqlitedrv "modernc.org/sqlite" // registers the "sqlite" database/sql driver
	sqlitelib "modernc.org/sqlite/lib"
)

// RealmRepository implements port.RealmRepository over SQLite.
type RealmRepository struct {
	db *sql.DB
}

// NewRealmRepository opens (creating if needed) the database at dsn and applies
// the schema. dsn is a file path, or ":memory:" for an ephemeral database.
func NewRealmRepository(dsn string) (*RealmRepository, error) {
	db, err := open(dsn)
	if err != nil {
		return nil, err
	}
	return &RealmRepository{db: db}, nil
}

// Close releases the underlying database handle.
func (r *RealmRepository) Close() error {
	return r.db.Close()
}

// Save inserts a realm, returning port.ErrRealmExists if the name is taken.
func (r *RealmRepository) Save(ctx context.Context, realm cultivation.Realm) error {
	const q = `
INSERT INTO realms (
	name, power_multiplier, power_adder,
	lifespan_multiplier, lifespan_adder,
	bottleneck_points, max_levels, main_max_levels, tier
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, q,
		realm.Name, realm.PowerMultiplier, realm.PowerAdder,
		realm.LifespanMultiplier, realm.LifespanAdder,
		realm.BottleneckPoints, realm.MaxLevels, realm.MainCharacterMaxLevels, realm.Tier,
	)
	if err != nil {
		if serr, ok := errors.AsType[*sqlitedrv.Error](err); ok {
			switch serr.Code() {
			case sqlitelib.SQLITE_CONSTRAINT_PRIMARYKEY, sqlitelib.SQLITE_CONSTRAINT_UNIQUE:
				return fmt.Errorf("%w: %q", port.ErrRealmExists, realm.Name)
			}
		}
		return fmt.Errorf("save realm: %w", err)
	}
	return nil
}

// FindByName returns the named realm or port.ErrRealmNotFound.
func (r *RealmRepository) FindByName(ctx context.Context, name string) (cultivation.Realm, error) {
	const q = `
SELECT name, power_multiplier, power_adder,
       lifespan_multiplier, lifespan_adder,
       bottleneck_points, max_levels, main_max_levels, tier
FROM realms WHERE name = ?`

	realm, err := scanRealm(r.db.QueryRowContext(ctx, q, name))
	if errors.Is(err, sql.ErrNoRows) {
		return cultivation.Realm{}, fmt.Errorf("%w: %q", port.ErrRealmNotFound, name)
	}
	if err != nil {
		return cultivation.Realm{}, fmt.Errorf("find realm: %w", err)
	}
	if err := r.loadLevels(ctx, &realm); err != nil {
		return cultivation.Realm{}, err
	}
	return realm, nil
}

// AddLevel inserts an ordered level into a realm, returning
// cultivation.ErrLevelNumberExists if the number is already used there.
func (r *RealmRepository) AddLevel(ctx context.Context, realm string, l cultivation.Level) error {
	const q = `INSERT INTO realm_levels (realm, number, name, breakthrough_points, bottleneck_points) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, q, realm, l.Number, l.Name, l.BreakthroughPoints, l.BottleneckPoints)
	if err != nil {
		if serr, ok := errors.AsType[*sqlitedrv.Error](err); ok {
			switch serr.Code() {
			case sqlitelib.SQLITE_CONSTRAINT_PRIMARYKEY, sqlitelib.SQLITE_CONSTRAINT_UNIQUE:
				return fmt.Errorf("%w: %d", cultivation.ErrLevelNumberExists, l.Number)
			}
		}
		return fmt.Errorf("add realm level: %w", err)
	}
	return nil
}

// loadLevels populates a realm's levels ordered by number.
func (r *RealmRepository) loadLevels(ctx context.Context, realm *cultivation.Realm) error {
	rows, err := r.db.QueryContext(ctx,
		`SELECT number, name, breakthrough_points, bottleneck_points FROM realm_levels WHERE realm = ? ORDER BY number`, realm.Name)
	if err != nil {
		return fmt.Errorf("load realm levels: %w", err)
	}
	defer func() { _ = rows.Close() }()

	realm.Levels = nil
	for rows.Next() {
		var l cultivation.Level
		if err := rows.Scan(&l.Number, &l.Name, &l.BreakthroughPoints, &l.BottleneckPoints); err != nil {
			return fmt.Errorf("scan realm level: %w", err)
		}
		realm.Levels = append(realm.Levels, l)
	}
	return rows.Err()
}

// List returns all realms ordered by name.
func (r *RealmRepository) List(ctx context.Context) ([]cultivation.Realm, error) {
	const q = `
SELECT name, power_multiplier, power_adder,
       lifespan_multiplier, lifespan_adder,
       bottleneck_points, max_levels, main_max_levels, tier
FROM realms ORDER BY tier, name`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list realms: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var realms []cultivation.Realm
	for rows.Next() {
		realm, err := scanRealm(rows)
		if err != nil {
			return nil, fmt.Errorf("scan realm: %w", err)
		}
		realms = append(realms, realm)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate realms: %w", err)
	}
	for i := range realms {
		if err := r.loadLevels(ctx, &realms[i]); err != nil {
			return nil, err
		}
	}
	return realms, nil
}

// scanner is satisfied by both *sql.Row and *sql.Rows.
type scanner interface {
	Scan(dest ...any) error
}

func scanRealm(s scanner) (cultivation.Realm, error) {
	var realm cultivation.Realm
	err := s.Scan(
		&realm.Name, &realm.PowerMultiplier, &realm.PowerAdder,
		&realm.LifespanMultiplier, &realm.LifespanAdder,
		&realm.BottleneckPoints, &realm.MaxLevels, &realm.MainCharacterMaxLevels, &realm.Tier,
	)
	return realm, err
}
