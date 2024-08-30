package pubsub

const (
	RedisClient            = "redis_client"
	ReadBufferSizeDefault  = 4096
	WriteBufferSizeDefault = 4096
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
	ErrorPublishRedis       = "Publish error on redis:"
	ErrorUnexpectedCloseWs  = "Unexpected close error on Ws:"
	ErrorReadDeadlineWs     = "Error while setting ReadDeadline on Ws:"
	ErrorWriteDeadlineWs    = "Error while setting WriteDeadline on Ws:"
	ErrorUnableToCloseWs    = "Unable to close ws error:"
	ErrorUnableToCloseRedis = "Unable to close redis error:"
	ErrorUnmarshalErrorJson = "Unmarshal error:"
	ErrorUnableToWriteWs    = "Unable to write message in Ws:"
	ErrorWriterOpenWs       = "Unable to open Ws Writer:"
	ErrorWriterCloseWs      = "Unable to close Ws Writer:"
)

const (
	ErrorUserUUIDMissing         = "User UUID is missing in headers"
	ErrorChatroomIDMissing       = "Chatroom ID is missing in request"
	ErrorCommunityIDMissing      = "Community ID is missing in request"
	ErrorTopicMessageTypeMissing = "Topic message type is missing from params"
	ErrorTopicMissing            = "Topic is missing from request"
	ErrorTopicInvalid            = "Invalid format of topic"
)
