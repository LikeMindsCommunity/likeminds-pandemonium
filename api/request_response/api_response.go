package requestresponse

type PSResponse struct {
	DeviceID         string `json:"device_id"`
	TopicMessageType string `json:"topic_message_type"`
	RawData          string `json:"raw_data"`
}
