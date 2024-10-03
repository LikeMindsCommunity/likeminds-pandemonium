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
	SenderID       string `json:"sender_id"`
	ConversationID string `json:"conversation_id"`
	Timestamp      int64  `json:"timestamp"`
	TotalCount     string `json:"total_count"`
}

// UpdateSentReport Function to update the cache and send payload for Sent Report
func UpdateSentReport(redisClient *redis.Client, wsServerParent *ws.WsServerParent, topic, deviceID string, rawData []byte) error {
	var response Response
	if err := json.Unmarshal(rawData, &response); err != nil {
		return fmt.Errorf(common.ErrorUnmarshalErrorJson, err)
	}
	var conversationResponse models.ConversationResponse
	if err := json.Unmarshal([]byte(response.RawData), &conversationResponse); err != nil {
		return fmt.Errorf(common.ErrorUnmarshalErrorJson, err)
	}

	senderUUID := conversationResponse.Conversation.Member.UUID
	conversationID := conversationResponse.Conversation.ID
	chatroomID := conversationResponse.Conversation.ChatroomID
	// Generate the cache key for the sent report
	cacheKey := fmt.Sprintf(common.ChatroomDeliveryReportPrefix, chatroomID)

	// Create a SentReport struct instance
	sentReport := SentReport{
		SenderID:       senderUUID,
		ConversationID: conversationID,
		Timestamp:      time.Now().UnixMilli(),
		TotalCount:     "", // This would be calculated dynamically, set to empty for now
	}
	// Marshal the payload into JSON bytes
	sentReportBytes, err := json.Marshal(sentReport)
	if err != nil {
		return fmt.Errorf(common.ErrorMarshalErrorJson, err)
	}
	sentReportResponse := NewResponse(deviceID, common.TopicMessageTypeSentReport, string(sentReportBytes))

	// Use the generic SaveToCacheGeneric function to save the value to Redis
	if err := SaveHashSet(redisClient, cacheKey, conversationID, sentReportResponse, 7*24*time.Hour); err != nil {
		return err
	}

	// Send payload to the user who sent the conversation
	client := wsServerParent.GetConnectionFromWsServer(topic, senderUUID)
	if client != nil {
		// Send payload to the client using WebSocket
		if err := client.SendPayloadToClientConnection(sentReport); err != nil {
			return err
		}
	}

	return nil
}
