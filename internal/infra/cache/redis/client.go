package redis

import (
	"context"

	"github.com/maadiii/taskmanager/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

func NewClient(lc fx.Lifecycle, ctx context.Context, cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{Addr: cfg.Redis.Addr})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return client.Ping(ctx).Err()
		},
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return client, nil
}
