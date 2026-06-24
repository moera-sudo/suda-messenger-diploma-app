import 'package:dio/dio.dart';
import 'package:injectable/injectable.dart';
import 'package:messenger_app_v2/app/config/app_config.dart';
import '../../auth/auth_event_bus.dart';
import '../../logger/app_logger.dart';
import '../../storage/secure_storage_client.dart';
import '../../../domain/models/auth_token_pair.dart';

@injectable
class AuthInterceptor extends Interceptor {
  final SecureStorageClient _storage;
  final AppLogger _logger;
  final AuthEventBus _authEvents;
  final Dio _refreshDio;

  AuthInterceptor(
    this._storage,
    this._logger,
    this._authEvents,
    AppConfig config,
  ) : _refreshDio = Dio(BaseOptions(baseUrl: config.baseUrl));

  /// Single-flight guard: while a refresh is in progress, all other 401s await
  /// the same future instead of firing their own refresh. This prevents the
  /// refresh-token rotation race (parallel media-upload + screen requests each
  /// refreshing, the 2nd reusing an already-rotated refresh token → 401 →
  /// clearTokens → cascade logout).
  Future<AuthTokenPair?>? _ongoingRefresh;

  @override
  Future<void> onRequest(
      RequestOptions options, RequestInterceptorHandler handler) async {
    // Auth endpoints (login/register/refresh) don't need a bearer token.
    if (options.path.contains('/auth/')) {
      return handler.next(options);
    }

    final token = await _storage.getAccessToken();
    if (token != null) {
      options.headers['Authorization'] = 'Bearer $token';
    }

    super.onRequest(options, handler);
  }

  @override
  Future<void> onError(DioException err, ErrorInterceptorHandler handler) async {
    if (err.response?.statusCode != 401) {
      return super.onError(err, handler);
    }

    final path = err.requestOptions.path;

    // A 401 on the refresh/login path itself means the session is dead.
    if (path.contains('/auth/refresh') || path.contains('/auth/login')) {
      return handler.next(err);
    }

    _logger.warning('⚠️ 401 detected on $path. Trying to refresh session...');

    try {
      // Share a single refresh across all concurrent 401s.
      final refreshFuture = _ongoingRefresh ??= _refreshToken();
      final AuthTokenPair? newTokens;
      try {
        newTokens = await refreshFuture;
      } finally {
        _ongoingRefresh = null;
      }

      if (newTokens == null) {
        // No refresh token stored — unrecoverable deauthorization.
        _logger.error('No refresh token available — forcing logout.');
        await _forceLogout();
        return handler.next(err);
      }

      // Replay the failed request with the fresh access token.
      final opts = err.requestOptions;
      opts.headers['Authorization'] = 'Bearer ${newTokens.accessToken}';
      try {
        final cloneReq = await _refreshDio.request(
          opts.path,
          options: Options(
            method: opts.method,
            headers: opts.headers,
            contentType: opts.contentType,
          ),
          data: opts.data,
          queryParameters: opts.queryParameters,
        );
        return handler.resolve(cloneReq);
      } on DioException catch (retryErr) {
        // Even with a fresh token the server rejects us — session is dead.
        if (retryErr.response?.statusCode == 401) {
          _logger.error(
            'Retry after token refresh still returned 401 — forcing logout.',
            retryErr,
          );
          await _forceLogout();
        } else {
          _logger.warning(
            'Request retry after refresh failed (${retryErr.response?.statusCode}) — propagating.',
          );
        }
        return handler.next(err);
      }
    } on DioException catch (e) {
      // Reaches here only if _refreshToken() itself threw — i.e. /auth/refresh rejected.
      if (e.response?.statusCode == 401 || e.response?.statusCode == 403) {
        _logger.error('Refresh rejected by server — forcing logout.', e);
        await _forceLogout();
      } else {
        _logger.warning(
            'Refresh failed (network/transient) — keeping tokens: $e');
      }
      return handler.next(err);
    } catch (e) {
      // Unexpected error — keep tokens, surface the original failure.
      _logger.error('Unexpected refresh error — keeping tokens.', e);
      return handler.next(err);
    }
  }

  /// Clears tokens and broadcasts a forced-logout so the UI redirects to login.
  Future<void> _forceLogout() async {
    await _storage.clearTokens();
    _authEvents.emitForcedLogout();
  }

  Future<AuthTokenPair?> _refreshToken() async {
    final refreshToken = await _storage.getRefreshToken();
    if (refreshToken == null || refreshToken.isEmpty) return null;

    // Use the clean Dio (no interceptors) to hit the refresh endpoint.
    final response = await _refreshDio.post(
      '/api/v1/messenger/auth/refresh',
      data: {'refresh_token': refreshToken},
    );

    if (response.statusCode == 200) {
      final tokens = AuthTokenPair.fromJson(response.data);
      await _storage.saveTokens(
        access: tokens.accessToken,
        refresh: tokens.refreshToken,
      );
      _logger.info('🔑 Session refreshed successfully.');
      return tokens;
    }
    return null;
  }
}
