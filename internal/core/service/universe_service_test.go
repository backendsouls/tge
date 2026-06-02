package service_test

import (
	"context"
	"errors"
	"testing"

	"tge/internal/core/domain/cosmology"
	"tge/internal/core/port"
	"tge/internal/core/service"
)

// noopTimelineRepo is a TimelineRepository double whose provisioning calls
// always succeed; the cosmology services auto-provision a timeline on creation.
type noopTimelineRepo struct{}

func (noopTimelineRepo) Save(context.Context, port.LocationRef, cosmology.Timeline) error { return nil }
func (noopTimelineRepo) AddEvent(context.Context, port.LocationRef, cosmology.Event) error {
	return nil
}
func (noopTimelineRepo) Find(context.Context, port.LocationRef) (cosmology.Timeline, error) {
	return cosmology.Timeline{}, nil
}

// fakeUniverseRepo is an in-memory UniverseRepository.
type fakeUniverseRepo struct {
	byName map[string]cosmology.Universe
}

func newFakeUniverseRepo() *fakeUniverseRepo {
	return &fakeUniverseRepo{byName: map[string]cosmology.Universe{}}
}

func (f *fakeUniverseRepo) Create(_ context.Context, name string) error {
	if _, ok := f.byName[name]; ok {
		return port.ErrUniverseExists
	}
	f.byName[name] = cosmology.Universe{Name: name}
	return nil
}

func (f *fakeUniverseRepo) FindByName(_ context.Context, name string) (cosmology.Universe, error) {
	u, ok := f.byName[name]
	if !ok {
		return cosmology.Universe{}, port.ErrUniverseNotFound
	}
	return u, nil
}

func (f *fakeUniverseRepo) List(context.Context) ([]cosmology.Universe, error) {
	out := make([]cosmology.Universe, 0, len(f.byName))
	for _, u := range f.byName {
		out = append(out, u)
	}
	return out, nil
}

func (f *fakeUniverseRepo) SaveSystems(_ context.Context, u cosmology.Universe) error {
	f.byName[u.Name] = u
	return nil
}

func (f *fakeUniverseRepo) SaveRealms(_ context.Context, u cosmology.Universe) error {
	f.byName[u.Name] = u
	return nil
}

func (f *fakeUniverseRepo) FindBySystem(_ context.Context, system string) (cosmology.Universe, error) {
	for _, u := range f.byName {
		for _, s := range u.Systems {
			if s.Name == system {
				return u, nil
			}
		}
	}
	return cosmology.Universe{}, port.ErrUniverseNotFound
}

func TestUniverseService_CreateUniverse(t *testing.T) {
	repo := newFakeUniverseRepo()
	svc := service.NewUniverseService(repo, knownSystems(), noopTimelineRepo{})

	if _, err := svc.CreateUniverse(context.Background(), "Universe A"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := svc.CreateUniverse(context.Background(), "Universe A"); !errors.Is(err, port.ErrUniverseExists) {
		t.Fatalf("err = %v, want %v", err, port.ErrUniverseExists)
	}
}

func TestUniverseService_AddSystem(t *testing.T) {
	systems := knownSystems("Cultivation", "Sorcery")

	t.Run("adds an existing unclaimed system", func(t *testing.T) {
		repo := newFakeUniverseRepo()
		svc := service.NewUniverseService(repo, systems, noopTimelineRepo{})
		if _, err := svc.CreateUniverse(context.Background(), "Universe A"); err != nil {
			t.Fatal(err)
		}
		u, err := svc.AddSystem(context.Background(), port.AddUniverseSystemInput{Universe: "Universe A", System: "Cultivation"})
		if err != nil {
			t.Fatalf("add system: %v", err)
		}
		if len(u.Systems) != 1 || u.Systems[0].Name != "Cultivation" {
			t.Errorf("systems = %+v", u.Systems)
		}
	})

	t.Run("rejects an unknown system", func(t *testing.T) {
		repo := newFakeUniverseRepo()
		svc := service.NewUniverseService(repo, systems, noopTimelineRepo{})
		_, _ = svc.CreateUniverse(context.Background(), "Universe A")
		_, err := svc.AddSystem(context.Background(), port.AddUniverseSystemInput{Universe: "Universe A", System: "Nope"})
		if !errors.Is(err, port.ErrPowerSystemNotFound) {
			t.Fatalf("err = %v, want %v", err, port.ErrPowerSystemNotFound)
		}
	})

	t.Run("rejects a system already in another universe", func(t *testing.T) {
		repo := newFakeUniverseRepo()
		svc := service.NewUniverseService(repo, systems, noopTimelineRepo{})
		_, _ = svc.CreateUniverse(context.Background(), "Universe A")
		_, _ = svc.CreateUniverse(context.Background(), "Universe B")
		if _, err := svc.AddSystem(context.Background(), port.AddUniverseSystemInput{Universe: "Universe A", System: "Cultivation"}); err != nil {
			t.Fatal(err)
		}
		_, err := svc.AddSystem(context.Background(), port.AddUniverseSystemInput{Universe: "Universe B", System: "Cultivation"})
		if !errors.Is(err, port.ErrPowerSystemTaken) {
			t.Fatalf("err = %v, want %v", err, port.ErrPowerSystemTaken)
		}
	})

	t.Run("rejects adding to a missing universe", func(t *testing.T) {
		repo := newFakeUniverseRepo()
		svc := service.NewUniverseService(repo, systems, noopTimelineRepo{})
		_, err := svc.AddSystem(context.Background(), port.AddUniverseSystemInput{Universe: "Ghost", System: "Cultivation"})
		if !errors.Is(err, port.ErrUniverseNotFound) {
			t.Fatalf("err = %v, want %v", err, port.ErrUniverseNotFound)
		}
	})
}

func TestUniverseService_AddRealms(t *testing.T) {
	repo := newFakeUniverseRepo()
	svc := service.NewUniverseService(repo, knownSystems(), noopTimelineRepo{})
	ctx := context.Background()
	if _, err := svc.CreateUniverse(ctx, "Universe A"); err != nil {
		t.Fatal(err)
	}

	t.Run("adds realms and persists them", func(t *testing.T) {
		u, err := svc.AddRealms(ctx, port.AddRealmsInput{Universe: "Universe A", Names: []string{"Hell Realm", "Mortal Realm", "Heaven Realm"}})
		if err != nil {
			t.Fatalf("add realms: %v", err)
		}
		if len(u.Realms) != 3 {
			t.Fatalf("want 3 realms, got %+v", u.Realms)
		}
		got, _ := svc.GetUniverse(ctx, "Universe A")
		if len(got.Realms) != 3 {
			t.Errorf("realms not persisted: %+v", got.Realms)
		}
	})

	t.Run("rejects a missing universe", func(t *testing.T) {
		_, err := svc.AddRealms(ctx, port.AddRealmsInput{Universe: "Ghost", Names: []string{"A", "B"}})
		if !errors.Is(err, port.ErrUniverseNotFound) {
			t.Fatalf("err = %v, want %v", err, port.ErrUniverseNotFound)
		}
	})
}
