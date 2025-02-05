package requestresponse

type CaravanCreateMessageTaskRequest struct {
	ApiVersion      int                   `json:"api_version,omitempty"`
	ChatroomID      int                   `json:"chatroom_id,omitempty"`
	ChatroomStateID int                   `json:"chatroom_state_id,omitempty"`
	MemberState     int                   `json:"member_state,omitempty"`
	MessageID       int                   `json:"message_id,omitempty"`
	RequestBody     *CreateMessageRequest `json:"request_body,omitempty"`
	UserID          int                   `json:"user_id,omitempty"`
}
