package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/maadiii/taskmanager/internal/app/dto"
	apperrors "github.com/maadiii/taskmanager/pkg/errors"
	"github.com/redis/go-redis/v9"
)

const taskTTL = 5 * time.Minute

type TaskCache struct {
	client *redis.Client
}

func NewTaskCache(client *redis.Client) *TaskCache {
	return &TaskCache{client: client}
}

func (c *TaskCache) Get(ctx context.Context, userID, taskID string) (*dto.GetByIdRs, error) {
	value, err := c.client.Get(ctx, taskKey(userID, taskID)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	response := new(dto.GetByIdRs)
	if err := json.Unmarshal([]byte(value), response); err != nil {
		return nil, apperrors.Cache(err, "decode task")
	}

	return response, nil
}

func (c *TaskCache) Set(ctx context.Context, userID, taskID string, task *dto.GetByIdRs) error {
	value, err := json.Marshal(task)
	if err != nil {
		return apperrors.Cache(err, "encode task")
	}

	return c.client.Set(ctx, taskKey(userID, taskID), value, taskTTL).Err()
}

func (c *TaskCache) Delete(ctx context.Context, userID, taskID string) error {
	return c.client.Del(ctx, taskKey(userID, taskID)).Err()
}

func taskKey(userID, taskID string) string {
	return fmt.Sprintf("task:%s:%s", userID, taskID)
}
