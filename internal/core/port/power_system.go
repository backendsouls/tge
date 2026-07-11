package port

import (
	"context"
	"errors"
	"tge/internal/core/domain/powersystem"
)

var (
	// ErrPowerSystemNotFound is returned when a lookup finds no matching system.
	ErrPowerSystemNotFound = errors.New("power system: not found")
	// ErrPowerSystemExists is returned when creating a system whose name is taken.
	ErrPowerSystemExists = errors.New("power system: already exists")
)

// AddNodeInput describes a power node to add to an existing system DAG.
type AddNodeInput struct {
	System   string
	NodeID   string
	Name     string
	Category string
	Tags     []string
}

// AddEdgeInput links two nodes in the DAG.
type AddEdgeInput struct {
	System   string
	NodeID   string
	TargetID string
	EdgeType powersystem.EdgeType
}

// PowerSystemRepository is a driven port persisting power systems and their DAGs.
type PowerSystemRepository interface {
	Save(ctx context.Context, system powersystem.PowerSystem) error
	FindByName(ctx context.Context, name string) (powersystem.PowerSystem, error)
	List(ctx context.Context) ([]powersystem.PowerSystem, error)
}

// PowerSystemService is a driving port for power-system use cases.
type PowerSystemService interface {
	CreateSystem(ctx context.Context, name string, kind powersystem.PowerSystemType) (powersystem.PowerSystem, error)
	AddNode(ctx context.Context, in AddNodeInput) (powersystem.PowerSystem, error)
	AddEdge(ctx context.Context, in AddEdgeInput) (powersystem.PowerSystem, error)
	GetSystem(ctx context.Context, name string) (powersystem.PowerSystem, error)
	ListSystems(ctx context.Context) ([]powersystem.PowerSystem, error)
}
