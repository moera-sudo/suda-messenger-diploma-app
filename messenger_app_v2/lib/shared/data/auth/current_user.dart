import 'package:injectable/injectable.dart';
import 'package:jwt_decoder/jwt_decoder.dart';

import '../storage/secure_storage_client.dart';

/// Resolves the current user's identity from the stored access token.
/// Centralizes the JWT decoding that was previously duplicated across BLoCs.
@lazySingleton
class CurrentUser {
  final SecureStorageClient _storage;

  CurrentUser(this._storage);

  /// Returns the current user's ID (`user_id` or `sub` claim), or null if
  /// there is no valid token.
  Future<String?> id() async {
    final token = await _storage.getAccessToken();
    if (token == null) return null;
    try {
      final decoded = JwtDecoder.decode(token);
      return decoded['user_id']?.toString() ?? decoded['sub']?.toString();
    } catch (_) {
      return null;
    }
  }
}
