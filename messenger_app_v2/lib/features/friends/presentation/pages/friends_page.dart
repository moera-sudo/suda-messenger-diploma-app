import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import '../../../../app/config/theme/app_theme.dart';
import '../../../../app/navigation/app_routes.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/placeholder/placeholder_page.dart';
import '../../../../shared/presentation/widgets/suda_avatar.dart';
import '../bloc/friends_cubit.dart';

class FriendsPage extends StatelessWidget {
  const FriendsPage({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) => sl<FriendsCubit>()..load(),
      child: const _FriendsView(),
    );
  }
}

class _FriendsView extends StatelessWidget {
  const _FriendsView();

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    return DefaultTabController(
      length: 2,
      child: Scaffold(
        backgroundColor: theme.scaffoldBackgroundColor,
        appBar: AppBar(
          backgroundColor: theme.scaffoldBackgroundColor,
          elevation: 0,
          scrolledUnderElevation: 0,
          automaticallyImplyLeading: false,
          title: Text(
            l10n.friendsTabFriends,
            style: TextStyle(
              fontFamily: 'Manrope',
              fontSize: 18,
              fontWeight: FontWeight.w700,
              color: theme.colorScheme.onSurface,
            ),
          ),
          actions: [
            IconButton(
              icon: Icon(Icons.person_search_rounded, color: palette.textAccent),
              tooltip: l10n.friendsAddFriend,
              onPressed: () => context.push('/search'),
            ),
          ],
          bottom: TabBar(
            labelStyle: const TextStyle(fontFamily: 'Manrope', fontWeight: FontWeight.w700),
            unselectedLabelStyle: const TextStyle(fontFamily: 'Manrope'),
            labelColor: theme.colorScheme.tertiary,
            unselectedLabelColor: palette.textSecondary,
            indicatorColor: theme.colorScheme.tertiary,
            dividerColor: palette.divider,
            tabs: [
              Tab(text: l10n.friendsTabFriends),
              Tab(text: l10n.friendsTabRequests),
            ],
          ),
        ),
        body: const TabBarView(children: [
          _FriendsTab(),
          _RequestsTab(),
        ]),
      ),
    );
  }
}

// ── Friends tab ────────────────────────────────────────────────

class _FriendsTab extends StatelessWidget {
  const _FriendsTab();

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<FriendsCubit, FriendsState>(
      buildWhen: (p, c) => p.friends != c.friends || p.status != c.status,
      builder: (context, state) {
        if (state.status == FriendsStatus.loading) {
          return const Center(child: CircularProgressIndicator());
        }
        if (state.friends.isEmpty) {
          return PlaceholderPage(
            type: PlaceholderType.noContent,
            scaffold: false,
            message: context.l10n.friendsEmptyFriends,
          );
        }
        return ListView.builder(
          itemCount: state.friends.length,
          itemBuilder: (_, i) => _FriendRow(friend: state.friends[i]),
        );
      },
    );
  }
}

class _FriendRow extends StatelessWidget {
  final dynamic friend;
  const _FriendRow({required this.friend});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    final initials = (friend.displayName as String).isNotEmpty
        ? (friend.displayName as String)[0].toUpperCase()
        : '?';

