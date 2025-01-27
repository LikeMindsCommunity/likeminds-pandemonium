package helpers

import (
	"likeminds-pandemonium/api/models"
	"likeminds-pandemonium/api/repository"
	"strconv"
)

func GetChatroomByID(chatroomID int) (*models.Chatroom, error) {
	chatroomIDStr := strconv.Itoa(chatroomID)
	chatroom, err := repository.GetChatroomByID(chatroomIDStr)
	if err != nil {
		return nil, err
	}

	return chatroom, nil
}
