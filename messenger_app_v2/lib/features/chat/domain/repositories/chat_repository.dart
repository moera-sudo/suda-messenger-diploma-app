import 'package:dartz/dartz.dart';

import '../../../../shared/domain/models/app_failure.dart';
import '../../data/models/chat_models.dart';

abstract class ChatRepository {
  Future<Either<AppFailure, List<ChatListItem>>> getChats();

  Future<Either<AppFailure, Chat>> createChat({
    required String type,
    String? targetId,
    String? name,
    String? description,
    List<String>? memberIds,
    String? visibility,
    String? username,
  });

  /// Cursor-based: pass [beforeMessageId] to load older messages.
  Future<Either<AppFailure, List<Message>>> getMessages(
    String chatId, {
    int limit = 50,
    int? beforeMessageId,
  });

  /// GET /chats/{id}/media — media shared in a chat, grouped by kind (§5.3).
  Future<Either<AppFailure, ChatMedia>> getChatMedia(String chatId);

  Future<Either<AppFailure, Message>> sendMessage({
    required String chatId,
    required String content,
    required String type,
    String? clientSideId,
    int? replyToMessageId,
    String? attachmentMediaId,
  });

  Future<Either<AppFailure, Unit>> editMessage(int messageId, String content);

  Future<Either<AppFailure, Unit>> deleteMessage(int messageId, {bool forEveryone = false});

  /// GET /user/id/:id/status
  Future<Either<AppFailure, UserStatus>> getUserStatus(String userId);

  /// GET /messages/:id/readers
  Future<Either<AppFailure, List<MessageReader>>> getMessageReaders(int messageId);

  // ─── Chat info & management ────────────────────────────────

  /// GET /chats/{id}/info
  Future<Either<AppFailure, ChatInfo>> getChatInfo(String chatId);

  /// PUT /chats/{id}
  Future<Either<AppFailure, Unit>> updateChat(String chatId,
      {String? name, String? description, String? avatarMediaId});

  /// DELETE /chats/{id}?scope=everyone|me
  Future<Either<AppFailure, Unit>> deleteChat(String chatId, {String scope = 'everyone'});

  /// POST /chats/{id}/leave
  Future<Either<AppFailure, Unit>> leaveChat(String chatId);

  // ─── Members ───────────────────────────────────────────────

  /// POST /chats/{id}/members
  Future<Either<AppFailure, Unit>> addMember(String chatId, String userId);

  /// DELETE /chats/{id}/members/{userId}
  Future<Either<AppFailure, Unit>> removeMember(String chatId, String userId);

  /// PUT /chats/{id}/members/{userId}/role
  Future<Either<AppFailure, Unit>> changeMemberRole(String chatId, String userId, String role);

  // ─── Mute ──────────────────────────────────────────────────

  /// POST /chats/{id}/mute
  Future<Either<AppFailure, Unit>> muteChat(String chatId);

  /// POST /chats/{id}/unmute
  Future<Either<AppFailure, Unit>> unmuteChat(String chatId);

  // ─── Forward ───────────────────────────────────────────────

  /// POST /chats/forward
  Future<Either<AppFailure, Unit>> forwardMessage({
    required String fromChatId,
    required String toChatId,
    required int messageId,
  });

  // ─── Pins ──────────────────────────────────────────────────

  /// POST /api/v1/messenger/user/pins
  Future<Either<AppFailure, String>> pinChat(String chatId);

  /// DELETE /api/v1/messenger/user/pins/{pinId}
  Future<Either<AppFailure, Unit>> unpinChat(String pinId);

  /// GET /api/v1/messenger/user/pins?type=CHATLIST → Map<chatId, pinId>
  Future<Either<AppFailure, Map<String, String>>> getPinIds();
}
