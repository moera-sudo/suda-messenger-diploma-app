import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';

import '../../../app/DI/get_it.dart';
import '../api/api_client.dart';
import '../logger/app_logger.dart';

/// Top-level FCM background handler — runs in a separate isolate,
/// must be a top-level function annotated with @pragma('vm:entry-point').
@pragma('vm:entry-point')
Future<void> firebaseBackgroundMessageHandler(RemoteMessage message) async {
  // Background isolate: no DI, no UI. Firebase is already initialized by the plugin.
}

class NotificationService {
  static const _channelId = 'messenger_messages';
  static const _channelName = 'Messages';

  static final _localNotifications = FlutterLocalNotificationsPlugin();
  static const _androidChannel = AndroidNotificationChannel(
    _channelId,
    _channelName,
    importance: Importance.high,
    enableVibration: true,
  );

  /// Chat ID stored when the user taps a notification (background or terminated).
  static String? _pendingChatId;

  /// Returns and clears the pending chat ID after a notification tap.
  static String? consumePendingChatId() {
    final id = _pendingChatId;
    _pendingChatId = null;
    return id;
  }

  /// Phase 1 — call from bootstrap BEFORE authentication.
  /// Creates the Android notification channel, requests permission, and wires
  /// foreground / background-tap handlers.
  static Future<void> setup() async {
    final logger = sl<AppLogger>();

    await _setupLocalNotifications(logger);
    await _requestPermission(logger);
    _handleForegroundMessages(logger);
    _handleBackgroundTap(logger);
    await _handleInitialMessage(logger);

    logger.info('[NotificationService] setup complete');
  }

  /// Phase 2 — call after the user is authenticated.
  /// Registers the FCM token with our backend and listens for token rotation.
  static Future<void> registerToken() async {
    final logger = sl<AppLogger>();
    final api = sl<ApiClient>();

    try {
      final token = await FirebaseMessaging.instance.getToken();
      if (token != null) {
        await api.post('/user/device', data: {'token': token, 'platform': 'android'});
        logger.info('[NotificationService] FCM token registered');
      }
    } catch (e) {
      logger.error('[NotificationService] Token registration failed: $e');
    }

    FirebaseMessaging.instance.onTokenRefresh.listen((newToken) async {
      try {
        await api.post('/user/device', data: {'token': newToken, 'platform': 'android'});
        logger.info('[NotificationService] FCM token refreshed');
      } catch (e) {
        logger.error('[NotificationService] Token refresh failed: $e');
      }
    });
  }

  // ── Private helpers ──────────────────────────────────────────

  static Future<void> _setupLocalNotifications(AppLogger logger) async {
    const initSettings = InitializationSettings(
      android: AndroidInitializationSettings('@mipmap/ic_launcher'),
    );
    await _localNotifications.initialize(
      settings: initSettings,
      onDidReceiveNotificationResponse: (details) {
        if (details.payload != null) {
          _pendingChatId = details.payload;
          logger.debug('[NotificationService] Local notification tapped → chatId=${details.payload}');
        }
      },
    );
    await _localNotifications
        .resolvePlatformSpecificImplementation<AndroidFlutterLocalNotificationsPlugin>()
        ?.createNotificationChannel(_androidChannel);
    logger.debug('[NotificationService] Android notification channel created');
  }

  static Future<void> _requestPermission(AppLogger logger) async {
    final settings = await FirebaseMessaging.instance.requestPermission(
      alert: true,
      badge: true,
      sound: true,
    );
    logger.info('[NotificationService] Permission status: ${settings.authorizationStatus}');
  }

  /// Shows a local notification while the app is in the foreground.
  static void _handleForegroundMessages(AppLogger logger) {
    FirebaseMessaging.onMessage.listen((message) {
      final notification = message.notification;
      if (notification == null) return;

      logger.debug('[NotificationService] Foreground message: ${notification.title}');

      _localNotifications.show(
        id: message.hashCode,
        title: notification.title,
        body: notification.body,
        notificationDetails: const NotificationDetails(
          android: AndroidNotificationDetails(
            _channelId,
            _channelName,
            importance: Importance.high,
            priority: Priority.high,
          ),
        ),
        payload: message.data['chat_id'],
      );
    });
  }

  /// Handles notification tap when the app is in the background (not terminated).
  static void _handleBackgroundTap(AppLogger logger) {
    FirebaseMessaging.onMessageOpenedApp.listen((message) {
      final chatId = message.data['chat_id'] as String?;
      if (chatId != null) {
        _pendingChatId = chatId;
        logger.info('[NotificationService] Background tap → chatId=$chatId');
      }
    });
  }

  /// Handles notification tap when the app was terminated (cold start).
  static Future<void> _handleInitialMessage(AppLogger logger) async {
    final message = await FirebaseMessaging.instance.getInitialMessage();
    if (message != null) {
      final chatId = message.data['chat_id'] as String?;
      if (chatId != null) {
        _pendingChatId = chatId;
        logger.info('[NotificationService] Cold-start tap → chatId=$chatId');
      }
    }
  }
}
