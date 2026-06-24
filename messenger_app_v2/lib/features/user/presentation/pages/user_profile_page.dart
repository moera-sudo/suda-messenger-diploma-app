import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import '../../../../app/config/theme/app_theme.dart';
import '../../../../app/navigation/app_routes.dart';
import '../../../../features/chat/domain/repositories/chat_repository.dart';
import '../../../../features/friends/domain/repositories/friends_repository.dart';
import '../../../../shared/presentation/format_utils.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/widgets/suda_avatar.dart';
import '../../data/models/user_models.dart';
import '../bloc/user_profile_bloc.dart';

class UserProfilePage extends StatelessWidget {
  final String userId;
  // When opened from a DIRECT chat, identifies that chat so the profile can
  // link to its shared media. Null when opened from search/friends/etc.
  final String? chatId;
  final String? chatName;

  const UserProfilePage({
    super.key,
    required this.userId,
    this.chatId,
    this.chatName,
  });

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) => sl<UserProfileBloc>()..add(LoadUserProfile(userId)),
      child: _UserProfileView(chatId: chatId, chatName: chatName),
    );
  }
}

class _UserProfileView extends StatelessWidget {
  final String? chatId;
  final String? chatName;
  const _UserProfileView({this.chatId, this.chatName});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Theme.of(context).scaffoldBackgroundColor,
      body: BlocBuilder<UserProfileBloc, UserProfileState>(
        buildWhen: (prev, next) =>
            prev.status != next.status || prev.isBlocked != next.isBlocked,
        builder: (context, state) {
          if (state.status == UserProfileStatus.loading ||
              state.status == UserProfileStatus.initial) {
            return const Center(child: CircularProgressIndicator());
          }
          if (state.status == UserProfileStatus.failure || state.profile == null) {
            return _ErrorView(error: state.error);
          }
          return _ProfileContent(
            profile: state.profile!,
            isBlocked: state.isBlocked,
            chatId: chatId,
            chatName: chatName,
          );
        },
      ),
    );
  }
}

