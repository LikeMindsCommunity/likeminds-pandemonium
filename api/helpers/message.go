package helpers

import (
	"encoding/json"
	"fmt"
	"likeminds-pandemonium/api/constant"
	"likeminds-pandemonium/api/models"
	"likeminds-pandemonium/api/repository"
	requestresponse "likeminds-pandemonium/api/request_response"
	"likeminds-pandemonium/api/utilities"
	"likeminds-pandemonium/external"
	"log"
	"regexp"
	"slices"
	"strconv"
	"time"
)

type CreateMessageRequestContext struct {
	UserInfo        models.UserInfo  `json:"userinfo"`
	Chatroom        models.Chatroom  `json:"chatroom"`
	Community       models.Community `json:"community"`
	OriginalMessage models.Message   `json:"original_message"`
}

func ValidateCreateMessageRequest(createMessageRequest *requestresponse.CreateMessageRequest, userID string, deviceID string, topic string) (*CreateMessageRequestContext, *constant.APIError) {

	createMessageRequestContext := &CreateMessageRequestContext{}

	err := validateWsChatroomID(createMessageRequest.ChatroomID, topic)
	if err != nil {
		return nil, constant.APIErrorBadRequest(fmt.Errorf("failed to validate ws chatroom id, err=%s", err))
	}

	userinfo, err := GetUserInfoByUUID(userID)
	if err != nil {
		return nil, constant.APIErrorBadRequest(fmt.Errorf("failed to get user, user id=%s, err=%s", userID, err))
	}
	createMessageRequestContext.UserInfo = *userinfo

	chatroom, err := GetChatroomByID(int64(createMessageRequest.ChatroomID))
	if err != nil {
		return nil, constant.APIErrorBadRequest(fmt.Errorf("failed to get chatroom, chatroom id=%d, err=%s", createMessageRequest.ChatroomID, err))
	}
	createMessageRequestContext.Chatroom = *chatroom

	if chatroom.IsSecret {
		validateUserSecretChatroomAccess, err := ValidateUserSecretChatroomAccess([]byte(*chatroom.SecretChatroomParticipants), userinfo.UserID)
		if err != nil {
			return nil, constant.APIErrorBadRequest(fmt.Errorf("failed to validate secret chatroom access, chatroom id=%d, user id=%d, err=%s", createMessageRequest.ChatroomID, userinfo.UserID, err))
		}
		if !validateUserSecretChatroomAccess {
			return nil, constant.APIErrorBadRequest(fmt.Errorf("failed to get chatroom for user, user id=%d, err=%s", userinfo.UserID, err))
		}
	}

	community, err := GetCommunityByID(int64(chatroom.CommunityID))
	if err != nil {
		return nil, constant.APIErrorBadRequest(fmt.Errorf("failed to get community, community id=%d, err=%s", chatroom.CommunityID, err))
	}
	createMessageRequestContext.Community = *community

	if createMessageRequest.RepliedConversationId != nil {
		originalMessage, err := GetMessageByID(createMessageRequest.RepliedConversationId.(float64))
		if err != nil {
			return nil, constant.APIErrorBadRequest(fmt.Errorf("failed to get replied message, message id=%f, err=%s", createMessageRequest.RepliedConversationId.(float64), err))
		}
		createMessageRequestContext.OriginalMessage = *originalMessage
	}

	if createMessageRequest.State == int32(constant.MessageStateMessagePoll) {
		err := validateCreatePollMessageRequest(createMessageRequest, community.ID)
		if err != nil {
			return nil, constant.APIErrorBadRequest(fmt.Errorf("failed to validate create poll message request, err=%s", err))
		}
	}

	return createMessageRequestContext, nil
}

func GetMessageByID(messageID float64) (*models.Message, error) {
	messageIDStr := fmt.Sprintf("%g", messageID)
	message, err := repository.GetMessageByID(messageIDStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get message, message id=%f, err=%s", messageID, err)
	}
	return message, nil
}

