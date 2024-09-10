package models

type ConversationResponse struct {
	Conversation Conversation `json:"conversation"`
}

type Conversation struct {
	Member       Member   `json:"member"`
	Participants []string `json:"participants"`
}
