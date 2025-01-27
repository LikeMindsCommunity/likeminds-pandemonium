package models

type MessageAttachment struct {
	ID           int      `gorm:"primaryKey;column:id" json:"id"`
	FileURL      string   `gorm:"column:file_url" json:"file_url"`
	Type         string   `gorm:"column:type;size:50;not null" json:"type"`
	AnswerID     int      `gorm:"column:answer_id;not null" json:"answer_id"`
	CreatedAt    int64    `gorm:"column:created_at;not null" json:"created_at"`
	LocationLat  *float64 `gorm:"column:location_lat" json:"location_lat,omitempty"`
	LocationLong *float64 `gorm:"column:location_long" json:"location_long,omitempty"`
	LocationName *string  `gorm:"column:location_name" json:"location_name,omitempty"`
	Index        *int     `gorm:"column:index" json:"index,omitempty"`
	Dimensions   *string  `gorm:"column:dimensions" json:"dimensions,omitempty"`
	Height       *int     `gorm:"column:height" json:"height,omitempty"`
	Width        *int     `gorm:"column:width" json:"width,omitempty"`
	ThumbnailURL *string  `gorm:"column:thumbnail_url" json:"thumbnail_url,omitempty"`
	Meta         *string  `gorm:"column:meta" json:"meta,omitempty"`
	Name         *string  `gorm:"column:name;size:200" json:"name,omitempty"`
}

func (messageAttachment *MessageAttachment) TableName() string {
	return "togther_answerattachment"
}