func validateCreatePollMessageRequest(createMessageRequest *requestresponse.CreateMessageRequest, communityID int64) error {
	communityConfigurationChatPollDefault := &constant.CommunityConfiguration{}
	configuration, err := GetCommunityConfiguration(int(communityID), constant.CommunityConfigurationChatPoll)
	if err != nil && err.Error() != "record not found" {
		log.Printf("failed to get community configuration, community=%d, configuration=%s, err=%s", communityID, constant.CommunityConfigurationChatPoll, err)
	}
	if err != nil && err.Error() == "record not found" {
		log.Printf("using default, failed to get community configuration, community=%d, configuration=%s, err=%s", communityID, constant.CommunityConfigurationChatPoll, err)
		communityConfigurationChatPollDefault = constant.CommunityConfigurationChatPollDefault
	}

	if communityConfigurationChatPollDefault.Type == constant.CommunityConfigurationChatPoll {
		err := validateCreatePollMessageDefaultCommunityConfiguration(createMessageRequest)
		if err != nil {
			return fmt.Errorf("failed to validate create poll message against default community configuration")
		}
	} else {
		communityConfigurationChatPollValue := configuration.Value
		var communityConfigurationValueStruct constant.CommunityConfigurationValue
		err := json.Unmarshal([]byte(communityConfigurationChatPollValue), &communityConfigurationValueStruct)
		if err != nil {
			return fmt.Errorf("failed to unmarshal community configuration chat poll")
		}
		err = validateCreatePollMessageCustomerCommunityConfiguration(createMessageRequest, &communityConfigurationValueStruct)
		if err != nil {
			return fmt.Errorf("failed to validate create poll message against community configuration")
		}
	}

	return nil
}

func validateCreatePollMessageDefaultCommunityConfiguration(createMessageRequest *requestresponse.CreateMessageRequest) error {
	if createMessageRequest.PollType == nil {
		return fmt.Errorf("poll_type is missing, required")
	}

	if createMessageRequest.ExpiryTime == nil && !createMessageRequest.NoPollExpiry {
		return fmt.Errorf("expiry_time is missing, required")
	}

	var pollType = new(int)
	*pollType = int(*createMessageRequest.PollType)
	if *pollType == constant.MessagePollTypeDeferred && !createMessageRequest.NoPollExpiry {
		return fmt.Errorf("expiry_time is missing, required")
	}

	if createMessageRequest.NoPollExpiry {
		createMessageRequest.ExpiryTime = nil
	}

	if !createMessageRequest.AllowVoteChange {
		if *pollType == constant.MessagePollTypeInstant {
			createMessageRequest.AllowVoteChange = false
		} else {
			createMessageRequest.AllowVoteChange = true
		}
	}

	return nil
}

func validateCreatePollMessageCustomerCommunityConfiguration(createMessageRequest *requestresponse.CreateMessageRequest, communityConfigurationValueStruct *constant.CommunityConfigurationValue) error {
	if communityConfigurationValueStruct.AllowOveride {
		return nil
	}

	if communityConfigurationValueStruct.PollType != "" {
		pollTypeInt := constant.GetMessagePollTypeFromEnum(communityConfigurationValueStruct.PollType)
		var pollTypeInt32 = new(int32)
		*pollTypeInt32 = int32(pollTypeInt)
		createMessageRequest.PollType = pollTypeInt32
	}

	if communityConfigurationValueStruct.MultipleSelectState != "" {
		multiSelectStateInt := constant.GetMessagePollMultiSelectStateFromEnum(communityConfigurationValueStruct.MultipleSelectState)
		var multiSelectStateInt64 = new(int64)
		*multiSelectStateInt64 = int64(multiSelectStateInt)
		createMessageRequest.MultilpleSelectState = multiSelectStateInt64

	}

	var multilpleSelectNoDefault = new(int64)
	*multilpleSelectNoDefault = int64(1)
	if createMessageRequest.MultilpleSelectNo = multilpleSelectNoDefault; communityConfigurationValueStruct.MultipleSelectNo != 0 {
		var multilpleSelectNo = new(int64)
		*multilpleSelectNo = int64(communityConfigurationValueStruct.MultipleSelectNo)
		createMessageRequest.MultilpleSelectNo = multilpleSelectNo
	}

	createMessageRequest.IsAnonymous = communityConfigurationValueStruct.IsAnonymous
	createMessageRequest.AllowAddOption = communityConfigurationValueStruct.AllowAddOption
	createMessageRequest.NoPollExpiry = communityConfigurationValueStruct.NoPollExpiry
	createMessageRequest.AllowVoteChange = communityConfigurationValueStruct.AllowVoteChange

	if createMessageRequest.NoPollExpiry {
		createMessageRequest.ExpiryTime = nil
	} else {
		if createMessageRequest.ExpiryTime == nil {
			return fmt.Errorf("expiry_time missing, required")
		}
	}

	return nil
}

