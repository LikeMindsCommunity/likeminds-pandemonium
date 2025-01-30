package repository

import (
	"errors"
	"fmt"
	"likeminds-pandemonium/api/models"
	requestresponse "likeminds-pandemonium/api/request_response"
	"likeminds-pandemonium/database"
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
		// TODO:
		// 1. update swarm widget request with above message id
		// 2. send request to swarm
		// 3. err -> rollback transaction
		// 4. success -> retrieve widget id and update message with widget id
	}

	tx.Rollback()

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
