import 'dart:async';

import 'package:injectable/injectable.dart';

/// App-wide bus for cross-cutting auth events.
///
/// The Dio [AuthInterceptor] lives below the BLoC layer and cannot navigate or
/// touch [AuthBloc] directly. When it detects an unrecoverable deauthorization
/// (refresh token is invalid), it broadcasts here. [AuthBloc] listens and flips
/// the app into the unauthenticated state, which the global listener in
/// `app.dart` turns into a redirect to the Welcome screen.
@lazySingleton
class AuthEventBus {
  final _forcedLogoutController = StreamController<void>.broadcast();

  /// Emits whenever the session is forcibly terminated (invalid refresh token,
  /// current session revoked, etc.).
  Stream<void> get onForcedLogout => _forcedLogoutController.stream;

  /// Signal an unrecoverable deauthorization. Safe to call multiple times.
  void emitForcedLogout() {
    if (!_forcedLogoutController.isClosed) {
      _forcedLogoutController.add(null);
    }
  }

  @disposeMethod
  void dispose() {
    _forcedLogoutController.close();
  }
}
