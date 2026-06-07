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

type MultiverseRepository struct {
	db *sql.DB
}

func NewMultiverseRepository(dsn string) (*MultiverseRepository, error) {
	db, err := open(dsn)
	if err != nil {
		return nil, err
	}
	return &MultiverseRepository{db: db}, nil
}

func (r *MultiverseRepository) Close() error {
	return r.db.Close()
}

func (r *MultiverseRepository) Save(ctx context.Context, m cosmology.Multiverse) error {
	const q = `INSERT INTO multiverses (name) VALUES (?)`
	_, err := r.db.ExecContext(ctx, q, m.Name)
	if err != nil {
		if serr, ok := errors.AsType[*sqlitedrv.Error](err); ok {
			if serr.Code() == sqlitelib.SQLITE_CONSTRAINT_PRIMARYKEY || serr.Code() == sqlitelib.SQLITE_CONSTRAINT_UNIQUE {
				return port.ErrMultiverseExists
			}
		}
		return fmt.Errorf("save multiverse: %w", err)
	}
	return nil
}

func (r *MultiverseRepository) AddUniverse(ctx context.Context, multiverse, universe string) error {
	const q = `INSERT INTO multiverse_universes (multiverse, universe) VALUES (?, ?)`
	_, err := r.db.ExecContext(ctx, q, multiverse, universe)
	if err != nil {
		if serr, ok := errors.AsType[*sqlitedrv.Error](err); ok {
			if serr.Code() == sqlitelib.SQLITE_CONSTRAINT_PRIMARYKEY || serr.Code() == sqlitelib.SQLITE_CONSTRAINT_UNIQUE {
				return cosmology.ErrMultiverseUniverseExists
			}
		}
		return fmt.Errorf("add universe to multiverse: %w", err)
	}
	return nil
}

func (r *MultiverseRepository) FindByName(ctx context.Context, name string) (cosmology.Multiverse, error) {
	var m cosmology.Multiverse
	err := r.db.QueryRowContext(ctx, `SELECT name FROM multiverses WHERE name = ?`, name).Scan(&m.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return cosmology.Multiverse{}, fmt.Errorf("%w: %q", port.ErrMultiverseNotFound, name)
	}
	if err != nil {
		return cosmology.Multiverse{}, fmt.Errorf("find multiverse: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `SELECT universe FROM multiverse_universes WHERE multiverse = ? ORDER BY rowid`, name)
	if err != nil {
		return cosmology.Multiverse{}, fmt.Errorf("load multiverse universes: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var u string
		if err := rows.Scan(&u); err != nil {
			return cosmology.Multiverse{}, err
		}
		m.Universes = append(m.Universes, cosmology.Universe{Name: u})
	}
	return m, rows.Err()
}

func (r *MultiverseRepository) List(ctx context.Context) ([]cosmology.Multiverse, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT name FROM multiverses ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list multiverses: %w", err)
	}
	defer rows.Close()

	var list []cosmology.Multiverse
	for rows.Next() {
		var m cosmology.Multiverse
		if err := rows.Scan(&m.Name); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for i := range list {
		ur, err := r.db.QueryContext(ctx, `SELECT universe FROM multiverse_universes WHERE multiverse = ? ORDER BY rowid`, list[i].Name)
		if err != nil {
			return nil, err
		}
		for ur.Next() {
			var u string
			if err := ur.Scan(&u); err != nil {
				ur.Close()
				return nil, err
			}
			list[i].Universes = append(list[i].Universes, cosmology.Universe{Name: u})
		}
		ur.Close()
	}
	return list, nil
}
