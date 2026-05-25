package port

import (
	"context"
	"errors"
	"tge/internal/core/domain/character"
)

var (
	ErrSpeciesNotFound = errors.New("species: not found")
	ErrSpeciesExists   = errors.New("species: already exists")
)

type CreateSpeciesInput struct {
	Name          string
	Power         float64
	Lifespan      int
	DefaultGender string
}

type SpeciesRepository interface {
	Save(ctx context.Context, s character.Species) error
	FindByName(ctx context.Context, name string) (character.Species, error)
	List(ctx context.Context) ([]character.Species, error)
}

type SpeciesService interface {
	CreateSpecies(ctx context.Context, in CreateSpeciesInput) (character.Species, error)
	Species(ctx context.Context, name string) (character.Species, error)
	ListSpecies(ctx context.Context) ([]character.Species, error)
}
