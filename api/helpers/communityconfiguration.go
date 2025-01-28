package helpers

import (
	"likeminds-pandemonium/api/models"
	"likeminds-pandemonium/api/repository"
)

func GetCommunityConfiguration(communityID int, communityConfiguration string) (*models.CommunityConfiguration, error) {
	configuration, err := repository.GetCommunityConfiguration(communityID, communityConfiguration)
	if err != nil {
		return nil, err
	}

	return configuration, nil
}
