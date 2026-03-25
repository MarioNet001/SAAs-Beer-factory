package inventory

import (
	"context"
	"database/sql"
	"fmt"
	"sistema-gestion-beer/src/audit"
)

type InventoryServiceInterface interface {
	AddStock(ctx context.Context, itemID string, amount float64) error
	AdjustStock(ctx context.Context, itemID string, amount float64) error
}

type inventoryService struct {
	repo         Repository
	auditService audit.AuditService
}

func NewInventoryService(repo Repository, db *sql.DB) InventoryServiceInterface {
	return &inventoryService{
		repo:         repo,
		auditService: audit.NewAuditService(db),
	}
}

func (s *inventoryService) AddStock(ctx context.Context, itemID string, amount float64) error {
	err := s.repo.UpdateStock(ctx, itemID, amount)
	if err != nil {
		return err
	}
	return s.auditService.LogEvent("ADD_STOCK", fmt.Sprintf("Item %s added %f", itemID, amount))
}

func (s *inventoryService) AdjustStock(ctx context.Context, itemID string, amount float64) error {
	err := s.repo.UpdateStock(ctx, itemID, amount)
	if err != nil {
		return err
	}
	return s.auditService.LogEvent("ADJUST_STOCK", fmt.Sprintf("Item %s adjusted by %f", itemID, amount))
}
