// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'user_preferences.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

UserPreferences _$UserPreferencesFromJson(Map<String, dynamic> json) =>
    UserPreferences(
      language: json['language'] as String? ?? 'en',
      notificationsEnabled: json['notifications_enabled'] as bool? ?? true,
      showLastSeen: json['show_last_seen'] as String? ?? 'EVERYONE',
      showOnline: json['show_online'] as bool? ?? true,
      readReceiptsEnabled: json['read_receipts_enabled'] as bool? ?? true,
      autoDownloadMedia: json['auto_download_media'] as String? ?? 'WIFI_ONLY',
      theme: json['theme'] as String? ?? 'AUTO',
    );

Map<String, dynamic> _$UserPreferencesToJson(UserPreferences instance) =>
    <String, dynamic>{
      'language': instance.language,
      'notifications_enabled': instance.notificationsEnabled,
      'show_last_seen': instance.showLastSeen,
      'show_online': instance.showOnline,
      'read_receipts_enabled': instance.readReceiptsEnabled,
      'auto_download_media': instance.autoDownloadMedia,
      'theme': instance.theme,
    };
