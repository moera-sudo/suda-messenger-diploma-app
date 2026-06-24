import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:injectable/injectable.dart';

import '../../../../shared/data/logger/app_logger.dart';
import '../../data/models/search_result.dart';
import '../../domain/repositories/search_repository.dart';
import '../bloc/search_bloc.dart' show debounce;

// --- EVENTS ---
abstract class ChatSearchEvent extends Equatable {}

class ChatSearchQueryChanged extends ChatSearchEvent {
  final String query;
  ChatSearchQueryChanged(this.query);
  @override
  List<Object> get props => [query];
}

class ClearChatSearch extends ChatSearchEvent {
  @override
  List<Object> get props => [];
}

// --- STATE ---
enum ChatSearchStatus { initial, loading, success, failure, empty }

class ChatSearchState extends Equatable {
  final ChatSearchStatus status;
  final List<SearchResult> results;
  final String query;

  const ChatSearchState({
    this.status = ChatSearchStatus.initial,
    this.results = const [],
    this.query = '',
  });

  @override
  List<Object> get props => [status, results, query];
}

// --- BLOC ---
@injectable
class ChatSearchBloc extends Bloc<ChatSearchEvent, ChatSearchState> {
  final SearchRepository _repo;
  final AppLogger _logger;

  String chatId = '';

  ChatSearchBloc(this._repo, this._logger) : super(const ChatSearchState()) {
    on<ChatSearchQueryChanged>(
      _onQueryChanged,
      transformer: debounce(const Duration(milliseconds: 400)),
    );
    on<ClearChatSearch>(_onClear);
  }

  Future<void> _onQueryChanged(
    ChatSearchQueryChanged event,
    Emitter<ChatSearchState> emit,
  ) async {
    final query = event.query.trim();
    if (query.length < 2) {
      emit(const ChatSearchState(status: ChatSearchStatus.initial));
      return;
    }
    emit(ChatSearchState(status: ChatSearchStatus.loading, query: query));
    _logger.debug('Searching "$query" in chat $chatId');

    final result = await _repo.searchChatMessages(chatId, query);
    result.fold(
      (failure) {
        _logger.error('searchChatMessages failed', failure);
        emit(ChatSearchState(status: ChatSearchStatus.failure, query: query));
      },
      (results) {
        _logger.info('Found ${results.length} messages for "$query" in $chatId');
        emit(ChatSearchState(
          status: results.isEmpty ? ChatSearchStatus.empty : ChatSearchStatus.success,
          results: results,
          query: query,
        ));
      },
    );
  }

  void _onClear(ClearChatSearch event, Emitter<ChatSearchState> emit) {
    emit(const ChatSearchState(status: ChatSearchStatus.initial));
  }
}
