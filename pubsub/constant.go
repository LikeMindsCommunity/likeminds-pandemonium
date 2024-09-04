package pubsub

const (
	RedisClient             = "redis_client"
	ReadBufferSizeDefault   = 4096
	WriteBufferSizeDefault  = 4096
	WsConnectionEstablished = "Connected to websocket server"
	PingWs                  = "Received ping from client"
	ReceivedMessageClientWs = "Received message from client having message type as %v"
	ReceivedMessageRedisWs  = "Received message from redis"
)

const (
	TopicTypeChatroom            = "chatroom"
	TopicTypeCommunity           = "community"
	TopicMessageTypeConversation = "conversation"
)

const (
	ParamTopic            = "topic"
	ParamTopicMessageType = "topic_message_type"
)

const (
	ErrorFailedUpgrader     = "Failed to upgrade connection: %v"
	ErrorPublishRedis       = "Publish error on redis:"
	ErrorReadClientWs       = "Error reading message from client: %v"
	ErrorReadDeadlineWs     = "Error while setting ReadDeadline on websocket:"
	ErrorWriteDeadlineWs    = "Error while setting WriteDeadline on websocket:"
	ErrorUnableToCloseWs    = "Unable to close ws error:"
	ErrorUnableToCloseRedis = "Unable to close redis error:"
	ErrorUnmarshalErrorJson = "Unmarshal error:"
	ErrorUnableToWriteWs    = "Unable to write message in websocket:"
	ErrorWriterOpenWs       = "Unable to open websocket Writer:"
	ErrorWriterCloseWs      = "Unable to close websocket Writer:"
)

const (
	ErrorUserUUIDMissing         = "User UUID is missing in headers"
	ErrorChatroomIDMissing       = "Chatroom ID is missing in request"
	ErrorCommunityIDMissing      = "Community ID is missing in request"
	ErrorTopicMessageTypeMissing = "Topic message type is missing from params"
	ErrorTopicMissing            = "Topic is missing from request"
	ErrorTopicInvalid            = "Invalid format of topic"
)
