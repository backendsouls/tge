package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"tge/internal/core/domain/character"
	"tge/internal/core/port"

	sqlitedrv "modernc.org/sqlite"
	sqlitelib "modernc.org/sqlite/lib"
)

type SpeciesRepository struct {
	db *sql.DB
}

func NewSpeciesRepository(dsn string) (*SpeciesRepository, error) {
	db, err := open(dsn)
	if err != nil {
		return nil, err
	}
	return &SpeciesRepository{db: db}, nil
}

func (r *SpeciesRepository) Close() error {
	return r.db.Close()
}

func (r *SpeciesRepository) Save(ctx context.Context, s character.Species) error {
	const q = `INSERT INTO species (name, power, lifespan, default_gender) VALUES (?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, q, s.Name, s.Power, s.Lifespan, string(s.DefaultGender))
	if err != nil {
		if serr, ok := errors.AsType[*sqlitedrv.Error](err); ok {
			if serr.Code() == sqlitelib.SQLITE_CONSTRAINT_PRIMARYKEY || serr.Code() == sqlitelib.SQLITE_CONSTRAINT_UNIQUE {
				return port.ErrSpeciesExists
			}
		}
		return fmt.Errorf("save species: %w", err)
	}
	return nil
}

func (r *SpeciesRepository) FindByName(ctx context.Context, name string) (character.Species, error) {
	const q = `SELECT name, power, lifespan, default_gender FROM species WHERE name = ?`
	var s character.Species
	var gender string
	err := r.db.QueryRowContext(ctx, q, name).Scan(&s.Name, &s.Power, &s.Lifespan, &gender)
	if errors.Is(err, sql.ErrNoRows) {
		return character.Species{}, fmt.Errorf("%w: %q", port.ErrSpeciesNotFound, name)
	}
	if err != nil {
		return character.Species{}, fmt.Errorf("find species: %w", err)
	}
	s.DefaultGender = character.Gender(gender)
	return s, nil
}

func (r *SpeciesRepository) List(ctx context.Context) ([]character.Species, error) {
	const q = `SELECT name, power, lifespan, default_gender FROM species ORDER BY name`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list species: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var list []character.Species
	for rows.Next() {
		var s character.Species
		var gender string
		if err := rows.Scan(&s.Name, &s.Power, &s.Lifespan, &gender); err != nil {
			return nil, err
		}
		s.DefaultGender = character.Gender(gender)
		list = append(list, s)
	}
	return list, rows.Err()
}
