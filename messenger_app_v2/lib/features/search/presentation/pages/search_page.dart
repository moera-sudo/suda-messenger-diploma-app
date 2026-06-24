import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import '../../../../app/config/theme/app_theme.dart';
import '../../../../app/navigation/app_routes.dart';
import '../../../../features/chat/data/models/chat_models.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/widgets/suda_avatar.dart';
import '../../data/models/search_result.dart';
import '../bloc/search_bloc.dart';
import '../widgets/highlight_text.dart';

class SearchPage extends StatefulWidget {
  const SearchPage({super.key});

  @override
  State<SearchPage> createState() => _SearchPageState();
}

class _SearchPageState extends State<SearchPage> {
  final _ctrl = TextEditingController();

  @override
  void dispose() {
    _ctrl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final recentChats =
        GoRouterState.of(context).extra as List<ChatListItem>? ?? [];

    return BlocProvider(
      create: (_) => sl<SearchBloc>(),
      child: _SearchView(ctrl: _ctrl, recentChats: recentChats),
    );
  }
}

class _SearchView extends StatelessWidget {
  final TextEditingController ctrl;
  final List<ChatListItem> recentChats;

  const _SearchView({required this.ctrl, required this.recentChats});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;

    return BlocListener<SearchBloc, SearchState>(
      listenWhen: (_, next) => next is ChatCreated,
      listener: (context, state) {
        if (state is ChatCreated) {
          context.pushReplacementNamed(
            AppRoutes.chatDetail,
            pathParameters: {'id': state.chatId},
            extra: {
              'name': state.chatName,
              'interlocutorId': '',
              'chatType': 'DIRECT',
            },
          );
        }
      },
      child: Scaffold(
        backgroundColor: theme.scaffoldBackgroundColor,
        appBar: AppBar(
          backgroundColor: theme.scaffoldBackgroundColor,
          elevation: 0,
          scrolledUnderElevation: 0,
          leading: IconButton(
            icon: Icon(Icons.arrow_back_rounded, color: theme.colorScheme.onSurface),
            onPressed: () => context.pop(),
          ),
          titleSpacing: 0,
          title: _InlineSearchField(ctrl: ctrl),
          bottom: PreferredSize(
            preferredSize: const Size.fromHeight(1),
            child: Divider(height: 1, color: palette.divider),
          ),
        ),
        body: BlocBuilder<SearchBloc, SearchState>(
          buildWhen: (prev, next) => next is! ChatCreated,
          builder: (context, state) {
            if (state.status == SearchStatus.loading) {
              return LinearProgressIndicator(
                color: theme.colorScheme.tertiary,
                backgroundColor: theme.colorScheme.surface,
              );
            }
            if (state.status == SearchStatus.failure) {
              return Center(
                child: Text(
                  context.l10n.errorGeneric,
                  style: TextStyle(color: palette.danger),
                ),
              );
            }
            if (state.status == SearchStatus.empty) {
              return _EmptyState(query: state.query);
            }
            if (state.status == SearchStatus.success) {
              return _ResultsSections(state: state);
            }
            // Initial state — show recent chats or hint
            return _RecentSection(recentChats: recentChats);
          },
        ),
      ),
    );
  }
}

// ─── Inline search field ──────────────────────────────────────

class _InlineSearchField extends StatelessWidget {
  final TextEditingController ctrl;
  const _InlineSearchField({required this.ctrl});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;

