package rpg

import (
	"errors"
	"strings"
)

var (
	// ErrInvalidInventoryItem is returned when an inventory item name is blank.
	ErrInvalidInventoryItem = errors.New("inventory: item must not be empty")
	// ErrInvalidQuantity is returned for a non-positive quantity.
	ErrInvalidQuantity = errors.New("inventory: quantity must be positive")
)

// ItemStack is a quantity of one item held in an Inventory.
type ItemStack struct {
	Item     string
	Quantity int
}

// Inventory is the collection of items a character carries.
type Inventory struct {
	Items []ItemStack
}

// Add puts quantity of item into the inventory, merging with an existing stack
// of the same item. It rejects a blank item or a non-positive quantity.
func (inv *Inventory) Add(item string, quantity int) error {
	item = strings.TrimSpace(item)
	if item == "" {
		return ErrInvalidInventoryItem
	}
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	for i := range inv.Items {
		if inv.Items[i].Item == item {
			inv.Items[i].Quantity += quantity
			return nil
		}
	}
	inv.Items = append(inv.Items, ItemStack{Item: item, Quantity: quantity})
	return nil
}
