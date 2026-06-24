import 'dart:async';

import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:injectable/injectable.dart';

import '../../../../shared/data/logger/app_logger.dart';
import '../../data/models/user_models.dart';
import '../../domain/repositories/user_repository.dart';

// --- EVENTS ---
abstract class UserProfileEvent extends Equatable {
  @override
  List<Object> get props => [];
}

class LoadUserProfile extends UserProfileEvent {
  final String userId;
  LoadUserProfile(this.userId);
  @override
  List<Object> get props => [userId];
}

class BlockUserEvent extends UserProfileEvent {}

class UnblockUserEvent extends UserProfileEvent {}

// --- STATE ---
enum UserProfileStatus { initial, loading, ready, failure }

class UserProfileState extends Equatable {
  final UserProfileStatus status;
  final UserProfileModel? profile;
  final bool isBlocked;
  final String? error;

  const UserProfileState({
    this.status = UserProfileStatus.initial,
    this.profile,
    this.isBlocked = false,
    this.error,
  });

  UserProfileState copyWith({
    UserProfileStatus? status,
    UserProfileModel? profile,
    bool? isBlocked,
    String? error,
  }) {
    return UserProfileState(
      status: status ?? this.status,
      profile: profile ?? this.profile,
      isBlocked: isBlocked ?? this.isBlocked,
      error: error ?? this.error,
    );
  }

  @override
  List<Object?> get props => [status, profile?.id, isBlocked, error];
}

// --- BLOC ---
@injectable
class UserProfileBloc extends Bloc<UserProfileEvent, UserProfileState> {
  final UserRepository _repo;
  final AppLogger _logger;

  String _userId = '';

  UserProfileBloc(this._repo, this._logger) : super(const UserProfileState()) {
    on<LoadUserProfile>(_onLoadUserProfile);
    on<BlockUserEvent>(_onBlockUser);
    on<UnblockUserEvent>(_onUnblockUser);
  }

  Future<void> _onLoadUserProfile(
    LoadUserProfile event,
    Emitter<UserProfileState> emit,
  ) async {
    _userId = event.userId;
    emit(state.copyWith(status: UserProfileStatus.loading));
    _logger.info('Loading user profile: $_userId');

    // Fetch profile, status, and block list in parallel
    final profileFuture = _repo.getUserProfile(_userId);
    final statusFuture = _repo.getUserStatus(_userId);
    final blockedFuture = _repo.getBlockedUsers();

    final profileResult = await profileFuture;
    final statusResult = await statusFuture;
    final blockedResult = await blockedFuture;

    // Handle profile result
    UserProfileModel? profile;
    profileResult.fold(
      (failure) {
        _logger.error('getUserProfile failed for $_userId', failure);
        emit(state.copyWith(
          status: UserProfileStatus.failure,
          error: failure.message,
        ));
        return;
      },
      (p) => profile = p,
    );

    if (profile == null) return;

    // Merge online status into profile
    statusResult.fold(
      (l) => _logger.warning('getUserStatus failed for $_userId: ${l.message}'),
      (statusModel) {
        profile = profile!.copyWith(
          isOnline: statusModel.isOnline,
          lastSeenAt: statusModel.lastSeenAt,
        );
      },
    );

    // Check if blocked
    bool isBlocked = false;
    blockedResult.fold(
      (l) => _logger.warning('getBlockedIds failed: ${l.message}'),
      (users) => isBlocked = users.any((b) => b.userId == _userId),
    );

    _logger.info('Profile loaded: $_userId, online=${profile!.isOnline}, blocked=$isBlocked');
    emit(state.copyWith(
      status: UserProfileStatus.ready,
      profile: profile,
      isBlocked: isBlocked,
    ));
  }

  Future<void> _onBlockUser(
    BlockUserEvent event,
    Emitter<UserProfileState> emit,
  ) async {
    _logger.info('Blocking user: $_userId');
    final result = await _repo.blockUser(_userId);
    result.fold(
      (failure) => _logger.error('blockUser failed for $_userId', failure),
      (_) {
        _logger.info('User $_userId blocked');
        emit(state.copyWith(isBlocked: true));
      },
    );
  }

  Future<void> _onUnblockUser(
    UnblockUserEvent event,
    Emitter<UserProfileState> emit,
  ) async {
    _logger.info('Unblocking user: $_userId');
    final result = await _repo.unblockUser(_userId);
    result.fold(
      (failure) => _logger.error('unblockUser failed for $_userId', failure),
      (_) {
        _logger.info('User $_userId unblocked');
        emit(state.copyWith(isBlocked: false));
      },
    );
  }
}
