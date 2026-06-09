package sqlite_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"tge/internal/adapter/sqlite"
	"tge/internal/core/domain/novel"
	"tge/internal/core/port"
)

func newNovelRepo(t *testing.T) *sqlite.NovelRepository {
	t.Helper()
	dsn := filepath.Join(t.TempDir(), "test.db")
	repo, err := sqlite.NewNovelRepository(dsn)
	if err != nil {
		t.Fatalf("open novel repo: %v", err)
	}
	t.Cleanup(func() { repo.Close() })
	return repo
}

func TestNovelRepository_SaveListFind(t *testing.T) {
	repo := newNovelRepo(t)
	ctx := context.Background()

	ascension, _ := novel.NewNovel("Ascension", "Lin Feng")
	sorcery, _ := novel.NewNovel("Sorcerer's Path", "Mu Chen")
	if err := repo.Save(ctx, ascension); err != nil {
		t.Fatalf("save: %v", err)
	}
	if err := repo.Save(ctx, sorcery); err != nil {
		t.Fatalf("save: %v", err)
	}

	got, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(got) != 2 || got[0].Title != "Ascension" || got[1].Title != "Sorcerer's Path" {
		t.Errorf("list (want ordered by title) = %+v", got)
	}

	n, err := repo.FindByMainCharacter(ctx, "Mu Chen")
	if err != nil {
		t.Fatalf("find by MC: %v", err)
	}
	if n.Title != "Sorcerer's Path" {
		t.Errorf("FindByMainCharacter = %+v, want Sorcerer's Path", n)
	}
}

func TestNovelRepository_DuplicateTitle(t *testing.T) {
	repo := newNovelRepo(t)
	ctx := context.Background()
	n, _ := novel.NewNovel("Ascension", "Lin Feng")
	if err := repo.Save(ctx, n); err != nil {
		t.Fatalf("first save: %v", err)
	}
	other, _ := novel.NewNovel("Ascension", "Mu Chen")
	if err := repo.Save(ctx, other); !errors.Is(err, port.ErrNovelExists) {
		t.Fatalf("err = %v, want %v", err, port.ErrNovelExists)
	}
}

func TestNovelRepository_FindByMainCharacterMissing(t *testing.T) {
	repo := newNovelRepo(t)
	_, err := repo.FindByMainCharacter(context.Background(), "Nobody")
	if !errors.Is(err, port.ErrNovelNotFound) {
		t.Fatalf("err = %v, want %v", err, port.ErrNovelNotFound)
	}
}

func TestNovelRepository_StructureRoundTrip(t *testing.T) {
	repo := newNovelRepo(t)
	ctx := context.Background()

	n, _ := novel.NewNovel("Ascension", "Lin Feng")
	if err := repo.Save(ctx, n); err != nil {
		t.Fatalf("save: %v", err)
	}
	_ = n.AddVolume(1, "Beginnings")
	_ = n.AddVolume(2, "Trials")
	_ = n.AddChapter(1, 1, "A Mortal's Dream")
	_ = n.AddChapter(1, 2, "First Steps")
	_ = n.AddChapter(2, 1, "The Sect")
	if err := repo.SaveStructure(ctx, n); err != nil {
		t.Fatalf("save structure: %v", err)
	}

	got, err := repo.FindByTitle(ctx, "Ascension")
	if err != nil {
		t.Fatalf("find by title: %v", err)
	}
	if len(got.Volumes) != 2 {
		t.Fatalf("want 2 volumes, got %+v", got.Volumes)
	}
	if got.Volumes[0].Number != 1 || got.Volumes[0].Title != "Beginnings" || len(got.Volumes[0].Chapters) != 2 {
		t.Errorf("volume 1 mismatch: %+v", got.Volumes[0])
	}
	if got.Volumes[0].Chapters[1].Title != "First Steps" {
		t.Errorf("chapter mismatch: %+v", got.Volumes[0].Chapters)
	}
	if got.Volumes[1].Number != 2 || len(got.Volumes[1].Chapters) != 1 {
		t.Errorf("volume 2 mismatch: %+v", got.Volumes[1])
	}

	// SaveStructure replaces rather than appends.
	n2, _ := repo.FindByTitle(ctx, "Ascension")
	_ = n2.AddVolume(3, "Ascension")
	if err := repo.SaveStructure(ctx, n2); err != nil {
		t.Fatalf("re-save: %v", err)
	}
	again, _ := repo.FindByTitle(ctx, "Ascension")
	if len(again.Volumes) != 3 {
		t.Errorf("want 3 volumes after replace, got %d", len(again.Volumes))
	}
}

func TestNovelRepository_FindByTitleMissing(t *testing.T) {
	repo := newNovelRepo(t)
	_, err := repo.FindByTitle(context.Background(), "Ghost")
	if !errors.Is(err, port.ErrNovelNotFound) {
		t.Fatalf("err = %v, want %v", err, port.ErrNovelNotFound)
	}
}
