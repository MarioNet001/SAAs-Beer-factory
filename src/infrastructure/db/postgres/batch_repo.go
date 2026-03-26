package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sistema-gestion-beer/src/domain/batch"
)

type sqlBatchRepository struct {
	db *sql.DB
}

func NewBatchRepo(db *sql.DB) batch.BatchRepository {
	return &sqlBatchRepository{db: db}
}

func (r *sqlBatchRepository) Save(ctx context.Context, b *batch.Batch) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO batches (recipe_id, state, created_at) VALUES ($1, $2, $3) RETURNING id`
	err = tx.QueryRowContext(ctx, query, b.RecipeID, b.State, b.CreatedAt).Scan(&b.ID)
	if err != nil {
		return fmt.Errorf("failed to insert batch: %w", err)
	}

	return tx.Commit()
}

func (r *sqlBatchRepository) Update(ctx context.Context, b *batch.Batch) error {
	query := `UPDATE batches SET state = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, b.State, b.ID)
	if err != nil {
		return fmt.Errorf("failed to update batch: %w", err)
	}
	return nil
}

func (r *sqlBatchRepository) SaveSnapshot(ctx context.Context, s *batch.BatchRecipeSnapshot) error {
	query := `INSERT INTO batch_recipe_snapshots (batch_id, recipe_id, snapshot_data) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, s.BatchID, s.RecipeID, s.SnapshotData).Scan(&s.ID)
	if err != nil {
		return fmt.Errorf("failed to insert snapshot: %w", err)
	}
	return nil
}

func (r *sqlBatchRepository) GetByID(ctx context.Context, id string) (*batch.Batch, error) {
	query := `SELECT id, recipe_id, state, created_at FROM batches WHERE id = $1`
	var b batch.Batch
	err := r.db.QueryRowContext(ctx, query, id).Scan(&b.ID, &b.RecipeID, &b.State, &b.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get batch: %w", err)
	}
	return &b, nil
}
