import 'package:dartz/dartz.dart';
import 'package:injectable/injectable.dart';

import '../../../../shared/data/api/api_client.dart';
import '../../../../shared/data/api/server_exception.dart';
import '../../../../shared/data/logger/app_logger.dart';
import '../../../../shared/domain/models/app_failure.dart';
import '../../domain/repositories/user_repository.dart';
import '../models/user_models.dart';

@LazySingleton(as: UserRepository)
class UserRepositoryImpl implements UserRepository {
  final ApiClient _api;
  final AppLogger _logger;

  static const _userBase = '/api/v1/messenger/user';
  static const _usersBase = '/api/v1/messenger/users';

  UserRepositoryImpl(this._api, this._logger);

  // ─── Profile ─────────────────────────────────────────────────

  @override
  Future<Either<AppFailure, UserProfileModel>> getUserProfile(String userId) async {
    try {
      final response = await _api.get('$_userBase/id/$userId');
      final map = _extractMap(response);
      return Right(UserProfileModel.fromJson(map));
    } catch (e) {
      _logger.error('getUserProfile failed (userId=$userId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, UserProfileModel>> getUserStatus(String userId) async {
    try {
      final response = await _api.get('$_userBase/id/$userId/status');
      final map = _extractMap(response);
      // Status endpoint returns { user_id, is_online, last_seen_at }
      // We build a minimal UserProfileModel to merge with profile data
      return Right(UserProfileModel(
        id: map['user_id'] as String? ?? userId,
        username: '',
        displayName: '',
        isOnline: map['is_online'] as bool? ?? false,
        lastSeenAt: map['last_seen_at'] as String?,
      ));
    } catch (e) {
      _logger.error('getUserStatus failed (userId=$userId)', e);
      return Left(_handleError(e));
    }
  }

  // ─── Contacts ─────────────────────────────────────────────────

  @override
  Future<Either<AppFailure, List<ContactModel>>> getContacts() async {
    try {
      final response = await _api.get('$_usersBase/contacts');
      final List raw = response as List? ?? [];
      return Right(raw
          .map((e) => ContactModel.fromJson(Map<String, dynamic>.from(e as Map)))
          .toList());
    } catch (e) {
      _logger.error('getContacts failed', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> addContact(
    String userId, {
    String? customName,
  }) async {
    try {
      final body = <String, dynamic>{'user_id': userId};
      if (customName != null && customName.isNotEmpty) body['custom_name'] = customName;
      await _api.post('$_usersBase/contacts', data: body);
      return const Right(unit);
    } catch (e) {
      _logger.error('addContact failed (userId=$userId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> removeContact(String userId) async {
    try {
      await _api.delete('$_usersBase/contacts/$userId');
      return const Right(unit);
    } catch (e) {
      _logger.error('removeContact failed (userId=$userId)', e);
      return Left(_handleError(e));
    }
  }

  // ─── Block list ───────────────────────────────────────────────

  @override
  Future<Either<AppFailure, List<BlockedUserModel>>> getBlockedUsers() async {
    try {
      final response = await _api.get('$_usersBase/blocked');
      final List raw = response as List? ?? [];
      final users = raw
          .whereType<Map<String, dynamic>>()
          .map(BlockedUserModel.fromJson)
          .toList();
      return Right(users);
    } catch (e) {
      _logger.error('getBlockedUsers failed', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> blockUser(String userId) async {
    try {
      await _api.post('$_usersBase/block', data: {'user_id': userId});
      return const Right(unit);
    } catch (e) {
      _logger.error('blockUser failed (userId=$userId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> unblockUser(String userId) async {
    try {
      await _api.post('$_usersBase/unblock', data: {'user_id': userId});
      return const Right(unit);
    } catch (e) {
      _logger.error('unblockUser failed (userId=$userId)', e);
      return Left(_handleError(e));
    }
  }

  // ─── Helpers ──────────────────────────────────────────────────

  Map<String, dynamic> _extractMap(dynamic response) {
    if (response is Map<String, dynamic>) {
      // Unwrap common wrapper keys
      for (final key in ['user', 'profile', 'data']) {
        if (response.containsKey(key) && response[key] is Map) {
          return Map<String, dynamic>.from(response[key] as Map);
        }
      }
      return response;
    }
    return {};
  }

  AppFailure _handleError(dynamic e) {
    if (e is ServerException) return ServerFailure(message: e.message, code: e.errorCode);
    return const UnknownFailure();
  }
}
