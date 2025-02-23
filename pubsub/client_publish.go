package pubsub

import (
	"fmt"
	"likeminds-pandemonium/api"
	"likeminds-pandemonium/api/constant"
	"likeminds-pandemonium/common"
	"likeminds-pandemonium/ws"
	"log"

	"github.com/gin-gonic/gin"
)

func Publish(c *gin.Context) {
	PublishWithMethod(c, constant.POSTMethod)
}

func PublishWithMethod(c *gin.Context, method int) {
	switch method {
	case constant.POSTMethod:
		topic := c.Param(common.ParamTopic)
		topicSplit, err := GetTopicSplit(topic)
		if err != nil {
			api.GeneralBadRequestError(c, err.Error())
			return
		}

		topicMessageType := c.Query(common.ParamTopicMessageType)
		if topicMessageType == "" || topicMessageType == "null" {
			api.GeneralBadRequestError(c, common.ErrorTopicMessageTypeMissing)
			return
		}

		switch topicSplit[0] {
		case common.TopicTypeChatroom:
			switch topicMessageType {
			case common.TopicMessageTypeDeliveredDR:
				go updateDeliveredDROnPublish(c, topicSplit)
			case common.TopicMessageTypeReadDR:
				go updateReadDROnPublish(c, topicSplit)
			}

		case common.TopicTypeCommunity:
			switch topicMessageType {
			case common.TopicMessageTypeDeliveredDR:
				go updateDeliveredDROnPublish(c, topicSplit)
			}
		}
		api.GenerateResponse(c, nil)
	}
}

type PublishDeliveredDRRequest struct {
	MinTimestamp string      `json:"min_timestamp"`
	MaxTimestamp string      `json:"max_timestamp"`
	CommunityID  interface{} `json:"community_id"`
}

func updateDeliveredDROnPublish(c *gin.Context, topicSplit []string) {
	deliveredUUID := c.GetHeader(constant.HeadersMemberID)
	if deliveredUUID == "" || deliveredUUID == "null" {
		api.GeneralUnauthorizedError(c, common.ErrorUserUUIDMissing)
		return
	}

	deliveredDeviceID := c.GetHeader(constant.HeadersDeviceID)

	// Get the min_timestamp and max_timestamp from the body.
	var publishDeliveredDRRequest PublishDeliveredDRRequest
	if err := c.ShouldBindJSON(&publishDeliveredDRRequest); err != nil {
		api.GeneralBadRequestError(c, fmt.Sprintf(common.ErrorInvalidJSONFormat, err))
		return
	}

	// Construct the Redis key for the chatroom delivery report.
	chatroomID := topicSplit[1]
	chatroomKey := fmt.Sprintf(common.DRChatroomPrefix, chatroomID)

	// Get the Redis client from the context.
	redisClient := GetRedisClientFromContext(c)
	wsServerParent := ws.GetWsServerParentFromContext(c)

	// Step 1: Fetch all member values having score between min and max timestamps from the Redis key dr_chatroom_<chatroom_id>
	chatroomKeyMembers, err := FetchMembersFromZSet(redisClient, chatroomKey, publishDeliveredDRRequest.MinTimestamp, publishDeliveredDRRequest.MaxTimestamp)
	if err != nil {
		api.GeneralAPIError(c, err.Error())
		return
	}

	// Step 2: Iterate over each member value, which corresponds to conversation IDs.
	for _, chatroomKeyMember := range chatroomKeyMembers {
		// Construct the key for the conversation delivery report.
		conversationKey := fmt.Sprintf("%s", chatroomKeyMember.Member)

		// Update the delivered report using the common function.
		if err := UpdateDeliveredDR(redisClient, wsServerParent, deliveredUUID, deliveredDeviceID, publishDeliveredDRRequest.CommunityID, chatroomID, conversationKey); err != nil {
			log.Println(err)
			continue
		}
	}
}

type PublishReadDRRequest struct {
	MinTimestamp string      `json:"min_timestamp"`
	MaxTimestamp string      `json:"max_timestamp"`
	CommunityID  interface{} `json:"community_id"`
}

func updateReadDROnPublish(c *gin.Context, topicSplit []string) {
	chatroomID := topicSplit[1]

	readUUID := c.GetHeader(constant.HeadersMemberID)
	if readUUID == "" || readUUID == "null" {
		api.GeneralUnauthorizedError(c, common.ErrorUserUUIDMissing)
		return
	}

	readDeviceID := c.GetHeader(constant.HeadersDeviceID)

	// Get the min_timestamp and max_timestamp from the body.
	var publishReadDRRequest PublishReadDRRequest
	if err := c.ShouldBindJSON(&publishReadDRRequest); err != nil {
		api.GeneralBadRequestError(c, fmt.Sprintf(common.ErrorInvalidJSONFormat, err))
		return
	}

	// Construct the Redis key for the chatroom delivery report.
	chatroomKey := fmt.Sprintf(common.DRChatroomPrefix, chatroomID)

	// Get the Redis client from the context.
	redisClient := GetRedisClientFromContext(c)
	wsServerParent := ws.GetWsServerParentFromContext(c)

	// Step 1: Fetch all member values between min and max timestamps from the Redis key.
	chatroomKeyMembers, err := FetchMembersFromZSet(redisClient, chatroomKey, publishReadDRRequest.MinTimestamp, publishReadDRRequest.MaxTimestamp)
	if err != nil {
		api.GeneralAPIError(c, err.Error())
		return
	}

	// Step 2: Iterate over each member value, which corresponds to conversation IDs.
	for _, chatroomKeyMember := range chatroomKeyMembers {
		// Construct the key for the conversation delivery report.
		conversationKey := fmt.Sprintf("%s", chatroomKeyMember.Member)

		// Update the read report using the common function.
		if err := UpdateReadDR(redisClient, wsServerParent, readUUID, readDeviceID, publishReadDRRequest.CommunityID, chatroomID, conversationKey); err != nil {
			log.Println(err)
		}
	}
}
