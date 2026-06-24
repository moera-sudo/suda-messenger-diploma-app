import 'package:dartz/dartz.dart';
import 'package:injectable/injectable.dart';
import 'package:uuid/uuid.dart';

import '../../../../shared/data/api/api_client.dart';
import '../../../../shared/data/api/server_exception.dart';
import '../../../../shared/data/logger/app_logger.dart';
import '../../../../shared/domain/models/app_failure.dart';
import '../../domain/repositories/chat_repository.dart';
import '../models/chat_models.dart';

@LazySingleton(as: ChatRepository)
class ChatRepositoryImpl implements ChatRepository {
  final ApiClient _api;
  final AppLogger _logger;
  static const _basePath = '/api/v1/messenger/chats';
  static const _msgPath = '/api/v1/messenger/messages';

  ChatRepositoryImpl(this._api, this._logger);

  // ─── Chats ─────────────────────────────────────────────────

  @override
  Future<Either<AppFailure, List<ChatListItem>>> getChats() async {
    try {
      final response = await _api.get(_basePath);
      final List data = response as List;
      return Right(data.map((e) => ChatListItem.fromJson(e)).toList());
    } catch (e) {
      _logger.error('getChats failed', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Chat>> createChat({
    required String type,
    String? targetId,
    String? name,
    String? description,
    List<String>? memberIds,
    String? visibility,
    String? username,
  }) async {
    try {
      final req = CreateChatRequest(
        type: type,
        targetId: targetId,
        name: name,
        description: description,
        memberIds: memberIds,
        visibility: visibility,
        username: username,
      );
      final response = await _api.post(_basePath, data: req.toJson());
      return Right(Chat.fromJson(response));
    } catch (e) {
      _logger.error('createChat failed (type=$type)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, ChatInfo>> getChatInfo(String chatId) async {
    try {
      final response = await _api.get('$_basePath/$chatId/info');
      return Right(ChatInfo.fromJson(Map<String, dynamic>.from(response)));
    } catch (e) {
      _logger.error('getChatInfo failed (chatId=$chatId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> updateChat(
    String chatId, {
    String? name,
    String? description,
    String? avatarMediaId,
  }) async {
    try {
      final body = <String, dynamic>{};
      if (name != null) body['name'] = name;
      if (description != null) body['description'] = description;
      if (avatarMediaId != null) body['avatar_media_id'] = avatarMediaId;
      await _api.put('$_basePath/$chatId', data: body);
      return const Right(unit);
    } catch (e) {
      _logger.error('updateChat failed (chatId=$chatId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> deleteChat(String chatId, {String scope = 'everyone'}) async {
    try {
      await _api.delete('$_basePath/$chatId', queryParameters: <String, dynamic>{'scope': scope});
      return const Right(unit);
    } catch (e) {
      _logger.error('deleteChat failed (chatId=$chatId, scope=$scope)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> leaveChat(String chatId) async {
    try {
      await _api.post('$_basePath/$chatId/leave', data: {});
      return const Right(unit);
    } catch (e) {
      _logger.error('leaveChat failed (chatId=$chatId)', e);
      return Left(_handleError(e));
    }
  }

  // ─── Messages ──────────────────────────────────────────────

  @override
  Future<Either<AppFailure, List<Message>>> getMessages(
    String chatId, {
    int limit = 50,
    int? beforeMessageId,
  }) async {
    try {
      final query = <String, dynamic>{'limit': limit};
      if (beforeMessageId != null) query['before_message_id'] = beforeMessageId;

      final response = await _api.get('$_basePath/$chatId/messages', query: query);
      if (response == null) return const Right([]);

      final List data = (response is Map && response.containsKey('messages'))
          ? response['messages'] as List
          : (response is List ? response : []);

      return Right(data.map((e) => Message.fromJson(e)).toList());
    } catch (e) {
      _logger.error('getMessages failed (chatId=$chatId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, ChatMedia>> getChatMedia(String chatId) async {
    try {
      final response = await _api.get('$_basePath/$chatId/media');
      final map = (response is Map)
          ? Map<String, dynamic>.from(response)
          : <String, dynamic>{};
      return Right(ChatMedia.fromJson(map));
    } catch (e) {
      _logger.error('getChatMedia failed (chatId=$chatId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Message>> sendMessage({
    required String chatId,
    required String content,
    required String type,
    String? clientSideId,
    int? replyToMessageId,
    String? attachmentMediaId,
  }) async {
    try {
      final req = SendMessageRequest(
        chatId: chatId,
        content: content,
        type: type,
        clientSideId: clientSideId ?? const Uuid().v4(),
        replyToMessageId: replyToMessageId,
        attachmentMediaId: attachmentMediaId,
      );
      final response = await _api.post('$_basePath/$chatId/messages', data: req.toJson());
      return Right(Message.fromJson(response));
    } catch (e) {
      _logger.error('sendMessage failed (chatId=$chatId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> editMessage(int messageId, String content) async {
    try {
      await _api.put('$_msgPath/$messageId', data: {'content': content});
      return const Right(unit);
    } catch (e) {
      _logger.error('editMessage failed (messageId=$messageId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> deleteMessage(int messageId, {bool forEveryone = false}) async {
    try {
      await _api.delete('$_msgPath/$messageId?for_everyone=$forEveryone');
      return const Right(unit);
    } catch (e) {
      _logger.error('deleteMessage failed (messageId=$messageId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> forwardMessage({
    required String fromChatId,
    required String toChatId,
    required int messageId,
  }) async {
    try {
      await _api.post('$_basePath/forward', data: {
        'from_chat_id': fromChatId,
        'to_chat_id': toChatId,
        'message_id': messageId,
      });
      return const Right(unit);
    } catch (e) {
      _logger.error('forwardMessage failed (messageId=$messageId)', e);
      return Left(_handleError(e));
    }
  }

  // ─── Members ───────────────────────────────────────────────

  @override
  Future<Either<AppFailure, Unit>> addMember(String chatId, String userId) async {
    try {
      await _api.post('$_basePath/$chatId/members', data: {'user_id': userId});
      return const Right(unit);
    } catch (e) {
      _logger.error('addMember failed (chatId=$chatId, userId=$userId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> removeMember(String chatId, String userId) async {
    try {
      await _api.delete('$_basePath/$chatId/members/$userId');
      return const Right(unit);
    } catch (e) {
      _logger.error('removeMember failed (chatId=$chatId, userId=$userId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> changeMemberRole(
    String chatId,
    String userId,
    String role,
  ) async {
    try {
      await _api.put('$_basePath/$chatId/members/$userId/role', data: {'role': role});
      return const Right(unit);
    } catch (e) {
      _logger.error('changeMemberRole failed (chatId=$chatId, userId=$userId)', e);
      return Left(_handleError(e));
    }
  }

  // ─── Mute ──────────────────────────────────────────────────

  @override
  Future<Either<AppFailure, Unit>> muteChat(String chatId) async {
    try {
      await _api.post('$_basePath/$chatId/mute', data: {'muted_until': null});
      return const Right(unit);
    } catch (e) {
      _logger.error('muteChat failed (chatId=$chatId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> unmuteChat(String chatId) async {
    try {
      await _api.post('$_basePath/$chatId/unmute', data: {});
      return const Right(unit);
    } catch (e) {
      _logger.error('unmuteChat failed (chatId=$chatId)', e);
      return Left(_handleError(e));
    }
  }

  // ─── Status & readers ──────────────────────────────────────

  @override
  Future<Either<AppFailure, UserStatus>> getUserStatus(String userId) async {
    try {
      final response = await _api.get('/api/v1/messenger/user/id/$userId/status');
      return Right(UserStatus.fromJson(Map<String, dynamic>.from(response)));
    } catch (e) {
      _logger.error('getUserStatus failed (userId=$userId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, List<MessageReader>>> getMessageReaders(int messageId) async {
    try {
      final response = await _api.get('$_msgPath/$messageId/readers');
      final List raw = response as List? ?? [];
      return Right(raw.map((e) => MessageReader.fromJson(Map<String, dynamic>.from(e))).toList());
    } catch (e) {
      _logger.error('getMessageReaders failed (messageId=$messageId)', e);
      return Left(_handleError(e));
    }
  }

  // ─── Pins ──────────────────────────────────────────────────

  @override
  Future<Either<AppFailure, String>> pinChat(String chatId) async {
    try {
      final response = await _api.post(
        '/api/v1/messenger/user/pins',
        data: {
          'pin_type': 'CHATLIST',
          'target_type': 'CHAT',
          'target_id': chatId,
        },
      );
      // Server returns { id, pin_type, target_type, target_id, order }
      final pinId = (response as Map<String, dynamic>?)?['id']?.toString() ?? '';
      return Right(pinId);
    } catch (e) {
      _logger.error('pinChat failed (chatId=$chatId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> unpinChat(String pinId) async {
    try {
      await _api.delete('/api/v1/messenger/user/pins/$pinId');
      return const Right(unit);
    } catch (e) {
      _logger.error('unpinChat failed (pinId=$pinId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Map<String, String>>> getPinIds() async {
    try {
      final response = await _api.get(
        '/api/v1/messenger/user/pins',
        query: {'type': 'CHATLIST'},
      );
      final List raw = response as List? ?? [];
      final map = <String, String>{};
      for (final item in raw) {
        final m = Map<String, dynamic>.from(item as Map);
        final pinId = m['id']?.toString();
        final chatId = m['target_id']?.toString();
        if (pinId != null && chatId != null) map[chatId] = pinId;
      }
      return Right(map);
    } catch (e) {
      _logger.error('getPinIds failed', e);
      return Left(_handleError(e));
    }
  }

  // ─── Error handling ────────────────────────────────────────

  AppFailure _handleError(dynamic e) {
    if (e is ServerException) return ServerFailure(message: e.message, code: e.errorCode);
    return const UnknownFailure();
  }
}
