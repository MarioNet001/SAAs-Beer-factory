package batchapi

import (
	"context"
	"sistema-gestion-beer/src/domain/batch"
)

type BatchHandler struct {
	service batch.BatchService
}

func NewHandler(service batch.BatchService) *BatchHandler {
	return &BatchHandler{service: service}
}

func (h *BatchHandler) HandleCreateBatch(ctx context.Context, recipeID string) (*batch.Batch, error) {
	return h.service.CreateBatch(ctx, recipeID)
}

func (h *BatchHandler) HandleTransitionState(ctx context.Context, batchID string, newState batch.BatchState) error {
	return h.service.TransitionState(ctx, batchID, newState)
}
