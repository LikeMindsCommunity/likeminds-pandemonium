package models

type UserInfo struct {
	ID               int64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name             string   `gorm:"column:name;type:varchar(200);not null" json:"name"`
	Email            string   `gorm:"column:email;type:varchar(200);not null" json:"email"`
	City             *string  `gorm:"column:city;type:varchar(100)" json:"city,omitempty"`
	ContactNumber    *int64   `gorm:"column:contact_number" json:"contact_number,omitempty"`
	Gender           *int     `gorm:"column:gender" json:"gender,omitempty"`
	ImageURL         *string  `gorm:"column:image_url;type:varchar(500)" json:"image_url,omitempty"`
	Interests        *string  `gorm:"column:interests;type:varchar(400)" json:"interests,omitempty"`
	About            *string  `gorm:"column:about;type:varchar(400)" json:"about,omitempty"`
	FBLink           *string  `gorm:"column:fb_link;type:varchar(400)" json:"fb_link,omitempty"`
	LinkedInLink     *string  `gorm:"column:linkedin_link;type:varchar(400)" json:"linkedin_link,omitempty"`
	Headline         *string  `gorm:"column:headline;type:varchar(100)" json:"headline,omitempty"`
	UserID           int      `gorm:"column:user_id_id;not null;uniqueIndex" json:"user_id"`
	ImageFile        *string  `gorm:"column:image_file;type:varchar(100)" json:"image_file,omitempty"`
	FCMToken         *string  `gorm:"column:fcm_token;type:varchar(1024)" json:"fcm_token,omitempty"`
	LoginJSON        *string  `gorm:"column:login_json;type:text" json:"login_json,omitempty"`
	LoginType        *string  `gorm:"column:login_type;type:varchar(50)" json:"login_type,omitempty"`
	Latitude         *float64 `gorm:"column:latitude" json:"latitude,omitempty"`
	Longitude        *float64 `gorm:"column:longitude" json:"longitude,omitempty"`
	Address          *string  `gorm:"column:address;type:varchar(1024)" json:"address,omitempty"`
	MobileOS         *string  `gorm:"column:mobile_os;type:varchar(200)" json:"mobile_os,omitempty"`
	SecondaryEmail   *string  `gorm:"column:secondary_email;type:varchar(200)" json:"secondary_email,omitempty"`
	CreatedAt        int64    `gorm:"column:created_at;not null" json:"created_at"`
	VersionCode      *int     `gorm:"column:version_code" json:"version_code,omitempty"`
	ImageLink        *string  `gorm:"column:image_link;type:varchar(5000)" json:"image_link,omitempty"`
	AppleID          *string  `gorm:"column:apple_id;type:varchar(100)" json:"apple_id,omitempty"`
	HasTags          bool     `gorm:"column:has_tags;not null" json:"has_tags"`
	UpdatedAt        int64    `gorm:"column:updated_at;not null" json:"updated_at"`
	UserUniqueID     *string  `gorm:"column:user_unique_id;type:varchar(255);uniqueIndex" json:"user_unique_id,omitempty"`
	IsBot            bool     `gorm:"column:is_bot;not null" json:"is_bot"`
	OrganisationName *string  `gorm:"column:organisation_name;type:varchar(255)" json:"organisation_name,omitempty"`
	IsGuest          bool     `gorm:"column:is_guest;not null" json:"is_guest"`
	LastActive       int64    `gorm:"column:last_active;not null" json:"last_active"`
	// Roles            *[]sql.NullString `gorm:"column:roles;type:text[]" json:"roles"`
	MetaInfo interface{} `gorm:"column:meta_info;type:jsonb;not null" json:"meta_info"`
}

func (userInfo *UserInfo) TableName() string {
	return "togther_userinfo"
}
