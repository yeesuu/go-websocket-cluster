package service

import (
	"context"
	"github.com/go-redis/redis/v8"
)
type RedisService struct {
	client *redis.Client
	ctx context.Context
}

func NewRedisService(addr string, password string, db int) *RedisService {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		Password: password,
		DB: db,
	})
	rs := &RedisService{
		client: client,
		ctx: context.Background(),
	}
	return rs
}