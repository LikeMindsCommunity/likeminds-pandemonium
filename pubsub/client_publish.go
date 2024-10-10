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
				publishRawDataOnTopic(c, topic, topicMessageType)
				updateSentDR(c, topic)

			case common.TopicMessageTypeDeliveredDR:
				updateDeliveredDROnPublish(c, topic)
			}
		case common.TopicTypeCommunity:
			switch topicMessageType {
			case common.TopicMessageTypeConversation:
				publishRawDataOnTopic(c, topic, topicMessageType)
			case common.TopicMessageTypeDeliveredDR:
				updateDeliveredDROnPublish(c, topic)
			}
		}
	}
}

func publishRawDataOnTopic(c *gin.Context, topic string, topicMessageType string) {
	deviceID := c.GetHeader(constant.HeadersDeviceID)
	rawData, _ := c.GetRawData()
	psResponse := NewResponse(deviceID, topicMessageType, string(rawData))
	responseBytes, err := json.Marshal(psResponse)
	if err != nil {
		return
	}

	redisClient := GetRedisClientFromContext(c)
	// publish rawData to pubsub channel:<chatroomID>
	if err := PublishMessageToRedis(redisClient, topic, responseBytes); err != nil {
		api.GeneralAPIError(c, err.Error())
		return
	}
	api.GenerateResponse(c, nil)
}

func updateSentDR(c *gin.Context, topic string) {
	deviceID := c.GetHeader(constant.HeadersDeviceID)
	rawData, _ := c.GetRawData()

	redisClient := GetRedisClientFromContext(c)
	wsServerParent := ws.GetWsServerParentFromContext(c)

	if err := UpdateSentDR(redisClient, wsServerParent, topic, deviceID, rawData); err != nil {
		log.Println(err)
	}
}

type PublishDeliveredDR struct {
	SenderUUIDs []string `json:"sender_uuids"`
}

func updateDeliveredDROnPublish(c *gin.Context, topic string) {
	// Get the community_id from the query parameters
	chatroomID := c.Param("chatroom_id")
	if chatroomID == "" || chatroomID == "null" {
		api.GeneralBadRequestError(c, common.ErrorChatroomIDMissing)
		return
	}
	receiverUUID := c.GetHeader(constant.HeadersMemberID)
	if receiverUUID == "" || receiverUUID == "null" {
		api.GeneralUnauthorizedError(c, common.ErrorUserUUIDMissing)
		return
	}
	deviceID := c.GetHeader(constant.HeadersDeviceID)
	// Get the list of receiver_uuids from the request body
	var deliveredDR PublishDeliveredDR
	rawData, _ := c.GetRawData()
	if err := json.Unmarshal(rawData, &deliveredDR); err != nil {
		log.Printf(common.ErrorUnmarshalErrorJson, err)
	}

	redisClient := GetRedisClientFromContext(c)
	wsServerParent := ws.GetWsServerParentFromContext(c)

	senderUUIDs := deliveredDR.SenderUUIDs
	if senderUUIDs == nil || len(senderUUIDs) == 0 {
		if err := UpdateDeliveredDR(redisClient, wsServerParent, topic, chatroomID, deviceID, "", receiverUUID); err != nil {
			log.Println(err)
		}
	} else {
		for _, senderUUID := range senderUUIDs {
			if err := UpdateDeliveredDR(redisClient, wsServerParent, topic, chatroomID, deviceID, senderUUID, receiverUUID); err != nil {
				log.Println(err)
			}
		}
	}
}