func ValidateCreateMessagePermission(createMessageRequest *requestresponse.CreateMessageRequest, chatroom *models.Chatroom, userIDInt int) *constant.APIError {
	if chatroom.Type == constant.ChatroomTypeMasterIntro {
		return constant.APIErrorForbidden(fmt.Errorf("failed to post in chatroom, type=%s", constant.ChatroomTypeMasterIntroEnum))
	}

	memberState, err := GetMemberStateInCommunity(chatroom.CommunityID, userIDInt)
	if err != nil {
		return constant.APIErrorForbidden(fmt.Errorf("failed to get member state in community, community id=%d, userId=%d, err=%s", chatroom.CommunityID, userIDInt, err))
	}

	isMemberAdminInCommunity := IsMemberStateAdminInCommunity(memberState)
	if !validateMessageGroupTags(createMessageRequest.Text, chatroom.IsSecret, chatroom.UserID, isMemberAdminInCommunity, userIDInt) {
		return constant.APIErrorForbidden(fmt.Errorf("invalid group tags in message, text=%s", createMessageRequest.Text))
	}

	_, err = ValidateUserRightInCommunity(chatroom.CommunityID, userIDInt, constant.MemberRightRespondInRoom)
	if err != nil {
		return constant.APIErrorForbidden(fmt.Errorf("user right missing in community, community id=%d, user id=%d, right=%s", chatroom.CommunityID, userIDInt, constant.MemberRightRespondInRoomEnum))
	}

	_, err = ValidateUserRightInChatroom(int(chatroom.ID), userIDInt, constant.UserChannelSettingChatroomMemberCanMessage, isMemberAdminInCommunity)
	if err != nil {
		return constant.APIErrorForbidden(fmt.Errorf("user right missing in chatroom, chatroom id=%d, user id=%d, right=%s", chatroom.ID, userIDInt, constant.UserChannelSettingChatroomMemberCanMessage))
	}

	if createMessageRequest.State == int32(constant.MessageStateMessagePoll) {
		_, err := ValidateUserRightInCommunity(chatroom.CommunityID, userIDInt, constant.MemberRightCreatePoll)
		if err != nil {
			return constant.APIErrorForbidden(fmt.Errorf("user right missing in community, community id=%d, user id=%d, right=%s", chatroom.CommunityID, userIDInt, constant.MemberRightCreatePollEnum))
		}
	}

	return nil
}

func ValidateUserSecretChatroomAccess(secretChatroomParticipants []byte, userID int) (bool, error) {
	var secretChatroomParticipantsList []int
	err := json.Unmarshal(secretChatroomParticipants, &secretChatroomParticipantsList)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal secret chatroom participants, err=%s", err)
	}

	if !slices.Contains(secretChatroomParticipantsList, userID) {
		return false, nil
	}

	return true, nil
}

