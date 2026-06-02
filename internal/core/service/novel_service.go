package service

import (
	"context"
	"errors"

	"tge/internal/core/domain/character"
	"tge/internal/core/domain/novel"
	"tge/internal/core/port"
)

// NovelService implements port.NovelService. It associates a novel with an
// existing main character that does not already lead another novel.
type NovelService struct {
	novels port.NovelRepository
	chars  port.CharacterRepository
}

// NewNovelService wires the service to its repositories (both interfaces).
func NewNovelService(novels port.NovelRepository, chars port.CharacterRepository) *NovelService {
	return &NovelService{novels: novels, chars: chars}
}

// CreateNovel validates the referenced character is a free main character and
// persists the novel.
func (s *NovelService) CreateNovel(ctx context.Context, in port.CreateNovelInput) (novel.Novel, error) {
	n, err := novel.NewNovel(in.Title, in.MainCharacter)
	if err != nil {
		return novel.Novel{}, err
	}

	c, err := s.chars.FindByName(ctx, n.MainCharacter)
	if err != nil {
		return novel.Novel{}, err
	}
	if c.Type != character.MainCharacter {
		return novel.Novel{}, port.ErrNotMainCharacter
	}

	switch _, err := s.novels.FindByMainCharacter(ctx, n.MainCharacter); {
	case err == nil:
		return novel.Novel{}, port.ErrMainCharacterTaken
	case errors.Is(err, port.ErrNovelNotFound):
		// free to use
	default:
		return novel.Novel{}, err
	}

	if err := s.novels.Save(ctx, n); err != nil {
		return novel.Novel{}, err
	}
	return n, nil
}

// AddVolume adds a volume to an existing novel and persists the structure.
func (s *NovelService) AddVolume(ctx context.Context, in port.AddVolumeInput) (novel.Novel, error) {
	n, err := s.novels.FindByTitle(ctx, in.Novel)
	if err != nil {
		return novel.Novel{}, err
	}
	if err := n.AddVolume(in.Number, in.Title); err != nil {
		return novel.Novel{}, err
	}
	if err := s.novels.SaveStructure(ctx, n); err != nil {
		return novel.Novel{}, err
	}
	return n, nil
}

// AddChapter adds a chapter to a volume of an existing novel and persists it.
func (s *NovelService) AddChapter(ctx context.Context, in port.AddChapterInput) (novel.Novel, error) {
	n, err := s.novels.FindByTitle(ctx, in.Novel)
	if err != nil {
		return novel.Novel{}, err
	}
	if err := n.AddChapter(in.Volume, in.Number, in.Title); err != nil {
		return novel.Novel{}, err
	}
	if err := s.novels.SaveStructure(ctx, n); err != nil {
		return novel.Novel{}, err
	}
	return n, nil
}

// GetNovel returns a single novel with its volumes and chapters.
func (s *NovelService) GetNovel(ctx context.Context, title string) (novel.Novel, error) {
	return s.novels.FindByTitle(ctx, title)
}

// ListNovels returns all novels.
func (s *NovelService) ListNovels(ctx context.Context) ([]novel.Novel, error) {
	return s.novels.List(ctx)
}
