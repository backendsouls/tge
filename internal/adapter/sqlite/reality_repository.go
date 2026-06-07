package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"tge/internal/core/domain/cosmology"
	"tge/internal/core/port"

	sqlitedrv "modernc.org/sqlite"
	sqlitelib "modernc.org/sqlite/lib"
)

type RealityRepository struct {
	db *sql.DB
}

func NewRealityRepository(dsn string) (*RealityRepository, error) {
	db, err := open(dsn)
	if err != nil {
		return nil, err
	}
	return &RealityRepository{db: db}, nil
}

func (r *RealityRepository) Close() error {
	return r.db.Close()
}

func (r *RealityRepository) Save(ctx context.Context, rl cosmology.Reality) error {
	const q = `INSERT INTO realities (name) VALUES (?)`
	_, err := r.db.ExecContext(ctx, q, rl.Name)
	if err != nil {
		if serr, ok := errors.AsType[*sqlitedrv.Error](err); ok {
			if serr.Code() == sqlitelib.SQLITE_CONSTRAINT_PRIMARYKEY || serr.Code() == sqlitelib.SQLITE_CONSTRAINT_UNIQUE {
				return port.ErrRealityExists
			}
		}
		return fmt.Errorf("save reality: %w", err)
	}
	return nil
}

func (r *RealityRepository) AddOmniverse(ctx context.Context, reality, omniverse string) error {
	const q = `INSERT INTO reality_omniverses (reality, omniverse) VALUES (?, ?)`
	_, err := r.db.ExecContext(ctx, q, reality, omniverse)
	if err != nil {
		if serr, ok := errors.AsType[*sqlitedrv.Error](err); ok {
			if serr.Code() == sqlitelib.SQLITE_CONSTRAINT_PRIMARYKEY || serr.Code() == sqlitelib.SQLITE_CONSTRAINT_UNIQUE {
				return cosmology.ErrRealityOmniverseExists
			}
		}
		return fmt.Errorf("add omniverse to reality: %w", err)
	}
	return nil
}

func (r *RealityRepository) FindByName(ctx context.Context, name string) (cosmology.Reality, error) {
	var rl cosmology.Reality
	err := r.db.QueryRowContext(ctx, `SELECT name FROM realities WHERE name = ?`, name).Scan(&rl.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return cosmology.Reality{}, fmt.Errorf("%w: %q", port.ErrRealityNotFound, name)
	}
	if err != nil {
		return cosmology.Reality{}, fmt.Errorf("find reality: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `SELECT omniverse FROM reality_omniverses WHERE reality = ? ORDER BY rowid`, name)
	if err != nil {
		return cosmology.Reality{}, fmt.Errorf("load reality omniverses: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var o string
		if err := rows.Scan(&o); err != nil {
			return cosmology.Reality{}, err
		}
		rl.Omniverses = append(rl.Omniverses, cosmology.Omniverse{Name: o})
	}
	return rl, rows.Err()
}

func (r *RealityRepository) List(ctx context.Context) ([]cosmology.Reality, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT name FROM realities ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list realities: %w", err)
	}
	defer rows.Close()

	var list []cosmology.Reality
	for rows.Next() {
		var rl cosmology.Reality
		if err := rows.Scan(&rl.Name); err != nil {
			return nil, err
		}
		list = append(list, rl)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for i := range list {
		or, err := r.db.QueryContext(ctx, `SELECT omniverse FROM reality_omniverses WHERE reality = ? ORDER BY rowid`, list[i].Name)
		if err != nil {
			return nil, err
		}
		for or.Next() {
			var o string
			if err := or.Scan(&o); err != nil {
				or.Close()
				return nil, err
			}
			list[i].Omniverses = append(list[i].Omniverses, cosmology.Omniverse{Name: o})
		}
		or.Close()
	}
	return list, nil
}
