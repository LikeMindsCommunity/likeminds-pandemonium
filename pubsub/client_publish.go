package pubsub

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"likeminds-pandemonium/api"
	"likeminds-pandemonium/api/constant"
	"likeminds-pandemonium/common"
	"likeminds-pandemonium/ws"
	"log"
)

func Publish(c *gin.Context) {
	PublishWithMethod(c, constant.POSTMethod)
}

func PublishWithMethod(c *gin.Context, method int) {
	switch method {
	case constant.POSTMethod:
		topic := c.Param(common.ParamTopic)
		topicMessageType := c.Query(common.ParamTopicMessageType)
		if topicMessageType == "" || topicMessageType == "null" {
			api.GeneralBadRequestError(c, common.ErrorTopicMessageTypeMissing)
			return
		}
		topicSplit, err := GetTopicSplit(topic)
		if err != nil {
			api.GeneralBadRequestError(c, err.Error())
			return
		}

		switch topicSplit[0] {
		case common.TopicTypeChatroom:
			switch topicMessageType {
			case common.TopicMessageTypeConversation:
				response := publishRawDataOnTopic(c, topic, topicMessageType)
				if response != nil {
					updateSentReport(c, topic, response)
				}
			}
		case common.TopicTypeCommunity:
			switch topicMessageType {
			case common.TopicMessageTypeConversation:
				_ = publishRawDataOnTopic(c, topic, topicMessageType)
			}
		}
	}
}

func publishRawDataOnTopic(c *gin.Context, topic string, topicMessageType string) *Response {
	deviceID := c.GetHeader(constant.HeadersDeviceId)
	rawData, _ := c.GetRawData()
	response := NewResponse(deviceID, topicMessageType, string(rawData))
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return nil
	}

	redisClient := GetRedisClientFromContext(c)
	// publish rawData to pubsub channel:<chatroomID>
	if err := redisClient.Publish(ctx, topic, responseBytes).Err(); err != nil {
		api.GeneralAPIError(c, err.Error())
		log.Println(common.ErrorPublishRedis, err)
		return nil
	}
	api.GenerateResponse(c, nil)
	return response
}

func updateSentReport(c *gin.Context, topic string, response *Response) {
	redisClient := GetRedisClientFromContext(c)
	wsServerParent := ws.GetWsServerParentFromContext(c)

	if err := UpdateSentReport(redisClient, wsServerParent, topic, response); err != nil {
		api.GeneralAPIError(c, err.Error())
		log.Println(err)
		return
	}
	api.GenerateResponse(c, nil)
}
