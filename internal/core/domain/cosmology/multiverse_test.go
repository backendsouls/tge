package cosmology_test

import (
	"testing"
	"tge/internal/core/domain/cosmology"
)

func TestNewMultiverse(t *testing.T) {
	_, err := cosmology.NewMultiverse("   ")
	if err != cosmology.ErrInvalidMultiverseName {
		t.Errorf("expected ErrInvalidMultiverseName, got %v", err)
	}

	m, err := cosmology.NewMultiverse(" Marvel ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Name != "Marvel" {
		t.Errorf("expected name Marvel, got %q", m.Name)
	}
}

func TestMultiverse_AddUniverse(t *testing.T) {
	m, _ := cosmology.NewMultiverse("Marvel")

	if err := m.AddUniverse(""); err != cosmology.ErrInvalidUniverseName {
		t.Errorf("expected ErrInvalidUniverseName for empty universe name, got %v", err)
	}

	if err := m.AddUniverse(" MCU "); err != nil {
		t.Fatalf("unexpected error adding MCU: %v", err)
	}

	if len(m.Universes) != 1 || m.Universes[0].Name != "MCU" {
		t.Fatalf("expected MCU to be added, got %v", m.Universes)
	}

	err := m.AddUniverse("MCU")
	if err == nil {
		t.Fatalf("expected error adding duplicate universe")
	}
}
