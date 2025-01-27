package helpers

import (
	"likeminds-pandemonium/api/repository"
)

func ValidateUserRight(communityID int, userID int, rightID int, chatroomID int) bool {
	return repository.UserHasRightInCommunity(communityID, userID, rightID) && repository.UserHasRightInChatroom(chatroomID, userID)
}
