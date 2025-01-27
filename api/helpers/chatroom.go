package helpers

import (
	"fmt"
	"likeminds-pandemonium/api/models"
	"likeminds-pandemonium/api/repository"
)

func GetChatroomByID(chatroomID float64) (*models.Chatroom, error) {
	chatroomIDStr := fmt.Sprintf("%g", chatroomID)
	chatroom, err := repository.GetChatroomByID(chatroomIDStr)
	if err != nil {
		return nil, err
	}

	return chatroom, nil
}
