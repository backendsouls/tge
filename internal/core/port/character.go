package port

import (
	"context"
	"errors"
	"tge/internal/core/domain/character"
)

var (
	// ErrCharacterNotFound is returned when a lookup finds no matching character.
	ErrCharacterNotFound = errors.New("character: not found")
	// ErrCharacterExists is returned when saving a character whose name is taken.
	ErrCharacterExists = errors.New("character: already exists")
)

// CreateCharacterInput describes the mortal character to create. Type and Gender
// are validated by the domain; each named system must already exist, as must a
// given Class or Profession (both optional).
type CreateCharacterInput struct {
	Name       string
	Type       string
	Gender     string
	Species    string
	Systems    []string
	Class      string
	Profession string
}

// GiveItemInput adds a quantity of an item to a character's inventory.
type GiveItemInput struct {
	Character string
	Item      string
	Quantity  int
}

// CultivateInput sets (or replaces) a character's cultivation state at one power
// node. System defaults to the character's first power system and Path to the
// System name; the state is anchored to the given Realm and Level.
type CultivateInput struct {
	Character   string
	System      string
	Path        string
	Realm       string
	LevelNumber int
	LevelName   string
	Points      int
	Progress    float64
}

// CultivationRecord is one persisted cultivation entry for a character: the
// progress at a single (System, Path) node.
type CultivationRecord struct {
	System      string
	Path        string
	Realm       string
	LevelNumber int
	LevelName   string
	Points      int
	Progress    float64
}

// CharacterRepository is a driven port persisting characters keyed by name.
type CharacterRepository interface {
	Save(ctx context.Context, c character.Character) error
	FindByName(ctx context.Context, name string) (character.Character, error)
	// MainCharacters returns every main-character-typed character, ordered by name.
	MainCharacters(ctx context.Context) ([]character.Character, error)
	List(ctx context.Context) ([]character.Character, error)
	// AddItem adds quantity of an item to a character's inventory, merging with
	// any existing stack of the same item.
	AddItem(ctx context.Context, character, item string, quantity int) error
	// SaveCultivation sets (upserts) a character's cultivation state at one
	// (System, Path) node.
	SaveCultivation(ctx context.Context, character string, rec CultivationRecord) error
}

// CharacterService is a driving port for character use cases.
type CharacterService interface {
	CreateCharacter(ctx context.Context, in CreateCharacterInput) (character.Character, error)
	MainCharacter(ctx context.Context) (character.Character, error)
	Character(ctx context.Context, name string) (character.Character, error)
	ListCharacters(ctx context.Context) ([]character.Character, error)
	GiveItem(ctx context.Context, in GiveItemInput) (character.Character, error)
	// Cultivate sets a character's cultivation state at a power node and returns
	// the updated character.
	Cultivate(ctx context.Context, in CultivateInput) (character.Character, error)
}
