package pubsub

import (
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"likeminds-pandemonium/common"
	"time"
)

// InitRedisClient creates a new PublishWithMethod Client
func InitRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     common.GoDotEnvVariable("REDIS_DSN"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

// SaveHashSet saves any value in an HSet with a specified TTL (Time to Live)
func SaveHashSet(redisClient *redis.Client, hashKey string, field interface{}, value interface{}, ttl time.Duration) error {
	// Marshal the struct into a JSON format to save in Redis HSET
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf(common.ErrorMarshalErrorJson, err)
	}

	// Save the value in Redis using HSet
	if err := redisClient.HSet(ctx, hashKey, field, valueBytes).Err(); err != nil {
		return fmt.Errorf(common.ErrorFailedCacheSaveRedis, err)
	}

	// Set the expiration time (TTL)
	if err := redisClient.Expire(ctx, hashKey, ttl).Err(); err != nil {
		return fmt.Errorf(common.ErrorFailedExpSaveRedis, err)
	}

	return nil
}

// FetchFieldFromHashSet Function to fetch data from a Redis Hash using a key and field
func FetchFieldFromHashSet(redisClient *redis.Client, key string, field string) (string, error) {
	// Use HGet to fetch the value of the field in the hash
	result, err := redisClient.HGet(ctx, key, field).Result()
	if err != nil {
		// Use the ErrorFailedRedisFetch constant in the error message
		return "", fmt.Errorf(common.ErrorFailedCacheFetchRedis, err)
	}
	return result, nil
}

// FetchHashSet Function to fetch data from a Redis Hash using a key
func FetchHashSet(redisClient *redis.Client, hashKey string, field string) (map[string]string, error) {
	value, err := redisClient.HGetAll(ctx, hashKey).Result()
	if err != nil {
		return nil, err
	}
	if len(value) == 0 {
		return nil, fmt.Errorf(common.ErrorCacheEmptyRedis, hashKey)
	}
	return value, nil
}
