package cosmology_test

import (
	"testing"
	"tge/internal/core/domain/cosmology"
)

func TestNewReality(t *testing.T) {
	_, err := cosmology.NewReality("   ")
	if err != cosmology.ErrInvalidRealityName {
		t.Errorf("expected ErrInvalidRealityName, got %v", err)
	}

	r, err := cosmology.NewReality(" The Box ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Name != "The Box" {
		t.Errorf("expected name The Box, got %q", r.Name)
	}
}

func TestReality_AddOmniverse(t *testing.T) {
	r, _ := cosmology.NewReality("The Box")

	if err := r.AddOmniverse(""); err != cosmology.ErrInvalidOmniverseName {
		t.Errorf("expected ErrInvalidOmniverseName for empty omniverse name, got %v", err)
	}

	if err := r.AddOmniverse(" The All "); err != nil {
		t.Fatalf("unexpected error adding The All: %v", err)
	}

	if len(r.Omniverses) != 1 || r.Omniverses[0].Name != "The All" {
		t.Fatalf("expected The All to be added, got %v", r.Omniverses)
	}

	err := r.AddOmniverse("The All")
	if err == nil {
		t.Fatalf("expected error adding duplicate omniverse")
	}
}
