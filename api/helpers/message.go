package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"likeminds-pandemonium/api/constant"
	"likeminds-pandemonium/api/models"
	"likeminds-pandemonium/api/repository"
	requestresponse "likeminds-pandemonium/api/request_response"
	"likeminds-pandemonium/api/utilities"
	"likeminds-pandemonium/common"
	"log"
	"regexp"
	"slices"
	"time"
)

type CreateMessageRequestContext struct {
	UserInfo        models.UserInfo  `json:"userinfo"`
	Chatroom        models.Chatroom  `json:"chatroom"`
	Community       models.Community `json:"community"`
	OriginalMessage models.Message   `json:"original_message"`
}

func ValidateCreateMessageRequest(createMessageRequest requestresponse.CreateMessageRequest, userID string, deviceID string, topic string) (*CreateMessageRequestContext, error) {

	createMessageRequestContext := &CreateMessageRequestContext{}

	err := validateWsChatroomID(createMessageRequest.ChatroomID.(float64), topic)
	if err != nil {
		return nil, err
	}

	userinfo, err := GetUserInfoByUUID(userID)
	if err != nil {
		return nil, err
	}
	createMessageRequestContext.UserInfo = *userinfo

	chatroom, err := GetChatroomByID(createMessageRequest.ChatroomID.(float64))
	if err != nil {
		return nil, err
	}
	createMessageRequestContext.Chatroom = *chatroom

	if chatroom.IsSecret && !ValidateUserSecretChatroomAccess([]byte(*chatroom.SecretChatroomParticipants), userinfo.UserID) {
		return nil, fmt.Errorf("chatroom not found")
	}

	community, err := GetCommunityByID(chatroom.CommunityID)
	if err != nil {
		return nil, err
	}
	createMessageRequestContext.Community = *community

	if createMessageRequest.RepliedConversationId != "" {
		originalMessage, err := GetMessageByID(createMessageRequest.RepliedConversationId.(float64))
		if err != nil {
			return nil, err
		}
		createMessageRequestContext.OriginalMessage = *originalMessage
	}

	return createMessageRequestContext, nil
}

func GetMessageByID(messageID float64) (*models.Message, error) {
	messageIDStr := fmt.Sprintf("%g", messageID)
	message, err := repository.GetMessageByID(messageIDStr)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func ValidateCreateMessagePermission(createMessageRequest *requestresponse.CreateMessageRequest, chatroom *models.Chatroom, userIDInt int) error {
	if chatroom.Type == constant.ChatroomTypeMasterIntro {
		return errors.New("cannot post message in master intro chatroom")
	}

	memberState, err := repository.GetMemberStateInCommunity(chatroom.CommunityID, userIDInt)
	if err != nil {
		return errors.New("cannot get member state in community")
	}

	isMemberAdminInCommunity := IsMemberStateAdminInCommunity(memberState)
	if !validateMessageGroupTags(createMessageRequest.Text, chatroom.IsSecret, chatroom.UserID, isMemberAdminInCommunity, userIDInt) {
		return errors.New("invalid message group tags")
	}

	if !ValidateUserRight(chatroom.CommunityID, userIDInt, constant.MemberRightRespondInRoom, int(chatroom.ID)) {
		return fmt.Errorf("user right missing, right=%s", constant.MemberRightRespondInRoomEnum)
	}

	return nil
}

func ValidateUserSecretChatroomAccess(secretChatroomParticipants []byte, userID int) bool {
	var secretChatroomParticipantsList []int
	err := json.Unmarshal(secretChatroomParticipants, &secretChatroomParticipantsList)
	if err != nil {
		log.Printf(common.ErrorUnmarshalErrorJson, err)
		return false
	}

	if !slices.Contains(secretChatroomParticipantsList, userID) {
		return false
	}

	return true
}

func FillCreateMessageModelInstance(
	createMessageModelInstance *models.Message,
	createMessageRequest *requestresponse.CreateMessageRequest,
	requestContext CreateMessageRequestContext,
	deviceID string,
	isGuest bool) error {

	createMessageModelInstance.Answer = createMessageRequest.Text
	createMessageModelInstance.APIVersion = 1
	createMessageModelInstance.AttachmentCount = int(createMessageRequest.AttachmentCount)
	createMessageModelInstance.AttachmentsUploaded = false
	createMessageModelInstance.CardID = int(requestContext.Chatroom.ID)
	createMessageModelInstance.CommunityID = int(requestContext.Community.ID)
	createMessageModelInstance.CreatedAt = time.Now().Unix()
	createMessageModelInstance.DeviceID = &deviceID
	createMessageModelInstance.HasFiles = createMessageRequest.HasFiles
	createMessageModelInstance.IsGuest = isGuest
	createMessageModelInstance.TemporaryID = &createMessageRequest.TemporaryID
	createMessageModelInstance.UserID = int(requestContext.UserInfo.UserID)
	// TODO: replied chatroom_id
	// TODO: replied platform_code

	if int(createMessageRequest.State) == constant.ConversationStateConversationPoll {
		err := FillCreatePollMessageModelInstance(createMessageModelInstance, createMessageRequest)
		if err != nil {
			return err
		}
	}

	if int(createMessageRequest.State) == constant.ConversationStateConversationEvent {
		err := FillCreateEventMessageModelInstance(createMessageModelInstance, createMessageRequest)
		if err != nil {
			return err
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
	createMessageModelInstance.PollAnswerText = constant.ConversationPollAnswerText

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

func CreateMessageInDB(createMessageModelInstance *models.Message) (int, error) {
	messageID, err := repository.CreateMessage(createMessageModelInstance)
	if err != nil {
		return 0, err
	}

	return messageID, nil
}

func FillDataResponse(dataResponse *requestresponse.CreateMessageResponse, createMessageModelInstance *models.Message) {
	dataResponse.Message = createMessageModelInstance
}

func validateWsChatroomID(chatroomID float64, topic string) error {
	topicChatroomID, err := utilities.GetTopicSplit(topic)
	if err != nil {
		return err
	}
	if fmt.Sprintf("%g", chatroomID) != topicChatroomID[1] {
		return errors.New(common.ErrorChatroomIdMismatch)
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
