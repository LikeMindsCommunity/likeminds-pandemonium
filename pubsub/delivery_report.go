package pubsub

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"likeminds-pandemonium/api"
	"likeminds-pandemonium/api/constant"
	"likeminds-pandemonium/common"
	"likeminds-pandemonium/common/models"
	"likeminds-pandemonium/ws"
	"log"
	"time"
)

// SentReport struct that holds the sent report data to be saved in Redis
type SentReport struct {
	Timestamp      int64       `json:"timestamp"`
	ConversationID interface{} `json:"conversation_id"`
}

// UpdateSentReport Function to update the cache and send payload for Sent Report
func UpdateSentReport(redisClient *redis.Client, wsServerParent *ws.WsServerParent, topic string, response *Response) error {
	var conversationResponse models.ConversationResponse
	if err := json.Unmarshal([]byte(response.RawData), &conversationResponse); err != nil {
		return fmt.Errorf(common.ErrorUnmarshalErrorJson, err)
	}

	communityID := conversationResponse.Conversation.CommunityID
	conversationID := conversationResponse.Conversation.ID
	// Generate the cache key for the sent report
	cacheKey := fmt.Sprintf(common.CommunityDeliveryReportPrefix, communityID)

	// Create a SentReport struct instance
	sentReport := SentReport{
		Timestamp:      time.Now().UnixMilli(),
		ConversationID: conversationID,
	}
	// Marshal the payload into JSON bytes
	sentReportBytes, err := json.Marshal(sentReport)
	if err != nil {
		return fmt.Errorf(common.ErrorMarshalErrorJson, err)
	}

	sentReportResponse := NewResponse(response.DeviceID, common.TopicMessageTypeSentReport, string(sentReportBytes))
	userUUID := conversationResponse.Conversation.Member.UUID
	// Use the generic SaveToCacheGeneric function to save the value to Redis
	if err := SaveHashSet(redisClient, cacheKey, fmt.Sprintf(common.UserDeliveryReportFieldPrefix, userUUID), sentReportResponse, 7*24*time.Hour); err != nil {
		return err
	}

	// Send payload to the user who sent the conversation
	client := wsServerParent.GetConnectionFromWsServer(topic, userUUID)
	if client != nil {
		// Send payload to the client using WebSocket
		if err := client.SendPayloadToClientConnection(sentReportResponse); err != nil {
			return err
		}
	}

	return nil
}

// UpdateDeliveredReport Function to update the cache and send payload for Delivered Report
func UpdateDeliveredReport(redisClient *redis.Client, wsServerParent *ws.WsServerParent, topic string, response *Response, receiverUUID string) error {
	var conversationResponse models.ConversationResponse
	if err := json.Unmarshal([]byte(response.RawData), &conversationResponse); err != nil {
		return fmt.Errorf(common.ErrorUnmarshalErrorJson, err)
	}

	// Create a DeliveryReport struct instance
	deliveryReport := SentReport{
		Timestamp: time.Now().UnixMilli(),
	}
	// Marshal the payload into JSON bytes
	reportBytes, err := json.Marshal(deliveryReport)
	if err != nil {
		return fmt.Errorf(common.ErrorMarshalErrorJson, err)
	}
	reportResponse := NewResponse(response.DeviceID, common.TopicMessageTypeDeliveredReport, string(reportBytes))
	// Use the generic SaveToCacheGeneric function to save the value to Redis
	communityID := conversationResponse.Conversation.CommunityID
	cacheKey := fmt.Sprintf(common.CommunityDeliveryReportPrefix, communityID)
	if err := SaveHashSet(redisClient, cacheKey, fmt.Sprintf(common.UserDeliveryReportFieldPrefix, receiverUUID), reportResponse, 7*24*time.Hour); err != nil {
		return err
	}

	// Send payload to the user who sent the conversation and is still active
	userUUID := conversationResponse.Conversation.Member.UUID
	client := wsServerParent.GetConnectionFromWsServer(topic, userUUID)
	if client != nil {
		// Send payload to the client using WebSocket
		if err := client.SendPayloadToClientConnection(reportResponse); err != nil {
			return err
		}
	}

	return nil
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
