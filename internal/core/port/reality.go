package port

import (
	"context"
	"errors"
	"tge/internal/core/domain/cosmology"
)

var (
	ErrRealityNotFound = errors.New("reality: not found")
	ErrRealityExists   = errors.New("reality: already exists")
)

type AddRealityOmniverseInput struct {
	Reality   string
	Omniverse string
}

type RealityRepository interface {
	Save(ctx context.Context, r cosmology.Reality) error
	FindByName(ctx context.Context, name string) (cosmology.Reality, error)
	List(ctx context.Context) ([]cosmology.Reality, error)
	AddOmniverse(ctx context.Context, reality, omniverse string) error
}

type RealityService interface {
	CreateReality(ctx context.Context, name string) (cosmology.Reality, error)
	GetReality(ctx context.Context, name string) (cosmology.Reality, error)
	ListRealities(ctx context.Context) ([]cosmology.Reality, error)
	AddOmniverse(ctx context.Context, in AddRealityOmniverseInput) (cosmology.Reality, error)
}
