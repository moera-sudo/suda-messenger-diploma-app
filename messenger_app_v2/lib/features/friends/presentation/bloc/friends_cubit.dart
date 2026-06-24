import 'dart:async';

import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:injectable/injectable.dart';

import '../../../../shared/data/api/socket_client.dart';
import '../../../../shared/data/logger/app_logger.dart';
import '../../../../shared/domain/models/socket_event.dart';
import '../../../../shared/presentation/feedback/app_feedback.dart';
import '../../data/models/friend_models.dart';
import '../../domain/repositories/friends_repository.dart';

enum FriendsStatus { initial, loading, loaded, failure }

class FriendsState extends Equatable {
  final FriendsStatus status;
  final List<FriendModel> friends;
  final List<FriendRequestModel> incoming;
  final List<FriendRequestModel> outgoing;

  const FriendsState({
    this.status = FriendsStatus.initial,
    this.friends = const [],
    this.incoming = const [],
    this.outgoing = const [],
  });

  FriendsState copyWith({
    FriendsStatus? status,
    List<FriendModel>? friends,
    List<FriendRequestModel>? incoming,
    List<FriendRequestModel>? outgoing,
  }) {
    return FriendsState(
      status: status ?? this.status,
      friends: friends ?? this.friends,
      incoming: incoming ?? this.incoming,
      outgoing: outgoing ?? this.outgoing,
    );
  }

  @override
  List<Object?> get props => [status, friends, incoming, outgoing];
}

@injectable
class FriendsCubit extends Cubit<FriendsState> {
  final FriendsRepository _repo;
  final SocketClient _socket;
  final AppLogger _logger;

  StreamSubscription? _socketSub;

  FriendsCubit(this._repo, this._socket, this._logger)
      : super(const FriendsState()) {
    _socketSub = _socket.events.listen((event) {
      if ([
        SocketEventType.friendRequestReceived,
        SocketEventType.friendRequestAccepted,
        SocketEventType.friendRequestRejected,
        SocketEventType.friendRequestCancelled,
        SocketEventType.friendRemoved,
      ].contains(event.type)) {
        _logger.info('Friends WS event: ${event.type}');
        _reload();
      }
    });
  }

  Future<void> load() async {
    emit(state.copyWith(status: FriendsStatus.loading));
    await Future.wait([_loadFriends(), _loadRequests()]);
  }

  Future<void> _reload() async {
    await Future.wait([_loadFriends(), _loadRequests()]);
  }

  Future<void> _loadFriends() async {
    final result = await _repo.getFriends();
    result.fold(
      (f) {
        _logger.error('getFriends failed: ${f.message}');
        emit(state.copyWith(status: FriendsStatus.failure));
      },
      (list) => emit(state.copyWith(status: FriendsStatus.loaded, friends: list)),
    );
  }

  Future<void> _loadRequests() async {
    final result = await _repo.getRequests();
    result.fold(
      (f) => _logger.error('getRequests failed: ${f.message}'),
      (list) {
        final inc = list.where((r) => r.direction == 'incoming').toList();
        final out = list.where((r) => r.direction == 'outgoing').toList();
        emit(state.copyWith(incoming: inc, outgoing: out));
      },
    );
  }

  Future<void> sendRequest(String userId) async {
    final result = await _repo.sendRequest(userId);
    result.fold((_) {}, (_) => _reload());
  }

  Future<void> cancelRequest(String requestId) async {
    final result = await _repo.cancelRequest(requestId);
    result.fold((_) {}, (_) => _reload());
  }

  Future<void> acceptRequest(String requestId) async {
    final result = await _repo.acceptRequest(requestId);
    result.fold(
      (f) {
        _logger.error('acceptRequest failed ($requestId): ${f.message}');
        AppFeedback.showError(f.message);
      },
      (_) => _reload(),
    );
  }

  Future<void> rejectRequest(String requestId) async {
    final result = await _repo.rejectRequest(requestId);
    result.fold(
      (f) {
        _logger.error('rejectRequest failed ($requestId): ${f.message}');
        AppFeedback.showError(f.message);
      },
      (_) => _reload(),
    );
  }

  Future<void> unfriend(String userId) async {
    final result = await _repo.unfriend(userId);
    result.fold((_) {}, (_) => _reload());
  }

  @override
  Future<void> close() {
    _socketSub?.cancel();
    return super.close();
  }
}
