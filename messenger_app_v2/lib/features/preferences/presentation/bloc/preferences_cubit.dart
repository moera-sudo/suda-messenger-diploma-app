import 'dart:async';

import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:injectable/injectable.dart';

import '../../../../shared/data/api/socket_client.dart';
import '../../../../shared/data/logger/app_logger.dart';
import '../../../../shared/domain/models/socket_event.dart';
import '../../data/models/user_preferences.dart';
import '../../domain/repositories/preferences_repository.dart';

enum PreferencesStatus { initial, loading, loaded, failure }

class PreferencesState extends Equatable {
  final PreferencesStatus status;
  final UserPreferences prefs;

  const PreferencesState({
    this.status = PreferencesStatus.initial,
    this.prefs = const UserPreferences(),
  });

  PreferencesState copyWith({PreferencesStatus? status, UserPreferences? prefs}) {
    return PreferencesState(
      status: status ?? this.status,
      prefs: prefs ?? this.prefs,
    );
  }

  @override
  List<Object?> get props => [status, prefs];
}

@injectable
class PreferencesCubit extends Cubit<PreferencesState> {
  final PreferencesRepository _repo;
  final SocketClient _socket;
  final AppLogger _logger;

  StreamSubscription? _socketSub;

  PreferencesCubit(this._repo, this._socket, this._logger)
      : super(const PreferencesState()) {
    // Sync across sessions: server pushes PREFERENCES_UPDATED on any change.
    _socketSub = _socket.events.listen((event) {
      if (event.type == SocketEventType.preferencesUpdated) {
        _onRemoteUpdate(event);
      }
    });
  }

  Future<void> load() async {
    emit(state.copyWith(status: PreferencesStatus.loading));
    final result = await _repo.getPreferences();
    result.fold(
      (failure) {
        _logger.error('Preferences load failed: ${failure.message}');
        emit(state.copyWith(status: PreferencesStatus.failure));
      },
      (prefs) => emit(PreferencesState(status: PreferencesStatus.loaded, prefs: prefs)),
    );
  }

  /// Optimistically applies [next], then persists [patch]. Rolls back on failure
  /// (ApiClient already surfaces the error message to the user).
  Future<void> update(UserPreferences next, Map<String, dynamic> patch) async {
    final previous = state.prefs;
    emit(state.copyWith(status: PreferencesStatus.loaded, prefs: next));

    final result = await _repo.updatePreferences(patch);
    result.fold(
      (failure) {
        _logger.error('Preferences update rolled back: ${failure.message}');
        emit(state.copyWith(prefs: previous));
      },
      (_) => _logger.debug('Preferences persisted: ${patch.keys.join(', ')}'),
    );
  }

  void _onRemoteUpdate(SocketEvent event) {
    final raw = event.payload?['preferences'];
    if (raw is! Map) return;
    try {
      final prefs = UserPreferences.fromJson(Map<String, dynamic>.from(raw));
      _logger.debug('Preferences synced from another session');
      emit(PreferencesState(status: PreferencesStatus.loaded, prefs: prefs));
    } catch (e) {
      _logger.warning('Failed to parse PREFERENCES_UPDATED payload: $e');
    }
  }

  @override
  Future<void> close() {
    _socketSub?.cancel();
    return super.close();
  }
}
