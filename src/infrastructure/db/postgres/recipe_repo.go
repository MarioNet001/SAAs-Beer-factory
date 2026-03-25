package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sistema-gestion-beer/src/domain/recipe"
)

type recipeRepo struct {
	db *sql.DB
}

func NewRecipeRepo(db *sql.DB) recipe.RecipeRepository {
	return &recipeRepo{db: db}
}

func (r *recipeRepo) Save(ctx context.Context, recipe *recipe.Recipe) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert Recipe
	query := `INSERT INTO recipes (name, style, version, batch_size_liters) VALUES ($1, $2, $3, $4) RETURNING id`
	var recipeID string
	err = tx.QueryRowContext(ctx, query, recipe.Name, "Default Style", recipe.Version, 20.0).Scan(&recipeID)
	if err != nil {
		return fmt.Errorf("failed to insert recipe: %w", err)
	}

	// Insert Stages
	for _, stage := range recipe.Stages {
		var stageID string
		queryStage := `INSERT INTO recipe_stages (recipe_id, name, sequence_order) VALUES ($1, $2, $3) RETURNING id`
		err = tx.QueryRowContext(ctx, queryStage, recipeID, stage.Name, stage.Order).Scan(&stageID)
		if err != nil {
			return fmt.Errorf("failed to insert stage: %w", err)
		}

		// Insert Ingredients
		for _, ing := range stage.Ingredients {
			queryIng := `INSERT INTO recipe_ingredients (recipe_stage_id, inventory_product_id, quantity, unit) VALUES ($1, $2, $3, $4)`
			_, err = tx.ExecContext(ctx, queryIng, stageID, ing.ProductID, ing.Quantity, ing.UnitOfMeasure)
			if err != nil {
				return fmt.Errorf("failed to insert ingredient: %w", err)
			}
		}
	}

	return tx.Commit()
}

func (r *recipeRepo) GetByID(ctx context.Context, id string) (*recipe.Recipe, error) {
	return nil, fmt.Errorf("GetByID not implemented")
}

func (r *recipeRepo) List(ctx context.Context) ([]*recipe.Recipe, error) {
	query := `SELECT id, name, version, created_at FROM recipes`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query recipes: %w", err)
	}
	defer rows.Close()

	var recipes []*recipe.Recipe
	for rows.Next() {
		var r recipe.Recipe
		if err := rows.Scan(&r.ID, &r.Name, &r.Version, &r.CreatedAt); err != nil {
			return nil, err
		}
		recipes = append(recipes, &r)
	}
	return recipes, nil
}
