import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import '../../../../app/config/theme/app_theme.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/placeholder/placeholder_page.dart';
import '../../../../shared/presentation/widgets/suda_avatar.dart';
import '../bloc/contacts_bloc.dart';

class BlockedUsersPage extends StatelessWidget {
  const BlockedUsersPage({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) => sl<ContactsBloc>()..add(LoadBlockedUsers()),
      child: const _BlockedUsersView(),
    );
  }
}

class _BlockedUsersView extends StatelessWidget {
  const _BlockedUsersView();

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
        title: Text(
          l10n.blockedUsersTitle,
          style: TextStyle(
            fontFamily: 'Manrope',
            fontSize: 18,
            fontWeight: FontWeight.w700,
            color: theme.colorScheme.onSurface,
          ),
        ),
        bottom: PreferredSize(
          preferredSize: const Size.fromHeight(1),
          child: Divider(height: 1, color: palette.divider),
        ),
      ),
      body: BlocBuilder<ContactsBloc, ContactsState>(
        buildWhen: (prev, next) =>
            prev.isLoadingBlocked != next.isLoadingBlocked ||
            prev.blockedUsers != next.blockedUsers,
        builder: (context, state) {
          if (state.isLoadingBlocked) {
            return const Center(child: CircularProgressIndicator());
          }
          if (state.blockedUsers.isEmpty) {
            return const PlaceholderPage(type: PlaceholderType.noContent, scaffold: false);
          }
          return ListView.builder(
            itemCount: state.blockedUsers.length,
            itemBuilder: (context, index) {
              final user = state.blockedUsers[index];
              return _BlockedUserRow(
                userId: user.id,
                displayName: user.displayName,
                username: user.username,
                avatarMediaId: user.avatarMediaId,
                onUnblock: () => context
                    .read<ContactsBloc>()
                    .add(UnblockFromList(user.id)),
              );
            },
          );
        },
      ),
    );
  }
}

class _BlockedUserRow extends StatelessWidget {
  final String userId;
  final String displayName;
  final String username;
  final String? avatarMediaId;
  final VoidCallback onUnblock;

  const _BlockedUserRow({
    required this.userId,
    required this.displayName,
    required this.username,
    this.avatarMediaId,
    required this.onUnblock,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    final initials = displayName.isNotEmpty ? displayName[0].toUpperCase() : '?';

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
      decoration: BoxDecoration(
        border: Border(
          bottom: BorderSide(color: palette.divider, width: 0.5),
        ),
      ),
      child: Row(
        children: [
          SudaAvatar(
            mediaUrl: avatarMediaId,
            initials: initials,
            size: 48,
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  displayName.isNotEmpty
                      ? displayName
                      : l10n.blockedUserFallbackName,
                  style: TextStyle(
                    fontFamily: 'Manrope',
                    fontSize: 15,
                    fontWeight: FontWeight.w600,
                    color: theme.colorScheme.onSurface,
                  ),
                ),
                if (username.isNotEmpty)
                  Text(
                    username,
                    style: TextStyle(
                      fontFamily: 'Manrope',
                      fontSize: 13,
                      color: palette.textSecondary,
                    ),
                  ),
              ],
            ),
          ),
          TextButton(
            onPressed: onUnblock,
            child: Text(
              l10n.blockedUsersUnblock,
              style: TextStyle(
                fontFamily: 'Manrope',
                fontWeight: FontWeight.w600,
                color: palette.danger,
              ),
            ),
          ),
        ],
      ),
    );
  }
}
