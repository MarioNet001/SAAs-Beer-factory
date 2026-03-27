package scheduling

import (
	"context"
	"sistema-gestion-beer/src/domain/recipe"
)

type InventoryClient interface {
	CheckIngredientAvailability(ctx context.Context, productID string, quantity float64) (bool, error)
}

type RecipeClient interface {
	GetRecipe(ctx context.Context, id string) (*recipe.Recipe, error)
}

type SchedulingService interface {
	CreateSchedule(ctx context.Context, schedule *Schedule) error
	UpdateScheduleStatus(ctx context.Context, id string, status ScheduleStatus) error
	GetSchedule(ctx context.Context, id string) (*Schedule, error)
	ListSchedulesByTank(ctx context.Context, tankID string) ([]*Schedule, error)
	RegisterTank(ctx context.Context, tank *Tank) error
}
