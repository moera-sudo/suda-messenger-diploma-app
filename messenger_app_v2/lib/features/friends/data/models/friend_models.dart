import 'package:json_annotation/json_annotation.dart';

part 'friend_models.g.dart';

@JsonSerializable()
class FriendModel {
  @JsonKey(name: 'user_id')
  final String userId;
  final String username;
  @JsonKey(name: 'display_name')
  final String displayName;
  @JsonKey(name: 'avatar_media_id')
  final String? avatarMediaId;
  @JsonKey(name: 'is_online', defaultValue: false)
  final bool isOnline;
  @JsonKey(name: 'became_at', defaultValue: '')
  final String becameAt;

  const FriendModel({
    required this.userId,
    required this.username,
    required this.displayName,
    this.avatarMediaId,
    this.isOnline = false,
    this.becameAt = '',
  });

  factory FriendModel.fromJson(Map<String, dynamic> json) =>
      _$FriendModelFromJson(json);
}

@JsonSerializable()
class FriendRequestModel {
  @JsonKey(name: 'request_id')
  final String requestId;
  final String direction; // 'incoming' | 'outgoing'
  final String status;
  @JsonKey(name: 'other_user_id')
  final String otherUserId;
  @JsonKey(name: 'other_username')
  final String otherUsername;
  @JsonKey(name: 'other_display_name')
  final String otherDisplayName;
  @JsonKey(name: 'other_avatar_media_id')
  final String? otherAvatarMediaId;
  @JsonKey(name: 'created_at', defaultValue: '')
  final String createdAt;

  const FriendRequestModel({
    required this.requestId,
    required this.direction,
    required this.status,
    required this.otherUserId,
    required this.otherUsername,
    required this.otherDisplayName,
    this.otherAvatarMediaId,
    this.createdAt = '',
  });

  factory FriendRequestModel.fromJson(Map<String, dynamic> json) =>
      _$FriendRequestModelFromJson(json);
}

@JsonSerializable()
class FriendStatusModel {
  @JsonKey(name: 'user_id')
  final String userId;
  /// NONE | PENDING_SENT | PENDING_RECEIVED | FRIENDS | BLOCKED
  final String relation;
  @JsonKey(name: 'request_id')
  final String? requestId;

  const FriendStatusModel({
    required this.userId,
    required this.relation,
    this.requestId,
  });

  factory FriendStatusModel.fromJson(Map<String, dynamic> json) =>
      _$FriendStatusModelFromJson(json);
}
