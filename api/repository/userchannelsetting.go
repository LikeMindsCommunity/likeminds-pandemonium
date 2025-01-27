package repository

import (
	"likeminds-pandemonium/api/models"
	"likeminds-pandemonium/database"
	"log"

	"gorm.io/gorm"
)

type UserChannelSettingRepository struct {
	userChannelSettingDatabase *gorm.DB
}

func NewUserChannelSettingRepository() *UserChannelSettingRepository {
	return &UserChannelSettingRepository{
		userChannelSettingDatabase: database.Postgres,
	}
}

func UserHasRightInChatroom(chatroomID int, userID int) bool {
	userHasRightInChatroom := &models.UserChannelSetting{}
	filter := map[string]interface{}{
		"chatroom_id": chatroomID,
		"user_id":     userID,
	}

	err := NewUserChannelSettingRepository().userChannelSettingDatabase.Where(filter).First(&userHasRightInChatroom).Error
	if err != nil {
		log.Print("error fetching user channel settings")
		return false
	}

	// if row id is default value
	if userHasRightInChatroom.ID == 0 {
		log.Print("user channel settings response error")
		return false
	}

	return true
}
