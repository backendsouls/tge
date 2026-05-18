package cosmology_test

import (
	"testing"
	"tge/internal/core/domain/cosmology"
)

func TestNewOmniverse(t *testing.T) {
	_, err := cosmology.NewOmniverse("   ")
	if err != cosmology.ErrInvalidOmniverseName {
		t.Errorf("expected ErrInvalidOmniverseName, got %v", err)
	}

	o, err := cosmology.NewOmniverse(" The All ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if o.Name != "The All" {
		t.Errorf("expected name The All, got %q", o.Name)
	}
}

func TestOmniverse_AddMultiverse(t *testing.T) {
	o, _ := cosmology.NewOmniverse("The All")

	if err := o.AddMultiverse(""); err != cosmology.ErrInvalidMultiverseName {
		t.Errorf("expected ErrInvalidMultiverseName for empty multiverse name, got %v", err)
	}

	if err := o.AddMultiverse(" Marvel "); err != nil {
		t.Fatalf("unexpected error adding Marvel: %v", err)
	}

	if len(o.Multiverses) != 1 || o.Multiverses[0].Name != "Marvel" {
		t.Fatalf("expected Marvel to be added, got %v", o.Multiverses)
	}

	err := o.AddMultiverse("Marvel")
	if err == nil {
		t.Fatalf("expected error adding duplicate multiverse")
	}
}
