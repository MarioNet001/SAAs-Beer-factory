package recipe

import (
	"context"
	"testing"
)

type MockRecipeRepository struct {
	SaveFunc    func(ctx context.Context, recipe *Recipe) error
	GetByIDFunc func(ctx context.Context, id string) (*Recipe, error)
}

func (m *MockRecipeRepository) Save(ctx context.Context, recipe *Recipe) error {
	return m.SaveFunc(ctx, recipe)
}
func (m *MockRecipeRepository) GetByID(ctx context.Context, id string) (*Recipe, error) {
	return m.GetByIDFunc(ctx, id)
}
func (m *MockRecipeRepository) List(ctx context.Context) ([]*Recipe, error) {
	return nil, nil
}

type MockInventoryClient struct {
	CheckIngredientAvailabilityFunc func(ctx context.Context, productID string, quantity float64) (bool, error)
}

func (m *MockInventoryClient) CheckIngredientAvailability(ctx context.Context, productID string, quantity float64) (bool, error) {
	return m.CheckIngredientAvailabilityFunc(ctx, productID, quantity)
}

func TestCreateRecipe_Success(t *testing.T) {
	mockRepo := &MockRecipeRepository{
		SaveFunc: func(ctx context.Context, recipe *Recipe) error {
			return nil
		},
	}
	mockInv := &MockInventoryClient{
		CheckIngredientAvailabilityFunc: func(ctx context.Context, productID string, quantity float64) (bool, error) {
			return true, nil
		},
	}

	service := NewRecipeService(mockRepo, mockInv)

	recipe := &Recipe{
		Name: "Test Recipe",
		Stages: []RecipeStage{
			{
				Name:  "Boil",
				Order: 1,
				Ingredients: []RecipeIngredient{
					{ProductID: "malt", Quantity: 5.0, UnitOfMeasure: "kg"},
				},
			},
		},
	}

	err := service.CreateRecipe(context.Background(), recipe)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if recipe.Version != 1 {
		t.Errorf("expected version 1, got %d", recipe.Version)
	}
}

func TestCreateRecipe_InsufficientStock(t *testing.T) {
	mockRepo := &MockRecipeRepository{}
	mockInv := &MockInventoryClient{
		CheckIngredientAvailabilityFunc: func(ctx context.Context, productID string, quantity float64) (bool, error) {
			return false, nil
		},
	}

	service := NewRecipeService(mockRepo, mockInv)

	recipe := &Recipe{
		Name: "Test Recipe",
		Stages: []RecipeStage{
			{
				Name:  "Boil",
				Order: 1,
				Ingredients: []RecipeIngredient{
					{ProductID: "malt", Quantity: 5.0, UnitOfMeasure: "kg"},
				},
			},
		},
	}

	err := service.CreateRecipe(context.Background(), recipe)
	if err == nil {
		t.Fatal("expected error due to insufficient stock, got nil")
	}
}
