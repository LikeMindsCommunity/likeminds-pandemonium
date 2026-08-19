package repository

import (
	"fmt"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/api/models"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/database"

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
		return nil, fmt.Errorf("failed to get collabcard state, chatroom id=%d, user id=%d, err=%s", chatroomID, userID, err)
	}

	return collabcardState, nil
}
