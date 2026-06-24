import 'package:json_annotation/json_annotation.dart';

part 'user_models.g.dart';

@JsonSerializable()
class UserProfileModel {
  final String id;
  final String username;
  @JsonKey(name: 'display_name')
  final String displayName;
  @JsonKey(name: 'first_name')
  final String? firstName;
  @JsonKey(name: 'last_name')
  final String? lastName;
  final String? bio;
  @JsonKey(name: 'avatar_media_id')
  final String? avatarMediaId;
  @JsonKey(name: 'wallet_address')
  final String? walletAddress;
  @JsonKey(name: 'is_verified', defaultValue: false)
  final bool isVerified;
  // Populated from the status endpoint, not the profile endpoint
  @JsonKey(name: 'is_online', defaultValue: false)
  final bool isOnline;
  @JsonKey(name: 'last_seen_at')
  final String? lastSeenAt;

  const UserProfileModel({
    required this.id,
    required this.username,
    required this.displayName,
    this.firstName,
    this.lastName,
    this.bio,
    this.avatarMediaId,
    this.walletAddress,
    this.isVerified = false,
    this.isOnline = false,
    this.lastSeenAt,
  });

  factory UserProfileModel.fromJson(Map<String, dynamic> json) =>
      _$UserProfileModelFromJson(json);

  Map<String, dynamic> toJson() => _$UserProfileModelToJson(this);

  UserProfileModel copyWith({
    String? id,
    String? username,
    String? displayName,
    String? firstName,
    String? lastName,
    String? bio,
    String? avatarMediaId,
    String? walletAddress,
    bool? isVerified,
    bool? isOnline,
    String? lastSeenAt,
  }) {
    return UserProfileModel(
      id: id ?? this.id,
      username: username ?? this.username,
      displayName: displayName ?? this.displayName,
      firstName: firstName ?? this.firstName,
      lastName: lastName ?? this.lastName,
      bio: bio ?? this.bio,
      avatarMediaId: avatarMediaId ?? this.avatarMediaId,
      walletAddress: walletAddress ?? this.walletAddress,
      isVerified: isVerified ?? this.isVerified,
      isOnline: isOnline ?? this.isOnline,
      lastSeenAt: lastSeenAt ?? this.lastSeenAt,
    );
  }
}

@JsonSerializable()
class ContactModel {
  @JsonKey(name: 'contact_id')
  final String contactId;
  @JsonKey(name: 'custom_name')
  final String? customName;
  // Populated by a secondary profile fetch — not from JSON
  @JsonKey(includeFromJson: false, includeToJson: false)
  final UserProfileModel? profile;

  const ContactModel({
    required this.contactId,
    this.customName,
    this.profile,
  });

  factory ContactModel.fromJson(Map<String, dynamic> json) =>
      _$ContactModelFromJson(json);

  Map<String, dynamic> toJson() => _$ContactModelToJson(this);

  ContactModel withProfile(UserProfileModel profile) => ContactModel(
        contactId: contactId,
        customName: customName,
        profile: profile,
      );

  String get displayLabel => customName?.isNotEmpty == true
      ? customName!
      : profile?.displayName ?? contactId;
}

/// Response model for GET /api/v1/messenger/users/blocked (array of objects).
class BlockedUserModel {
  final String userId;
  final String username;
  final String displayName;
  final String? avatarMediaId;

  const BlockedUserModel({
    required this.userId,
    required this.username,
    required this.displayName,
    this.avatarMediaId,
  });

  factory BlockedUserModel.fromJson(Map<String, dynamic> json) => BlockedUserModel(
        userId: json['user_id'] as String,
        username: json['username'] as String? ?? '',
        displayName: json['display_name'] as String? ?? '',
        avatarMediaId: json['avatar_media_id'] as String?,
      );
}
