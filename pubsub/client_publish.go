package pubsub

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"likeminds-pandemonium/api"
	"likeminds-pandemonium/api/constant"
	"likeminds-pandemonium/common"
	"likeminds-pandemonium/ws"
	"log"
	"strconv"
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
				go updateSentDR(c, topic)

			case common.TopicMessageTypeDeliveredDR:
				updateDeliveredDROnPublish(c, topic)
			}
		case common.TopicTypeCommunity:
			switch topicMessageType {
			case common.TopicMessageTypeConversation:
				publishRawDataOnTopic(c, topic, topicMessageType)
			}
		}
	}
}

func publishRawDataOnTopic(c *gin.Context, topic string, topicMessageType string) {
	deviceID := c.GetHeader(constant.HeadersDeviceID)
	rawData, _ := c.GetRawData()
	//To reuse rawData
	c.Set(common.RawData, rawData)
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
	//Reusing it from gin context
	rawData, _ := c.Get(common.RawData)

	redisClient := GetRedisClientFromContext(c)
	wsServerParent := ws.GetWsServerParentFromContext(c)

	if err := UpdateSentDR(redisClient, wsServerParent, topic, deviceID, rawData.([]byte)); err != nil {
		log.Println(err)
	}
}

type PublishDeliveredDR struct {
	MinTimestamp string `json:"min_timestamp"`
	MaxTimestamp string `json:"max_timestamp"`
}

func updateDeliveredDROnPublish(c *gin.Context, topic string) {
	// Extract chatroom_id from the topic by splitting.
	topicSplit, err := GetTopicSplit(topic)
	if err != nil || len(topicSplit) < 2 {
		api.GeneralBadRequestError(c, common.ErrorTopicInvalid)
		return
	}
	chatroomID := topicSplit[1]

	deliveredUUID := c.GetHeader(constant.HeadersMemberID)
	if deliveredUUID == "" || deliveredUUID == "null" {
		api.GeneralUnauthorizedError(c, common.ErrorUserUUIDMissing)
		return
	}

	deliveredDeviceID := c.GetHeader(constant.HeadersDeviceID)

	// Get the min_timestamp and max_timestamp from the body.
	var deliveredDR PublishDeliveredDR
	if err := c.ShouldBindJSON(&deliveredDR); err != nil {
		api.GeneralBadRequestError(c, fmt.Sprintf(common.ErrorInvalidJSONFormat, err))
		return
	}

	// Parse timestamps to float64 for Redis.
	minTimestampString, err := strconv.ParseFloat(deliveredDR.MinTimestamp, 64)
	if err != nil {
		api.GeneralBadRequestError(c, common.ErrorInvalidTimestamp)
		return
	}
	maxTimestampString, err := strconv.ParseFloat(deliveredDR.MaxTimestamp, 64)
	if err != nil {
		api.GeneralBadRequestError(c, common.ErrorInvalidTimestamp)
		return
	}

	// Construct the Redis key for the chatroom delivery report.
	redisKey := fmt.Sprintf(common.DRChatroomPrefix, chatroomID)

	// Get the Redis client from the context.
	redisClient := GetRedisClientFromContext(c)
	wsServerParent := ws.GetWsServerParentFromContext(c)

	// Step 1: Fetch all member values between min and max timestamps from the Redis key.
	drConversations, err := FetchMembersFromZSet(redisClient, redisKey, minTimestampString, maxTimestampString)
	if err != nil {
		api.GeneralAPIError(c, err.Error())
		return
	}

	// Step 2: Iterate over each member value, which corresponds to conversation IDs.
	for _, drConversation := range drConversations {
		// Construct the key for the conversation delivery report.
		conversationKey := fmt.Sprintf("%s", drConversation.Member)

		// Fetch the conversation delivery report from Redis.
		conversationData, err := FetchFieldFromHashSet(redisClient, conversationKey, common.DRConversationMetaPrefix)
		if err != nil {
			log.Println(err)
			continue
		}

		// Unmarshal the fetched data into a map.
		var conversationMap map[string]interface{}
		if err := json.Unmarshal([]byte(conversationData), &conversationMap); err != nil {
			log.Println(fmt.Sprintf(common.ErrorUnmarshalErrorJson, err))
			continue
		}

		// Extract the sender UUID from the fetched conversation data.
		senderUUID, _ := conversationMap["sender_uuid"].(string)

		// Update the delivered report using the common function.
		if err := UpdateDeliveredDR(redisClient, wsServerParent, topic, conversationKey, deliveredDeviceID, senderUUID, deliveredUUID); err != nil {
			log.Println(err)
		}
	}
}
