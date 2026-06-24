/// Response model for GET /api/v1/messenger/channels/{id}/view
class ChannelViewModel {
  final String id;
  final String name;
  final String? description;
  final String? username;
  final String? avatarMediaId;
  final String visibility; // "PUBLIC" | "PRIVATE"
  final int subscriberCount;
  final bool commentsEnabled;
  final bool isMember;
  final String myRole; // "OWNER"|"ADMIN"|"SUBSCRIBER"|""
  final String joinStatus; // "none"|"pending_request"|"pending_invite"|"member"
  final bool canRead;
  final String createdAt;

  const ChannelViewModel({
    required this.id,
    required this.name,
    this.description,
    this.username,
    this.avatarMediaId,
    required this.visibility,
    required this.subscriberCount,
    required this.commentsEnabled,
    required this.isMember,
    required this.myRole,
    required this.joinStatus,
    required this.canRead,
    required this.createdAt,
  });

  bool get isOwner => myRole == 'OWNER';
  bool get isAdmin => myRole == 'ADMIN' || myRole == 'OWNER';
  bool get isPublic => visibility == 'PUBLIC';

  factory ChannelViewModel.fromJson(Map<String, dynamic> json) => ChannelViewModel(
        id: json['id'] as String,
        name: json['name'] as String? ?? '',
        description: json['description'] as String?,
        username: json['username'] as String?,
        avatarMediaId: json['avatar_media_id'] as String?,
        visibility: json['visibility'] as String? ?? 'PUBLIC',
        subscriberCount: json['subscriber_count'] as int? ?? 0,
        commentsEnabled: json['comments_enabled'] as bool? ?? false,
        isMember: json['is_member'] as bool? ?? false,
        myRole: json['my_role'] as String? ?? '',
        joinStatus: json['join_status'] as String? ?? 'none',
        canRead: json['can_read'] as bool? ?? false,
        createdAt: json['created_at'] as String? ?? '',
      );
}

/// One item from GET /api/v1/messenger/channels/invites
class ChannelInviteItem {
  final String channelId;
  final String name;
  final String? username;
  final String? avatarMediaId;
  final String? description;
  final int subscriberCount;
  final String createdAt;

  const ChannelInviteItem({
    required this.channelId,
    required this.name,
    this.username,
    this.avatarMediaId,
    this.description,
    required this.subscriberCount,
    required this.createdAt,
  });

  factory ChannelInviteItem.fromJson(Map<String, dynamic> json) => ChannelInviteItem(
        channelId: json['channel_id'] as String,
        name: json['name'] as String? ?? '',
        username: json['username'] as String?,
        avatarMediaId: json['avatar_media_id'] as String?,
        description: json['description'] as String?,
        subscriberCount: json['subscriber_count'] as int? ?? 0,
        createdAt: json['created_at'] as String? ?? '',
      );
}
