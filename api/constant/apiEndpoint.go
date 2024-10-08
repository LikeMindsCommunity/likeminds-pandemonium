package constant

const (
	POSTMethod = 1
)

const (
	RedisPublish   = "/publish/:topic"
	RedisSubscribe = "/subscribe/:topic"
	SentDR         = "/sent_dr"
	DeliveredDR    = "/delivered_dr"
)

const (
	HeadersMemberID = "x-member-id"
	HeadersDeviceID = "x-device-id"
)
