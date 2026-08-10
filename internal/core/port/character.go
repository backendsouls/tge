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

// CreateCharacterInput describes the mortal character to create.
type CreateCharacterInput struct {
	Name       string
	Type       string
	Gender     string
	Species    string
	Class      string
	Profession string
	Systems    []string
}

// GiveItemInput adds a quantity of an item to a character's inventory.
type GiveItemInput struct {
	Character string
	Item      string
	Quantity  int
}

// TrainNodeInput requests progression on a specific PowerNode in a PowerSystem.
type TrainNodeInput struct {
	Character string
	System    string
	NodeID    string
}

// AddPowerInput requests adding a power system to a character.
type AddPowerInput struct {
	Character string
	System    string
}

// CharacterRepository is a driven port persisting characters keyed by name.
// Thanks to aggregate-root serialization (Flat Files/JSON), we only need high-level Save/Find operations.
type CharacterRepository interface {
	Save(ctx context.Context, c character.Character) error
	FindByName(ctx context.Context, name string) (character.Character, error)
	// MainCharacters returns every main-character-typed character, ordered by name.
	MainCharacters(ctx context.Context) ([]character.Character, error)
	List(ctx context.Context) ([]character.Character, error)
	// Clean soft deletes all characters.
	Clean(ctx context.Context) error
}

// CharacterService is a driving port for character use cases.
type CharacterService interface {
	CreateCharacter(ctx context.Context, in CreateCharacterInput) (character.Character, error)
	MainCharacter(ctx context.Context) (character.Character, error)
	Character(ctx context.Context, name string) (character.Character, error)
	ListCharacters(ctx context.Context) ([]character.Character, error)
	CleanCharacters(ctx context.Context) error
	GiveItem(ctx context.Context, in GiveItemInput) (character.Character, error)

	// TrainNode attempts to progress or unlock a node within a character's power system graph.
	TrainNode(ctx context.Context, in TrainNodeInput) (character.Character, error)
	AddPower(ctx context.Context, in AddPowerInput) (character.Character, error)
	AssignIdleActivity(ctx context.Context, charName string, systemName string, powerName string, duration float64) (character.Character, error)
	PassTime(ctx context.Context, charName string, seconds int64) (character.Character, error)
}
