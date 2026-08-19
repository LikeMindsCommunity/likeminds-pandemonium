package repository

import (
	"fmt"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/api/models"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/database"

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

func UserHasRightInChatroom(chatroomID int, userID int, userRight string) (bool, error) {
	userHasRightInChatroom := &models.UserChannelSetting{}
	filter := map[string]interface{}{
		"chatroom_id":  chatroomID,
		"user_id":      userID,
		"setting_type": userRight,
	}

	userRightInChatroom := NewUserChannelSettingRepository().userChannelSettingDatabase.Where(filter).First(&userHasRightInChatroom)
	if userRightInChatroom.Error != nil {
		return false, fmt.Errorf("failed to validate chatroom user right, chatroom id=%d, userId=%d, user right=%s, err=%s", chatroomID, userID, userRight, userRightInChatroom.Error)
	}

	return userHasRightInChatroom.Enabled, nil
}
