// Package redis internal/repository/redis/task_cache.go
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"tdl/types"
	"time"

	"github.com/go-redis/redis/v8"
)

type TaskCache struct {
	client *redis.Client
}

func NewTaskCache(client *redis.Client) *TaskCache {
	return &TaskCache{client: client}
}

func (c *TaskCache) CacheUserTasks(userID uint, tasks []types.Task) error {
	key := fmt.Sprintf("user:%d:tasks", userID)
	data, err := json.Marshal(tasks)
	if err != nil {
		return err
	}

	// 缓存1小时
	return c.client.Set(context.Background(), key, data, time.Hour).Err()
}

func (c *TaskCache) GetUserTasks(userID uint) ([]types.Task, error) {
	key := fmt.Sprintf("user:%d:tasks", userID)

	data, err := c.client.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, err
	}

	var tasks []types.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (c *TaskCache) InvalidateUserTasks(userID uint) error {
	key := fmt.Sprintf("user:%d:tasks", userID)
	return c.client.Del(context.Background(), key).Err()
}
