class SessionModel {
  final String id;
  final String userAgent;
  final String clientIp;
  final String createdAt;
  final String? lastUsedAt;
  final String deviceName;

  /// Whether this is the session the request was made from. Backend must
  /// populate `is_current` (see PreDefence-Fix §5 / plan C3). Defaults to false
  /// until the backend adds the flag — no client-side guessing.
  final bool isCurrent;

  const SessionModel({
    required this.id,
    required this.userAgent,
    required this.clientIp,
    required this.createdAt,
    this.lastUsedAt,
    this.deviceName = '',
    this.isCurrent = false,
  });

  factory SessionModel.fromJson(Map<String, dynamic> json) {
    return SessionModel(
      id: json['id'] as String? ?? '',
      userAgent: json['user_agent'] as String? ?? '',
      clientIp: json['client_ip'] as String? ?? '',
      createdAt: json['created_at'] as String? ?? '',
      lastUsedAt: json['last_used_at'] as String?,
      deviceName: json['device_name'] as String? ?? '',
      isCurrent: json['is_current'] as bool? ?? false,
    );
  }
}
