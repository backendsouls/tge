package service

import (
	"context"

	"tge/internal/core/domain/progression"
	"tge/internal/core/port"
)

// PowerSystemService implements port.PowerSystemService, building power trees on
// top of a PowerSystemRepository.
type PowerSystemService struct {
	repo port.PowerSystemRepository
}

// NewPowerSystemService wires the service to its repository.
func NewPowerSystemService(repo port.PowerSystemRepository) *PowerSystemService {
	return &PowerSystemService{repo: repo}
}

// CreateSystem validates the name and persists a new empty system.
func (s *PowerSystemService) CreateSystem(ctx context.Context, name string, kind progression.PowerSystemType) (progression.PowerSystem, error) {
	ps, err := progression.NewPowerSystem(name, kind)
	if err != nil {
		return progression.PowerSystem{}, err
	}
	if err := s.repo.Create(ctx, ps.Name); err != nil {
		return progression.PowerSystem{}, err
	}
	return ps, nil
}

// AddPower loads the system, adds the power to its tree and persists it.
func (s *PowerSystemService) AddPower(ctx context.Context, in port.AddPowerInput) (progression.PowerSystem, error) {
	ps, err := s.repo.FindByName(ctx, in.System)
	if err != nil {
		return progression.PowerSystem{}, err
	}
	if err := ps.AddPower(in.Name, in.Parent); err != nil {
		return progression.PowerSystem{}, err
	}
	if err := s.repo.SavePowers(ctx, ps); err != nil {
		return progression.PowerSystem{}, err
	}
	return ps, nil
}

// GetSystem returns a single system with its power tree.
func (s *PowerSystemService) GetSystem(ctx context.Context, name string) (progression.PowerSystem, error) {
	return s.repo.FindByName(ctx, name)
}

// ListSystems returns all systems.
func (s *PowerSystemService) ListSystems(ctx context.Context) ([]progression.PowerSystem, error) {
	return s.repo.List(ctx)
}
