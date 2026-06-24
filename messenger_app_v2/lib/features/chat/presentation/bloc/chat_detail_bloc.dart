import 'dart:async';

import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:injectable/injectable.dart';
import 'package:jwt_decoder/jwt_decoder.dart';

import '../../../../features/media/domain/repositories/media_repository.dart';
import '../../../../features/preferences/domain/repositories/preferences_repository.dart';
import '../../../../shared/data/api/socket_client.dart';
import '../../../../shared/data/logger/app_logger.dart';
import '../../../../shared/data/storage/secure_storage_client.dart';
import '../../../../shared/domain/models/app_failure.dart';
import '../../../../shared/domain/models/socket_event.dart';
import '../../data/models/chat_models.dart';
import '../../domain/repositories/chat_repository.dart';

// --- EVENTS ---
abstract class ChatDetailEvent extends Equatable {
  @override
  List<Object> get props => [];
}

class LoadMessages extends ChatDetailEvent {}

/// Loads group members so messages can show sender name + avatar.
class LoadGroupMembers extends ChatDetailEvent {}

class SendMessage extends ChatDetailEvent {
  final String content;
  final String type;
  SendMessage(this.content, {this.type = 'text'});
}

class EditMessageEvent extends ChatDetailEvent {
  final int messageId;
  final String newContent;
  EditMessageEvent(this.messageId, this.newContent);
  @override
  List<Object> get props => [messageId, newContent];
}

class DeleteMessageEvent extends ChatDetailEvent {
  final int messageId;
  final bool forEveryone;
  DeleteMessageEvent(this.messageId, {this.forEveryone = false});
  @override
  List<Object> get props => [messageId, forEveryone];
}

class _OnSocketEvent extends ChatDetailEvent {
  final SocketEvent event;
  _OnSocketEvent(this.event);
}

class OnTypingInput extends ChatDetailEvent {
  final bool isTyping;
  OnTypingInput(this.isTyping);
}

class SetReplyTarget extends ChatDetailEvent {
  final Message message;
  SetReplyTarget(this.message);
  @override
  List<Object> get props => [message.id];
}

class ClearReplyTarget extends ChatDetailEvent {}

class LoadMoreMessages extends ChatDetailEvent {
  final int beforeMessageId;
  LoadMoreMessages(this.beforeMessageId);
  @override
  List<Object> get props => [beforeMessageId];
}

class SendMediaMessage extends ChatDetailEvent {
  final String filePath;
  final String messageType; // IMAGE, FILE, VOICE, VIDEO
  final String? content;   // optional: for VOICE — duration in seconds as string
  SendMediaMessage({
    required this.filePath,
    required this.messageType,
    this.content,
  });
  @override
  List<Object> get props => [filePath, messageType];
}

/// Search jump-to: ensure [messageId] is loaded (paging older messages if
/// needed) and mark it as the highlight target so the page can scroll to it.
class JumpToMessage extends ChatDetailEvent {
  final int messageId;
  JumpToMessage(this.messageId);
  @override
  List<Object> get props => [messageId];
}

/// Clears the active jump-to highlight after the brief flash.
class ClearHighlight extends ChatDetailEvent {}

// --- STATE ---
enum ChatStatus { initial, loading, success, failure }

class ChatDetailState extends Equatable {
  final ChatStatus status;
  final List<Message> messages;
  final String? currentUserId;
  final bool isInterlocutorTyping;
  final bool isInterlocutorOnline;
  final String? interlocutorLastSeenAt;
  final Message? replyTarget;
  final bool hasMoreMessages;
  final bool isLoadingMore;
  // Block status for DIRECT chats (§4.1).
  // blockedMe  = interlocutor has blocked the current user (cannot write).
  // blockedByMe = current user has blocked the interlocutor.
  final bool blockedMe;
  final bool blockedByMe;
  // Convenience getter — input is disabled when either block flag is true.
  bool get canWrite => !blockedMe && !blockedByMe;
  // Token-gated public channel — user lacks sufficient SUDA balance.
  final bool isGatingRequired;
  // Group chats only: senderId → member (display name) and resolved avatar URL.
  final Map<String, ChatMember> members;
  final Map<String, String?> memberAvatarUrls;
  // Search jump-to: id of the message to briefly highlight, null when inactive.
  final int? highlightedMessageId;
  // DIRECT chats: interlocutor's avatar media id (resolved by SudaAvatar).
  final String? interlocutorAvatarMediaId;

