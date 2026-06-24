import 'dart:io';

import 'package:dartz/dartz.dart';
import 'package:device_info_plus/device_info_plus.dart';
import 'package:injectable/injectable.dart';
import 'package:messenger_app_v2/app/config/const.dart';
import '../../../../shared/data/api/api_client.dart';
import '../../../../shared/data/api/server_exception.dart';
import '../../../../shared/data/storage/secure_storage_client.dart';
import '../../../../shared/domain/models/app_failure.dart';
import '../../../../shared/domain/models/auth_token_pair.dart';
import '../../domain/repositories/auth_repository.dart';
import '../models/auth_request_models.dart';

@LazySingleton(as: AuthRepository)
class AuthRepositoryImpl implements AuthRepository {
  final ApiClient _api;
  final SecureStorageClient _storage;
  static const _basePath = '${Consts.PrefixMessenger}/auth';

  AuthRepositoryImpl(this._api, this._storage);

  @override
  Future<Either<AppFailure, Unit>> register({
    required String username,
    required String displayName,
    required String email,
    required String password,
  }) async {
    try {
      final req = RegisterRequest(
        username: username,
        displayName: displayName,
        email: email,
        password: password,
      );
      await _api.post('$_basePath/register', data: req.toJson());
      return const Right(unit);
    } on ServerException catch (e) {
      return Left(ServerFailure(message: e.message, code: e.errorCode));
    } catch (e) {
      return const Left(UnknownFailure());
    }
  }

  @override
  Future<Either<AppFailure, Unit>> login({
    required String email,
    required String password,
  }) async {
    try {
      String userAgent = 'FlutterApp';
      try {
        final plugin = DeviceInfoPlugin();
        if (Platform.isAndroid) {
          final info = await plugin.androidInfo;
          userAgent = 'Android ${info.version.release} / ${info.manufacturer} ${info.model}';
        } else if (Platform.isIOS) {
          final info = await plugin.iosInfo;
          userAgent = 'iOS ${info.systemVersion} / ${info.name}';
        }
      } catch (_) {}

      final req = LoginRequest(email: email, password: password, userAgent: userAgent);
      final response = await _api.post('$_basePath/login', data: req.toJson());
      final tokens = AuthTokenPair.fromJson(response);
      await _storage.saveTokens(
        access: tokens.accessToken,
        refresh: tokens.refreshToken,
      );
      return const Right(unit);
    } on ServerException catch (e) {
      return Left(ServerFailure(message: e.message, code: e.errorCode));
    } catch (e) {
      return const Left(UnknownFailure());
    }
  }

  @override
  Future<Either<AppFailure, Unit>> verify({
    required String email,
    required String code,
  }) async {
    try {
      final req = VerifyRequest(email: email, code: code);
      await _api.post('$_basePath/verify', data: req.toJson());
      return const Right(unit);
    } on ServerException catch (e) {
      return Left(ServerFailure(message: e.message, code: e.errorCode));
    } catch (e) {
      return const Left(UnknownFailure());
    }
  }

  @override
  Future<Either<AppFailure, Unit>> forgotPassword({
    required String email,
  }) async {
    try {
      final req = ForgotPasswordRequest(email: email.trim().toLowerCase());
      await _api.post('$_basePath/forgot-password', data: req.toJson());
      return const Right(unit);
    } on ServerException catch (e) {
      return Left(ServerFailure(message: e.message, code: e.errorCode));
    } catch (_) {
      return const Left(UnknownFailure());
    }
  }

  @override
  Future<Either<AppFailure, Unit>> resetPassword({
    required String email,
    required String code,
    required String newPassword,
  }) async {
    try {
      final req = ResetPasswordRequest(
        email: email.trim().toLowerCase(),
        code: code.trim(),
        newPassword: newPassword.trim(),
      );
      await _api.post('$_basePath/reset-password', data: req.toJson());
      return const Right(unit);
    } on ServerException catch (e) {
      return Left(ServerFailure(message: e.message, code: e.errorCode));
    } catch (_) {
      return const Left(UnknownFailure());
    }
  }

  @override
  Future<bool> isAuthenticated() async {
    final token = await _storage.getAccessToken();
    return token != null && token.isNotEmpty;
  }
}
