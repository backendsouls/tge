package rpg

import (
	"errors"
	"strings"
)

// ErrInvalidItemName is returned when an item name is blank.
var ErrInvalidItemName = errors.New("item: name must not be empty")

// Item is a thing that can exist in the world and be held in an inventory.
type Item struct {
	Name        string
	Description string
	Grade       string
}

// NewItem validates and builds an item.
func NewItem(name, description, grade string) (Item, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Item{}, ErrInvalidItemName
	}
	return Item{Name: name, Description: strings.TrimSpace(description), Grade: strings.TrimSpace(grade)}, nil
}
