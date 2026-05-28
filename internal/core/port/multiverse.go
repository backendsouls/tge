package port

import (
	"context"
	"errors"
	"tge/internal/core/domain/cosmology"
)

var (
	ErrMultiverseNotFound = errors.New("multiverse: not found")
	ErrMultiverseExists   = errors.New("multiverse: already exists")
)

type AddMultiverseUniverseInput struct {
	Multiverse string
	Universe   string
}

type MultiverseRepository interface {
	Save(ctx context.Context, m cosmology.Multiverse) error
	FindByName(ctx context.Context, name string) (cosmology.Multiverse, error)
	List(ctx context.Context) ([]cosmology.Multiverse, error)
	AddUniverse(ctx context.Context, multiverse, universe string) error
}

type MultiverseService interface {
	CreateMultiverse(ctx context.Context, name string) (cosmology.Multiverse, error)
	GetMultiverse(ctx context.Context, name string) (cosmology.Multiverse, error)
	ListMultiverses(ctx context.Context) ([]cosmology.Multiverse, error)
	AddUniverse(ctx context.Context, in AddMultiverseUniverseInput) (cosmology.Multiverse, error)
}
