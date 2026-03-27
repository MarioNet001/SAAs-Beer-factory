package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sistema-gestion-beer/src/domain/scheduling"
)

type sqlTankRepository struct {
	db *sql.DB
}

func NewTankRepo(db *sql.DB) scheduling.TankRepository {
	return &sqlTankRepository{db: db}
}

func (r *sqlTankRepository) Save(ctx context.Context, t *scheduling.Tank) error {
	query := `INSERT INTO tanks (name, capacity, status) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, t.Name, t.Capacity, t.Status).Scan(&t.ID)
	if err != nil {
		return fmt.Errorf("failed to insert tank: %w", err)
	}
	return nil
}

func (r *sqlTankRepository) GetByID(ctx context.Context, id string) (*scheduling.Tank, error) {
	query := `SELECT id, name, capacity, status FROM tanks WHERE id = $1`
	var t scheduling.Tank
	err := r.db.QueryRowContext(ctx, query, id).Scan(&t.ID, &t.Name, &t.Capacity, &t.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to get tank: %w", err)
	}
	return &t, nil
}

func (r *sqlTankRepository) ListAll(ctx context.Context) ([]*scheduling.Tank, error) {
	query := `SELECT id, name, capacity, status FROM tanks`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list tanks: %w", err)
	}
	defer rows.Close()

	var tanks []*scheduling.Tank
	for rows.Next() {
		var t scheduling.Tank
		if err := rows.Scan(&t.ID, &t.Name, &t.Capacity, &t.Status); err != nil {
			return nil, fmt.Errorf("failed to scan tank: %w", err)
		}
		tanks = append(tanks, &t)
	}
	return tanks, nil
}
