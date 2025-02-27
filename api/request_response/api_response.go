package requestresponse

type PSRequest struct {
	TopicMessageType string `json:"topic_message_type"`
	RawData          []byte `json:"raw_data"`
}

type PSResponse struct {
	DeviceID         string `json:"device_id"`
	TopicMessageType string `json:"topic_message_type"`
	RawData          []byte `json:"raw_data"`
}