class _ErrorView extends StatelessWidget {
  final String? error;
  const _ErrorView({this.error});

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    return SafeArea(
      child: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(8),
            child: Align(
              alignment: Alignment.centerLeft,
              child: IconButton(
                icon: Icon(Icons.arrow_back_rounded, color: Theme.of(context).colorScheme.onSurface),
                onPressed: () => context.pop(),
              ),
            ),
          ),
          Expanded(
            child: Center(
              child: Text(
                error ?? context.l10n.errorGeneric,
                style: TextStyle(color: palette.danger),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _ProfileContent extends StatelessWidget {
  final UserProfileModel profile;
  final bool isBlocked;
  final String? chatId;
  final String? chatName;

  const _ProfileContent({
    required this.profile,
    required this.isBlocked,
    this.chatId,
    this.chatName,
  });

  @override
  Widget build(BuildContext context) {
    return CustomScrollView(
      slivers: [
        SliverToBoxAdapter(child: _UserHero(profile: profile)),
        SliverToBoxAdapter(child: _UserActions(profile: profile)),
        SliverToBoxAdapter(
          child: _UserInfoSection(
            profile: profile,
            isBlocked: isBlocked,
            chatId: chatId,
            chatName: chatName,
          ),
        ),
        const SliverToBoxAdapter(child: SizedBox(height: 32)),
      ],
    );
  }
}

// ─── Hero ─────────────────────────────────────────────────────

class _UserHero extends StatelessWidget {
  final UserProfileModel profile;
  const _UserHero({required this.profile});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final l10n = context.l10n;

    final statusText = profile.isOnline
        ? l10n.userProfileOnline
        : _formatLastSeen(profile.lastSeenAt, l10n);

    final initials = profile.displayName.isNotEmpty
        ? profile.displayName[0].toUpperCase()
        : '?';

    return Container(
      constraints: const BoxConstraints(minHeight: 220),
      child: Stack(
        children: [
          // Gradient background
          Positioned.fill(
            child: DecoratedBox(
              decoration: BoxDecoration(
                gradient: RadialGradient(
                  center: const Alignment(0.2, -0.2),
                  radius: 1.2,
                  colors: [
                    theme.colorScheme.primary.withValues(alpha: 0.5),
                    theme.colorScheme.surface,
                  ],
                ),
              ),
            ),
          ),
          // Back button
          Positioned(
            top: 44,
            left: 4,
            child: SafeArea(
              child: IconButton(
                icon: Icon(Icons.arrow_back_rounded, color: theme.colorScheme.onSurface),
                onPressed: () => context.pop(),
              ),
            ),
          ),
          // More button
          Positioned(
            top: 44,
            right: 4,
            child: SafeArea(
              child: IconButton(
                icon: Icon(Icons.more_vert_rounded, color: theme.colorScheme.onSurface),
                onPressed: () => _showMoreSheet(context),
              ),
            ),
          ),
          // Content
          Align(
            alignment: Alignment.bottomCenter,
            child: Padding(
              padding: const EdgeInsets.fromLTRB(18, 80, 18, 20),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  SudaAvatar(
                    mediaId: profile.avatarMediaId,
                    initials: initials,
                    size: 96,
                    ring: true,
                    showOnline: profile.isOnline,
                  ),
                  const SizedBox(height: 10),
                  Text(
                    profile.displayName,
                    style: const TextStyle(
                      fontFamily: 'Manrope',
                      fontSize: 22,
                      fontWeight: FontWeight.w800,
                      letterSpacing: -0.03 * 22,
                    ),
                    textAlign: TextAlign.center,
                  ),
                  const SizedBox(height: 4),
                  Text(
                    '${formatHandle(profile.username)} · $statusText',
                    style: TextStyle(
                      fontFamily: 'Manrope',
                      fontSize: 13,
                      color: Theme.of(context).extension<AppPaletteTheme>()!.textSecondary,
                    ),
                    textAlign: TextAlign.center,
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  String _formatLastSeen(String? lastSeenAt, dynamic l10n) {
    if (lastSeenAt == null) return l10n.chatStatusOffline;
    try {
      final dt = DateTime.parse(lastSeenAt).toLocal();
      final now = DateTime.now();
      final diff = now.difference(dt);
      // Guard against unrealistic values (clock skew → future, epoch/0 → huge).
      if (diff.isNegative) return l10n.chatLastSeenJustNow;
      if (diff.inDays > 3650) return l10n.chatStatusOffline;
      if (diff.inMinutes < 1) return l10n.chatLastSeenJustNow;
      if (diff.inMinutes < 60) return '${l10n.chatLastSeenPrefix} ${diff.inMinutes} ${l10n.chatLastSeenMinutes}';
      if (diff.inHours < 24) return '${l10n.chatLastSeenPrefix} ${diff.inHours} ${l10n.chatLastSeenHours}';
      final day = dt.day.toString().padLeft(2, '0');
      final month = dt.month.toString().padLeft(2, '0');
      return '${l10n.chatLastSeenPrefix} $day.$month';
    } catch (_) {
      return l10n.chatStatusOffline;
    }
  }

  void _showMoreSheet(BuildContext context) {
    // More options sheet — to be extended in future iterations
  }
}

// ─── Actions ──────────────────────────────────────────────────

class _UserActions extends StatelessWidget {
  final UserProfileModel profile;
  const _UserActions({required this.profile});

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 12),
      decoration: BoxDecoration(
        color: Theme.of(context).scaffoldBackgroundColor,
        border: Border(bottom: BorderSide(color: palette.divider)),
      ),
      child: Row(
        children: [
          _ActionButton(
            icon: Icons.chat_rounded,
            label: l10n.userProfileMessage,
            onTap: () => _onMessage(context),
          ),
          const SizedBox(width: 8),
          _ActionButton(
            icon: Icons.notifications_none_rounded,
            label: l10n.userProfileMute,
            onTap: () => _onComingSoon(context),
          ),
        ],
      ),
    );
  }

  void _onMessage(BuildContext context) {
    sl<ChatRepository>()
        .createChat(type: 'DIRECT', targetId: profile.id)
        .then((result) {
      result.fold(
        (_) {},
        (chat) {
          if (context.mounted) {
            context.pushNamed(
              AppRoutes.chatDetail,
              pathParameters: {'id': chat.id},
              extra: {
                'name': profile.displayName,
                'interlocutorId': profile.id,
                'chatType': 'DIRECT',
              },
            );
          }
        },
      );
    });
  }

  void _onComingSoon(BuildContext context) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(context.l10n.comingSoon),
        duration: const Duration(seconds: 2),
      ),
    );
  }
}

class _ActionButton extends StatelessWidget {
  final IconData icon;
  final String label;
  final VoidCallback onTap;

  const _ActionButton({
    required this.icon,
    required this.label,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;

    return Expanded(
      child: Material(
        color: theme.colorScheme.surface,
        borderRadius: BorderRadius.circular(14),
        child: InkWell(
          onTap: onTap,
          borderRadius: BorderRadius.circular(14),
          child: Padding(
            padding: const EdgeInsets.symmetric(vertical: 10, horizontal: 6),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Icon(icon, size: 18, color: theme.colorScheme.tertiary),
                const SizedBox(height: 6),
                Text(
                  label,
                  style: TextStyle(
                    fontFamily: 'Manrope',
                    fontSize: 12,
                    fontWeight: FontWeight.w600,
                    color: palette.textSecondary,
                  ),
                  textAlign: TextAlign.center,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

// ─── Info section ─────────────────────────────────────────────

class _UserInfoSection extends StatelessWidget {
  final UserProfileModel profile;
  final bool isBlocked;
  final String? chatId;
  final String? chatName;

  const _UserInfoSection({
    required this.profile,
    required this.isBlocked,
    this.chatId,
    this.chatName,
  });

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 8),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // BIO card
          if (profile.bio != null && profile.bio!.isNotEmpty) ...[
            _InfoCard(
              label: l10n.userProfileBioSection,
              child: Text(
                profile.bio!,
                style: const TextStyle(
                  fontFamily: 'Manrope',
                  fontSize: 14,
                  height: 1.6,
                ),
              ),
            ),
            const SizedBox(height: 12),
          ],

          // Wallet address card
          if (profile.walletAddress != null &&
              profile.walletAddress!.isNotEmpty) ...[
            _InfoCard(
              label: l10n.userProfileWalletAddress,
              child: _WalletAddressRow(address: profile.walletAddress!),
            ),
            const SizedBox(height: 12),
          ],

          // Action list
          _InfoListCard(
            children: [
              // Shared media — only when opened from a chat (needs chatId).
              if (chatId != null)
                _InfoListRow(
                  icon: Icons.image_outlined,
                  label: l10n.userProfileSharedMedia,
                  onTap: () => context.pushNamed(
                    AppRoutes.sharedMedia,
                    pathParameters: {'id': chatId!},
                    extra: {'chatName': chatName ?? profile.displayName},
                  ),
                ),
              _InfoListRow(
                icon: Icons.notifications_none_rounded,
                label: l10n.userProfileNotifications,
                trailing: l10n.chatStatusOnline,
              ),
              _FriendStatusRow(userId: profile.id),
              _InfoListRow(
                icon: Icons.block_rounded,
                label: isBlocked
                    ? l10n.userProfileUnblockUser
                    : l10n.userProfileBlockUser,
                isDanger: true,
                onTap: () => _onBlockToggle(context, isBlocked),
                showChevron: false,
              ),
            ],
          ),
        ],
      ),
    );
  }

  void _onBlockToggle(BuildContext context, bool currentlyBlocked) {
    final l10n = context.l10n;
    if (currentlyBlocked) {
      context.read<UserProfileBloc>().add(UnblockUserEvent());
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(l10n.userProfileUnblockSuccess)),
      );
    } else {
      showDialog<bool>(
        context: context,
        builder: (ctx) => AlertDialog(
          title: Text(l10n.userProfileBlockConfirm(profile.displayName)),
          content: Text(l10n.userProfileBlockConfirmMsg),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(ctx).pop(false),
              child: Text(l10n.buttonRetry),
            ),
            TextButton(
              onPressed: () => Navigator.of(ctx).pop(true),
              child: Text(
                l10n.userProfileBlockUser,
                style: TextStyle(
                  color: Theme.of(context).extension<AppPaletteTheme>()!.danger,
                ),
              ),
            ),
          ],
        ),
      ).then((confirmed) {
        if (confirmed == true && context.mounted) {
          context.read<UserProfileBloc>().add(BlockUserEvent());
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text(l10n.userProfileBlockSuccess)),
          );
        }
      });
    }
  }
}

class _InfoCard extends StatelessWidget {
  final String label;
  final Widget child;

