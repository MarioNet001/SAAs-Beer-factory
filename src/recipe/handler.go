package recipeapi

import (
	"context"
	"fmt"
	"sistema-gestion-beer/src/domain/recipe"
)

type RecipeHandler struct {
	service recipe.RecipeService
}

func NewHandler(service recipe.RecipeService) *RecipeHandler {
	return &RecipeHandler{service: service}
}

func (h *RecipeHandler) HandleCreateRecipe(ctx context.Context, r *recipe.Recipe) error {
	err := h.service.CreateRecipe(ctx, r)
	if err != nil {
		return fmt.Errorf("error handling create recipe: %w", err)
	}
	return nil
}

func (h *RecipeHandler) HandleListRecipes(ctx context.Context) ([]*recipe.Recipe, error) {
	recipes, err := h.service.ListRecipes(ctx)
	if err != nil {
		return nil, fmt.Errorf("error handling list recipes: %w", err)
	}
	return recipes, nil
}
