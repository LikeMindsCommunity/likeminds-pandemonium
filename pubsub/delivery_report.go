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

func UpdateSentDR(redisClient *redis.Client, wsServerParent *ws.WsServerParent, senderUUID, senderDeviceID string, chatroomID, communityID int, conversationID int64, conversationCreatedAt float64, participantsCount int) error {
	// Create the cache keys
	chatroomKey := fmt.Sprintf(common.DRChatroomPrefix, chatroomID)
	conversationKey := fmt.Sprintf(common.DRConversationPrefix, conversationID)

	// Update Zset key dr_chatroom_<chatroom_id> where score will be <conversation created at timestamp> and member will be dr_conversation_<conversation_id>
	err := SaveZSet(redisClient, chatroomKey, conversationCreatedAt, conversationKey, common.DeliveryReportTTL)
	if err != nil {
		return err
	}

	// Create ConversationMetaCache object will participantsCount and senderUUID and save in HastSet key dr_conversation_<conversation_id> where field will be dr_conversation_meta and value will be ConversationMetaCache
	var conversationMetaCache = models.ConversationMetaCache{DeliveryCount: participantsCount, SenderUUID: senderUUID}
	err = SaveHashSet(redisClient, conversationKey, common.DRConversationMetaPrefix, conversationMetaCache, common.DeliveryReportTTL)
	if err != nil {
		return err
	}

	// Fine client connection who sent the message to send him sent_dr topic message type payload. He can be connected either to chatroom:<chatroom_id> or community:<community_id> topic
	topicChatroom := fmt.Sprintf(common.TopicTypeChatroomDynamic, chatroomID)
	topicCommunity := fmt.Sprintf(common.TopicTypeCommunityDynamic, communityID)
	clientConnectedToChatroom := wsServerParent.GetConnectionFromWsServer(topicChatroom, senderUUID)
	clientConnectedToCommunity := wsServerParent.GetConnectionFromWsServer(topicCommunity, senderUUID)
	finalClient := clientConnectedToChatroom
	if finalClient == nil {
		finalClient = clientConnectedToCommunity
	}
	if finalClient != nil {
		// Marshal the updated sent report for the response.
		sentReportBytes, err := json.Marshal(conversationMetaCache)
		if err != nil {
			return err
		}

		payload := NewPSResponse(senderDeviceID, common.TopicMessageTypeSentDR, string(sentReportBytes))
		if err := finalClient.SendPayloadToClientConnection(payload); err != nil {
			return err
		}
	}

	return nil
}

// UpdateDeliveredDRWithConversationID updates the delivered report in Redis and sends a payload to the conversation creator.
func UpdateDeliveredDRWithConversationID(redisClient *redis.Client, wsServerParent *ws.WsServerParent, communityID, chatroomID, conversationID interface{}, deliveredUUID, deliveredDeviceID string) error {
	// Fetch and update the dr_conversation_<conversation_id> with delivered report
	conversationKey := fmt.Sprintf(common.DRConversationPrefix, conversationID)
	return UpdateDeliveredDR(redisClient, wsServerParent, communityID, chatroomID, conversationKey, deliveredUUID, deliveredDeviceID)
}

