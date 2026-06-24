/// Response model for GET /api/v1/messenger/channels/{id}/settings (OWNER/ADMIN).
class ChannelSettings {
  final String id;
  final String name;
  final String? description;
  final String? username;
  final String? avatarMediaId;
  final String visibility; // "PUBLIC" | "PRIVATE"
  final bool commentsEnabled;

  const ChannelSettings({
    required this.id,
    required this.name,
    this.description,
    this.username,
    this.avatarMediaId,
    required this.visibility,
    required this.commentsEnabled,
  });

  bool get isPublic => visibility == 'PUBLIC';

  factory ChannelSettings.fromJson(Map<String, dynamic> json) => ChannelSettings(
        id: json['id'] as String? ?? '',
        name: json['name'] as String? ?? '',
        description: json['description'] as String?,
        username: json['username'] as String?,
        avatarMediaId: json['avatar_media_id'] as String?,
        visibility: json['visibility'] as String? ?? 'PUBLIC',
        commentsEnabled: json['comments_enabled'] as bool? ?? false,
      );
}
