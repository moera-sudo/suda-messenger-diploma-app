import 'package:json_annotation/json_annotation.dart';

part 'user_preferences.g.dart';

/// User-level preferences. Matches backend contract:
/// GET/PUT /api/v1/messenger/user/preferences (Frontend.md §5.10).
@JsonSerializable()
class UserPreferences {
  @JsonKey(defaultValue: 'en')
  final String language;

  @JsonKey(name: 'notifications_enabled', defaultValue: true)
  final bool notificationsEnabled;

  /// EVERYONE / CONTACTS / NOBODY
  @JsonKey(name: 'show_last_seen', defaultValue: 'EVERYONE')
  final String showLastSeen;

  @JsonKey(name: 'show_online', defaultValue: true)
  final bool showOnline;

  @JsonKey(name: 'read_receipts_enabled', defaultValue: true)
  final bool readReceiptsEnabled;

  /// ALWAYS / WIFI_ONLY / NEVER
  @JsonKey(name: 'auto_download_media', defaultValue: 'WIFI_ONLY')
  final String autoDownloadMedia;

  /// LIGHT / DARK / AUTO — kept as-is, color theme is managed locally by SettingsCubit.
  @JsonKey(defaultValue: 'AUTO')
  final String theme;

  const UserPreferences({
    this.language = 'en',
    this.notificationsEnabled = true,
    this.showLastSeen = 'EVERYONE',
    this.showOnline = true,
    this.readReceiptsEnabled = true,
    this.autoDownloadMedia = 'WIFI_ONLY',
    this.theme = 'AUTO',
  });

  UserPreferences copyWith({
    String? language,
    bool? notificationsEnabled,
    String? showLastSeen,
    bool? showOnline,
    bool? readReceiptsEnabled,
    String? autoDownloadMedia,
    String? theme,
  }) {
    return UserPreferences(
      language: language ?? this.language,
      notificationsEnabled: notificationsEnabled ?? this.notificationsEnabled,
      showLastSeen: showLastSeen ?? this.showLastSeen,
      showOnline: showOnline ?? this.showOnline,
      readReceiptsEnabled: readReceiptsEnabled ?? this.readReceiptsEnabled,
      autoDownloadMedia: autoDownloadMedia ?? this.autoDownloadMedia,
      theme: theme ?? this.theme,
    );
  }

  factory UserPreferences.fromJson(Map<String, dynamic> json) =>
      _$UserPreferencesFromJson(json);

  Map<String, dynamic> toJson() => _$UserPreferencesToJson(this);
}
