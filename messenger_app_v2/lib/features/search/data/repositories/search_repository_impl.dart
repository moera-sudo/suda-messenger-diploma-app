import 'package:dartz/dartz.dart';
import 'package:injectable/injectable.dart';

import '../../../../shared/data/api/api_client.dart';
import '../../../../shared/data/api/server_exception.dart';
import '../../../../shared/data/logger/app_logger.dart';
import '../../../../shared/domain/models/app_failure.dart';
import '../../domain/repositories/search_repository.dart';
import '../models/search_result.dart';

@LazySingleton(as: SearchRepository)
class SearchRepositoryImpl implements SearchRepository {
  final ApiClient _api;
  final AppLogger _logger;

  static const _searchPath = '/api/v1/messenger/search';
  static const _chatsPath = '/api/v1/messenger/chats';

  SearchRepositoryImpl(this._api, this._logger);

  @override
  Future<Either<AppFailure, List<SearchResult>>> search(String query) async {
    try {
      final response = await _api.get(_searchPath, query: {'q': query});

      final List rawUsers    = response['users']    as List? ?? [];
      final List rawChats    = response['chats']    as List? ?? [];
      final List rawMessages = response['messages'] as List? ?? [];

      final users = rawUsers.map((e) {
        final m = Map<String, dynamic>.from(e as Map);
        return SearchResult(
          type: SearchResultType.user,
          id: m['id']?.toString() ?? '',
          title: m['title'] as String? ?? '',
          description: m['subtitle'] as String? ?? '',
          imageUrl: m['image_url'] as String? ?? '',
        );
      }).toList();

      final chats = rawChats.map((e) {
        final m = Map<String, dynamic>.from(e as Map);
        return SearchResult(
          type: SearchResultType.chat,
          id: m['id']?.toString() ?? '',
          title: m['title'] as String? ?? '',
          description: m['subtitle'] as String? ?? '',
          imageUrl: m['image_url'] as String? ?? '',
          // Real chat type (DIRECT/GROUP/CHANNEL/SAVED) so the UI can route
          // channels to the channel screen instead of treating them as groups.
          chatType: m['type'] as String?,
        );
      }).toList();

      final msgs = rawMessages.map((e) {
        final m = Map<String, dynamic>.from(e as Map);
        return SearchResult(
          type: SearchResultType.message,
          id: m['id']?.toString() ?? m['message_id']?.toString() ?? '',
          title: m['sender_name'] as String? ?? m['sender_id']?.toString() ?? '',
          description: m['content'] as String? ?? m['text'] as String? ?? '',
          chatId: m['chat_id']?.toString(),
          messageId: m['id'] is int
              ? m['id'] as int
              : int.tryParse(m['id']?.toString() ?? ''),
        );
      }).toList();

      final results = [...users, ...chats, ...msgs];
      _logger.info('Search "$query" returned ${results.length} results'
          ' (users:${users.length} chats:${chats.length} msgs:${msgs.length})');
      return Right(results);
    } catch (e) {
      _logger.error('search failed (query=$query)', e);
      if (e is ServerException) return Left(ServerFailure(message: e.message));
      return const Left(UnknownFailure());
    }
  }

  @override
  Future<Either<AppFailure, List<SearchResult>>> searchChatMessages(
    String chatId,
    String query,
  ) async {
    try {
      final response = await _api.get(
        '$_chatsPath/$chatId/search',
        query: {'q': query},
      );
      // Server may return {"messages": [...]}, {"results": [...]} or a plain list.
      final List raw;
      if (response is Map && response.containsKey('messages')) {
        raw = response['messages'] as List? ?? [];
      } else if (response is Map && response.containsKey('results')) {
        raw = response['results'] as List? ?? [];
      } else {
        raw = response as List? ?? [];
      }

      final results = raw.map((e) {
        final m = Map<String, dynamic>.from(e as Map);
        // Normalise message search results into SearchResult shape. The server
        // shape uses title (sender name) + subtitle (matched text); keep older
        // field names (sender_name/content) as fallbacks.
        return SearchResult(
          type: SearchResultType.message,
          id: m['id']?.toString() ?? m['message_id']?.toString() ?? '',
          title: m['title'] as String? ??
              m['sender_name'] as String? ??
              m['display_name'] as String? ??
              m['sender_id']?.toString() ?? '',
          description: m['subtitle'] as String? ??
              m['content'] as String? ??
              m['text'] as String? ?? '',
          chatId: chatId,
          messageId: m['message_id'] is int
              ? m['message_id'] as int
              : m['id'] is int
                  ? m['id'] as int
                  : int.tryParse(
                      (m['message_id'] ?? m['id'])?.toString() ?? ''),
        );
      }).toList();

      _logger.info(
          'Chat search "$query" in $chatId returned ${results.length} messages');
      return Right(results);
    } catch (e) {
      _logger.error('searchChatMessages failed (chatId=$chatId, q=$query)', e);
      if (e is ServerException) return Left(ServerFailure(message: e.message));
      return const Left(UnknownFailure());
    }
  }
}
