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
	err = SaveHashSet(redisClient, conversationKey, "", sentDRValue, 7*24*time.Hour)
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

// UpdateDeliveredDR Function to update the cache and send payload for Delivered DR using ZSet
func UpdateDeliveredDR(redisClient *redis.Client, wsServerParent *ws.WsServerParent, topic string, chatroomID interface{}, deviceID, senderUUID, receiverUUID string) error {
	// Create a SentDR struct instance with the current timestamp
	deliveredDR := DeliveryReport{
		Timestamp: time.Now().UnixMilli(),
		UserUUID:  receiverUUID,
	}

	// Generate the cache key for the delivered DR
	cacheKey := fmt.Sprintf(common.DeliveredDRPrefix, chatroomID)
	// Use the generic SaveToCacheGeneric function to save the value to Redis (ZSet equivalent)
	err := SaveZSet(redisClient, cacheKey, float64(deliveredDR.Timestamp), receiverUUID, 7*24*time.Hour)
	if err != nil {
		return err
	}

	// Send payload to the conversation creator if the user is still active
	client := wsServerParent.GetConnectionFromWsServer(topic, senderUUID)
	if client != nil {
		// Marshal the SentDR struct into JSON bytes
		deliveredDRBytes, err := json.Marshal(deliveredDR)
		if err != nil {
			return fmt.Errorf(common.ErrorMarshalErrorJson, err)
		}
		// Create a response for the delivered DR
		deliveredDRResponse := NewResponse(deviceID, common.TopicMessageTypeDeliveredDR, string(deliveredDRBytes))
		// Send the delivered DR response via WebSocket to the user who created the conversation
		if err := client.SendPayloadToClientConnection(deliveredDRResponse); err != nil {
			return err
		}
	}

	return nil
}

type SentDRRequest struct {
	CommunityID interface{} `json:"community_id"`
}

func SentDR(c *gin.Context) {
	// Step 1: Parse the request body to get the chatroom_id, user_uuids, page, and page_size
	var sentDRRequest SentDRRequest
	if err := c.ShouldBindJSON(&sentDRRequest); err != nil {
		api.GeneralBadRequestError(c, common.ErrorInvalidJSONFormat)
		return
	}

	// Step 2: Get chatroom_id from body
	communityID := sentDRRequest.CommunityID
	if communityID == "" || communityID == "null" {
		api.GeneralBadRequestError(c, common.ErrorCommunityIDMissing)
		return
	}

	// Get x-member-id from headers
	memberID := c.GetHeader(constant.HeadersMemberID)
	if memberID == "" || memberID == "null" {
		api.GeneralUnauthorizedError(c, common.ErrorUserUUIDMissing)
		return
	}

	// Create Redis key for delivered DR
	sentDRKey := fmt.Sprintf(common.SentDRPrefix, communityID)
	userDRField := fmt.Sprintf(common.UserDRFieldPrefix, memberID)

	// Get Redis client from context
	redisClient := GetRedisClientFromContext(c)

	// Fetch the delivered DR from Redis using HGET
	sentDRCacheValue, err := FetchFieldFromHashSet(redisClient, sentDRKey, userDRField)
	if err != nil {
		api.GeneralStatusNotFoundError(c, fmt.Sprintf(common.ErrorFailedCacheFetchRedis, err))
		return
	}
	if sentDRCacheValue == "" {
		api.GeneralStatusNotFoundError(c, common.ErrorNoDRFound)
		return
	}

	// Unmarshal the sentDRCacheValue into the sentDRResponse structure
	var sentDRResponse DeliveryReport
	err = json.Unmarshal([]byte(sentDRCacheValue), &sentDRResponse)
	if err != nil {
		api.GeneralAPIError(c, fmt.Sprintf(common.ErrorUnmarshalErrorJson, err))
		return
	}

	api.GenerateResponse(c, sentDRResponse)
}

