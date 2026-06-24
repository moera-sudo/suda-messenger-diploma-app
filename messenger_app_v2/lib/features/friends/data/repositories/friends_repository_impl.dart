import 'package:dartz/dartz.dart';
import 'package:injectable/injectable.dart';

import '../../../../shared/data/api/api_client.dart';
import '../../../../shared/data/api/server_exception.dart';
import '../../../../shared/data/logger/app_logger.dart';
import '../../../../shared/domain/models/app_failure.dart';
import '../../domain/repositories/friends_repository.dart';
import '../models/friend_models.dart';

@LazySingleton(as: FriendsRepository)
class FriendsRepositoryImpl implements FriendsRepository {
  final ApiClient _api;
  final AppLogger _logger;

  static const _base = '/api/v1/messenger/users/friends';

  FriendsRepositoryImpl(this._api, this._logger);

  @override
  Future<Either<AppFailure, List<FriendModel>>> getFriends() async {
    try {
      final res = await _api.get(_base);
      final list = (res as List).map((e) => FriendModel.fromJson(Map<String, dynamic>.from(e as Map))).toList();
      _logger.debug('getFriends: ${list.length} friends');
      return Right(list);
    } catch (e) {
      _logger.error('getFriends failed', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, List<FriendRequestModel>>> getRequests() async {
    try {
      final res = await _api.get('$_base/requests');
      final list = (res as List).map((e) => FriendRequestModel.fromJson(Map<String, dynamic>.from(e as Map))).toList();
      _logger.debug('getRequests: ${list.length} requests');
      return Right(list);
    } catch (e) {
      _logger.error('getRequests failed', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, FriendStatusModel>> getStatus(String userId) async {
    try {
      final res = await _api.get('$_base/$userId/status');
      return Right(FriendStatusModel.fromJson(Map<String, dynamic>.from(res as Map)));
    } catch (e) {
      _logger.error('getStatus failed (userId=$userId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> sendRequest(String userId) async {
    try {
      await _api.post('$_base/requests', data: {'user_id': userId});
      _logger.info('Friend request sent to $userId');
      return const Right(unit);
    } catch (e) {
      _logger.error('sendRequest failed (userId=$userId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> cancelRequest(String requestId) async {
    try {
      await _api.delete('$_base/requests/$requestId');
      _logger.info('Friend request $requestId cancelled');
      return const Right(unit);
    } catch (e) {
      _logger.error('cancelRequest failed ($requestId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> acceptRequest(String requestId) async {
    try {
      await _api.post('$_base/requests/$requestId/accept', data: {});
      _logger.info('Friend request $requestId accepted');
      return const Right(unit);
    } catch (e) {
      _logger.error('acceptRequest failed ($requestId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> rejectRequest(String requestId) async {
    try {
      await _api.post('$_base/requests/$requestId/reject', data: {});
      _logger.info('Friend request $requestId rejected');
      return const Right(unit);
    } catch (e) {
      _logger.error('rejectRequest failed ($requestId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> unfriend(String userId) async {
    try {
      await _api.delete('$_base/$userId');
      _logger.info('Unfriended $userId');
      return const Right(unit);
    } catch (e) {
      _logger.error('unfriend failed (userId=$userId)', e);
      return Left(_handleError(e));
    }
  }

  AppFailure _handleError(dynamic e) {
    if (e is ServerException) return ServerFailure(message: e.message, code: e.errorCode);
    return const UnknownFailure();
  }
}
