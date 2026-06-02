package service

import (
	"context"
	"errors"

	"tge/internal/core/domain/cosmology"
	"tge/internal/core/port"
)

// UniverseService implements port.UniverseService. It groups existing power
// systems under universes, keeping a power system in at most one universe.
type UniverseService struct {
	universes port.UniverseRepository
	systems   port.PowerSystemRepository
	timelines port.TimelineRepository
}

// NewUniverseService wires the service to its repositories (all interfaces).
func NewUniverseService(universes port.UniverseRepository, systems port.PowerSystemRepository, timelines port.TimelineRepository) *UniverseService {
	return &UniverseService{universes: universes, systems: systems, timelines: timelines}
}

// CreateUniverse validates the name and persists a new empty universe.
func (s *UniverseService) CreateUniverse(ctx context.Context, name string) (cosmology.Universe, error) {
	u, err := cosmology.NewUniverse(name)
	if err != nil {
		return cosmology.Universe{}, err
	}
	if err := s.universes.Create(ctx, u.Name); err != nil {
		return cosmology.Universe{}, err
	}
	if err := ensureTimeline(ctx, s.timelines, port.LocationRef{Kind: port.LocationUniverse, Name: u.Name}); err != nil {
		return cosmology.Universe{}, err
	}
	return u, nil
}

// AddSystem adds an existing, unclaimed power system to a universe.
func (s *UniverseService) AddSystem(ctx context.Context, in port.AddUniverseSystemInput) (cosmology.Universe, error) {
	u, err := s.universes.FindByName(ctx, in.Universe)
	if err != nil {
		return cosmology.Universe{}, err
	}
	if _, err := s.systems.FindByName(ctx, in.System); err != nil {
		return cosmology.Universe{}, err
	}

	switch owner, err := s.universes.FindBySystem(ctx, in.System); {
	case err == nil:
		if owner.Name != u.Name {
			return cosmology.Universe{}, port.ErrPowerSystemTaken
		}
	case errors.Is(err, port.ErrUniverseNotFound):
		// unclaimed
	default:
		return cosmology.Universe{}, err
	}

	if err := u.AddSystem(in.System); err != nil {
		return cosmology.Universe{}, err
	}
	if err := s.universes.SaveSystems(ctx, u); err != nil {
		return cosmology.Universe{}, err
	}
	return u, nil
}

// AddRealms adds one or more in-universe realms (locations) to a universe.
func (s *UniverseService) AddRealms(ctx context.Context, in port.AddRealmsInput) (cosmology.Universe, error) {
	u, err := s.universes.FindByName(ctx, in.Universe)
	if err != nil {
		return cosmology.Universe{}, err
	}
	if err := u.AddRealms(in.Names...); err != nil {
		return cosmology.Universe{}, err
	}
	if err := s.universes.SaveRealms(ctx, u); err != nil {
		return cosmology.Universe{}, err
	}
	// Each realm is itself a location, so it owns a timeline too.
	for _, name := range in.Names {
		owner := port.LocationRef{Kind: port.LocationRealm, Name: name, Universe: u.Name}
		if err := ensureTimeline(ctx, s.timelines, owner); err != nil {
			return cosmology.Universe{}, err
		}
	}
	return u, nil
}

// GetUniverse returns a single universe with its member systems and realms.
func (s *UniverseService) GetUniverse(ctx context.Context, name string) (cosmology.Universe, error) {
	return s.universes.FindByName(ctx, name)
}

// ListUniverses returns all universes.
func (s *UniverseService) ListUniverses(ctx context.Context) ([]cosmology.Universe, error) {
	return s.universes.List(ctx)
}
