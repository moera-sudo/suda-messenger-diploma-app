import 'package:json_annotation/json_annotation.dart';

part 'chat_models.g.dart';

/// Hydrated preview of a replied-to or forwarded message.
/// Both reply_preview and forwarded_from use the same shape (§3.1).
@JsonSerializable()
class MessagePreview {
  @JsonKey(name: 'message_id')
  final int messageId;
  @JsonKey(name: 'sender_id')
  final String? senderId;
  @JsonKey(name: 'sender_name')
  final String? senderName;
  final String content;
  final String type;
  @JsonKey(defaultValue: false)
  final bool deleted;

  const MessagePreview({
    required this.messageId,
    this.senderId,
    this.senderName,
    required this.content,
    required this.type,
    this.deleted = false,
  });

  factory MessagePreview.fromJson(Map<String, dynamic> json) =>
      _$MessagePreviewFromJson(json);

  Map<String, dynamic> toJson() => _$MessagePreviewToJson(this);
}

/// Canonical message status constants — always uppercase, matching server values.
abstract class MessageStatus {
  static const sending = 'SENDING';
  static const sent    = 'SENT';
  static const read    = 'READ';
  static const failed  = 'FAILED';
}

@JsonSerializable()
class Chat {
  final String id;
  final String type; // DIRECT, GROUP, CHANNEL, SAVED
  final String? name;
  @JsonKey(name: 'created_at')
  final String createdAt;

  Chat({required this.id, required this.type, this.name, required this.createdAt});

  factory Chat.fromJson(Map<String, dynamic> json) => _$ChatFromJson(json);
}

@JsonSerializable()
class ChatListItem {
  final String id;
  final String type; // DIRECT, GROUP, CHANNEL, SAVED
  final String? name;
  @JsonKey(name: 'last_message_content')
  final String lastMessageContent;
  @JsonKey(name: 'last_message_time')
  final String lastMessageTime;
  @JsonKey(name: 'unread_count')
  final int unreadCount;
  @JsonKey(name: 'interlocutor_id')
  final String? interlocutorId;
  @JsonKey(name: 'interlocutor_name')
  final String? interlocutorName;
  @JsonKey(name: 'interlocutor_avatar')
  final String? interlocutorAvatar;
  // DIRECT chats: interlocutor's avatar media id (resolve via /media/{id}/url).
  @JsonKey(name: 'interlocutor_avatar_media_id')
  final String? interlocutorAvatarMediaId;
  @JsonKey(name: 'avatar_media_id')
  final String? avatarMediaId;
  @JsonKey(name: 'is_muted', defaultValue: false)
  final bool isMuted;
  @JsonKey(name: 'is_pinned', defaultValue: false)
  final bool isPinned;
  @JsonKey(name: 'member_count')
  final int? memberCount;
  @JsonKey(name: 'online_count')
  final int? onlineCount;
  @JsonKey(name: 'member_role')
  final String? memberRole; // OWNER, ADMIN, MEMBER
  // Type of last message — used for smart preview rendering (§2.3).
  @JsonKey(name: 'last_message_type')
  final String? lastMessageType;
  // Block status for DIRECT chats (always false for non-DIRECT).
  @JsonKey(name: 'blocked_by_me', defaultValue: false)
  final bool blockedByMe;
  @JsonKey(name: 'blocked_me', defaultValue: false)
  final bool blockedMe;

  ChatListItem({
    required this.id,
    required this.type,
    this.name,
    required this.lastMessageContent,
    required this.lastMessageTime,
    required this.unreadCount,
    this.interlocutorId,
    this.interlocutorName,
    this.interlocutorAvatar,
    this.interlocutorAvatarMediaId,
    this.avatarMediaId,
    this.isMuted = false,
    this.isPinned = false,
    this.memberCount,
    this.onlineCount,
    this.memberRole,
    this.lastMessageType,
    this.blockedByMe = false,
    this.blockedMe = false,
  });

