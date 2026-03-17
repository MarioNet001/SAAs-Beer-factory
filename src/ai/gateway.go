package ai

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// AIGateway implements the AIQueue interface defined in inventory.
type AIGateway struct {
	client    *redis.Client
	queueName string
}

func NewAIGateway(client *redis.Client, queueName string) *AIGateway {
	return &AIGateway{
		client:    client,
		queueName: queueName,
	}
}

// Enqueue adds a batch ID to the Redis list for the AI worker to process.
func (g *AIGateway) Enqueue(ctx context.Context, batchID string) error {
	err := g.client.LPush(ctx, g.queueName, batchID).Err()
	if err != nil {
		return fmt.Errorf("failed to enqueue batch %s: %w", batchID, err)
	}
	return nil
}
