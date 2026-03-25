package recipe

import (
	"time"
)

type Recipe struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Version     int           `json:"version"`
	Stages      []RecipeStage `json:"stages"`
	CreatedAt   time.Time     `json:"created_at"`
}

type RecipeStage struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Order       int                `json:"order"`
	Ingredients []RecipeIngredient `json:"ingredients"`
}

type RecipeIngredient struct {
	ID            string  `json:"id"`
	ProductID     string  `json:"product_id"`
	Quantity      float64 `json:"quantity"`
	UnitOfMeasure string  `json:"unit_of_measure"`
}
