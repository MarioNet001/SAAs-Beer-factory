package recipe

import (
	"context"
)

type InventoryClient interface {
	CheckIngredientAvailability(ctx context.Context, productID string, quantity float64) (bool, error)
}
