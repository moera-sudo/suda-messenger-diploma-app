import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:injectable/injectable.dart';
import 'package:stream_transform/stream_transform.dart';

import '../../../../features/chat/domain/repositories/chat_repository.dart';
import '../../data/models/search_result.dart';
import '../../domain/repositories/search_repository.dart';

EventTransformer<E> debounce<E>(Duration duration) {
  return (events, mapper) => events.debounce(duration).switchMap(mapper);
}

// --- EVENTS ---
abstract class SearchEvent extends Equatable {}

class SearchQueryChanged extends SearchEvent {
  final String query;
  SearchQueryChanged(this.query);
  @override
  List<Object> get props => [query];
}

class ClearSearch extends SearchEvent {
  @override
  List<Object> get props => [];
}

class CreateChatWithUser extends SearchEvent {
  final String userId;
  final String userName;
  CreateChatWithUser(this.userId, this.userName);
  @override
  List<Object> get props => [userId, userName];
}

// --- STATE ---
enum SearchStatus { initial, loading, success, failure, empty }

class SearchState extends Equatable {
  final SearchStatus status;
  final List<SearchResult> userResults;
  final List<SearchResult> chatResults;
  final List<SearchResult> messageResults;
  final String query;

  const SearchState({
    this.status = SearchStatus.initial,
    this.userResults = const [],
    this.chatResults = const [],
    this.messageResults = const [],
    this.query = '',
  });

  bool get hasResults =>
      userResults.isNotEmpty ||
      chatResults.isNotEmpty ||
      messageResults.isNotEmpty;

  @override
  List<Object> get props =>
      [status, userResults, chatResults, messageResults, query];
}

class ChatCreated extends SearchState {
  final String chatId;
  final String chatName;
  const ChatCreated(this.chatId, this.chatName);
  @override
  List<Object> get props => [chatId, chatName];
}

// --- BLOC ---
@injectable
class SearchBloc extends Bloc<SearchEvent, SearchState> {
  final SearchRepository _searchRepo;
  final ChatRepository _chatRepo;

  SearchBloc(this._searchRepo, this._chatRepo) : super(const SearchState()) {
    on<SearchQueryChanged>(
      _onSearchChanged,
      transformer: debounce(const Duration(milliseconds: 400)),
    );
    on<ClearSearch>(_onClearSearch);
    on<CreateChatWithUser>(_onCreateChat);
  }

  Future<void> _onSearchChanged(
    SearchQueryChanged event,
    Emitter<SearchState> emit,
  ) async {
    final query = event.query.trim();
    if (query.length < 2) {
      emit(const SearchState(status: SearchStatus.initial));
      return;
    }
    emit(SearchState(status: SearchStatus.loading, query: query));

    final result = await _searchRepo.search(query);
    result.fold(
      (l) => emit(SearchState(
        status: SearchStatus.failure,
        query: query,
      )),
      (all) {
        final users = all
            .where((r) => r.type == SearchResultType.user)
            .toList();
        final chats = all
            .where((r) => r.type == SearchResultType.chat)
            .toList();
        final messages = all
            .where((r) => r.type == SearchResultType.message)
            .toList();

        final isEmpty = users.isEmpty && chats.isEmpty && messages.isEmpty;
        emit(SearchState(
          status: isEmpty ? SearchStatus.empty : SearchStatus.success,
          userResults: users,
          chatResults: chats,
          messageResults: messages,
          query: query,
        ));
      },
    );
  }

  void _onClearSearch(ClearSearch event, Emitter<SearchState> emit) {
    emit(const SearchState(status: SearchStatus.initial));
  }

  Future<void> _onCreateChat(
    CreateChatWithUser event,
    Emitter<SearchState> emit,
  ) async {
    final currentState = state;
    emit(SearchState(
      status: SearchStatus.loading,
      query: currentState.query,
    ));

    final result = await _chatRepo.createChat(
      type: 'DIRECT',
      targetId: event.userId,
      name: event.userName,
    );
    result.fold(
      (l) => emit(currentState),
      (chat) {
        final name = (chat.name != null && chat.name!.isNotEmpty)
            ? chat.name!
            : event.userName;
        emit(ChatCreated(chat.id, name));
      },
    );
  }
}
