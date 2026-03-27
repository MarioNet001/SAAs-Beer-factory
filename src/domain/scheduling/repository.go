package scheduling

import (
	"context"
)

type TankRepository interface {
	Save(ctx context.Context, tank *Tank) error
	GetByID(ctx context.Context, id string) (*Tank, error)
	ListAll(ctx context.Context) ([]*Tank, error)
}

type ScheduleRepository interface {
	Save(ctx context.Context, schedule *Schedule) error
	Update(ctx context.Context, schedule *Schedule) error
	GetByID(ctx context.Context, id string) (*Schedule, error)
	ListByTank(ctx context.Context, tankID string) ([]*Schedule, error)
}
