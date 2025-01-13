package Engine

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	ModuleLogPrefix  = "ESCAPE-ENGINE"
	PackageLogPrefix = "Engine"
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

// Generates a random 4-character room code. For use in creating lobbies
func generateRoomCode() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	code := make([]byte, 4)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}

	return string(code)
}

// Logs the error in the format of "[funcLogPrefix] ERROR! [err]"
func LogError(funcLogPrefix string, err error) {
	log.Printf("%s ERROR! %s", funcLogPrefix, err)
}
