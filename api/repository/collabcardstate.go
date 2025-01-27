package repository

import (
	"likeminds-pandemonium/api/models"
	"likeminds-pandemonium/database"

	"gorm.io/gorm"
)

type CollabcardStateRepository struct {
	collabcardStateDatabase *gorm.DB
}

func NewCollabcardStateRepository() *CollabcardStateRepository {
	return &CollabcardStateRepository{
		collabcardStateDatabase: database.Postgres,
	}
}

func GetUserCollabcardStateForChatroom(chatroomID int, userID int) (*models.CollabCardState, error) {
	collabcardState := &models.CollabCardState{}

	filter := map[string]interface{}{
		"card_id": chatroomID,
		"user_id": userID,
	}

	err := NewCollabcardStateRepository().collabcardStateDatabase.Where(filter).First(&collabcardState).Error
	if err != nil {
		return nil, err
	}

	return collabcardState, nil
}
