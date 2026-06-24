// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'friend_models.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

FriendModel _$FriendModelFromJson(Map<String, dynamic> json) => FriendModel(
  userId: json['user_id'] as String,
  username: json['username'] as String,
  displayName: json['display_name'] as String,
  avatarMediaId: json['avatar_media_id'] as String?,
  isOnline: json['is_online'] as bool? ?? false,
  becameAt: json['became_at'] as String? ?? '',
);

Map<String, dynamic> _$FriendModelToJson(FriendModel instance) =>
    <String, dynamic>{
      'user_id': instance.userId,
      'username': instance.username,
      'display_name': instance.displayName,
      'avatar_media_id': instance.avatarMediaId,
      'is_online': instance.isOnline,
      'became_at': instance.becameAt,
    };

FriendRequestModel _$FriendRequestModelFromJson(Map<String, dynamic> json) =>
    FriendRequestModel(
      requestId: json['request_id'] as String,
      direction: json['direction'] as String,
      status: json['status'] as String,
      otherUserId: json['other_user_id'] as String,
      otherUsername: json['other_username'] as String,
      otherDisplayName: json['other_display_name'] as String,
      otherAvatarMediaId: json['other_avatar_media_id'] as String?,
      createdAt: json['created_at'] as String? ?? '',
    );

Map<String, dynamic> _$FriendRequestModelToJson(FriendRequestModel instance) =>
    <String, dynamic>{
      'request_id': instance.requestId,
      'direction': instance.direction,
      'status': instance.status,
      'other_user_id': instance.otherUserId,
      'other_username': instance.otherUsername,
      'other_display_name': instance.otherDisplayName,
      'other_avatar_media_id': instance.otherAvatarMediaId,
      'created_at': instance.createdAt,
    };

FriendStatusModel _$FriendStatusModelFromJson(Map<String, dynamic> json) =>
    FriendStatusModel(
      userId: json['user_id'] as String,
      relation: json['relation'] as String,
      requestId: json['request_id'] as String?,
    );

Map<String, dynamic> _$FriendStatusModelToJson(FriendStatusModel instance) =>
    <String, dynamic>{
      'user_id': instance.userId,
      'relation': instance.relation,
      'request_id': instance.requestId,
    };
