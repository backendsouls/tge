package progression

import (
	"errors"
	"testing"
)

func TestNewRealm(t *testing.T) {
	t.Run("creates a realm with the given attributes", func(t *testing.T) {
		r, err := NewRealm(RealmConfig{
			Name:               "Qi Condensation",
			PowerMultiplier:    2,
			PowerAdder:         10,
			LifespanMultiplier: 5,
			LifespanAdder:      100,
			BottleneckPoints:   250,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if r.Name != "Qi Condensation" {
			t.Errorf("Name = %q, want %q", r.Name, "Qi Condensation")
		}
		if r.PowerMultiplier != 2 {
			t.Errorf("PowerMultiplier = %v, want %v", r.PowerMultiplier, 2.0)
		}
		if r.PowerAdder != 10 {
			t.Errorf("PowerAdder = %v, want %v", r.PowerAdder, 10.0)
		}
		if r.LifespanMultiplier != 5 {
			t.Errorf("LifespanMultiplier = %v, want %v", r.LifespanMultiplier, 5.0)
		}
		if r.LifespanAdder != 100 {
			t.Errorf("LifespanAdder = %v, want %v", r.LifespanAdder, 100.0)
		}
		if r.BottleneckPoints != 250 {
			t.Errorf("BottleneckPoints = %v, want %v", r.BottleneckPoints, 250)
		}
	})

	t.Run("rejects an empty name", func(t *testing.T) {
		_, err := NewRealm(RealmConfig{Name: "  "})
		if !errors.Is(err, ErrInvalidName) {
			t.Fatalf("err = %v, want %v", err, ErrInvalidName)
		}
	})

	t.Run("rejects negative bottleneck points", func(t *testing.T) {
		_, err := NewRealm(RealmConfig{Name: "Foundation", BottleneckPoints: -1})
		if !errors.Is(err, ErrInvalidPoints) {
			t.Fatalf("err = %v, want %v", err, ErrInvalidPoints)
		}
	})

	t.Run("trims surrounding whitespace from the name", func(t *testing.T) {
		r, err := NewRealm(RealmConfig{Name: "  Core Formation  "})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if r.Name != "Core Formation" {
			t.Errorf("Name = %q, want %q", r.Name, "Core Formation")
		}
	})
}

func TestRealmPower(t *testing.T) {
	// Power follows ax + b where a = PowerMultiplier, b = PowerAdder.
	r, err := NewRealm(RealmConfig{
		Name:            "Qi Condensation",
		PowerMultiplier: 3,
		PowerAdder:      7,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cases := []struct {
		x    float64
		want float64
	}{
		{x: 0, want: 7},
		{x: 1, want: 10},
		{x: 10, want: 37},
	}
	for _, c := range cases {
		if got := r.Power(c.x); got != c.want {
			t.Errorf("Power(%v) = %v, want %v", c.x, got, c.want)
		}
	}
}

func TestRealmLifespan(t *testing.T) {
	// Lifespan follows ax + b where a = LifespanMultiplier, b = LifespanAdder.
	r, err := NewRealm(RealmConfig{
		Name:               "Qi Condensation",
		LifespanMultiplier: 5,
		LifespanAdder:      100,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cases := []struct {
		x    float64
		want float64
	}{
		{x: 0, want: 100},
		{x: 1, want: 105},
		{x: 20, want: 200},
	}
	for _, c := range cases {
		if got := r.Lifespan(c.x); got != c.want {
			t.Errorf("Lifespan(%v) = %v, want %v", c.x, got, c.want)
		}
	}
}
