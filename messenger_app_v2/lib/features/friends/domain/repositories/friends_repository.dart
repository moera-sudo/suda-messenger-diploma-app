import 'package:dartz/dartz.dart';

import '../../../../shared/domain/models/app_failure.dart';
import '../../data/models/friend_models.dart';

abstract class FriendsRepository {
  /// POST /api/v1/messenger/users/friends/requests {user_id}
  Future<Either<AppFailure, Unit>> sendRequest(String userId);

  /// DELETE /api/v1/messenger/users/friends/requests/{id}
  Future<Either<AppFailure, Unit>> cancelRequest(String requestId);

  /// POST /api/v1/messenger/users/friends/requests/{id}/accept
  Future<Either<AppFailure, Unit>> acceptRequest(String requestId);

  /// POST /api/v1/messenger/users/friends/requests/{id}/reject
  Future<Either<AppFailure, Unit>> rejectRequest(String requestId);

  /// GET /api/v1/messenger/users/friends
  Future<Either<AppFailure, List<FriendModel>>> getFriends();

  /// DELETE /api/v1/messenger/users/friends/{user_id}
  Future<Either<AppFailure, Unit>> unfriend(String userId);

  /// GET /api/v1/messenger/users/friends/requests (both directions)
  Future<Either<AppFailure, List<FriendRequestModel>>> getRequests();

  /// GET /api/v1/messenger/users/friends/{user_id}/status
  Future<Either<AppFailure, FriendStatusModel>> getStatus(String userId);
}