    return InkWell(
      onTap: () => context.pushNamed(
        AppRoutes.chatProfile,
        pathParameters: {'id': friend.userId as String},
      ),
      onLongPress: () => _showUnfriendSheet(context, l10n, palette, theme),
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
        decoration: BoxDecoration(
          border: Border(bottom: BorderSide(color: palette.divider, width: 0.5)),
        ),
        child: Row(children: [
          SudaAvatar(
            mediaUrl: friend.avatarMediaId as String?,
            initials: initials,
            size: 48,
            showOnline: friend.isOnline as bool? ?? false,
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Text(
                friend.displayName as String,
                style: TextStyle(
                    fontFamily: 'Manrope',
                    fontSize: 15,
                    fontWeight: FontWeight.w600,
                    color: theme.colorScheme.onSurface),
              ),
              Text(
                friend.username,
                style: TextStyle(fontFamily: 'Manrope', fontSize: 13, color: palette.textSecondary),
              ),
            ]),
          ),
          IconButton(
            icon: Icon(Icons.chat_rounded, color: theme.colorScheme.tertiary),
            onPressed: () => context.pushNamed(
              AppRoutes.chatProfile,
              pathParameters: {'id': friend.userId as String},
            ),
          ),
        ]),
      ),
    );
  }

  void _showUnfriendSheet(
      BuildContext context, dynamic l10n, AppPaletteTheme palette, ThemeData theme) {
    showModalBottomSheet(
      context: context,
      backgroundColor: theme.colorScheme.surface,
      shape: const RoundedRectangleBorder(
          borderRadius: BorderRadius.vertical(top: Radius.circular(24))),
      builder: (ctx) => SafeArea(
        child: Column(mainAxisSize: MainAxisSize.min, children: [
          Container(
            margin: const EdgeInsets.symmetric(vertical: 8),
            width: 38,
            height: 4,
            decoration: BoxDecoration(
                color: palette.textSecondary.withValues(alpha: 0.4),
                borderRadius: BorderRadius.circular(2)),
          ),
          ListTile(
            leading: Icon(Icons.person_remove_outlined, color: palette.danger),
            title: Text(l10n.friendsUnfriend,
                style: TextStyle(
                    color: palette.danger,
                    fontFamily: 'Manrope',
                    fontWeight: FontWeight.w600)),
            onTap: () {
              Navigator.pop(ctx);
              context.read<FriendsCubit>().unfriend(friend.userId as String);
            },
          ),
          const SizedBox(height: 8),
        ]),
      ),
    );
  }
}

// ── Requests tab ───────────────────────────────────────────────

class _RequestsTab extends StatelessWidget {
  const _RequestsTab();

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<FriendsCubit, FriendsState>(
      buildWhen: (p, c) => p.incoming != c.incoming || p.outgoing != c.outgoing,
      builder: (context, state) {
        if (state.incoming.isEmpty && state.outgoing.isEmpty) {
          return PlaceholderPage(
            type: PlaceholderType.noContent,
            scaffold: false,
            message: context.l10n.friendsEmptyRequests,
          );
        }
        return ListView(children: [
          if (state.incoming.isNotEmpty) ...[
            _SectionLabel(context.l10n.friendsIncoming),
            ...state.incoming.map((r) => _RequestRow(request: r, isIncoming: true)),
          ],
          if (state.outgoing.isNotEmpty) ...[
            _SectionLabel(context.l10n.friendsOutgoing),
            ...state.outgoing.map((r) => _RequestRow(request: r, isIncoming: false)),
          ],
        ]);
      },
    );
  }
}

class _SectionLabel extends StatelessWidget {
  final String label;
  const _SectionLabel(this.label);

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    return Padding(
      padding: const EdgeInsets.fromLTRB(14, 14, 14, 6),
      child: Text(
        label.toUpperCase(),
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
}

class _RequestRow extends StatelessWidget {
  final dynamic request;
  final bool isIncoming;
  const _RequestRow({required this.request, required this.isIncoming});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;
    final cubit = context.read<FriendsCubit>();

    final displayName = request.otherDisplayName as String;
    final initials = displayName.isNotEmpty ? displayName[0].toUpperCase() : '?';

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
      decoration: BoxDecoration(
          border: Border(bottom: BorderSide(color: palette.divider, width: 0.5))),
      child: Row(children: [
        SudaAvatar(
          mediaUrl: request.otherAvatarMediaId as String?,
          initials: initials,
          size: 48,
        ),
        const SizedBox(width: 12),
        Expanded(
          child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            Text(displayName,
                style: TextStyle(
                    fontFamily: 'Manrope',
                    fontSize: 15,
                    fontWeight: FontWeight.w600,
                    color: theme.colorScheme.onSurface)),
            Text(request.otherUsername,
                style: TextStyle(
                    fontFamily: 'Manrope', fontSize: 13, color: palette.textSecondary)),
          ]),
        ),
        if (isIncoming) ...[
          TextButton(
            onPressed: () => cubit.acceptRequest(request.requestId as String),
            child: Text(l10n.friendsAccept,
                style: TextStyle(fontFamily: 'Manrope', color: theme.colorScheme.tertiary)),
          ),
          TextButton(
            onPressed: () => cubit.rejectRequest(request.requestId as String),
            child: Text(l10n.friendsReject,
                style: TextStyle(fontFamily: 'Manrope', color: palette.danger)),
          ),
        ] else
          TextButton(
            onPressed: () => cubit.cancelRequest(request.requestId as String),
            child: Text(l10n.friendsCancelRequest,
                style: TextStyle(fontFamily: 'Manrope', color: palette.textSecondary)),
          ),
      ]),
    );
  }
}
