package sqlite_test

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"tge/internal/adapter/sqlite"
	"tge/internal/core/domain/cosmology"
	"tge/internal/core/port"
)

func newUniverseRepo(t *testing.T) *sqlite.UniverseRepository {
	t.Helper()
	dsn := filepath.Join(t.TempDir(), "test.db")
	repo, err := sqlite.NewUniverseRepository(dsn)
	if err != nil {
		t.Fatalf("open universe repo: %v", err)
	}
	t.Cleanup(func() { _ = repo.Close() })
	return repo
}

func TestUniverseRepository_CreateAndFind(t *testing.T) {
	repo := newUniverseRepo(t)
	ctx := context.Background()

	if err := repo.Create(ctx, "Universe A"); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := repo.Create(ctx, "Universe A"); !errors.Is(err, port.ErrUniverseExists) {
		t.Fatalf("dup: err = %v, want %v", err, port.ErrUniverseExists)
	}

	got, err := repo.FindByName(ctx, "Universe A")
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if got.Name != "Universe A" || len(got.Systems) != 0 {
		t.Errorf("got %+v", got)
	}

	if _, err := repo.FindByName(ctx, "Nope"); !errors.Is(err, port.ErrUniverseNotFound) {
		t.Fatalf("missing: err = %v, want %v", err, port.ErrUniverseNotFound)
	}
}

func TestUniverseRepository_SaveSystemsAndFindBySystem(t *testing.T) {
	repo := newUniverseRepo(t)
	ctx := context.Background()
	if err := repo.Create(ctx, "Universe A"); err != nil {
		t.Fatal(err)
	}

	u, _ := cosmology.NewUniverse("Universe A")
	_ = u.AddSystem("Cultivation")
	_ = u.AddSystem("Sorcery")
	if err := repo.SaveSystems(ctx, u); err != nil {
		t.Fatalf("save systems: %v", err)
	}

	got, err := repo.FindByName(ctx, "Universe A")
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	names := []string{got.Systems[0].Name, got.Systems[1].Name}
	if !reflect.DeepEqual(names, []string{"Cultivation", "Sorcery"}) {
		t.Errorf("systems = %v, want [Cultivation Sorcery]", names)
	}

	owner, err := repo.FindBySystem(ctx, "Cultivation")
	if err != nil {
		t.Fatalf("find by system: %v", err)
	}
	if owner.Name != "Universe A" {
		t.Errorf("owner = %q, want Universe A", owner.Name)
	}

	if _, err := repo.FindBySystem(ctx, "Unclaimed"); !errors.Is(err, port.ErrUniverseNotFound) {
		t.Fatalf("unclaimed: err = %v, want %v", err, port.ErrUniverseNotFound)
	}
}

func TestUniverseRepository_SaveRealms(t *testing.T) {
	repo := newUniverseRepo(t)
	ctx := context.Background()
	if err := repo.Create(ctx, "Universe A"); err != nil {
		t.Fatal(err)
	}

	u, _ := cosmology.NewUniverse("Universe A")
	_ = u.AddRealms("Hell Realm", "Mortal Realm", "Heaven Realm")
	if err := repo.SaveRealms(ctx, u); err != nil {
		t.Fatalf("save realms: %v", err)
	}

	got, err := repo.FindByName(ctx, "Universe A")
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	names := []string{got.Realms[0].Name, got.Realms[1].Name, got.Realms[2].Name}
	if !reflect.DeepEqual(names, []string{"Hell Realm", "Mortal Realm", "Heaven Realm"}) {
		t.Errorf("realms = %v, want [Hell Realm Mortal Realm Heaven Realm]", names)
	}
}

func TestUniverseRepository_List(t *testing.T) {
	repo := newUniverseRepo(t)
	ctx := context.Background()
	if err := repo.Create(ctx, "Universe B"); err != nil {
		t.Fatal(err)
	}
	if err := repo.Create(ctx, "Universe A"); err != nil {
		t.Fatal(err)
	}
	got, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(got) != 2 || got[0].Name != "Universe A" || got[1].Name != "Universe B" {
		t.Errorf("list (want ordered by name) = %+v", got)
	}
}
