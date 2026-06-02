package service

import (
	"context"

	"tge/internal/core/domain/cosmology"
	"tge/internal/core/port"
)

type RealityService struct {
	realities  port.RealityRepository
	omniverses port.OmniverseRepository
	timelines  port.TimelineRepository
}

func NewRealityService(realities port.RealityRepository, omniverses port.OmniverseRepository, timelines port.TimelineRepository) *RealityService {
	return &RealityService{realities: realities, omniverses: omniverses, timelines: timelines}
}

func (s *RealityService) CreateReality(ctx context.Context, name string) (cosmology.Reality, error) {
	r, err := cosmology.NewReality(name)
	if err != nil {
		return cosmology.Reality{}, err
	}
	if err := s.realities.Save(ctx, r); err != nil {
		return cosmology.Reality{}, err
	}
	if err := ensureTimeline(ctx, s.timelines, port.LocationRef{Kind: port.LocationBox, Name: r.Name}); err != nil {
		return cosmology.Reality{}, err
	}
	return r, nil
}

func (s *RealityService) GetReality(ctx context.Context, name string) (cosmology.Reality, error) {
	return s.realities.FindByName(ctx, name)
}

func (s *RealityService) ListRealities(ctx context.Context) ([]cosmology.Reality, error) {
	return s.realities.List(ctx)
}

func (s *RealityService) AddOmniverse(ctx context.Context, in port.AddRealityOmniverseInput) (cosmology.Reality, error) {
	r, err := s.realities.FindByName(ctx, in.Reality)
	if err != nil {
		return cosmology.Reality{}, err
	}
	if _, err := s.omniverses.FindByName(ctx, in.Omniverse); err != nil {
		return cosmology.Reality{}, err
	}
	if err := r.AddOmniverse(in.Omniverse); err != nil {
		return cosmology.Reality{}, err
	}
	if err := s.realities.AddOmniverse(ctx, in.Reality, in.Omniverse); err != nil {
		return cosmology.Reality{}, err
	}
	return r, nil
}