func FillCreateMessageModelInstance(
	createMessageModelInstance *models.Message,
	createMessageRequest *requestresponse.CreateMessageRequest,
	requestContext CreateMessageRequestContext,
	deviceID string,
	isGuest bool,
	platformCode string) *constant.APIError {

	createMessageModelInstance.Answer = createMessageRequest.Text
	createMessageModelInstance.APIVersion = 1

	createMessageModelInstance.HasFiles = createMessageRequest.HasFiles
	createMessageModelInstance.AttachmentCount = int(createMessageRequest.AttachmentCount)
	createMessageModelInstance.AttachmentsUploaded = false
	if createMessageModelInstance.AttachmentCount > 0 {
		createMessageRequest.HasFiles = true
		createMessageModelInstance.HasFiles = true
	}

	createMessageModelInstance.CardID = int(requestContext.Chatroom.ID)
	createMessageModelInstance.CommunityID = int(requestContext.Community.ID)
	createMessageModelInstance.CreatedAt = time.Now().Unix()
	createMessageModelInstance.DeviceID = &deviceID
	createMessageModelInstance.IsGuest = isGuest
	createMessageModelInstance.TemporaryID = &createMessageRequest.TemporaryID
	createMessageModelInstance.UserID = int(requestContext.UserInfo.UserID)
	createMessageModelInstance.Platform = &platformCode
	createMessageRequest.RepliedConversationId = requestContext.OriginalMessage.ID

	if createMessageRequest.OGTags != nil {
		ogTagsMap := createMessageRequest.OGTags.(map[string]interface{})
		ogTagsByte, err := json.Marshal(ogTagsMap)
		if err != nil {
			return constant.APIErrorBadRequest(fmt.Errorf("failed to marshal og_tags"))
		}
		ogTagsString := string(ogTagsByte)
		createMessageModelInstance.OgTags = &ogTagsString
	}

	if createMessageRequest.RepliedChatroomID != "" {
		repliedChatroomID, err := strconv.Atoi(createMessageRequest.RepliedChatroomID)
		if err != nil {
			return constant.APIErrorBadRequest(fmt.Errorf("failed to parse replied_chatroom_id as integer, err=%s", err))
		}
		createMessageModelInstance.ReplyChatroomID = &repliedChatroomID
	}

	if int(createMessageRequest.State) == constant.MessageStateMessagePoll {
		err := FillCreatePollMessageModelInstance(createMessageModelInstance, createMessageRequest)
		if err != nil {
			return constant.APIErrorBadRequest(fmt.Errorf("failed to fill create poll message model instance, err=%s", err))
		}
	}

	if int(createMessageRequest.State) == constant.MessageStateMessageEvent {
		err := FillCreateEventMessageModelInstance(createMessageModelInstance, createMessageRequest)
		if err != nil {
			return constant.APIErrorBadRequest(fmt.Errorf("failed to fill create event message model instance, err=%s", err))
		}
	}

	return nil
}

func FillCreatePollMessageModelInstance(createMessageModelInstance *models.Message, createMessageRequest *requestresponse.CreateMessageRequest) error {

	pollTypeDefault := int(0)

	createMessageModelInstance.State = int(createMessageRequest.State)
	if createMessageModelInstance.PollType = &pollTypeDefault; createMessageRequest.PollType != nil {
		var pollType = new(int)
		*pollType = int(*createMessageRequest.PollType)
		createMessageModelInstance.PollType = pollType
	}

	if createMessageModelInstance.MultipleSelectState = nil; createMessageRequest.MultilpleSelectState != nil {
		var multilpleSelectState = new(int)
		*multilpleSelectState = int(*createMessageRequest.MultilpleSelectState)
		createMessageModelInstance.PollType = multilpleSelectState
	}

	if createMessageModelInstance.MultipleSelectNo = nil; createMessageRequest.MultilpleSelectNo != nil {
		var multilpleSelectNo = new(int)
		*multilpleSelectNo = int(*createMessageRequest.MultilpleSelectNo)
		createMessageModelInstance.PollType = multilpleSelectNo
	}

	createMessageModelInstance.IsAnonymous = createMessageRequest.IsAnonymous
	createMessageModelInstance.AllowAddOption = createMessageRequest.AllowAddOption

	if createMessageModelInstance.ExpiryTime = nil; createMessageRequest.ExpiryTime != nil {
		var expiryTime = new(int)
		*expiryTime = int(*createMessageRequest.ExpiryTime)
		createMessageModelInstance.PollType = expiryTime
	}

	createMessageModelInstance.NoPollExpiry = createMessageRequest.NoPollExpiry
	createMessageModelInstance.AllowVoteChange = createMessageRequest.AllowVoteChange
	createMessageModelInstance.PollAnswerText = constant.MessagePollAnswerText

	return nil
}

