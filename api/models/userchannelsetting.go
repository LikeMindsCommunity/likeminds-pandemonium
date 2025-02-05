package models

type UserChannelSetting struct {
	ID          int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	SettingType string `gorm:"column:setting_type;type:varchar(255);not null" json:"setting_type"`
	Enabled     bool   `gorm:"column:enabled;not null" json:"enabled"`
	CreatedAt   int64  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt   int64  `gorm:"column:updated_at;not null" json:"updated_at"`
	ChangedByID *int   `gorm:"column:changed_by_id" json:"changed_by_id,omitempty"`
	ChatroomID  int    `gorm:"column:chatroom_id;not null" json:"chatroom_id"`
	UserID      int    `gorm:"column:user_id;not null" json:"user_id"`
}

func (userChannelSetting *UserChannelSetting) TableName() string {
	return "togther_userchannelsettings"
}
