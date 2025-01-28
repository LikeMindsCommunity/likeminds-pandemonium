package repository

import (
	"likeminds-pandemonium/api/models"
	"likeminds-pandemonium/database"

	"gorm.io/gorm"
)

type CommunityConfigurationRepository struct {
	communityConfigurationDatabase *gorm.DB
}

func NewCommunityConfigurationRepository() *CommunityConfigurationRepository {
	return &CommunityConfigurationRepository{
		communityConfigurationDatabase: database.Postgres,
	}
}

func GetCommunityConfiguration(communityID int, communityConfiguration string) (*models.CommunityConfiguration, error) {
	configuration := &models.CommunityConfiguration{}

	filter := map[string]interface{}{
		"community_id": communityID,
		"type":         communityConfiguration,
	}

	err := NewCommunityConfigurationRepository().communityConfigurationDatabase.Where(filter).First(&configuration).Error
	if err != nil {
		return nil, err
	}

	return configuration, nil
}
