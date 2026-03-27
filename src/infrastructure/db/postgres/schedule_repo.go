package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sistema-gestion-beer/src/domain/scheduling"
)

type sqlScheduleRepository struct {
	db *sql.DB
}

func NewScheduleRepo(db *sql.DB) scheduling.ScheduleRepository {
	return &sqlScheduleRepository{db: db}
}

func (r *sqlScheduleRepository) Save(ctx context.Context, s *scheduling.Schedule) error {
	// Handle empty string as NULL for batch_id
	var batchID *string
	if s.BatchID != "" {
		batchID = &s.BatchID
	}

	query := `INSERT INTO schedules (tank_id, batch_id, quantity, start_time, end_time, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, s.TankID, batchID, s.Quantity, s.StartTime, s.EndTime, s.Status).Scan(&s.ID)
	if err != nil {
		return fmt.Errorf("failed to insert schedule: %w", err)
	}
	return nil
}

func (r *sqlScheduleRepository) Update(ctx context.Context, s *scheduling.Schedule) error {
	query := `UPDATE schedules SET status = $1, start_time = $2, end_time = $3, quantity = $4 WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query, s.Status, s.StartTime, s.EndTime, s.Quantity, s.ID)
	if err != nil {
		return fmt.Errorf("failed to update schedule: %w", err)
	}
	return nil
}

func (r *sqlScheduleRepository) GetByID(ctx context.Context, id string) (*scheduling.Schedule, error) {
	query := `SELECT id, tank_id, batch_id, quantity, start_time, end_time, status FROM schedules WHERE id = $1`
	var s scheduling.Schedule
	err := r.db.QueryRowContext(ctx, query, id).Scan(&s.ID, &s.TankID, &s.BatchID, &s.Quantity, &s.StartTime, &s.EndTime, &s.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}
	return &s, nil
}

func (r *sqlScheduleRepository) ListByTank(ctx context.Context, tankID string) ([]*scheduling.Schedule, error) {
	query := `SELECT id, tank_id, batch_id, quantity, start_time, end_time, status FROM schedules WHERE tank_id = $1`
	rows, err := r.db.QueryContext(ctx, query, tankID)
	if err != nil {
		return nil, fmt.Errorf("failed to list schedules: %w", err)
	}
	defer rows.Close()

	var schedules []*scheduling.Schedule
	for rows.Next() {
		var s scheduling.Schedule
		if err := rows.Scan(&s.ID, &s.TankID, &s.BatchID, &s.Quantity, &s.StartTime, &s.EndTime, &s.Status); err != nil {
			return nil, fmt.Errorf("failed to scan schedule: %w", err)
		}
		schedules = append(schedules, &s)
	}
	return schedules, nil
}
