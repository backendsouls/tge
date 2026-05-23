package novel

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidNovelTitle is returned when a novel title is blank.
	ErrInvalidNovelTitle = errors.New("novel: title must not be empty")
	// ErrNovelMainCharacterRequired is returned when a novel has no main character.
	ErrNovelMainCharacterRequired = errors.New("novel: a main character is required")
	// ErrVolumeExists is returned when a volume number is already used in the novel.
	ErrVolumeExists = errors.New("novel: volume number already exists")
	// ErrVolumeNotFound is returned when the named volume does not exist.
	ErrVolumeNotFound = errors.New("novel: volume not found")
)

// Novel is a story with a single associated main character (referenced by name)
// and ordered volumes, each containing chapters.
type Novel struct {
	Title         string
	MainCharacter string
	Volumes       []Volume
}

// AddVolume appends a volume, rejecting a duplicate number within the novel.
func (n *Novel) AddVolume(number int, title string) error {
	v, err := NewVolume(number, title)
	if err != nil {
		return err
	}
	for _, ex := range n.Volumes {
		if ex.Number == number {
			return fmt.Errorf("%w: %d", ErrVolumeExists, number)
		}
	}
	n.Volumes = append(n.Volumes, v)
	return nil
}

// AddChapter appends a chapter to the named volume.
func (n *Novel) AddChapter(volumeNumber, chapterNumber int, title string) error {
	for i := range n.Volumes {
		if n.Volumes[i].Number == volumeNumber {
			return n.Volumes[i].addChapter(chapterNumber, title)
		}
	}
	return fmt.Errorf("%w: %d", ErrVolumeNotFound, volumeNumber)
}

// NewNovel validates and builds a novel. Whether the referenced character exists
// and is actually a main character is enforced by the application service.
func NewNovel(title, mainCharacter string) (Novel, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Novel{}, ErrInvalidNovelTitle
	}
	mainCharacter = strings.TrimSpace(mainCharacter)
	if mainCharacter == "" {
		return Novel{}, ErrNovelMainCharacterRequired
	}
	return Novel{Title: title, MainCharacter: mainCharacter}, nil
}
