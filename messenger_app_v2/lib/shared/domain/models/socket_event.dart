import 'package:json_annotation/json_annotation.dart';

part 'socket_event.g.dart';

@JsonSerializable()
class SocketEvent {
  final String type;
  final Map<String, dynamic>? payload; // Payload arrives as JSON object

  SocketEvent({required this.type, this.payload});

  factory SocketEvent.fromJson(Map<String, dynamic> json) => _$SocketEventFromJson(json);
}

// Event type constants
class SocketEventType {
  static const sendMessage = 'SEND_MESSAGE';
  static const readMessage = 'READ_MESSAGE';
  static const typing = 'TYPING';
  static const newMessage = 'NEW_MESSAGE'; // Arrives from server
  static const editMessage = 'EDIT_MESSAGE'; // Client sends to server
  static const deleteMessage = 'DELETE_MESSAGE'; // Client sends to server
  static const messageEdited = 'MESSAGE_EDITED'; // Arrives from server
  static const messageDeleted = 'MESSAGE_DELETED'; // Arrives from server
  static const messagesRead = 'MESSAGES_READ'; // Arrives from server: {chat_id, last_read_id, by_user}
  static const userOnline = 'USER_ONLINE'; // Arrives from server: {user_id}
  static const userOffline = 'USER_OFFLINE'; // Arrives from server: {user_id}
  static const error = 'ERROR';

  // Group/channel membership events
  static const chatUpdated = 'CHAT_UPDATED'; // {chat: Chat}
  static const memberAdded = 'MEMBER_ADDED'; // {chat_id, member: ChatMember}
  static const memberRemoved = 'MEMBER_REMOVED'; // {chat_id, user_id}
  static const memberLeft = 'MEMBER_LEFT'; // {chat_id, user_id}

  // Message pin events
  static const messagePinned = 'MESSAGE_PINNED'; // {chat_id, message_id}
  static const messageUnpinned = 'MESSAGE_UNPINNED'; // {chat_id, message_id}

  // Reaction events
  static const reactionAdded = 'REACTION_ADDED'; // {message_id, user_id, emoji}
  static const reactionRemoved = 'REACTION_REMOVED'; // {message_id, user_id, emoji}

  // Preferences sync across sessions
  static const preferencesUpdated = 'PREFERENCES_UPDATED'; // {preferences}

  // Channel comment events
  static const channelNewComment = 'CHANNEL_NEW_COMMENT'; // payload = ChannelComment

  // Transaction events (in-chat transfer / donation)
  static const sudaSent          = 'SUDA_SENT';       // server → sender
  static const sudaReceived      = 'SUDA_RECEIVED';    // server → recipient
  static const donationSent      = 'DONATION_SENT';
  static const donationReceived  = 'DONATION_RECEIVED';

  // Chat deletion event — sent to all participants (scope=everyone) or only to the deleting user (scope=me).
  static const chatDeleted = 'CHAT_DELETED'; // payload: {chat_id, scope}

  // Block/unblock events (§4.4) — sent to the affected user in real time.
  static const userBlocked   = 'USER_BLOCKED';   // payload: {by_user: uuid}
  static const userUnblocked = 'USER_UNBLOCKED'; // payload: {by_user: uuid}

  // Friends events
  static const friendRequestReceived = 'FRIEND_REQUEST_RECEIVED';
  static const friendRequestAccepted = 'FRIEND_REQUEST_ACCEPTED';
  static const friendRequestRejected = 'FRIEND_REQUEST_REJECTED';
  static const friendRequestCancelled = 'FRIEND_REQUEST_CANCELLED';
  static const friendRemoved = 'FRIEND_REMOVED';
}