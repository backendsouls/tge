package service_test

import (
	"context"
	"errors"
	"testing"

	"tge/internal/core/domain/character"
	"tge/internal/core/domain/novel"
	"tge/internal/core/port"
	"tge/internal/core/service"
)

// fakeNovelRepo is an in-memory NovelRepository.
type fakeNovelRepo struct {
	byTitle map[string]novel.Novel
	saveErr error
}

func newFakeNovelRepo() *fakeNovelRepo {
	return &fakeNovelRepo{byTitle: map[string]novel.Novel{}}
}

func (f *fakeNovelRepo) Save(_ context.Context, n novel.Novel) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	if _, ok := f.byTitle[n.Title]; ok {
		return port.ErrNovelExists
	}
	f.byTitle[n.Title] = n
	return nil
}

func (f *fakeNovelRepo) List(context.Context) ([]novel.Novel, error) {
	out := make([]novel.Novel, 0, len(f.byTitle))
	for _, n := range f.byTitle {
		out = append(out, n)
	}
	return out, nil
}

func (f *fakeNovelRepo) FindByMainCharacter(_ context.Context, mc string) (novel.Novel, error) {
	for _, n := range f.byTitle {
		if n.MainCharacter == mc {
			return n, nil
		}
	}
	return novel.Novel{}, port.ErrNovelNotFound
}

func (f *fakeNovelRepo) FindByTitle(_ context.Context, title string) (novel.Novel, error) {
	n, ok := f.byTitle[title]
	if !ok {
		return novel.Novel{}, port.ErrNovelNotFound
	}
	return n, nil
}

func (f *fakeNovelRepo) SaveStructure(_ context.Context, n novel.Novel) error {
	f.byTitle[n.Title] = n
	return nil
}

func charsByName(cs ...character.Character) *stubCharRepo {
	m := map[string]character.Character{}
	for _, c := range cs {
		m[c.Name] = c
	}
	return &stubCharRepo{byName: m}
}

func TestNovelService_CreateNovel(t *testing.T) {
	mc := character.Character{Name: "Lin Feng", Type: character.MainCharacter, Gender: character.Female}
	side := character.Character{Name: "Bai Li", Type: character.SideCharacter, Gender: character.Female}

	t.Run("creates a novel for a free main character", func(t *testing.T) {
		novels := newFakeNovelRepo()
		svc := service.NewNovelService(novels, charsByName(mc))

		n, err := svc.CreateNovel(context.Background(), port.CreateNovelInput{Title: "Ascension", MainCharacter: "Lin Feng"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n.Title != "Ascension" || n.MainCharacter != "Lin Feng" {
			t.Errorf("got %+v", n)
		}
	})

	t.Run("rejects an unknown character", func(t *testing.T) {
		svc := service.NewNovelService(newFakeNovelRepo(), charsByName())
		_, err := svc.CreateNovel(context.Background(), port.CreateNovelInput{Title: "X", MainCharacter: "Ghost"})
		if !errors.Is(err, port.ErrCharacterNotFound) {
			t.Fatalf("err = %v, want %v", err, port.ErrCharacterNotFound)
		}
	})

	t.Run("rejects a non-main character", func(t *testing.T) {
		svc := service.NewNovelService(newFakeNovelRepo(), charsByName(side))
		_, err := svc.CreateNovel(context.Background(), port.CreateNovelInput{Title: "X", MainCharacter: "Bai Li"})
		if !errors.Is(err, port.ErrNotMainCharacter) {
			t.Fatalf("err = %v, want %v", err, port.ErrNotMainCharacter)
		}
	})

	t.Run("rejects a main character already leading a novel", func(t *testing.T) {
		novels := newFakeNovelRepo()
		svc := service.NewNovelService(novels, charsByName(mc))
		if _, err := svc.CreateNovel(context.Background(), port.CreateNovelInput{Title: "Ascension", MainCharacter: "Lin Feng"}); err != nil {
			t.Fatal(err)
		}
		_, err := svc.CreateNovel(context.Background(), port.CreateNovelInput{Title: "Another", MainCharacter: "Lin Feng"})
		if !errors.Is(err, port.ErrMainCharacterTaken) {
			t.Fatalf("err = %v, want %v", err, port.ErrMainCharacterTaken)
		}
	})
}

func TestNovelService_VolumesAndChapters(t *testing.T) {
	mc := character.Character{Name: "Lin Feng", Type: character.MainCharacter, Gender: character.Female}
	novels := newFakeNovelRepo()
	svc := service.NewNovelService(novels, charsByName(mc))
	ctx := context.Background()
	if _, err := svc.CreateNovel(ctx, port.CreateNovelInput{Title: "Ascension", MainCharacter: "Lin Feng"}); err != nil {
		t.Fatal(err)
	}

	t.Run("adds a volume then a chapter, persisting both", func(t *testing.T) {
		if _, err := svc.AddVolume(ctx, port.AddVolumeInput{Novel: "Ascension", Number: 1, Title: "Beginnings"}); err != nil {
			t.Fatalf("add volume: %v", err)
		}
		n, err := svc.AddChapter(ctx, port.AddChapterInput{Novel: "Ascension", Volume: 1, Number: 1, Title: "A Mortal's Dream"})
		if err != nil {
			t.Fatalf("add chapter: %v", err)
		}
		if len(n.Volumes) != 1 || len(n.Volumes[0].Chapters) != 1 {
			t.Fatalf("structure not built: %+v", n.Volumes)
		}

		got, _ := svc.GetNovel(ctx, "Ascension")
		if len(got.Volumes) != 1 || got.Volumes[0].Chapters[0].Title != "A Mortal's Dream" {
			t.Errorf("structure not persisted: %+v", got.Volumes)
		}
	})

	t.Run("adding a volume to a missing novel fails", func(t *testing.T) {
		_, err := svc.AddVolume(ctx, port.AddVolumeInput{Novel: "Ghost", Number: 1})
		if !errors.Is(err, port.ErrNovelNotFound) {
			t.Fatalf("err = %v, want %v", err, port.ErrNovelNotFound)
		}
	})
}
