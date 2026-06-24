import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import '../../../../app/config/theme/app_theme.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/widgets/suda_avatar.dart';
import '../../data/models/search_result.dart';
import '../bloc/chat_search_bloc.dart';
import '../widgets/highlight_text.dart';

class ChatSearchPage extends StatefulWidget {
  final String chatId;
  final String chatName;

  const ChatSearchPage({
    super.key,
    required this.chatId,
    required this.chatName,
  });

  @override
  State<ChatSearchPage> createState() => _ChatSearchPageState();
}

class _ChatSearchPageState extends State<ChatSearchPage> {
  final _ctrl = TextEditingController();

  @override
  void dispose() {
    _ctrl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) {
        final bloc = sl<ChatSearchBloc>();
        bloc.chatId = widget.chatId;
        return bloc;
      },
      child: _ChatSearchView(ctrl: _ctrl, chatName: widget.chatName),
    );
  }
}

class _ChatSearchView extends StatelessWidget {
  final TextEditingController ctrl;
  final String chatName;

  const _ChatSearchView({required this.ctrl, required this.chatName});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    return Scaffold(
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
        title: _ChatSearchField(ctrl: ctrl, chatName: chatName),
        bottom: PreferredSize(
          preferredSize: const Size.fromHeight(1),
          child: Divider(height: 1, color: palette.divider),
        ),
      ),
      body: BlocBuilder<ChatSearchBloc, ChatSearchState>(
        builder: (context, state) {
          switch (state.status) {
            case ChatSearchStatus.loading:
              return LinearProgressIndicator(
                color: theme.colorScheme.tertiary,
                backgroundColor: theme.colorScheme.surface,
              );

            case ChatSearchStatus.failure:
              return Center(
                child: Text(
                  l10n.errorGeneric,
                  style: TextStyle(color: palette.danger),
                ),
              );

            case ChatSearchStatus.empty:
              return Center(
                child: Text(
                  l10n.searchNoResults(state.query),
                  style: TextStyle(
                    fontFamily: 'Manrope',
                    color: palette.textSecondary,
                  ),
                ),
              );

            case ChatSearchStatus.success:
              return ListView.builder(
                itemCount: state.results.length,
                itemBuilder: (ctx, index) => _MessageSearchRow(
                  result: state.results[index],
                  query: state.query,
                  palette: palette,
                  // Return the tapped message id so the chat can jump to it.
                  onTap: () => context.pop(state.results[index].messageId),
                ),
              );

            case ChatSearchStatus.initial:
              return Center(
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Icon(
                      Icons.manage_search_rounded,
                      size: 64,
                      color: theme.colorScheme.onSurface.withValues(alpha: 0.1),
                    ),
                    const SizedBox(height: 16),
                    Text(
                      l10n.chatSearchHint,
                      style: TextStyle(
                        fontFamily: 'Manrope',
                        color: theme.colorScheme.onSurface.withValues(alpha: 0.3),
                      ),
                    ),
                  ],
                ),
              );
          }
        },
      ),
    );
  }
}

// ─── Search field ─────────────────────────────────────────────

class _ChatSearchField extends StatelessWidget {
  final TextEditingController ctrl;
  final String chatName;

  const _ChatSearchField({required this.ctrl, required this.chatName});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

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
                  hintText: l10n.chatSearchHint,
                  hintStyle: TextStyle(color: palette.textSecondary),
                  border: InputBorder.none,
                  isDense: true,
                  contentPadding: EdgeInsets.zero,
                ),
                onChanged: (value) {
                  context
                      .read<ChatSearchBloc>()
                      .add(ChatSearchQueryChanged(value));
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
                  constraints: const BoxConstraints(minWidth: 36, minHeight: 36),
                  onPressed: () {
                    ctrl.clear();
                    context.read<ChatSearchBloc>().add(ClearChatSearch());
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

// ─── Message result row ───────────────────────────────────────

class _MessageSearchRow extends StatelessWidget {
  final SearchResult result;
  final String query;
  final AppPaletteTheme palette;
  final VoidCallback onTap;

  const _MessageSearchRow({
    required this.result,
    required this.query,
    required this.palette,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final initials =
        result.title.isNotEmpty ? result.title[0].toUpperCase() : '?';

    return InkWell(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
        decoration: BoxDecoration(
          border: Border(
            bottom: BorderSide(color: palette.divider, width: 0.5),
          ),
        ),
        child: Row(
          children: [
            SudaAvatar(initials: initials, size: 36),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    result.title,
                    style: TextStyle(
                      fontFamily: 'Manrope',
                      fontSize: 14,
                      fontWeight: FontWeight.w600,
                      color: theme.colorScheme.onSurface,
                    ),
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
          ],
        ),
      ),
    );
  }
}
