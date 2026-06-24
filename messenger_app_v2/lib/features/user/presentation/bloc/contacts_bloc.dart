import 'dart:async';

import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:injectable/injectable.dart';

import '../../../../shared/data/logger/app_logger.dart';
import '../../data/models/user_models.dart';
import '../../domain/repositories/user_repository.dart';

// --- EVENTS ---
abstract class ContactsEvent extends Equatable {
  @override
  List<Object> get props => [];
}

class LoadContacts extends ContactsEvent {}

class RemoveContact extends ContactsEvent {
  final String userId;
  RemoveContact(this.userId);
  @override
  List<Object> get props => [userId];
}

// Blocked users events (managed by same bloc for simplicity)
class LoadBlockedUsers extends ContactsEvent {}

class UnblockFromList extends ContactsEvent {
  final String userId;
  UnblockFromList(this.userId);
  @override
  List<Object> get props => [userId];
}

// --- STATE ---
enum ContactsStatus { initial, loading, ready, failure }

class ContactsState extends Equatable {
  final ContactsStatus status;
  final List<ContactModel> contacts;
  final List<UserProfileModel> blockedUsers;
  final bool isLoadingBlocked;
  final String? error;

  const ContactsState({
    this.status = ContactsStatus.initial,
    this.contacts = const [],
    this.blockedUsers = const [],
    this.isLoadingBlocked = false,
    this.error,
  });

  ContactsState copyWith({
    ContactsStatus? status,
    List<ContactModel>? contacts,
    List<UserProfileModel>? blockedUsers,
    bool? isLoadingBlocked,
    String? error,
  }) {
    return ContactsState(
      status: status ?? this.status,
      contacts: contacts ?? this.contacts,
      blockedUsers: blockedUsers ?? this.blockedUsers,
      isLoadingBlocked: isLoadingBlocked ?? this.isLoadingBlocked,
      error: error ?? this.error,
    );
  }

  @override
  List<Object?> get props => [status, contacts, blockedUsers, isLoadingBlocked, error];
}

// --- BLOC ---
@injectable
class ContactsBloc extends Bloc<ContactsEvent, ContactsState> {
  final UserRepository _repo;
  final AppLogger _logger;

  ContactsBloc(this._repo, this._logger) : super(const ContactsState()) {
    on<LoadContacts>(_onLoadContacts);
    on<RemoveContact>(_onRemoveContact);
    on<LoadBlockedUsers>(_onLoadBlockedUsers);
    on<UnblockFromList>(_onUnblockFromList);
  }

  Future<void> _onLoadContacts(
    LoadContacts event,
    Emitter<ContactsState> emit,
  ) async {
    emit(state.copyWith(status: ContactsStatus.loading));
    _logger.info('Loading contacts');

    final result = await _repo.getContacts();
    await result.fold(
      (failure) async {
        _logger.error('getContacts failed', failure);
        emit(state.copyWith(status: ContactsStatus.failure, error: failure.message));
      },
      (contacts) async {
        _logger.info('Loaded ${contacts.length} contacts — fetching profiles');

        // Fetch profiles concurrently for display
        final enriched = await Future.wait(
          contacts.map((c) async {
            final profileResult = await _repo.getUserProfile(c.contactId);
            return profileResult.fold(
              (_) => c,
              (profile) => c.withProfile(profile),
            );
          }),
        );

        emit(state.copyWith(
          status: ContactsStatus.ready,
          contacts: enriched,
        ));
      },
    );
  }

  Future<void> _onRemoveContact(
    RemoveContact event,
    Emitter<ContactsState> emit,
  ) async {
    _logger.info('Removing contact: ${event.userId}');
    final result = await _repo.removeContact(event.userId);
    result.fold(
      (failure) => _logger.error('removeContact failed (${event.userId})', failure),
      (_) {
        _logger.info('Contact ${event.userId} removed');
        final updated = state.contacts
            .where((c) => c.contactId != event.userId)
            .toList();
        emit(state.copyWith(contacts: updated));
      },
    );
  }

  Future<void> _onLoadBlockedUsers(
    LoadBlockedUsers event,
    Emitter<ContactsState> emit,
  ) async {
    emit(state.copyWith(isLoadingBlocked: true));
    _logger.info('Loading blocked users');

    final result = await _repo.getBlockedUsers();
    result.fold(
      (failure) {
        _logger.error('getBlockedUsers failed', failure);
        emit(state.copyWith(isLoadingBlocked: false));
      },
      (users) {
        // Map BlockedUserModel → UserProfileModel to keep ContactsState unchanged.
        final profiles = users.map((b) => UserProfileModel(
          id: b.userId,
          username: b.username,
          displayName: b.displayName,
          avatarMediaId: b.avatarMediaId,
        )).toList();
        _logger.info('Loaded ${profiles.length} blocked users');
        emit(state.copyWith(isLoadingBlocked: false, blockedUsers: profiles));
      },
    );
  }

  Future<void> _onUnblockFromList(
    UnblockFromList event,
    Emitter<ContactsState> emit,
  ) async {
    _logger.info('Unblocking user from list: ${event.userId}');
    final result = await _repo.unblockUser(event.userId);
    result.fold(
      (failure) => _logger.error('unblockUser failed (${event.userId})', failure),
      (_) {
        _logger.info('User ${event.userId} unblocked');
        final updated = state.blockedUsers
            .where((u) => u.id != event.userId)
            .toList();
        emit(state.copyWith(blockedUsers: updated));
      },
    );
  }
}
