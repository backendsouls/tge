package cultivation_test

import (
	"errors"
	"testing"

	"tge/internal/core/domain/cultivation"
)

func TestNewLevel(t *testing.T) {
	if _, err := cultivation.NewLevel(1, "  ", 0, 0); err != cultivation.ErrInvalidLevelName {
		t.Errorf("blank name: got %v", err)
	}
	if _, err := cultivation.NewLevel(0, "First", 0, 0); err == nil {
		t.Error("expected error for non-positive number")
	}
	if _, err := cultivation.NewLevel(1, "First", -1, 0); err != cultivation.ErrInvalidLevelPoints {
		t.Errorf("negative points: got %v", err)
	}
	l, err := cultivation.NewLevel(1, " First Level ", 100, 20)
	if err != nil || l.Name != "First Level" || l.Number != 1 {
		t.Fatalf("level: %+v err=%v", l, err)
	}
	if l.BreakthroughPoints != 100 || l.BottleneckPoints != 20 {
		t.Errorf("points = (%d, %d), want (100, 20)", l.BreakthroughPoints, l.BottleneckPoints)
	}
}

func TestRealm_AddLevel(t *testing.T) {
	r, err := cultivation.NewRealm(cultivation.RealmConfig{Name: "Qi Condensation"})
	if err != nil {
		t.Fatalf("new realm: %v", err)
	}
	// Insert out of order; levels should come back sorted by Number.
	if err := r.AddLevel(3, "Third Level", 0, 0); err != nil {
		t.Fatalf("add level 3: %v", err)
	}
	if err := r.AddLevel(1, "First Level", 0, 0); err != nil {
		t.Fatalf("add level 1: %v", err)
	}
	if len(r.Levels) != 2 || r.Levels[0].Number != 1 || r.Levels[1].Number != 3 {
		t.Fatalf("levels not sorted: %+v", r.Levels)
	}
	if err := r.AddLevel(1, "dup", 0, 0); !errors.Is(err, cultivation.ErrLevelNumberExists) {
		t.Fatalf("duplicate number: got %v", err)
	}
}

func TestRealm_AddLevel_MaxLevels(t *testing.T) {
	r, err := cultivation.NewRealm(cultivation.RealmConfig{Name: "Qi Condensation", MaxLevels: 2})
	if err != nil {
		t.Fatalf("new realm: %v", err)
	}
	if err := r.AddLevel(1, "First Level", 0, 0); err != nil {
		t.Fatalf("add level 1: %v", err)
	}
	if err := r.AddLevel(2, "Second Level", 0, 0); err != nil {
		t.Fatalf("add level 2: %v", err)
	}
	if err := r.AddLevel(3, "Third Level", 0, 0); !errors.Is(err, cultivation.ErrLevelNumberExceedsMax) {
		t.Fatalf("over cap: got %v, want cultivation.ErrLevelNumberExceedsMax", err)
	}
	// MaxLevels 0 means unlimited.
	unlimited, _ := cultivation.NewRealm(cultivation.RealmConfig{Name: "Boundless"})
	if err := unlimited.AddLevel(999, "Far Level", 0, 0); err != nil {
		t.Fatalf("unlimited realm rejected a high level: %v", err)
	}
}

func TestRealm_MaxLevelsFor(t *testing.T) {
	r, err := cultivation.NewRealm(cultivation.RealmConfig{
		Name: "Qi Condensation", MaxLevels: 9, MainCharacterMaxLevels: 13,
	})
	if err != nil {
		t.Fatalf("new realm: %v", err)
	}
	if got := r.MaxLevelsFor(false); got != 9 {
		t.Errorf("normal cap = %d, want 9", got)
	}
	if got := r.MaxLevelsFor(true); got != 13 {
		t.Errorf("main cap = %d, want 13", got)
	}
	// The realm can define levels up to the main-character cap (13)...
	if err := r.AddLevel(13, "Thirteenth Level", 0, 0); err != nil {
		t.Fatalf("level 13 within main cap rejected: %v", err)
	}
	// ...but not beyond it.
	if err := r.AddLevel(14, "Fourteenth Level", 0, 0); !errors.Is(err, cultivation.ErrLevelNumberExceedsMax) {
		t.Fatalf("level 14 over main cap: got %v", err)
	}

	// When the main cap is unset, the main character inherits the normal cap.
	r2, _ := cultivation.NewRealm(cultivation.RealmConfig{Name: "Foundation", MaxLevels: 9})
	if got := r2.MaxLevelsFor(true); got != 9 {
		t.Errorf("inherited main cap = %d, want 9", got)
	}
}
