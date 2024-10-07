package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"likeminds-pandemonium/common"
	"log"
	"time"
)

var ctx = context.Background()

// InitRedisClient creates a new PublishWithMethod Client
func InitRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     common.GoDotEnvVariable("REDIS_DSN"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

// PublishMessageToRedis is a wrapper method for publishing a message to a Redis topic
func PublishMessageToRedis(redisClient *redis.Client, topic string, message []byte) error {
	// Publish the message to the specified Redis topic
	if err := redisClient.Publish(ctx, topic, message).Err(); err != nil {
		// Log the error and return a formatted error message
		log.Printf(common.ErrorPublishRedis, topic, err)
		return fmt.Errorf(common.ErrorPublishRedis, topic, err)
	}
	// If successful, return nil to indicate success
	return nil
}

// SubscribeToRedisTopic is a wrapper method for subscribing to a Redis topic
func SubscribeToRedisTopic(redisClient *redis.Client, topic string) (*redis.PubSub, error) {
	// Subscribe to the specified Redis topic
	sub := redisClient.Subscribe(ctx, topic)

	// Check for errors in the subscription
	if _, err := sub.Receive(ctx); err != nil {
		// Log the error and return a formatted error message
		log.Printf(common.ErrorSubscribeRedis, topic, err)
		return nil, fmt.Errorf(common.ErrorSubscribeRedis, topic, err)
	}

	// If successful, return the PubSub instance
	return sub, nil
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

// SaveZSet adds a member with a score to a Redis ZSet (sorted set) with a specified TTL (Time to Live)
func SaveZSet(redisClient *redis.Client, zsetKey string, score float64, member interface{}, ttl time.Duration) error {
	// Save the member with its score in the ZSet
	if err := redisClient.ZAdd(ctx, zsetKey, redis.Z{
		Score:  score,
		Member: member,
	}).Err(); err != nil {
		return fmt.Errorf(common.ErrorFailedCacheSaveRedis, err)
	}

	// Set the expiration time (TTL)
	if err := redisClient.Expire(ctx, zsetKey, ttl).Err(); err != nil {
		return fmt.Errorf(common.ErrorFailedExpSaveRedis, err)
	}

	return nil
}

// FetchMembersFromZSet Function to fetch members from a Redis ZSet (sorted set) within a score range
func FetchMembersFromZSet(redisClient *redis.Client, zsetKey string, minScore, maxScore float64) ([]string, error) {
	// Use ZRangeByScore to fetch the members within the score range
	result, err := redisClient.ZRangeByScore(ctx, zsetKey, &redis.ZRangeBy{
		Min: fmt.Sprintf("%f", minScore),
		Max: fmt.Sprintf("%f", maxScore),
	}).Result()
	if err != nil {
		return nil, fmt.Errorf(common.ErrorFailedCacheFetchRedis, err)
	}
	return result, nil
}
