package service

import (
	"context"
	"errors"
	"fmt"

	"tge/internal/core/domain/character"
	"tge/internal/core/port"
)

type IdleService struct {
	chars port.CharacterRepository
}

func NewIdleService(chars port.CharacterRepository) *IdleService {
	return &IdleService{chars: chars}
}

// CommitOfflineGains bakes finished idle slots into the base energy pools and removes them.
func (s *IdleService) CommitOfflineGains(char *character.Character) {
	if char.MechanicState.EnergyPools == nil {
		char.MechanicState.EnergyPools = make(map[string]int)
	}

	var activeSlots []character.IdleSlot
	for _, slot := range char.IdleState.Slots {
		deltaHours := float64(char.NovelTime-slot.StartTime) / 3600.0

		if slot.Duration > 0 {
			if deltaHours >= slot.Duration {
				if slot.Rate > 0 {
					remainder := char.AdvanceNode(slot.System, slot.Power, slot.Duration*slot.Rate)
					if remainder > 0 {
						poolName := fmt.Sprintf("%s_%s", slot.System, slot.Power)
						char.MechanicState.EnergyPools[poolName] += int(remainder)
					}
				}
			} else {
				if slot.Rate > 0 {
					remainder := char.AdvanceNode(slot.System, slot.Power, deltaHours*slot.Rate)
					if remainder > 0 {
						poolName := fmt.Sprintf("%s_%s", slot.System, slot.Power)
						char.MechanicState.EnergyPools[poolName] += int(remainder)
					}
				}
				slot.Duration -= deltaHours
				slot.StartTime = char.NovelTime
				activeSlots = append(activeSlots, slot)
			}
		} else {
			// Indefinite duration
			if slot.Rate > 0 {
				remainder := char.AdvanceNode(slot.System, slot.Power, deltaHours*slot.Rate)
				if remainder > 0 {
					poolName := fmt.Sprintf("%s_%s", slot.System, slot.Power)
					char.MechanicState.EnergyPools[poolName] += int(remainder)
				}
			}
			slot.StartTime = char.NovelTime
			activeSlots = append(activeSlots, slot)
		}
	}
	char.IdleState.Slots = activeSlots
}

// AssignActivity assigns an idle activity to a character if a slot is available.
func (s *IdleService) AssignActivity(ctx context.Context, charName string, systemName string, powerName string, duration float64) (character.Character, error) {
	char, err := s.chars.FindByName(ctx, charName)
	if err != nil {
		return character.Character{}, err
	}

	s.CommitOfflineGains(&char)

	if systemName == "none" || systemName == "" {
		char.IdleState.Slots = []character.IdleSlot{}
		if err := s.chars.Save(ctx, char); err != nil {
			return character.Character{}, err
		}
		return char, nil
	}

	totalSlots := char.IdleState.TotalSlots
	if totalSlots <= 0 {
		totalSlots = 1
	}

	if len(char.IdleState.Slots) >= totalSlots {
		return character.Character{}, errors.New("no available idle slots")
	}

	slot := character.IdleSlot{
		StartTime: char.NovelTime,
		Duration:  duration,
		System:    systemName,
		Power:     powerName,
		Rate:      10.0, // Default 10 XP/hr
	}

	char.IdleState.Slots = append(char.IdleState.Slots, slot)

	if err := s.chars.Save(ctx, char); err != nil {
		return character.Character{}, err
	}

	return char, nil
}
