package helpers

import (
	"likeminds-pandemonium/api/models"
	"likeminds-pandemonium/api/repository"
)

func GetCommunityByID(communityID int) (*models.Community, error) {
	community, err := repository.GetCommunityByID(communityID)
	if err != nil {
		return nil, err
	}

	return community, nil
}