  const ChatDetailState({
    this.status = ChatStatus.initial,
    this.messages = const [],
    this.currentUserId,
    this.isInterlocutorTyping = false,
    this.isInterlocutorOnline = false,
    this.interlocutorLastSeenAt,
    this.replyTarget,
    this.hasMoreMessages = false,
    this.isLoadingMore = false,
    this.blockedMe = false,
    this.blockedByMe = false,
    this.isGatingRequired = false,
    this.members = const {},
    this.memberAvatarUrls = const {},
    this.highlightedMessageId,
    this.interlocutorAvatarMediaId,
  });

  ChatDetailState copyWith({
    ChatStatus? status,
    List<Message>? messages,
    String? currentUserId,
    bool? isInterlocutorTyping,
    bool? isInterlocutorOnline,
    String? interlocutorLastSeenAt,
    bool clearLastSeen = false,
    Object? replyTarget = _sentinel,
    bool? hasMoreMessages,
    bool? isLoadingMore,
    bool? blockedMe,
    bool? blockedByMe,
    bool? isGatingRequired,
    Map<String, ChatMember>? members,
    Map<String, String?>? memberAvatarUrls,
    Object? highlightedMessageId = _sentinel,
    String? interlocutorAvatarMediaId,
  }) {
    return ChatDetailState(
      status: status ?? this.status,
      messages: messages ?? this.messages,
      currentUserId: currentUserId ?? this.currentUserId,
      isInterlocutorTyping: isInterlocutorTyping ?? this.isInterlocutorTyping,
      isInterlocutorOnline: isInterlocutorOnline ?? this.isInterlocutorOnline,
      interlocutorLastSeenAt: clearLastSeen ? null : (interlocutorLastSeenAt ?? this.interlocutorLastSeenAt),
      replyTarget: replyTarget == _sentinel ? this.replyTarget : (replyTarget as Message?),
      hasMoreMessages: hasMoreMessages ?? this.hasMoreMessages,
      isLoadingMore: isLoadingMore ?? this.isLoadingMore,
      blockedMe: blockedMe ?? this.blockedMe,
      blockedByMe: blockedByMe ?? this.blockedByMe,
      isGatingRequired: isGatingRequired ?? this.isGatingRequired,
      members: members ?? this.members,
      memberAvatarUrls: memberAvatarUrls ?? this.memberAvatarUrls,
      highlightedMessageId: highlightedMessageId == _sentinel
          ? this.highlightedMessageId
          : (highlightedMessageId as int?),
      interlocutorAvatarMediaId:
          interlocutorAvatarMediaId ?? this.interlocutorAvatarMediaId,
    );
  }

  @override
  List<Object?> get props => [
        status,
        messages,
        currentUserId,
        isInterlocutorTyping,
        isInterlocutorOnline,
        interlocutorLastSeenAt,
        replyTarget?.id,
        hasMoreMessages,
        isLoadingMore,
        blockedMe,
        blockedByMe,
        isGatingRequired,
        members,
        memberAvatarUrls,
        highlightedMessageId,
        interlocutorAvatarMediaId,
      ];
}

const _sentinel = Object();

// --- BLOC ---
@injectable
class ChatDetailBloc extends Bloc<ChatDetailEvent, ChatDetailState> {
  final ChatRepository _repo;
  final MediaRepository _mediaRepo;
  final SocketClient _socket;
  final SecureStorageClient _storage;
  final AppLogger _logger;
  final PreferencesRepository _prefsRepo;

  String chatId = '';
  String _interlocutorId = '';
  String _chatType = 'DIRECT';
  StreamSubscription? _socketSub;
  Timer? _typingDebounce;
  // Loaded once per chat session from user preferences.
  bool _readReceiptsEnabled = true;

