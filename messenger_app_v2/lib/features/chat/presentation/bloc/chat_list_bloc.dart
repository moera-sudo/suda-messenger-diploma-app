import 'dart:async';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';
import 'package:injectable/injectable.dart';
import 'package:jwt_decoder/jwt_decoder.dart';
import '../../../../shared/data/api/socket_client.dart';
import '../../../../shared/domain/models/socket_event.dart';
import '../../data/models/chat_models.dart';
import '../../domain/repositories/chat_repository.dart';
import '../../../../shared/data/storage/secure_storage_client.dart';
import '../../../../shared/data/logger/app_logger.dart';

// Events
abstract class ChatListEvent extends Equatable {
  @override
  List<Object> get props => [];
}
class LoadChats extends ChatListEvent {}
class _OnSocketMessage extends ChatListEvent {
  final SocketEvent event;
  _OnSocketMessage(this.event);
}
class PinChatEvent extends ChatListEvent {
  final String chatId;
  PinChatEvent(this.chatId);
  @override
  List<Object> get props => [chatId];
}
class UnpinChatEvent extends ChatListEvent {
  final String chatId;
  UnpinChatEvent(this.chatId);
  @override
  List<Object> get props => [chatId];
}

/// Reset the unread badge for a chat (dispatched when the user opens it).
class MarkChatRead extends ChatListEvent {
  final String chatId;
  MarkChatRead(this.chatId);
  @override
  List<Object> get props => [chatId];
}

class MuteChatEvent extends ChatListEvent {
  final String chatId;
  MuteChatEvent(this.chatId);
  @override
  List<Object> get props => [chatId];
}

class UnmuteChatEvent extends ChatListEvent {
  final String chatId;
  UnmuteChatEvent(this.chatId);
  @override
  List<Object> get props => [chatId];
}

class DeleteChatEvent extends ChatListEvent {
  final String chatId;
  final String scope;
  DeleteChatEvent(this.chatId, {this.scope = 'everyone'});
  @override
  List<Object> get props => [chatId, scope];
}

// State
enum ChatListStatus { initial, loading, success, failure }
class ChatListState extends Equatable {
  final ChatListStatus status;
  final List<ChatListItem> chats;
  final String? error;

  const ChatListState({
    this.status = ChatListStatus.initial,
    this.chats = const [],
    this.error,
  });

  @override
  List<Object?> get props => [status, chats, error];
}

@injectable
class ChatListBloc extends Bloc<ChatListEvent, ChatListState> {
  final ChatRepository _repo;
  final SocketClient _socket;
  final SecureStorageClient _storage;
  final AppLogger _logger;
  StreamSubscription? _socketSub;
  String? _myUserId;
  final Map<String, String> _chatPinIds = {}; // chatId → pinId from server

  ChatListBloc(this._repo, this._socket, this._storage, this._logger) : super(const ChatListState()) {
    on<LoadChats>(_onLoadChats);
    on<_OnSocketMessage>(_onSocketMessage);
    on<PinChatEvent>(_onPinChat);
    on<UnpinChatEvent>(_onUnpinChat);
    on<MarkChatRead>(_onMarkChatRead);
    on<MuteChatEvent>(_onMuteChat);
    on<UnmuteChatEvent>(_onUnmuteChat);
    on<DeleteChatEvent>(_onDeleteChat);

    // Подписываемся на сокет сразу
    _init();
  }

  Future<void> _init() async {
    final token = await _storage.getAccessToken();
    if (token != null) {
      final decoded = JwtDecoder.decode(token);
      _myUserId = decoded['user_id'] ?? decoded['sub'];
    }

    _socketSub = _socket.events.listen((event) {
      if (event.type == SocketEventType.newMessage ||
          event.type == SocketEventType.chatUpdated ||
          event.type == SocketEventType.memberLeft ||
          event.type == SocketEventType.memberRemoved ||
          event.type == SocketEventType.userBlocked ||
          event.type == SocketEventType.userUnblocked ||
          event.type == SocketEventType.chatDeleted) {
        add(_OnSocketMessage(event));
      }
    });
  }

  Future<void> _onLoadChats(LoadChats event, Emitter<ChatListState> emit) async {
    emit(const ChatListState(status: ChatListStatus.loading));
    final result = await _repo.getChats();
    result.fold(
      (l) {
        _logger.error('getChats failed', l);
        emit(ChatListState(status: ChatListStatus.failure, error: l.message));
      },
      (r) async {
        _logger.info('Loaded ${r.length} chats');
        emit(ChatListState(status: ChatListStatus.success, chats: r));
        // Populate pin-ID map so unpin knows the server-side pin UUID
        final pinsResult = await _repo.getPinIds();
        pinsResult.fold(
          (l) => _logger.warning('getPinIds failed: ${l.message}'),
          (ids) {
            _chatPinIds.clear();
            _chatPinIds.addAll(ids);
            _logger.debug('Loaded ${ids.length} pin mappings');
          },
        );
      },
    );
  }

