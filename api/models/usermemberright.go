package models

type UserMemberRight struct {
	ID          int `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CommunityID int `gorm:"column:community_id;not null" json:"community_id"`
	RightID     int `gorm:"column:right_id;not null" json:"right_id"`
	UserID      int `gorm:"column:user_id;not null" json:"user_id"`
}

// TableName overrides the default table name
func (userMemberRight *UserMemberRight) TableName() string {
	return "togther_usermemberrights"
}
