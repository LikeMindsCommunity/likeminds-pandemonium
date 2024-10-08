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
			case common.TopicMessageTypeDeliveredReport:
				// Handle the delivered_report case
				updateDeliveredReportOnPublish(c, topic)
			}
		case common.TopicTypeCommunity:
			switch topicMessageType {
			case common.TopicMessageTypeConversation:
				_ = publishRawDataOnTopic(c, topic, topicMessageType)
			case common.TopicMessageTypeDeliveredReport:
				// Handle the delivered_report case
				updateDeliveredReportOnPublish(c, topic)
			}
		}
	}
}

func publishRawDataOnTopic(c *gin.Context, topic string, topicMessageType string) *Response {
	deviceID := c.GetHeader(constant.HeadersDeviceID)
	rawData, _ := c.GetRawData()
	response := NewResponse(deviceID, topicMessageType, string(rawData))
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return nil
	}

	redisClient := GetRedisClientFromContext(c)
	// publish rawData to pubsub channel:<chatroomID>
	if err := PublishMessageToRedis(redisClient, topic, responseBytes); err != nil {
		api.GeneralAPIError(c, err.Error())
		return nil
	}
	api.GenerateResponse(c, nil)
	return response
}

func updateSentReport(c *gin.Context, topic string, response *Response) {
	redisClient := GetRedisClientFromContext(c)
	wsServerParent := ws.GetWsServerParentFromContext(c)

	if err := UpdateSentReport(redisClient, wsServerParent, topic, response); err != nil {
		log.Println(err)
	}
}

type PublishDeliveredRequest struct {
	ReceiverUUID []string `json:"receiver_uuids"`
}

func updateDeliveredReportOnPublish(c *gin.Context, topic string) {
	// Get the community_id from the query parameters
	chatroomID := c.Param("chatroom_id")
	if chatroomID == "" || chatroomID == "null" {
		api.GeneralBadRequestError(c, common.ErrorChatroomIDMissing)
		return
	}
	senderUUID := c.GetHeader(constant.HeadersMemberID)
	deviceID := c.GetHeader(constant.HeadersDeviceID)
	if senderUUID == "" || senderUUID == "null" {
		api.GeneralUnauthorizedError(c, common.ErrorUserUUIDMissing)
		return
	}

	// Get the list of receiver_uuids from the request body
	var deliveredRequest PublishDeliveredRequest
	rawData, _ := c.GetRawData()
	if err := json.Unmarshal(rawData, &deliveredRequest); err != nil {
		log.Printf(common.ErrorUnmarshalErrorJson, err)
	}
	redisClient := GetRedisClientFromContext(c)
	wsServerParent := ws.GetWsServerParentFromContext(c)
	for _, receiverUUID := range deliveredRequest.ReceiverUUID {
		if err := UpdateDeliveredReport(redisClient, wsServerParent, topic, chatroomID, deviceID, senderUUID, receiverUUID); err != nil {
			log.Println(err)
		}
	}
}