  const _InfoCard({required this.label, required this.child});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;

    return Container(
      padding: const EdgeInsets.all(14),
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        border: Border.all(color: palette.divider),
        borderRadius: BorderRadius.circular(16),
      ),
      width: double.infinity,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            label,
            style: TextStyle(
              fontFamily: 'Manrope',
              fontSize: 11,
              fontWeight: FontWeight.w700,
              letterSpacing: 0.1 * 11,
              color: palette.textSecondary,
            ),
          ),
          const SizedBox(height: 8),
          child,
        ],
      ),
    );
  }
}

class _InfoListCard extends StatelessWidget {
  final List<Widget> children;
  const _InfoListCard({required this.children});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;

    return Container(
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        border: Border.all(color: palette.divider),
        borderRadius: BorderRadius.circular(16),
      ),
      clipBehavior: Clip.hardEdge,
      child: Column(children: children),
    );
  }
}

class _InfoListRow extends StatelessWidget {
  final IconData icon;
  final String label;
  final String? trailing;
  final bool isDanger;
  final bool showChevron;
  final VoidCallback? onTap;

  const _InfoListRow({
    required this.icon,
    required this.label,
    this.trailing,
    this.isDanger = false,
    this.showChevron = true,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final iconColor = isDanger ? palette.danger : palette.textAccent;
    final textColor = isDanger ? palette.danger : theme.colorScheme.onSurface;

    return InkWell(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 14),
        decoration: BoxDecoration(
          border: Border(bottom: BorderSide(color: palette.divider, width: 0.5)),
        ),
        child: Row(
          children: [
            Icon(icon, size: 20, color: iconColor),
            const SizedBox(width: 12),
            Expanded(
              child: Text(
                label,
                style: TextStyle(
                  fontFamily: 'Manrope',
                  fontSize: 14,
                  fontWeight: FontWeight.w500,
                  color: textColor,
                ),
              ),
            ),
            if (trailing != null) ...[
              Text(
                trailing!,
                style: TextStyle(
                  fontFamily: 'Manrope',
                  fontSize: 13,
                  color: palette.textSecondary,
                ),
              ),
              const SizedBox(width: 4),
            ],
            if (showChevron)
              Icon(Icons.chevron_right_rounded, size: 18, color: palette.textSecondary),
          ],
        ),
      ),
    );
  }
}