  void _onSocketMessage(_OnSocketMessage event, Emitter<ChatListState> emit) {
    final payload = event.event.payload;
    if (payload == null) return;

    try {
      final eventType = event.event.type;

      // CHAT_UPDATED: update name/avatar in list
      if (eventType == SocketEventType.chatUpdated) {
        final chatData = payload['chat'] as Map<String, dynamic>?;
        if (chatData == null) return;
        final chatId = chatData['id'] as String?;
        if (chatId == null) return;

        final currentChats = List<ChatListItem>.from(state.chats);
        final idx = currentChats.indexWhere((c) => c.id == chatId);
        if (idx != -1) {
          currentChats[idx] = currentChats[idx].copyWith(
            name: chatData['name'] as String?,
            avatarMediaId: chatData['avatar_media_id'] as String?,
          );
          emit(ChatListState(status: ChatListStatus.success, chats: currentChats));
        }
        return;
      }

      // CHAT_DELETED: remove chat from list for all affected users.
      if (eventType == SocketEventType.chatDeleted) {
        final chatId = payload['chat_id'] as String?;
        if (chatId != null) {
          final currentChats = List<ChatListItem>.from(state.chats)
            ..removeWhere((c) => c.id == chatId);
          _logger.info('Chat $chatId removed from list via CHAT_DELETED (scope=${payload['scope']})');
          emit(ChatListState(status: ChatListStatus.success, chats: currentChats));
        }
        return;
      }

      // MEMBER_LEFT / MEMBER_REMOVED: if it's me → remove chat from list
      if (eventType == SocketEventType.memberLeft ||
          eventType == SocketEventType.memberRemoved) {
        final userId = payload['user_id'] as String?;
        final chatId = payload['chat_id'] as String?;
        if (userId == _myUserId && chatId != null) {
          final currentChats = List<ChatListItem>.from(state.chats)
            ..removeWhere((c) => c.id == chatId);
          emit(ChatListState(status: ChatListStatus.success, chats: currentChats));
        }
        return;
      }

      // USER_BLOCKED / USER_UNBLOCKED: update block status on the DIRECT chat (§4.4).
      if (eventType == SocketEventType.userBlocked ||
          eventType == SocketEventType.userUnblocked) {
        final byUser = payload['by_user'] as String?;
        if (byUser == null) return;
        final isBlocked = eventType == SocketEventType.userBlocked;
        final currentChats = List<ChatListItem>.from(state.chats);
        // Find the DIRECT chat where the other participant is byUser.
        final idx = currentChats.indexWhere(
          (c) => c.type == 'DIRECT' && c.interlocutorId == byUser,
        );
        if (idx != -1) {
          currentChats[idx] = currentChats[idx].copyWith(blockedMe: isBlocked);
          _logger.info('Chat list: ${isBlocked ? "BLOCKED" : "UNBLOCKED"} by $byUser');
          emit(ChatListState(status: ChatListStatus.success, chats: currentChats));
        }
        return;
      }

      // NEW_MESSAGE: update last message preview + unread count
      final incomingMsg = Message.fromJson(payload);
      final targetChatId = incomingMsg.chatId;
      final currentChats = List<ChatListItem>.from(state.chats);
      final chatIndex = currentChats.indexWhere((c) => c.id == targetChatId);

      if (chatIndex != -1) {
        final existingChat = currentChats[chatIndex];
        final isFromMe = incomingMsg.senderId == _myUserId;
        final newUnread = isFromMe ? existingChat.unreadCount : existingChat.unreadCount + 1;

        final updatedChat = existingChat.copyWith(
          lastMessageContent: incomingMsg.content,
          lastMessageTime: incomingMsg.createdAt,
          unreadCount: newUnread,
          lastMessageType: incomingMsg.type,
        );

        if (existingChat.isPinned) {
          // Pinned chats keep their server-assigned position; only content updates.
          currentChats[chatIndex] = updatedChat;
        } else {
          // Non-pinned chats bubble to the top of the unpinned section.
          currentChats.removeAt(chatIndex);
          final firstNonPinned = currentChats.indexWhere((c) => !c.isPinned);
          currentChats.insert(firstNonPinned == -1 ? 0 : firstNonPinned, updatedChat);
        }
        emit(ChatListState(status: ChatListStatus.success, chats: currentChats));
      } else {
        _logger.debug('Chat $targetChatId not in list — reloading all chats');
        if (state.status != ChatListStatus.loading) add(LoadChats());
      }
    } catch (e, stackTrace) {
      _logger.error('Failed to handle socket event in ChatListBloc', e, stackTrace);
    }
  }

  void _onMarkChatRead(MarkChatRead event, Emitter<ChatListState> emit) {
    final idx = state.chats.indexWhere((c) => c.id == event.chatId);
    if (idx == -1 || state.chats[idx].unreadCount == 0) return;
    final updated = List<ChatListItem>.from(state.chats);
    updated[idx] = updated[idx].copyWith(unreadCount: 0);
    emit(ChatListState(status: ChatListStatus.success, chats: updated));
  }

