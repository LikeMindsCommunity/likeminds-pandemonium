package pubsub

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"likeminds-pandemonium/api"
	"likeminds-pandemonium/api/constant"
	requestresponse "likeminds-pandemonium/api/request_response"
	"likeminds-pandemonium/common"
	"likeminds-pandemonium/ws"
	"log"
	"strings"
	"time"
)

func UpdateSentDR(redisClient *redis.Client, wsServerParent *ws.WsServerParent, deviceID string, rawData []byte) error {
	var createMessagePSResponse requestresponse.PSResponse
	if err := json.Unmarshal(rawData, &createMessagePSResponse); err != nil {
		return fmt.Errorf(common.ErrorUnmarshalErrorJson, err)
	}
	var createMessageResponse requestresponse.CreateMessageResponse
	if err := json.Unmarshal([]byte(createMessagePSResponse.RawData), &createMessageResponse); err != nil {
		return fmt.Errorf(common.ErrorUnmarshalErrorJson, err)
	}

	conversationID := createMessageResponse.Data.Message.ID
	chatroomID := createMessageResponse.Data.Message.CardID
	userUUID := *createMessageResponse.Data.User.UserUniqueID
	participantsCount := createMessageResponse.TotalParticipantsCount
	communityID := createMessageResponse.Data.Message.CommunityID

	// Create the cache keys
	chatroomKey := fmt.Sprintf(common.DRChatroomPrefix, chatroomID)
	conversationKey := fmt.Sprintf(common.DRConversationPrefix, conversationID)

	// Update the cache for chatroom and conversation
	err := SaveZSet(redisClient, chatroomKey, float64(createMessageResponse.Data.Message.CreatedAt), conversationKey, common.DeliveryReportTTL)
	if err != nil {
		return err
	}

	sentDRValue := map[string]interface{}{
		common.DeliveryCount: participantsCount,
		common.SenderUUID:    userUUID,
	}
	err = SaveHashSet(redisClient, conversationKey, common.DRConversationMetaPrefix, sentDRValue, common.DeliveryReportTTL)
	if err != nil {
		return err
	}

	topicChatroom := fmt.Sprintf(common.TopicTypeChatroomDynamic, chatroomID)
	topicCommunity := fmt.Sprintf(common.TopicTypeCommunityDynamic, communityID)
	// Send the updated delivered report to the conversation creator's connection.
	clientConnectedToChatroom := wsServerParent.GetConnectionFromWsServer(topicChatroom, userUUID)
	clientConnectedToCommunity := wsServerParent.GetConnectionFromWsServer(topicCommunity, userUUID)
	finalClient := clientConnectedToChatroom
	if finalClient == nil {
		finalClient = clientConnectedToCommunity
	}
	if finalClient != nil {
		payload := NewResponse(deviceID, common.TopicMessageTypeSentDR, string(rawData))
		if err := finalClient.SendPayloadToClientConnection(payload); err != nil {
			return err
		}
	}

	return nil
}

// UpdateDeliveredDRWithConversationID updates the delivered report in Redis and sends a payload to the conversation creator.
func UpdateDeliveredDRWithConversationID(redisClient *redis.Client, wsServerParent *ws.WsServerParent, chatroomID, conversationID interface{}, deliveredDeviceID, senderUUID, deliveredUUID string, communityID interface{}) error {
	// Fetch and update the dr_conversation_<conversation_id>
	conversationKey := fmt.Sprintf(common.DRConversationPrefix, conversationID)
	return UpdateDeliveredDR(redisClient, wsServerParent, chatroomID, conversationKey, deliveredDeviceID, senderUUID, deliveredUUID, communityID)
}

