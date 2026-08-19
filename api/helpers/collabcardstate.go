package helpers

import (
	"fmt"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/api/models"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/api/repository"
)

func GetUserCollabcardStateForChatroom(chatroomID int, userID int) (*models.CollabCardState, error) {
	collabcardState, err := repository.GetUserCollabcardStateForChatroom(chatroomID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get collabcard state, chatroom id=%d, user id=%d, err=%s", chatroomID, userID, err)
	}

	return collabcardState, nil
}