func FillCreateEventMessageModelInstance(createMessageModelInstance *models.Message, createMessageRequest *requestresponse.CreateMessageRequest) error {
	createMessageModelInstance.State = int(createMessageRequest.State)
	createMessageModelInstance.Header = createMessageRequest.Header
	createMessageModelInstance.OnlineLink = createMessageRequest.OnlineLink
	createMessageModelInstance.OnlineLinkID = createMessageRequest.OnlineLinkID
	createMessageModelInstance.OnlineLinkPassword = createMessageRequest.OnlineLinkPassword
	createMessageModelInstance.Location = createMessageRequest.Location
	createMessageModelInstance.LocationLat = createMessageRequest.LocationLat
	createMessageModelInstance.LocationLong = createMessageRequest.LocationLong

	if createMessageModelInstance.StartTime = 0; createMessageRequest.StartTime != 0 {
		createMessageModelInstance.StartTime = createMessageRequest.StartTime
	}

	if createMessageModelInstance.EndTime = 0; createMessageRequest.EndTime != 0 {
		createMessageModelInstance.EndTime = createMessageRequest.EndTime
	}

	if createMessageModelInstance.OnlineLinkEnableBefore = utilities.GetMillisecondsInMinutes(constant.OnlineEventLinkEnableBeforeMinutes); createMessageRequest.OnlineLinkEnableBefore != 0 {
		createMessageModelInstance.OnlineLinkEnableBefore = createMessageRequest.OnlineLinkEnableBefore
	}

	if createMessageModelInstance.CoHosts = nil; createMessageRequest.CoHosts != nil {
		createMessageModelInstance.CoHosts = createMessageRequest.CoHosts
	}

	return nil
}

func FillCreateMessageAttachmentsModelInstances(createMessageAttachmentModelInstances *[]models.MessageAttachment, createMessageRequestAttachments []requestresponse.MessageAttachment) *constant.APIError {

	for i := range createMessageRequestAttachments {
		createMessageAttachmentModelInstance := &models.MessageAttachment{}
		createMessageAttachmentModelInstance.Type = createMessageRequestAttachments[i].Type
		createMessageAttachmentModelInstance.FileURL = createMessageRequestAttachments[i].FileURL
		createMessageAttachmentModelInstance.LocationName = &createMessageRequestAttachments[i].LocationName
		createMessageAttachmentModelInstance.LocationLat = createMessageRequestAttachments[i].LocationLat
		createMessageAttachmentModelInstance.LocationLong = createMessageRequestAttachments[i].LocationLong
		createMessageAttachmentModelInstance.Width = &createMessageRequestAttachments[i].Width
		createMessageAttachmentModelInstance.Height = &createMessageRequestAttachments[i].Height
		createMessageAttachmentModelInstance.ThumbnailURL = &createMessageRequestAttachments[i].ThumbnailURL
		createMessageAttachmentModelInstance.CreatedAt = time.Now().Unix()

		attachmentMetaMap := createMessageRequestAttachments[i].Meta.(map[string]interface{})
		attachmentMetaByte, err := json.Marshal(attachmentMetaMap)
		if err != nil {
			return constant.APIErrorBadRequest(fmt.Errorf("failed to marshal attachment meta, err=%s", err))
		}
		attachmentMetaString := string(attachmentMetaByte)

		createMessageAttachmentModelInstance.Meta = &attachmentMetaString

		createMessageAttachmentModelInstance.Name = &createMessageRequestAttachments[i].Name

		*createMessageAttachmentModelInstances = append(*createMessageAttachmentModelInstances, *createMessageAttachmentModelInstance)
	}

	return nil
}

func FillCreateMessagePollModelInstances(createMessagePollsModelInstances *[]models.MessagePoll, createMessageRequestPolls []requestresponse.PollObject, userID int) *constant.APIError {
	for i := range createMessageRequestPolls {
		createMessagePollsModelInstance := &models.MessagePoll{}
		createMessagePollsModelInstance.ConversationID = createMessageRequestPolls[i].ConversationID
		createMessagePollsModelInstance.UserID = userID
		createMessagePollsModelInstance.Text = createMessageRequestPolls[i].Text
		createMessagePollsModelInstance.CreatedAt = time.Now().Unix()
		createMessagePollsModelInstance.UpdatedAt = time.Now().Unix()

		*createMessagePollsModelInstances = append(*createMessagePollsModelInstances, *createMessagePollsModelInstance)
	}

	return nil
}

func CreateMessageInDB(
	createMessageModelInstance *models.Message,
	createMessageAttachmentModelInstances []models.MessageAttachment,
	createMessagePollModelInstances []models.MessagePoll,
) (int64, *constant.APIError) {
	messageID, err := repository.CreateMessage(createMessageModelInstance, createMessageAttachmentModelInstances, createMessagePollModelInstances)
	if err != nil {
		return 0, constant.APIErrorInternalServerError(fmt.Errorf("failed to create message in database, err=%s", err))
	}

	return messageID, nil
}

