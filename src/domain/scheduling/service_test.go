package scheduling

import (
	"context"
	"sistema-gestion-beer/src/domain/recipe"
	"testing"
	"time"
)

type MockTankRepository struct {
	SaveFunc    func(ctx context.Context, tank *Tank) error
	GetByIDFunc func(ctx context.Context, id string) (*Tank, error)
	ListAllFunc func(ctx context.Context) ([]*Tank, error)
}

func (m *MockTankRepository) Save(ctx context.Context, tank *Tank) error {
	return m.SaveFunc(ctx, tank)
}
func (m *MockTankRepository) GetByID(ctx context.Context, id string) (*Tank, error) {
	return m.GetByIDFunc(ctx, id)
}
func (m *MockTankRepository) ListAll(ctx context.Context) ([]*Tank, error) { return m.ListAllFunc(ctx) }

type MockScheduleRepository struct {
	SaveFunc       func(ctx context.Context, schedule *Schedule) error
	UpdateFunc     func(ctx context.Context, schedule *Schedule) error
	GetByIDFunc    func(ctx context.Context, id string) (*Schedule, error)
	ListByTankFunc func(ctx context.Context, tankID string) ([]*Schedule, error)
}

func (m *MockScheduleRepository) Save(ctx context.Context, schedule *Schedule) error {
	return m.SaveFunc(ctx, schedule)
}
func (m *MockScheduleRepository) Update(ctx context.Context, schedule *Schedule) error {
	return m.UpdateFunc(ctx, schedule)
}
func (m *MockScheduleRepository) GetByID(ctx context.Context, id string) (*Schedule, error) {
	return m.GetByIDFunc(ctx, id)
}
func (m *MockScheduleRepository) ListByTank(ctx context.Context, tankID string) ([]*Schedule, error) {
	return m.ListByTankFunc(ctx, tankID)
}

type MockInventoryClient struct {
	CheckIngredientAvailabilityFunc func(ctx context.Context, productID string, quantity float64) (bool, error)
}

func (m *MockInventoryClient) CheckIngredientAvailability(ctx context.Context, productID string, quantity float64) (bool, error) {
	return m.CheckIngredientAvailabilityFunc(ctx, productID, quantity)
}

type MockRecipeClient struct {
	GetRecipeFunc func(ctx context.Context, id string) (*recipe.Recipe, error)
}

func (m *MockRecipeClient) GetRecipe(ctx context.Context, id string) (*recipe.Recipe, error) {
	return m.GetRecipeFunc(ctx, id)
}
func TestCreateSchedule(t *testing.T) {
	ctx := context.Background()
	tankID := "tank-1"

	mockInventory := &MockInventoryClient{
		CheckIngredientAvailabilityFunc: func(ctx context.Context, productID string, quantity float64) (bool, error) {
			return true, nil
		},
	}
	mockRecipe := &MockRecipeClient{
		GetRecipeFunc: func(ctx context.Context, id string) (*recipe.Recipe, error) {
			return &recipe.Recipe{Stages: []recipe.RecipeStage{}}, nil
		},
	}

	t.Run("Happy path - Successful scheduling", func(t *testing.T) {
		mockTankRepo := &MockTankRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*Tank, error) {
				return &Tank{ID: tankID, Capacity: 100, Status: TankAvailable}, nil
			},
		}
		mockScheduleRepo := &MockScheduleRepository{
			SaveFunc: func(ctx context.Context, schedule *Schedule) error {
				return nil
			},
			ListByTankFunc: func(ctx context.Context, tankID string) ([]*Schedule, error) {
				return []*Schedule{}, nil
			},
		}

		service := NewSchedulingService(mockTankRepo, mockScheduleRepo, mockInventory, mockRecipe)

		schedule := &Schedule{
			TankID:    tankID,
			Quantity:  50,
			RecipeID:  "recipe-1",
			StartTime: time.Now(),
			EndTime:   time.Now().Add(1 * time.Hour),
		}

		err := service.CreateSchedule(ctx, schedule)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Conflict detection - Scheduling two batches in same time", func(t *testing.T) {
		existingSchedule := &Schedule{
			TankID:    tankID,
			StartTime: time.Now(),
			EndTime:   time.Now().Add(1 * time.Hour),
		}

		mockTankRepo := &MockTankRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*Tank, error) {
				return &Tank{ID: tankID, Capacity: 100, Status: TankAvailable}, nil
			},
		}
		mockScheduleRepo := &MockScheduleRepository{
			ListByTankFunc: func(ctx context.Context, tankID string) ([]*Schedule, error) {
				return []*Schedule{existingSchedule}, nil
			},
			SaveFunc: func(ctx context.Context, schedule *Schedule) error {
				return nil
			},
		}

		service := NewSchedulingService(mockTankRepo, mockScheduleRepo, mockInventory, mockRecipe)

		newSchedule := &Schedule{
			TankID:    tankID,
			Quantity:  50,
			RecipeID:  "recipe-1",
			StartTime: time.Now().Add(30 * time.Minute),
			EndTime:   time.Now().Add(1 * time.Hour).Add(30 * time.Minute),
		}

		err := service.CreateSchedule(ctx, newSchedule)
		if err == nil {
			t.Errorf("Expected conflict error, got nil")
		}
	})
}
