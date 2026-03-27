package scheduling

import (
	"time"
)

type ScheduleStatus string

const (
	StatusScheduled  ScheduleStatus = "Scheduled"
	StatusInProgress ScheduleStatus = "InProgress"
	StatusCompleted  ScheduleStatus = "Completed"
	StatusCancelled  ScheduleStatus = "Cancelled"
)

type TankStatus string

const (
	TankAvailable   TankStatus = "Available"
	TankOccupied    TankStatus = "Occupied"
	TankMaintenance TankStatus = "Maintenance"
)

type Tank struct {
	ID       string
	Name     string
	Capacity int
	Status   TankStatus
}

type Schedule struct {
	ID        string         `json:"id"`
	TankID    string         `json:"tank_id"`
	BatchID   string         `json:"batch_id"`
	Quantity  int            `json:"quantity"`
	RecipeID  string         `json:"recipe_id"`
	StartTime time.Time      `json:"start_time"`
	EndTime   time.Time      `json:"end_time"`
	Status    ScheduleStatus `json:"status"`
}