// UpdateDeliveredDR updates the delivered report in Redis and sends a payload to the conversation creator.
func UpdateDeliveredDR(redisClient *redis.Client, wsServerParent *ws.WsServerParent, chatroomID interface{}, conversationKey string, deliveredDeviceID, senderUUID, deliveredUUID string, communityID interface{}) error {
	// If the message sender is the same as the message delivered user, no need to update
	if senderUUID == deliveredUUID {
		return nil
	}

	// Construct the field for the delivered report using the new key format.
	deliveredUUIDField := fmt.Sprintf(common.DRUserDeliveredPrefix, deliveredUUID)

	// Check if the delivered eport for this user already exists.
	existingDeliveredReport, err := FetchFieldFromHashSet(redisClient, conversationKey, deliveredUUIDField)

	// If the delivered report already exists, no need to update.
	if existingDeliveredReport != "" {
		return nil
	}

	// Set the current timestamp as the delivered timestamp.
	currentTimestamp := time.Now().UnixMilli()
	// Update the Redis key with the new delivered report.
	err = SaveHashSet(redisClient, conversationKey, deliveredUUIDField, currentTimestamp, common.DeliveryReportTTL)
	if err != nil {
		return err
	}

	topicChatroom := fmt.Sprintf(common.TopicTypeChatroomDynamic, chatroomID)
	topicCommunity := fmt.Sprintf(common.TopicTypeCommunityDynamic, communityID)
	// Send the updated delivered report to the conversation creator's connection.
	clientConnectedToChatroom := wsServerParent.GetConnectionFromWsServer(topicChatroom, senderUUID)
	clientConnectedToCommunity := wsServerParent.GetConnectionFromWsServer(topicCommunity, senderUUID)
	finalClient := clientConnectedToChatroom
	if finalClient == nil {
		finalClient = clientConnectedToCommunity
	}
	if finalClient != nil {
		// Fetch the existing conversation metadata to include in the delivered report.
		conversationMeta, err := FetchFieldFromHashSet(redisClient, conversationKey, common.DRConversationMetaPrefix)
		if err != nil || conversationMeta == "" {
			return fmt.Errorf("failed to fetch conversation meta: %v", err)
		}

		// Unmarshal the conversation metadata into a map.
		var deliveredReport map[string]interface{}
		if err := json.Unmarshal([]byte(conversationMeta), &deliveredReport); err != nil {
			return fmt.Errorf("failed to unmarshal conversation meta: %v", err)
		}

		// Include the new delivered report field in the response.
		deliveredReport[deliveredUUIDField] = currentTimestamp

		// Marshal the updated delivered report for the response.
		deliveredReportBytes, err := json.Marshal(deliveredReport)
		if err != nil {
			return err
		}

		// Create the response payload for the delivered report.
		deliveredReportResponse := NewResponse(deliveredDeviceID, common.TopicMessageTypeDeliveredDR, string(deliveredReportBytes))

		// Send the payload via WebSocket to the conversation creator.
		if err := finalClient.SendPayloadToClientConnection(deliveredReportResponse); err != nil {
			return err
		}
	}

	return nil
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

	// Extract chatroom_id and conversation_ids from GET parameters.
	chatroomID := c.Query(common.ParamChatroomID)
	if chatroomID == "" {
		api.GeneralBadRequestError(c, common.ErrorChatroomIDMissing)
		return
	}

	conversationIDsParam := c.Query(common.ParamConversationIDs)
	if conversationIDsParam == "" || conversationIDsParam == "null" {
		api.GeneralBadRequestError(c, common.ErrorConversationIDsMissing)
		return
	}

	// Parse conversation_ids into an array of strings using Unmarshal.
	var conversationIDs []string
	if err := json.Unmarshal([]byte(conversationIDsParam), &conversationIDs); err != nil {
		api.GeneralBadRequestError(c, common.ErrorConversationIDsMissing)
		return
	}

	// Get Redis client from context.
	redisClient := GetRedisClientFromContext(c)

	// Fetch the conversation delivery report data for all provided conversation IDs.
	conversationKeys := make([]string, len(conversationIDs))
	for i, conversationID := range conversationIDs {
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
		api.GeneralAPIError(c, fmt.Sprintf(common.ErrorFailedCacheFetchRedis, err))
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
			common.DeliveryCount:               metaMap[common.DeliveryCount],
			common.SenderUUID:                  metaMap[common.SenderUUID],
			common.TopicMessageTypeDeliveredDR: extractDeliveredDRFields(data),
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
		if strings.HasPrefix(key, common.DRUserDelivered) {
			// Remove the prefix from the key before adding it to the map.
			strippedKey := strings.TrimPrefix(key, common.DRUserDelivered)
			deliveredDR[strippedKey] = value
		}
	}
	return deliveredDR
}

