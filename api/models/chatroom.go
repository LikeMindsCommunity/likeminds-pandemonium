package models

type Chatroom struct {
	ID                          int64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title                       string   `gorm:"column:title;type:text;not null" json:"title"`
	LikesCount                  int      `gorm:"column:likes_count;not null" json:"likes_count"`
	ShareCount                  int      `gorm:"column:share_count;not null" json:"share_count"`
	AnswersCount                int      `gorm:"column:answers_count;not null" json:"answers_count"`
	CommunityID                 int      `gorm:"column:community_id;not null;index" json:"community_id"`
	UserID                      int      `gorm:"column:user_id;not null;index" json:"user_id"`
	DateEpoch                   int64    `gorm:"column:date_epoch;not null" json:"date_epoch"`
	AnswerText                  string   `gorm:"column:answer_text;type:varchar(100);not null" json:"answer_text"`
	ShareLink                   string   `gorm:"column:share_link;type:varchar(2048);not null" json:"share_link"`
	OGTags                      string   `gorm:"column:og_tags;type:text;not null" json:"og_tags"`
	ImageCount                  *int     `gorm:"column:image_count" json:"image_count,omitempty"`
	PDFCount                    *int     `gorm:"column:pdf_count" json:"pdf_count,omitempty"`
	Type                        int      `gorm:"column:type;not null" json:"type"`
	DateTime                    int64    `gorm:"column:date_time;not null" json:"date_time"`
	Duration                    int64    `gorm:"column:duration;not null" json:"duration"`
	AttendingCount              int      `gorm:"column:attending_count;not null" json:"attending_count"`
	PollsCount                  int      `gorm:"column:polls_count;not null" json:"polls_count"`
	Location                    *string  `gorm:"column:location;type:text" json:"location,omitempty"`
	LocationLat                 *float64 `gorm:"column:location_lat" json:"location_lat,omitempty"`
	LocationLong                *float64 `gorm:"column:location_long" json:"location_long,omitempty"`
	About                       *string  `gorm:"column:about;type:text" json:"about,omitempty"`
	CoHosts                     *string  `gorm:"column:co_hosts;type:text" json:"co_hosts,omitempty"`
	EndDate                     *int64   `gorm:"column:end_date" json:"end_date,omitempty"`
	OnlineLink                  *string  `gorm:"column:online_link;type:text" json:"online_link,omitempty"`
	StartDate                   *int64   `gorm:"column:start_date" json:"start_date,omitempty"`
	UpdatedMemberID             *int     `gorm:"column:updated_member_id;index" json:"updated_member_id,omitempty"`
	UpdatedTime                 int64    `gorm:"column:updated_time;not null" json:"updated_time"`
	MultipleSelect              bool     `gorm:"column:multiple_select;not null" json:"multiple_select"`
	MultipleSelectNo            *int     `gorm:"column:multiple_select_no" json:"multiple_select_no,omitempty"`
	MultipleSelectState         *int     `gorm:"column:multiple_select_state" json:"multiple_select_state,omitempty"`
	Header                      *string  `gorm:"column:header;type:text" json:"header,omitempty"`
	HasBeenNamed                bool     `gorm:"column:has_been_named;not null" json:"has_been_named"`
	AllowAddOption              *bool    `gorm:"column:allow_add_option" json:"allow_add_option,omitempty"`
	IsPollAnonymous             *bool    `gorm:"column:is_poll_anonymous" json:"is_poll_anonymous,omitempty"`
	PollType                    *int     `gorm:"column:poll_type" json:"poll_type,omitempty"`
	InternalLink                *string  `gorm:"column:internal_link;type:text" json:"internal_link,omitempty"`
	PreviewChatroomID           *int     `gorm:"column:preview_chatroom_id;index" json:"preview_chatroom_id,omitempty"`
	PreviewCommunityID          *int     `gorm:"column:preview_community_id;index" json:"preview_community_id,omitempty"`
	PreviewType                 *string  `gorm:"column:preview_type;type:text" json:"preview_type,omitempty"`
	DeletedByUserID             *int     `gorm:"column:deleted_by_user_id;index" json:"deleted_by_user_id,omitempty"`
	IsDeleted                   bool     `gorm:"column:is_deleted;not null" json:"is_deleted"`
	IsPending                   bool     `gorm:"column:is_pending;not null" json:"is_pending"`
	Reason                      *string  `gorm:"column:reason;type:varchar(512)" json:"reason,omitempty"`
	TagID                       *int     `gorm:"column:tag_id;index" json:"tag_id,omitempty"`
	MemberState                 *int     `gorm:"column:member_state" json:"member_state,omitempty"`
	DisablePollAnnouncementMail bool     `gorm:"column:disable_poll_announcement_mail;not null" json:"disable_poll_announcement_mail"`
	HasFiles                    bool     `gorm:"column:has_files;not null" json:"has_files"`
	AudioCount                  *int     `gorm:"column:audio_count" json:"audio_count,omitempty"`
	VideoCount                  *int     `gorm:"column:video_count" json:"video_count,omitempty"`
	AttachmentCount             int      `gorm:"column:attachment_count;not null" json:"attachment_count"`
	AttachmentsUploaded         *bool    `gorm:"column:attachments_uploaded" json:"attachments_uploaded,omitempty"`
	IsPinned                    bool     `gorm:"column:is_pinned;not null" json:"is_pinned"`
	PinningTime                 int64    `gorm:"column:pinning_time;not null" json:"pinning_time"`
	IsSecret                    bool     `gorm:"column:is_secret;not null" json:"is_secret"`
	SecretChatroomParticipants  *string  `gorm:"column:secret_chatroom_participants;type:text" json:"secret_chatroom_participants,omitempty"`
	HasReactions                bool     `gorm:"column:has_reactions;not null" json:"has_reactions"`
	DeviceID                    *string  `gorm:"column:device_id;type:text" json:"device_id,omitempty"`
	Platform                    *string  `gorm:"column:platform;type:text" json:"platform,omitempty"`
	TopicID                     *int     `gorm:"column:topic_id;index" json:"topic_id,omitempty"`
	AutoFollowDone              bool     `gorm:"column:auto_follow_done;not null" json:"auto_follow_done"`
	IsEdited                    bool     `gorm:"column:is_edited;not null" json:"is_edited"`
	Access                      *int     `gorm:"column:access" json:"access,omitempty"`
	CreatedAt                   int64    `gorm:"column:created_at;not null" json:"created_at"`
	EventPaymentLink            *string  `gorm:"column:event_payment_link;type:text" json:"event_payment_link,omitempty"`
	IsPaid                      bool     `gorm:"column:is_paid;not null" json:"is_paid"`
	OnlineLinkEnableBefore      int64    `gorm:"column:online_link_enable_before;not null" json:"online_link_enable_before"`
	OnlineLinkID                *string  `gorm:"column:online_link_id;type:text" json:"online_link_id,omitempty"`
	OnlineLinkPassword          *string  `gorm:"column:online_link_password;type:text" json:"online_link_password,omitempty"`
	UpdatedAt                   int64    `gorm:"column:updated_at;not null" json:"updated_at"`
	EventWebPage                *string  `gorm:"column:event_web_page;type:text" json:"event_web_page,omitempty"`
	WebflowItemID               *string  `gorm:"column:webflow_item_id;type:text" json:"webflow_item_id,omitempty"`
	MemberCanMessage            bool     `gorm:"column:member_can_message;not null" json:"member_can_message"`
	ChatroomWithUserID          *int     `gorm:"column:chatroom_with_user_id;index" json:"chatroom_with_user_id,omitempty"`
	IsPrivate                   bool     `gorm:"column:is_private;not null" json:"is_private"`
	IncludeMembersLater         bool     `gorm:"column:include_members_later;not null" json:"include_members_later"`
	AboutRecording              *string  `gorm:"column:about_recording;type:text" json:"about_recording,omitempty"`
	HasEventRecording           bool     `gorm:"column:has_event_recording;not null" json:"has_event_recording"`
	RecordingURLOGTags          *string  `gorm:"column:recording_url_og_tags;type:text" json:"recording_url_og_tags,omitempty"`
	AccessWithoutSubscription   bool     `gorm:"column:access_without_subscription;not null" json:"access_without_subscription"`
	OnlineLinkType              *int     `gorm:"column:online_link_type" json:"online_link_type,omitempty"`
	IsPrivateMember             bool     `gorm:"column:is_private_member;not null" json:"is_private_member"`
	ThirdPartyUniqueID          *string  `gorm:"column:third_party_unique_id;type:text" json:"third_party_unique_id,omitempty"`
	SingleEventURL              *string  `gorm:"column:single_event_url;type:text" json:"single_event_url,omitempty"`
	ChatroomImageURL            *string  `gorm:"column:chatroom_image_url;type:text" json:"chatroom_image_url,omitempty"`
	TagOnlyParticipants         bool     `gorm:"column:tag_only_participants;not null" json:"tag_only_participants"`
	CustomTag                   *string  `gorm:"column:custom_tag;type:text" json:"custom_tag,omitempty"`
	EventKind                   string   `gorm:"column:event_kind;type:text;not null" json:"event_kind"`
}

func (chatroom *Chatroom) TableName() string {
	return "togther_collabcard"
}
