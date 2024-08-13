package pubsub

const (
	RedisClient = "redis_client"
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
	ErrorPublish                 = "Publish error:"
	ErrorTopicMessageTypeMissing = "topic_message_type is required"
	ErrorTopicMissing            = "topic is required"
	ErrorTopicInvalid            = "topic is invalid"
)
