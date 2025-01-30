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
	HeadersApiKey       = "x-api-key"
	HeadersDeviceID     = "x-device-id"
	HeadersPlatformCode = "x-platform-code"
	HeadersVersionCode  = "x-version-code"
	HeadersApiVersion   = "x-api-version"
	HeadersPlatformType = "x-platform-type"
)

const (
	HTTPResponseCodeOK                  = 200
	HTTPResonseCodeCreated              = 201
	HTTPResponseCodeBadRequest          = 404
	HTTPResponseCodeUnauthorised        = 401
	HTTPResponseCodeForbidden           = 403
	HTTPResponseCodeInternalServerError = 500
)

const (
	PlatformTypePandemoniumService = "pandemonium-service"
)

type ApiHeaders struct {
	HeadersMemberID     string `json:"x-member-id,omitempty"`
	HeadersApiKey       string `json:"x-api-key,omitempty"`
	HeadersDeviceID     string `json:"x-device-id,omitempty"`
	HeadersPlatformCode string `json:"x-platform-code,omitempty"`
	HeadersVersionCode  string `json:"x-version-code,omitempty"`
	HeadersApiVersion   string `json:"x-api-version,omitempty"`
	HeadersPlatformType string `json:"x-platform-type,omitempty"`
}
