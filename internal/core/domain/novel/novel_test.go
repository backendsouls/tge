package novel

import (
	"errors"
	"testing"
)

func TestNewNovel(t *testing.T) {
	t.Run("creates a novel bound to a main character", func(t *testing.T) {
		n, err := NewNovel("  Ascension  ", "  Lin Feng  ")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n.Title != "Ascension" {
			t.Errorf("Title = %q, want %q", n.Title, "Ascension")
		}
		if n.MainCharacter != "Lin Feng" {
			t.Errorf("character.MainCharacter = %q, want %q", n.MainCharacter, "Lin Feng")
		}
	})

	t.Run("rejects a blank title", func(t *testing.T) {
		_, err := NewNovel("  ", "Lin Feng")
		if !errors.Is(err, ErrInvalidNovelTitle) {
			t.Fatalf("err = %v, want %v", err, ErrInvalidNovelTitle)
		}
	})

	t.Run("requires a main character", func(t *testing.T) {
		_, err := NewNovel("Ascension", "  ")
		if !errors.Is(err, ErrNovelMainCharacterRequired) {
			t.Fatalf("err = %v, want %v", err, ErrNovelMainCharacterRequired)
		}
	})
}

func TestNovelVolumesAndChapters(t *testing.T) {
	t.Run("adds volumes and nested chapters", func(t *testing.T) {
		n, _ := NewNovel("Ascension", "Lin Feng")
		if err := n.AddVolume(1, "Beginnings"); err != nil {
			t.Fatalf("add volume: %v", err)
		}
		if err := n.AddChapter(1, 1, "A character.Mortal's Dream"); err != nil {
			t.Fatalf("add chapter: %v", err)
		}
		if err := n.AddChapter(1, 2, "First Steps"); err != nil {
			t.Fatalf("add chapter: %v", err)
		}

		if len(n.Volumes) != 1 || n.Volumes[0].Title != "Beginnings" {
			t.Fatalf("volume not added: %+v", n.Volumes)
		}
		if len(n.Volumes[0].Chapters) != 2 || n.Volumes[0].Chapters[1].Title != "First Steps" {
			t.Errorf("chapters not added: %+v", n.Volumes[0].Chapters)
		}
	})

	t.Run("rejects a duplicate volume number", func(t *testing.T) {
		n, _ := NewNovel("Ascension", "Lin Feng")
		_ = n.AddVolume(1, "A")
		if err := n.AddVolume(1, "B"); !errors.Is(err, ErrVolumeExists) {
			t.Fatalf("err = %v, want %v", err, ErrVolumeExists)
		}
	})

	t.Run("rejects a non-positive volume number", func(t *testing.T) {
		n, _ := NewNovel("Ascension", "Lin Feng")
		if err := n.AddVolume(0, "A"); !errors.Is(err, ErrInvalidVolumeNumber) {
			t.Fatalf("err = %v, want %v", err, ErrInvalidVolumeNumber)
		}
	})

	t.Run("rejects a chapter in a missing volume", func(t *testing.T) {
		n, _ := NewNovel("Ascension", "Lin Feng")
		if err := n.AddChapter(9, 1, "X"); !errors.Is(err, ErrVolumeNotFound) {
			t.Fatalf("err = %v, want %v", err, ErrVolumeNotFound)
		}
	})

	t.Run("rejects a duplicate chapter number within a volume", func(t *testing.T) {
		n, _ := NewNovel("Ascension", "Lin Feng")
		_ = n.AddVolume(1, "A")
		_ = n.AddChapter(1, 1, "X")
		if err := n.AddChapter(1, 1, "Y"); !errors.Is(err, ErrChapterExists) {
			t.Fatalf("err = %v, want %v", err, ErrChapterExists)
		}
	})
}