  ChatDetailBloc(
    this._repo,
    this._mediaRepo,
    this._socket,
    this._storage,
    this._logger,
    this._prefsRepo,
  ) : super(const ChatDetailState()) {
    on<LoadMessages>(_onLoadMessages);
    on<LoadGroupMembers>(_onLoadGroupMembers);
    on<SendMessage>(_onSendMessage);
    on<SendMediaMessage>(_onSendMediaMessage);
    on<EditMessageEvent>(_onEditMessage);
    on<DeleteMessageEvent>(_onDeleteMessage);
    on<_OnSocketEvent>(_onSocketEvent);
    on<OnTypingInput>(_onTypingInput);
    on<SetReplyTarget>(_onSetReplyTarget);
    on<ClearReplyTarget>(_onClearReplyTarget);
    on<LoadMoreMessages>(_onLoadMoreMessages);
    on<JumpToMessage>(_onJumpToMessage);
    on<ClearHighlight>(_onClearHighlight);
  }

  void init(String id, {String interlocutorId = '', String chatType = 'DIRECT'}) {
    chatId = id;
    _interlocutorId = interlocutorId;
    _chatType = chatType;
    _logger.info('ChatDetailBloc init: chatId=$chatId, type=$chatType, interlocutorId=$interlocutorId');
    add(LoadMessages());
    if (_chatType == 'GROUP') add(LoadGroupMembers());

    _socketSub = _socket.events.listen((event) {
      add(_OnSocketEvent(event));
    });
  }

  /// Loads group members and resolves their avatar URLs so each message can
  /// show the sender's display name + avatar.
  Future<void> _onLoadGroupMembers(LoadGroupMembers event, Emitter<ChatDetailState> emit) async {
    final result = await _repo.getChatInfo(chatId);
    await result.fold(
      (f) async => _logger.warning('Failed to load group members: ${f.message}'),
      (info) async {
        final membersById = {for (final m in info.members) m.userId: m};
        // Resolve avatar URLs in parallel for members that have an avatar.
        final withAvatars = info.members.where((m) => m.avatarMediaId != null).toList();
        final urls = await Future.wait(withAvatars.map((m) async {
          final r = await _mediaRepo.getMediaUrl(m.avatarMediaId!);
          return MapEntry(m.userId, r.fold((_) => null, (u) => u));
        }));
        if (emit.isDone) return;
        emit(state.copyWith(
          members: membersById,
          memberAvatarUrls: {for (final e in urls) e.key: e.value},
        ));
      },
    );
  }

