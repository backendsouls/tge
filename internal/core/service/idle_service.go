package service

import (
	"context"

	"tge/internal/core/domain/character"
	"tge/internal/core/port"
)

type IdleService struct {
	chars port.CharacterRepository
}

func NewIdleService(chars port.CharacterRepository) *IdleService {
	return &IdleService{chars: chars}
}

// CommitOfflineGains commits any dynamically generated offline gains to the character's base energy pools,
// and resets the StartTime to the current NovelTime.
func (s *IdleService) CommitOfflineGains(char *character.Character) {
	if char.IdleState.StartTime == 0 {
		char.IdleState.StartTime = char.NovelTime
		return
	}

	// Calculate points generated up to now and commit them to the base pool
	currentPools := char.CurrentEnergyPools(char.NovelTime)
	if char.MechanicState.EnergyPools == nil {
		char.MechanicState.EnergyPools = make(map[string]int)
	}

	for k, v := range currentPools {
		char.MechanicState.EnergyPools[k] = v
	}

	// Reset the start time for the next cycle
	char.IdleState.StartTime = char.NovelTime
}

// AssignActivity assigns an idle activity to a character, modifying their regeneration rates.
func (s *IdleService) AssignActivity(ctx context.Context, charName string, activity string) (character.Character, error) {
	char, err := s.chars.FindByName(ctx, charName)
	if err != nil {
		return character.Character{}, err
	}

	// Make sure we commit gains before switching activities
	s.CommitOfflineGains(&char)

	char.IdleState.ActiveActivity = activity

	// Reset rates based on activity
	char.IdleState.Rates = character.IdleRates{}

	switch activity {
	case "rest":
		char.IdleState.Rates.CultivationPointsPerHour = 10.0
	case "secluded_cultivation":
		char.IdleState.Rates.CultivationPointsPerHour = 100.0
	case "training_skills":
		char.IdleState.Rates.SkillPointsPerHour = 100.0
	case "studying_ability":
		char.IdleState.Rates.AbilityPointsPerHour = 100.0
	case "working_profession":
		char.IdleState.Rates.ProfessionPointsPerHour = 100.0
	case "none", "":
		// default rates
	default:
		// generic fallback or custom handling
		char.IdleState.Rates.CultivationPointsPerHour = 5.0
	}

	// StartTime is already reset by CommitOfflineGains, but we can ensure it is fresh
	char.IdleState.StartTime = char.NovelTime

	if err := s.chars.Save(ctx, char); err != nil {
		return character.Character{}, err
	}

	return char, nil
}
