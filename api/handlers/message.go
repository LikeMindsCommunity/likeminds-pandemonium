package handlers

import (
	"encoding/json"
	"errors"
	"likeminds-pandemonium/api/constant"
	"likeminds-pandemonium/api/helpers"
	"likeminds-pandemonium/api/models"
	requestresponse "likeminds-pandemonium/api/request_response"
	"likeminds-pandemonium/api/utilities"
	"likeminds-pandemonium/common"
	"log"
	"strconv"
)

func CreateMessage(messageData map[string]interface{}, userID string, deviceID string, topic string, platformCode string, versionCode string, apiVersion string) requestresponse.PSResponse {

	psResponse := &requestresponse.PSResponse{
		DeviceID:         deviceID,
		TopicMessageType: common.TopicMessageTypeCreateConversationResponse,
		RawData:          "",
	}

	createMessageResponse := &requestresponse.CreateMessageResponse{}

	CreateMessageData, err := json.Marshal(messageData["data"])
	if err != nil {
		var apiError = constant.APIErrorBadRequest(errors.New(common.ErrorMarshalErrorJson))
		return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError)
	}

	var createMessageRequest requestresponse.CreateMessageRequest
	if err := json.Unmarshal(CreateMessageData, &createMessageRequest); err != nil {
		var apiError = constant.APIErrorBadRequest(errors.New(common.ErrorMarshalErrorJson))
		return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError)
	}

	requestContext, apiError := helpers.ValidateCreateMessageRequest(createMessageRequest, userID, deviceID, topic)
	if apiError != nil {
		return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError)
	}
	log.Print(requestContext)

	apiError = helpers.ValidateCreateMessagePermission(&createMessageRequest, &requestContext.Chatroom, requestContext.UserInfo.UserID)
	if apiError != nil {
		return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError)
	}

	collabcardState, err := helpers.GetUserCollabcardStateForChatroom(int(requestContext.Chatroom.ID), requestContext.UserInfo.UserID)
	if err != nil {
		log.Print(err)
	}

	if collabcardState.ID != 0 {
		// TODO: add m2cm v2 check
	}

	var isGuest bool
	if requestContext.Chatroom.AccessWithoutSubscription &&
		collabcardState.ID == 0 &&
		!helpers.IsMemberVerifiedInCommunity(int(requestContext.Community.ID), int(requestContext.UserInfo.ID)) {
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

	messageID, apiError := helpers.CreateMessageInDB(createMessageModelInstance, *createMessageAttachmentModelInstances)
	if apiError != nil {
		return helpers.CreateMessageErrorResponse(psResponse, createMessageResponse, apiError)
	}
	log.Printf("created_message id=%d", messageID)

	apiVersionInt, err := strconv.Atoi(apiVersion)
	if err != nil {
		log.Print(err)
	}

	utilities.SafeGo(func() {
		helpers.CreateMessageAsnycTasks(
			int(requestContext.Chatroom.ID),
			messageID,
			createMessageRequest.ShouldStreamChatbotResponse,
			apiVersionInt,
		)
	})

	return helpers.CreateMessageSuccessResponse(psResponse, createMessageResponse, createMessageModelInstance)
}