  Future<void> _onLoadMessages(LoadMessages event, Emitter<ChatDetailState> emit) async {
    emit(state.copyWith(status: ChatStatus.loading));
    _logger.debug('Loading messages for chat $chatId');

    final token = await _storage.getAccessToken();
    String? myId;
    if (token != null) {
      try {
        final decoded = JwtDecoder.decode(token);
        myId = decoded['user_id']?.toString() ?? decoded['sub']?.toString();
      } catch (e) {
        _logger.error('Failed to decode JWT in ChatDetailBloc', e);
      }
    }

    final messagesFuture = _repo.getMessages(chatId, limit: 50);
    final statusFuture = _interlocutorId.isNotEmpty
        ? _repo.getUserStatus(_interlocutorId)
        : Future.value(null);
    // DIRECT: fetch block status proactively (§4.2).
    // GROUP: fetch to resolve the group avatar for the header.
    final chatInfoFuture = (_chatType == 'DIRECT' || _chatType == 'GROUP')
        ? _repo.getChatInfo(chatId)
        : Future.value(null);
    final prefsFuture = _prefsRepo.getPreferences();

    final results = await Future.wait([messagesFuture, statusFuture, chatInfoFuture, prefsFuture]);
    final messagesResult = results[0] as dynamic;
    final statusResult = results[1] as dynamic;
    final chatInfoResult = results[2] as dynamic;
    final prefsResult = results[3] as dynamic;
    prefsResult.fold(
      (f) => _logger.warning('Failed to load preferences for read-receipt check: ${f.message}'),
      (prefs) => _readReceiptsEnabled = prefs.readReceiptsEnabled,
    );

    messagesResult.fold(
      (l) {
        if (l is ServerFailure && l.code == 'GATING_REQUIRED') {
          // 403 GATING_REQUIRED = token-gated public channel, insufficient balance.
          _logger.info('Chat $chatId returned 403 GATING_REQUIRED — token gating blocks access');
          emit(state.copyWith(
            status: ChatStatus.success,
            messages: const [],
            currentUserId: myId,
            isGatingRequired: true,
          ));
          return;
        }
        // BLOCKED = interlocutor has blocked the current user (new backend code §4.2).
        // FORBIDDEN kept as fallback for older server versions.
        if (l is ServerFailure && (l.code == 'BLOCKED' || l.code == 'FORBIDDEN')) {
          _logger.info('Chat $chatId returned ${l.code} — showing blocked state');
          emit(state.copyWith(
            status: ChatStatus.success,
            messages: const [],
            currentUserId: myId,
            blockedMe: true,
          ));
          return;
        }
        if (l is ServerFailure &&
            (l.code == '404' || l.code == 'NOT_FOUND' || l.message.contains('not found'))) {
          _logger.info('Chat $chatId returned 404 — treating as empty');
          emit(state.copyWith(status: ChatStatus.success, messages: [], currentUserId: myId));
        } else {
          _logger.error('Failed to load messages for chat $chatId', l);
          emit(state.copyWith(status: ChatStatus.failure));
          return;
        }
      },
      (msgs) {
        _logger.info('Loaded ${msgs.length} messages for chat $chatId');
        emit(state.copyWith(
          status: ChatStatus.success,
          messages: msgs,
          currentUserId: myId,
          hasMoreMessages: msgs.length == 50,
        ));
        if (msgs.isNotEmpty) {
          _sendReadReceipt(msgs.first.id);
        }
      },
    );

    if (statusResult != null) {
      statusResult.fold(
        (l) => _logger.warning('Failed to get interlocutor status: ${l.message}'),
        (s) {
          _logger.debug('Interlocutor online=${s.isOnline}, lastSeen=${s.lastSeenAt}');
          emit(state.copyWith(
            isInterlocutorOnline: s.isOnline,
            interlocutorLastSeenAt: s.lastSeenAt,
            clearLastSeen: s.isOnline,
          ));
        },
      );
    }

    // Apply data from ChatInfo.
    if (chatInfoResult != null) {
      chatInfoResult.fold(
        (l) => _logger.warning('Failed to load chat info: ${l.message}'),
        (info) {
          if (_chatType == 'GROUP') {
            // For GROUP chats: propagate the group's own avatar to the header.
            _logger.debug('Chat $chatId group avatar: ${info.avatarMediaId}');
            emit(state.copyWith(interlocutorAvatarMediaId: info.avatarMediaId));
            return;
          }
          // DIRECT chats — block status + interlocutor avatar (§4.2).
          _logger.debug('Chat $chatId block status: blockedByMe=${info.blockedByMe}, blockedMe=${info.blockedMe}');
          String? interlocutorAvatarId;
          for (final m in info.members) {
            if (m.userId == _interlocutorId) {
              interlocutorAvatarId = m.avatarMediaId;
              break;
            }
          }
          emit(state.copyWith(
            blockedByMe: info.blockedByMe,
            blockedMe: info.blockedMe,
            interlocutorAvatarMediaId: interlocutorAvatarId,
          ));
        },
      );
    }
  }

  void _sendReadReceipt(int messageId) {
    if (!_readReceiptsEnabled) {
      _logger.debug('Read receipts disabled — skipping READ_MESSAGE for msg $messageId');
      return;
    }
    _socket.send(SocketEventType.readMessage, {
      'chat_id': chatId,
      'message_id': messageId,
    });
    _logger.debug('Sent READ_MESSAGE for msg $messageId in chat $chatId');
  }

