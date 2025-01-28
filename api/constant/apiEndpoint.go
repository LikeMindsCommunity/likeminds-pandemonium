package constant

const (
	POSTMethod = 1
)

const (
	RedisPublish   = "/publish/:topic"
	RedisSubscribe = "/subscribe/:topic"
	DeliveryReport = "/delivery_report"
)

const (
	HeadersMemberID     = "x-member-id"
	HeadersDeviceID     = "x-device-id"
	HeadersPlatformCode = "x-platform-code"
	HeadersVersionCode  = "x-version-code"
	HeadersApiVersion   = "x-api-version"
)

const (
	HTTPResponseCodeOK                  = 200
	HTTPResonseCodeCreated              = 201
	HTTPResponseCodeBadRequest          = 404
	HTTPResponseCodeUnauthorised        = 401
	HTTPResponseCodeForbidden           = 403
	HTTPResponseCodeInternalServerError = 500
)
