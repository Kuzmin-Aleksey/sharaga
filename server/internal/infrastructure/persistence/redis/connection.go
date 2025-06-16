package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"sharaga/internal/config"
)

type Connection struct {
	*redis.Client
}

func NewConnection(cfg config.RedisConfig) (*Connection, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Host,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &Connection{
		Client: client,
	}, nil
}
