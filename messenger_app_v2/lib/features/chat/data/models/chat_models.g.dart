// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'chat_models.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

MessagePreview _$MessagePreviewFromJson(Map<String, dynamic> json) =>
    MessagePreview(
      messageId: (json['message_id'] as num).toInt(),
      senderId: json['sender_id'] as String?,
      senderName: json['sender_name'] as String?,
      content: json['content'] as String,
      type: json['type'] as String,
      deleted: json['deleted'] as bool? ?? false,
    );

Map<String, dynamic> _$MessagePreviewToJson(MessagePreview instance) =>
    <String, dynamic>{
      'message_id': instance.messageId,
      'sender_id': instance.senderId,
      'sender_name': instance.senderName,
      'content': instance.content,
      'type': instance.type,
      'deleted': instance.deleted,
    };

Chat _$ChatFromJson(Map<String, dynamic> json) => Chat(
  id: json['id'] as String,
  type: json['type'] as String,
  name: json['name'] as String?,
  createdAt: json['created_at'] as String,
);

Map<String, dynamic> _$ChatToJson(Chat instance) => <String, dynamic>{
  'id': instance.id,
  'type': instance.type,
  'name': instance.name,
  'created_at': instance.createdAt,
};

ChatListItem _$ChatListItemFromJson(Map<String, dynamic> json) => ChatListItem(
  id: json['id'] as String,
  type: json['type'] as String,
  name: json['name'] as String?,
  lastMessageContent: json['last_message_content'] as String,
  lastMessageTime: json['last_message_time'] as String,
  unreadCount: (json['unread_count'] as num).toInt(),
  interlocutorId: json['interlocutor_id'] as String?,
  interlocutorName: json['interlocutor_name'] as String?,
  interlocutorAvatar: json['interlocutor_avatar'] as String?,
  interlocutorAvatarMediaId: json['interlocutor_avatar_media_id'] as String?,
  avatarMediaId: json['avatar_media_id'] as String?,
  isMuted: json['is_muted'] as bool? ?? false,
  isPinned: json['is_pinned'] as bool? ?? false,
  memberCount: (json['member_count'] as num?)?.toInt(),
  onlineCount: (json['online_count'] as num?)?.toInt(),
  memberRole: json['member_role'] as String?,
  lastMessageType: json['last_message_type'] as String?,
  blockedByMe: json['blocked_by_me'] as bool? ?? false,
  blockedMe: json['blocked_me'] as bool? ?? false,
);

Map<String, dynamic> _$ChatListItemToJson(ChatListItem instance) =>
    <String, dynamic>{
      'id': instance.id,
      'type': instance.type,
      'name': instance.name,
      'last_message_content': instance.lastMessageContent,
      'last_message_time': instance.lastMessageTime,
      'unread_count': instance.unreadCount,
      'interlocutor_id': instance.interlocutorId,
      'interlocutor_name': instance.interlocutorName,
      'interlocutor_avatar': instance.interlocutorAvatar,
      'interlocutor_avatar_media_id': instance.interlocutorAvatarMediaId,
      'avatar_media_id': instance.avatarMediaId,
      'is_muted': instance.isMuted,
      'is_pinned': instance.isPinned,
      'member_count': instance.memberCount,
      'online_count': instance.onlineCount,
      'member_role': instance.memberRole,
      'last_message_type': instance.lastMessageType,
      'blocked_by_me': instance.blockedByMe,
      'blocked_me': instance.blockedMe,
    };

Message _$MessageFromJson(Map<String, dynamic> json) => Message(
  id: (json['id'] as num).toInt(),
  chatId: json['chat_id'] as String,
  senderId: json['sender_id'] as String?,
  content: json['content'] as String,
  type: json['type'] as String,
  status: json['status'] as String,
  createdAt: json['created_at'] as String,
  clientSideId: json['client_side_id'] as String?,
  editedAt: json['edited_at'] as String?,
  attachmentMediaId: json['attachment_media_id'] as String?,
  replyToMessageId: (json['reply_to_message_id'] as num?)?.toInt(),
  forwardedFromChat: json['forwarded_from_chat'] as String?,
  forwardedFromMsg: (json['forwarded_from_msg'] as num?)?.toInt(),
  deletedAt: json['deleted_at'] as String?,
  commentCount: (json['comment_count'] as num?)?.toInt(),
  replyPreview: json['reply_preview'] == null
      ? null
      : MessagePreview.fromJson(json['reply_preview'] as Map<String, dynamic>),
  forwardedFrom: json['forwarded_from'] == null
      ? null
      : MessagePreview.fromJson(json['forwarded_from'] as Map<String, dynamic>),
);

