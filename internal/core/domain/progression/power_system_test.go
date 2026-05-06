package progression

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewPowerSystem(t *testing.T) {
	t.Run("creates a named system", func(t *testing.T) {
		ps, err := NewPowerSystem("cosmology.Universe A Cultivation")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ps.Name != "cosmology.Universe A Cultivation" {
			t.Errorf("Name = %q, want %q", ps.Name, "cosmology.Universe A Cultivation")
		}
		if len(ps.Powers) != 0 {
			t.Errorf("new system should have no powers, got %v", ps.Powers)
		}
	})

	t.Run("rejects a blank name", func(t *testing.T) {
		_, err := NewPowerSystem("  ")
		if !errors.Is(err, ErrInvalidSystemName) {
			t.Fatalf("err = %v, want %v", err, ErrInvalidSystemName)
		}
	})
}

func TestPowerSystem_AddPower(t *testing.T) {
	t.Run("adds root and nested powers forming a tree", func(t *testing.T) {
		ps, _ := NewPowerSystem("cosmology.Universe A Cultivation")
		if err := ps.AddPower("Body", ""); err != nil {
			t.Fatalf("add Body: %v", err)
		}
		if err := ps.AddPower("Iron Skin", "Body"); err != nil {
			t.Fatalf("add Iron Skin: %v", err)
		}
		if err := ps.AddPower("Soul", ""); err != nil {
			t.Fatalf("add Soul: %v", err)
		}

		want := []string{"Body", "Iron Skin", "Soul"}
		if got := ps.Names(); !reflect.DeepEqual(got, want) {
			t.Errorf("Names() = %v, want %v", got, want)
		}
		if len(ps.Powers) != 2 {
			t.Fatalf("want 2 root powers, got %d", len(ps.Powers))
		}
		if len(ps.Powers[0].Children) != 1 || ps.Powers[0].Children[0].Name != "Iron Skin" {
			t.Errorf("Iron Skin not nested under Body: %+v", ps.Powers[0])
		}
	})

	t.Run("rejects a duplicate power name anywhere in the system", func(t *testing.T) {
		ps, _ := NewPowerSystem("S")
		_ = ps.AddPower("Body", "")
		err := ps.AddPower("Body", "")
		if !errors.Is(err, ErrPowerExists) {
			t.Fatalf("err = %v, want %v", err, ErrPowerExists)
		}
	})

	t.Run("rejects an unknown parent", func(t *testing.T) {
		ps, _ := NewPowerSystem("S")
		err := ps.AddPower("Iron Skin", "Body")
		if !errors.Is(err, ErrPowerParentNotFound) {
			t.Fatalf("err = %v, want %v", err, ErrPowerParentNotFound)
		}
	})

	t.Run("rejects a blank power name", func(t *testing.T) {
		ps, _ := NewPowerSystem("S")
		err := ps.AddPower("  ", "")
		if !errors.Is(err, ErrInvalidPowerName) {
			t.Fatalf("err = %v, want %v", err, ErrInvalidPowerName)
		}
	})
}
