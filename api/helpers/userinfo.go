package helpers

import (
	"fmt"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/api/models"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/api/repository"
)

func GetUserInfoByUUID(UUID string) (*models.UserInfo, error) {
	userinfo, err := repository.GetUserInfoByUUID(UUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info, user id=%s, err=%s", UUID, err)
	}

	return userinfo, nil
}