type DeliveredDRRequest struct {
	ChatroomID interface{} `json:"chatroom_id"`
	UserUUIDs  []string    `json:"user_uuids"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
}

func DeliveredDR(c *gin.Context) {
	// Step 1: Parse the request body to get the chatroom_id, user_uuids, page, and page_size
	var deliveredDRRequest DeliveredDRRequest
	if err := c.ShouldBindJSON(&deliveredDRRequest); err != nil {
		api.GeneralBadRequestError(c, common.ErrorInvalidJSONFormat)
		return
	}

	// Step 2: Get chatroom_id from body
	chatroomID := deliveredDRRequest.ChatroomID
	if chatroomID == "" || chatroomID == "null" {
		api.GeneralBadRequestError(c, common.ErrorChatroomIDMissing)
		return
	}

	// Step 3: Set default values for page and page size
	page := deliveredDRRequest.Page
	pageSize := deliveredDRRequest.PageSize
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	// Step 4: Construct the Redis key for the chatroom delivered DR
	redisKey := fmt.Sprintf(common.DeliveredDRPrefix, chatroomID)

	// Step 5: Get Redis client from context
	redisClient := GetRedisClientFromContext(c)

	// Step 6: Fetch all members from the ZSet (delivered DR)
	// We use the FetchMembersFromZSet helper function to get members within the time range
	allMembers, err := FetchMembersFromZSet(redisClient, redisKey, 0, float64(time.Now().UnixMilli()))
	if err != nil {
		api.GeneralAPIError(c, err.Error())
		return
	}

	// Step 7: Get the member with the least timestamp (first member in the ZSet, which is allMembers[0])
	var leastMember *redis.Z = nil
	if len(allMembers) > 0 {
		leastMember = &allMembers[0] // The first member in the ZSet has the least timestamp
	}

	// Step 8: Filter members if user_uuids were provided in the request body
	var filteredMembers []redis.Z
	if len(deliveredDRRequest.UserUUIDs) > 0 {
		filteredMembers = filterMembersByUUIDs(allMembers, deliveredDRRequest.UserUUIDs)
	} else {
		filteredMembers = allMembers
	}

	// Step 9: Apply pagination to the filtered members
	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize
	if endIndex > len(filteredMembers) {
		endIndex = len(filteredMembers)
	}
	paginatedMembers := filteredMembers[startIndex:endIndex]

	// Step 10: Prepare the deliveredDRResponse data as a map of user_uuid: timestamp
	deliveredDRMap := make(map[string]map[string]interface{})
	for _, member := range paginatedMembers {
		deliveredDRMap[member.Member.(string)] = map[string]interface{}{
			"timestamp": member.Score,
		}
	}

	// Step 11: Return the paginated result as a JSON array
	deliveredDRResponse := map[string]interface{}{
		"delivered_dr": deliveredDRMap,
		"page":         page,
		"page_size":    pageSize,
		"total":        len(filteredMembers),
	}
	// Step 12: Prepare the least_delivered_dr data (if applicable)
	if leastMember != nil {
		leastMemberMap := make(map[string]map[string]interface{})
		leastMemberMap[leastMember.Member.(string)] = map[string]interface{}{
			"timestamp": leastMember.Score,
		}
		deliveredDRResponse["least_delivered_dr"] = leastMemberMap
	}

	// Step 13: Send the deliveredDRResponse back
	api.GenerateResponse(c, deliveredDRResponse)
}

// Helper function to filter members by provided user UUIDs
func filterMembersByUUIDs(members []redis.Z, userUUIDs []string) []redis.Z {
	// Create a set of user UUIDs for faster lookup
	uuidSet := make(map[string]bool)
	for _, uuid := range userUUIDs {
		uuidSet[uuid] = true
	}

	// Filter members based on user UUIDs
	var filteredMembers []redis.Z
	for _, member := range members {
		if uuidSet[member.Member.(string)] {
			filteredMembers = append(filteredMembers, member)
		}
	}

	return filteredMembers
}