    return Padding(
      padding: const EdgeInsets.only(right: 8),
      child: Container(
        height: 40,
        decoration: BoxDecoration(
          color: theme.colorScheme.surface,
          borderRadius: BorderRadius.circular(12),
        ),
        child: Row(
          children: [
            const SizedBox(width: 12),
            Icon(Icons.search_rounded, size: 18, color: palette.textSecondary),
            const SizedBox(width: 8),
            Expanded(
              child: TextField(
                controller: ctrl,
                autofocus: true,
                style: TextStyle(
                  fontFamily: 'Manrope',
                  fontSize: 15,
                  color: theme.colorScheme.onSurface,
                ),
                cursorColor: theme.colorScheme.tertiary,
                decoration: InputDecoration(
                  hintText: context.l10n.searchHint,
                  hintStyle: TextStyle(color: palette.textSecondary),
                  border: InputBorder.none,
                  isDense: true,
                  contentPadding: EdgeInsets.zero,
                ),
                onChanged: (value) {
                  context.read<SearchBloc>().add(SearchQueryChanged(value));
                },
              ),
            ),
            ValueListenableBuilder<TextEditingValue>(
              valueListenable: ctrl,
              builder: (context, value, _) {
                if (value.text.isEmpty) return const SizedBox.shrink();
                return IconButton(
                  icon: Icon(Icons.close_rounded,
                      size: 18, color: palette.textSecondary),
                  padding: EdgeInsets.zero,
                  constraints: const BoxConstraints(
                    minWidth: 36,
                    minHeight: 36,
                  ),
                  onPressed: () {
                    ctrl.clear();
                    context.read<SearchBloc>().add(ClearSearch());
                  },
                );
              },
            ),
          ],
        ),
      ),
    );
  }
}

// ─── Recent section ───────────────────────────────────────────

class _RecentSection extends StatelessWidget {
  final List<ChatListItem> recentChats;
  const _RecentSection({required this.recentChats});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    if (recentChats.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              Icons.search_rounded,
              size: 64,
              color: theme.colorScheme.onSurface.withValues(alpha: 0.1),
            ),
            const SizedBox(height: 16),
            Text(
              l10n.searchEnterUsername,
              style: TextStyle(
                color: theme.colorScheme.onSurface.withValues(alpha: 0.3),
                fontFamily: 'Manrope',
              ),
            ),
          ],
        ),
      );
    }

    return ListView(
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(14, 12, 14, 8),
          child: Text(
            l10n.searchRecent,
            style: TextStyle(
              fontFamily: 'Manrope',
              fontSize: 12,
              fontWeight: FontWeight.w700,
              letterSpacing: 0.1 * 12,
              color: palette.textSecondary,
            ),
          ),
        ),
        ...recentChats.take(4).map((chat) {
          final title = switch (chat.type) {
            'SAVED' => l10n.chatSavedMessages,
            'GROUP' || 'CHANNEL' => chat.name ?? l10n.chatUnknownUser,
            _ => chat.interlocutorName ?? l10n.chatUnknownUser,
          };
          final initials = title.isNotEmpty ? title[0].toUpperCase() : '?';
          return InkWell(
            onTap: () => context.pushNamed(
              AppRoutes.chatDetail,
              pathParameters: {'id': chat.id},
              extra: {
                'name': title,
                'interlocutorId': chat.interlocutorId ?? '',
                'chatType': chat.type,
              },
            ),
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
              child: Row(
                children: [
                  SudaAvatar(initials: initials, size: 40),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          title,
                          style: TextStyle(
                            fontFamily: 'Manrope',
                            fontSize: 15,
                            fontWeight: FontWeight.w600,
                            color: theme.colorScheme.onSurface,
                          ),
                          overflow: TextOverflow.ellipsis,
                        ),
                        Text(
                          chat.lastMessageContent.isEmpty
                              ? l10n.chatNoMessagesYet
                              : chat.lastMessageContent,
                          style: TextStyle(
                            fontFamily: 'Manrope',
                            fontSize: 13,
                            color: palette.textSecondary,
                          ),
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          );
        }),
      ],
    );
  }
}

// ─── Results sections ─────────────────────────────────────────

class _ResultsSections extends StatelessWidget {
  final SearchState state;
  const _ResultsSections({required this.state});

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;

