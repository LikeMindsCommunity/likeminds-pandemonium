package repository

import (
	"fmt"
	"likeminds-pandemonium/api/models"
	"likeminds-pandemonium/database"

	"gorm.io/gorm"
)

type ChatroomRepository struct {
	chatroomDatabase *gorm.DB
}

func NewChatroomRepository() *ChatroomRepository {
	return &ChatroomRepository{
		chatroomDatabase: database.Postgres,
	}
}

func GetChatroomByID(ID int64) (*models.Chatroom, error) {
	chatroom := &models.Chatroom{}

	filter := map[string]interface{}{
		"id": ID,
	}

	err := NewChatroomRepository().chatroomDatabase.Where(filter).First(&chatroom).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get chatroom, id=%d, err=%s", ID, err)
	}

	return chatroom, nil
}
