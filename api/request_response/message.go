package requestresponse

import (
	"likeminds-pandemonium/api/models"
)

type CreateMessageRequest struct {
	ChatroomID                  int                 `json:"chatroom_id"`
	Text                        string              `json:"text"`
	Header                      *string             `json:"header,omitempty"`
	OnlineLink                  *string             `json:"online_link,omitempty"`
	OnlineLinkID                *string             `json:"online_link_id,omitempty"`
	OnlineLinkPassword          *string             `json:"online_link_password,omitempty"`
	OnlineLinkEnableBefore      int64               `json:"online_link_enable_before,omitempty"`
	Location                    *string             `json:"location,omitempty"`
	LocationLat                 *string             `json:"location_lat,omitempty"`
	LocationLong                *string             `json:"location_long,omitempty"`
	StartTime                   int64               `json:"start_time,omitempty"`
	EndTime                     int64               `json:"end_time,omitempty"`
	CoHosts                     *string             `json:"co_hosts,omitempty"`
	PollType                    *int32              `json:"poll_type,omitempty"`
	AllowAddOption              bool                `json:"allow_add_option,omitempty"`
	ExpiryTime                  *int64              `json:"expiry_time,omitempty"`
	Polls                       []PollObject        `json:"polls,omitempty"`
	MultilpleSelectState        *int64              `json:"multiple_select_state,omitempty"`
	MultilpleSelectNo           *int64              `json:"multiple_select_no,omitempty"`
	NoPollExpiry                bool                `json:"no_poll_expiry,omitempty"`
	AllowVoteChange             bool                `json:"allow_vote_change,omitempty"`
	AttachmentCount             int64               `json:"attachment_count,omitempty"`
	RepliedConversationId       interface{}         `json:"replied_conversation_id,omitempty"`
	RepliedChatroomID           string              `json:"replied_chatroom_id,omitempty"`
	InternalLink                string              `json:"internal_link,omitempty"`
	Preview                     ConversationPreview `json:"preview,omitempty"`
	IsAnonymous                 bool                `json:"is_anonymous,omitempty"`
	State                       int32               `json:"state"`
	HasFiles                    bool                `json:"has_files,omitempty"`
	TemporaryID                 string              `json:"temporary_id,omitempty"`
	OGTags                      interface{}         `json:"og_tags,omitempty"`
	ShareLink                   string              `json:"share_link,omitempty"`
	Attachments                 []MessageAttachment `json:"attachments,omitempty"`
	Metadata                    interface{}         `json:"metadata,omitempty"`
	TriggerBot                  bool                `json:"trigger_bot,omitempty"`
	ShouldStreamChatbotResponse bool                `json:"should_stream_chatbot_response,omitempty"`
}

type CreateMessageResponse struct {
	HTTPStatusCode int              `json:"http_status_code"`
	Success        bool             `json:"success"`
	Data           *MessageResponse `json:"data"`
	Error          string           `json:"error"`
}

type MessageResponse struct {
	Message        *models.Message             `json:"message"`
	Attachments    *[]models.MessageAttachment `json:"attachments"`
	Polls          *[]models.MessagePoll       `json:"polls"`
	RepliedMessage *models.Message             `json:"message_replied"`
	Widget         *SwarmWidgetResponse        `json:"widget"`
	User           *models.UserInfo            `json:"user"`
}

type PollObject struct {
	Text           string `json:"text"`
	ConversationID int    `json:"conversation_id"`
	UserID         int    `json:"user_id"`
}

type ConversationPreview struct {
	InternalLink string          `json:"internal_link"`
	PreviewType  string          `json:"preview_type"`
	PreviewText  string          `json:"preview_text"`
	Title        string          `json:"title"`
	Community    CommunityObject `json:"community"`
	Action       string          `json:"action"`
	ActionRoute  string          `json:"action_route"`
}

type CommunityObject struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Purpose        string `json:"purpose"`
	ImageUrl       string `json:"image_url"`
	ImageUrlRound  string `json:"image_url_round"`
	CreatedBy      string `json:"created_by"`
	PromotersCount int32  `json:"promoters_count"`
	MembersCount   int32  `json:"members_count"`
	MemberState    int32  `json:"member_state"`
}

type MessageAttachment struct {
	Name         string      `json:"name,omitempty"`
	FileURL      string      `json:"url,omitempty"`
	Type         string      `json:"type,omitempty"`
	ThumbnailURL string      `json:"thumbnail_url,omitempty"`
	Index        int         `json:"index,omitempty"`
	Height       int         `json:"height,omitempty"`
	Width        int         `json:"width,omitempty"`
	Meta         interface{} `json:"meta,omitempty"`
	LocationName string      `json:"location_name,omitempty"`
	LocationLat  *float64    `json:"location_lat,omitempty"`
	LocationLong *float64    `json:"location_long,omitempty"`
}