  Future<void> _onPinChat(PinChatEvent event, Emitter<ChatListState> emit) async {
    // Optimistic update
    final currentChats = List<ChatListItem>.from(state.chats);
    final idx = currentChats.indexWhere((c) => c.id == event.chatId);
    if (idx == -1) return;

    currentChats[idx] = currentChats[idx].copyWith(isPinned: true);
    // Move pinned chat to the top of the list
    final pinned = currentChats.removeAt(idx);
    currentChats.insert(0, pinned);
    emit(ChatListState(status: ChatListStatus.success, chats: currentChats));

    final result = await _repo.pinChat(event.chatId);
    result.fold(
      (failure) {
        _logger.error('pinChat failed (chatId=${event.chatId})', failure);
        // Revert on failure
        if (state.status != ChatListStatus.loading) add(LoadChats());
      },
      (pinId) {
        _chatPinIds[event.chatId] = pinId;
        _logger.info('Chat ${event.chatId} pinned (pinId=$pinId)');
      },
    );
  }

  Future<void> _onUnpinChat(UnpinChatEvent event, Emitter<ChatListState> emit) async {
    final pinId = _chatPinIds[event.chatId];
    if (pinId == null) {
      _logger.warning('unpinChat: no pinId for chatId=${event.chatId} — reloading');
      if (state.status != ChatListStatus.loading) add(LoadChats());
      return;
    }

    // Optimistic update
    final currentChats = List<ChatListItem>.from(state.chats);
    final idx = currentChats.indexWhere((c) => c.id == event.chatId);
    if (idx != -1) {
      currentChats[idx] = currentChats[idx].copyWith(isPinned: false);
      emit(ChatListState(status: ChatListStatus.success, chats: currentChats));
    }

    final result = await _repo.unpinChat(pinId);
    result.fold(
      (failure) {
        _logger.error('unpinChat failed (pinId=$pinId)', failure);
        if (state.status != ChatListStatus.loading) add(LoadChats());
      },
      (_) {
        _chatPinIds.remove(event.chatId);
        _logger.info('Chat ${event.chatId} unpinned');
      },
    );
  }

  Future<void> _onMuteChat(MuteChatEvent event, Emitter<ChatListState> emit) async {
    final idx = state.chats.indexWhere((c) => c.id == event.chatId);
    if (idx == -1) return;

    // Optimistic update
    final updated = List<ChatListItem>.from(state.chats);
    updated[idx] = updated[idx].copyWith(isMuted: true);
    emit(ChatListState(status: ChatListStatus.success, chats: updated));

    final result = await _repo.muteChat(event.chatId);
    result.fold(
      (f) {
        _logger.error('muteChat failed (chatId=${event.chatId})', f);
        if (state.status != ChatListStatus.loading) add(LoadChats());
      },
      (_) => _logger.info('Chat ${event.chatId} muted'),
    );
  }

  Future<void> _onUnmuteChat(UnmuteChatEvent event, Emitter<ChatListState> emit) async {
    final idx = state.chats.indexWhere((c) => c.id == event.chatId);
    if (idx == -1) return;

    // Optimistic update
    final updated = List<ChatListItem>.from(state.chats);
    updated[idx] = updated[idx].copyWith(isMuted: false);
    emit(ChatListState(status: ChatListStatus.success, chats: updated));

    final result = await _repo.unmuteChat(event.chatId);
    result.fold(
      (f) {
        _logger.error('unmuteChat failed (chatId=${event.chatId})', f);
        if (state.status != ChatListStatus.loading) add(LoadChats());
      },
      (_) => _logger.info('Chat ${event.chatId} unmuted'),
    );
  }

  Future<void> _onDeleteChat(DeleteChatEvent event, Emitter<ChatListState> emit) async {
    final result = await _repo.deleteChat(event.chatId, scope: event.scope);
    result.fold(
      (f) {
        // scope=me: user may have left the chat already — the backend rejects with
        // 403 but we still want to hide the chat locally (it's the user's intent).
        if (event.scope == 'me') {
          _logger.info('deleteChat scope=me failed (likely already left): ${f.message} — removing locally');
          final updated = List<ChatListItem>.from(state.chats)
            ..removeWhere((c) => c.id == event.chatId);
          emit(ChatListState(status: ChatListStatus.success, chats: updated));
        } else {
          _logger.error('deleteChat failed (chatId=${event.chatId}, scope=${event.scope})', f);
        }
      },
      (_) {
        _logger.info('Chat ${event.chatId} deleted (scope=${event.scope})');
        final updated = List<ChatListItem>.from(state.chats)
          ..removeWhere((c) => c.id == event.chatId);
        emit(ChatListState(status: ChatListStatus.success, chats: updated));
      },
    );
  }

  @override
  Future<void> close() {
    _socketSub?.cancel();
    _logger.info("ChatListBloc closed and socket sub canceled.");

    return super.close();
  }
}