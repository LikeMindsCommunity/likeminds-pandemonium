package models

type Member struct {
	ID               int64   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CommunityIDID    int     `gorm:"column:community_id_id;not null" json:"community_id_id"`
	MemberIDID       int     `gorm:"column:member_id_id;not null" json:"member_id_id"`
	State            *int    `gorm:"column:state" json:"state,omitempty"`
	CreatedAt        int64   `gorm:"column:created_at;not null" json:"created_at"`
	ToolState        int     `gorm:"column:tool_state;not null" json:"tool_state"`
	ApprovedMemberID *int    `gorm:"column:approved_member_id" json:"approved_member_id,omitempty"`
	AskMemberID      *int    `gorm:"column:ask_member_id" json:"ask_member_id,omitempty"`
	EditRequired     bool    `gorm:"column:edit_required;not null" json:"edit_required"`
	ActionsRequired  *bool   `gorm:"column:actions_required" json:"actions_required,omitempty"`
	ImageURL         *string `gorm:"column:image_url;type:text" json:"image_url,omitempty"`
	UpdatedAt        int64   `gorm:"column:updated_at;not null" json:"updated_at"`
	ApprovedByID     *int    `gorm:"column:approved_by_id" json:"approved_by_id,omitempty"`
	CustomTitle      *string `gorm:"column:custom_title;type:text" json:"custom_title,omitempty"`
	IsOwner          bool    `gorm:"column:is_owner;not null" json:"is_owner"`
	JoinedByID       *int    `gorm:"column:joined_by_id" json:"joined_by_id,omitempty"`
	ParentCMID       *int    `gorm:"column:parent_cm_id" json:"parent_cm_id,omitempty"`
	ParentCMList     *string `gorm:"column:parent_cm_list;type:text" json:"parent_cm_list,omitempty"`
	BecameMemberAt   int64   `gorm:"column:became_member_at;not null" json:"became_member_at"`
	HasOnboarded     bool    `gorm:"column:has_onboarded;not null" json:"has_onboarded"`
}

func (member *Member) TableName() string {
	return "togther_members"
}
