package constant

// constants for user states in community
const (
	UserStateGuest int = iota
	UserStateAdmin
	UserStateTempAdmin
	UserStatePendingMember
	UserStateMember
	UserStateDeclinedMember
	UserStateUnknownNominatedPromoter
	UserStateKnownNominatedPromoter
	UserStateInterestedMember
	UserStateProfileUnavailable
)

// enums for user states in community
const (
	UserStateGuestEnum                    = "guest"
	UserStateAdminEnum                    = "admin"
	UserStateTempAdminEnum                = "temp_admin"
	UserStatePendingMemberEnum            = "pending_member"
	UserStateMemberEnum                   = "member"
	UserStateDeclinedMemberEnum           = "declined_member"
	UserStateUnknownNominatedPromoterEnum = "unknown_nominated_promoter"
	UserStateKnownNominatedPromoterEnum   = "known_nominated_promoter"
	UserStateInterestedMemberEnum         = "interested_member"
	UserStateProfileUnavailableEnum       = "profile_unavailable"
)

const (
	MemberRightCreateRooms int = iota
	MemberRightCreatePoll
	MemberRightCreateEvent
	MemberRightRespondInRoom
	MemberRightInvitePrivateLink
	MemberRightAutoApprove
	MemberRightCreateSecretRoom
	ManagerRightEnableDirectMessages
	MemberRightEnableMembersCanDM
	MemberRightCreatePosts
	MemberRightCommentAndReplyOnPosts
	MemberRightCreateFeedPoll
)

const (
	MemberRightCreateRoomsEnum            = "create_rooms"
	MemberRightCreatePollEnum             = "create_poll"
	MemberRightCreateEventEnum            = "create_event"
	MemberRightRespondInRoomEnum          = "respond_in_room"
	MemberRightInvitePrivateLinkEnum      = "invite_private_link"
	MemberRightAutoApproveEnum            = "auto_approve"
	MemberRightCreateSecretRoomEnum       = "create_secret_room"
	ManagerRightEnableDirectMessagesEnum  = "enable_direct_messages"
	MemberRightEnableMembersCanDMEnum     = "enable_members_can_dm"
	MemberRightCreatePostsEnum            = "create_posts"
	MemberRightCommentAndReplyOnPostsEnum = "comment_and_reply_on_posts"
	MemberRightCreateFeedPollEnum         = "create_feed_poll"
)
