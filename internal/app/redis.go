package app

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var redisAddr = os.Getenv("REDIS_ADDR")

func RedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	ctx := context.Background()

	// Test połączenia
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	fmt.Println("✅ Connected to Redis")
	return rdb
}
