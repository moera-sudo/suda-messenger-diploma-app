import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:injectable/injectable.dart';
import '../logger/app_logger.dart';

@lazySingleton
class SecureStorageClient {
  final FlutterSecureStorage _storage;
  final AppLogger _logger;

  SecureStorageClient(this._storage, this._logger);

  static const _accessTokenKey = 'access_token';
  static const _refreshTokenKey = 'refresh_token';

  Future<String?> getAccessToken() async {
    return await _storage.read(key: _accessTokenKey);
  }

  Future<String?> getRefreshToken() async {
    return await _storage.read(key: _refreshTokenKey);
  }

  Future<void> saveTokens({
    required String access,
    required String refresh,
  }) async {
    _logger.debug("Saving tokens to secure storage");
    await _storage.write(key: _accessTokenKey, value: access);
    await _storage.write(key: _refreshTokenKey, value: refresh);
    _logger.debug("Token successfully saved");
  }

  Future<void> clearTokens() async {
    await _storage.delete(key: _accessTokenKey);
    await _storage.delete(key: _refreshTokenKey);
  }

  Future<void> clearStorageForDebug() async {
    _logger.warning("Clearing all data from Secure Storage");
    await _storage.deleteAll();
  }
}