// UpdateReadDRWithConversationID updates the delivered report in Redis and sends a payload to the conversation creator.
func UpdateReadDRWithConversationID(redisClient *redis.Client, wsServerParent *ws.WsServerParent, chatroomID, conversationID interface{}, deliveredDeviceID, deliveredUUID string, communityID interface{}) error {
	// Fetch and update the dr_conversation_<conversation_id>
	conversationKey := fmt.Sprintf(common.DRConversationPrefix, conversationID)
	return UpdateReadDR(redisClient, wsServerParent, chatroomID, conversationKey, deliveredDeviceID, deliveredUUID, communityID)
}

// UpdateReadDR updates the read report in Redis and sends a payload to the conversation creator.
func UpdateReadDR(redisClient *redis.Client, wsServerParent *ws.WsServerParent, chatroomID interface{}, conversationKey, readDeviceID, readUUID string, communityID interface{}) error {
	// Fetch the conversation delivery report from Redis.
	conversationData, err := FetchFieldFromHashSet(redisClient, conversationKey, common.DRConversationMetaPrefix)
	if err != nil {
		return err
	}

	// Unmarshal the fetched data into a map.
	var conversationMap map[string]interface{}
	if err := json.Unmarshal([]byte(conversationData), &conversationMap); err != nil {
		return fmt.Errorf(common.ErrorUnmarshalErrorJson, err)
	}

	// Extract the sender UUID from the fetched conversation data.
	senderUUID, _ := conversationMap["sender_uuid"].(string)
	// If the message sender is the same as the message read user, no need to update
	if senderUUID == readUUID {
		return nil
	}

	// Construct the field for the read report using the new key format.
	readUUIDField := fmt.Sprintf(common.DRUserReadPrefix, readUUID)
	// Check if the read report for this user already exists.
	existingReadReport, err := FetchFieldFromHashSet(redisClient, conversationKey, readUUIDField)

	// If the read report already exists, no need to update.
	if existingReadReport != "" {
		return nil
	}

	// Set the current timestamp as the read timestamp.
	currentTimestamp := time.Now().UnixMilli()
	// Update the Redis key with the new read report.
	err = SaveHashSet(redisClient, conversationKey, readUUIDField, currentTimestamp, common.DeliveryReportTTL)
	if err != nil {
		return err
	}

	topicChatroom := fmt.Sprintf(common.TopicTypeChatroomDynamic, chatroomID)
	topicCommunity := fmt.Sprintf(common.TopicTypeCommunityDynamic, communityID)
	// Send the updated read report to the conversation creator's connection.
	clientConnectedToChatroom := wsServerParent.GetConnectionFromWsServer(topicChatroom, senderUUID)
	clientConnectedToCommunity := wsServerParent.GetConnectionFromWsServer(topicCommunity, senderUUID)
	finalClient := clientConnectedToChatroom
	if finalClient == nil {
		finalClient = clientConnectedToCommunity
	}
	if finalClient != nil {
		// Include the new read report field in the response.
		conversationMap[readUUIDField] = currentTimestamp

		// Marshal the updated read report for the response.
		readReportBytes, err := json.Marshal(conversationMap)
		if err != nil {
			return err
		}

		// Create the response payload for the read report.
		readReportResponse := NewResponse(readDeviceID, common.TopicMessageTypeReadDR, string(readReportBytes))

		// Send the payload via WebSocket to the conversation creator.
		if err := finalClient.SendPayloadToClientConnection(readReportResponse); err != nil {
			return err
		}
	}

	return nil
}
