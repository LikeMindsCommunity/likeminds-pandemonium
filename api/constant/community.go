package constant

const (
	CommunityConfigurationMediaLimitsEnum         = "media_limits"
	CommunityConfigurationFeedMetaData            = "feed_metadata"
	CommunityConfigurationProfileMetaData         = "profile_metadata"
	CommunityConfigurationNSFWFiltering           = "nsfw_filtering"
	CommunityConfigurationWidgetMetadata          = "widgets_metadata"
	CommunityConfigurationGuestFlowMetaData       = "guest_flow_metadata"
	CommunityConfigurationFeedSettings            = "feed_settings"
	CommunityConfigurationPersonalisedFeedWeights = "personalised_feed_weights"
	CommunityConfigurationChatbot                 = "chatbot"
	CommunityConfigurationChatPoll                = "chat_poll"
)

type CommunityConfiguration struct {
	Type        string
	Description string
	Value       CommunityConfigurationValue
}

type CommunityConfigurationValue struct {
	AllowOveride        bool
	PollType            string
	NoPollExpiry        bool
	AllowVoteChange     bool
	MultipleSelectState string
	MultipleSelectNo    int
	IsAnonymous         bool
	AllowAddOption      bool
}

var CommunityConfigurationChatPollDefault = &CommunityConfiguration{
	Type:        CommunityConfigurationChatPoll,
	Description: "Chat poll configurations for the community.",
	Value: CommunityConfigurationValue{
		AllowOveride:        true,
		PollType:            "instant", // values - "instant" | "deferred" | "open"
		NoPollExpiry:        false,
		AllowVoteChange:     false,
		MultipleSelectState: "exactly", // values - "exactly" | "at_max" | "at_most" | "at_least"
		MultipleSelectNo:    1,
		IsAnonymous:         false,
		AllowAddOption:      false,
	},
}
