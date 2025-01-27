package middleware

import (
	"likeminds-pandemonium/common"
	"likeminds-pandemonium/ws"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// ApiMiddleware will add the pubsub client  to the context
func ApiMiddleware(client *redis.Client, wsServerParent *ws.WsServerParent) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(common.RedisClient, client)
		c.Set(common.WsServerKey, wsServerParent)
		c.Next()
	}
}
