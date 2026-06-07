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

type OmniverseRepository struct {
	db *sql.DB
}

func NewOmniverseRepository(dsn string) (*OmniverseRepository, error) {
	db, err := open(dsn)
	if err != nil {
		return nil, err
	}
	return &OmniverseRepository{db: db}, nil
}

func (r *OmniverseRepository) Close() error {
	return r.db.Close()
}

func (r *OmniverseRepository) Save(ctx context.Context, o cosmology.Omniverse) error {
	const q = `INSERT INTO omniverses (name) VALUES (?)`
	_, err := r.db.ExecContext(ctx, q, o.Name)
	if err != nil {
		if serr, ok := errors.AsType[*sqlitedrv.Error](err); ok {
			if serr.Code() == sqlitelib.SQLITE_CONSTRAINT_PRIMARYKEY || serr.Code() == sqlitelib.SQLITE_CONSTRAINT_UNIQUE {
				return port.ErrOmniverseExists
			}
		}
		return fmt.Errorf("save omniverse: %w", err)
	}
	return nil
}

func (r *OmniverseRepository) AddMultiverse(ctx context.Context, omniverse, multiverse string) error {
	const q = `INSERT INTO omniverse_multiverses (omniverse, multiverse) VALUES (?, ?)`
	_, err := r.db.ExecContext(ctx, q, omniverse, multiverse)
	if err != nil {
		if serr, ok := errors.AsType[*sqlitedrv.Error](err); ok {
			if serr.Code() == sqlitelib.SQLITE_CONSTRAINT_PRIMARYKEY || serr.Code() == sqlitelib.SQLITE_CONSTRAINT_UNIQUE {
				return cosmology.ErrOmniverseMultiverseExists
			}
		}
		return fmt.Errorf("add multiverse to omniverse: %w", err)
	}
	return nil
}

func (r *OmniverseRepository) FindByName(ctx context.Context, name string) (cosmology.Omniverse, error) {
	var o cosmology.Omniverse
	err := r.db.QueryRowContext(ctx, `SELECT name FROM omniverses WHERE name = ?`, name).Scan(&o.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return cosmology.Omniverse{}, fmt.Errorf("%w: %q", port.ErrOmniverseNotFound, name)
	}
	if err != nil {
		return cosmology.Omniverse{}, fmt.Errorf("find omniverse: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `SELECT multiverse FROM omniverse_multiverses WHERE omniverse = ? ORDER BY rowid`, name)
	if err != nil {
		return cosmology.Omniverse{}, fmt.Errorf("load omniverse multiverses: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var m string
		if err := rows.Scan(&m); err != nil {
			return cosmology.Omniverse{}, err
		}
		o.Multiverses = append(o.Multiverses, cosmology.Multiverse{Name: m})
	}
	return o, rows.Err()
}

func (r *OmniverseRepository) List(ctx context.Context) ([]cosmology.Omniverse, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT name FROM omniverses ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list omniverses: %w", err)
	}
	defer rows.Close()

	var list []cosmology.Omniverse
	for rows.Next() {
		var o cosmology.Omniverse
		if err := rows.Scan(&o.Name); err != nil {
			return nil, err
		}
		list = append(list, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for i := range list {
		mr, err := r.db.QueryContext(ctx, `SELECT multiverse FROM omniverse_multiverses WHERE omniverse = ? ORDER BY rowid`, list[i].Name)
		if err != nil {
			return nil, err
		}
		for mr.Next() {
			var m string
			if err := mr.Scan(&m); err != nil {
				mr.Close()
				return nil, err
			}
			list[i].Multiverses = append(list[i].Multiverses, cosmology.Multiverse{Name: m})
		}
		mr.Close()
	}
	return list, nil
}
