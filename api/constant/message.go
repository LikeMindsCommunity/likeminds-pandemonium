package constant

const (
	RegexTagEveryone     = `@everyone`
	RegexTagParticipants = `@participants`
)

const (
	MessageStateAnswer int = iota
	MessageStateMessageHeader
	MessageStateMessageFollow
	MessageStateMessageUnfollow
	MessageStateMessageCreator
	MessageStateMessageCommunityEdit
	MessageStateMessageGuest
	MessageStateMessageAddParticipant
	MessageStateMessageLeaveChatroom
	MessageStateMessageRemovedFromChatroom
	MessageStateMessagePoll
	MessageStateMessageAddAllMembers
	MessageStateChatroomTopic
	MessageStateMessageDirectMessageMemberRemovedOrLeft
	MessageStateMessageDirectMessageCMRemoved
	MessageStateMessageDirectMessageMemberBecomesCMDisableChat
	MessageStateMessageDirectMessageCMBecomesMemberEnableChat
	MessageStateMessageDirectMessageMemberBecomesCMEnableChat
	MessageStateMessageEvent
	MessageStateMessageDirectMessageBlockMemberDisableChat
	MessageStateMessageDirectMessageUnblockMemberEnableChat
	MessageStateChatroomDelete
)

const (
	MessagePollAnswerText = "Be the first to vote"
)

const (
	OnlineEventLinkEnableBeforeMinutes = 15
)

const (
	MessagePollTypeInstant int = iota
	MessagePollTypeDeferred
	MessagePollTypeOpen
)

const (
	MessagePollTypeInstantEnum  = "instant"
	MessagePollTypeDeferredEnum = "deferred"
	MessagePollTypeOpenEnum     = "open"
)

func GetMessagePollTypeFromEnum(pollType string) int {
	switch pollType {
	case MessagePollTypeInstantEnum:
		return MessagePollTypeInstant
	case MessagePollTypeDeferredEnum:
		return MessagePollTypeDeferred
	case MessagePollTypeOpenEnum:
		return MessagePollTypeOpen
	}

	return -1
}

const (
	MessagePollMultiSelectStateExactly int = iota
	MessagePollMultiSelectStateAtMax
	MessagePollMultiSelectStateAtMost
	MessagePollMultiSelectStateAtLeast
)

const (
	MessagePollMultiSelectStateExactlyEnum = "exactly"
	MessagePollMultiSelectStateAtMaxEnum   = "at_max"
	MessagePollMultiSelectStateAtMostEnum  = "at_most"
	MessagePollMultiSelectStateAtLeastEnum = "at_least"
)

func GetMessagePollMultiSelectStateFromEnum(multiSelectState string) int {
	switch multiSelectState {
	case MessagePollMultiSelectStateExactlyEnum:
		return MessagePollMultiSelectStateExactly
	case MessagePollMultiSelectStateAtMaxEnum:
		return MessagePollMultiSelectStateAtMax
	case MessagePollMultiSelectStateAtMostEnum:
		return MessagePollMultiSelectStateAtMost
	case MessagePollMultiSelectStateAtLeastEnum:
		return MessagePollMultiSelectStateAtLeast
	}

	return -1
}
