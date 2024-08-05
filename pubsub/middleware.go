package pubsub

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// ApiMiddleware will add the pubsub client  to the context
func ApiMiddleware(client *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(RedisClient, client)
		c.Next()
	}
}
