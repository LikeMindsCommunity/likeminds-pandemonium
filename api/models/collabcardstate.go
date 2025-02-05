package models

type CollabCardState struct {
	ID                       int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	State                    *int   `gorm:"column:state" json:"state,omitempty"`
	CreatedAt                *int64 `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedAt                *int64 `gorm:"column:updated_at" json:"updated_at,omitempty"`
	CardID                   int    `gorm:"column:card_id;not null" json:"card_id"`
	CommunityID              int    `gorm:"column:community_id;not null" json:"community_id"`
	UserID                   int    `gorm:"column:user_id;not null" json:"user_id"`
	MuteStatus               bool   `gorm:"column:mute_status;not null" json:"mute_status"`
	FollowStatus             bool   `gorm:"column:follow_status;not null" json:"follow_status"`
	IsGuest                  bool   `gorm:"column:is_guest;not null" json:"is_guest"`
	RemoveID                 *int   `gorm:"column:remove_id" json:"remove_id,omitempty"`
	SourceID                 *int   `gorm:"column:source_id" json:"source_id,omitempty"`
	ExpiryTime               *int64 `gorm:"column:expiry_time" json:"expiry_time,omitempty"`
	ExternalSeen             bool   `gorm:"column:external_seen;not null" json:"external_seen"`
	IsTagged                 bool   `gorm:"column:is_tagged;not null" json:"is_tagged"`
	ExternalFollow           bool   `gorm:"column:external_follow;not null" json:"external_follow"`
	AttendingStatus          *bool  `gorm:"column:attending_status" json:"attending_status,omitempty"`
	ManualSetActive          *int64 `gorm:"column:manual_set_active" json:"manual_set_active,omitempty"`
	LastSeenConversationID   *int   `gorm:"column:last_seen_conversation_id" json:"last_seen_conversation_id,omitempty"`
	SecretChatroomLeft       bool   `gorm:"column:secret_chatroom_left;not null" json:"secret_chatroom_left"`
	Attended                 bool   `gorm:"column:attended;not null" json:"attended"`
	ChatRequestCreatedAt     *int64 `gorm:"column:chat_request_created_at" json:"chat_request_created_at,omitempty"`
	ChatRequestState         *int   `gorm:"column:chat_request_state" json:"chat_request_state,omitempty"`
	ChatRequestedByID        *int   `gorm:"column:chat_requested_by_id" json:"chat_requested_by_id,omitempty"`
	IsNotificationPaused     bool   `gorm:"column:is_noti_paused;not null" json:"is_noti_paused"`
	NotificationState        int    `gorm:"column:noti_state;not null" json:"noti_state"`
	UnpauseNotificationAt    *int64 `gorm:"column:unpause_noti_at" json:"unpause_noti_at,omitempty"`
	ChatRequestInitiatedByID *int   `gorm:"column:chat_request_initiated_by_id" json:"chat_request_initiated_by_id,omitempty"`
}

func (collabCardState *CollabCardState) TableName() string {
	return "togther_collabcardstate"
}
