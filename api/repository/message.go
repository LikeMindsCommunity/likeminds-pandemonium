package repository

import (
	"errors"
	"likeminds-pandemonium/api/models"
	"likeminds-pandemonium/database"

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
		return nil, err
	}

	return message, nil
}

func CreateMessage(
	createMessageModelInstance *models.Message,
	createMessageAttachmentModelInstances []models.MessageAttachment,
	createMessagePollModelInstances []models.MessagePoll,
) (int, error) {
	messageID := int64(0)

	tx := NewMessageRepository().messageDatabase.Begin()

	message := tx.Create(createMessageModelInstance)
	if message.Error != nil {
		tx.Rollback()
		return 0, errors.New("failed to create message in database, rollingback")
	}
	messageID = createMessageModelInstance.ID

	if len(createMessageAttachmentModelInstances) > 0 {
		createMessageAttachmentModelInstances, err := updateAttachmentsWithMessageID(createMessageAttachmentModelInstances, int(messageID))
		if err != nil {
			tx.Rollback()
			return 0, errors.New("failed to attach messageID to attachments, rollingback")
		}
		attachments := tx.Create(createMessageAttachmentModelInstances)
		if attachments.Error != nil {
			tx.Rollback()
			return 0, errors.New("failed to create message attachments in database, rollingback")
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
			return 0, errors.New("failed to attach messageID to polls, rollingback")
		}
		polls := tx.Create(createMessagePollModelInstances)
		if polls.Error != nil {
			tx.Rollback()
			return 0, errors.New("failed to create message polls in database, rollingback")
		}
		if int(polls.RowsAffected) != len(createMessagePollModelInstances) {
			tx.Rollback()
			return 0, errors.New("failed to create all message attachments in database, rollingback")
		}
	}

	tx.Commit()

	return int(messageID), nil
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
