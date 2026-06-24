import 'package:flutter/material.dart';

import '../../domain/models/app_failure.dart';
import '../../domain/models/error_dictionary.dart';

class AppFeedback {
  AppFeedback._();

  static final GlobalKey<ScaffoldMessengerState> messengerKey =
      GlobalKey<ScaffoldMessengerState>();
  static String? _lastMessage;
  static DateTime? _lastShownAt;

  static void showError(
    String? message, {
    String fallback = 'Something went wrong',
  }) {
    _show(
      message: safeErrorMessage(message, fallback: fallback),
      backgroundColor: const Color(0xFFB3261E),
      icon: Icons.error_outline,
    );
  }

  static void showSuccess(String message) {
    _show(
      message: _sanitizeMessage(message, fallback: 'Done'),
      backgroundColor: const Color(0xFF1E7F4D),
      icon: Icons.check_circle_outline,
    );
  }

  static String failureMessage(
    AppFailure? failure, {
    String fallback = 'Something went wrong',
  }) {
    if (failure == null) return fallback;

    final code = failure.code?.toUpperCase();
    if (code == 'NETWORK_ERROR') return 'No internet connection';

    // Try ErrorDictionary first for precise user-facing message
    final statusCode = _parseStatusCode(failure.code);
    final dictMessage = ErrorDictionary.humanMessage(
      statusCode: statusCode,
      errorKey: failure.code,
    );
    if (dictMessage != 'Something went wrong. Please try again.') {
      return dictMessage;
    }

    return safeErrorMessage(failure.message, fallback: fallback);
  }

  static int? _parseStatusCode(String? code) {
    if (code == null) return null;
    return int.tryParse(code);
  }

  static String safeErrorMessage(
    String? message, {
    String fallback = 'Something went wrong',
  }) {
    final cleaned = _sanitizeMessage(message, fallback: fallback);
    final lower = cleaned.toLowerCase();

    final noisy = [
      'this exception was thrown',
      'dioexception',
      'stack trace',
      'socketexception',
      'xmlhttprequest',
      'html',
    ];
    for (final pattern in noisy) {
      if (lower.contains(pattern)) {
        return fallback;
      }
    }

    if (cleaned.length > 90) {
      return fallback;
    }
    return cleaned;
  }

  static String _sanitizeMessage(String? message, {required String fallback}) {
    if (message == null) return fallback;

    final compact = message
        .replaceAll(RegExp(r'\s+'), ' ')
        .replaceAll('\n', ' ')
        .trim();

    if (compact.isEmpty) {
      return fallback;
    }
    return compact;
  }

  static void _show({
    required String message,
    required Color backgroundColor,
    required IconData icon,
  }) {
    final now = DateTime.now();
    if (_lastMessage == message &&
        _lastShownAt != null &&
        now.difference(_lastShownAt!).inMilliseconds < 1200) {
      return;
    }
    _lastMessage = message;
    _lastShownAt = now;

    final messenger = messengerKey.currentState;
    if (messenger == null) {
      return;
    }

    messenger
      ..hideCurrentSnackBar()
      ..showSnackBar(
        SnackBar(
          behavior: SnackBarBehavior.floating,
          margin: const EdgeInsets.fromLTRB(16, 0, 16, 16),
          backgroundColor: backgroundColor,
          duration: const Duration(seconds: 3),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(12),
          ),
          content: Row(
            children: [
              Icon(icon, color: Colors.white),
              const SizedBox(width: 10),
              Expanded(
                child: Text(
                  message,
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                  style: const TextStyle(
                    color: Colors.white,
                    fontWeight: FontWeight.w500,
                  ),
                ),
              ),
            ],
          ),
        ),
      );
  }
}
