package file_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"tge/internal/adapter/file"
	"tge/internal/core/domain/powersystem"
	"tge/internal/core/port"
)

func TestPowerSystemRepository_SaveAndFind(t *testing.T) {
	tempDir := t.TempDir()
	repo := file.NewPowerSystemRepository(tempDir)
	ctx := context.Background()

	// Create a DAG system
	sys, _ := powersystem.NewPowerSystem("Nasuverse", powersystem.Magic)
	nodeA, _ := powersystem.NewPowerNode("Circuits", "Base", nil)
	_ = sys.AddNode(&nodeA)

	// Save it to disk
	err := repo.Save(ctx, sys)
	if err != nil {
		t.Fatalf("expected no error saving system, got %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(filepath.Join(tempDir, "nasuverse.json")); os.IsNotExist(err) {
		t.Fatalf("expected nasuverse.json to exist on disk")
	}

	// Read it back
	loadedSys, err := repo.FindByName(ctx, "Nasuverse")
	if err != nil {
		t.Fatalf("expected no error reading system, got %v", err)
	}

	if loadedSys.Name != "Nasuverse" {
		t.Errorf("expected loaded system name Nasuverse, got %s", loadedSys.Name)
	}

	if _, ok := loadedSys.Nodes["circuits"]; !ok {
		t.Errorf("expected node 'circuits' to be loaded into the map")
	}
}

func TestPowerSystemRepository_NotFound(t *testing.T) {
	tempDir := t.TempDir()
	repo := file.NewPowerSystemRepository(tempDir)
	ctx := context.Background()

	_, err := repo.FindByName(ctx, "GhostSystem")
	if err != port.ErrPowerSystemNotFound {
		t.Errorf("expected ErrPowerSystemNotFound, got %v", err)
	}
}
