package recipe

import (
	"context"
	"fmt"
	"time"
)

type RecipeService interface {
	CreateRecipe(ctx context.Context, recipe *Recipe) error
	GetRecipe(ctx context.Context, id string) (*Recipe, error)
	ListRecipes(ctx context.Context) ([]*Recipe, error)
}

type recipeService struct {
	repo      RecipeRepository
	inventory InventoryClient
}

func NewRecipeService(repo RecipeRepository, inventory InventoryClient) RecipeService {
	return &recipeService{repo: repo, inventory: inventory}
}

func (s *recipeService) CreateRecipe(ctx context.Context, recipe *Recipe) error {
	// 1. Verify ingredients via InventoryClient
	for _, stage := range recipe.Stages {
		for _, ingredient := range stage.Ingredients {
			available, err := s.inventory.CheckIngredientAvailability(ctx, ingredient.ProductID, ingredient.Quantity)
			if err != nil {
				return fmt.Errorf("failed to check inventory for product %s: %w", ingredient.ProductID, err)
			}
			if !available {
				return fmt.Errorf("insufficient stock for product %s", ingredient.ProductID)
			}
		}
	}

	// 2. Versioning logic
	if recipe.ID == "" {
		recipe.Version = 1
		recipe.CreatedAt = time.Now()
	} else {
		existing, err := s.repo.GetByID(ctx, recipe.ID)
		if err != nil {
			return fmt.Errorf("failed to fetch existing recipe: %w", err)
		}

		recipe.Version = existing.Version + 1
		recipe.CreatedAt = time.Now()
		// Ensure it's treated as a new record
		recipe.ID = ""
	}

	// 3. Save
	if err := s.repo.Save(ctx, recipe); err != nil {
		return fmt.Errorf("failed to save recipe: %w", err)
	}

	return nil
}

func (s *recipeService) GetRecipe(ctx context.Context, id string) (*Recipe, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *recipeService) ListRecipes(ctx context.Context) ([]*Recipe, error) {
	return s.repo.List(ctx)
}
