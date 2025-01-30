package repository

import (
	"errors"
	"fmt"
	"likeminds-pandemonium/api/models"
	requestresponse "likeminds-pandemonium/api/request_response"
	"likeminds-pandemonium/database"
	"likeminds-pandemonium/external"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type MessageRepository struct {
	messageDatabase *gorm.DB
}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{
		messageDatabase: database.Postgres,
	}
}

func GetMessageByID(ID string) (*models.Message, error) {
	message := &models.Message{}

	filter := map[string]interface{}{
		"id": ID,
	}

	err := NewMessageRepository().messageDatabase.Where(filter).First(&message).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get message, message id=%s, err=%s", ID, err)
	}

	return message, nil
}

func CreateMessage(
	createMessageModelInstance *models.Message,
	createMessageAttachmentModelInstances []models.MessageAttachment,
	createMessagePollModelInstances []models.MessagePoll,
	swarmCreateWidgetRequest *requestresponse.SwarmCreateWidgetRequest,
) (int64, error) {

	tx := NewMessageRepository().messageDatabase.Begin()

	message := tx.Create(createMessageModelInstance)
	if message.Error != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to create message in database, rollingback, err=%s", message.Error)
	}
	messageID := createMessageModelInstance.ID

	if swarmCreateWidgetRequest != nil {
		swarmCreateWidgetRequest.RequestBody.ParentEntityID = strconv.Itoa(int(messageID))
		swarmResponse, err := external.NewAPIClientSwarm().Post(external.EnpointSwarmCreateWidget, swarmCreateWidgetRequest.Headers, swarmCreateWidgetRequest.RequestBody)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to create widget in swarm, rollingback, err=%s", err)
		}
		widgetID := swarmResponse.Data["widget"].(map[string]interface{})["_id"]
		message := tx.Model(&models.Message{}).Where("id = ?", messageID).Update("widget_id", widgetID)
		if message.Error != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to update message widget in database, rollingback, err=%s", message.Error)
		}
		if int(message.RowsAffected) != 1 {
			tx.Rollback()
			return 0, errors.New("failed to update message widget in database, rollingback")
		}
	}

	if len(createMessageAttachmentModelInstances) > 0 {
		createMessageAttachmentModelInstances, err := updateAttachmentsWithMessageID(createMessageAttachmentModelInstances, int(messageID))
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to attach messageID to attachments, rollingback, err=%s", err)
		}
		attachments := tx.Create(createMessageAttachmentModelInstances)
		if attachments.Error != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to create message attachments in database, rollingback, err=%s", attachments.Error)
		}
		if int(attachments.RowsAffected) != len(createMessageAttachmentModelInstances) {
			tx.Rollback()
			return 0, errors.New("failed to create all message attachments in database, rollingback")
		}
	}

	if len(createMessagePollModelInstances) > 0 {
		createMessagePollModelInstances, err := updatePollsWithMessageID(createMessagePollModelInstances, int(messageID))
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to attach messageID to polls, rollingback, err=%s", err)
		}
		polls := tx.Create(createMessagePollModelInstances)
		if polls.Error != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to create message polls in database, rollingback, err=%s", polls.Error)
		}
		if int(polls.RowsAffected) != len(createMessagePollModelInstances) {
			tx.Rollback()
			return 0, errors.New("failed to create all message polls in database, rollingback")
		}
	}

	chatroom := tx.Model(&models.Chatroom{}).Where("id = ?", createMessageModelInstance.CardID).Update("updated_at", time.Now().Unix())
	if chatroom.Error != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to update chatroom in database, rollingback, err=%s", chatroom.Error)
	}
	if int(chatroom.RowsAffected) != 1 {
		tx.Rollback()
		return 0, errors.New("failed to update chatroom in database, rollingback")
	}

	tx.Commit()

	return messageID, nil
}

func updateAttachmentsWithMessageID(createMessageAttachmentModelInstances []models.MessageAttachment, messageID int) ([]models.MessageAttachment, error) {
	for i := range createMessageAttachmentModelInstances {
		createMessageAttachmentModelInstances[i].AnswerID = messageID
	}
	return createMessageAttachmentModelInstances, nil
}

func updatePollsWithMessageID(createMessagePollModelInstances []models.MessagePoll, messageID int) ([]models.MessagePoll, error) {
	for i := range createMessagePollModelInstances {
		createMessagePollModelInstances[i].ConversationID = messageID
	}
	return createMessagePollModelInstances, nil
}
