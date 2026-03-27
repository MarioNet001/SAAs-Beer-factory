package batch

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"sistema-gestion-beer/src/domain/inventory"
	"sistema-gestion-beer/src/domain/recipe"
)

type BatchService interface {
	CreateBatch(ctx context.Context, recipeID string) (*Batch, error)
	TransitionState(ctx context.Context, batchID string, newState BatchState) error
}

type batchService struct {
	repo             BatchRepository
	recipeService    recipe.RecipeService
	inventoryService inventory.InventoryService
}

func NewBatchService(repo BatchRepository, recipeService recipe.RecipeService, inventoryService inventory.InventoryService) BatchService {
	return &batchService{
		repo:             repo,
		recipeService:    recipeService,
		inventoryService: inventoryService,
	}
}

func (s *batchService) CreateBatch(ctx context.Context, recipeID string) (*Batch, error) {
	// 1. Fetch recipe
	rec, err := s.recipeService.GetRecipe(ctx, recipeID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch recipe: %w", err)
	}

	// 2. Consume ingredients
	for _, stage := range rec.Stages {
		for _, ing := range stage.Ingredients {
			err := s.inventoryService.AdjustStock(ctx, ing.ProductID, -ing.Quantity)
			if err != nil {
				return nil, fmt.Errorf("failed to consume ingredient %s: %w", ing.ProductID, err)
			}
		}
	}

	// 3. Create batch
	batch := &Batch{
		RecipeID:  recipeID,
		State:     StatePlanned,
		CreatedAt: time.Now(),
	}

	// Save
	if err := s.repo.Save(ctx, batch); err != nil {
		return nil, fmt.Errorf("failed to save batch: %w", err)
	}

	// 4. Snapshot
	recipeData, err := json.Marshal(rec)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal recipe: %w", err)
	}

	snapshot := &BatchRecipeSnapshot{
		BatchID:      batch.ID,
		RecipeID:     recipeID,
		SnapshotData: recipeData,
	}

	if err := s.repo.SaveSnapshot(ctx, snapshot); err != nil {
		return nil, fmt.Errorf("failed to save snapshot: %w", err)
	}

	return batch, nil
}

func (s *batchService) TransitionState(ctx context.Context, batchID string, newState BatchState) error {
	batch, err := s.repo.GetByID(ctx, batchID)
	if err != nil {
		return fmt.Errorf("failed to fetch batch: %w", err)
	}

	if !isValidTransition(batch.State, newState) {
		return fmt.Errorf("invalid transition from %s to %s", batch.State, newState)
	}

	batch.State = newState
	return s.repo.Update(ctx, batch)
}

func isValidTransition(currentState, newState BatchState) bool {
	if currentState == newState {
		return true
	}

	// Allow cancelling from any state except completed (or maybe from planned/brewing)
	if newState == StateCancelled {
		return currentState != StateCompleted
	}

	allowed := map[BatchState][]BatchState{
		StatePlanned:    {StateBrewing},
		StateBrewing:    {StateFermenting},
		StateFermenting: {StateMaturation},
		StateMaturation: {StatePackaging},
		StatePackaging:  {StateCompleted},
		StateCompleted:  {},
		StateCancelled:  {},
	}

	for _, allowedState := range allowed[currentState] {
		if allowedState == newState {
			return true
		}
	}
	return false
}
