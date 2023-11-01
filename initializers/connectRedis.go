package initializers

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	ctx         context.Context
)

// Establishes a connection between the Redis Server and the application.
func ConnectRedis(config *Config) {
	ctx = context.TODO()

	RedisClient = redis.NewClient(&redis.Options{
		Addr: config.RedisUri,
	})

	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	err := RedisClient.Set(ctx, "test", "How to Refresh Access Tokens the right way in Golang", 0).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ Redis client connected successfully...")
}
