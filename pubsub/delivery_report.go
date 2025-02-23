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

func UpdateSentDR(redisClient *redis.Client, wsServerParent *ws.WsServerParent, senderUUID, senderDeviceID string, communityID, chatroomID int, conversationID int64, conversationCreatedAt float64, participantsCount int) error {
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

	// Fine client connection who sent the message to send him sent_dr topic_message_type payload. He can be connected either to chatroom:<chatroom_id> or community:<community_id> topic
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
		sentDRBytes, err := json.Marshal(conversationMetaCache)
		if err != nil {
			return err
		}
		sentDRPSResponse := NewPSResponse(senderDeviceID, common.TopicMessageTypeSentDR, string(sentDRBytes))

		if err := finalClient.SendPayloadToClientConnection(sentDRPSResponse); err != nil {
			return err
		}
	}

	return nil
}

// UpdateDeliveredDR updates the delivered report in Redis and sends a payload to the conversation creator.
func UpdateDeliveredDR(redisClient *redis.Client, wsServerParent *ws.WsServerParent, deliveredUUID, deliveredDeviceID string, communityID, chatroomID interface{}, conversationKey string) error {
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
	existingDeliveredUUIDFieldValue, err := FetchFieldFromHashSet(redisClient, conversationKey, deliveredUUIDField)

	// If the delivered report already exists, no need to update.
	if existingDeliveredUUIDFieldValue != "" {
		return nil
	}

	// Set the current timestamp as the delivered timestamp.
	currentTimestamp := time.Now().UnixMilli()
	// Update the Redis key dr_conversation_<conversation_id> where field will be dr_user_delivered_<uuid> and value will be <current time in millis>
	err = SaveHashSet(redisClient, conversationKey, deliveredUUIDField, currentTimestamp, common.DeliveryReportTTL)
	if err != nil {
		return err
	}

	// Fine client connection who sent the message to send him delivered_dr topic_message_type payload. He can be connected either to chatroom:<chatroom_id> or community:<community_id> topic
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
		var deliveredDRMap map[string]interface{}
		deliveredDRMap[common.DRConversationMetaPrefix] = conversationMetaCache
		// Include the new delivered report uuid
		deliveredDRMap[deliveredUUIDField] = currentTimestamp

		// Marshal the updated delivered report for the response.
		deliveredDRMapBytes, err := json.Marshal(deliveredDRMap)
		if err != nil {
			return err
		}
		deliveredDRPSResponse := NewPSResponse(deliveredDeviceID, common.TopicMessageTypeDeliveredDR, string(deliveredDRMapBytes))

		// Send the payload via WebSocket to the conversation creator.
		if err := finalClient.SendPayloadToClientConnection(deliveredDRPSResponse); err != nil {
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
	userUUID := c.GetHeader(constant.HeadersMemberID)
	if userUUID == "" || userUUID == "null" {
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
	conversationKeyCMDs := make([]*redis.MapStringStringCmd, len(conversationKeys))
	for i, conversationKey := range conversationKeys {
		conversationKeyCMDs[i] = pipe.HGetAll(c, conversationKey)
	}
	_, err := pipe.Exec(c)
	if err != nil {
		api.GeneralAPIError(c, fmt.Sprintf(common.ErrorFailedCacheFetchRedis, err))
		return
	}

	// Construct the deliveryReportResponse.
	deliveryReportMap := make(map[string]map[string]interface{})
	for i, conversationKeyCMD := range conversationKeyCMDs {
		conversationID := conversationIDs[i]

		// Get the result of the HGETALL command.
		conversationKeyValue, err := conversationKeyCMD.Result()
		if err != nil {
			log.Printf("Error fetching conversationKeyValue for conversation %s: %v", conversationID, err)
			continue
		}

		// Extract the "dr_conversation_meta" field.
		conversationMetaCacheValue, ok := conversationKeyValue[common.DRConversationMetaPrefix]
		if !ok || conversationMetaCacheValue == "" {
			log.Printf("Missing or empty dr_conversation_meta for conversation %s", conversationID)
			continue
		}

		/*// Unmarshal the metadata field. todo might not required this
		var conversationMetaCache models.ConversationMetaCache
		if err := json.Unmarshal([]byte(conversationMetaCacheValue), &conversationMetaCache); err != nil {
			log.Printf("Error unmarshalling conversation meta for %s: %v", conversationID, err)
			continue
		}*/

		// Add the conversation conversationKeyValue to the delivery report map directly.
		deliveryReportMap[conversationID] = map[string]interface{}{
			common.DRConversationMetaPrefix:    conversationMetaCacheValue,
			common.TopicMessageTypeDeliveredDR: extractDeliveredDRFieldValue(conversationKeyValue),
			common.TopicMessageTypeReadDR:      extractReadDRFieldValue(conversationKeyValue),
		}
	}

	// Create and send the deliveryReportResponse.
	deliveryReportResponse := DeliveryReportResponse{
		DeliveryReport: deliveryReportMap,
	}
	api.GenerateResponse(c, deliveryReportResponse)
}

// extractDeliveredDRFieldValue extracts the delivered reports from the Redis data map.
func extractDeliveredDRFieldValue(conversationKeyValue map[string]string) map[string]interface{} {
	deliveredUUIDMap := make(map[string]interface{})

	// Iterate through all fields in the Redis hash and find delivered report fields.
	for field, value := range conversationKeyValue {
		if strings.HasPrefix(field, common.DRUserDelivered) {
			// Remove the prefix from the field before adding it to the map.
			deliveredUUIDFieldStripped := strings.TrimPrefix(field, common.DRUserDelivered)
			deliveredUUIDMap[deliveredUUIDFieldStripped] = value
		}
	}
	return deliveredUUIDMap
}

// extractReadDRFieldValue extracts the delivered reports from the Redis data map.
func extractReadDRFieldValue(conversationKeyValue map[string]string) map[string]interface{} {
	readUUIDMap := make(map[string]interface{})

	// Iterate through all fields in the Redis hash and find delivered report fields.
	for key, value := range conversationKeyValue {
		if strings.HasPrefix(key, common.DRUserRead) {
			// Remove the prefix from the key before adding it to the map.
			readUUIDFieldStripped := strings.TrimPrefix(key, common.DRUserRead)
			readUUIDMap[readUUIDFieldStripped] = value
		}
	}
	return readUUIDMap
}

// UpdateReadDR updates the read report in Redis and sends a payload to the conversation creator.
func UpdateReadDR(redisClient *redis.Client, wsServerParent *ws.WsServerParent, readUUID, readDeviceID string, communityID, chatroomID interface{}, conversationKey string) error {
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
	existingReadUUIDFieldValue, err := FetchFieldFromHashSet(redisClient, conversationKey, readUUIDField)

	// If the read report already exists, no need to update.
	if existingReadUUIDFieldValue != "" {
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
	// Fine client connection who sent the message to send him read_dr topic_message_type payload. He can be connected either to chatroom:<chatroom_id> or community:<community_id> topic
	clientConnectedToChatroom := wsServerParent.GetConnectionFromWsServer(topicChatroom, senderUUID)
	clientConnectedToCommunity := wsServerParent.GetConnectionFromWsServer(topicCommunity, senderUUID)
	finalClient := clientConnectedToChatroom
	if finalClient == nil {
		finalClient = clientConnectedToCommunity
	}
	if finalClient != nil {
		var readDRMap map[string]interface{}
		readDRMap[common.DRConversationMetaPrefix] = conversationMetaCache

		// Include the new read report field in the response.
		readDRMap[readUUIDField] = currentTimestamp

		// Marshal the updated read report for the response.
		readDRMapBytes, err := json.Marshal(readDRMap)
		if err != nil {
			return err
		}
		readDRPSResponse := NewPSResponse(readDeviceID, common.TopicMessageTypeReadDR, string(readDRMapBytes))

		// Send the payload via WebSocket to the conversation creator.
		if err := finalClient.SendPayloadToClientConnection(readDRPSResponse); err != nil {
			return err
		}
	}

	return nil
}
