package helpers

import (
	"fmt"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/api/models"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/api/repository"
)

func GetCommunityByID(communityID int64) (*models.Community, error) {
	community, err := repository.GetCommunityByID(communityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get community, id=%d, err=%s", communityID, err)
	}

	return community, nil
}