  factory ChatListItem.fromJson(Map<String, dynamic> json) => _$ChatListItemFromJson(json);

  ChatListItem copyWith({
    String? id,
    String? type,
    String? name,
    String? lastMessageContent,
    String? lastMessageTime,
    int? unreadCount,
    String? interlocutorId,
    String? interlocutorName,
    String? interlocutorAvatar,
    String? interlocutorAvatarMediaId,
    String? avatarMediaId,
    bool? isMuted,
    bool? isPinned,
    int? memberCount,
    int? onlineCount,
    String? memberRole,
    String? lastMessageType,
    bool? blockedByMe,
    bool? blockedMe,
  }) {
    return ChatListItem(
      id: id ?? this.id,
      type: type ?? this.type,
      name: name ?? this.name,
      lastMessageContent: lastMessageContent ?? this.lastMessageContent,
      lastMessageTime: lastMessageTime ?? this.lastMessageTime,
      unreadCount: unreadCount ?? this.unreadCount,
      interlocutorId: interlocutorId ?? this.interlocutorId,
      interlocutorName: interlocutorName ?? this.interlocutorName,
      interlocutorAvatar: interlocutorAvatar ?? this.interlocutorAvatar,
      interlocutorAvatarMediaId:
          interlocutorAvatarMediaId ?? this.interlocutorAvatarMediaId,
      avatarMediaId: avatarMediaId ?? this.avatarMediaId,
      isMuted: isMuted ?? this.isMuted,
      isPinned: isPinned ?? this.isPinned,
      memberCount: memberCount ?? this.memberCount,
      onlineCount: onlineCount ?? this.onlineCount,
      memberRole: memberRole ?? this.memberRole,
      lastMessageType: lastMessageType ?? this.lastMessageType,
      blockedByMe: blockedByMe ?? this.blockedByMe,
      blockedMe: blockedMe ?? this.blockedMe,
    );
  }
}

@JsonSerializable()
class Message {
  final int id;
  @JsonKey(name: 'chat_id')
  final String chatId;
  @JsonKey(name: 'sender_id')
  final String? senderId; // null for SYSTEM messages (no author)
  final String content;
  final String type; // TEXT, IMAGE, VIDEO, AUDIO, FILE, VOICE, SYSTEM, SUDA_TRANSFER, DONATION
  final String status; // SENT, READ, SENDING, FAILED
  @JsonKey(name: 'created_at')
  final String createdAt;
  @JsonKey(name: 'client_side_id')
  final String? clientSideId;
  @JsonKey(name: 'edited_at')
  final String? editedAt;
  @JsonKey(name: 'attachment_media_id')
  final String? attachmentMediaId;
  @JsonKey(name: 'reply_to_message_id')
  final int? replyToMessageId;
  @JsonKey(name: 'forwarded_from_chat')
  final String? forwardedFromChat;
  @JsonKey(name: 'forwarded_from_msg')
  final int? forwardedFromMsg;
  @JsonKey(name: 'deleted_at')
  final String? deletedAt;
  // Channel posts only: number of comments (server-provided), null if absent.
  @JsonKey(name: 'comment_count')
  final int? commentCount;
  // Hydrated reply preview (§3.1) — populated by backend when isReply is true.
  @JsonKey(name: 'reply_preview')
  final MessagePreview? replyPreview;
  // Hydrated forward source (§3.1) — populated by backend when isForwarded is true.
  @JsonKey(name: 'forwarded_from')
  final MessagePreview? forwardedFrom;

  bool get isEdited => editedAt != null;
  bool get isDeleted => deletedAt != null;
  bool get isForwarded => forwardedFromChat != null || forwardedFrom != null;
  bool get isReply => replyToMessageId != null;

  Message({
    required this.id,
    required this.chatId,
    this.senderId,
    required this.content,
    required this.type,
    required this.status,
    required this.createdAt,
    this.clientSideId,
    this.editedAt,
    this.attachmentMediaId,
    this.replyToMessageId,
    this.forwardedFromChat,
    this.forwardedFromMsg,
    this.deletedAt,
    this.commentCount,
    this.replyPreview,
    this.forwardedFrom,
  });

