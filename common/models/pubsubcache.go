package models

type ConversationMetaCache struct {
	DeliveryCount int    `json:"delivery_count"`
	SenderUUID    string `json:"sender_uuid"`
}
