package common

import "time"

const (
	RedisClient                  = "redis_client"
	ReadBufferSizeDefault        = 4096
	WriteBufferSizeDefault       = 4096
	ParamTopic                   = "topic"
	ParamTopicMessageType        = "topic_message_type"
	TopicTypeChatroom            = "chatroom"
	TopicTypeCommunity           = "community"
	TopicMessageTypeConversation = "conversation"
	TopicMessageTypeSentReport   = "sent_report"
	ChatroomDeliveryReportPrefix = "chatroom_%s_delivery_report" // Redis key prefix for chatroom delivery reports
	// WriteWait Max wait time when writing message to peer
	WriteWait = 10 * time.Second
	// PongWait Max time till next pong from peer
	PongWait = 60 * time.Second
	// PingPeriod should be less than PongWait
	PingPeriod  = (PongWait * 9) / 10
	WsServerKey = "ws_server"
)

const (
	WsConnectionEstablished  = "connected to websocket server"
	PingReceivedClient       = "received ping from client"
	PongReceivedClient       = "received pong from client"
	PingSendClient           = "sending ping to client"
	PongSendClient           = "sending pong to client"
	ReceivedMessageClientWs  = "received message from client having message type as %v"
	ReceivedMessageRedisWs   = "received message from redis"
	ConnectionClosed         = "connection closed"
	LogSuccessCacheSaveRedis = "successfully set hash: %s with field: %s"
)

const (
	ErrorFailedUpgrader        = "failed to upgrade connection: %v"
	ErrorPublishRedis          = "publish error on redis:"
	ErrorReadClientWs          = "error reading message from client: %v"
	ErrorReadDeadlineWs        = "error while setting ReadDeadline on websocket:"
	ErrorWriteDeadlineWs       = "error while setting WriteDeadline on websocket:"
	ErrorUnableToCloseWs       = "unable to close ws error:"
	ErrorUnableToCloseRedis    = "unable to close redis error:"
	ErrorUnmarshalErrorJson    = "unmarshal error: %v"
	ErrorUnableToWriteWs       = "unable to write message in websocket: %v"
	ErrorWriterOpenWs          = "unable to open websocket writer: %v"
	ErrorWriterCloseWs         = "unable to close websocket writer: %v"
	ErrorPingSentClient        = "error sending ping to client: %v"
	ErrorPongSentClient        = "error sending pong to client: %v"
	ErrorFailedCacheSaveRedis  = "failed to save cache: %v"
	ErrorFailedExpSaveRedis    = "failed to save cache exp: %v"
	ErrorFailedCacheFetchRedis = "failed to fetch from redis cache: %v"
	ErrorCacheEmptyRedis       = "no data found for key: %s"
	ErrorMarshalErrorJson      = "marshal error: %v"
)

const (
	ErrorUserUUIDMissing         = "user UUID is missing in headers"
	ErrorChatroomIDMissing       = "chatroom ID is missing in request"
	ErrorCommunityIDMissing      = "community ID is missing in request"
	ErrorTopicMessageTypeMissing = "topic message type is missing from params"
	ErrorTopicMissing            = "topic is missing from request"
	ErrorTopicInvalid            = "invalid format of topic"
	ErrorSenderIDMissing         = "sender ID is missing"       // New constant for missing sender ID
	ErrorDeviceIDMissing         = "device ID is missing"       // New constant for missing device ID
	ErrorConversationIDMissing   = "conversation ID is missing" // New constant for missing conversation ID
	ErrorTimestampMissing        = "timestamp is missing"       // New constant for missing timestamp
)
