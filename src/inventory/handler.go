package inventory

import (
	"context"
	"fmt"
)

type Handler struct {
	service InventoryServiceInterface
}

func NewHandler(service InventoryServiceInterface) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleAddStock(ctx context.Context, itemID string, amount float64) error {
	err := h.service.AddStock(ctx, itemID, amount)
	if err != nil {
		return fmt.Errorf("error handling add stock: %w", err)
	}
	return nil
}