Map<String, dynamic> _$MessageToJson(Message instance) => <String, dynamic>{
  'id': instance.id,
  'chat_id': instance.chatId,
  'sender_id': instance.senderId,
  'content': instance.content,
  'type': instance.type,
  'status': instance.status,
  'created_at': instance.createdAt,
  'client_side_id': instance.clientSideId,
  'edited_at': instance.editedAt,
  'attachment_media_id': instance.attachmentMediaId,
  'reply_to_message_id': instance.replyToMessageId,
  'forwarded_from_chat': instance.forwardedFromChat,
  'forwarded_from_msg': instance.forwardedFromMsg,
  'deleted_at': instance.deletedAt,
  'comment_count': instance.commentCount,
  'reply_preview': instance.replyPreview,
  'forwarded_from': instance.forwardedFrom,
};

SendMessageRequest _$SendMessageRequestFromJson(Map<String, dynamic> json) =>
    SendMessageRequest(
      chatId: json['chat_id'] as String,
      content: json['content'] as String,
      type: json['type'] as String,
      clientSideId: json['client_side_id'] as String,
      attachmentMediaId: json['attachment_media_id'] as String?,
      replyToMessageId: (json['reply_to_message_id'] as num?)?.toInt(),
    );

Map<String, dynamic> _$SendMessageRequestToJson(SendMessageRequest instance) =>
    <String, dynamic>{
      'chat_id': instance.chatId,
      'content': instance.content,
      'type': instance.type,
      'client_side_id': instance.clientSideId,
      'attachment_media_id': instance.attachmentMediaId,
      'reply_to_message_id': instance.replyToMessageId,
    };

UserStatus _$UserStatusFromJson(Map<String, dynamic> json) => UserStatus(
  userId: json['user_id'] as String,
  isOnline: json['is_online'] as bool,
  lastSeenAt: json['last_seen_at'] as String?,
);

Map<String, dynamic> _$UserStatusToJson(UserStatus instance) =>
    <String, dynamic>{
      'user_id': instance.userId,
      'is_online': instance.isOnline,
      'last_seen_at': instance.lastSeenAt,
    };

MessageReader _$MessageReaderFromJson(Map<String, dynamic> json) =>
    MessageReader(
      userId: json['user_id'] as String,
      username: json['username'] as String,
      displayName: json['display_name'] as String,
      readAt: json['read_at'] as String,
    );

Map<String, dynamic> _$MessageReaderToJson(MessageReader instance) =>
    <String, dynamic>{
      'user_id': instance.userId,
      'username': instance.username,
      'display_name': instance.displayName,
      'read_at': instance.readAt,
    };

CreateChatRequest _$CreateChatRequestFromJson(Map<String, dynamic> json) =>
    CreateChatRequest(
      type: json['type'] as String,
      targetId: json['target_id'] as String?,
      name: json['name'] as String?,
      description: json['description'] as String?,
      memberIds: (json['member_ids'] as List<dynamic>?)
          ?.map((e) => e as String)
          .toList(),
      visibility: json['visibility'] as String?,
      username: json['username'] as String?,
    );

Map<String, dynamic> _$CreateChatRequestToJson(CreateChatRequest instance) =>
    <String, dynamic>{
      'type': instance.type,
      'target_id': instance.targetId,
      'name': instance.name,
      'description': instance.description,
      'member_ids': instance.memberIds,
      'visibility': instance.visibility,
      'username': instance.username,
    };

ChatMember _$ChatMemberFromJson(Map<String, dynamic> json) => ChatMember(
  userId: json['user_id'] as String,
  username: json['username'] as String,
  displayName: json['display_name'] as String,
  avatarMediaId: json['avatar_media_id'] as String?,
  role: json['role'] as String,
  joinedAt: json['joined_at'] as String?,
);

Map<String, dynamic> _$ChatMemberToJson(ChatMember instance) =>
    <String, dynamic>{
      'user_id': instance.userId,
      'username': instance.username,
      'display_name': instance.displayName,
      'avatar_media_id': instance.avatarMediaId,
      'role': instance.role,
      'joined_at': instance.joinedAt,
    };

ChatInfo _$ChatInfoFromJson(Map<String, dynamic> json) => ChatInfo(
  id: json['id'] as String,
  type: json['type'] as String,
  name: json['name'] as String?,
  description: json['description'] as String?,
  avatarMediaId: json['avatar_media_id'] as String?,
  members:
      (json['members'] as List<dynamic>?)
          ?.map((e) => ChatMember.fromJson(e as Map<String, dynamic>))
          .toList() ??
      [],
  memberCount: (json['member_count'] as num?)?.toInt() ?? 0,
  username: json['username'] as String?,
  isVerified: json['is_verified'] as bool? ?? false,
  blockedByMe: json['blocked_by_me'] as bool? ?? false,
  blockedMe: json['blocked_me'] as bool? ?? false,
);

Map<String, dynamic> _$ChatInfoToJson(ChatInfo instance) => <String, dynamic>{
  'id': instance.id,
  'type': instance.type,
  'name': instance.name,
  'description': instance.description,
  'avatar_media_id': instance.avatarMediaId,
  'members': instance.members,
  'member_count': instance.memberCount,
  'username': instance.username,
  'is_verified': instance.isVerified,
  'blocked_by_me': instance.blockedByMe,
  'blocked_me': instance.blockedMe,
};
