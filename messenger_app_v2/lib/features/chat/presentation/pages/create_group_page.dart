import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import '../../../../app/config/theme/app_theme.dart';
import '../../../../app/navigation/app_routes.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/widgets/suda_avatar.dart';
import '../../../../shared/presentation/widgets/suda_button.dart';
import '../../../../shared/presentation/widgets/suda_text_field.dart';
import '../../../friends/presentation/bloc/friends_cubit.dart';
import '../../../search/data/models/search_result.dart';
import '../../../search/presentation/bloc/search_bloc.dart';
import '../../domain/repositories/chat_repository.dart';

class CreateGroupPage extends StatefulWidget {
  const CreateGroupPage({super.key});

  @override
  State<CreateGroupPage> createState() => _CreateGroupPageState();
}

class _CreateGroupPageState extends State<CreateGroupPage> {
  final _groupNameCtrl = TextEditingController();
  final _searchCtrl = TextEditingController();
  final Set<SearchResult> _selected = {};
  bool _isCreating = false;
  int _step = 0; // 0 = select members, 1 = set name

  @override
  void dispose() {
    _groupNameCtrl.dispose();
    _searchCtrl.dispose();
    super.dispose();
  }

  void _toggleUser(SearchResult user) {
    setState(() {
      if (_selected.any((u) => u.id == user.id)) {
        _selected.removeWhere((u) => u.id == user.id);
      } else {
        _selected.add(user);
      }
    });
  }

  Future<void> _createGroup(BuildContext context) async {
    final name = _groupNameCtrl.text.trim();
    if (name.isEmpty || _selected.isEmpty) return;

    setState(() => _isCreating = true);
    final repo = sl<ChatRepository>();
    final result = await repo.createChat(
      type: 'GROUP',
      name: name,
      memberIds: _selected.map((u) => u.id).toList(),
    );
    if (!context.mounted) return;
    setState(() => _isCreating = false);

    result.fold(
      (failure) => ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(failure.message)),
      ),
      (chat) => context.pushReplacementNamed(
        AppRoutes.chatDetail,
        pathParameters: {'id': chat.id},
        extra: {'name': name, 'interlocutorId': '', 'chatType': 'GROUP'},
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);
    final l10n = context.l10n;

