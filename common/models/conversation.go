package models

type ConversationResponse struct {
	Conversation Conversation `json:"conversation"`
	Participants []string     `json:"participants"`
}

type Conversation struct {
	ID          interface{} `json:"id"`
	ChatroomID  interface{} `json:"chatroom_id"`
	CommunityID interface{} `json:"community_id"`
	Member      Member      `json:"member"`
}
