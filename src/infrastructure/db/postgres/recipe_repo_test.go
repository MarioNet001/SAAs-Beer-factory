package postgres

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"os"
	"sistema-gestion-beer/src/domain/recipe"
	"testing"
)

func TestRecipeRepo_Save_Transaction(t *testing.T) {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		t.Skip("DB_URL not set, skipping integration test")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewRecipeRepo(db)

	r := &recipe.Recipe{
		Name:    "Integration Test Recipe",
		Version: 1,
		Stages: []recipe.RecipeStage{
			{
				Name:  "Stage 1",
				Order: 1,
				Ingredients: []recipe.RecipeIngredient{
					{ProductID: "550e8400-e29b-41d4-a716-446655440000", Quantity: 10.0, UnitOfMeasure: "kg"},
				},
			},
		},
	}

	err = repo.Save(context.Background(), r)
	if err != nil {
		t.Fatalf("expected no error saving recipe, got %v", err)
	}

	// Verify record exists
	var id string
	err = db.QueryRow("SELECT id FROM recipes WHERE name = $1", "Integration Test Recipe").Scan(&id)
	if err != nil {
		t.Fatalf("expected to find recipe in database, got error %v", err)
	}
}
