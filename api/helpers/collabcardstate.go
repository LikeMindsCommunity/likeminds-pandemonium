package helpers

import (
	"likeminds-pandemonium/api/models"
	"likeminds-pandemonium/api/repository"
)

func GetUserCollabcardStateForChatroom(chatroomID int, userID int) (*models.CollabCardState, error) {
	collabcardState, err := repository.GetUserCollabcardStateForChatroom(chatroomID, userID)
	if err != nil {
		return nil, err
	}

	return collabcardState, nil
}
