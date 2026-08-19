package repository

import (
	"fmt"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/api/models"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/database"

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
		return nil, fmt.Errorf("failed to get community configuration, community id=%d, configuration=%s, err=%s", communityID, communityConfiguration, err)
	}

	return configuration, nil
}
