package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"tge/internal/config"
)

func TestLoad_Embedded(t *testing.T) {
	t.Setenv(config.EnvPath, "") // ensure no override
	d, err := config.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if d.World.Reality != "The Box" {
		t.Errorf("World.Reality = %q, want %q", d.World.Reality, "The Box")
	}
	if d.Character.Species.Name != "Human" || d.Character.Species.Lifespan != 80 {
		t.Errorf("Human base = %+v", d.Character.Species)
	}
	if d.Character.Stats.STR != 5 {
		t.Errorf("base STR = %d, want 5", d.Character.Stats.STR)
	}
	if len(d.Catalog.Realms) == 0 || len(d.Catalog.Classes) == 0 {
		t.Fatalf("catalog looks empty: %+v", d.Catalog)
	}
}

func TestLoad_Override(t *testing.T) {
	override := filepath.Join(t.TempDir(), "over.yml")
	// Only override a couple of keys; the rest must fall back to the embedded values.
	const body = `
world:
  reality: "The Omniarch"
character:
  gender: "Female"
`
	if err := os.WriteFile(override, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv(config.EnvPath, override)

	d, err := config.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if d.World.Reality != "The Omniarch" {
		t.Errorf("override World.Reality = %q, want %q", d.World.Reality, "The Omniarch")
	}
	if d.Character.Gender != "Female" {
		t.Errorf("override gender = %q, want Female", d.Character.Gender)
	}
	// Untouched keys keep the embedded defaults.
	if d.World.Universe != "Origin Universe" {
		t.Errorf("untouched World.Universe = %q, want %q", d.World.Universe, "Origin Universe")
	}
	if d.Character.Age != 16 {
		t.Errorf("untouched age = %d, want 16", d.Character.Age)
	}
}
