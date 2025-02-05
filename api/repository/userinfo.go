package repository

import (
	"fmt"
	"likeminds-pandemonium/api/models"
	"likeminds-pandemonium/database"

	"gorm.io/gorm"
)

type UserInfoRepository struct {
	userInfoDatabase *gorm.DB
}

func NewUserInfoRepository() *UserInfoRepository {
	return &UserInfoRepository{
		userInfoDatabase: database.Postgres,
	}
}

func GetUserInfoByUUID(UUID string) (*models.UserInfo, error) {
	userInfo := &models.UserInfo{}

	filter := map[string]interface{}{
		"user_unique_id": UUID,
	}

	err := NewUserInfoRepository().userInfoDatabase.Where(filter).First(&userInfo).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user info, user id=%s, err=%s", UUID, err)
	}

	return userInfo, nil
}
