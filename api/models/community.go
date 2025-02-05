package models

type Community struct {
	ID                    int64   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name                  string  `gorm:"column:name;type:text;not null" json:"name"`
	About                 *string `gorm:"column:about;type:text" json:"about,omitempty"`
	Purpose               *string `gorm:"column:purpose;type:varchar(2048)" json:"purpose,omitempty"`
	Location              *string `gorm:"column:location;type:varchar(200)" json:"location,omitempty"`
	ImageURL              *string `gorm:"column:image_url;type:varchar(100)" json:"image_url,omitempty"`
	MembersCount          int     `gorm:"column:members_count;not null" json:"members_count"`
	ActiveSince           string  `gorm:"column:active_since;type:date;not null" json:"active_since"`
	WhatsappGroupLink     *string `gorm:"column:whatsapp_group_link;type:varchar(400)" json:"whatsapp_group_link,omitempty"`
	UpdatedAt             int64   `gorm:"column:updated_at;not null" json:"updated_at"`
	CreatedAt             int64   `gorm:"column:created_at;not null" json:"created_at"`
	PurposeCollabCard     *int    `gorm:"column:purpose_collabcard" json:"purpose_collabcard,omitempty"`
	HideCommunity         *string `gorm:"column:hide_community;type:bpchar(1);default:0" json:"hide_community,omitempty"`
	IntroductionText      *string `gorm:"column:introduction_text;type:varchar(2048)" json:"introduction_text,omitempty"`
	ImageLink             *string `gorm:"column:image_link;type:varchar(500)" json:"image_link,omitempty"`
	IntroductionTextState int     `gorm:"column:introduction_text_state;not null" json:"introduction_text_state"`
	Thumbnail             *string `gorm:"column:thumbnail;type:varchar(500)" json:"thumbnail,omitempty"`
	Type                  *int    `gorm:"column:type" json:"type,omitempty"`
	SubType               *int    `gorm:"column:sub_type" json:"sub_type,omitempty"`
	AttributeType         int     `gorm:"column:attribute_type;not null" json:"attribute_type"`
	ImageLinkRound        *string `gorm:"column:image_link_round;type:text" json:"image_link_round,omitempty"`
	AutoApproval          bool    `gorm:"column:auto_approval;not null" json:"auto_approval"`
	GracePeriod           int64   `gorm:"column:grace_period;not null" json:"grace_period"`
	IsDiscoverable        bool    `gorm:"column:is_discoverable;not null" json:"is_discoverable"`
	IsPaid                bool    `gorm:"column:is_paid;not null" json:"is_paid"`
	WebsiteURL            *string `gorm:"column:website_url;type:text" json:"website_url,omitempty"`
	CommunityCategory     *string `gorm:"column:community_category;type:text" json:"community_category,omitempty"`
	ReferralEnabled       bool    `gorm:"column:referral_enabled;not null" json:"referral_enabled"`
	DashboardLink         *string `gorm:"column:dashboard_link;type:text" json:"dashboard_link,omitempty"`
	FeeEvent              int     `gorm:"column:fee_event;not null" json:"fee_event"`
	FeeMembership         int     `gorm:"column:fee_membership;not null" json:"fee_membership"`
	FeePaymentPages       int     `gorm:"column:fee_payment_pages;not null" json:"fee_payment_pages"`
	BrandColor            *string `gorm:"column:brand_color;type:text" json:"brand_color,omitempty"`
	LikemindsPlan         *string `gorm:"column:likeminds_plan;type:text" json:"likeminds_plan,omitempty"`
	Branding              *string `gorm:"column:branding;type:text" json:"branding,omitempty"`
	IsWhiteLabel          bool    `gorm:"column:is_whitelabel;not null" json:"is_whitelabel"`
	WhiteLabelInfo        *string `gorm:"column:whitelabel_info;type:text" json:"whitelabel_info,omitempty"`
	HideDMTab             bool    `gorm:"column:hide_dm_tab;not null" json:"hide_dm_tab"`
	IsFreemiumCommunity   bool    `gorm:"column:is_freemium_community;not null" json:"is_freemium_community"`
}

func (community *Community) TableName() string {
	return "togther_community"
}
