import 'package:dartz/dartz.dart';
import 'package:injectable/injectable.dart';

import '../../../../features/media/domain/repositories/media_repository.dart';
import '../../../../shared/data/api/api_client.dart';
import '../../../../shared/data/api/server_exception.dart';
import '../../../../shared/data/logger/app_logger.dart';
import '../../../../shared/data/storage/secure_storage_client.dart';
import '../../../../shared/domain/models/app_failure.dart';
import '../../domain/repositories/profile_repository.dart';
import '../models/my_profile.dart';
import '../models/session_model.dart';

@LazySingleton(as: ProfileRepository)
class ProfileRepositoryImpl implements ProfileRepository {
  final ApiClient _api;
  final SecureStorageClient _storage;
  final MediaRepository _mediaRepo;
  final AppLogger _logger;

  static const _basePath = '/api/v1/messenger/user';

  const ProfileRepositoryImpl(
    this._api,
    this._storage,
    this._mediaRepo,
    this._logger,
  );

  @override
  Future<Either<AppFailure, MyProfile>> getMe() async {
    final paths = <String>[
      '$_basePath/me',
      '/api/v1/messenger/users/me',
      '/api/v1/user/me',
      '/api/v1/users/me',
    ];

    ServerException? lastServerError;

    for (final path in paths) {
      try {
        final response = await _api.get(path);
        final raw = _extractMap(response);
        if (raw == null) {
          continue;
        }
        return Right(await _withResolvedAvatar(_normalizeProfile(raw), raw));
      } on ServerException catch (e) {
        lastServerError = e;
        if (e.statusCode == 404) {
          continue;
        }
        return Left(ServerFailure(message: e.message, code: e.errorCode));
      } catch (_) {
        continue;
      }
    }

    if (lastServerError != null) {
      return Left(
        ServerFailure(
          message: lastServerError.message,
          code: lastServerError.errorCode,
        ),
      );
    }

    return const Left(UnknownFailure(message: 'Failed to load profile'));
  }

  @override
  Future<Either<AppFailure, Unit>> updateProfile(
    UpdateMyProfilePayload payload,
  ) async {
    final paths = <String>[
      '$_basePath/profile',
      '$_basePath/me',
      _basePath,
      '/api/v1/user/profile',
    ];

    ServerException? lastServerError;

    for (final path in paths) {
      try {
        await _api.put(path, data: payload.toJson());
        return const Right(unit);
      } on ServerException catch (e) {
        if (e.statusCode == 404) {
          lastServerError = e;
          continue;
        }
        if (e.statusCode == 405) {
          try {
            await _api.post(path, data: payload.toJson());
            return const Right(unit);
          } on ServerException catch (nested) {
            lastServerError = nested;
            if (nested.statusCode == 404) {
              continue;
            }
            return Left(
              ServerFailure(message: nested.message, code: nested.errorCode),
            );
          } catch (_) {
            continue;
          }
        }
        return Left(ServerFailure(message: e.message, code: e.errorCode));
      } catch (_) {
        continue;
      }
    }

    if (lastServerError != null) {
      return Left(
        ServerFailure(
          message: lastServerError.message,
          code: lastServerError.errorCode,
        ),
      );
    }

    return const Left(UnknownFailure(message: 'Profile update failed'));
  }

