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

func CreateMessage(messageData map[string]interface{}, UUID string, apiKey string, deviceID string, topic string, sdkSource string, platformCode string, versionCode int, apiVersion int) requestresponse.PSResponse {

	psResponse := &requestresponse.PSResponse{
		DeviceID:         deviceID,
		TopicMessageType: common.TopicMessageTypeCreateConversationResponse,
		RawData:          "",
	}

	createMessageResponse := &requestresponse.CreateMessageResponse{}
	apiError := &constant.APIError{}

	CreateMessageData, err := json.Marshal(messageData["data"])
	if err != nil {
		apiError = constant.APIErrorBadRequest(fmt.Errorf("failed to marshal create message json data, err=%s", err))
		return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError)
	}

	var createMessageRequest requestresponse.CreateMessageRequest
	if err := json.Unmarshal(CreateMessageData, &createMessageRequest); err != nil {
		apiError = constant.APIErrorBadRequest(fmt.Errorf("failed to unmarshal create message request, err=%s", err))
		return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError)
	}

	requestContext, apiError := helpers.ValidateCreateMessageRequest(&createMessageRequest, UUID, apiKey, deviceID, topic)
	if apiError != nil {
		return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError)
	}

	apiError = helpers.ValidateCreateMessagePermission(&createMessageRequest, &requestContext.Chatroom, requestContext.UserInfo.UserID, requestContext.MemberState)
	if apiError != nil {
		return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError)
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
			return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError)
		}
	}

	isMemberVerifiedInCommunity, err := helpers.IsMemberVerifiedInCommunity(int(requestContext.Community.ID), int(requestContext.UserInfo.UserID))
	if err != nil {
		log.Printf("failed to verify member in community, community id=%d, user id=%d, err=%s", requestContext.Community.ID, requestContext.UserInfo.UserID, err)
	}

	var isGuest bool
	if requestContext.Chatroom.AccessWithoutSubscription &&
		collabcardState.ID == 0 &&
		!isMemberVerifiedInCommunity {
		isGuest = true
	}

	createMessageModelInstance := &models.Message{}
	apiError = helpers.FillCreateMessageModelInstance(createMessageModelInstance, &createMessageRequest, *requestContext, deviceID, isGuest, platformCode)
	if apiError != nil {
		return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError)
	}

	createMessageAttachmentModelInstances := &[]models.MessageAttachment{}
	apiError = helpers.FillCreateMessageAttachmentsModelInstances(createMessageAttachmentModelInstances, createMessageRequest.Attachments)
	if apiError != nil {
		return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError)
	}

	createMessagePollModelInstances := &[]models.MessagePoll{}
	apiError = helpers.FillCreateMessagePollModelInstances(createMessagePollModelInstances, createMessageRequest.Polls, requestContext.UserInfo.UserID)
	if apiError != nil {
		return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError)
	}

	var swarmCreateWidgetRequest *requestresponse.SwarmCreateWidgetRequest
	if requestContext.CreateWidget {
		swarmCreateWidgetRequest, apiError = helpers.GetSwarmCreateWidgetRequest(UUID, apiKey, int(requestContext.Community.ID), constant.MessageWidgetTypeMessageEnum, createMessageRequest.Metadata)
		if apiError != nil {
			return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError)
		}
	}

	messageID, apiError := helpers.CreateMessageInDB(createMessageModelInstance, *createMessageAttachmentModelInstances, *createMessagePollModelInstances, swarmCreateWidgetRequest)
	if apiError != nil {
		return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError)
	}
	log.Printf("created message, id=%d", messageID)

	utilities.SafeGo(func() {
		helpers.CreateMessageCaravanTasks(
			int(requestContext.Chatroom.ID),
			int(messageID),
			apiVersion,
			collabcardState.ID,
			requestContext.UserInfo.UserID,
			requestContext.MemberState,
			&createMessageRequest,
		)
	})

	return helpers.CreateMessageSuccessResponse(psResponse, createMessageResponse, createMessageModelInstance)
}
