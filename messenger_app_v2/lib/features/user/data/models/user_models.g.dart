// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'user_models.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

UserProfileModel _$UserProfileModelFromJson(Map<String, dynamic> json) =>
    UserProfileModel(
      id: json['id'] as String,
      username: json['username'] as String,
      displayName: json['display_name'] as String,
      firstName: json['first_name'] as String?,
      lastName: json['last_name'] as String?,
      bio: json['bio'] as String?,
      avatarMediaId: json['avatar_media_id'] as String?,
      walletAddress: json['wallet_address'] as String?,
      isVerified: json['is_verified'] as bool? ?? false,
      isOnline: json['is_online'] as bool? ?? false,
      lastSeenAt: json['last_seen_at'] as String?,
    );

Map<String, dynamic> _$UserProfileModelToJson(UserProfileModel instance) =>
    <String, dynamic>{
      'id': instance.id,
      'username': instance.username,
      'display_name': instance.displayName,
      'first_name': instance.firstName,
      'last_name': instance.lastName,
      'bio': instance.bio,
      'avatar_media_id': instance.avatarMediaId,
      'wallet_address': instance.walletAddress,
      'is_verified': instance.isVerified,
      'is_online': instance.isOnline,
      'last_seen_at': instance.lastSeenAt,
    };

ContactModel _$ContactModelFromJson(Map<String, dynamic> json) => ContactModel(
  contactId: json['contact_id'] as String,
  customName: json['custom_name'] as String?,
);

Map<String, dynamic> _$ContactModelToJson(ContactModel instance) =>
    <String, dynamic>{
      'contact_id': instance.contactId,
      'custom_name': instance.customName,
    };
