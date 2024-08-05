package pubsub

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"likeminds-pandemonium/api"
	"likeminds-pandemonium/api/constant"
	"likeminds-pandemonium/common"
	"log"
)

var ctx = context.Background()

// InitRedisClient creates a new PubSub Client
func InitRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     common.GoDotEnvVariable("REDIS_DSN"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func Publish(c *gin.Context) {
	PubSub(c, constant.POSTMethod)
}
func PubSub(c *gin.Context, method int) {
	switch method {
	case constant.POSTMethod:
		topic := c.Param("topic")
		if topic == "" || topic == "null" {
			api.GeneralBadRequestError(c, "topic is required")
			return
		}
		rawData, _ := c.GetRawData()
		redisClient := GetRedisClientFromContext(c)
		// publish rawData to pubsub channel:<chatroomID>
		if err := redisClient.Publish(ctx, topic, rawData).Err(); err != nil {
			api.GeneralAPIError(c, err.Error())
			log.Println("Publish error:", err)
			return
		}
		api.GenerateResponse(c, nil)
	}
}
