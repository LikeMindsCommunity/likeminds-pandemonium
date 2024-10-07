package pubsub

import (
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"likeminds-pandemonium/common"
	"likeminds-pandemonium/common/models"
	"likeminds-pandemonium/ws"
	"time"
)

// SentReport struct that holds the sent report data to be saved in Redis
type SentReport struct {
	SenderUUID     string      `json:"sender_uuid"`
	Timestamp      int64       `json:"timestamp"`
	ConversationID interface{} `json:"conversation_id"`
}

// UpdateSentReport Function to update the cache and send payload for Sent Report
func UpdateSentReport(redisClient *redis.Client, wsServerParent *ws.WsServerParent, topic string, response *Response) error {
	var conversationResponse models.ConversationResponse
	if err := json.Unmarshal([]byte(response.RawData), &conversationResponse); err != nil {
		return fmt.Errorf(common.ErrorUnmarshalErrorJson, err)
	}

	senderUUID := conversationResponse.Conversation.Member.UUID
	communityID := conversationResponse.Conversation.CommunityID
	conversationID := conversationResponse.Conversation.ID
	// Generate the cache key for the sent report
	cacheKey := fmt.Sprintf(common.CommunityDeliveryReportPrefix, communityID)

	// Create a SentReport struct instance
	sentReport := SentReport{
		SenderUUID:     senderUUID,
		Timestamp:      time.Now().UnixMilli(),
		ConversationID: conversationID,
	}
	// Marshal the payload into JSON bytes
	sentReportBytes, err := json.Marshal(sentReport)
	if err != nil {
		return fmt.Errorf(common.ErrorMarshalErrorJson, err)
	}
	sentReportResponse := NewResponse(response.DeviceID, common.TopicMessageTypeSentReport, string(sentReportBytes))

	// Use the generic SaveToCacheGeneric function to save the value to Redis
	if err := SaveHashSet(redisClient, cacheKey, fmt.Sprintf(common.UserDeliveryReportFieldPrefix, senderUUID), sentReportResponse, 7*24*time.Hour); err != nil {
		return err
	}

	// Send payload to the user who sent the conversation
	client := wsServerParent.GetConnectionFromWsServer(topic, senderUUID)
	if client != nil {
		// Send payload to the client using WebSocket
		if err := client.SendPayloadToClientConnection(sentReportResponse); err != nil {
			return err
		}
	}

	return nil
}
