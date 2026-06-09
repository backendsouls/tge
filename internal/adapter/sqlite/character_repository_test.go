package sqlite_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"tge/internal/adapter/sqlite"
	"tge/internal/core/domain/character"
	"tge/internal/core/domain/progression"
	"tge/internal/core/port"
)

func newCharRepo(t *testing.T, dsn string) *sqlite.CharacterRepository {
	t.Helper()
	repo, err := sqlite.NewCharacterRepository(dsn)
	if err != nil {
		t.Fatalf("open character repo: %v", err)
	}
	t.Cleanup(func() { repo.Close() })
	return repo
}

func mortalFixture(t *testing.T, name string, ctype character.CharacterType, gender character.Gender, systems ...string) character.Character {
	t.Helper()
	ps := make([]progression.PowerSystem, 0, len(systems))
	for _, s := range systems {
		ps = append(ps, progression.PowerSystem{Name: s})
	}
	ch, err := character.NewMortalCharacter(character.CharacterConfig{
		Name:         name,
		Type:         ctype,
		Gender:       gender,
		PowerSystems: ps,
	})
	if err != nil {
		t.Fatal(err)
	}
	return ch
}

func TestCharacterRepository_SaveAndFind(t *testing.T) {
	dsn := filepath.Join(t.TempDir(), "test.db")
	ctx := context.Background()
	repo := newCharRepo(t, dsn)

	want := mortalFixture(t, "Yun Hai", character.Hero, character.Male, "Universe A Cultivation", "Universe B Sorcery")
	if err := repo.Save(ctx, want); err != nil {
		t.Fatalf("save: %v", err)
	}

	got, err := repo.FindByName(ctx, "Yun Hai")
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if got.Name != "Yun Hai" || got.Type != character.Hero || got.Gender != character.Male {
		t.Errorf("got %+v", got)
	}
	if got.PowerValue != want.PowerValue || got.Mortal != want.Mortal {
		t.Errorf("got %+v, want power %q mortal %+v", got, want.PowerValue, want.Mortal)
	}
	names := []string{}
	for _, s := range got.PowerSystems {
		names = append(names, s.Name)
	}
	if len(names) != 2 || names[0] != "Universe A Cultivation" || names[1] != "Universe B Sorcery" {
		t.Errorf("systems = %v, want [Universe A Cultivation, Universe B Sorcery]", names)
	}
}

func TestCharacterRepository_MainCharacters(t *testing.T) {
	dsn := filepath.Join(t.TempDir(), "test.db")
	ctx := context.Background()
	repo := newCharRepo(t, dsn)

	if mcs, err := repo.MainCharacters(ctx); err != nil || len(mcs) != 0 {
		t.Fatalf("empty: mcs = %v, err = %v", mcs, err)
	}

	if err := repo.Save(ctx, mortalFixture(t, "Bai Li", character.SideCharacter, character.Female, "Universe A Cultivation")); err != nil {
		t.Fatal(err)
	}
	if err := repo.Save(ctx, mortalFixture(t, "Mu Chen", character.MainCharacter, character.Male, "Universe A Cultivation")); err != nil {
		t.Fatal(err)
	}
	if err := repo.Save(ctx, mortalFixture(t, "Lin Feng", character.MainCharacter, character.Female, "Universe A Cultivation")); err != nil {
		t.Fatal(err)
	}

	mcs, err := repo.MainCharacters(ctx)
	if err != nil {
		t.Fatalf("main characters: %v", err)
	}
	// Only main characters, ordered by name; the side character is excluded.
	if len(mcs) != 2 || mcs[0].Name != "Lin Feng" || mcs[1].Name != "Mu Chen" {
		t.Errorf("main characters = %+v, want [Lin Feng, Mu Chen]", mcs)
	}
}

func TestCharacterRepository_DuplicateName(t *testing.T) {
	dsn := filepath.Join(t.TempDir(), "test.db")
	ctx := context.Background()
	repo := newCharRepo(t, dsn)

	ch := mortalFixture(t, "Lin Feng", character.MainCharacter, character.Female, "Universe A Cultivation")
	if err := repo.Save(ctx, ch); err != nil {
		t.Fatalf("first save: %v", err)
	}
	if err := repo.Save(ctx, ch); !errors.Is(err, port.ErrCharacterExists) {
		t.Fatalf("err = %v, want %v", err, port.ErrCharacterExists)
	}
}

func TestCharacterRepository_List(t *testing.T) {
	dsn := filepath.Join(t.TempDir(), "test.db")
	ctx := context.Background()
	repo := newCharRepo(t, dsn)

	if err := repo.Save(ctx, mortalFixture(t, "Yun Hai", character.Hero, character.Male, "Universe A Cultivation")); err != nil {
		t.Fatal(err)
	}
	if err := repo.Save(ctx, mortalFixture(t, "Lin Feng", character.MainCharacter, character.Female, "Universe A Cultivation")); err != nil {
		t.Fatal(err)
	}

	got, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(got) != 2 || got[0].Name != "Lin Feng" || got[1].Name != "Yun Hai" {
		t.Errorf("unexpected list (want ordered by name): %+v", got)
	}
}

var _ port.CharacterRepository = (*sqlite.CharacterRepository)(nil)
