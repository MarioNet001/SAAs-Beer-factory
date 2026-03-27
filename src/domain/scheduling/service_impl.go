package scheduling

import (
	"context"
	"errors"
	"fmt"
)

type schedulingService struct {
	tankRepo     TankRepository
	scheduleRepo ScheduleRepository
	inventory    InventoryClient
	recipe       RecipeClient
}

func NewSchedulingService(tankRepo TankRepository, scheduleRepo ScheduleRepository, inventory InventoryClient, recipe RecipeClient) SchedulingService {
	return &schedulingService{
		tankRepo:     tankRepo,
		scheduleRepo: scheduleRepo,
		inventory:    inventory,
		recipe:       recipe,
	}
}

func (s *schedulingService) CreateSchedule(ctx context.Context, schedule *Schedule) error {
	tank, err := s.tankRepo.GetByID(ctx, schedule.TankID)
	if err != nil {
		return err
	}
	if tank == nil {
		return errors.New("tank not found")
	}
	if tank.Capacity < schedule.Quantity {
		return errors.New("tank capacity exceeded")
	}

	// 1. Fetch recipe
	rec, err := s.recipe.GetRecipe(ctx, schedule.RecipeID)
	if err != nil {
		return fmt.Errorf("failed to fetch recipe: %w", err)
	}

	// 2. Ingredient check
	for _, stage := range rec.Stages {
		for _, ing := range stage.Ingredients {
			available, err := s.inventory.CheckIngredientAvailability(ctx, ing.ProductID, ing.Quantity)
			if err != nil {
				return err
			}
			if !available {
				return fmt.Errorf("insufficient stock for ingredient %s", ing.ProductID)
			}
		}
	}

	// 3. Conflict detection
	existingSchedules, err := s.scheduleRepo.ListByTank(ctx, schedule.TankID)
	if err != nil {
		return err
	}

	for _, es := range existingSchedules {
		if schedule.StartTime.Before(es.EndTime) && schedule.EndTime.After(es.StartTime) {
			return errors.New("conflict detected")
		}
	}

	return s.scheduleRepo.Save(ctx, schedule)
}

func (s *schedulingService) UpdateScheduleStatus(ctx context.Context, id string, status ScheduleStatus) error {
	schedule, err := s.scheduleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	schedule.Status = status
	return s.scheduleRepo.Update(ctx, schedule)
}

func (s *schedulingService) GetSchedule(ctx context.Context, id string) (*Schedule, error) {
	return s.scheduleRepo.GetByID(ctx, id)
}

func (s *schedulingService) ListSchedulesByTank(ctx context.Context, tankID string) ([]*Schedule, error) {
	return s.scheduleRepo.ListByTank(ctx, tankID)
}

func (s *schedulingService) RegisterTank(ctx context.Context, tank *Tank) error {
	return s.tankRepo.Save(ctx, tank)
}
