package repository

import (
	"likeminds-pandemonium/api/constant"
	"likeminds-pandemonium/api/models"
	"likeminds-pandemonium/database"
	"log"

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
		return nil, err
	}

	return &memberState, nil
}

func IsMemberVerifiedInCommunity(communityID int, userID int) bool {
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
		log.Print("error fetching member")
		return false
	}

	// if row id is default value
	if member.ID == 0 {
		log.Print("member response error")
		return false
	}

	return true
}
