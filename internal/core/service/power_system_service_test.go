package service_test

import (
	"context"
	"errors"
	"testing"

	"tge/internal/core/domain/progression"
	"tge/internal/core/port"
	"tge/internal/core/service"
)

// fakePSRepo is an in-memory PowerSystemRepository for exercising the service.
type fakePSRepo struct {
	systems map[string]progression.PowerSystem
}

func newFakePSRepo() *fakePSRepo {
	return &fakePSRepo{systems: map[string]progression.PowerSystem{}}
}

func (f *fakePSRepo) Create(_ context.Context, name string) error {
	if _, ok := f.systems[name]; ok {
		return port.ErrPowerSystemExists
	}
	f.systems[name] = progression.PowerSystem{Name: name}
	return nil
}

func (f *fakePSRepo) FindByName(_ context.Context, name string) (progression.PowerSystem, error) {
	ps, ok := f.systems[name]
	if !ok {
		return progression.PowerSystem{}, port.ErrPowerSystemNotFound
	}
	return ps, nil
}

func (f *fakePSRepo) List(context.Context) ([]progression.PowerSystem, error) {
	out := make([]progression.PowerSystem, 0, len(f.systems))
	for _, ps := range f.systems {
		out = append(out, ps)
	}
	return out, nil
}

func (f *fakePSRepo) SavePowers(_ context.Context, system progression.PowerSystem) error {
	f.systems[system.Name] = system
	return nil
}

func TestPowerSystemService_CreateSystem(t *testing.T) {
	repo := newFakePSRepo()
	svc := service.NewPowerSystemService(repo)

	ps, err := svc.CreateSystem(context.Background(), "Universe A Cultivation", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ps.Name != "Universe A Cultivation" {
		t.Errorf("Name = %q", ps.Name)
	}
	if _, err := svc.CreateSystem(context.Background(), "Universe A Cultivation", ""); !errors.Is(err, port.ErrPowerSystemExists) {
		t.Fatalf("err = %v, want %v", err, port.ErrPowerSystemExists)
	}
}

func TestPowerSystemService_AddPower(t *testing.T) {
	repo := newFakePSRepo()
	svc := service.NewPowerSystemService(repo)
	ctx := context.Background()
	if _, err := svc.CreateSystem(ctx, "Universe A Cultivation", ""); err != nil {
		t.Fatal(err)
	}

	t.Run("adds a power and persists the tree", func(t *testing.T) {
		ps, err := svc.AddPower(ctx, port.AddPowerInput{System: "Universe A Cultivation", Name: "Body"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(ps.Powers) != 1 || ps.Powers[0].Name != "Body" {
			t.Errorf("tree = %+v", ps.Powers)
		}
		// Persisted, so a follow-up nested add sees the parent.
		if _, err := svc.AddPower(ctx, port.AddPowerInput{System: "Universe A Cultivation", Name: "Iron Skin", Parent: "Body"}); err != nil {
			t.Fatalf("nested add failed: %v", err)
		}
	})

	t.Run("fails for an unknown system", func(t *testing.T) {
		_, err := svc.AddPower(ctx, port.AddPowerInput{System: "Nope", Name: "Body"})
		if !errors.Is(err, port.ErrPowerSystemNotFound) {
			t.Fatalf("err = %v, want %v", err, port.ErrPowerSystemNotFound)
		}
	})
}
