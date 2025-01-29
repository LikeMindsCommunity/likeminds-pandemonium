package helpers

import (
	"likeminds-pandemonium/api/repository"
)

func ValidateUserRightInChatroom(chatroomID int, userID int, userRight string, isMemberAdminInCommunity bool) bool {
	return isMemberAdminInCommunity || repository.UserHasRightInChatroom(chatroomID, userID, userRight)
}
