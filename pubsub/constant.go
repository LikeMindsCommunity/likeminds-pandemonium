package pubsub

const (
	RedisClient = "redis_client"
)

const (
	TopicTypeChatroom            = "chatroom"
	TopicNameChatroom            = TopicTypeChatroom + ":%s"
	TopicMessageTypeConversation = "conversation"
)

const (
	TopicCommunityType = "community"
	TopicCommunity     = TopicCommunityType + ":%s"
)

const (
	ParamTopic            = "topic"
	ParamTopicMessageType = "topic_message_type"
)
