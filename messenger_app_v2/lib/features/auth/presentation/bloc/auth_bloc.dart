import 'dart:async';

import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:injectable/injectable.dart';
import '../../../../shared/data/auth/auth_event_bus.dart';
import '../../../../shared/domain/models/app_failure.dart';
import '../../domain/repositories/auth_repository.dart';

// --- EVENTS ---
abstract class AuthEvent extends Equatable {
  @override
  List<Object?> get props => [];
}

class AuthLoginRequested extends AuthEvent {
  final String email;
  final String password;
  AuthLoginRequested(this.email, this.password);
}

class AuthRegisterRequested extends AuthEvent {
  final String username;
  final String displayName;
  final String email;
  final String password;

  AuthRegisterRequested({
    required this.username,
    required this.displayName,
    required this.email,
    required this.password,
  });
}

class AuthVerifyRequested extends AuthEvent {
  final String email;
  final String code;
  AuthVerifyRequested(this.email, this.code);
}

class AuthForgotPasswordRequested extends AuthEvent {
  final String email;
  AuthForgotPasswordRequested(this.email);
}

class AuthResetPasswordRequested extends AuthEvent {
  final String email;
  final String code;
  final String newPassword;

  AuthResetPasswordRequested({
    required this.email,
    required this.code,
    required this.newPassword,
  });
}

class AuthCheckRequested extends AuthEvent {}

/// Fired by the [AuthEventBus] when the Dio layer detects an unrecoverable
/// deauthorization (invalid refresh token / revoked session). Flips the app to
/// the unauthenticated state so the global listener redirects to Welcome.
class AuthForcedLogout extends AuthEvent {}

/// Resend email verification code.
/// NOTE: Requires POST /auth/resend-verify endpoint on backend (not yet available).
/// Currently emits [AuthStatus.forgotCodeSent] so the UI can reset its timer.
class AuthResendCodeRequested extends AuthEvent {
  final String email;
  AuthResendCodeRequested(this.email);
  @override
  List<Object?> get props => [email];
}

// --- STATE ---
enum AuthStatus {
  initial,
  loading,
  success,
  forgotCodeSent,
  passwordReset,
  failure,
  authenticated,
  unauthenticated,
}

class AuthState extends Equatable {
  final AuthStatus status;
  final AppFailure? failure;
  final String? emailToVerify; // Храним email, если нужна верификация

  const AuthState({
    this.status = AuthStatus.initial,
    this.failure,
    this.emailToVerify,
  });

  // Метод copyWith ОБЯЗАТЕЛЕН, если ты используешь state.copyWith(...)
  AuthState copyWith({
    AuthStatus? status,
    AppFailure? failure,
    String? emailToVerify,
  }) {
    return AuthState(
      status: status ?? this.status,
      failure: failure ?? this.failure,
      emailToVerify: emailToVerify ?? this.emailToVerify,
    );
  }

  @override
  List<Object?> get props => [status, failure, emailToVerify];
}

// --- BLOC ---
@injectable
class AuthBloc extends Bloc<AuthEvent, AuthState> {
  final AuthRepository _repo;
  final AuthEventBus _authEvents;
  StreamSubscription<void>? _forcedLogoutSub;

  AuthBloc(this._repo, this._authEvents) : super(const AuthState()) {
    on<AuthLoginRequested>(_onLogin);
    on<AuthRegisterRequested>(_onRegister);
    on<AuthVerifyRequested>(_onVerify);
    on<AuthForgotPasswordRequested>(_onForgotPassword);
    on<AuthResetPasswordRequested>(_onResetPassword);
    on<AuthCheckRequested>(_onAuthCheck);
    on<AuthResendCodeRequested>(_onResendCode);
    on<AuthForcedLogout>(_onForcedLogout);

    // Bridge the Dio-layer deauthorization signal into BLoC state.
    _forcedLogoutSub =
        _authEvents.onForcedLogout.listen((_) => add(AuthForcedLogout()));
  }

  void _onForcedLogout(AuthForcedLogout event, Emitter<AuthState> emit) {
    if (state.status == AuthStatus.unauthenticated) return;
    emit(state.copyWith(status: AuthStatus.unauthenticated));
  }

