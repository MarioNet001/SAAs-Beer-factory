package inventory

import (
	"context"
	"fmt"
	"time"
)

type InventoryService interface {
	AddStock(ctx context.Context, itemID string, amount float64) error
	AdjustStock(ctx context.Context, itemID string, amount float64) error
	Get(ctx context.Context, productID string) (*ProductInventory, error)
	CheckIngredientAvailability(ctx context.Context, productID string, quantity float64) (bool, error)
}

type inventoryService struct {
	repo InventoryRepository
}

func NewInventoryService(repo InventoryRepository) InventoryService {
	return &inventoryService{repo: repo}
}

func (s *inventoryService) Get(ctx context.Context, productID string) (*ProductInventory, error) {
	return s.repo.Get(productID)
}

func (s *inventoryService) AdjustStock(ctx context.Context, productID string, amount float64) error {
	inv, err := s.repo.Get(productID)
	if err != nil {
		return fmt.Errorf("failed to get inventory: %w", err)
	}

	newStock := float64(inv.StockLevel) + amount
	if newStock < 0 {
		return fmt.Errorf("insufficient stock: current %d, required adjustment %f", inv.StockLevel, amount)
	}

	inv.StockLevel = int(newStock)
	inv.UpdatedAt = time.Now()

	if err := s.repo.Update(inv); err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}

	return nil
}

func (s *inventoryService) AddStock(ctx context.Context, itemID string, amount float64) error {
	inv, err := s.repo.Get(itemID)
	if err != nil {
		return fmt.Errorf("failed to get inventory: %w", err)
	}
	inv.StockLevel += int(amount)
	return s.repo.Update(inv)
}

func (s *inventoryService) CheckIngredientAvailability(ctx context.Context, productID string, quantity float64) (bool, error) {
	inv, err := s.repo.Get(productID)
	if err != nil {
		return false, fmt.Errorf("failed to get inventory: %w", err)
	}
	return float64(inv.StockLevel) >= quantity, nil
}
