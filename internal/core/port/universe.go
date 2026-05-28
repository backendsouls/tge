package port

import (
	"context"
	"errors"
	"tge/internal/core/domain/cosmology"
)

var (
	// ErrUniverseNotFound is returned when a lookup finds no matching universe.
	ErrUniverseNotFound = errors.New("universe: not found")
	// ErrUniverseExists is returned when creating a universe whose name is taken.
	ErrUniverseExists = errors.New("universe: already exists")
	// ErrPowerSystemTaken is returned when a power system already belongs to another universe.
	ErrPowerSystemTaken = errors.New("universe: power system already belongs to another universe")
)

// AddUniverseSystemInput describes a power system to add to a universe.
type AddUniverseSystemInput struct {
	Universe string
	System   string
}

// AddRealmsInput describes one or more in-universe realms (locations) to add.
type AddRealmsInput struct {
	Universe string
	Names    []string
}

// UniverseRepository is a driven port persisting universes and their membership.
type UniverseRepository interface {
	Create(ctx context.Context, name string) error
	FindByName(ctx context.Context, name string) (cosmology.Universe, error)
	List(ctx context.Context) ([]cosmology.Universe, error)
	// SaveSystems replaces the stored member systems of a universe.
	SaveSystems(ctx context.Context, u cosmology.Universe) error
	// SaveRealms replaces the stored realms (locations) of a universe.
	SaveRealms(ctx context.Context, u cosmology.Universe) error
	// FindBySystem returns the universe a power system belongs to, or ErrUniverseNotFound.
	FindBySystem(ctx context.Context, system string) (cosmology.Universe, error)
}

// UniverseService is a driving port for universe use cases.
type UniverseService interface {
	CreateUniverse(ctx context.Context, name string) (cosmology.Universe, error)
	AddSystem(ctx context.Context, in AddUniverseSystemInput) (cosmology.Universe, error)
	AddRealms(ctx context.Context, in AddRealmsInput) (cosmology.Universe, error)
	GetUniverse(ctx context.Context, name string) (cosmology.Universe, error)
	ListUniverses(ctx context.Context) ([]cosmology.Universe, error)
}
