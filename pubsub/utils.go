package pubsub

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// GetRedisClientFromContext Exposed api method to get pubsub client from context
func GetRedisClientFromContext(c *gin.Context) *redis.Client {
	redisClient, exists := c.Get(RedisClient)
	if !exists {
		return nil
	}
	return redisClient.(*redis.Client)
}