  Future<void> _onSendMessage(SendMessage event, Emitter<ChatDetailState> emit) async {
    // Negative timestamp — guaranteed never to collide with server-assigned positive int IDs
    final tempId = -DateTime.now().millisecondsSinceEpoch;
    final tempMsg = Message(
      id: tempId,
      chatId: chatId,
      senderId: state.currentUserId ?? '',
      content: event.content,
      type: event.type,
      status: MessageStatus.sending,
      createdAt: DateTime.now().toUtc().toIso8601String(),
    );

    final tempMessages = List<Message>.from(state.messages)..insert(0, tempMsg);
    emit(state.copyWith(messages: tempMessages));

    final replyId = state.replyTarget?.id;
    final result = await _repo.sendMessage(
      chatId: chatId,
      content: event.content,
      type: event.type,
      replyToMessageId: replyId,
    );
    result.fold(
      (failure) {
        _logger.error('Send message failed for chat $chatId', failure);
        final cleanMessages = state.messages.where((m) => m.id != tempId).toList();

        // 400 BAD_REQUEST on reply = original was deleted or belongs to another chat (§3.2).
        if (failure.code == 'BAD_REQUEST' && state.replyTarget != null) {
          _logger.info('Reply target invalid (deleted/wrong chat) — clearing reply target');
          emit(state.copyWith(messages: cleanMessages, replyTarget: null));
          return;
        }

        // 403 BLOCKED/FORBIDDEN = interlocutor has blocked the current user (§4.2).
        if (failure.code == 'BLOCKED' || failure.code == 'FORBIDDEN') {
          _logger.info('Send failed with ${failure.code} — showing blocked state');
          emit(state.copyWith(messages: cleanMessages, blockedMe: true));
          return;
        }

        final failedMessages = state.messages
            .map((m) => m.id == tempId ? m.copyWith(status: MessageStatus.failed) : m)
            .toList();
        emit(state.copyWith(messages: failedMessages));
      },
      (realMessage) {
        _logger.debug('Message sent, replacing temp $tempId with real ${realMessage.id}');
        final successMessages = state.messages.map((m) => m.id == tempId ? realMessage : m).toList();
        emit(state.copyWith(messages: successMessages, replyTarget: null));
      },
    );
  }

  Future<void> _onSendMediaMessage(
    SendMediaMessage event,
    Emitter<ChatDetailState> emit,
  ) async {
    _logger.info('Sending ${event.messageType} from ${event.filePath}');

    // For FILE messages preserve the original filename in `content` so the
    // bubble can display it (the server echoes content back).
    final isFile = event.messageType.toUpperCase() == 'FILE';
    final fileName = isFile ? Uri.file(event.filePath).pathSegments.last : null;
    final messageContent = event.content ?? fileName ?? '';

    // Optimistic bubble — shows "Uploading…" while upload is in progress
    final tempId = -DateTime.now().millisecondsSinceEpoch;
    final tempMsg = Message(
      id: tempId,
      chatId: chatId,
      senderId: state.currentUserId ?? '',
      content: isFile ? fileName! : (event.content ?? 'Uploading…'),
      type: event.messageType,
      status: MessageStatus.sending,
      createdAt: DateTime.now().toUtc().toIso8601String(),
    );
    emit(state.copyWith(
      messages: [tempMsg, ...state.messages],
    ));

    final uploadResult = await _mediaRepo.uploadFile(
      filePath: event.filePath,
      mediaType: event.messageType,
    );

    await uploadResult.fold(
      (failure) async {
        _logger.error('Media upload failed for ${event.messageType}', failure);
        final failed = state.messages
            .map((m) => m.id == tempId ? m.copyWith(status: MessageStatus.failed) : m)
            .toList();
        emit(state.copyWith(messages: failed));
      },
      (mediaId) async {
        _logger.info('Upload done: mediaId=$mediaId, sending ${event.messageType} message');
        final result = await _repo.sendMessage(
          chatId: chatId,
          content: messageContent,
          type: event.messageType.toLowerCase(),
          attachmentMediaId: mediaId,
        );
        result.fold(
          (failure) {
            _logger.error('Send ${event.messageType} message failed', failure);
            final failed = state.messages
                .map((m) => m.id == tempId ? m.copyWith(status: MessageStatus.failed) : m)
                .toList();
            emit(state.copyWith(messages: failed));
          },
          (realMessage) {
            final updated = state.messages
                .map((m) => m.id == tempId ? realMessage : m)
                .toList();
            emit(state.copyWith(messages: updated));
          },
        );
      },
    );
  }

