import 'dart:async';

import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:injectable/injectable.dart';

import '../../../../features/wallet/data/models/wallet_models.dart';
import '../../../../features/wallet/domain/repositories/wallet_repository.dart';
import '../../../../shared/domain/models/app_failure.dart';
import '../../data/models/my_profile.dart';
import '../../domain/repositories/profile_repository.dart';

abstract class ProfileEvent extends Equatable {
  const ProfileEvent();

  @override
  List<Object?> get props => [];
}

class ProfileLoadRequested extends ProfileEvent {
  const ProfileLoadRequested();
}

class ProfileSaveRequested extends ProfileEvent {
  final String displayName;
  final String firstName;
  final String lastName;
  final String bio;

  const ProfileSaveRequested({
    required this.displayName,
    required this.firstName,
    required this.lastName,
    required this.bio,
  });

  @override
  List<Object?> get props => [displayName, firstName, lastName, bio];
}

class ProfileAvatarUploadRequested extends ProfileEvent {
  final String filePath;
  final String? fileName;

  const ProfileAvatarUploadRequested({required this.filePath, this.fileName});

  @override
  List<Object?> get props => [filePath, fileName];
}

class ProfileLogoutRequested extends ProfileEvent {
  const ProfileLogoutRequested();
}

/// Loads (or reloads) the SUDA wallet balance shown on the profile.
class ProfileWalletRequested extends ProfileEvent {
  const ProfileWalletRequested();
}

enum ProfileViewStatus { initial, loading, ready, failure }

enum ProfileActionStatus {
  idle,
  saveSuccess,
  uploadSuccess,
  logoutSuccess,
  failure,
}

class ProfileState extends Equatable {
  final ProfileViewStatus status;
  final ProfileActionStatus actionStatus;
  final MyProfile? profile;
  final AppFailure? failure;
  final String? actionMessage;
  final bool isSaving;
  final bool isUploadingAvatar;
  final bool isLoggingOut;
  final int avatarVersion;
  // SUDA wallet (address + balance); null until loaded / on error.
  final MyWallet? wallet;

  const ProfileState({
    this.status = ProfileViewStatus.initial,
    this.actionStatus = ProfileActionStatus.idle,
    this.profile,
    this.failure,
    this.actionMessage,
    this.isSaving = false,
    this.isUploadingAvatar = false,
    this.isLoggingOut = false,
    this.avatarVersion = 0,
    this.wallet,
  });

  ProfileState copyWith({
    ProfileViewStatus? status,
    ProfileActionStatus? actionStatus,
    MyProfile? profile,
    AppFailure? failure,
    String? actionMessage,
    bool clearActionMessage = false,
    bool? isSaving,
    bool? isUploadingAvatar,
    bool? isLoggingOut,
    int? avatarVersion,
    MyWallet? wallet,
  }) {
    return ProfileState(
      status: status ?? this.status,
      actionStatus: actionStatus ?? this.actionStatus,
      profile: profile ?? this.profile,
      failure: failure ?? this.failure,
      actionMessage: clearActionMessage
          ? null
          : (actionMessage ?? this.actionMessage),
      isSaving: isSaving ?? this.isSaving,
      isUploadingAvatar: isUploadingAvatar ?? this.isUploadingAvatar,
      isLoggingOut: isLoggingOut ?? this.isLoggingOut,
      avatarVersion: avatarVersion ?? this.avatarVersion,
      wallet: wallet ?? this.wallet,
    );
  }

  @override
  List<Object?> get props => [
    status,
    actionStatus,
    profile,
    failure,
    actionMessage,
    isSaving,
    isUploadingAvatar,
    isLoggingOut,
    avatarVersion,
    wallet,
  ];
}

@injectable
class ProfileBloc extends Bloc<ProfileEvent, ProfileState> {
  final ProfileRepository _repository;
  final WalletRepository _walletRepository;

  ProfileBloc(this._repository, this._walletRepository)
      : super(const ProfileState()) {
    on<ProfileLoadRequested>(_onLoad);
    on<ProfileSaveRequested>(_onSave);
    on<ProfileAvatarUploadRequested>(_onUploadAvatar);
    on<ProfileLogoutRequested>(_onLogout);
    on<ProfileWalletRequested>(_onLoadWallet);
  }

  /// Loads the SUDA wallet balance. Retries a few times on 404 because the
  /// wallet is created asynchronously after signup (Frontend.md §7.1).
  Future<void> _onLoadWallet(
    ProfileWalletRequested event,
    Emitter<ProfileState> emit,
  ) async {
    for (var attempt = 0; attempt < 3; attempt++) {
      final result = await _walletRepository.getMyWallet();
      final wallet = result.fold((_) => null, (w) => w);
      if (wallet != null) {
        emit(state.copyWith(wallet: wallet));
        return;
      }
      // Likely 404 while the wallet is still being generated — back off briefly.
      await Future.delayed(const Duration(seconds: 1));
    }
  }

