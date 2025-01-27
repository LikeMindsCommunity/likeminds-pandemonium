package repository

import (
	"likeminds-pandemonium/api/models"
	"likeminds-pandemonium/database"

	"gorm.io/gorm"
)

type CommunityRepository struct {
	communityDatabase *gorm.DB
}

func NewCommunityRepository() *CommunityRepository {
	return &CommunityRepository{
		communityDatabase: database.Postgres,
	}
}

func GetCommunityByID(ID int) (*models.Community, error) {
	community := &models.Community{}

	filter := map[string]interface{}{
		"id": ID,
	}

	err := NewCommunityRepository().communityDatabase.Where(filter).First(&community).Error
	if err != nil {
		return nil, err
	}

	return community, nil
}
