package helpers

import (
	"fmt"
	"likeminds-pandemonium/api/repository"
)

func ValidateUserRightInChatroom(chatroomID int, userID int, userRight string, isMemberAdminInCommunity bool) (bool, error) {
	if isMemberAdminInCommunity {
		return isMemberAdminInCommunity, nil
	}

	userHasRightInChatroom, err := repository.UserHasRightInChatroom(chatroomID, userID, userRight)
	if err != nil {
		return false, fmt.Errorf("failed to validate chatroom user right, chatroom id=%d, userId=%d, user right=%s, err=%s", chatroomID, userID, userRight, err)
	}

	return userHasRightInChatroom, nil
}
