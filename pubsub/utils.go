package pubsub

import (
	"errors"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/common"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
)

type PSResponse struct {
	DeviceID         string `json:"device_id"`
	TopicMessageType string `json:"topic_message_type"`
	RawData          string `json:"raw_data"`
}

// NewResponse to create PSResponse from deviceID (device ID of user), topicMessageType (message type publish to connection), rawData (raw data publish to the connection)
func NewResponse(deviceID string, topicMessageType string, rawData string) *PSResponse {
	return &PSResponse{DeviceID: deviceID, TopicMessageType: topicMessageType, RawData: rawData}
}

// GetRedisClientFromContext Exposed api method to get pubsub client from context
func GetRedisClientFromContext(c *gin.Context) *redis.Client {
	redisClient, exists := c.Get(common.RedisClient)
	if !exists {
		return nil
	}
	return redisClient.(*redis.Client)
}

// GetTopicSplit will decode topic and return split
func GetTopicSplit(topic string) ([]string, error) {
	if topic == "" || topic == "null" {
		return nil, errors.New(common.ErrorTopicMissing)
	}
	topicSplit := strings.Split(topic, ":")
	if len(topicSplit) <= 1 {
		return nil, errors.New(common.ErrorTopicInvalid)
	}
	return topicSplit, nil
}
