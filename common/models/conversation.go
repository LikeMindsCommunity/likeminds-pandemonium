package models

type ConversationResponse struct {
	DeviceID     string       `json:"device_id"`
	Conversation Conversation `json:"conversation"`
}

type Conversation struct {
	Member Member `json:"member"`
}
