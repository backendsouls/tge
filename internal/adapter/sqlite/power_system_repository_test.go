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

func newPSRepo(t *testing.T) *sqlite.PowerSystemRepository {
	t.Helper()
	dsn := filepath.Join(t.TempDir(), "test.db")
	repo, err := sqlite.NewPowerSystemRepository(dsn)
	if err != nil {
		t.Fatalf("open power system repo: %v", err)
	}
	t.Cleanup(func() { repo.Close() })
	return repo
}

func TestPowerSystemRepository_CreateAndFind(t *testing.T) {
	repo := newPSRepo(t)
	ctx := context.Background()

	if err := repo.Create(ctx, "Universe A Cultivation"); err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := repo.FindByName(ctx, "Universe A Cultivation")
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if got.Name != "Universe A Cultivation" {
		t.Errorf("Name = %q, want %q", got.Name, "Universe A Cultivation")
	}
	if len(got.Powers) != 0 {
		t.Errorf("new system should have no powers, got %v", got.Powers)
	}
}

func TestPowerSystemRepository_Duplicate(t *testing.T) {
	repo := newPSRepo(t)
	ctx := context.Background()
	if err := repo.Create(ctx, "S"); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := repo.Create(ctx, "S"); !errors.Is(err, port.ErrPowerSystemExists) {
		t.Fatalf("err = %v, want %v", err, port.ErrPowerSystemExists)
	}
}

func TestPowerSystemRepository_FindMissing(t *testing.T) {
	repo := newPSRepo(t)
	_, err := repo.FindByName(context.Background(), "Nope")
	if !errors.Is(err, port.ErrPowerSystemNotFound) {
		t.Fatalf("err = %v, want %v", err, port.ErrPowerSystemNotFound)
	}
}

func TestPowerSystemRepository_SaveAndLoadTree(t *testing.T) {
	repo := newPSRepo(t)
	ctx := context.Background()
	if err := repo.Create(ctx, "Universe A Cultivation"); err != nil {
		t.Fatalf("create: %v", err)
	}

	ps, _ := progression.NewPowerSystem("Universe A Cultivation")
	if err := ps.AddPower("Body", ""); err != nil {
		t.Fatal(err)
	}
	if err := ps.AddPower("Iron Skin", "Body"); err != nil {
		t.Fatal(err)
	}
	if err := ps.AddPower("Soul", ""); err != nil {
		t.Fatal(err)
	}
	if err := repo.SavePowers(ctx, ps); err != nil {
		t.Fatalf("save powers: %v", err)
	}

	got, err := repo.FindByName(ctx, "Universe A Cultivation")
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if !reflect.DeepEqual(got.Names(), []string{"Body", "Iron Skin", "Soul"}) {
		t.Errorf("Names() = %v, want [Body Iron Skin Soul]", got.Names())
	}
	if len(got.Powers) != 2 || got.Powers[0].Name != "Body" || len(got.Powers[0].Children) != 1 {
		t.Errorf("tree not reconstructed correctly: %+v", got.Powers)
	}
	if got.Powers[0].Children[0].Name != "Iron Skin" {
		t.Errorf("Iron Skin not nested under Body: %+v", got.Powers[0])
	}
}

func TestPowerSystemRepository_List(t *testing.T) {
	repo := newPSRepo(t)
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
		t.Errorf("unexpected list (want ordered by name): %+v", got)
	}
}