  Future<void> _onLoad(
    ProfileLoadRequested event,
    Emitter<ProfileState> emit,
  ) async {
    emit(
      state.copyWith(
        status: ProfileViewStatus.loading,
        actionStatus: ProfileActionStatus.idle,
        clearActionMessage: true,
      ),
    );

    final result = await _repository.getMe();
    result.fold(
      (failure) {
        emit(
          state.copyWith(
            status: ProfileViewStatus.failure,
            actionStatus: ProfileActionStatus.failure,
            failure: failure,
            actionMessage: failure.message,
          ),
        );
      },
      (profile) {
        emit(
          state.copyWith(
            status: ProfileViewStatus.ready,
            actionStatus: ProfileActionStatus.idle,
            profile: profile,
            failure: null,
            clearActionMessage: true,
          ),
        );
        // Load the SUDA balance in the background (non-blocking).
        add(const ProfileWalletRequested());
      },
    );
  }

  Future<void> _onSave(
    ProfileSaveRequested event,
    Emitter<ProfileState> emit,
  ) async {
    if (state.isSaving || state.profile == null) {
      return;
    }

    emit(
      state.copyWith(
        isSaving: true,
        actionStatus: ProfileActionStatus.idle,
        clearActionMessage: true,
      ),
    );

    final payload = UpdateMyProfilePayload(
      displayName: event.displayName.trim(),
      firstName: event.firstName.trim(),
      lastName: event.lastName.trim(),
      bio: event.bio.trim(),
    );

    final result = await _repository.updateProfile(payload);
    if (result.isLeft()) {
      final failure = result.fold((f) => f, (_) => null)!;
      emit(
        state.copyWith(
          isSaving: false,
          actionStatus: ProfileActionStatus.failure,
          failure: failure,
          actionMessage: failure.message,
        ),
      );
      return;
    }

    final refresh = await _repository.getMe();
    refresh.fold(
      (failure) {
        emit(
          state.copyWith(
            isSaving: false,
            actionStatus: ProfileActionStatus.failure,
            failure: failure,
            actionMessage: failure.message,
          ),
        );
      },
      (profile) {
        emit(
          state.copyWith(
            status: ProfileViewStatus.ready,
            isSaving: false,
            actionStatus: ProfileActionStatus.saveSuccess,
            profile: profile,
            failure: null,
          ),
        );
      },
    );
  }

  Future<void> _onUploadAvatar(
    ProfileAvatarUploadRequested event,
    Emitter<ProfileState> emit,
  ) async {
    if (state.isUploadingAvatar || state.profile == null) {
      return;
    }

    emit(
      state.copyWith(
        isUploadingAvatar: true,
        actionStatus: ProfileActionStatus.idle,
        clearActionMessage: true,
      ),
    );

    final result = await _repository.uploadAvatar(
      filePath: event.filePath,
      fileName: event.fileName,
    );

    result.fold(
      (failure) {
        emit(
          state.copyWith(
            isUploadingAvatar: false,
            actionStatus: ProfileActionStatus.failure,
            failure: failure,
            actionMessage: failure.message,
          ),
        );
      },
      (avatarUrl) {
        final current = state.profile!;
        emit(
          state.copyWith(
            profile: current.copyWith(avatarUrl: avatarUrl),
            isUploadingAvatar: false,
            actionStatus: ProfileActionStatus.uploadSuccess,
            avatarVersion: state.avatarVersion + 1,
            failure: null,
          ),
        );
      },
    );
  }

  Future<void> _onLogout(
    ProfileLogoutRequested event,
    Emitter<ProfileState> emit,
  ) async {
    if (state.isLoggingOut) {
      return;
    }

    emit(
      state.copyWith(
        isLoggingOut: true,
        actionStatus: ProfileActionStatus.idle,
        clearActionMessage: true,
      ),
    );

    final result = await _repository.logout();
    result.fold(
      (failure) {
        emit(
          state.copyWith(
            isLoggingOut: false,
            actionStatus: ProfileActionStatus.failure,
            actionMessage: failure.message,
            failure: failure,
          ),
        );
      },
      (_) {
        emit(
          state.copyWith(
            isLoggingOut: false,
            actionStatus: ProfileActionStatus.logoutSuccess,
            failure: null,
          ),
        );
      },
    );
  }
}
