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
	"strings"
	"time"
)

// DeliveryReport struct that holds the delivery report data to be saved in Redis
type DeliveryReport struct {
	Timestamp      int64       `json:"timestamp"`
	ConversationID interface{} `json:"conversation_id,omitempty"`
	UserUUID       interface{} `json:"user_uuid"`
}

func UpdateSentDR(redisClient *redis.Client, wsServerParent *ws.WsServerParent, topic, deviceID string, rawData []byte) error {
	var conversationResponse models.ConversationResponse
	if err := json.Unmarshal(rawData, &conversationResponse); err != nil {
		return fmt.Errorf(common.ErrorUnmarshalErrorJson, err)
	}

	conversationID := conversationResponse.Conversation.ID
	chatroomID := conversationResponse.Conversation.ChatroomID
	userUUID := conversationResponse.Conversation.Member.UUID
	participantsCount := len(conversationResponse.Participants)

	// Create the cache keys
	chatroomKey := fmt.Sprintf(common.DRChatroomPrefix, chatroomID)
	conversationKey := fmt.Sprintf(common.DRConversationPrefix, conversationID)

	// Update the cache for chatroom and conversation
	err := SaveZSet(redisClient, chatroomKey, float64(conversationResponse.Conversation.CreatedAt), conversationKey, 7*24*time.Hour)
	if err != nil {
		return err
	}

	sentDRValue := map[string]interface{}{
		"delivery_count": participantsCount,
		"sender_uuid":    userUUID,
	}
	err = SaveHashSet(redisClient, conversationKey, common.DRConversationMetaPrefix, sentDRValue, 7*24*time.Hour)
	if err != nil {
		return err
	}

	// Send the payload to the conversation creator
	client := wsServerParent.GetConnectionFromWsServer(topic, userUUID)
	if client != nil {
		payload := NewResponse(deviceID, common.TopicMessageTypeSentDR, string(rawData))
		if err := client.SendPayloadToClientConnection(payload); err != nil {
			return err
		}
	}

	return nil
}

// UpdateDeliveredDRWithConversationID updates the delivered report in Redis and sends a payload to the conversation creator.
func UpdateDeliveredDRWithConversationID(redisClient *redis.Client, wsServerParent *ws.WsServerParent, chatroomID, conversationID interface{}, deliveredDeviceID, senderUUID, deliveredUUID string) error {
	// Fetch and update the dr_conversation_<conversation_id>
	conversationKey := fmt.Sprintf(common.DRConversationPrefix, conversationID)
	return UpdateDeliveredDR(redisClient, wsServerParent, chatroomID, conversationKey, deliveredDeviceID, senderUUID, deliveredUUID)
}

// UpdateDeliveredDR updates the delivered report in Redis and sends a payload to the conversation creator.
func UpdateDeliveredDR(redisClient *redis.Client, wsServerParent *ws.WsServerParent, chatroomID interface{}, conversationKey string, deliveredDeviceID, senderUUID, deliveredUUID string) error {
	// If the message sender is same as message delivered user, no need to update
	if senderUUID == deliveredUUID {
		return nil
	}

	// Construct the field for the delivered report using the new key format.
	deliveredUUIDField := fmt.Sprintf(common.DRUserPrefix, deliveredUUID)

	// Check if the delivered report for this user already exists.
	existingDeliveredReport, err := FetchFieldFromHashSet(redisClient, conversationKey, deliveredUUIDField)

	// If the delivered report already exists, no need to update.
	if existingDeliveredReport != "" {
		return nil
	}

	// Set the current timestamp as the delivered timestamp.
	currentTimestamp := time.Now().UnixMilli()
	// Update the Redis key with the new delivered report.
	err = SaveHashSet(redisClient, conversationKey, deliveredUUIDField, currentTimestamp, 7*24*time.Hour)
	if err != nil {
		return err
	}

	topicChatroom := fmt.Sprintf(common.TopicTypeChatroomDynamic, chatroomID)
	// Send the updated delivered report to the conversation creator's connection.
	client := wsServerParent.GetConnectionFromWsServer(topicChatroom, senderUUID)
	if client != nil {
		// Create the payload for the delivered report.
		deliveredReport := map[string]interface{}{
			deliveredUUIDField: currentTimestamp,
		}

		// Marshal the updated delivered report for the response.
		deliveredReportBytes, err := json.Marshal(deliveredReport)
		if err != nil {
			return err
		}

		// Create the response payload for the delivered report.
		deliveredReportResponse := NewResponse(deliveredDeviceID, common.TopicMessageTypeDeliveredDR, string(deliveredReportBytes))

		// Send the payload via WebSocket to the conversation creator.
		if err := client.SendPayloadToClientConnection(deliveredReportResponse); err != nil {
			return err
		}
	}

	return nil
}

