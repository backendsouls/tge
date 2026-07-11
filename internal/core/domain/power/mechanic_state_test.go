package power_test

import (
	"testing"

	"tge/internal/core/domain/power"
)

func TestMechanicState_Initialization(t *testing.T) {
	state, err := power.NewMechanicState(0, 100.0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if state.Tier != 0 {
		t.Errorf("expected Tier 0, got %d", state.Tier)
	}

	if state.BasePower != 100.0 {
		t.Errorf("expected BasePower 100.0, got %f", state.BasePower)
	}

	// Tier cannot be negative
	_, err = power.NewMechanicState(-1, 50.0)
	if err == nil {
		t.Errorf("expected error when creating state with negative tier")
	}
}

func TestMechanicState_EnergyPools(t *testing.T) {
	state, _ := power.NewMechanicState(1, 10.0)
	state.AddEnergyPool("mana", 500)

	if state.EnergyPools["mana"] != 500 {
		t.Errorf("expected mana pool of 500, got %d", state.EnergyPools["mana"])
	}
}

func TestMechanicState_Alignment(t *testing.T) {
	state, _ := power.NewMechanicState(1, 10.0)
	err := state.SetAlignment(150.0)
	if err == nil {
		t.Errorf("expected error when setting alignment > 100")
	}

	err = state.SetAlignment(-150.0)
	if err == nil {
		t.Errorf("expected error when setting alignment < -100")
	}

	err = state.SetAlignment(50.0)
	if err != nil {
		t.Errorf("expected no error setting valid alignment, got %v", err)
	}
	if state.Alignment != 50.0 {
		t.Errorf("expected alignment 50, got %f", state.Alignment)
	}
}
