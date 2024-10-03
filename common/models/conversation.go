package models

type ConversationResponse struct {
	Conversation Conversation `json:"conversation"`
	Participants []string     `json:"participants"`
}

type Conversation struct {
	ID         string `json:"id"`
	ChatroomID string `json:"chatroom_id"`
	Member     Member `json:"member"`
}
