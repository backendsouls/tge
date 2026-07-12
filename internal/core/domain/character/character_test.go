package character

import (
	"errors"
	"testing"
	"tge/internal/core/domain/progression"
)

func sampleSystems(t *testing.T) []progression.PowerSystem {
	t.Helper()
	ps, err := progression.NewPowerSystem("cosmology.Universe A progression.Cultivation", "")
	if err != nil {
		t.Fatal(err)
	}
	if err := ps.AddPower("Body", ""); err != nil {
		t.Fatal(err)
	}
	return []progression.PowerSystem{ps}
}

func TestNewMortalCharacter(t *testing.T) {
	systems := sampleSystems(t)

	t.Run("creates a mortal in one or more systems with defaults", func(t *testing.T) {
		ch, err := NewMortalCharacter(CharacterConfig{
			Name:         "Lin Feng",
			Type:         MainCharacter,
			Gender:       Female,
			Species:      Species{Name: "Human", Power: 1, Lifespan: 80},
			PowerSystems: systems,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ch.Name != "Lin Feng" || ch.Type != MainCharacter || ch.Gender != Female {
			t.Errorf("got %+v", ch)
		}
		if len(ch.Species) != 1 || ch.Species[0].Name != "Human" {
			t.Errorf("Species = %+v, want [Human]", ch.Species)
		}
		if ch.PowerValue != "1" {
			t.Errorf("PowerValue = %q, want default %q", ch.PowerValue, "1")
		}
		if len(ch.PowerSystems) != 1 {
			t.Errorf("PowerSystems = %v, want 1", ch.PowerSystems)
		}
		if ch.Mortal.Age != 16 || ch.Mortal.Lifespan != 80 {
			t.Errorf("Mortal = %+v", ch.Mortal)
		}
	})

	t.Run("rejects an empty name", func(t *testing.T) {
		_, err := NewMortalCharacter(CharacterConfig{Name: " ", Type: SideCharacter, Gender: Male, Species: Species{Name: "Human", Power: 1, Lifespan: 80}, PowerSystems: systems})
		if !errors.Is(err, ErrInvalidName) {
			t.Fatalf("err = %v, want %v", err, ErrInvalidName)
		}
	})

	t.Run("rejects an invalid type", func(t *testing.T) {
		_, err := NewMortalCharacter(CharacterConfig{Name: "X", Type: "Villain", Gender: Male, Species: Species{Name: "Human", Power: 1, Lifespan: 80}, PowerSystems: systems})
		if !errors.Is(err, ErrInvalidCharacterType) {
			t.Fatalf("err = %v, want %v", err, ErrInvalidCharacterType)
		}
	})

	t.Run("rejects an invalid gender", func(t *testing.T) {
		_, err := NewMortalCharacter(CharacterConfig{Name: "X", Type: SideCharacter, Gender: "Other", Species: Species{Name: "Human", Power: 1, Lifespan: 80}, PowerSystems: systems})
		if !errors.Is(err, ErrInvalidGender) {
			t.Fatalf("err = %v, want %v", err, ErrInvalidGender)
		}
	})

	t.Run("requires at least one system", func(t *testing.T) {
		_, err := NewMortalCharacter(CharacterConfig{Name: "X", Type: SideCharacter, Gender: Male, Species: Species{Name: "Human", Power: 1, Lifespan: 80}})
		if !errors.Is(err, ErrMissingSystem) {
			t.Fatalf("err = %v, want %v", err, ErrMissingSystem)
		}
	})
}

func TestCheckRole(t *testing.T) {
	female := []Character{{Name: "Lin Feng", Type: MainCharacter, Gender: Female}}
	male := []Character{{Name: "Mu Chen", Type: MainCharacter, Gender: Male}}
	both := append(append([]Character{}, female...), male...)

	t.Run("non-lead roles always pass", func(t *testing.T) {
		for _, ty := range []CharacterType{MainCharacter, SideCharacter, SupportCharacter} {
			if err := CheckRole(ty, nil); err != nil {
				t.Errorf("CheckRole(%s, nil) = %v, want nil", ty, err)
			}
		}
	})

	t.Run("Hero requires some female main character", func(t *testing.T) {
		if err := CheckRole(Hero, female); err != nil {
			t.Errorf("Hero with female MC: %v", err)
		}
		if err := CheckRole(Hero, both); err != nil {
			t.Errorf("Hero with a female MC present: %v", err)
		}
		if err := CheckRole(Hero, male); !errors.Is(err, ErrRoleGenderMismatch) {
			t.Errorf("Hero with only male MCs: err = %v, want %v", err, ErrRoleGenderMismatch)
		}
	})

	t.Run("Heroine requires some male main character", func(t *testing.T) {
		if err := CheckRole(Heroine, male); err != nil {
			t.Errorf("Heroine with male MC: %v", err)
		}
		if err := CheckRole(Heroine, female); !errors.Is(err, ErrRoleGenderMismatch) {
			t.Errorf("Heroine with only female MCs: err = %v, want %v", err, ErrRoleGenderMismatch)
		}
	})

	t.Run("lead roles require a main character to exist", func(t *testing.T) {
		if err := CheckRole(Hero, nil); !errors.Is(err, ErrMainCharacterRequired) {
			t.Errorf("Hero with no MC: err = %v, want %v", err, ErrMainCharacterRequired)
		}
	})
}