  @override
  Future<void> close() {
    _forcedLogoutSub?.cancel();
    return super.close();
  }

  Future<void> _onAuthCheck(
    AuthCheckRequested event,
    Emitter<AuthState> emit,
  ) async {
    // Этот метод должен быть в AuthRepository
    final isAuth = await _repo.isAuthenticated();
    if (isAuth) {
      emit(state.copyWith(status: AuthStatus.authenticated));
    } else {
      emit(state.copyWith(status: AuthStatus.unauthenticated));
    }
  }

  Future<void> _onLogin(
    AuthLoginRequested event,
    Emitter<AuthState> emit,
  ) async {
    emit(state.copyWith(status: AuthStatus.loading));
    final result = await _repo.login(
      email: event.email,
      password: event.password,
    );

    result.fold(
      (failure) =>
          emit(state.copyWith(status: AuthStatus.failure, failure: failure)),
      (_) => emit(state.copyWith(status: AuthStatus.success)),
    );
  }

  Future<void> _onRegister(
    AuthRegisterRequested event,
    Emitter<AuthState> emit,
  ) async {
    emit(state.copyWith(status: AuthStatus.loading));
    final result = await _repo.register(
      username: event.username,
      displayName: event.displayName,
      email: event.email,
      password: event.password,
    );

    result.fold(
      (failure) =>
          emit(state.copyWith(status: AuthStatus.failure, failure: failure)),
      (_) => emit(
        state.copyWith(status: AuthStatus.success, emailToVerify: event.email),
      ),
    );
  }

  Future<void> _onVerify(
    AuthVerifyRequested event,
    Emitter<AuthState> emit,
  ) async {
    emit(state.copyWith(status: AuthStatus.loading));
    // Вызываем метод verify из репозитория
    // (Убедись, что он там есть: Future<Either<AppFailure, Unit>> verify(...))
    final result = await _repo.verify(email: event.email, code: event.code);

    result.fold(
      (failure) =>
          emit(state.copyWith(status: AuthStatus.failure, failure: failure)),
      (_) => emit(state.copyWith(status: AuthStatus.success)),
    );
  }

  Future<void> _onForgotPassword(
    AuthForgotPasswordRequested event,
    Emitter<AuthState> emit,
  ) async {
    final email = event.email.trim().toLowerCase();
    if (email.isEmpty) {
      emit(
        state.copyWith(
          status: AuthStatus.failure,
          failure: const ServerFailure(
            message: 'Email is required',
            code: 'BAD_REQUEST',
          ),
        ),
      );
      return;
    }

    emit(state.copyWith(status: AuthStatus.loading));
    final result = await _repo.forgotPassword(email: email);

    result.fold(
      (failure) =>
          emit(state.copyWith(status: AuthStatus.failure, failure: failure)),
      (_) => emit(state.copyWith(status: AuthStatus.forgotCodeSent)),
    );
  }

  Future<void> _onResetPassword(
    AuthResetPasswordRequested event,
    Emitter<AuthState> emit,
  ) async {
    final email = event.email.trim().toLowerCase();
    final code = event.code.trim();
    final newPassword = event.newPassword.trim();

    if (email.isEmpty || code.length != 6 || newPassword.length < 6) {
      emit(
        state.copyWith(
          status: AuthStatus.failure,
          failure: const ServerFailure(
            message: 'Invalid reset payload',
            code: 'BAD_REQUEST',
          ),
        ),
      );
      return;
    }

    emit(state.copyWith(status: AuthStatus.loading));
    final result = await _repo.resetPassword(
      email: email,
      code: code,
      newPassword: newPassword,
    );

    result.fold(
      (failure) =>
          emit(state.copyWith(status: AuthStatus.failure, failure: failure)),
      (_) => emit(state.copyWith(status: AuthStatus.passwordReset)),
    );
  }

  // NOTE: No /auth/resend-verify endpoint on backend yet.
  // Emits forgotCodeSent so the UI resets its countdown timer.
  // Replace body with actual repo call when endpoint is added.
  Future<void> _onResendCode(
    AuthResendCodeRequested event,
    Emitter<AuthState> emit,
  ) async {
    emit(state.copyWith(status: AuthStatus.forgotCodeSent));
  }
}
