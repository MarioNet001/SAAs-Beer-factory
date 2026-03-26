package batch

import (
	"time"
)

type BatchState string

const (
	StatePlanned    BatchState = "Planned"
	StateBrewing    BatchState = "Brewing"
	StateFermenting BatchState = "Fermenting"
	StateMaturation BatchState = "Maturation"
	StatePackaging  BatchState = "Packaging"
	StateCompleted  BatchState = "Completed"
	StateCancelled  BatchState = "Cancelled"
)

type Batch struct {
	ID        string
	RecipeID  string
	State     BatchState
	CreatedAt time.Time
}

type BatchRecipeSnapshot struct {
	ID           string
	BatchID      string
	RecipeID     string
	SnapshotData []byte
}
