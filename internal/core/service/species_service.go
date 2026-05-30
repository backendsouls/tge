package service

import (
	"context"

	"tge/internal/core/domain/character"
	"tge/internal/core/port"
)

type SpeciesService struct {
	repo port.SpeciesRepository
}

func NewSpeciesService(repo port.SpeciesRepository) *SpeciesService {
	return &SpeciesService{repo: repo}
}

func (s *SpeciesService) CreateSpecies(ctx context.Context, in port.CreateSpeciesInput) (character.Species, error) {
	sp, err := character.NewSpecies(in.Name, in.Power, in.Lifespan, character.Gender(in.DefaultGender))
	if err != nil {
		return character.Species{}, err
	}
	if err := s.repo.Save(ctx, sp); err != nil {
		return character.Species{}, err
	}
	return sp, nil
}

func (s *SpeciesService) Species(ctx context.Context, name string) (character.Species, error) {
	return s.repo.FindByName(ctx, name)
}

func (s *SpeciesService) ListSpecies(ctx context.Context) ([]character.Species, error) {
	return s.repo.List(ctx)
}
