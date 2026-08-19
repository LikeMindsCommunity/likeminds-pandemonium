package helpers

import (
	"fmt"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/api/repository"
)

func ValidateUserRightInCommunity(communityID int, userID int, rightID int) (bool, error) {
	userHasRightInCommunity, err := repository.UserHasRightInCommunity(communityID, userID, rightID)
	if err != nil {
		return false, fmt.Errorf("failed to validate community user right, community id=%d, user id=%d, right id=%d, err=%s", communityID, userID, rightID, err)
	}
	return userHasRightInCommunity, nil
}
