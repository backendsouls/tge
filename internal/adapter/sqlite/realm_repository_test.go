package sqlite_test

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"tge/internal/adapter/sqlite"
	"tge/internal/core/domain/progression"
	"tge/internal/core/port"
)

// newRepo opens a fresh, migrated database in a temp file for each test.
func newRepo(t *testing.T) *sqlite.RealmRepository {
	t.Helper()
	dsn := filepath.Join(t.TempDir(), "test.db")
	repo, err := sqlite.NewRealmRepository(dsn)
	if err != nil {
		t.Fatalf("open repository: %v", err)
	}
	t.Cleanup(func() { _ = repo.Close() })
	return repo
}

func sampleRealm(t *testing.T, name string) progression.Realm {
	t.Helper()
	r, err := progression.NewRealm(progression.RealmConfig{
		Name:               name,
		PowerMultiplier:    2,
		PowerAdder:         10,
		LifespanMultiplier: 5,
		LifespanAdder:      100,
		BottleneckPoints:   250,
	})
	if err != nil {
		t.Fatalf("build realm: %v", err)
	}
	return r
}

func TestRealmRepository_SaveAndFind(t *testing.T) {
	repo := newRepo(t)
	ctx := context.Background()
	want := sampleRealm(t, "Qi Condensation")

	if err := repo.Save(ctx, want); err != nil {
		t.Fatalf("save: %v", err)
	}

	got, err := repo.FindByName(ctx, "Qi Condensation")
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("round-trip mismatch:\n got %+v\nwant %+v", got, want)
	}
}

func TestRealmRepository_FindByNameMissing(t *testing.T) {
	repo := newRepo(t)
	_, err := repo.FindByName(context.Background(), "Nope")
	if !errors.Is(err, port.ErrRealmNotFound) {
		t.Fatalf("err = %v, want %v", err, port.ErrRealmNotFound)
	}
}

func TestRealmRepository_DuplicateName(t *testing.T) {
	repo := newRepo(t)
	ctx := context.Background()
	r := sampleRealm(t, "Foundation")

	if err := repo.Save(ctx, r); err != nil {
		t.Fatalf("first save: %v", err)
	}
	err := repo.Save(ctx, r)
	if !errors.Is(err, port.ErrRealmExists) {
		t.Fatalf("err = %v, want %v", err, port.ErrRealmExists)
	}
}

func TestRealmRepository_List(t *testing.T) {
	repo := newRepo(t)
	ctx := context.Background()
	if err := repo.Save(ctx, sampleRealm(t, "Qi Condensation")); err != nil {
		t.Fatalf("save: %v", err)
	}
	if err := repo.Save(ctx, sampleRealm(t, "Foundation")); err != nil {
		t.Fatalf("save: %v", err)
	}

	got, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d realms, want 2", len(got))
	}
	// Ordered by name for a stable, predictable listing.
	if got[0].Name != "Foundation" || got[1].Name != "Qi Condensation" {
		t.Errorf("unexpected order: %q, %q", got[0].Name, got[1].Name)
	}
}

// Ensures the adapter satisfies the driven port (compile-time check).
var _ port.RealmRepository = (*sqlite.RealmRepository)(nil)