  Future<void> _onEditMessage(EditMessageEvent event, Emitter<ChatDetailState> emit) async {
    _logger.debug('Editing message ${event.messageId} in chat $chatId');

    final now = DateTime.now().toUtc().toIso8601String();
    final optimisticMessages = state.messages.map((m) {
      if (m.id == event.messageId) return m.copyWith(content: event.newContent, editedAt: now);
      return m;
    }).toList();
    emit(state.copyWith(messages: optimisticMessages));

    final result = await _repo.editMessage(event.messageId, event.newContent);
    result.fold(
      (failure) {
        _logger.error('Edit message ${event.messageId} failed', failure);
        add(LoadMessages());
      },
      (_) => _logger.info('Message ${event.messageId} edited successfully'),
    );
  }

  Future<void> _onDeleteMessage(DeleteMessageEvent event, Emitter<ChatDetailState> emit) async {
    _logger.debug('Deleting message ${event.messageId}, forEveryone=${event.forEveryone}');

    final optimisticMessages = state.messages.where((m) => m.id != event.messageId).toList();
    emit(state.copyWith(messages: optimisticMessages));

    final result = await _repo.deleteMessage(event.messageId, forEveryone: event.forEveryone);
    result.fold(
      (failure) {
        _logger.error('Delete message ${event.messageId} failed', failure);
        add(LoadMessages());
      },
      (_) => _logger.info('Message ${event.messageId} deleted successfully'),
    );
  }

