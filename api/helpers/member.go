package helpers

import (
	"fmt"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/api/constant"
	"github.com/LikeMindsCommunity/likeminds-pandemonium/api/repository"
)

func GetMemberStateInCommunity(communityID int, userID int) (*int, error) {
	memberState, err := repository.GetMemberStateInCommunity(communityID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get member state in community, community id=%d, user id=%d, err=%s", communityID, userID, err)
	}

	return memberState, nil
}

func IsMemberStateAdminInCommunity(memberState *int) bool {
	return *memberState == constant.UserStateAdmin
}

func IsMemberVerifiedInCommunity(communityID int, userID int) (bool, error) {
	memberIsVerifiedInCommunity, err := repository.IsMemberVerifiedInCommunity(communityID, userID)
	if err != nil {
		return false, fmt.Errorf("failed to verify member in community, community id=%d, user id=%d, err=%s", communityID, userID, err)
	}

	return memberIsVerifiedInCommunity, nil
}
