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
	deviceID := c.GetHeader(constant.HeadersDeviceID)
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

func DeliverReports(c *gin.Context) {
	// Get the community_id from the query parameters
	communityID := c.Query("community_id")
	if communityID == "" || communityID == "null" {
		api.GeneralBadRequestError(c, "community_id is required")
		return
	}

	// Get x-member-id from headers
	memberID := c.GetHeader(constant.HeadersMemberID)
	if memberID == "" || memberID == "null" {
		api.GeneralBadRequestError(c, "x-member-id is required")
		return
	}

	// Create Redis key for delivery report
	communityDRKey := fmt.Sprintf(common.CommunityDeliveryReportPrefix, communityID)
	userDRField := fmt.Sprintf(common.UserDeliveryReportFieldPrefix, memberID)

	// Get Redis client from context
	redisClient := GetRedisClientFromContext(c)

	// Fetch the delivery report from Redis using HGET
	drCacheValue, err := redisClient.HGet(c, communityDRKey, userDRField).Result()
	if err != nil {
		api.GeneralStatusNotFoundError(c, "Failed to fetch delivery report: "+err.Error())
		log.Println("Error fetching delivery report:", err)
		return
	}

	if drCacheValue == "" {
		api.GeneralStatusNotFoundError(c, "No delivery report found for this user")
		return
	}

	// Unmarshal the drCacheValue into the drResponse structure
	var drResponse Response
	err = json.Unmarshal([]byte(drCacheValue), &drResponse)
	if err != nil {
		api.GeneralAPIError(c, "Error parsing delivery report: "+err.Error())
		log.Println("Error unmarshalling delivery report:", err)
		return
	}

	api.GenerateResponse(c, drResponse)
}
