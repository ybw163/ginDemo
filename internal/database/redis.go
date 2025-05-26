package database

import (
	"context"
	"fmt"
	"gin-web-project/internal/config"
	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func ConnectRedis() error {
	redisConfig := config.Cfg.Redis
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	// 测试连接
	_, err := RedisClient.Ping(Ctx).Result()
	return err
}