  void _onSocketEvent(_OnSocketEvent event, Emitter<ChatDetailState> emit) {
    final payload = event.event.payload;
    if (payload == null) return;

    switch (event.event.type) {
      case SocketEventType.newMessage:
        final msgChatId = payload['chat_id'] as String?;
        if (msgChatId != chatId) return;

        try {
          final incomingMessage = Message.fromJson(payload);
          // Skip own messages — the REST response already added them optimistically.
          if (incomingMessage.senderId == state.currentUserId) return;
          // Dedup: skip if we already have a message with this server id.
          if (state.messages.any((m) => m.id == incomingMessage.id)) {
            _logger.debug('WS NEW_MESSAGE dedup: id=${incomingMessage.id} already in list');
            return;
          }
          _logger.info('WS NEW_MESSAGE appending id=${incomingMessage.id} to chat $chatId');
          final updatedMessages = List<Message>.from(state.messages)..insert(0, incomingMessage);
          // The sender just sent a message → they are no longer typing. Clear the
          // indicator immediately instead of waiting for the debounce to expire.
          _typingDebounce?.cancel();
          emit(state.copyWith(messages: updatedMessages, isInterlocutorTyping: false));
          _sendReadReceipt(incomingMessage.id);
        } catch (e) {
          _logger.error('Failed to parse NEW_MESSAGE payload', e);
        }
        break;

      case SocketEventType.messagesRead:
        final readChatId = payload['chat_id'] as String?;
        if (readChatId != null && readChatId != chatId) return;

        final lastReadId = payload['last_read_id'] as int?;
        final byUser = payload['by_user'] as String?;
        if (lastReadId == null || byUser == state.currentUserId) return;

        final updatedMessages = state.messages.map((m) {
          if (m.senderId == state.currentUserId &&
              m.id <= lastReadId &&
              m.status != MessageStatus.read) {
            return m.copyWith(status: MessageStatus.read);
          }
          return m;
        }).toList();
        emit(state.copyWith(messages: updatedMessages));
        _logger.debug('Messages up to $lastReadId marked READ by $byUser');
        break;

      case SocketEventType.messageEdited:
        final msgChatId = payload['chat_id'] as String?;
        if (msgChatId != null && msgChatId != chatId) return;

        final messageId = payload['message_id'] as int?;
        final newContent = payload['content'] as String?;
        final editedAt = payload['edited_at'] as String?;
        if (messageId == null || newContent == null) return;

        final updatedMessages = state.messages.map((m) {
          if (m.id == messageId) {
            return m.copyWith(
              content: newContent,
              editedAt: editedAt ?? DateTime.now().toUtc().toIso8601String(),
            );
          }
          return m;
        }).toList();
        emit(state.copyWith(messages: updatedMessages));
        break;

      case SocketEventType.messageDeleted:
        final msgChatId = payload['chat_id'] as String?;
        if (msgChatId != null && msgChatId != chatId) return;

        final messageId = payload['message_id'] as int?;
        if (messageId == null) return;

        emit(state.copyWith(
          messages: state.messages.where((m) => m.id != messageId).toList(),
        ));
        break;

      case SocketEventType.userOnline:
        final userId = payload['user_id'] as String?;
        if (userId != _interlocutorId || _interlocutorId.isEmpty) return;

        emit(state.copyWith(isInterlocutorOnline: true, clearLastSeen: true));
        _logger.debug('Interlocutor $_interlocutorId came online');
        break;

      case SocketEventType.userOffline:
        final userId = payload['user_id'] as String?;
        if (userId != _interlocutorId || _interlocutorId.isEmpty) return;

        _logger.debug('Interlocutor $_interlocutorId went offline');
        _repo.getUserStatus(_interlocutorId).then((result) {
          if (isClosed) return; // guard must be first
          result.fold(
            (l) => _logger.warning('getUserStatus after offline failed: ${l.message}'),
            (s) {
              if (isClosed) return;
              emit(state.copyWith(
                isInterlocutorOnline: false,
                interlocutorLastSeenAt: s.lastSeenAt,
              ));
            },
          );
        });
        break;

      case SocketEventType.typing:
        final typingChatId = payload['chat_id'] as String?;
        if (typingChatId != chatId) return;

        // Honor an explicit stop flag from the peer.
        if (payload['stopped'] == true) {
          _typingDebounce?.cancel();
          emit(state.copyWith(isInterlocutorTyping: false));
          break;
        }

        emit(state.copyWith(isInterlocutorTyping: true));
        _typingDebounce?.cancel();
        _typingDebounce = Timer(const Duration(seconds: 3), () {
          if (!isClosed) emit(state.copyWith(isInterlocutorTyping: false));
        });
        break;

      // Real-time block/unblock events (§4.4).
      case SocketEventType.userBlocked:
        final blockedBy = payload['by_user'] as String?;
        if (blockedBy == _interlocutorId && _interlocutorId.isNotEmpty) {
          _logger.info('USER_BLOCKED by interlocutor $_interlocutorId — locking input');
          emit(state.copyWith(blockedMe: true));
        }
        break;

      case SocketEventType.userUnblocked:
        final unblockedBy = payload['by_user'] as String?;
        if (unblockedBy == _interlocutorId && _interlocutorId.isNotEmpty) {
          _logger.info('USER_UNBLOCKED by interlocutor $_interlocutorId — unlocking input');
          emit(state.copyWith(blockedMe: false));
        }
        break;

      // In-chat transfer confirmed by observer — reload so the SUDA_TRANSFER
      // system message (created by messenger) appears in the list.
      case SocketEventType.sudaSent:
      case SocketEventType.sudaReceived:
        final txChatId = payload['chat_id'] as String?;
        if (txChatId != null && txChatId == chatId) {
          _logger.info('WS ${ event.event.type} received for chat $chatId — reloading messages');
          add(LoadMessages());
        }
        break;

      case SocketEventType.donationSent:
      case SocketEventType.donationReceived:
        final donChatId = payload['chat_id'] as String?;
        if (donChatId != null && donChatId == chatId) {
          _logger.info('WS ${event.event.type} received for chat $chatId — reloading messages');
          add(LoadMessages());
        }
        break;
    }
  }

