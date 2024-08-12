package pubsub

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"likeminds-pandemonium/api"
	"likeminds-pandemonium/api/constant"
	"likeminds-pandemonium/common"
	"likeminds-pandemonium/common/models"
	"log"
	"strings"
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
		topic := c.Param(ParamTopic)
		if topic == "" || topic == "null" {
			api.GeneralBadRequestError(c, "topic is required")
			return
		}
		topicSplit := strings.Split(topic, ":")
		if len(topicSplit) < 1 {
			api.GeneralBadRequestError(c, "invalid topic")
			return
		}
		if topicSplit[0] == TopicChatroomType {
			publishOnChatroomTopic(c, topic)
		}
		api.GenerateResponse(c, nil)
	}
}

func publishOnChatroomTopic(c *gin.Context, topic string) {
	deviceID := c.GetHeader(constant.HeadersDeviceId)
	rawData, _ := c.GetRawData()

	var conversationResponse models.ConversationResponse
	if err := json.Unmarshal(rawData, &conversationResponse); err != nil {
		return
	}
	conversationResponse.DeviceID = deviceID
	conversationResponseString, err := json.Marshal(conversationResponse)
	if err != nil {
		return
	}

	redisClient := GetRedisClientFromContext(c)
	// publish rawData to pubsub channel:<chatroomID>
	if err := redisClient.Publish(ctx, topic, conversationResponseString).Err(); err != nil {
		api.GeneralAPIError(c, err.Error())
		log.Println("Publish error:", err)
		return
	}
}
