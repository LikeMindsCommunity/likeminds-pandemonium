package repository

import (
	"likeminds-pandemonium/api/models"
	"likeminds-pandemonium/database"
	"log"

	"gorm.io/gorm"
)

type UserMemberRightRepository struct {
	userMemberRightDatabase *gorm.DB
}

func NewUserMemberRightRepository() *UserMemberRightRepository {
	return &UserMemberRightRepository{
		userMemberRightDatabase: database.Postgres,
	}
}

func UserHasRightInCommunity(communityID int, userID int, rightID int) bool {
	userHasRightInCommunity := &models.UserMemberRight{}
	filter := map[string]interface{}{
		"community_id": communityID,
		"user_id":      userID,
		"right_id":     rightID,
	}

	err := NewUserMemberRightRepository().userMemberRightDatabase.Where(filter).First(&userHasRightInCommunity).Error
	if err != nil {
		log.Print("error fetching user member rights")
		return false
	}

	// if row id is default value
	if userHasRightInCommunity.ID == 0 {
		log.Print("user member rights response error")
		return false
	}

	return true
}
