package file_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"tge/internal/adapter/file"
	"tge/internal/core/domain/character"
	"tge/internal/core/port"
)

func TestCharacterRepository_SaveAndFind(t *testing.T) {
	tempDir := t.TempDir()
	repo := file.NewCharacterRepository(tempDir)
	ctx := context.Background()

	cfg := character.CharacterConfig{
		Name:    "Klein",
		Type:    character.MainCharacter,
		Gender:  character.Male,
		Species: character.Species{Name: "Human"},
	}
	char, _ := character.NewMortalCharacter(cfg)

	// Add some mechanics to test deep serialization
	char.MechanicState.Tier = 3
	char.UnlockedNodes = append(char.UnlockedNodes, character.NodeProgress{
		System:    "Lord of the Mysteries",
		NodeID:    "seer",
		Level:     9,
		BasePower: 100.0,
	})

	err := repo.Save(ctx, char)
	if err != nil {
		t.Fatalf("expected no error saving character, got %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(filepath.Join(tempDir, "klein.json")); os.IsNotExist(err) {
		t.Fatalf("expected klein.json to exist on disk")
	}

	loadedChar, err := repo.FindByName(ctx, "Klein")
	if err != nil {
		t.Fatalf("expected no error loading character, got %v", err)
	}

	if loadedChar.MechanicState.Tier != 3 {
		t.Errorf("expected loaded tier 3, got %d", loadedChar.MechanicState.Tier)
	}

	if len(loadedChar.UnlockedNodes) != 1 || loadedChar.UnlockedNodes[0].NodeID != "seer" {
		t.Errorf("expected unlocked node 'seer' to be loaded properly")
	}
}

func TestCharacterRepository_NotFound(t *testing.T) {
	tempDir := t.TempDir()
	repo := file.NewCharacterRepository(tempDir)
	ctx := context.Background()

	_, err := repo.FindByName(ctx, "Amon")
	if err != port.ErrCharacterNotFound {
		t.Errorf("expected ErrCharacterNotFound, got %v", err)
	}
}