// DeliveryReportRequest defines the structure of the request body for the /delivery_report API.
type DeliveryReportRequest struct {
	ChatroomID      string   `json:"chatroom_id" binding:"required"`
	ConversationIDs []string `json:"conversation_ids" binding:"required"`
}

// DeliveryReportResponse defines the structure of the response for the /delivery_report API.
type DeliveryReportResponse struct {
	DeliveryReport map[string]map[string]interface{} `json:"delivery_report"`
}

// DeliveryReportHandler handles the request for the /delivery_report API.
func DeliveryReportHandler(c *gin.Context) {
	// Parse x-member-id from headers.
	memberID := c.GetHeader(constant.HeadersMemberID)
	if memberID == "" || memberID == "null" {
		api.GeneralUnauthorizedError(c, common.ErrorUserUUIDMissing)
		return
	}

	// Parse the request body to extract chatroom_id and conversation_ids.
	var request DeliveryReportRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		api.GeneralBadRequestError(c, common.ErrorInvalidJSONFormat)
		return
	}

	// Validate that chatroom_id and conversation_ids are present.
	chatroomID := request.ChatroomID
	if chatroomID == "" {
		api.GeneralBadRequestError(c, common.ErrorChatroomIDMissing)
		return
	}
	conversationIDs := request.ConversationIDs
	if len(conversationIDs) == 0 {
		api.GeneralBadRequestError(c, common.ErrorConversationIDsMissing)
		return
	}

	// Get Redis client from context.
	redisClient := GetRedisClientFromContext(c)

	// Fetch the conversation delivery report data for all provided conversation IDs.
	conversationKeys := make([]string, len(conversationIDs))
	for i, conversationID := range request.ConversationIDs {
		conversationKeys[i] = fmt.Sprintf(common.DRConversationPrefix, conversationID)
	}

	// Perform a pipeline GET operation for all conversation keys to get the delivery reports.
	pipe := redisClient.Pipeline()
	cmds := make([]*redis.MapStringStringCmd, len(conversationKeys))
	for i, conversationKey := range conversationKeys {
		cmds[i] = pipe.HGetAll(c, conversationKey)
	}
	_, err := pipe.Exec(c)
	if err != nil {
		api.GeneralAPIError(c, err.Error())
		return
	}

	// Construct the response.
	deliveryReport := make(map[string]map[string]interface{})
	for i, cmd := range cmds {
		conversationID := conversationIDs[i]

		// Get the result of the HGETALL command.
		data, err := cmd.Result()
		if err != nil {
			log.Printf("Error fetching data for conversation %s: %v", conversationID, err)
			continue
		}

		// Extract the "dr_conversation_meta" field.
		metaData, ok := data[common.DRConversationMetaPrefix]
		if !ok || metaData == "" {
			log.Printf("Missing or empty dr_conversation_meta for conversation %s", conversationID)
			continue
		}

		// Unmarshal the metadata field.
		var metaMap map[string]interface{}
		if err := json.Unmarshal([]byte(metaData), &metaMap); err != nil {
			log.Printf("Error unmarshalling conversation meta for %s: %v", conversationID, err)
			continue
		}

		// Add the conversation data to the delivery report map directly.
		deliveryReport[conversationID] = map[string]interface{}{
			"delivery_count": metaMap["delivery_count"],
			"sender_uuid":    metaMap["sender_uuid"],
			"delivered_dr":   extractDeliveredDRFields(data),
		}
	}

	// Create and send the response.
	response := DeliveryReportResponse{
		DeliveryReport: deliveryReport,
	}
	api.GenerateResponse(c, response)
}

// extractDeliveredDRFields extracts the delivered reports from the Redis data map.
func extractDeliveredDRFields(data map[string]string) map[string]interface{} {
	deliveredDR := make(map[string]interface{})

	// Iterate through all fields in the Redis hash and find delivered report fields.
	for key, value := range data {
		if strings.HasPrefix(key, common.DRUser) {
			deliveredDR[key] = value
		}
	}
	return deliveredDR
}
