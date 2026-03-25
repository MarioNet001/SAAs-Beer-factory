package recipe

import "context"

type RecipeRepository interface {
	Save(ctx context.Context, recipe *Recipe) error
	GetByID(ctx context.Context, id string) (*Recipe, error)
	List(ctx context.Context) ([]*Recipe, error)
}
