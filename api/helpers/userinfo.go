package helpers

import (
	"likeminds-pandemonium/api/models"
	"likeminds-pandemonium/api/repository"
)

func GetUserInfoByUUID(userID string) (*models.UserInfo, error) {
	userinfo, err := repository.GetUserInfoByUUID(userID)
	if err != nil {
		return nil, err
	}

	return userinfo, nil
}
