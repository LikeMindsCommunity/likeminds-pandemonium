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

func UserHasRightInChatroom(chatroomID int, userID int, userRight string) bool {
	userHasRightInChatroom := &models.UserChannelSetting{}
	filter := map[string]interface{}{
		"chatroom_id":  chatroomID,
		"user_id":      userID,
		"setting_type": userRight,
	}

	userRightInChatroom := NewUserChannelSettingRepository().userChannelSettingDatabase.Where(filter).First(&userHasRightInChatroom)
	if userRightInChatroom.Error != nil {
		log.Print("error fetching user channel settings")
		return false
	}

	return userHasRightInChatroom.Enabled
}
