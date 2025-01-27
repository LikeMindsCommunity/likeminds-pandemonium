package helpers

import (
	"likeminds-pandemonium/api/constant"
	"likeminds-pandemonium/api/repository"
)

func GetMemberStateInCommunity(communityID int, userID int) (*int, error) {
	memberState, err := repository.GetMemberStateInCommunity(communityID, userID)
	if err != nil {
		return nil, err
	}

	return memberState, nil
}

func IsMemberStateAdminInCommunity(memberState *int) bool {
	return *memberState == constant.UserStateAdmin
}

func IsMemberVerifiedInCommunity(communityID int, userID int) bool {
	return repository.IsMemberVerifiedInCommunity(communityID, userID)
}
