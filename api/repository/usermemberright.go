package repository

import (
	"fmt"
	"likeminds-pandemonium/api/models"
	"likeminds-pandemonium/database"

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

func UserHasRightInCommunity(communityID int, userID int, rightID int) (bool, error) {
	userHasRightInCommunity := &models.UserMemberRight{}
	filter := map[string]interface{}{
		"community_id": communityID,
		"user_id":      userID,
		"right_id":     rightID,
	}

	err := NewUserMemberRightRepository().userMemberRightDatabase.Where(filter).First(&userHasRightInCommunity).Error
	if err != nil {
		return false, fmt.Errorf("failed to validate community user right, community id=%d, user id=%d, right id=%d, err=%s", communityID, userID, rightID, err)
	}

	return true, nil
}
