package models

type Message struct {
	ID                     int64   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Answer                 string  `gorm:"column:answer;type:text;not null" json:"answer"`
	CreatedAt              int64   `gorm:"column:created_at;not null" json:"created_at"`
	State                  int     `gorm:"column:state;not null" json:"state"`
	CardID                 int     `gorm:"column:card_id;not null" json:"card_id"`
	UserID                 int     `gorm:"column:user_id;not null" json:"user_id"`
	OgTags                 *string `gorm:"column:og_tags;type:text" json:"og_tags,omitempty"`
	CommunityID            int     `gorm:"column:community_id" json:"community_id"`
	IsGuest                bool    `gorm:"column:is_guest;not null" json:"is_guest"`
	RemoveID               *int    `gorm:"column:remove_id" json:"remove_id,omitempty"`
	IsDeleted              bool    `gorm:"column:is_deleted;not null" json:"is_deleted"`
	IsEdited               bool    `gorm:"column:is_edited;not null" json:"is_edited"`
	ReplyID                *int64  `gorm:"column:reply_id" json:"reply_id,omitempty"`
	InternalLink           *string `gorm:"column:internal_link;type:text" json:"internal_link,omitempty"`
	PreviewChatroomID      *int    `gorm:"column:preview_chatroom_id" json:"preview_chatroom_id,omitempty"`
	PreviewCommunityID     *int    `gorm:"column:preview_community_id" json:"preview_community_id,omitempty"`
	PreviewType            *string `gorm:"column:preview_type;type:text" json:"preview_type,omitempty"`
	HasFiles               bool    `gorm:"column:has_files;not null" json:"has_files"`
	DeletedByUserID        *int    `gorm:"column:deleted_by_user_id" json:"deleted_by_user_id,omitempty"`
	LastUpdated            int64   `gorm:"column:last_updated;not null" json:"last_updated"`
	AttachmentCount        int     `gorm:"column:attachment_count;not null" json:"attachment_count"`
	AttachmentsUploaded    bool    `gorm:"column:attachments_uploaded" json:"attachments_uploaded"`
	APIVersion             int     `gorm:"column:api_version;not null" json:"api_version"`
	DeviceID               *string `gorm:"column:device_id;type:text" json:"device_id,omitempty"`
	Platform               *string `gorm:"column:platform;type:text" json:"platform,omitempty"`
	TemporaryID            *string `gorm:"column:temporary_id;type:text" json:"temporary_id,omitempty"`
	AllowAddOption         bool    `gorm:"column:allow_add_option;not null" json:"allow_add_option"`
	ExpiryTime             *int64  `gorm:"column:expiry_time" json:"expiry_time,omitempty"`
	IsAnonymous            bool    `gorm:"column:is_anonymous;not null" json:"is_anonymous"`
	MultipleSelectNo       *int    `gorm:"column:multiple_select_no" json:"multiple_select_no,omitempty"`
	MultipleSelectState    *int    `gorm:"column:multiple_select_state" json:"multiple_select_state,omitempty"`
	PollType               *int    `gorm:"column:poll_type" json:"poll_type,omitempty"`
	HasReactions           bool    `gorm:"column:has_reactions;not null" json:"has_reactions"`
	PollAnswerText         string  `gorm:"column:poll_answer_text;type:text;not null" json:"poll_answer_text"`
	ReplyChatroomID        *int    `gorm:"column:reply_chatroom_id" json:"reply_chatroom_id,omitempty"`
	CoHosts                *string `gorm:"column:co_hosts;type:text" json:"co_hosts,omitempty"`
	EndTime                int64   `gorm:"column:end_time;not null" json:"end_time"`
	Header                 *string `gorm:"column:header;type:text" json:"header,omitempty"`
	Location               *string `gorm:"column:location;type:text" json:"location,omitempty"`
	LocationLat            *string `gorm:"column:location_lat;type:text" json:"location_lat,omitempty"`
	LocationLong           *string `gorm:"column:location_long;type:text" json:"location_long,omitempty"`
	OnlineLink             *string `gorm:"column:online_link;type:text" json:"online_link,omitempty"`
	OnlineLinkEnableBefore int64   `gorm:"column:online_link_enable_before;not null" json:"online_link_enable_before"`
	OnlineLinkID           *string `gorm:"column:online_link_id;type:text" json:"online_link_id,omitempty"`
	OnlineLinkPassword     *string `gorm:"column:online_link_password;type:text" json:"online_link_password,omitempty"`
	StartTime              int64   `gorm:"column:start_time;not null" json:"start_time"`
	AboutRecording         *string `gorm:"column:about_recording;type:text" json:"about_recording,omitempty"`
	HasEventRecording      bool    `gorm:"column:has_event_recording;not null" json:"has_event_recording"`
	RecordingURLTags       *string `gorm:"column:recording_url_og_tags;type:text" json:"recording_url_og_tags,omitempty"`
	WidgetID               string  `gorm:"column:widget_id;type:varchar(255);not null" json:"widget_id"`
	AllowVoteChange        bool    `gorm:"column:allow_vote_change;not null" json:"allow_vote_change"`
	NoPollExpiry           bool    `gorm:"column:no_poll_expiry;not null" json:"no_poll_expiry"`
}

func (message *Message) TableName() string {
	return "togther_card_answers"
}
