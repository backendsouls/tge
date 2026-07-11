package service

import (
	"context"
	"sync"

	"tge/internal/core/domain/powersystem"
	"tge/internal/core/port"
)

type PowerSystemService struct {
	repo port.PowerSystemRepository
	mu   sync.Mutex
}

func NewPowerSystemService(repo port.PowerSystemRepository) *PowerSystemService {
	return &PowerSystemService{repo: repo}
}

func (s *PowerSystemService) CreateSystem(ctx context.Context, name string, kind powersystem.PowerSystemType) (powersystem.PowerSystem, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ps, err := powersystem.NewPowerSystem(name, kind)
	if err != nil {
		return powersystem.PowerSystem{}, err
	}

	// Check if it already exists to emulate SQL UNIQUE constraint
	if _, err := s.repo.FindByName(ctx, ps.Name); err == nil {
		return powersystem.PowerSystem{}, port.ErrPowerSystemExists
	}

	if err := s.repo.Save(ctx, ps); err != nil {
		return powersystem.PowerSystem{}, err
	}
	return ps, nil
}

func (s *PowerSystemService) AddNode(ctx context.Context, in port.AddNodeInput) (powersystem.PowerSystem, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ps, err := s.repo.FindByName(ctx, in.System)
	if err != nil {
		return powersystem.PowerSystem{}, err
	}

	node, err := powersystem.NewPowerNode(in.Name, in.Category, in.Tags)
	if err != nil {
		return powersystem.PowerSystem{}, err
	}
	if in.NodeID != "" {
		node.ID = in.NodeID
	}

	if err := ps.AddNode(&node); err != nil {
		return powersystem.PowerSystem{}, err
	}

	if err := s.repo.Save(ctx, ps); err != nil {
		return powersystem.PowerSystem{}, err
	}
	return ps, nil
}

func (s *PowerSystemService) AddEdge(ctx context.Context, in port.AddEdgeInput) (powersystem.PowerSystem, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ps, err := s.repo.FindByName(ctx, in.System)
	if err != nil {
		return powersystem.PowerSystem{}, err
	}

	if err := ps.AddEdge(in.NodeID, in.TargetID, in.EdgeType); err != nil {
		return powersystem.PowerSystem{}, err
	}

	if err := s.repo.Save(ctx, ps); err != nil {
		return powersystem.PowerSystem{}, err
	}
	return ps, nil
}

func (s *PowerSystemService) GetSystem(ctx context.Context, name string) (powersystem.PowerSystem, error) {
	return s.repo.FindByName(ctx, name)
}

func (s *PowerSystemService) ListSystems(ctx context.Context) ([]powersystem.PowerSystem, error) {
	return s.repo.List(ctx)
}