  void _onTypingInput(OnTypingInput event, Emitter<ChatDetailState> emit) {
    // Send a stop flag too (best-effort) so the peer can clear instantly if the
    // backend forwards it; otherwise the peer falls back to its debounce timer.
    _socket.send(SocketEventType.typing, {
      'chat_id': chatId,
      'stopped': !event.isTyping,
    });
  }

  void _onSetReplyTarget(SetReplyTarget event, Emitter<ChatDetailState> emit) {
    // Guard: optimistic messages have negative ids — they haven't reached the
    // server yet and cannot be used as a reply target (§3.2).
    if (event.message.id <= 0) {
      _logger.warning('SetReplyTarget ignored: message id=${event.message.id} is not a server id');
      return;
    }
    emit(state.copyWith(replyTarget: event.message));
  }

  void _onClearReplyTarget(ClearReplyTarget event, Emitter<ChatDetailState> emit) {
    emit(state.copyWith(replyTarget: null));
  }

  Future<void> _onLoadMoreMessages(
    LoadMoreMessages event,
    Emitter<ChatDetailState> emit,
  ) async {
    if (state.isLoadingMore) return;
    emit(state.copyWith(isLoadingMore: true));
    _logger.debug('Loading more messages before ${event.beforeMessageId} in chat $chatId');

    final result = await _repo.getMessages(chatId, limit: 30, beforeMessageId: event.beforeMessageId);
    result.fold(
      (failure) {
        _logger.error('LoadMore failed for chat $chatId', failure);
        emit(state.copyWith(isLoadingMore: false));
      },
      (older) {
        _logger.info('Loaded ${older.length} older messages for chat $chatId');
        emit(state.copyWith(
          messages: [...state.messages, ...older],
          hasMoreMessages: older.length == 30,
          isLoadingMore: false,
        ));
      },
    );
  }

  /// Search jump-to: pages older messages until [event.messageId] is loaded
  /// (bounded), then flags it as the highlight target. The page scrolls to it.
  Future<void> _onJumpToMessage(
    JumpToMessage event,
    Emitter<ChatDetailState> emit,
  ) async {
    _logger.debug('JumpToMessage ${event.messageId} in chat $chatId');

    // Page backwards until the target is loaded or there is nothing more.
    // Bounded to avoid an unbounded request loop on a stale/invalid id.
    const maxPages = 20;
    var pages = 0;
    while (!state.messages.any((m) => m.id == event.messageId) &&
        state.hasMoreMessages &&
        pages < maxPages) {
      final oldestId = state.messages.isNotEmpty ? state.messages.last.id : null;
      if (oldestId == null) break;
      final result =
          await _repo.getMessages(chatId, limit: 30, beforeMessageId: oldestId);
      final shouldStop = result.fold(
        (failure) {
          _logger.error('JumpToMessage paging failed for chat $chatId', failure);
          return true;
        },
        (older) {
          if (emit.isDone) return true;
          emit(state.copyWith(
            messages: [...state.messages, ...older],
            hasMoreMessages: older.length == 30,
          ));
          return older.isEmpty;
        },
      );
      if (shouldStop) break;
      pages++;
    }

    if (emit.isDone) return;
    if (!state.messages.any((m) => m.id == event.messageId)) {
      _logger.warning('JumpToMessage ${event.messageId} not found in chat $chatId');
      return;
    }

    emit(state.copyWith(highlightedMessageId: event.messageId));
    // Clear the flash after a short delay.
    Future.delayed(const Duration(milliseconds: 1500), () {
      if (!isClosed) add(ClearHighlight());
    });
  }

  void _onClearHighlight(ClearHighlight event, Emitter<ChatDetailState> emit) {
    emit(state.copyWith(highlightedMessageId: null));
  }

  @override
  Future<void> close() {
    _logger.info('ChatDetailBloc closed: chatId=$chatId');
    _socketSub?.cancel();
    _typingDebounce?.cancel();
    return super.close();
  }
}
