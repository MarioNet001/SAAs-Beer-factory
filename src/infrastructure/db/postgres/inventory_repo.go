package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sistema-gestion-beer/src/domain/inventory"
	"sistema-gestion-beer/src/domain/recipe"
)

type inventoryRepo struct {
	db *sql.DB
}

func NewInventoryRepo(db *sql.DB) *inventoryRepo {
	return &inventoryRepo{db: db}
}

// Ensure it implements both interfaces
var _ inventory.InventoryRepository = &inventoryRepo{}
var _ recipe.InventoryClient = &inventoryRepo{}

func (r *inventoryRepo) CheckIngredientAvailability(ctx context.Context, productID string, quantity float64) (bool, error) {
	inv, err := r.Get(productID)
	if err != nil {
		return false, fmt.Errorf("failed to get inventory: %w", err)
	}
	return float64(inv.StockLevel) >= quantity, nil
}

func (r *inventoryRepo) Get(productID string) (*inventory.ProductInventory, error) {
	query := `SELECT product_id, stock_level, minimum_threshold, updated_at FROM inventory WHERE product_id = $1`
	var inv inventory.ProductInventory
	err := r.db.QueryRow(query, productID).Scan(&inv.ProductID, &inv.StockLevel, &inv.MinimumThreshold, &inv.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}
	return &inv, nil
}

func (r *inventoryRepo) Update(inv *inventory.ProductInventory) error {
	query := `UPDATE inventory SET stock_level = $1, updated_at = $2 WHERE product_id = $3`
	_, err := r.db.Exec(query, inv.StockLevel, inv.UpdatedAt, inv.ProductID)
	if err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}
	return nil
}

func (r *inventoryRepo) UpdateStock(ctx context.Context, itemID string, amount float64) error {
	inv, err := r.Get(itemID)
	if err != nil {
		return err
	}
	inv.StockLevel += int(amount)
	return r.Update(inv)
}
