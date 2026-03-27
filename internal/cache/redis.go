package cache

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var Client *redis.Client

func Connect() {
	Client = redis.NewClient(&redis.Options{
		Addr: getEnv("REDIS_HOST", "localhost") + ":" + getEnv("REDIS_PORT", "6379"),
	})
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
