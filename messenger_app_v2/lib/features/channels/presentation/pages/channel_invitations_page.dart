import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import '../../../../app/config/theme/app_theme.dart';
import '../../../../shared/presentation/feedback/app_feedback.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/placeholder/placeholder_page.dart';
import '../../../../shared/presentation/widgets/suda_avatar.dart';
import '../../../../shared/presentation/widgets/suda_button.dart';
import '../../data/models/channel_view_model.dart';
import '../../domain/repositories/channel_repository.dart';

class ChannelInvitationsPage extends StatefulWidget {
  const ChannelInvitationsPage({super.key});

  @override
  State<ChannelInvitationsPage> createState() => _ChannelInvitationsPageState();
}

class _ChannelInvitationsPageState extends State<ChannelInvitationsPage> {
  List<ChannelInviteItem>? _invites;
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    if (mounted) setState(() => _isLoading = true);
    final result = await sl<ChannelRepository>().getMyInvites();
    if (!mounted) return;
    setState(() {
      _isLoading = false;
      _invites = result.fold((_) => [], (list) => list);
    });
  }

  Future<void> _accept(ChannelInviteItem item) async {
    setState(() => _invites?.removeWhere((i) => i.channelId == item.channelId));
    final result = await sl<ChannelRepository>().acceptInvite(item.channelId);
    if (!mounted) return;
    result.fold(
      (f) {
        AppFeedback.showError(f.message);
        _load();
      },
      (_) => AppFeedback.showSuccess(context.l10n.channelSubscribeSuccess),
    );
  }

  Future<void> _decline(ChannelInviteItem item) async {
    setState(() => _invites?.removeWhere((i) => i.channelId == item.channelId));
    final result = await sl<ChannelRepository>().declineInvite(item.channelId);
    if (!mounted) return;
    result.fold(
      (f) {
        AppFeedback.showError(f.message);
        _load();
      },
      (_) {},
    );
  }

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
          l10n.channelInvitationsTitle,
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
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : (_invites?.isEmpty ?? true)
              ? PlaceholderPage(
                  type: PlaceholderType.noContent,
                  scaffold: false,
                  message: l10n.channelInvitationsEmpty,
                )
              : ListView.builder(
                  itemCount: _invites!.length,
                  itemBuilder: (context, i) {
                    final item = _invites![i];
                    return _InviteRow(
                      item: item,
                      onAccept: () => _accept(item),
                      onDecline: () => _decline(item),
                    );
                  },
                ),
    );
  }
}

class _InviteRow extends StatelessWidget {
  final ChannelInviteItem item;
  final VoidCallback onAccept;
  final VoidCallback onDecline;

  const _InviteRow({
    required this.item,
    required this.onAccept,
    required this.onDecline,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    final handle = item.username?.isNotEmpty == true
        ? (item.username!.startsWith('@') ? item.username! : '@${item.username}')
        : null;

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 12),
      decoration: BoxDecoration(
        border: Border(bottom: BorderSide(color: palette.divider, width: 0.5)),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SudaAvatar(
            mediaId: item.avatarMediaId,
            initials: item.name.isNotEmpty ? item.name[0] : '?',
            size: 48,
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  item.name,
                  style: TextStyle(
                    fontFamily: 'Manrope',
                    fontSize: 15,
                    fontWeight: FontWeight.w600,
                    color: theme.colorScheme.onSurface,
                  ),
                ),
                if (handle != null)
                  Text(
                    handle,
                    style: TextStyle(
                      fontFamily: 'Manrope',
                      fontSize: 13,
                      color: palette.textSecondary,
                    ),
                  ),
                if (item.description != null && item.description!.isNotEmpty)
                  Text(
                    item.description!,
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                    style: TextStyle(
                      fontFamily: 'Manrope',
                      fontSize: 13,
                      color: palette.textSecondary,
                      height: 1.4,
                    ),
                  ),
                const SizedBox(height: 10),
                Row(
                  children: [
                    Expanded(
                      child: SudaButton(
                        label: l10n.channelAcceptInvite,
                        variant: SudaButtonVariant.primary,
                        size: SudaButtonSize.sm,
                        fullWidth: true,
                        onPressed: onAccept,
                      ),
                    ),
                    const SizedBox(width: 8),
                    Expanded(
                      child: SudaButton(
                        label: l10n.channelDeclineInvite,
                        variant: SudaButtonVariant.dangerOutline,
                        size: SudaButtonSize.sm,
                        fullWidth: true,
                        onPressed: onDecline,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
