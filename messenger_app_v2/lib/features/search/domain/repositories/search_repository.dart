import 'package:dartz/dartz.dart';

import '../../../../shared/domain/models/app_failure.dart';
import '../../data/models/search_result.dart';

abstract class SearchRepository {
  /// GET /api/v1/messenger/search?q={query}
  Future<Either<AppFailure, List<SearchResult>>> search(String query);

  /// GET /api/v1/messenger/chats/{chatId}/search?q={query}
  Future<Either<AppFailure, List<SearchResult>>> searchChatMessages(
    String chatId,
    String query,
  );
}
