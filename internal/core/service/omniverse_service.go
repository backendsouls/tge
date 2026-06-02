package service

import (
	"context"

	"tge/internal/core/domain/cosmology"
	"tge/internal/core/port"
)

type OmniverseService struct {
	omniverses  port.OmniverseRepository
	multiverses port.MultiverseRepository
	timelines   port.TimelineRepository
}

func NewOmniverseService(omniverses port.OmniverseRepository, multiverses port.MultiverseRepository, timelines port.TimelineRepository) *OmniverseService {
	return &OmniverseService{omniverses: omniverses, multiverses: multiverses, timelines: timelines}
}

func (s *OmniverseService) CreateOmniverse(ctx context.Context, name string) (cosmology.Omniverse, error) {
	o, err := cosmology.NewOmniverse(name)
	if err != nil {
		return cosmology.Omniverse{}, err
	}
	if err := s.omniverses.Save(ctx, o); err != nil {
		return cosmology.Omniverse{}, err
	}
	if err := ensureTimeline(ctx, s.timelines, port.LocationRef{Kind: port.LocationOmniverse, Name: o.Name}); err != nil {
		return cosmology.Omniverse{}, err
	}
	return o, nil
}

func (s *OmniverseService) GetOmniverse(ctx context.Context, name string) (cosmology.Omniverse, error) {
	return s.omniverses.FindByName(ctx, name)
}

func (s *OmniverseService) ListOmniverses(ctx context.Context) ([]cosmology.Omniverse, error) {
	return s.omniverses.List(ctx)
}

func (s *OmniverseService) AddMultiverse(ctx context.Context, in port.AddOmniverseMultiverseInput) (cosmology.Omniverse, error) {
	o, err := s.omniverses.FindByName(ctx, in.Omniverse)
	if err != nil {
		return cosmology.Omniverse{}, err
	}
	if _, err := s.multiverses.FindByName(ctx, in.Multiverse); err != nil {
		return cosmology.Omniverse{}, err
	}
	if err := o.AddMultiverse(in.Multiverse); err != nil {
		return cosmology.Omniverse{}, err
	}
	if err := s.omniverses.AddMultiverse(ctx, in.Omniverse, in.Multiverse); err != nil {
		return cosmology.Omniverse{}, err
	}
	return o, nil
}
