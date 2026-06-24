package mocks

import (
	"github.com/redis/go-redis/v9"
)


func NewTestRedis(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{Addr: addr})
}