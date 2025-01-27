package repository

import (
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

func CreateMessage(createMessageModelInstance *models.Message) (int, error) {
	messageID := int64(0)
	err := NewMessageRepository().messageDatabase.Transaction(func(tx *gorm.DB) error {
		message := tx.Create(createMessageModelInstance)
		if message.Error != nil {
			return message.Error
		}
		messageID = createMessageModelInstance.ID

		// return nil will commit the whole transaction
		return nil
	})
	if err != nil {
		return int(messageID), err
	}

	return int(messageID), nil
}