func CreateMessageErrorResponse(psResponse *requestresponse.PSResponse, createMessageResponse *requestresponse.CreateMessageResponse, apiError *constant.APIError) requestresponse.PSResponse {
	fillCreateMessageErrorResponse(createMessageResponse, apiError)

	createMessageResponseBytes, err := json.Marshal(createMessageResponse)
	if err != nil {
		log.Printf("failed to marshal create message response, err=%s", err)
	}
	psResponse.RawData = string(createMessageResponseBytes)

	return *psResponse
}

func fillCreateMessageErrorResponse(createMessageResponse *requestresponse.CreateMessageResponse, apiError *constant.APIError) {
	createMessageResponse.HTTPStatusCode = apiError.HTTPStatusCode
	createMessageResponse.Success = false
	createMessageResponse.Message = nil
	createMessageResponse.Error = apiError.Error()
}

func CreateMessageSuccessResponse(psResponse *requestresponse.PSResponse, createMessageResponse *requestresponse.CreateMessageResponse, createMessageModelInstance *models.Message) requestresponse.PSResponse {
	fillCreateMessageSuccessResponse(createMessageResponse, createMessageModelInstance)

	createMessageResponseBytes, err := json.Marshal(createMessageResponse)
	if err != nil {
		log.Printf("failed to marshal create message response, err=%s", err)
	}
	psResponse.RawData = string(createMessageResponseBytes)

	return *psResponse
}

func CreateMessageAsnycTasks(
	chatroomId int,
	messageID int64,
	shouldStreamChatbotResponse bool,
	apiVersion int,
) {
	utilities.SafeGo(func() { createMessageAsnycTasksInCaravan() })
	utilities.SafeGo(func() { triggerChatbotInCaravan(chatroomId, messageID, shouldStreamChatbotResponse, apiVersion) })
}

func createMessageAsnycTasksInCaravan() {
	// TODO: send async task to caravan
}

func triggerChatbotInCaravan(chatroomId int, messageID int64, shouldStreamChatbotResponse bool, apiVersion int) {
	requestPostBody := &requestresponse.TriggerChatbot{
		ChatroomID:                  chatroomId,
		MessageID:                   messageID,
		ShouldStreamChatbotResponse: shouldStreamChatbotResponse,
		ApiVersion:                  apiVersion,
	}
	enpoint := external.EndpointTriggerChatbot

	external.NewAPIClientCaravan().Post(enpoint, requestPostBody)
}

func fillCreateMessageSuccessResponse(createMessageResponse *requestresponse.CreateMessageResponse, createMessageModelInstance *models.Message) {
	createMessageResponse.HTTPStatusCode = constant.HTTPResponseCodeOK
	createMessageResponse.Success = true
	createMessageResponse.Message = createMessageModelInstance
	createMessageResponse.Error = ""
}

func validateWsChatroomID(chatroomID int, topic string) error {
	chatroomIDString := strconv.Itoa(chatroomID)

	topicChatroomID, err := utilities.GetTopicSplit(topic)
	if err != nil {
		return fmt.Errorf("failed to split subscribed topic, err=%s", err)
	}

	if chatroomIDString != topicChatroomID[1] {
		return fmt.Errorf("failed to create message in chatroom, non-subscribed topic")
	}

	return nil
}

func validateMessageGroupTags(messageText string, isSecretChatroom bool, chatroomCreatorID int, isMemberAdminInCommunity bool, userID int) bool {
	everyoneTagRegex := regexp.MustCompile(constant.RegexTagEveryone)
	participantsTagRegex := regexp.MustCompile(constant.RegexTagParticipants)

	isEveryoneTag := everyoneTagRegex.FindAllString(messageText, -1)

	if len(isEveryoneTag) > 0 && isSecretChatroom {
		return false
	}

	if len(isEveryoneTag) > 0 && !isMemberAdminInCommunity {
		return false
	}

	isParticipantsTag := participantsTagRegex.FindAllString(messageText, -1)

	if len(isParticipantsTag) > 0 && !isMemberAdminInCommunity && userID != chatroomCreatorID {
		return false
	}

	return true
}
