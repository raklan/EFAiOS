package Engine

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

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

// Generates an ID for something. To ensure it's unique, I'm just using the current UNIX time in
// milliseconds with a random set of 10 characters appended to the end. Will probably need to change to something more random later
func GenerateId() string {
	const lettersAndNumbers = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 10)
	for i := range code {
		code[i] = lettersAndNumbers[rand.Intn(len(lettersAndNumbers))]
	}

	return fmt.Sprint(time.Now().UnixMilli(), string(code))
}