class _WalletAddressRow extends StatelessWidget {
  final String address;
  const _WalletAddressRow({required this.address});

  String get _truncated {
    if (address.length <= 12) return address;
    return '${address.substring(0, 6)}…${address.substring(address.length - 4)}';
  }

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;

    return Row(
      children: [
        Expanded(
          child: Text(
            _truncated,
            style: TextStyle(
              fontFamily: 'JetBrainsMono',
              fontSize: 13,
              color: palette.textSecondary,
            ),
          ),
        ),
        IconButton(
          icon: Icon(Icons.copy_rounded, size: 18, color: palette.textSecondary),
          onPressed: () {
            Clipboard.setData(ClipboardData(text: address));
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(
                content: Text('Copied'),
                duration: Duration(seconds: 1),
              ),
            );
          },
          padding: const EdgeInsets.all(4),
          constraints: const BoxConstraints(),
        ),
      ],
    );
  }
}

/// Friend status row — loads relation via FutureBuilder and renders the correct action.
class _FriendStatusRow extends StatefulWidget {
  final String userId;
  const _FriendStatusRow({required this.userId});

  @override
  State<_FriendStatusRow> createState() => _FriendStatusRowState();
}

class _FriendStatusRowState extends State<_FriendStatusRow> {
  late Future<String?> _statusFuture;

  @override
  void initState() {
    super.initState();
    _reload();
  }

  void _reload() {
    _statusFuture = sl<FriendsRepository>()
        .getStatus(widget.userId)
        .then((r) => r.fold((_) => null, (s) => s.relation));
  }

  Future<void> _act(String relation, String? requestId) async {
    final repo = sl<FriendsRepository>();
    switch (relation) {
      case 'NONE':
        await repo.sendRequest(widget.userId);
      case 'PENDING_SENT':
        if (requestId != null) await repo.cancelRequest(requestId);
      case 'PENDING_RECEIVED':
        if (requestId != null) await repo.acceptRequest(requestId);
      case 'FRIENDS':
        await repo.unfriend(widget.userId);
    }
    if (mounted) setState(() => _reload());
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    return FutureBuilder<String?>(
      future: _statusFuture,
      builder: (_, snap) {
        final relation = snap.data;
        if (relation == null || relation == 'BLOCKED') return const SizedBox.shrink();

        final (IconData icon, String label, bool isDanger) = switch (relation) {
          'NONE' => (Icons.person_add_rounded, l10n.friendsAddFriend, false),
          'PENDING_SENT' => (Icons.schedule_rounded, l10n.friendsCancelRequest, false),
          'PENDING_RECEIVED' => (Icons.check_circle_outline_rounded, l10n.friendsAccept, false),
          'FRIENDS' => (Icons.people_rounded, l10n.friendsYouAreFriends, false),
          _ => (Icons.person_add_rounded, l10n.friendsAddFriend, false),
        };

        return InkWell(
          onTap: () async {
            final statusResult = await sl<FriendsRepository>().getStatus(widget.userId);
            final requestId = statusResult.fold((_) => null, (s) => s.requestId);
            if (mounted) _act(relation, requestId);
          },
          child: Container(
            padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 14),
            decoration: BoxDecoration(
              border: Border(bottom: BorderSide(color: palette.divider, width: 0.5)),
            ),
            child: Row(children: [
              Icon(icon, size: 20, color: isDanger ? palette.danger : palette.textAccent),
              const SizedBox(width: 12),
              Expanded(
                child: Text(
                  label,
                  style: TextStyle(
                    fontFamily: 'Manrope',
                    fontSize: 14,
                    fontWeight: FontWeight.w500,
                    color: isDanger ? palette.danger : theme.colorScheme.onSurface,
                  ),
                ),
              ),
            ]),
          ),
        );
      },
    );
  }
}
