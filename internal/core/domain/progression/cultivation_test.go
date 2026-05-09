package progression

import (
	"errors"
	"testing"
)

func mustRealm(t *testing.T) Realm {
	t.Helper()
	r, err := NewRealm(RealmConfig{
		Name:               "Qi Condensation",
		PowerMultiplier:    2,
		PowerAdder:         10,
		LifespanMultiplier: 5,
		LifespanAdder:      100,
	})
	if err != nil {
		t.Fatalf("build realm: %v", err)
	}
	return r
}

func TestNewCultivationState(t *testing.T) {
	realm := mustRealm(t)
	level := Level{Number: 1, Name: "First Level", BreakthroughPoints: 100}

	t.Run("creates a state at a realm and level", func(t *testing.T) {
		c, err := NewCultivationState(realm, level, 10, 3)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c.Realm.Name != "Qi Condensation" {
			t.Errorf("Realm = %q, want %q", c.Realm.Name, "Qi Condensation")
		}
		if c.Level.Number != 1 {
			t.Errorf("Level = %d, want 1", c.Level.Number)
		}
		if c.Points != 10 {
			t.Errorf("Points = %d, want 10", c.Points)
		}
		if c.Progress != 3 {
			t.Errorf("Progress = %v, want 3", c.Progress)
		}
		if c.Kind() != Cultivation {
			t.Errorf("Kind() = %q, want %q", c.Kind(), Cultivation)
		}
	})

	t.Run("rejects negative progress", func(t *testing.T) {
		_, err := NewCultivationState(realm, level, 0, -1)
		if !errors.Is(err, ErrInvalidProgress) {
			t.Fatalf("err = %v, want %v", err, ErrInvalidProgress)
		}
	})
}

func TestCultivationStateDerivedStats(t *testing.T) {
	// Power/Lifespan are the realm's ax+b formulas evaluated at Progress.
	c, err := NewCultivationState(mustRealm(t), Level{Number: 1, Name: "First Level"}, 0, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := c.Power(); got != 16 { // 2*3 + 10
		t.Errorf("Power() = %v, want 16", got)
	}
	if got := c.Lifespan(); got != 115 { // 5*3 + 100
		t.Errorf("Lifespan() = %v, want 115", got)
	}
}
