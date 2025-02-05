package models

type MessagePoll struct {
	ID             int    `gorm:"primaryKey;column:id" json:"id"`
	Text           string `gorm:"not null;column:text" json:"text"`
	CreatedAt      int64  `gorm:"not null;column:created_at" json:"created_at"`
	UpdatedAt      int64  `gorm:"not null;column:updated_at" json:"updated_at"`
	ConversationID int    `gorm:"not null;column:conversation_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"conversation_id"`
	UserID         int    `gorm:"not null;column:user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user_id"`
}

func (messagePoll *MessagePoll) TableName() string {
	return "togther_conversationpolls"
}
