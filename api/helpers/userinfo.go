package helpers

import (
	"fmt"
	"likeminds-pandemonium/api/models"
	"likeminds-pandemonium/api/repository"
)

func GetUserInfoByUUID(userID string) (*models.UserInfo, error) {
	userinfo, err := repository.GetUserInfoByUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info, user id=%s, err=%s", userID, err)
	}

	return userinfo, nil
}
