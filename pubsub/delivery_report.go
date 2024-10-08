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

// UpdateSentReport Function to update the cache and send payload for Sent Report
func UpdateSentReport(redisClient *redis.Client, wsServerParent *ws.WsServerParent, topic string, response *Response) error {
	var conversationResponse models.ConversationResponse
	if err := json.Unmarshal([]byte(response.RawData), &conversationResponse); err != nil {
		return fmt.Errorf(common.ErrorUnmarshalErrorJson, err)
	}

	conversationID := conversationResponse.Conversation.ID
	userUUID := conversationResponse.Conversation.Member.UUID
	// Create a CommunityDeliveryReport struct instance
	sentReport := DeliveryReport{
		Timestamp:      time.Now().UnixMilli(),
		ConversationID: conversationID,
		UserUUID:       userUUID,
	}
	// Marshal the payload into JSON bytes
	sentReportBytes, err := json.Marshal(sentReport)
	if err != nil {
		return fmt.Errorf(common.ErrorMarshalErrorJson, err)
	}
	sentReportResponse := NewResponse(response.DeviceID, common.TopicMessageTypeSentReport, string(sentReportBytes))
	// Generate the cache key for the sent report
	communityID := conversationResponse.Conversation.CommunityID
	cacheKey := fmt.Sprintf(common.CommunityDeliveryReportPrefix, communityID)
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

// UpdateDeliveredReport Function to update the cache and send payload for Delivered Report using ZSet
func UpdateDeliveredReport(redisClient *redis.Client, wsServerParent *ws.WsServerParent, topic string, chatroomID interface{}, deviceID, senderUUID, receiverUUID string) error {
	// Create a CommunityDeliveryReport struct instance with the current timestamp
	deliveryReport := DeliveryReport{
		Timestamp: time.Now().UnixMilli(),
		UserUUID:  receiverUUID,
	}
	// Marshal the CommunityDeliveryReport struct into JSON bytes
	reportBytes, err := json.Marshal(deliveryReport)
	if err != nil {
		return fmt.Errorf(common.ErrorMarshalErrorJson, err)
	}

	// Create a response for the delivered report
	reportResponse := NewResponse(deviceID, common.TopicMessageTypeDeliveredReport, string(reportBytes))
	// Generate the cache key for the delivered report
	cacheKey := fmt.Sprintf(common.ChatroomDeliveryReportPrefix, chatroomID)
	// Use the generic SaveToCacheGeneric function to save the value to Redis (ZSet equivalent)
	err = SaveZSet(redisClient, cacheKey, float64(deliveryReport.Timestamp), receiverUUID, 7*24*time.Hour)
	if err != nil {
		return err
	}

	// Send payload to the conversation creator if the user is still active
	client := wsServerParent.GetConnectionFromWsServer(topic, senderUUID)
	if client != nil {
		// Send the report response via WebSocket to the user who created the conversation
		if err := client.SendPayloadToClientConnection(reportResponse); err != nil {
			return err
		}
	}

	return nil
}

type CommunityDeliveryReportRequest struct {
	CommunityID interface{} `json:"community_id"`
}

func CommunityDeliveryReport(c *gin.Context) {
	// Step 1: Parse the request body to get the chatroom_id, user_uuids, page, and page_size
	var requestBody CommunityDeliveryReportRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		api.GeneralBadRequestError(c, common.ErrorInvalidJSONFormat)
		return
	}

	// Step 2: Get chatroom_id from body
	communityID := requestBody.CommunityID
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

	// Create Redis key for delivery report
	communityDRKey := fmt.Sprintf(common.CommunityDeliveryReportPrefix, communityID)
	userDRField := fmt.Sprintf(common.UserDeliveryReportFieldPrefix, memberID)

	// Get Redis client from context
	redisClient := GetRedisClientFromContext(c)

	// Fetch the delivery report from Redis using HGET
	drCacheValue, err := FetchFieldFromHashSet(redisClient, communityDRKey, userDRField)
	if err != nil {
		api.GeneralStatusNotFoundError(c, fmt.Sprintf(common.ErrorFailedCacheFetchRedis, err))
		return
	}

	if drCacheValue == "" {
		api.GeneralStatusNotFoundError(c, common.ErrorNoDeliveryReportFound)
		return
	}

	// Unmarshal the drCacheValue into the drResponse structure
	var drResponse Response
	err = json.Unmarshal([]byte(drCacheValue), &drResponse)
	if err != nil {
		api.GeneralAPIError(c, fmt.Sprintf(common.ErrorUnmarshalErrorJson, err))
		return
	}

	api.GenerateResponse(c, drResponse)
}

type ChatroomDeliveryReportRequest struct {
	ChatroomID interface{} `json:"chatroom_id"`
	UserUUIDs  []string    `json:"user_uuids"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
}

type ChatroomDeliveryReportResponse struct {
	Reports  []map[string]interface{} `json:"reports"`
	Page     int                      `json:"page"`
	PageSize int                      `json:"page_size"`
	Total    int                      `json:"total"`
}

func ChatroomDeliveryReport(c *gin.Context) {
	// Step 1: Parse the request body to get the chatroom_id, user_uuids, page, and page_size
	var requestBody ChatroomDeliveryReportRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		api.GeneralBadRequestError(c, common.ErrorInvalidJSONFormat)
		return
	}

	// Step 2: Get chatroom_id from body
	chatroomID := requestBody.ChatroomID
	if chatroomID == "" || chatroomID == "null" {
		api.GeneralBadRequestError(c, common.ErrorChatroomIDMissing)
		return
	}

	// Step 3: Set default values for page and page size
	page := requestBody.Page
	pageSize := requestBody.PageSize
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	// Step 4: Construct the Redis key for the chatroom delivery report
	redisKey := fmt.Sprintf(common.ChatroomDeliveryReportPrefix, chatroomID)

	// Step 5: Get Redis client from context
	redisClient := GetRedisClientFromContext(c)

	// Step 6: Fetch all members from the ZSet (chatroom delivery reports)
	// We use the FetchMembersFromZSet helper function to get members within the time range
	allMembers, err := FetchMembersFromZSet(redisClient, redisKey, 0, float64(time.Now().UnixMilli()))
	if err != nil {
		api.GeneralAPIError(c, err.Error())
		return
	}

	// Step 7: Filter members if user_uuids were provided in the request body
	var filteredMembers []redis.Z
	if len(requestBody.UserUUIDs) > 0 {
		filteredMembers = filterMembersByUUIDs(allMembers, requestBody.UserUUIDs)
	} else {
		filteredMembers = allMembers
	}

	// Step 8: Apply pagination to the filtered members
	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize
	if endIndex > len(filteredMembers) {
		endIndex = len(filteredMembers)
	}
	paginatedMembers := filteredMembers[startIndex:endIndex]

	// Step 9: Prepare the response data
	var reports []map[string]interface{}
	for _, member := range paginatedMembers {
		reports = append(reports, map[string]interface{}{
			"timestamp": member.Score,
			"user_uuid": member.Member,
		})
	}

	// Step 10: Return the paginated result
	response := ChatroomDeliveryReportResponse{
		Reports:  reports,
		Page:     page,
		PageSize: pageSize,
		Total:    len(filteredMembers), // Total is the size of the filtered members
	}

	// Step 11: Send the response back
	c.JSON(200, response)
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
