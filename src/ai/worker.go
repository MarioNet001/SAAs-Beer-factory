package ai

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
)

type AIWorker struct {
	redisClient *redis.Client
	aiClient    *NanoBananaClient
	db          *sql.DB
	queueName   string
}

func NewAIWorker(redisClient *redis.Client, aiClient *NanoBananaClient, db *sql.DB, queueName string) *AIWorker {
	return &AIWorker{redisClient: redisClient, aiClient: aiClient, db: db, queueName: queueName}
}

func (w *AIWorker) Start(ctx context.Context) {
	log.Println("Starting AI Worker...")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			res, err := w.redisClient.BLPop(ctx, 0, w.queueName).Result()
			if err != nil {
				log.Printf("Error popping from queue: %v", err)
				continue
			}

			batchID := res[1] // res is [listName, value]

			if _, err := w.aiClient.GenerateLabel(ctx, batchID); err != nil {
				log.Printf("Error generating label for batch %s: %v", batchID, err)
				w.updateBatchStatus(ctx, batchID, "ERROR")
				continue
			}

			w.updateBatchStatus(ctx, batchID, "LABELLED")
		}
	}
}

func (w *AIWorker) updateBatchStatus(ctx context.Context, batchID string, status string) {
	_, err := w.db.ExecContext(ctx, "UPDATE batches SET status = $1 WHERE batch_id = $2", status, batchID)
	if err != nil {
		log.Printf("Failed to update status for batch %s to %s: %v", batchID, status, err)
	}
}