  factory Message.fromJson(Map<String, dynamic> json) => _$MessageFromJson(json);

  Message copyWith({
    int? id,
    String? chatId,
    String? senderId,
    String? content,
    String? type,
    String? status,
    String? createdAt,
    String? clientSideId,
    String? editedAt,
    String? attachmentMediaId,
    int? replyToMessageId,
    String? forwardedFromChat,
    int? forwardedFromMsg,
    String? deletedAt,
    int? commentCount,
    MessagePreview? replyPreview,
    MessagePreview? forwardedFrom,
  }) {
    return Message(
      id: id ?? this.id,
      chatId: chatId ?? this.chatId,
      senderId: senderId ?? this.senderId,
      content: content ?? this.content,
      type: type ?? this.type,
      status: status ?? this.status,
      createdAt: createdAt ?? this.createdAt,
      clientSideId: clientSideId ?? this.clientSideId,
      editedAt: editedAt ?? this.editedAt,
      attachmentMediaId: attachmentMediaId ?? this.attachmentMediaId,
      replyToMessageId: replyToMessageId ?? this.replyToMessageId,
      forwardedFromChat: forwardedFromChat ?? this.forwardedFromChat,
      forwardedFromMsg: forwardedFromMsg ?? this.forwardedFromMsg,
      deletedAt: deletedAt ?? this.deletedAt,
      commentCount: commentCount ?? this.commentCount,
      replyPreview: replyPreview ?? this.replyPreview,
      forwardedFrom: forwardedFrom ?? this.forwardedFrom,
    );
  }
}

@JsonSerializable()
class SendMessageRequest {
  @JsonKey(name: 'chat_id')
  final String chatId;
  final String content;
  final String type;
  @JsonKey(name: 'client_side_id')
  final String clientSideId;
  @JsonKey(name: 'attachment_media_id')
  final String? attachmentMediaId;
  @JsonKey(name: 'reply_to_message_id')
  final int? replyToMessageId;

  SendMessageRequest({
    required this.chatId,
    required this.content,
    required this.type,
    required this.clientSideId,
    this.attachmentMediaId,
    this.replyToMessageId,
  });

  Map<String, dynamic> toJson() => _$SendMessageRequestToJson(this);
}

@JsonSerializable()
class UserStatus {
  @JsonKey(name: 'user_id')
  final String userId;
  @JsonKey(name: 'is_online')
  final bool isOnline;
  @JsonKey(name: 'last_seen_at')
  final String? lastSeenAt;

  const UserStatus({
    required this.userId,
    required this.isOnline,
    this.lastSeenAt,
  });

  factory UserStatus.fromJson(Map<String, dynamic> json) => _$UserStatusFromJson(json);
}

@JsonSerializable()
class MessageReader {
  @JsonKey(name: 'user_id')
  final String userId;
  final String username;
  @JsonKey(name: 'display_name')
  final String displayName;
  @JsonKey(name: 'read_at')
  final String readAt;

  const MessageReader({
    required this.userId,
    required this.username,
    required this.displayName,
    required this.readAt,
  });

  factory MessageReader.fromJson(Map<String, dynamic> json) => _$MessageReaderFromJson(json);
}

@JsonSerializable()
class CreateChatRequest {
  final String type; // DIRECT, GROUP, CHANNEL, SAVED
  @JsonKey(name: 'target_id')
  final String? targetId;
  final String? name;
  final String? description;
  @JsonKey(name: 'member_ids')
  final List<String>? memberIds;
  final String? visibility; // PUBLIC, PRIVATE
  final String? username; // @handle для каналов

  CreateChatRequest({
    required this.type,
    this.targetId,
    this.name,
    this.description,
    this.memberIds,
    this.visibility,
    this.username,
  });

  Map<String, dynamic> toJson() => _$CreateChatRequestToJson(this);
}

