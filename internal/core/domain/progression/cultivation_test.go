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
		c, err := NewCultivationState(realm, level, 10, 0, 3)
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
		_, err := NewCultivationState(realm, level, 0, 0, -1)
		if !errors.Is(err, ErrInvalidProgress) {
			t.Fatalf("err = %v, want %v", err, ErrInvalidProgress)
		}
	})
}

func TestCultivationStateAdvanceWithin(t *testing.T) {
	realm := mustRealm(t)
	if err := realm.AddLevel(1, "First Level", 100, 20); err != nil {
		t.Fatal(err)
	}
	if err := realm.AddLevel(2, "Second Level", 300, 60); err != nil {
		t.Fatal(err)
	}
	start := CultivationState{Realm: realm, Level: realm.Levels[0]}

	t.Run("fills the breakthrough gate before the bottleneck", func(t *testing.T) {
		c, left := start.AdvanceWithin(60)
		if left != 0 || c.Points != 60 || c.Bottleneck != 0 {
			t.Fatalf("points=%d bottleneck=%d left=%d", c.Points, c.Bottleneck, left)
		}
	})

	t.Run("overflow spills into the bottleneck once breakthrough is full", func(t *testing.T) {
		c, left := start.AdvanceWithin(110) // 100 breakthrough + 10 bottleneck
		if left != 0 || c.Points != 100 || c.Bottleneck != 10 || c.Level.Number != 1 {
			t.Fatalf("level=%d points=%d bottleneck=%d left=%d", c.Level.Number, c.Points, c.Bottleneck, left)
		}
	})

	t.Run("advances to the next level when both gates fill", func(t *testing.T) {
		c, left := start.AdvanceWithin(150) // fills L1 (100+20), 30 into L2 breakthrough
		if left != 0 || c.Level.Number != 2 || c.Points != 30 || c.Bottleneck != 0 {
			t.Fatalf("level=%d points=%d bottleneck=%d left=%d", c.Level.Number, c.Points, c.Bottleneck, left)
		}
	})

	t.Run("returns leftover points at the realm ceiling", func(t *testing.T) {
		c, left := start.AdvanceWithin(120 + 360 + 25) // fill L1 and L2 fully, 25 beyond
		if c.Level.Number != 2 || !c.Ready() || left != 25 {
			t.Fatalf("level=%d ready=%v left=%d", c.Level.Number, c.Ready(), left)
		}
	})
}

func TestCultivationStateDerivedStats(t *testing.T) {
	// Power/Lifespan are the realm's ax+b formulas evaluated at Progress.
	c, err := NewCultivationState(mustRealm(t), Level{Number: 1, Name: "First Level"}, 0, 0, 3)
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
