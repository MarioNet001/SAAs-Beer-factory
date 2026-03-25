package inventory

import (
	"fmt"
	"time"
)

type InventoryService interface {
	Get(productID string) (*ProductInventory, error)
	AdjustStock(productID string, amount int) error
}

type inventoryService struct {
	repo InventoryRepository
}

func NewInventoryService(repo InventoryRepository) InventoryService {
	return &inventoryService{repo: repo}
}

func (s *inventoryService) Get(productID string) (*ProductInventory, error) {
	return s.repo.Get(productID)
}

func (s *inventoryService) AdjustStock(productID string, amount int) error {
	inv, err := s.repo.Get(productID)
	if err != nil {
		return fmt.Errorf("failed to get inventory: %w", err)
	}

	newStock := inv.StockLevel + amount
	if newStock < 0 {
		return fmt.Errorf("insufficient stock: current %d, required adjustment %d", inv.StockLevel, amount)
	}

	inv.StockLevel = newStock
	inv.UpdatedAt = time.Now()

	if err := s.repo.Update(inv); err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}

	if inv.StockLevel < inv.MinimumThreshold {
		// Log or trigger alert
	}

	return nil
}