    return MultiBlocProvider(
      providers: [
        BlocProvider(create: (_) => sl<SearchBloc>()),
        BlocProvider(create: (_) => sl<FriendsCubit>()..load()),
      ],
      child: Scaffold(
        backgroundColor: theme.scaffoldBackgroundColor,
        appBar: AppBar(
          backgroundColor: theme.scaffoldBackgroundColor,
          elevation: 0,
          leading: IconButton(
            icon: Icon(Icons.arrow_back_rounded, color: theme.colorScheme.onSurface),
            onPressed: _step == 1 ? () => setState(() => _step = 0) : () => context.pop(),
          ),
          title: Text(
            _step == 0 ? l10n.selectMembers : l10n.groupName,
            style: TextStyle(fontFamily: 'Manrope', fontWeight: FontWeight.w700, color: theme.colorScheme.onSurface),
          ),
          actions: [
            if (_step == 0 && _selected.isNotEmpty)
              TextButton(
                onPressed: () => setState(() => _step = 1),
                child: Text(
                  'Next',
                  style: TextStyle(fontFamily: 'Manrope', fontWeight: FontWeight.w700, color: theme.colorScheme.tertiary),
                ),
              ),
          ],
        ),
        body: _step == 0
            ? _buildSelectStep(context, palette, theme, l10n)
            : _buildNameStep(context, palette, theme, l10n),
      ),
    );
  }

  Widget _buildSelectStep(BuildContext context, AppPaletteTheme palette, ThemeData theme, dynamic l10n) {
    return Column(
      children: [
        // Selected users chips row
        if (_selected.isNotEmpty)
          Container(
            height: 56,
            decoration: BoxDecoration(
              color: theme.colorScheme.surface,
              border: Border(bottom: BorderSide(color: palette.divider)),
            ),
            child: ListView(
              scrollDirection: Axis.horizontal,
              padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
              children: _selected.map((u) {
                return Container(
                  margin: const EdgeInsets.only(right: 8),
                  padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                  decoration: BoxDecoration(
                    color: theme.colorScheme.primary.withValues(alpha: 0.15),
                    borderRadius: BorderRadius.circular(100),
                    border: Border.all(color: theme.colorScheme.primary.withValues(alpha: 0.4)),
                  ),
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Text(u.title, style: TextStyle(fontFamily: 'Manrope', fontSize: 13, color: theme.colorScheme.onSurface)),
                      const SizedBox(width: 6),
                      GestureDetector(
                        onTap: () => _toggleUser(u),
                        child: Icon(Icons.close_rounded, size: 14, color: palette.textSecondary),
                      ),
                    ],
                  ),
                );
              }).toList(),
            ),
          ),

        // Search field
        Padding(
          padding: const EdgeInsets.all(14),
          child: BlocBuilder<SearchBloc, SearchState>(
            builder: (ctx, _) => TextField(
              controller: _searchCtrl,
              onChanged: (q) => ctx.read<SearchBloc>().add(SearchQueryChanged(q)),
              decoration: InputDecoration(
                hintText: l10n.searchHint,
                hintStyle: TextStyle(fontFamily: 'Manrope', color: palette.textSecondary),
                prefixIcon: Icon(Icons.search_rounded, color: palette.textSecondary),
                filled: true,
                fillColor: theme.colorScheme.surface,
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(14),
                  borderSide: BorderSide(color: palette.divider),
                ),
                enabledBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(14),
                  borderSide: BorderSide(color: palette.divider),
                ),
                focusedBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(14),
                  borderSide: BorderSide(color: theme.colorScheme.tertiary),
                ),
                contentPadding: EdgeInsets.zero,
              ),
            ),
          ),
        ),

        // Results: contacts when no query, search results otherwise
        Expanded(
          child: BlocBuilder<SearchBloc, SearchState>(
            builder: (_, state) {
              final hasQuery = state.query.trim().isNotEmpty;

              if (!hasQuery) {
                // No query — show server contacts as the default pool
                return _buildContactsList(palette, theme, l10n);
              }
              if (state.status == SearchStatus.loading) {
                return const Center(child: CircularProgressIndicator());
              }
              if (state.status == SearchStatus.success && state.userResults.isNotEmpty) {
                final users = state.userResults;
                return ListView.builder(
                  itemCount: users.length,
                  itemBuilder: (_, i) =>
                      _selectableRow(users[i], palette, theme),
                );
              }
              return const SizedBox.shrink();
            },
          ),
        ),
      ],
    );
  }

  /// Server contacts shown before the user types a search query.
  Widget _buildContactsList(AppPaletteTheme palette, ThemeData theme, dynamic l10n) {
    return BlocBuilder<FriendsCubit, FriendsState>(
      buildWhen: (p, c) => p.friends != c.friends || p.status != c.status,
      builder: (_, state) {
        if (state.status == FriendsStatus.loading ||
            state.status == FriendsStatus.initial) {
          return const Center(child: CircularProgressIndicator());
        }
        if (state.friends.isEmpty) return const SizedBox.shrink();
        return ListView.builder(
          itemCount: state.friends.length + 1,
          itemBuilder: (_, i) {
            if (i == 0) {
              return Padding(
                padding: const EdgeInsets.fromLTRB(16, 12, 16, 6),
                child: Text(
                  l10n.navFriends.toUpperCase(),
                  style: TextStyle(
                    fontFamily: 'Manrope',
                    fontSize: 11,
                    fontWeight: FontWeight.w700,
                    letterSpacing: 0.1 * 11,
                    color: palette.textSecondary,
                  ),
                ),
              );
            }
            final friend = state.friends[i - 1];
            final result = SearchResult(
              type: SearchResultType.user,
              id: friend.userId,
              title: friend.displayName,
              description: friend.username.isNotEmpty
                  ? (friend.username.startsWith('@')
                      ? friend.username
                      : '@${friend.username}')
                  : '',
            );
            return _selectableRow(result, palette, theme);
          },
        );
      },
    );
  }

  Widget _selectableRow(SearchResult user, AppPaletteTheme palette, ThemeData theme) {
    final isSelected = _selected.any((u) => u.id == user.id);
    return ListTile(
      leading: SudaAvatar(initials: user.title, size: 40),
      title: Text(user.title,
          style: TextStyle(fontFamily: 'Manrope', fontWeight: FontWeight.w600, color: theme.colorScheme.onSurface)),
      subtitle: user.description.isEmpty
          ? null
          : Text(user.description,
              style: TextStyle(fontFamily: 'Manrope', fontSize: 13, color: palette.textSecondary)),
      trailing: isSelected
          ? Icon(Icons.check_circle_rounded, color: theme.colorScheme.tertiary)
          : Icon(Icons.circle_outlined, color: palette.textSecondary),
      onTap: () => _toggleUser(user),
    );
  }

  Widget _buildNameStep(BuildContext context, AppPaletteTheme palette, ThemeData theme, dynamic l10n) {
    return Padding(
      padding: const EdgeInsets.all(20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Selected members preview
          Wrap(
            spacing: 8,
            runSpacing: 8,
            children: _selected.map((u) => Chip(
              avatar: SudaAvatar(initials: u.title, size: 20),
              label: Text(u.title, style: const TextStyle(fontFamily: 'Manrope', fontSize: 13)),
              backgroundColor: theme.colorScheme.surface,
              side: BorderSide(color: palette.divider),
            )).toList(),
          ),

          const SizedBox(height: 28),

          SudaTextField(
            label: l10n.groupName,
            hint: l10n.groupNameHint,
            controller: _groupNameCtrl,
          ),

          const SizedBox(height: 28),

          ListenableBuilder(
            listenable: _groupNameCtrl,
            builder: (_, __) => SudaButton(
              label: l10n.createGroup,
              variant: SudaButtonVariant.primary,
              size: SudaButtonSize.lg,
              fullWidth: true,
              loading: _isCreating,
              onPressed: (_groupNameCtrl.text.trim().isNotEmpty && !_isCreating)
                  ? () => _createGroup(context)
                  : null,
            ),
          ),
        ],
      ),
    );
  }
}
