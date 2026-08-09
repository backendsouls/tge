package file

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"tge/internal/core/domain/character"
	"tge/internal/core/port"
)

type CharacterRepository struct {
	basePath string
}

// NewCharacterRepository creates a flat-file JSON repository for characters.
func NewCharacterRepository(basePath string) *CharacterRepository {
	_ = os.MkdirAll(basePath, 0755)
	return &CharacterRepository{
		basePath: basePath,
	}
}

func (r *CharacterRepository) filename(name string) string {
	cleanName := strings.ToLower(strings.ReplaceAll(name, " ", "_"))
	cleanName = filepath.Base(filepath.Clean(cleanName))
	return filepath.Join(r.basePath, cleanName+".json")
}

func (r *CharacterRepository) Save(ctx context.Context, c character.Character) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	targetPath := r.filename(c.Name)
	tempPath := targetPath + ".tmp"
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return err
	}
	return os.Rename(tempPath, targetPath)
}

func (r *CharacterRepository) FindByName(ctx context.Context, name string) (character.Character, error) {
	data, err := os.ReadFile(r.filename(name))
	if err != nil {
		if os.IsNotExist(err) {
			return character.Character{}, port.ErrCharacterNotFound
		}
		return character.Character{}, err
	}

	var c character.Character
	if err := json.Unmarshal(data, &c); err != nil {
		return character.Character{}, err
	}

	// Ensure slices/maps are not nil
	if c.UnlockedNodes == nil {
		c.UnlockedNodes = []character.NodeProgress{}
	}
	if c.MechanicState.EnergyPools == nil {
		c.MechanicState.EnergyPools = make(map[string]int)
	}

	return c, nil
}

func (r *CharacterRepository) MainCharacters(ctx context.Context) ([]character.Character, error) {
	all, err := r.List(ctx)
	if err != nil {
		return nil, err
	}

	var mains []character.Character
	for _, c := range all {
		if c.Type == character.MainCharacter {
			mains = append(mains, c)
		}
	}
	return mains, nil
}

func (r *CharacterRepository) List(ctx context.Context) ([]character.Character, error) {
	entries, err := os.ReadDir(r.basePath)
	if err != nil {
		return nil, err
	}

	var chars []character.Character
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(r.basePath, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		var c character.Character
		if err := json.Unmarshal(data, &c); err != nil {
			return nil, err
		}

		// Ensure slices/maps are not nil
		if c.UnlockedNodes == nil {
			c.UnlockedNodes = []character.NodeProgress{}
		}
		if c.MechanicState.EnergyPools == nil {
			c.MechanicState.EnergyPools = make(map[string]int)
		}

		chars = append(chars, c)
	}
	return chars, nil
}

func (r *CharacterRepository) Clean(ctx context.Context) error {
	entries, err := os.ReadDir(r.basePath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			oldPath := filepath.Join(r.basePath, entry.Name())
			newPath := oldPath + ".deleted"
			if err := os.Rename(oldPath, newPath); err != nil {
				return err
			}
		}
	}
	return nil
}
