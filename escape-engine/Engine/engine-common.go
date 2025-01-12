package Engine

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

func getRedisAddress() string {
	environ := os.Getenv("REDIS_ADDRESS")
	if environ == "" {
		return "localhost:6379"
	}
	return environ
}

var RDB = redis.NewClient(&redis.Options{
	Addr: getRedisAddress(),
})

var ctx = context.Background()
