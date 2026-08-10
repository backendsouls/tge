package character_test

import (
	"testing"
	"tge/internal/core/domain/character"
)

func TestNewSpecies(t *testing.T) {
	_, err := character.NewSpecies("   ", 1, 80, "")
	if err != character.ErrInvalidSpeciesName {
		t.Errorf("expected ErrInvalidSpeciesName, got %v", err)
	}

	if _, err := character.NewSpecies("Human", 0.65, 80, "Neither"); err == nil {
		t.Error("expected an error for an invalid default gender")
	}

	s, err := character.NewSpecies(" Human ", 0.65, 80, character.Male)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Name != "Human" {
		t.Errorf("expected name Human, got %q", s.Name)
	}
	if s.Power != 0.65 {
		t.Errorf("expected power 0.65, got %v", s.Power)
	}
	if s.Lifespan != 80 {
		t.Errorf("expected lifespan 80, got %v", s.Lifespan)
	}
	if s.DefaultGender != character.Male {
		t.Errorf("expected default gender Male, got %q", s.DefaultGender)
	}
}
