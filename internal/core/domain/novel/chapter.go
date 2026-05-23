package novel

import (
	"errors"
	"strings"
)

// ErrInvalidChapterNumber is returned when a chapter number is not positive.
var ErrInvalidChapterNumber = errors.New("chapter: number must be positive")

// Chapter is a single chapter within a volume. Title is optional.
type Chapter struct {
	Number int
	Title  string
}

// NewChapter validates and builds a chapter.
func NewChapter(number int, title string) (Chapter, error) {
	if number <= 0 {
		return Chapter{}, ErrInvalidChapterNumber
	}
	return Chapter{Number: number, Title: strings.TrimSpace(title)}, nil
}
