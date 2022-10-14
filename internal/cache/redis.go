package cache

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
)

var (
	// RedisClient export to communicate with Redis Client
	RedisClient *redis.Client
)

// InitRedisClient init Redis Client
func InitRedisClient(redisHost string, redisPort int) error {
	context := context.Background()
	redisAddr := fmt.Sprintf("%s:%d", redisHost, redisPort)
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := RedisClient.Ping(context).Result()
	if err != nil {
		return err
	}

	fmt.Printf("🗂 redis connected on %s. %s received.\n", redisAddr, pong)
	return nil
}
