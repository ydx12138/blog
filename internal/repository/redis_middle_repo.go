package repository

import "github.com/redis/go-redis/v9"

type redisMiddleRepo struct {
	Re *redis.Client
}

func NewRedisMiddleRepo(redis *redis.Client) RedisMiddleRepository {
	return &redisMiddleRepo{Re: redis}
}
