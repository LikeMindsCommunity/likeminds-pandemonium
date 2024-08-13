package pubsub

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"likeminds-pandemonium/api"
	"likeminds-pandemonium/api/constant"
	"log"
)

func Publish(c *gin.Context) {
	PublishWithMethod(c, constant.POSTMethod)
}

func PublishWithMethod(c *gin.Context, method int) {
	switch method {
	case constant.POSTMethod:
		topic := c.Param(ParamTopic)
		topicMessageType := c.Query(ParamTopicMessageType)
		if topicMessageType == "" || topicMessageType == "null" {
			api.GeneralBadRequestError(c, "topic_message_type is required")
			return
		}
		topicSplit, err := GetTopicSplit(topic)
		if err != nil {
			api.GeneralBadRequestError(c, err.Error())
			return
		}

		switch topicSplit[0] {
		case TopicTypeChatroom, TopicTypeCommunity:
			switch topicMessageType {
			case TopicMessageTypeConversation:
				publishRawDataOnTopic(c, topic, topicMessageType)

			}
		}
		api.GenerateResponse(c, nil)
	}
}

func publishRawDataOnTopic(c *gin.Context, topic string, topicMessageType string) {
	deviceID := c.GetHeader(constant.HeadersDeviceId)
	rawData, _ := c.GetRawData()
	responseBytes, err := json.Marshal(NewResponse(deviceID, topicMessageType, string(rawData)))
	if err != nil {
		return
	}

	redisClient := GetRedisClientFromContext(c)
	// publish rawData to pubsub channel:<chatroomID>
	if err := redisClient.Publish(ctx, topic, responseBytes).Err(); err != nil {
		api.GeneralAPIError(c, err.Error())
		log.Println("Publish error:", err)
		return
	}
}
