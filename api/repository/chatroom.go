package repository

import (
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

func GetChatroomByID(ID string) (*models.Chatroom, error) {
	chatroom := &models.Chatroom{}

	filter := map[string]interface{}{
		"id": ID,
	}

	err := NewChatroomRepository().chatroomDatabase.Where(filter).First(&chatroom).Error
	if err != nil {
		return nil, err
	}

	return chatroom, nil
}
