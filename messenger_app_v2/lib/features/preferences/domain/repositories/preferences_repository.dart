import 'package:dartz/dartz.dart';

import '../../../../shared/domain/models/app_failure.dart';
import '../../data/models/user_preferences.dart';

abstract class PreferencesRepository {
  /// GET /api/v1/messenger/user/preferences
  Future<Either<AppFailure, UserPreferences>> getPreferences();

  /// PUT /api/v1/messenger/user/preferences — partial update (only passed fields).
  Future<Either<AppFailure, Unit>> updatePreferences(Map<String, dynamic> patch);
}
