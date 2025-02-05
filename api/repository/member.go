package repository

import (
	"fmt"
	"likeminds-pandemonium/api/constant"
	"likeminds-pandemonium/api/models"
	"likeminds-pandemonium/database"

	"gorm.io/gorm"
)

type MemberRepository struct {
	memberDatabase *gorm.DB
}

func NewMemberRepository() *MemberRepository {
	return &MemberRepository{
		memberDatabase: database.Postgres,
	}
}

func GetMemberStateInCommunity(communityID int, userID int) (*int, error) {
	memberState := int(0)

	filter := map[string]interface{}{
		"community_id_id": communityID,
		"member_id_id":    userID,
	}

	err := NewMemberRepository().memberDatabase.Model(&models.Member{}).Select("state").Where(filter).First(&memberState).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get member state in community, community id=%d, user id=%d, err=%s", communityID, userID, err)
	}

	return &memberState, nil
}

func IsMemberVerifiedInCommunity(communityID int, userID int) (bool, error) {
	member := &models.Member{}
	filter := map[string]interface{}{
		"community_id_id": communityID,
		"member_id_id":    userID,
		"state": []int{
			constant.UserStateAdmin,
			constant.UserStateMember,
			constant.UserStateProfileUnavailable,
			constant.UserStateKnownNominatedPromoter},
	}

	err := NewMemberRepository().memberDatabase.Where(filter).First(&member).Error
	if err != nil {
		return false, fmt.Errorf("failed to verify member in community, community id=%d, user id=%d, err=%s", communityID, userID, err)
	}

	return true, nil
}
