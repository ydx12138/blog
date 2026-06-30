package core

import (
	"blog/config"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var RDB *redis.Client

func InitRedis() error {
	redis1 := config.Cfg.Redis
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redis1.Host, redis1.Port),
		Password: redis1.Password,
		DB:       redis1.DB,

		// 连接池（很重要）
		PoolSize:     redis1.PoolSize,
		MinIdleConns: 5,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 测试连接
	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("redis connect failed: %w", err)
	}

	fmt.Println("Redis connected successfully")
	return nil
}
