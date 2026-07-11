package cosmology

import (
	"errors"
	"reflect"
	"testing"
	"tge/internal/core/domain/powersystem"
)

func TestNewUniverse(t *testing.T) {
	t.Run("creates a named universe", func(t *testing.T) {
		u, err := NewUniverse("  Universe A  ")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if u.Name != "Universe A" {
			t.Errorf("Name = %q, want %q", u.Name, "Universe A")
		}
		if len(u.Systems) != 0 {
			t.Errorf("new universe should have no systems, got %v", u.Systems)
		}
	})

	t.Run("rejects a blank name", func(t *testing.T) {
		_, err := NewUniverse("  ")
		if !errors.Is(err, ErrInvalidUniverseName) {
			t.Fatalf("err = %v, want %v", err, ErrInvalidUniverseName)
		}
	})
}

func TestUniverse_AddSystem(t *testing.T) {
	t.Run("adds power systems by name", func(t *testing.T) {
		u, _ := NewUniverse("Universe A")
		if err := u.AddSystem("powersystem.Cultivation"); err != nil {
			t.Fatalf("add: %v", err)
		}
		if err := u.AddSystem("Sorcery"); err != nil {
			t.Fatalf("add: %v", err)
		}
		got := []string{u.Systems[0].Name, u.Systems[1].Name}
		if !reflect.DeepEqual(got, []string{"powersystem.Cultivation", "Sorcery"}) {
			t.Errorf("systems = %v, want [powersystem.Cultivation Sorcery]", got)
		}
	})

	t.Run("rejects a duplicate within the universe", func(t *testing.T) {
		u, _ := NewUniverse("Universe A")
		_ = u.AddSystem("powersystem.Cultivation")
		if err := u.AddSystem("powersystem.Cultivation"); !errors.Is(err, ErrUniverseSystemExists) {
			t.Fatalf("err = %v, want %v", err, ErrUniverseSystemExists)
		}
	})

	t.Run("rejects a blank system name", func(t *testing.T) {
		u, _ := NewUniverse("Universe A")
		if err := u.AddSystem("  "); !errors.Is(err, powersystem.ErrInvalidSystemName) {
			t.Fatalf("err = %v, want %v", err, powersystem.ErrInvalidSystemName)
		}
	})
}

func TestUniverse_AddRealms(t *testing.T) {
	t.Run("adds two or more realms", func(t *testing.T) {
		u, _ := NewUniverse("Universe A")
		if err := u.AddRealms("Hell cultivation.Realm", "character.Mortal cultivation.Realm", "Heaven cultivation.Realm"); err != nil {
			t.Fatalf("add realms: %v", err)
		}
		got := []string{u.Realms[0].Name, u.Realms[1].Name, u.Realms[2].Name}
		want := []string{"Hell cultivation.Realm", "character.Mortal cultivation.Realm", "Heaven cultivation.Realm"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("realms = %v, want %v", got, want)
		}
	})

	t.Run("allows a single bubble realm", func(t *testing.T) {
		u, _ := NewUniverse("Universe A")
		if err := u.AddRealms("Pocket cultivation.Realm"); err != nil {
			t.Fatalf("add single realm: %v", err)
		}
		if len(u.Realms) != 1 || u.Realms[0].Name != "Pocket cultivation.Realm" {
			t.Errorf("realms = %+v, want one Pocket cultivation.Realm", u.Realms)
		}
	})

	t.Run("allows adding more incrementally", func(t *testing.T) {
		u, _ := NewUniverse("Universe A")
		_ = u.AddRealms("Hell cultivation.Realm", "character.Mortal cultivation.Realm")
		if err := u.AddRealms("Heaven cultivation.Realm"); err != nil {
			t.Fatalf("add third realm: %v", err)
		}
		if len(u.Realms) != 3 {
			t.Errorf("want 3 realms, got %d", len(u.Realms))
		}
	})

	t.Run("rejects duplicates", func(t *testing.T) {
		u, _ := NewUniverse("Universe A")
		_ = u.AddRealms("Hell cultivation.Realm", "character.Mortal cultivation.Realm")
		if err := u.AddRealms("Hell cultivation.Realm"); !errors.Is(err, ErrRealmExistsInUniverse) {
			t.Fatalf("err = %v, want %v", err, ErrRealmExistsInUniverse)
		}
		// Also rejects duplicates within the same call.
		u2, _ := NewUniverse("Universe B")
		if err := u2.AddRealms("Hell cultivation.Realm", "Hell cultivation.Realm"); !errors.Is(err, ErrRealmExistsInUniverse) {
			t.Fatalf("err = %v, want %v", err, ErrRealmExistsInUniverse)
		}
	})
}
