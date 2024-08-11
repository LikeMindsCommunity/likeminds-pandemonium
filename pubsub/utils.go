package pubsub

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"likeminds-pandemonium/common"
	"strings"
)

var ctx = context.Background()

type Response struct {
	DeviceID         string `json:"device_id"`
	TopicMessageType string `json:"topic_message_type"`
	RawData          []byte `json:"raw_data"`
}

func NewResponse(deviceID string, topicMessageType string, rawData []byte) *Response {
	return &Response{DeviceID: deviceID, TopicMessageType: topicMessageType, RawData: rawData}
}

// InitRedisClient creates a new PublishWithMethod Client
func InitRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     common.GoDotEnvVariable("REDIS_DSN"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

// GetRedisClientFromContext Exposed api method to get pubsub client from context
func GetRedisClientFromContext(c *gin.Context) *redis.Client {
	redisClient, exists := c.Get(RedisClient)
	if !exists {
		return nil
	}
	return redisClient.(*redis.Client)
}

// GetTopicSplit will decode topic and return split
func GetTopicSplit(topic string) ([]string, error) {
	if topic == "" || topic == "null" {
		return nil, errors.New("topic is required")
	}
	topicSplit := strings.Split(topic, ":")
	if len(topicSplit) < 1 {
		return nil, errors.New("invalid topic")
	}
	return topicSplit, nil
}
