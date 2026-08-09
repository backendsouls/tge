package file

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"tge/internal/core/domain/powersystem"
	"tge/internal/core/port"
)

type PowerSystemRepository struct {
	basePath string
}

// NewPowerSystemRepository creates a flat-file JSON repository.
// Automatically creates the base directory if it doesn't exist.
func NewPowerSystemRepository(basePath string) *PowerSystemRepository {
	_ = os.MkdirAll(basePath, 0755)
	return &PowerSystemRepository{
		basePath: basePath,
	}
}

func (r *PowerSystemRepository) filename(name string) string {
	cleanName := strings.ToLower(strings.ReplaceAll(name, " ", "_"))
	cleanName = filepath.Base(filepath.Clean(cleanName))
	return filepath.Join(r.basePath, cleanName+".json")
}

func (r *PowerSystemRepository) Save(ctx context.Context, system powersystem.PowerSystem) error {
	data, err := json.MarshalIndent(system, "", "  ")
	if err != nil {
		return err
	}
	targetPath := r.filename(system.Name)
	tempPath := targetPath + ".tmp"
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return err
	}
	return os.Rename(tempPath, targetPath)
}

func (r *PowerSystemRepository) FindByName(ctx context.Context, name string) (powersystem.PowerSystem, error) {
	data, err := os.ReadFile(r.filename(name))
	if err != nil {
		if os.IsNotExist(err) {
			return powersystem.PowerSystem{}, port.ErrPowerSystemNotFound
		}
		return powersystem.PowerSystem{}, err
	}

	var system powersystem.PowerSystem
	if err := json.Unmarshal(data, &system); err != nil {
		return powersystem.PowerSystem{}, err
	}

	// Unmarshal creates nil maps if they were empty in JSON. Initialize them for safety.
	if system.Nodes == nil {
		system.Nodes = make(map[string]*powersystem.PowerNode)
	}

	return system, nil
}

func (r *PowerSystemRepository) List(ctx context.Context) ([]powersystem.PowerSystem, error) {
	entries, err := os.ReadDir(r.basePath)
	if err != nil {
		return nil, err
	}

	var systems []powersystem.PowerSystem
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(r.basePath, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		var system powersystem.PowerSystem
		if err := json.Unmarshal(data, &system); err != nil {
			return nil, err
		}

		if system.Nodes == nil {
			system.Nodes = make(map[string]*powersystem.PowerNode)
		}

		systems = append(systems, system)
	}
	return systems, nil
}
