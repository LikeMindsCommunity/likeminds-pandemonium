package models

type SDKClient struct {
	ID                int     `gorm:"primaryKey;autoIncrement" json:"id"`
	APIKey            *string `gorm:"size:64;unique;index" json:"api_key,omitempty"`
	CreatedAt         int64   `gorm:"not null" json:"created_at"`
	UpdatedAt         int64   `gorm:"not null" json:"updated_at"`
	CommunityID       int     `gorm:"not null;index" json:"community_id"`
	ProjectCreatorID  *int    `gorm:"index" json:"project_creator_id,omitempty"`
	IsDeleted         bool    `gorm:"not null" json:"is_deleted"`
	FirebaseServerKey *string `json:"firebase_server_key,omitempty"`
	IsJoinFormEnabled bool    `gorm:"not null" json:"is_join_form_enabled"`
	// GCPServiceAccountFile *map[string]any `gorm:"type:jsonb" json:"gcp_service_account_file,omitempty"`
}

// TableName overrides the default table name for the model
func (sdkClient *SDKClient) TableName() string {
	return "collabmates_api_sdkclient"
}
