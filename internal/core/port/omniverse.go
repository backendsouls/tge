package port

import (
	"context"
	"errors"
	"tge/internal/core/domain/cosmology"
)

var (
	ErrOmniverseNotFound = errors.New("omniverse: not found")
	ErrOmniverseExists   = errors.New("omniverse: already exists")
)

type AddOmniverseMultiverseInput struct {
	Omniverse  string
	Multiverse string
}

type OmniverseRepository interface {
	Save(ctx context.Context, o cosmology.Omniverse) error
	FindByName(ctx context.Context, name string) (cosmology.Omniverse, error)
	List(ctx context.Context) ([]cosmology.Omniverse, error)
	AddMultiverse(ctx context.Context, omniverse, multiverse string) error
}

type OmniverseService interface {
	CreateOmniverse(ctx context.Context, name string) (cosmology.Omniverse, error)
	GetOmniverse(ctx context.Context, name string) (cosmology.Omniverse, error)
	ListOmniverses(ctx context.Context) ([]cosmology.Omniverse, error)
	AddMultiverse(ctx context.Context, in AddOmniverseMultiverseInput) (cosmology.Omniverse, error)
}
