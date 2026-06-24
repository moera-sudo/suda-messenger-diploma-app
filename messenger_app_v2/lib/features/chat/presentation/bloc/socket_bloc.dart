import 'dart:async';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';
import 'package:injectable/injectable.dart';
import '../../../../shared/data/api/socket_client.dart';
import '../../../../shared/domain/models/socket_event.dart';

// Events
abstract class SocketEventBase extends Equatable {
  @override
  List<Object> get props => [];
}
class SocketConnect extends SocketEventBase {}
class SocketDisconnect extends SocketEventBase {}
class _SocketOnEvent extends SocketEventBase {
  final SocketEvent event;
  _SocketOnEvent(this.event);
}

// State
abstract class SocketState extends Equatable {
  @override
  List<Object> get props => [];
}
class SocketInitial extends SocketState {}
class SocketConnected extends SocketState {}
class SocketDisconnected extends SocketState {}

@injectable
class SocketBloc extends Bloc<SocketEventBase, SocketState> {
  final SocketClient _socketClient;
  StreamSubscription? _socketSub;

  SocketBloc(this._socketClient) : super(SocketInitial()) {
    on<SocketConnect>(_onConnect);
    on<SocketDisconnect>(_onDisconnect);
    on<_SocketOnEvent>(_onEvent);
  }

  Future<void> _onConnect(SocketConnect event, Emitter<SocketState> emit) async {
    await _socketClient.connect();
    _socketSub?.cancel();
    _socketSub = _socketClient.events.listen((event) {
      add(_SocketOnEvent(event));
    });
    emit(SocketConnected());
  }

  Future<void> _onDisconnect(SocketDisconnect event, Emitter<SocketState> emit) async {
    _socketClient.disconnect();
    _socketSub?.cancel();
    emit(SocketDisconnected());
  }

  void _onEvent(_SocketOnEvent event, Emitter<SocketState> emit) {
    // Этот BLoC может быть "хабом", который пересылает события другим блокам
    // Или другие блоки могут сами подписаться на SocketClient.
    // Пока оставим пустым, логику обработки сообщений вынесем в ChatListBloc и ChatDetailBloc.
  }

  @override
  Future<void> close() {
    _socketSub?.cancel();
    return super.close();
  }
}