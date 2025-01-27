package handlers

import (
	"encoding/json"
	"likeminds-pandemonium/api/helpers"
	"likeminds-pandemonium/api/models"
	requestresponse "likeminds-pandemonium/api/request_response"
	"likeminds-pandemonium/common"
	"log"
)

func CreateMessage(messageData map[string]interface{}, userID string, deviceID string, topic string) requestresponse.PSResponse {

	CreateMessageResponse := requestresponse.PSResponse{
		DeviceID:         deviceID,
		TopicMessageType: common.TopicMessageTypeCreateConversationResponse,
		RawData:          "",
	}

	CreateMessageData, err := json.Marshal(messageData["data"])
	if err != nil {
		log.Printf(common.ErrorMarshalErrorJson, err)
		return CreateMessageResponse
	}

	var createMessageRequest requestresponse.CreateMessageRequest
	if err := json.Unmarshal(CreateMessageData, &createMessageRequest); err != nil {
		log.Printf(common.ErrorUnmarshalErrorJson, err)
		return CreateMessageResponse
	}

	requestContext, err := helpers.ValidateCreateMessageRequest(createMessageRequest, userID, deviceID, topic)
	if err != nil {
		log.Print(err)
		return CreateMessageResponse
	}
	log.Print(requestContext)

	err = helpers.ValidateCreateMessagePermission(&createMessageRequest, &requestContext.Chatroom, requestContext.UserInfo.UserID)
	if err != nil {
		log.Print(err)
		return CreateMessageResponse
	}

	collabcardState, err := helpers.GetUserCollabcardStateForChatroom(int(requestContext.Chatroom.ID), requestContext.UserInfo.UserID)
	if err != nil {
		log.Print(err)
		return CreateMessageResponse
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
	err = helpers.FillCreateMessageModelInstance(createMessageModelInstance, &createMessageRequest, *requestContext, deviceID, isGuest)
	if err != nil {
		log.Print(err)
		return CreateMessageResponse
	}

	messageID, err := helpers.CreateMessageInDB(createMessageModelInstance)
	if err != nil {
		log.Print(err)
		return CreateMessageResponse
	}
	log.Printf("created_message id=%d", messageID)

	dataResponse := &requestresponse.CreateMessageResponse{}
	helpers.FillDataResponse(dataResponse, createMessageModelInstance)

	dataResponseBytes, err := json.Marshal(*dataResponse)
	if err != nil {
		log.Printf(common.ErrorMarshalErrorJson, err)
		return CreateMessageResponse
	}
	CreateMessageResponse.RawData = string(dataResponseBytes)

	return CreateMessageResponse
}
