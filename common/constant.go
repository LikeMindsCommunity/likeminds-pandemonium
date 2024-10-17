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
	// WriteWait Max wait time when writing message to peer
	WriteWait = 10 * time.Second
	// PongWait Max time till next pong from peer
	PongWait = 60 * time.Second
	// PingPeriod should be less than PongWait
	PingPeriod                  = (PongWait * 9) / 10
	WsServerKey                 = "ws_server"
	TopicMessageTypeSentDR      = "sent_dr"
	TopicMessageTypeDeliveredDR = "delivered_dr"
	DRChatroomPrefix            = "dr_chatroom_%v"
	DRConversationPrefix        = "dr_conversation_%v"
	DRConversationMetaPrefix    = "dr_conversation_meta"
	DRUserPrefix                = DRUser + "%v"
	TopicTypeChatroomDynamic    = "chatroom:%v"
	RawData                     = "raw_data"
	DRUser                      = "dr_user_"
	DeliveryCount               = "delivery_count"
	SenderUUID                  = "sender_uuid"
	ParamChatroomID             = "chatroom_id"
	ParamConversationIDs        = "conversation_ids"
	TopicTypeCommunityDynamic   = "community:%v"
)

const (
	WsConnectionEstablished = "connected to websocket server"
	PingReceivedClient      = "received ping from client"
	PongSendClient          = "sending pong to client"
	ReceivedMessageClientWs = "received message from client having message type as %v"
	ReceivedMessageRedisWs  = "received message from redis"
	ConnectionClosed        = "connection closed"
)

const (
	ErrorFailedUpgrader        = "failed to upgrade connection: %v"
	ErrorPublishRedis          = "failed to publish message to topic %s: %v"
	ErrorSubscribeRedis        = "failed to subscribe to topic %s: %v"
	ErrorReadClientWs          = "error reading message from client: %v"
	ErrorReadDeadlineWs        = "error while setting ReadDeadline on websocket:"
	ErrorWriteDeadlineWs       = "error while setting WriteDeadline on websocket:"
	ErrorUnableToCloseWs       = "unable to close ws error:"
	ErrorUnableToCloseRedis    = "unable to close redis error:"
	ErrorUnmarshalErrorJson    = "unmarshal error: %v"
	ErrorUnableToWriteWs       = "unable to write message in websocket: %v"
	ErrorWriterOpenWs          = "unable to open websocket writer: %v"
	ErrorWriterCloseWs         = "unable to close websocket writer: %v"
	ErrorPongSentClient        = "error sending pong to client: %v"
	ErrorFailedCacheSaveRedis  = "failed to save cache: %v"
	ErrorFailedExpSaveRedis    = "failed to save cache exp: %v"
	ErrorFailedCacheFetchRedis = "failed to fetch from redis cache: %v"
	ErrorMarshalErrorJson      = "marshal error: %v"
	ErrorInvalidJSONFormat     = "invalid JSON format: %v"
)

const (
	ErrorUserUUIDMissing         = "user UUID is missing in headers"
	ErrorChatroomIDMissing       = "chatroom ID is missing in request"
	ErrorCommunityIDMissing      = "community ID is missing in request"
	ErrorTopicMessageTypeMissing = "topic message type is missing from params"
	ErrorTopicMissing            = "topic is missing from request"
	ErrorTopicInvalid            = "invalid format of topic"
	ErrorConversationIDsMissing  = "conversation IDs is missing in request"
	ErrorConversationIDsInvalid  = "invalid conversation IDs in request"
)

const (
	DeliveryReportTTL = 7 * 24 * time.Hour // 7 days TTL for delivery reports
)
