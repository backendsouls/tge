package novel

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidVolumeNumber is returned when a volume number is not positive.
	ErrInvalidVolumeNumber = errors.New("volume: number must be positive")
	// ErrChapterExists is returned when a chapter number is already used in the volume.
	ErrChapterExists = errors.New("volume: chapter number already exists")
)

// Volume is a volume within a novel, containing ordered chapters. Title is optional.
type Volume struct {
	Number   int
	Title    string
	Chapters []Chapter
}

// NewVolume validates and builds an (empty) volume.
func NewVolume(number int, title string) (Volume, error) {
	if number <= 0 {
		return Volume{}, ErrInvalidVolumeNumber
	}
	return Volume{Number: number, Title: strings.TrimSpace(title)}, nil
}

// addChapter appends a chapter, rejecting a duplicate number within the volume.
func (v *Volume) addChapter(number int, title string) error {
	ch, err := NewChapter(number, title)
	if err != nil {
		return err
	}
	for _, c := range v.Chapters {
		if c.Number == number {
			return fmt.Errorf("%w: %d", ErrChapterExists, number)
		}
	}
	v.Chapters = append(v.Chapters, ch)
	return nil
}
