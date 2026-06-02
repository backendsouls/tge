package service

import (
	"context"

	"tge/internal/core/domain/cosmology"
	"tge/internal/core/port"
)

type MultiverseService struct {
	multiverses port.MultiverseRepository
	universes   port.UniverseRepository
	timelines   port.TimelineRepository
}

func NewMultiverseService(multiverses port.MultiverseRepository, universes port.UniverseRepository, timelines port.TimelineRepository) *MultiverseService {
	return &MultiverseService{multiverses: multiverses, universes: universes, timelines: timelines}
}

func (s *MultiverseService) CreateMultiverse(ctx context.Context, name string) (cosmology.Multiverse, error) {
	m, err := cosmology.NewMultiverse(name)
	if err != nil {
		return cosmology.Multiverse{}, err
	}
	if err := s.multiverses.Save(ctx, m); err != nil {
		return cosmology.Multiverse{}, err
	}
	if err := ensureTimeline(ctx, s.timelines, port.LocationRef{Kind: port.LocationMultiverse, Name: m.Name}); err != nil {
		return cosmology.Multiverse{}, err
	}
	return m, nil
}

func (s *MultiverseService) GetMultiverse(ctx context.Context, name string) (cosmology.Multiverse, error) {
	return s.multiverses.FindByName(ctx, name)
}

func (s *MultiverseService) ListMultiverses(ctx context.Context) ([]cosmology.Multiverse, error) {
	return s.multiverses.List(ctx)
}

func (s *MultiverseService) AddUniverse(ctx context.Context, in port.AddMultiverseUniverseInput) (cosmology.Multiverse, error) {
	m, err := s.multiverses.FindByName(ctx, in.Multiverse)
	if err != nil {
		return cosmology.Multiverse{}, err
	}
	if _, err := s.universes.FindByName(ctx, in.Universe); err != nil {
		return cosmology.Multiverse{}, err
	}
	if err := m.AddUniverse(in.Universe); err != nil {
		return cosmology.Multiverse{}, err
	}
	if err := s.multiverses.AddUniverse(ctx, in.Multiverse, in.Universe); err != nil {
		return cosmology.Multiverse{}, err
	}
	return m, nil
}
