import 'package:dartz/dartz.dart';

import '../../../../shared/domain/models/app_failure.dart';
import '../../data/models/my_profile.dart';
import '../../data/models/session_model.dart';

abstract class ProfileRepository {
  Future<Either<AppFailure, MyProfile>> getMe();

  Future<Either<AppFailure, Unit>> updateProfile(UpdateMyProfilePayload payload);

  Future<Either<AppFailure, String>> uploadAvatar({
    required String filePath,
    String? fileName,
  });

  Future<Either<AppFailure, Unit>> logout();

  /// PUT /api/v1/messenger/user/password
  Future<Either<AppFailure, Unit>> changePassword({
    required String oldPassword,
    required String newPassword,
  });

  /// GET /api/v1/messenger/user/sessions
  Future<Either<AppFailure, List<SessionModel>>> getSessions();

  /// DELETE /api/v1/messenger/user/sessions/{id}
  Future<Either<AppFailure, Unit>> deleteSession(String sessionId);

  /// POST /api/v1/messenger/user/sessions/terminate-others
  Future<Either<AppFailure, Unit>> terminateOtherSessions(String currentRefreshToken);
}