    return ListView(
      children: [
        if (state.userResults.isNotEmpty) ...[
          _SectionHeader('${l10n.searchUsers} · ${state.userResults.length}'),
          ...state.userResults.map((r) => _SearchResultRow(
                result: r,
                query: state.query,
                // Open the user profile (Add friend / Message live there) rather
                // than silently creating a DIRECT chat on tap.
                onTap: () => context.pushNamed(
                  AppRoutes.chatProfile,
                  pathParameters: {'id': r.id},
                ),
              )),
        ],
        if (state.chatResults.isNotEmpty) ...[
          _SectionHeader('${l10n.searchChats} · ${state.chatResults.length}'),
          ...state.chatResults.map((r) => _SearchResultRow(
                result: r,
                query: state.query,
                onTap: () => context.pushNamed(
                  AppRoutes.chatDetail,
                  pathParameters: {'id': r.id},
                  extra: {
                    'name': r.title,
                    'interlocutorId': '',
                    // Use the real chat type so channels open as CHANNEL (not GROUP).
                    'chatType': r.chatType ?? 'GROUP',
                  },
                ),
              )),
        ],
        if (state.messageResults.isNotEmpty) ...[
          _SectionHeader(
              '${l10n.searchMessages} · ${state.messageResults.length}'),
          ...state.messageResults.take(5).map((r) => _SearchResultRow(
                result: r,
                query: state.query,
                isMessage: true,
                onTap: () {
                  final chatId = r.chatId ?? r.id;
                  context.pushNamed(
                    AppRoutes.chatDetail,
                    pathParameters: {'id': chatId},
                    extra: {
                      'name': r.title,
                      'interlocutorId': '',
                      'chatType': 'DIRECT',
                    },
                  );
                },
              )),
        ],
        const SizedBox(height: 24),
      ],
    );
  }
}

// ─── Section header ───────────────────────────────────────────

class _SectionHeader extends StatelessWidget {
  final String text;
  const _SectionHeader(this.text);

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    return Padding(
      padding: const EdgeInsets.fromLTRB(14, 14, 14, 6),
      child: Text(
        text,
        style: TextStyle(
          fontFamily: 'Manrope',
          fontSize: 12,
          fontWeight: FontWeight.w700,
          letterSpacing: 0.1 * 12,
          color: palette.textSecondary,
        ),
      ),
    );
  }
}

// ─── Search result row ────────────────────────────────────────

class _SearchResultRow extends StatelessWidget {
  final SearchResult result;
  final String query;
  final bool isMessage;
  final VoidCallback onTap;

  const _SearchResultRow({
    required this.result,
    required this.query,
    required this.onTap,
    this.isMessage = false,
  });

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);

    final initials =
        result.title.isNotEmpty ? result.title[0].toUpperCase() : '?';

    return InkWell(
      onTap: onTap,
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
        child: Row(
          children: [
            isMessage
                ? Container(
                    width: 40,
                    height: 40,
                    decoration: BoxDecoration(
                      color: theme.colorScheme.surface,
                      shape: BoxShape.circle,
                    ),
                    child: Icon(Icons.message_outlined,
                        size: 18, color: palette.textAccent),
                  )
                : SudaAvatar(
                    mediaUrl: result.imageUrl.isNotEmpty ? result.imageUrl : null,
                    initials: initials,
                    size: 40,
                  ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    result.title,
                    style: TextStyle(
                      fontFamily: 'Manrope',
                      fontSize: 15,
                      fontWeight: FontWeight.w600,
                      color: theme.colorScheme.onSurface,
                    ),
                    overflow: TextOverflow.ellipsis,
                  ),
                  if (result.description.isNotEmpty)
                    HighlightText(
                      text: result.description,
                      query: query,
                      palette: palette,
                    ),
                ],
              ),
            ),
            Icon(
              Icons.arrow_outward_rounded,
              size: 16,
              color: palette.textSecondary,
            ),
          ],
        ),
      ),
    );
  }
}

// ─── Empty state ──────────────────────────────────────────────

class _EmptyState extends StatelessWidget {
  final String query;
  const _EmptyState({required this.query});

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    return Center(
      child: Text(
        context.l10n.searchNoResults(query),
        style: TextStyle(
          fontFamily: 'Manrope',
          color: palette.textSecondary,
          fontSize: 14,
        ),
        textAlign: TextAlign.center,
      ),
    );
  }
}
