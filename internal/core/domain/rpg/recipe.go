package rpg

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidRecipeName is returned when a recipe name is blank.
	ErrInvalidRecipeName = errors.New("recipe: name must not be empty")
	// ErrInvalidRecipeOutput is returned when a recipe output item is blank.
	ErrInvalidRecipeOutput = errors.New("recipe: output item must not be empty")
	// ErrInvalidIngredientItem is returned when an ingredient item is blank.
	ErrInvalidIngredientItem = errors.New("recipe: ingredient item must not be empty")
	// ErrInvalidIngredientQuantity is returned for a non-positive quantity.
	ErrInvalidIngredientQuantity = errors.New("recipe: ingredient quantity must be positive")
	// ErrIngredientExists is returned when an ingredient item is added twice.
	ErrIngredientExists = errors.New("recipe: ingredient already in this recipe")
)

// Ingredient is a quantity of an input item consumed by a Recipe.
type Ingredient struct {
	Item     string
	Quantity int
}

// Recipe turns input ingredients into an output item.
type Recipe struct {
	Name   string
	Output string
	Inputs []Ingredient
}

// NewRecipe validates and builds a recipe with no inputs.
func NewRecipe(name, output string) (Recipe, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Recipe{}, ErrInvalidRecipeName
	}
	output = strings.TrimSpace(output)
	if output == "" {
		return Recipe{}, ErrInvalidRecipeOutput
	}
	return Recipe{Name: name, Output: output}, nil
}

// AddInput adds an input ingredient, rejecting a blank item, a non-positive
// quantity, or a duplicate item.
func (r *Recipe) AddInput(item string, quantity int) error {
	item = strings.TrimSpace(item)
	if item == "" {
		return ErrInvalidIngredientItem
	}
	if quantity <= 0 {
		return ErrInvalidIngredientQuantity
	}
	for _, in := range r.Inputs {
		if in.Item == item {
			return fmt.Errorf("%w: %q", ErrIngredientExists, item)
		}
	}
	r.Inputs = append(r.Inputs, Ingredient{Item: item, Quantity: quantity})
	return nil
}
