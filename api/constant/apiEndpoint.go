package constant

const (
	POSTMethod = 1
)

const (
	RedisPublish   = "/publish/:topic"
	RedisSubscribe = "/subscribe/:topic"
	DeliveryReport = "/deliver_report"
)

const (
	HeadersMemberID = "x-member-id"
	HeadersDeviceID = "x-device-id"
)
