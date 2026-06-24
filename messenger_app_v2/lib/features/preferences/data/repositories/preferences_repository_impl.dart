import 'package:dartz/dartz.dart';
import 'package:injectable/injectable.dart';

import '../../../../shared/data/api/api_client.dart';
import '../../../../shared/data/api/server_exception.dart';
import '../../../../shared/data/logger/app_logger.dart';
import '../../../../shared/domain/models/app_failure.dart';
import '../../domain/repositories/preferences_repository.dart';
import '../models/user_preferences.dart';

@LazySingleton(as: PreferencesRepository)
class PreferencesRepositoryImpl implements PreferencesRepository {
  final ApiClient _api;
  final AppLogger _logger;

  static const _base = '/api/v1/messenger/user/preferences';

  PreferencesRepositoryImpl(this._api, this._logger);

  @override
  Future<Either<AppFailure, UserPreferences>> getPreferences() async {
    try {
      final response = await _api.get(_base);
      final prefs = UserPreferences.fromJson(Map<String, dynamic>.from(response as Map));
      _logger.debug('Loaded preferences: lang=${prefs.language}, notif=${prefs.notificationsEnabled}');
      return Right(prefs);
    } catch (e) {
      _logger.error('getPreferences failed', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> updatePreferences(Map<String, dynamic> patch) async {
    try {
      await _api.put(_base, data: patch);
      _logger.info('Updated preferences: ${patch.keys.join(', ')}');
      return const Right(unit);
    } catch (e) {
      _logger.error('updatePreferences failed (fields=${patch.keys.join(', ')})', e);
      return Left(_handleError(e));
    }
  }

  AppFailure _handleError(dynamic e) {
    if (e is ServerException) return ServerFailure(message: e.message, code: e.errorCode);
    return const UnknownFailure();
  }
}
