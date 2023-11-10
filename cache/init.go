package cache

import (
	dotenv "github.com/NateshR/Likeminds-Real-Time/environment"
	log "github.com/NateshR/Likeminds-Real-Time/logging"
	"github.com/go-redis/redis/v7"
)

// InitRedis | initialises a connection pool to redis cache
func InitRedis() *redis.Client {
	//Initializing Redis
	dsn := dotenv.GetDotEnvVar("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr: dsn,
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	log.Info("Cache(Redis): Successfully connected and pinged.")

	return client
}
