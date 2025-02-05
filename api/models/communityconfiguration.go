package models

type CommunityConfiguration struct {
	ID          int    `gorm:"primaryKey;column:id" json:"id"`
	Type        string `gorm:"column:type;size:255;not null" json:"type"`
	Description string `gorm:"column:description;not null" json:"description"`
	Value       string `gorm:"column:value;type:jsonb;not null" json:"value"`
	CreatedAt   int64  `gorm:"column:created_at;not_null" json:"created_at"`
	UpdatedAt   int64  `gorm:"column:updated_at;not_null" json:"updated_at"`
	CommunityID int    `gorm:"column:community_id;not null" json:"community_id"`
}

func (communityConfiguration *CommunityConfiguration) TableName() string {
	return "togther_communityconfigurations"
}
