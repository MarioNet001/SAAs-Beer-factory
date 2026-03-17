package inventory

import (
	"context"
	"sistema-gestion-beer/src/audit"
)

type InventoryServiceInterface interface {
	AddStock(ctx context.Context, itemID string, amount float64) error
}

type inventoryService struct {
	repo         Repository
	auditService audit.AuditService
}

func NewInventoryService(repo Repository, auditService audit.AuditService) InventoryServiceInterface {
	return &inventoryService{
		repo:         repo,
		auditService: auditService,
	}
}

func (s *inventoryService) AddStock(ctx context.Context, itemID string, amount float64) error {
	err := s.repo.UpdateStock(ctx, itemID, amount)
	if err != nil {
		return err
	}
	return s.auditService.LogEvent("ADD_STOCK", fmt.Sprintf("Item %s added %f", itemID, amount))
}

import "fmt"
