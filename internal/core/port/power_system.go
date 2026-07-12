package port

import (
	"context"
	"errors"
	"tge/internal/core/domain/progression"
)

var (
	// ErrPowerSystemNotFound is returned when a lookup finds no matching system.
	ErrPowerSystemNotFound = errors.New("power system: not found")
	// ErrPowerSystemExists is returned when creating a system whose name is taken.
	ErrPowerSystemExists = errors.New("power system: already exists")
)

// AddPowerInput describes a power to add to an existing system. An empty Parent
// makes the power a root.
type AddPowerInput struct {
	System string
	Name   string
	Parent string
}

// PowerSystemRepository is a driven port persisting power systems and their
// power trees.
type PowerSystemRepository interface {
	Create(ctx context.Context, name string) error
	FindByName(ctx context.Context, name string) (progression.PowerSystem, error)
	List(ctx context.Context) ([]progression.PowerSystem, error)
	// SavePowers replaces the stored powers of a system with the given tree.
	SavePowers(ctx context.Context, system progression.PowerSystem) error
}

// PowerSystemService is a driving port for power-system use cases.
type PowerSystemService interface {
	CreateSystem(ctx context.Context, name string, kind progression.PowerSystemType) (progression.PowerSystem, error)
	AddPower(ctx context.Context, in AddPowerInput) (progression.PowerSystem, error)
	GetSystem(ctx context.Context, name string) (progression.PowerSystem, error)
	ListSystems(ctx context.Context) ([]progression.PowerSystem, error)
}
