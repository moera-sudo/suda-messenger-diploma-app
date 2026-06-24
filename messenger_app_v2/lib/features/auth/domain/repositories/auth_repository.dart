import 'package:dartz/dartz.dart';
import '../../../../shared/domain/models/app_failure.dart';

abstract class AuthRepository {
  Future<Either<AppFailure, Unit>> register({
    required String username,
    required String displayName,
    required String email,
    required String password,
  });

  Future<Either<AppFailure, Unit>> login({
    required String email,
    required String password,
  });

  Future<Either<AppFailure, Unit>> verify({
    required String email,
    required String code,
  });

  Future<Either<AppFailure, Unit>> forgotPassword({required String email});

  Future<Either<AppFailure, Unit>> resetPassword({
    required String email,
    required String code,
    required String newPassword,
  });

  Future<bool> isAuthenticated();
}
