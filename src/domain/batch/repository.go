package batch

import (
	"context"
)

type BatchRepository interface {
	Save(ctx context.Context, batch *Batch) error
	Update(ctx context.Context, batch *Batch) error
	SaveSnapshot(ctx context.Context, snapshot *BatchRecipeSnapshot) error
	GetByID(ctx context.Context, id string) (*Batch, error)
}
