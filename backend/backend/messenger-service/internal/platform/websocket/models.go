package websocket

import "encoding/json"

type WSEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// * Client -> Server events
const (
	EventSendMessage   = "SEND_MESSAGE"
	EventEditMessage   = "EDIT_MESSAGE"
	EventDeleteMessage = "DELETE_MESSAGE"
	EventReadMessage   = "READ_MESSAGE"
	EventTyping        = "TYPING"
	EventForward       = "FORWARD_MESSAGE"
	EventReaction      = "REACTION"
)

// * Server -> Client events
const (
	EventNewMessage         = "NEW_MESSAGE"
	EventMessageEdited      = "MESSAGE_EDITED"
	EventMessageDeleted     = "MESSAGE_DELETED"
	EventMessagesRead       = "MESSAGES_READ"
	EventTypingNotify       = "TYPING"
	EventUserOnline         = "USER_ONLINE"
	EventUserOffline        = "USER_OFFLINE"
	EventUserBlocked        = "USER_BLOCKED"   // -> the blocked user (so UI locks the chat live)
	EventUserUnblocked      = "USER_UNBLOCKED" // -> the unblocked user
	EventMemberAdded        = "MEMBER_ADDED"
	EventMemberRemoved      = "MEMBER_REMOVED"
	EventMemberLeft         = "MEMBER_LEFT"
	EventChatUpdated        = "CHAT_UPDATED"
	EventChatDeleted        = "CHAT_DELETED" // payload scope: "everyone" (all members) | "me" (caller devices)
	EventMessagePinned      = "MESSAGE_PINNED"
	EventMessageUnpinned    = "MESSAGE_UNPINNED"
	EventReactionAdded      = "REACTION_ADDED"
	EventReactionRemoved    = "REACTION_REMOVED"
	EventError              = "ERROR"
	EventPreferencesUpdated = "PREFERENCES_UPDATED"
	EventChannelNewPost      = "CHANNEL_NEW_POST"
	EventChannelNewComment   = "CHANNEL_NEW_COMMENT"
	EventChannelSubscribed   = "CHANNEL_SUBSCRIBED"
	EventChannelUnsubscribed = "CHANNEL_UNSUBSCRIBED"

	// Private channel join requests / invites
	EventChannelJoinRequest  = "CHANNEL_JOIN_REQUEST"  // -> channel admins
	EventChannelJoinApproved = "CHANNEL_JOIN_APPROVED" // -> requesting user
	EventChannelJoinRejected = "CHANNEL_JOIN_REJECTED" // -> requesting user
	EventChannelInvite       = "CHANNEL_INVITE"        // -> invited user
	EventChannelInviteAccepted = "CHANNEL_INVITE_ACCEPTED" // -> channel admins

	// Friends feature
	EventFriendRequestReceived = "FRIEND_REQUEST_RECEIVED"
	EventFriendRequestAccepted = "FRIEND_REQUEST_ACCEPTED"
	EventFriendRequestRejected = "FRIEND_REQUEST_REJECTED"
	EventFriendRequestCancelled = "FRIEND_REQUEST_CANCELLED"
	EventFriendRemoved         = "FRIEND_REMOVED"
)
