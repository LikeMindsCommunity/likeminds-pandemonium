package repository

import (
	"likeminds-pandemonium/database"

	"gorm.io/gorm"
)

type MessageAttachmentRepository struct {
	messageAttachmentDatabase *gorm.DB
}

func NewMessageAttachmentRepository() *MessageAttachmentRepository {
	return &MessageAttachmentRepository{
		messageAttachmentDatabase: database.Postgres,
	}
}
