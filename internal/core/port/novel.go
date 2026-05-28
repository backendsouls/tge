package port

import (
	"context"
	"errors"
	"tge/internal/core/domain/novel"
)

var (
	// ErrNovelNotFound is returned when a lookup finds no matching novel.
	ErrNovelNotFound = errors.New("novel: not found")
	// ErrNovelExists is returned when saving a novel whose title is taken.
	ErrNovelExists = errors.New("novel: already exists")
	// ErrNotMainCharacter is returned when a novel references a non-main character.
	ErrNotMainCharacter = errors.New("novel: referenced character is not a main character")
	// ErrMainCharacterTaken is returned when the character already leads another novel.
	ErrMainCharacterTaken = errors.New("novel: main character already belongs to another novel")
)

// CreateNovelInput describes the novel to create. MainCharacter must name an
// existing character of type MainCharacter not already leading another novel.
type CreateNovelInput struct {
	Title         string
	MainCharacter string
}

// AddVolumeInput describes a volume to add to a novel.
type AddVolumeInput struct {
	Novel  string
	Number int
	Title  string
}

// AddChapterInput describes a chapter to add to a novel's volume.
type AddChapterInput struct {
	Novel  string
	Volume int
	Number int
	Title  string
}

// NovelRepository is a driven port persisting novels keyed by title.
type NovelRepository interface {
	Save(ctx context.Context, n novel.Novel) error
	List(ctx context.Context) ([]novel.Novel, error)
	FindByTitle(ctx context.Context, title string) (novel.Novel, error)
	// FindByMainCharacter returns the novel led by the named character, or ErrNovelNotFound.
	FindByMainCharacter(ctx context.Context, mainCharacter string) (novel.Novel, error)
	// SaveStructure replaces the stored volumes and chapters of a novel.
	SaveStructure(ctx context.Context, n novel.Novel) error
}

// NovelService is a driving port for novel use cases.
type NovelService interface {
	CreateNovel(ctx context.Context, in CreateNovelInput) (novel.Novel, error)
	AddVolume(ctx context.Context, in AddVolumeInput) (novel.Novel, error)
	AddChapter(ctx context.Context, in AddChapterInput) (novel.Novel, error)
	GetNovel(ctx context.Context, title string) (novel.Novel, error)
	ListNovels(ctx context.Context) ([]novel.Novel, error)
}
