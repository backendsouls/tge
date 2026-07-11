package functional_test

import (
	"context"
	"testing"

	"tge/internal/adapter/file"
	"tge/internal/core/domain/character"
	"tge/internal/core/port"
	"tge/internal/core/service"
)

func TestFunctional_MechanicStateEvolution(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	charRepo := file.NewCharacterRepository(tempDir + "/characters")
	charSvc := service.NewCharacterService(
		charRepo, nil, mockSpecies{}, nil, nil, nil, mockWorld{}, service.CharacterDefaults{},
	)

	// Create Character
	charSvc.CreateCharacter(ctx, port.CreateCharacterInput{
		Name:    "Darth Vader",
		Type:    string(character.MainCharacter),
		Gender:  string(character.Male),
		Species: "Human",
	})

	// Fetch Character to modify MechanicState (simulating an external hook or combat event)
	c, _ := charRepo.FindByName(ctx, "Darth Vader")

	// Evolve State
	c.MechanicState.Tier = 5
	c.MechanicState.BasePower = 1000.0
	c.MechanicState.AddEnergyPool("The Force", 5000)
	c.MechanicState.SetAlignment(-100.0) // Dark Side

	// Recalculate power internally
	c.CalculateTotalPower()
	charRepo.Save(ctx, c)

	// Reload and Verify
	loaded, _ := charRepo.FindByName(ctx, "Darth Vader")

	if loaded.MechanicState.Tier != 5 {
		t.Errorf("expected Tier 5, got %d", loaded.MechanicState.Tier)
	}

	if loaded.MechanicState.EnergyPools["The Force"] != 5000 {
		t.Errorf("expected The Force pool 5000, got %d", loaded.MechanicState.EnergyPools["The Force"])
	}

	if loaded.MechanicState.Alignment != -100.0 {
		t.Errorf("expected alignment -100.0, got %f", loaded.MechanicState.Alignment)
	}

	if loaded.CalculateTotalPower() != 1000.0 {
		t.Errorf("expected total power 1000.0, got %f", loaded.CalculateTotalPower())
	}
}
