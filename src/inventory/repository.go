package inventory

import (
	"context"
)

type Repository interface {
	UpdateStock(ctx context.Context, itemID string, amount float64) error
}

type inventoryRepository struct{}

func NewInventoryRepository() Repository {
	return &inventoryRepository{}
}

func (r *inventoryRepository) UpdateStock(ctx context.Context, itemID string, amount float64) error {
	// DB logic here
	return nil
}
