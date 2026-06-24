package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewClient(ctx context.Context, addr string, password string) (*redis.Client, error){
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		Password: password, // "" если нет пароля
		DB: 0, // дефолтная бд
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return client, nil
}

