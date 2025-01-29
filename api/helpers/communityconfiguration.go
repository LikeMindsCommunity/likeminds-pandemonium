package helpers

import (
	"fmt"
	"likeminds-pandemonium/api/models"
	"likeminds-pandemonium/api/repository"
)

func GetCommunityConfiguration(communityID int, communityConfiguration string) (*models.CommunityConfiguration, error) {
	configuration, err := repository.GetCommunityConfiguration(communityID, communityConfiguration)
	if err != nil {
		return nil, fmt.Errorf("failed to get community configuration, community id=%d, configuration=%s, err=%s", communityID, communityConfiguration, err)
	}

	return configuration, nil
}
