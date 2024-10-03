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
		case common.TopicTypeChatroom, common.TopicTypeCommunity:
			switch topicMessageType {
			case common.TopicMessageTypeConversation:
				{
					conversationPublished := publishRawDataOnTopic(c, topic, topicMessageType)
					if conversationPublished {
						updateSentReport(c, topic)
					}
				}
			}
		}
	}
}

func publishRawDataOnTopic(c *gin.Context, topic string, topicMessageType string) bool {
	deviceID := c.GetHeader(constant.HeadersDeviceId)
	rawData, _ := c.GetRawData()
	responseBytes, err := json.Marshal(NewResponse(deviceID, topicMessageType, string(rawData)))
	if err != nil {
		return false
	}

	redisClient := GetRedisClientFromContext(c)
	// publish rawData to pubsub channel:<chatroomID>
	if err := redisClient.Publish(ctx, topic, responseBytes).Err(); err != nil {
		api.GeneralAPIError(c, err.Error())
		log.Println(common.ErrorPublishRedis, err)
		return false
	}
	api.GenerateResponse(c, nil)
	return true
}

func updateSentReport(c *gin.Context, topic string) {
	redisClient := GetRedisClientFromContext(c)
	wsServerParent := ws.GetWsServerParentFromContext(c)
	deviceID := c.GetHeader(constant.HeadersDeviceId)
	rawData, _ := c.GetRawData()

	if err := UpdateSentReport(redisClient, wsServerParent, topic, deviceID, rawData); err != nil {
		api.GeneralAPIError(c, err.Error())
		log.Println(err)
		return
	}
	api.GenerateResponse(c, nil)
}
