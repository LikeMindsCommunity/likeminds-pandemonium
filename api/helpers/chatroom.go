package helpers

import (
	"fmt"
	"likeminds-pandemonium/api/models"
	"likeminds-pandemonium/api/repository"
)

func GetChatroomByID(chatroomID int64) (*models.Chatroom, error) {
	chatroom, err := repository.GetChatroomByID(chatroomID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chatroom, id=%d, err=%s", chatroomID, err)
	}

	return chatroom, nil
}
