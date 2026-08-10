package character_test

import (
	"testing"
	"tge/internal/core/domain/character"
	"tge/internal/core/domain/power"
)

func TestCharacter_CalculateTotalPower(t *testing.T) {
	cfg := character.CharacterConfig{
		Name:    "Shirou",
		Type:    character.MainCharacter,
		Gender:  character.Male,
		Species: character.Species{Name: "Human", Power: 0.65},
	}
	char, err := character.NewMortalCharacter(cfg)
	if err != nil {
		t.Fatalf("expected no error creating character, got %v", err)
	}

	// Initialize mechanic state
	state, _ := power.NewMechanicState(1, 100.0)
	char.MechanicState = state

	// Add an unlocked node that grants +50 base power
	char.UnlockedNodes = append(char.UnlockedNodes, character.NodeProgress{
		System:    "Nasuverse",
		NodeID:    "magic_circuits",
		Level:     1,
		BasePower: 50.0,
	})

	// Expected power: MechanicState.BasePower (100) + UnlockedNodes BasePower (50) = 150
	totalPower := char.CalculateTotalPower()
	if totalPower != 97.5 {
		t.Errorf("expected total power to be 97.5, got %f", totalPower)
	}
}

func TestCharacter_CalculateTotalPower_WithSpeciesMultiplier(t *testing.T) {
	cfg := character.CharacterConfig{
		Name:    "Goku",
		Type:    character.MainCharacter,
		Gender:  character.Male,
		Species: character.Species{Name: "Saiyan", Power: 10.0}, // 10x multiplier
	}
	char, _ := character.NewMortalCharacter(cfg)

	state, _ := power.NewMechanicState(1, 100.0)
	char.MechanicState = state

	// Expected power: (MechanicState.BasePower (100) + Nodes (0)) * Species.Power (10) = 1000
	totalPower := char.CalculateTotalPower()
	if totalPower != 1000.0 {
		t.Errorf("expected total power to be 1000.0, got %f", totalPower)
	}
}
