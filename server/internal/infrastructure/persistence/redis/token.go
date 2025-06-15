package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"sharaga/pkg/failure"
	"strconv"
	"time"
)

func refreshTokenKey(token string) string {
	return "refresh_token_" + token
}

func (c *Connection) SaveRefreshToken(ctx context.Context, refreshToken string, userId int, ttl time.Duration) error {
	if err := c.Set(ctx, refreshTokenKey(refreshToken), userId, ttl).Err(); err != nil {
		return failure.NewInternalError(err.Error())
	}
	return nil
}

func (c *Connection) GetRefreshToken(ctx context.Context, refreshToken string) (int, error) {
	res := c.Get(ctx, refreshTokenKey(refreshToken))
	if err := res.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, failure.NewInternalError(err.Error())
	}

	v, _ := strconv.Atoi(res.Val())

	return v, nil
}

func (c *Connection) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	if err := c.Del(ctx, refreshTokenKey(refreshToken)).Err(); err != nil {
		return failure.NewInternalError(err.Error())
	}
	return nil
}
