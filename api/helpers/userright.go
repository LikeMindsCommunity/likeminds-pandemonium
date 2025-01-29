package helpers

import (
	"likeminds-pandemonium/api/repository"
)

func ValidateUserRightInCommunity(communityID int, userID int, rightID int) bool {
	return repository.UserHasRightInCommunity(communityID, userID, rightID)
}
