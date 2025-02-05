package constant

const (
	CommunityConfigurationMediaLimitsEnum             = "media_limits"
	CommunityConfigurationFeedMetaDataEnum            = "feed_metadata"
	CommunityConfigurationProfileMetaDataEnum         = "profile_metadata"
	CommunityConfigurationNSFWFilteringEnum           = "nsfw_filtering"
	CommunityConfigurationWidgetMetadataEnum          = "widgets_metadata"
	CommunityConfigurationGuestFlowMetaDataEnum       = "guest_flow_metadata"
	CommunityConfigurationFeedSettingsEnum            = "feed_settings"
	CommunityConfigurationPersonalisedFeedWeightsEnum = "personalised_feed_weights"
	CommunityConfigurationChatbotEnum                 = "chatbot"
	CommunityConfigurationChatPollEnum                = "chat_poll"
)

type CommunityConfiguration struct {
	Type        string
	Description string
	Value       interface{}
}

type CommunityConfigurationChatPollValue struct {
	AllowOveride        bool
	PollType            string
	NoPollExpiry        bool
	AllowVoteChange     bool
	MultipleSelectState string
	MultipleSelectNo    int
	IsAnonymous         bool
	AllowAddOption      bool
}

type CommunityConfigurationWidgetMetadataValue struct {
	Message bool
}

var CommunityConfigurationWidgetMetadataDefault = &CommunityConfiguration{
	Type:        CommunityConfigurationWidgetMetadataEnum,
	Description: "Widget metadata configurations for the community",
	Value: CommunityConfigurationWidgetMetadataValue{
		Message: false,
	},
}
