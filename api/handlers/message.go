package handlers

import (
	"encoding/json"
	"fmt"
	"likeminds-pandemonium/api/constant"
	"likeminds-pandemonium/api/helpers"
	"likeminds-pandemonium/api/models"
	requestresponse "likeminds-pandemonium/api/request_response"
	"likeminds-pandemonium/api/utilities"
	"likeminds-pandemonium/common"
	"log"
)

func CreateMessage(psRequest requestresponse.PSRequest, UUID string, apiKey string, senderDeviceID string, topic string, sdkSource string, platformCode string, versionCode int, apiVersion int) (requestresponse.PSResponse, *requestresponse.CreateMessageResponse) {
	psResponse := &requestresponse.PSResponse{
		DeviceID:         senderDeviceID,
		TopicMessageType: common.TopicMessageTypeCreateConversationResponse,
		RawData:          make([]byte, 0),
	}

	createMessageResponse := &requestresponse.CreateMessageResponse{}
	apiError := &constant.APIError{}

	var psRequestRawDataMap map[string]interface{}
	err := json.Unmarshal(psRequest.RawData, &psRequestRawDataMap)
	if err != nil {
		apiError = constant.APIErrorBadRequest(fmt.Errorf("failed to marshal PSRequest.RawData, err=%s", err))
		return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError), nil
	}

	CreateMessageData, err := json.Marshal(psRequestRawDataMap["data"])
	if err != nil {
		apiError = constant.APIErrorBadRequest(fmt.Errorf("failed to marshal create message json data, err=%s", err))
		return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError), nil
	}

	var createMessageRequest requestresponse.CreateMessageRequest
	if err := json.Unmarshal(CreateMessageData, &createMessageRequest); err != nil {
		apiError = constant.APIErrorBadRequest(fmt.Errorf("failed to unmarshal create message request, err=%s", err))
		return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError), nil
	}

	requestContext, apiError := helpers.ValidateCreateMessageRequest(&createMessageRequest, UUID, apiKey, senderDeviceID, topic)
	if apiError != nil {
		return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError), nil
	}

	apiError = helpers.ValidateCreateMessagePermission(&createMessageRequest, requestContext.Chatroom, requestContext.UserInfo.UserID, requestContext.MemberState)
	if apiError != nil {
		return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError), nil
	}

	collabcardState, err := helpers.GetUserCollabcardStateForChatroom(int(requestContext.Chatroom.ID), requestContext.UserInfo.UserID)
	if err != nil {
		log.Printf("failed to get collabcard state, chatroom=%d, user=%d", requestContext.Chatroom.ID, requestContext.UserInfo.UserID)
	}

	if collabcardState != nil {
		if constant.VersionCheck(sdkSource, platformCode, versionCode, apiVersion, constant.FeatureM2CMEnum) &&
			requestContext.Chatroom.IsPrivate &&
			requestContext.Chatroom.IsPrivateMember &&
			requestContext.Chatroom.Type == constant.ChatroomTypeDirectMessage &&
			*collabcardState.ChatRequestState == constant.DMChatRequestStatesRejected {
			apiError = constant.APIErrorForbidden(fmt.Errorf("failed to create message, user is blocked"))
			return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError), nil
		}
	}

	isMemberVerifiedInCommunity, err := helpers.IsMemberVerifiedInCommunity(int(requestContext.Community.ID), int(requestContext.UserInfo.UserID))
	if err != nil {
		log.Printf("failed to verify member in community, community id=%d, user id=%d, err=%s", requestContext.Community.ID, requestContext.UserInfo.UserID, err)
	}

	var isGuest bool
	if requestContext.Chatroom.AccessWithoutSubscription &&
		collabcardState == nil &&
		!isMemberVerifiedInCommunity {
		isGuest = true
	}

	createMessageModelInstance := &models.Message{}
	apiError = helpers.FillCreateMessageModelInstance(createMessageModelInstance, &createMessageRequest, *requestContext, senderDeviceID, isGuest, platformCode)
	if apiError != nil {
		return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError), nil
	}

	createMessageAttachmentModelInstances := &[]models.MessageAttachment{}
	apiError = helpers.FillCreateMessageAttachmentsModelInstances(createMessageAttachmentModelInstances, createMessageRequest.Attachments)
	if apiError != nil {
		return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError), nil
	}

	createMessagePollModelInstances := &[]models.MessagePoll{}
	apiError = helpers.FillCreateMessagePollModelInstances(createMessagePollModelInstances, createMessageRequest.Polls, requestContext.UserInfo.UserID)
	if apiError != nil {
		return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError), nil
	}

	var swarmCreateWidgetRequest *requestresponse.SwarmCreateWidgetRequest
	if requestContext.CreateWidget {
		swarmCreateWidgetRequest, apiError = helpers.GetSwarmCreateWidgetRequest(UUID, apiKey, int(requestContext.Community.ID), constant.MessageWidgetTypeMessageEnum, createMessageRequest.Metadata)
		if apiError != nil {
			return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError), nil
		}
	}

	messageResponse, apiError := helpers.CreateMessageInDB(createMessageModelInstance, *createMessageAttachmentModelInstances, *createMessagePollModelInstances, swarmCreateWidgetRequest)
	if apiError != nil {
		return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError), nil
	}
	messageResponse.User = requestContext.UserInfo
	messageResponse.RepliedMessage = requestContext.OriginalMessage
	log.Printf("created message, id=%d", messageResponse.Message.ID)

	utilities.SafeGo(func() {
		helpers.CreateMessageCaravanTasks(
			int(requestContext.Chatroom.ID),
			int(messageResponse.Message.ID),
			apiVersion,
			collabcardState.ID,
			requestContext.UserInfo.UserID,
			requestContext.MemberState,
			&createMessageRequest,
		)
	})

	participants := psRequestRawDataMap[common.ParamParticipantsType]
	participantsStringList, ok := participants.([]string)
	if !ok {
		apiError = constant.APIErrorBadRequest(fmt.Errorf("failed to get %s from PSRequest.RawData, err=%s", common.ParamParticipantsType, err))
		return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError), nil
	}

	totalParticipantsCount := psRequestRawDataMap[common.ParamTotalParticipantsCountType]
	totalParticipantsCountInt, ok := totalParticipantsCount.(int)
	if !ok {
		apiError = constant.APIErrorBadRequest(fmt.Errorf("failed to get %s from PSRequest.RawData, err=%s", common.ParamTotalParticipantsCountType, err))
		return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError), nil
	}

	return helpers.CreateMessageSuccessResponse(psResponse, createMessageResponse, messageResponse, participantsStringList, totalParticipantsCountInt), createMessageResponse
}