@JsonSerializable()
class ChatMember {
  @JsonKey(name: 'user_id')
  final String userId;
  final String username;
  @JsonKey(name: 'display_name')
  final String displayName;
  @JsonKey(name: 'avatar_media_id')
  final String? avatarMediaId;
  final String role; // OWNER, ADMIN, MEMBER
  @JsonKey(name: 'joined_at')
  final String? joinedAt;

  const ChatMember({
    required this.userId,
    required this.username,
    required this.displayName,
    this.avatarMediaId,
    required this.role,
    this.joinedAt,
  });

  factory ChatMember.fromJson(Map<String, dynamic> json) => _$ChatMemberFromJson(json);
}

@JsonSerializable()
class ChatInfo {
  final String id;
  final String type;
  final String? name;
  final String? description;
  @JsonKey(name: 'avatar_media_id')
  final String? avatarMediaId;
  @JsonKey(defaultValue: [])
  final List<ChatMember> members;
  @JsonKey(name: 'member_count', defaultValue: 0)
  final int memberCount;
  // Channel-specific fields
  final String? username;
  @JsonKey(name: 'is_verified', defaultValue: false)
  final bool isVerified;
  // Block status for DIRECT chats (§4.1) — always false for non-DIRECT.
  @JsonKey(name: 'blocked_by_me', defaultValue: false)
  final bool blockedByMe;
  @JsonKey(name: 'blocked_me', defaultValue: false)
  final bool blockedMe;

  const ChatInfo({
    required this.id,
    required this.type,
    this.name,
    this.description,
    this.avatarMediaId,
    this.members = const [],
    this.memberCount = 0,
    this.username,
    this.isVerified = false,
    this.blockedByMe = false,
    this.blockedMe = false,
  });

  factory ChatInfo.fromJson(Map<String, dynamic> json) => _$ChatInfoFromJson(json);
}

/// One shared-media entry (image/video/document/audio) from
/// GET /chats/{id}/media. Parsed leniently — the element may be a bare media id
/// string or an object with a media id plus optional metadata.
class MediaItem {
  final String mediaId;
  final int? messageId;
  final String? createdAt;
  final String? name;

  const MediaItem({
    required this.mediaId,
    this.messageId,
    this.createdAt,
    this.name,
  });

  static MediaItem? tryParse(dynamic raw) {
    if (raw is String) {
      return raw.isEmpty ? null : MediaItem(mediaId: raw);
    }
    if (raw is Map) {
      final m = Map<String, dynamic>.from(raw);
      final id = (m['media_id'] ?? m['attachment_media_id'] ?? m['id'])?.toString();
      if (id == null || id.isEmpty) return null;
      return MediaItem(
        mediaId: id,
        messageId: m['message_id'] is int
            ? m['message_id'] as int
            : int.tryParse(m['message_id']?.toString() ?? ''),
        createdAt: m['created_at'] as String?,
        name: (m['original_name'] ?? m['name'] ?? m['file_name'] ?? m['content'])
            as String?,
      );
    }
    return null;
  }
}

/// Shared media of a chat grouped by kind (GET /chats/{id}/media, §5.3).
class ChatMedia {
  final List<MediaItem> images;
  final List<MediaItem> videos;
  final List<MediaItem> documents;
  final List<MediaItem> audio;

  const ChatMedia({
    this.images = const [],
    this.videos = const [],
    this.documents = const [],
    this.audio = const [],
  });

  bool get isEmpty =>
      images.isEmpty && videos.isEmpty && documents.isEmpty && audio.isEmpty;

  static List<MediaItem> _list(dynamic raw) {
    if (raw is! List) return const [];
    return raw
        .map(MediaItem.tryParse)
        .whereType<MediaItem>()
        .toList(growable: false);
  }

  factory ChatMedia.fromJson(Map<String, dynamic> json) => ChatMedia(
        images: _list(json['images']),
        videos: _list(json['videos']),
        documents: _list(json['documents'] ?? json['files']),
        audio: _list(json['audio']),
      );
}