  @override
  Future<Either<AppFailure, String>> uploadAvatar({
    required String filePath,
    String? fileName,
  }) async {
    // New contract (Frontend.md §5/§6): upload file to media-service, then bind
    // the media_id to the user via PUT /user/avatar (the old multipart POST is
    // now 405). Returns a presigned view URL for immediate display.
    try {
      // Steps 1–3 — upload bytes through media-service (kind=AVATAR, public).
      final uploadResult = await _mediaRepo.uploadFile(
        filePath: filePath,
        mediaType: 'AVATAR',
        filename: fileName,
        isPrivate: false,
      );
      final mediaId = uploadResult.fold((_) => null, (id) => id);
      if (mediaId == null) {
        _logger.error('uploadAvatar: media upload failed');
        return uploadResult.fold(
          (f) => Left(f),
          (_) => const Left(UnknownFailure(message: 'Avatar upload failed')),
        );
      }

      // Step 4 — bind media_id to the user (PUT + JSON).
      await _api.put('$_basePath/avatar', data: {'media_id': mediaId});
      _logger.info('Avatar bound to user (media_id=$mediaId)');

      // Resolve a presigned view URL for immediate display.
      final urlResult = await _mediaRepo.getMediaUrl(mediaId);
      return urlResult.fold(
        (_) => const Left(
          UnknownFailure(message: 'Avatar updated, but URL missing'),
        ),
        (url) => Right(url),
      );
    } on ServerException catch (e) {
      _logger.error('uploadAvatar failed (bind step)', e);
      return Left(ServerFailure(message: e.message, code: e.errorCode));
    } catch (e) {
      _logger.error('uploadAvatar failed', e);
      return const Left(UnknownFailure(message: 'Avatar upload failed'));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> logout() async {
    final refreshToken = await _storage.getRefreshToken();
    final payload = <String, dynamic>{
      if (refreshToken != null && refreshToken.isNotEmpty)
        'refresh_token': refreshToken,
    };

    final paths = <String>[
      '$_basePath/logout',
      '/api/v1/messenger/auth/logout',
      '/api/v1/user/logout',
    ];

    ServerException? lastServerError;

    for (final path in paths) {
      try {
        await _api.post(path, data: payload);
        await _storage.clearTokens();
        return const Right(unit);
      } on ServerException catch (e) {
        lastServerError = e;
        if (e.statusCode == 404) {
          continue;
        }
        break;
      } catch (_) {
        continue;
      }
    }

    // Best effort local logout even if server route failed.
    await _storage.clearTokens();

    if (lastServerError != null) {
      return Left(
        ServerFailure(
          message: lastServerError.message,
          code: lastServerError.errorCode,
        ),
      );
    }

    return const Right(unit);
  }

  /// Backend returns `avatar_media_id` (a UUID), not a URL. Resolve it to a
  /// presigned view URL so [SudaAvatar] can render it. Failures are non-fatal —
  /// the avatar just falls back to initials.
  Future<MyProfile> _withResolvedAvatar(
    MyProfile profile,
    Map<String, dynamic> raw,
  ) async {
    if (profile.avatarUrl != null && profile.avatarUrl!.isNotEmpty) {
      return profile;
    }
    final mediaId = _readNullableString(raw, const [
      'avatar_media_id',
      'avatarMediaId',
    ]);
    if (mediaId == null) return profile;

    final result = await _mediaRepo.getMediaUrl(mediaId);
    return result.fold(
      (failure) {
        _logger.warning('Could not resolve avatar URL: ${failure.message}');
        return profile;
      },
      (url) => url.isEmpty ? profile : profile.copyWith(avatarUrl: url),
    );
  }

  Map<String, dynamic>? _extractMap(dynamic response) {
    if (response is Map<String, dynamic>) {
      final nested =
          response['user'] ?? response['profile'] ?? response['data'];
      if (nested is Map<String, dynamic>) {
        return nested;
      }
      return response;
    }

    if (response is Map) {
      final normalized = response.map(
        (key, value) => MapEntry(key.toString(), value),
      );
      final nested =
          normalized['user'] ?? normalized['profile'] ?? normalized['data'];
      if (nested is Map) {
        return nested.map((key, value) => MapEntry(key.toString(), value));
      }
      return normalized;
    }

    return null;
  }

  MyProfile _normalizeProfile(Map<String, dynamic> source) {
    final username = _readString(source, const [
      'username',
      'user_name',
      'login',
    ], fallback: 'username');

    return MyProfile(
      id: _readString(source, const ['id', 'user_id'], fallback: ''),
      username: username,
      email: _readString(source, const ['email', 'mail'], fallback: ''),
      displayName: _readString(source, const [
        'display_name',
        'displayName',
        'name',
      ], fallback: username),
      firstName: _readString(source, const [
        'first_name',
        'firstName',
      ], fallback: ''),
      lastName: _readString(source, const [
        'last_name',
        'lastName',
      ], fallback: ''),
      bio: _readString(source, const ['bio', 'about'], fallback: ''),
      avatarUrl: _readNullableString(source, const [
        'avatar_url',
        'avatarUrl',
        'avatar',
      ]),
      isOnline: _readBool(source, const ['is_online', 'isOnline', 'online']),
    );
  }

  String _readString(
    Map<String, dynamic> source,
    List<String> keys, {
    required String fallback,
  }) {
    for (final key in keys) {
      final raw = source[key];
      if (raw is String && raw.trim().isNotEmpty) {
        return raw.trim();
      }
      if (raw != null) {
        final text = raw.toString().trim();
        if (text.isNotEmpty) {
          return text;
        }
      }
    }
    return fallback;
  }

  String? _readNullableString(Map<String, dynamic> source, List<String> keys) {
    final value = _readString(source, keys, fallback: '');
    return value.isEmpty ? null : value;
  }

  bool _readBool(Map<String, dynamic> source, List<String> keys) {
    for (final key in keys) {
      final raw = source[key];
      if (raw is bool) {
        return raw;
      }
      if (raw is num) {
        return raw != 0;
      }
      if (raw is String) {
        final value = raw.trim().toLowerCase();
        if (value == '1' || value == 'true' || value == 'online') {
          return true;
        }
        if (value == '0' || value == 'false' || value == 'offline') {
          return false;
        }
      }
    }
    return false;
  }

  @override
  Future<Either<AppFailure, Unit>> changePassword({
    required String oldPassword,
    required String newPassword,
  }) async {
    try {
      await _api.put('$_basePath/password', data: {
        'old_password': oldPassword,
        'new_password': newPassword,
      });
      return const Right(unit);
    } on ServerException catch (e) {
      return Left(ServerFailure(message: e.message, code: e.errorCode));
    } catch (_) {
      return const Left(UnknownFailure());
    }
  }

  @override
  Future<Either<AppFailure, List<SessionModel>>> getSessions() async {
    try {
      final res = await _api.get('$_basePath/sessions');
      final list = (res as List)
          .map((e) => SessionModel.fromJson(Map<String, dynamic>.from(e as Map)))
          .toList();
      return Right(list);
    } on ServerException catch (e) {
      return Left(ServerFailure(message: e.message, code: e.errorCode));
    } catch (_) {
      return const Left(UnknownFailure());
    }
  }

  @override
  Future<Either<AppFailure, Unit>> deleteSession(String sessionId) async {
    try {
      await _api.delete('$_basePath/sessions/$sessionId');
      return const Right(unit);
    } on ServerException catch (e) {
      return Left(ServerFailure(message: e.message, code: e.errorCode));
    } catch (_) {
      return const Left(UnknownFailure());
    }
  }

  @override
  Future<Either<AppFailure, Unit>> terminateOtherSessions(
      String currentRefreshToken) async {
    try {
      await _api.post('$_basePath/sessions/terminate-others',
          data: {'current_refresh_token': currentRefreshToken});
      return const Right(unit);
    } on ServerException catch (e) {
      return Left(ServerFailure(message: e.message, code: e.errorCode));
    } catch (_) {
      return const Left(UnknownFailure());
    }
  }
}
