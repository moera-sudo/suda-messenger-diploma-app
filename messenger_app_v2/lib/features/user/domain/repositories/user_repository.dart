import 'package:dartz/dartz.dart';

import '../../../../shared/domain/models/app_failure.dart';
import '../../data/models/user_models.dart';

abstract class UserRepository {
  /// GET /api/v1/messenger/user/id/{userId}
  Future<Either<AppFailure, UserProfileModel>> getUserProfile(String userId);

  /// GET /api/v1/messenger/user/id/{userId}/status
  Future<Either<AppFailure, UserProfileModel>> getUserStatus(String userId);

  // ─── Contacts ────────────────────────────────────────────────

  /// GET /api/v1/messenger/users/contacts
  Future<Either<AppFailure, List<ContactModel>>> getContacts();

  /// POST /api/v1/messenger/users/contacts
  Future<Either<AppFailure, Unit>> addContact(String userId, {String? customName});

  /// DELETE /api/v1/messenger/users/contacts/{userId}
  Future<Either<AppFailure, Unit>> removeContact(String userId);

  // ─── Block list ───────────────────────────────────────────────

  /// GET /api/v1/messenger/users/blocked → returns array of blocked user objects
  Future<Either<AppFailure, List<BlockedUserModel>>> getBlockedUsers();

  /// POST /api/v1/messenger/users/block { user_id }
  Future<Either<AppFailure, Unit>> blockUser(String userId);

  /// POST /api/v1/messenger/users/unblock { user_id }
  Future<Either<AppFailure, Unit>> unblockUser(String userId);
}
