package models

type ConversationResponse struct {
	Conversation Conversation `json:"conversation"`
	Participants []string     `json:"participants"`
}

type Conversation struct {
	Member Member `json:"member"`
}
