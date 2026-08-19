package repository

import (
	"fmt"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/api/models"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/database"

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

func GetCommunityByID(ID int64) (*models.Community, error) {
	community := &models.Community{}

	filter := map[string]interface{}{
		"id": ID,
	}

	err := NewCommunityRepository().communityDatabase.Where(filter).First(&community).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get community, id=%d, err=%s", ID, err)
	}

	return community, nil
}
