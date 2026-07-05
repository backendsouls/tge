// Package config loads default values for the game's entities from YAML.
//
// A baseline defaults.yml is embedded in the binary so the game always works.
// Setting the TGE_DEFAULTS environment variable to a file path overrides it:
// the external file is decoded on top of the embedded defaults, so it only needs
// to specify the keys it wants to change.
//
// The defaults serve two purposes:
//   - fill-in values used when the caller does not provide one (the default
//     cosmology names, the Human base species, a new character's gender/age/stats);
//   - a starter Catalog of entity instances that are seeded into the store on
//     startup if absent.
package config

import (
	_ "embed"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// EnvPath is the environment variable naming an external defaults file that
// overrides the embedded one.
const EnvPath = "TGE_DEFAULTS"

//go:embed defaults.yml
var embedded []byte

// Stats mirrors rpg.Stats for (de)serialization.
type Stats struct {
	STR int `yaml:"str"`
	AGI int `yaml:"agi"`
	INT int `yaml:"int"`
	VIT int `yaml:"vit"`
	DEX int `yaml:"dex"`
	WIS int `yaml:"wis"`
	CHA int `yaml:"cha"`
	LUK int `yaml:"luk"`
}

// Species is a base biological template.
type Species struct {
	Name          string  `yaml:"name"`
	Power         float64 `yaml:"power"`
	Lifespan      int     `yaml:"lifespan"`
	DefaultGender string  `yaml:"default_gender"` // Male/Female, or empty for none
}

// World names the entities of the default cosmology.
type World struct {
	Reality     string `yaml:"reality"`
	Omniverse   string `yaml:"omniverse"`
	Multiverse  string `yaml:"multiverse"`
	Universe    string `yaml:"universe"`
	Realm       string `yaml:"realm"`
	PowerSystem string `yaml:"power_system"`
}

// Character holds the fill-in defaults for a new character.
type Character struct {
	Gender  string  `yaml:"gender"`
	Age     int     `yaml:"age"`
	Species Species `yaml:"species"` // the Human base template
	Stats   Stats   `yaml:"stats"`
}

// Grade defines a rarity level for items, skills, and abilities.
type Grade struct {
	Name            string  `yaml:"name"`
	Level           int     `yaml:"level"`
	PowerMultiplier float64 `yaml:"power_multiplier"`
}

// RPGDefaults holds config for the RPG mechanics.
type RPGDefaults struct {
	Grades []Grade `yaml:"grades"`
}

// Named is the common name+description shape of several catalog entities.
type Named struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Grade       string `yaml:"grade"`
}

// Level is a sub-stage of a realm. It carries both progression concepts.
type Level struct {
	Number             int    `yaml:"number"`
	Name               string `yaml:"name"`
	BreakthroughPoints int    `yaml:"breakthrough_points"`
	BottleneckPoints   int    `yaml:"bottleneck_points"`
}

// Realm is a cultivation stage with its formulas, caps and levels.
type Realm struct {
	Name                   string  `yaml:"name"`
	PowerMultiplier        float64 `yaml:"power_multiplier"`
	PowerAdder             float64 `yaml:"power_adder"`
	LifespanMultiplier     float64 `yaml:"lifespan_multiplier"`
	LifespanAdder          float64 `yaml:"lifespan_adder"`
	BottleneckPoints       int     `yaml:"bottleneck_points"` // realm-wide bottleneck
	MaxLevels              int     `yaml:"max_levels"`
	MainCharacterMaxLevels int     `yaml:"main_character_max_levels"`
	// LevelCount auto-generates this many ordinally-named levels ("First Level"
	// ..) when Levels is not given explicitly. The generated levels take their
	// points from LevelBreakthroughPoints and LevelBottleneckPoints.
	LevelCount              int     `yaml:"level_count"`
	LevelBreakthroughPoints int     `yaml:"level_breakthrough_points"`
	LevelBottleneckPoints   int     `yaml:"level_bottleneck_points"`
	Levels                  []Level `yaml:"levels"`
}

// Effect is a buff/debuff/status modifier.
type Effect struct {
	Name        string `yaml:"name"`
	Kind        string `yaml:"kind"`
	Description string `yaml:"description"`
}

// Equipment is an equippable item with a slot and stat bonus.
type Equipment struct {
	Name  string `yaml:"name"`
	Slot  string `yaml:"slot"`
	Bonus Stats  `yaml:"bonus"`
}

// Objective is a step of a quest.
type Objective struct {
	Order       int    `yaml:"order"`
	Description string `yaml:"description"`
}

// Quest is content with ordered objectives.
type Quest struct {
	Name        string      `yaml:"name"`
	Description string      `yaml:"description"`
	Objectives  []Objective `yaml:"objectives"`
}

// Ingredient is an input of a recipe.
type Ingredient struct {
	Item     string `yaml:"item"`
	Quantity int    `yaml:"quantity"`
}

// Recipe turns input ingredients into an output item.
type Recipe struct {
	Name   string       `yaml:"name"`
	Output string       `yaml:"output"`
	Inputs []Ingredient `yaml:"inputs"`
}

// Catalog is the starter content seeded into the store on startup.
type Catalog struct {
	Species     []Species   `yaml:"species"`
	Realms      []Realm     `yaml:"realms"`
	Classes     []Named     `yaml:"classes"`
	Professions []Named     `yaml:"professions"`
	Items       []Named     `yaml:"items"`
	Abilities   []Named     `yaml:"abilities"`
	Skills      []Named     `yaml:"skills"`
	Effects     []Effect    `yaml:"effects"`
	Equipment   []Equipment `yaml:"equipment"`
	Quests      []Quest     `yaml:"quests"`
	Recipes     []Recipe    `yaml:"recipes"`
}

// SystemLog configures the author-facing system messages.
type SystemLog struct {
	Template string `yaml:"template"`
}

// Defaults is the whole defaults document.
type Defaults struct {
	World     World       `yaml:"world"`
	Character Character   `yaml:"character"`
	RPG       RPGDefaults `yaml:"rpg"`
	Catalog   Catalog     `yaml:"catalog"`
	SystemLog SystemLog   `yaml:"system_log"`
}

// Load decodes the embedded defaults, then applies the TGE_DEFAULTS override
// file on top of them when that variable is set.
func Load() (Defaults, error) {
	var d Defaults
	if err := yaml.Unmarshal(embedded, &d); err != nil {
		return Defaults{}, fmt.Errorf("parse embedded defaults: %w", err)
	}
	if path := os.Getenv(EnvPath); path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return Defaults{}, fmt.Errorf("read %s=%q: %w", EnvPath, path, err)
		}
		if err := yaml.Unmarshal(data, &d); err != nil {
			return Defaults{}, fmt.Errorf("parse %q: %w", path, err)
		}
	}
	return d, nil
}
