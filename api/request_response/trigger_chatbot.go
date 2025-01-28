package requestresponse

type TriggerChatbot struct {
	ChatroomID                  int  `json:"chatroom_id,omitempty"`
	MessageID                   int  `json:"message_id,omitempty"`
	ShouldStreamChatbotResponse bool `json:"should_stream_chatbot_response,omitempty"`
	ApiVersion                  int  `json:"api_version,omitempty"`
}