// UpdateDeliveredDR updates the delivered report in Redis and sends a payload to the conversation creator.
func UpdateDeliveredDR(redisClient *redis.Client, wsServerParent *ws.WsServerParent, communityID, chatroomID interface{}, conversationKey, deliveredUUID, deliveredDeviceID string) error {
	// Fetch the dr_conversation_meta field from Redis.
	conversationMetaCacheValue, err := FetchFieldFromHashSet(redisClient, conversationKey, common.DRConversationMetaPrefix)
	if err != nil {
		return err
	}

	// Unmarshal the fetched data into ConversationMetaCache
	var conversationMetaCache models.ConversationMetaCache
	if err := json.Unmarshal([]byte(conversationMetaCacheValue), &conversationMetaCache); err != nil {
		return fmt.Errorf(common.ErrorUnmarshalErrorJson, err)
	}

	// Extract the sender UUID from the fetched dr_conversation_meta
	senderUUID := conversationMetaCache.SenderUUID
	if senderUUID == "" {
		return fmt.Errorf(common.ErrorSenderUUIDMissing, conversationKey)
	}
	// If the message sender is the same as the message delivered user, no need to update
	if senderUUID == deliveredUUID {
		return nil
	}

	// Construct the field for the delivered report using the new key format.
	deliveredUUIDField := fmt.Sprintf(common.DRUserDeliveredPrefix, deliveredUUID)
	// Check if the delivered report for this user already exists.
	existingDeliveredReport, err := FetchFieldFromHashSet(redisClient, conversationKey, deliveredUUIDField)

	// If the delivered report already exists, no need to update.
	if existingDeliveredReport != "" {
		return nil
	}

	// Set the current timestamp as the delivered timestamp.
	currentTimestamp := time.Now().UnixMilli()
	// Update the Redis key dr_conversation_<conversation_id> where field will be dr_user_delivered_<uuid> and value will be <current time in millis>
	err = SaveHashSet(redisClient, conversationKey, deliveredUUIDField, currentTimestamp, common.DeliveryReportTTL)
	if err != nil {
		return err
	}

	// Fine client connection who sent the message to send him delivered_dr topic message type payload. He can be connected either to chatroom:<chatroom_id> or community:<community_id> topic
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
		var deliveredReportMap map[string]interface{}
		deliveredReportMap[common.DRConversationMetaPrefix] = conversationMetaCache
		// Include the new delivered report uuid
		deliveredReportMap[deliveredUUIDField] = currentTimestamp

		// Marshal the updated delivered report for the response.
		deliveredReportMapBytes, err := json.Marshal(deliveredReportMap)
		if err != nil {
			return err
		}

		deliveredReportPSResponse := NewPSResponse(deliveredDeviceID, common.TopicMessageTypeDeliveredDR, string(deliveredReportMapBytes))

		// Send the payload via WebSocket to the conversation creator.
		if err := finalClient.SendPayloadToClientConnection(deliveredReportPSResponse); err != nil {
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
		conversationMetaCacheValue, ok := data[common.DRConversationMetaPrefix]
		if !ok || conversationMetaCacheValue == "" {
			log.Printf("Missing or empty dr_conversation_meta for conversation %s", conversationID)
			continue
		}

		// Unmarshal the metadata field.
		var conversationMetaCache models.ConversationMetaCache
		if err := json.Unmarshal([]byte(conversationMetaCacheValue), &conversationMetaCache); err != nil {
			log.Printf("Error unmarshalling conversation meta for %s: %v", conversationID, err)
			continue
		}

		// Add the conversation data to the delivery report map directly.
		deliveryReport[conversationID] = map[string]interface{}{
			common.DRConversationMetaPrefix:    conversationMetaCache,
			common.TopicMessageTypeDeliveredDR: extractDeliveredDRFields(data),
			common.TopicMessageTypeReadDR:      extractReadDRFields(data),
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

// extractReadDRFields extracts the delivered reports from the Redis data map.
func extractReadDRFields(data map[string]string) map[string]interface{} {
	readDR := make(map[string]interface{})

	// Iterate through all fields in the Redis hash and find delivered report fields.
	for key, value := range data {
		if strings.HasPrefix(key, common.DRUserRead) {
			// Remove the prefix from the key before adding it to the map.
			strippedKey := strings.TrimPrefix(key, common.DRUserRead)
			readDR[strippedKey] = value
		}
	}
	return readDR
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
	conversationMetaCacheValue, err := FetchFieldFromHashSet(redisClient, conversationKey, common.DRConversationMetaPrefix)
	if err != nil {
		return err
	}

	// Unmarshal the fetched data into a map.
	var conversationMetaCache models.ConversationMetaCache
	if err := json.Unmarshal([]byte(conversationMetaCacheValue), &conversationMetaCache); err != nil {
		return fmt.Errorf(common.ErrorUnmarshalErrorJson, err)
	}

	// Extract the sender UUID from the fetched conversation data.
	senderUUID := conversationMetaCache.SenderUUID
	if senderUUID == "" {
		return fmt.Errorf(common.ErrorSenderUUIDMissing, conversationKey)
	}

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
	// Update the Redis key dr_conversation_<conversation_id> where field will be dr_user_read_<uuid> and value will be <current time in millis>
	err = SaveHashSet(redisClient, conversationKey, readUUIDField, currentTimestamp, common.DeliveryReportTTL)
	if err != nil {
		return err
	}

	topicChatroom := fmt.Sprintf(common.TopicTypeChatroomDynamic, chatroomID)
	topicCommunity := fmt.Sprintf(common.TopicTypeCommunityDynamic, communityID)
	// Fine client connection who sent the message to send him read_dr topic message type payload. He can be connected either to chatroom:<chatroom_id> or community:<community_id> topic
	clientConnectedToChatroom := wsServerParent.GetConnectionFromWsServer(topicChatroom, senderUUID)
	clientConnectedToCommunity := wsServerParent.GetConnectionFromWsServer(topicCommunity, senderUUID)
	finalClient := clientConnectedToChatroom
	if finalClient == nil {
		finalClient = clientConnectedToCommunity
	}
	if finalClient != nil {
		var readReportMap map[string]interface{}
		readReportMap[common.DRConversationMetaPrefix] = conversationMetaCache

		// Include the new read report field in the response.
		readReportMap[readUUIDField] = currentTimestamp

		// Marshal the updated read report for the response.
		readReportBytes, err := json.Marshal(readReportMap)
		if err != nil {
			return err
		}

		readReportPSResponse := NewPSResponse(readDeviceID, common.TopicMessageTypeReadDR, string(readReportBytes))

		// Send the payload via WebSocket to the conversation creator.
		if err := finalClient.SendPayloadToClientConnection(readReportPSResponse); err != nil {
			return err
		}
	}

	return nil
}
