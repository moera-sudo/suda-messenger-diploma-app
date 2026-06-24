import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import '../../../../app/config/theme/app_theme.dart';
import '../../../../app/navigation/app_routes.dart';
import '../../../../features/channels/data/models/channel_view_model.dart';
import '../../../../features/channels/data/models/gating_rule.dart';
import '../../../../features/channels/domain/repositories/channel_repository.dart';
import '../../../../shared/domain/models/app_failure.dart';
import '../../../../shared/presentation/feedback/app_feedback.dart';
import '../../../../shared/presentation/format_utils.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/widgets/suda_avatar.dart';
import '../../../../shared/presentation/widgets/suda_button.dart';
import '../../../../features/channels/presentation/widgets/donate_sheet.dart';
import '../../domain/repositories/chat_repository.dart';

/// Channel info screen (broadcast channels).
class ChannelInfoPage extends StatefulWidget {
  final String chatId;
  final bool initialIsMuted;

  const ChannelInfoPage({super.key, required this.chatId, this.initialIsMuted = false});

  @override
  State<ChannelInfoPage> createState() => _ChannelInfoPageState();
}

class _ChannelInfoPageState extends State<ChannelInfoPage> {
  late Future<ChannelViewModel?> _dataFuture;
  // Resolved gating rule — drives whether the channel is paid (hides profile subscribe).
  GatingRule? _gatingRule;
  late bool _isMuted;

  @override
  void initState() {
    super.initState();
    _isMuted = widget.initialIsMuted;
    _dataFuture = _loadData();
    _loadGating();
  }

  Future<ChannelViewModel?> _loadData() async {
    final result = await sl<ChannelRepository>().getChannelView(widget.chatId);
    return result.fold((_) => null, (view) => view);
  }

  Future<void> _loadGating() async {
    final result = await sl<ChannelRepository>().getGatingRule(widget.chatId);
    final rule = result.fold((_) => null, (r) => r);
    if (mounted) setState(() => _gatingRule = rule);
  }

  void _reload() {
    setState(() {
      _dataFuture = _loadData();
    });
    _loadGating();
  }

  // ── Subscribe ──────────────────────────────────────────────────

  /// Free-channel subscribe only. Paid channels are subscribed via the opening
  /// paywall in the channel chat, never from the profile.
  Future<void> _onSubscribe() async {
    final result = await sl<ChannelRepository>().subscribe(widget.chatId);
    if (!mounted) return;
    result.fold(
      (failure) {
        final l10n = context.l10n;
        final code = failure is ServerFailure ? failure.code : null;
        if (code == 'INSUFFICIENT_BALANCE') {
          AppFeedback.showError(l10n.channelInsufficientBalance);
        } else if (code == 'NO_WALLET') {
          AppFeedback.showError(l10n.channelNoWallet);
        } else {
          AppFeedback.showError(l10n.channelSubscribeFailed);
        }
      },
      (_) {
        AppFeedback.showSuccess(context.l10n.channelSubscribeSuccess);
        _reload();
      },
    );
  }

  Future<void> _onJoinRequest() async {
    final result = await sl<ChannelRepository>().sendJoinRequest(widget.chatId);
    if (!mounted) return;
    result.fold(
      (f) => AppFeedback.showError(f.message),
      (_) {
        AppFeedback.showSuccess(context.l10n.channelRequestSent);
        _reload();
      },
    );
  }

  Future<void> _onCancelRequest() async {
    final result = await sl<ChannelRepository>().cancelJoinRequest(widget.chatId);
    if (!mounted) return;
    result.fold(
      (f) => AppFeedback.showError(f.message),
      (_) => _reload(),
    );
  }

  Future<void> _onAcceptInvite() async {
    final result = await sl<ChannelRepository>().acceptInvite(widget.chatId);
    if (!mounted) return;
    result.fold(
      (f) => AppFeedback.showError(f.message),
      (_) {
        AppFeedback.showSuccess(context.l10n.channelSubscribeSuccess);
        _reload();
      },
    );
  }

  Future<void> _onDeclineInvite() async {
    final result = await sl<ChannelRepository>().declineInvite(widget.chatId);
    if (!mounted) return;
    result.fold(
      (f) => AppFeedback.showError(f.message),
      (_) => context.pop(),
    );
  }

  Future<void> _onInviteUser() async {
    final l10n = context.l10n;
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final ctrl = TextEditingController();

    final confirmed = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        backgroundColor: theme.colorScheme.surface,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(18)),
        title: Text(l10n.channelInvite,
            style: TextStyle(fontFamily: 'Manrope', fontWeight: FontWeight.w700,
                color: theme.colorScheme.onSurface)),
        content: TextField(
          controller: ctrl,
          autofocus: true,
          style: TextStyle(fontFamily: 'Manrope', color: theme.colorScheme.onSurface),
          decoration: InputDecoration(
            hintText: l10n.channelInviteUsernameHint,
            hintStyle: TextStyle(color: palette.textSecondary),
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx, false),
            child: Text(l10n.buttonRetry,
                style: TextStyle(fontFamily: 'Manrope', color: palette.textSecondary)),
          ),
          TextButton(
            onPressed: () => Navigator.pop(ctx, true),
            child: Text(l10n.channelInvite,
                style: TextStyle(fontFamily: 'Manrope', color: theme.colorScheme.tertiary,
                    fontWeight: FontWeight.w700)),
          ),
        ],
      ),
    );
    if (confirmed != true || ctrl.text.trim().isEmpty) return;
    final username = ctrl.text.trim().replaceFirst('@', '');
    final result = await sl<ChannelRepository>().inviteUser(widget.chatId, username);
    if (!mounted) return;
    result.fold(
      (f) => AppFeedback.showError(f.message),
      (_) => AppFeedback.showSuccess(context.l10n.channelInviteSent),
    );
  }


  void _confirmLeave() {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        backgroundColor: theme.colorScheme.surface,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(18)),
        content: Text(l10n.channelLeaveConfirm,
            style: TextStyle(fontFamily: 'Manrope', color: theme.colorScheme.onSurface)),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: Text(l10n.buttonRetry,
                style: TextStyle(fontFamily: 'Manrope', color: palette.textSecondary)),
          ),
          TextButton(
            onPressed: () { Navigator.pop(ctx); _doLeave(); },
            child: Text(l10n.channelLeave,
                style: TextStyle(fontFamily: 'Manrope', color: palette.danger,
                    fontWeight: FontWeight.w700)),
          ),
        ],
      ),
    );
  }

  Future<void> _doLeave() async {
    final result = await sl<ChannelRepository>().unsubscribe(widget.chatId);
    if (!mounted) return;
    result.fold((_) {}, (_) {
      AppFeedback.showSuccess(context.l10n.channelUnsubscribeSuccess);
      context.go('/chats');
    });
  }

  void _confirmDelete() {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        backgroundColor: theme.colorScheme.surface,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(18)),
        content: Text(l10n.channelDeleteConfirm,
            style: TextStyle(fontFamily: 'Manrope', color: theme.colorScheme.onSurface)),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: Text(l10n.buttonRetry,
                style: TextStyle(fontFamily: 'Manrope', color: palette.textSecondary)),
          ),
          TextButton(
            onPressed: () { Navigator.pop(ctx); _doDelete(); },
            child: Text(l10n.channelDelete,
                style: TextStyle(fontFamily: 'Manrope', color: palette.danger,
                    fontWeight: FontWeight.w700)),
          ),
        ],
      ),
    );
  }

  Future<void> _doDelete() async {
    final result = await sl<ChatRepository>().deleteChat(widget.chatId, scope: 'everyone');
    if (!mounted) return;
    result.fold((_) {}, (_) => context.go('/chats'));
  }

  Future<void> _onToggleMute() async {
    final repo = sl<ChatRepository>();
    final result = _isMuted
        ? await repo.unmuteChat(widget.chatId)
        : await repo.muteChat(widget.chatId);
    if (!mounted) return;
    result.fold(
      (f) => AppFeedback.showError(f.message),
      (_) => setState(() => _isMuted = !_isMuted),
    );
  }

  void _onDonate() => DonateSheet.show(context, channelId: widget.chatId);

  Future<void> _onSettings() async {
    await context.pushNamed(
      AppRoutes.channelSettings,
      pathParameters: {'id': widget.chatId},
    );
    _reload(); // reflect any saved changes
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      body: FutureBuilder<ChannelViewModel?>(
        future: _dataFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          }
          final view = snapshot.data;
          if (view == null) {
            return _ErrorScaffoldBody(message: l10n.errorGeneric, danger: palette.danger);
          }

          // Pending invite: show Accept/Decline before anything else
          if (view.joinStatus == 'pending_invite') {
            return _PendingInviteScaffoldBody(
              view: view,
              onAccept: _onAcceptInvite,
              onDecline: _onDeclineInvite,
            );
          }

          // Private channel, not yet a member (none | pending_request)
          if (!view.canRead) {
            return _NonMemberScaffoldBody(
              view: view,
              onJoinRequest: _onJoinRequest,
              onCancelRequest: _onCancelRequest,
            );
          }

          // Can read (public non-member OR member OR owner)
          return CustomScrollView(
            slivers: [
              SliverToBoxAdapter(child: _ChannelHero(view: view)),
              if (view.isOwner)
                SliverToBoxAdapter(
                  child: _OwnerActionsBar(
                    onInvite: _onInviteUser,
                    onSettings: _onSettings,
                    onTreasury: () => context.pushNamed(
                      AppRoutes.channelTreasury,
                      pathParameters: {'id': widget.chatId},
                      extra: {'isOwner': true},
                    ),
                  ),
                )
              else if (view.isMember)
                SliverToBoxAdapter(
                  child: _SubscriberActionsBar(
                    onDonate: _onDonate,
                    isMuted: _isMuted,
                    onMute: _onToggleMute,
                  ),
                )
              else
                // Public visitor (can_read=true, not subscribed). Paid channels
                // hide the subscribe action — payment happens via the opening
                // paywall in the channel chat, not here.
                SliverToBoxAdapter(
                  child: _PublicVisitorActionsBar(
                    onSubscribe: _onSubscribe,
                    isPaidGated: (_gatingRule?.subscriptionPrice ?? 0) > 0,
                    onDonate: _onDonate,
                    isMuted: _isMuted,
                    onMute: _onToggleMute,
                  ),
                ),
              SliverToBoxAdapter(
                child: _ChannelContent(
                  view: view,
                  chatId: widget.chatId,
                  chatName: view.name,
                  onLeave: _confirmLeave,
                  onDelete: _confirmDelete,
                ),
              ),
              const SliverToBoxAdapter(child: SizedBox(height: 32)),
            ],
          );
        },
      ),
    );
  }
}

// ── Non-member body ─────────────────────────────────────────────

class _NonMemberScaffoldBody extends StatelessWidget {
  final ChannelViewModel view;
  final VoidCallback onJoinRequest;
  final VoidCallback onCancelRequest;

  const _NonMemberScaffoldBody({
    required this.view,
    required this.onJoinRequest,
    required this.onCancelRequest,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    final isPendingRequest = view.joinStatus == 'pending_request';

    return SafeArea(
      child: Column(
        children: [
          Align(
            alignment: Alignment.centerLeft,
            child: IconButton(
              icon: Icon(Icons.arrow_back_rounded, color: theme.colorScheme.onSurface),
              onPressed: () => context.pop(),
            ),
          ),
          Expanded(
            child: Center(
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 32),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    SudaAvatar(
                      mediaId: view.avatarMediaId,
                      initials: view.name,
                      size: 72,
                      ring: true,
                    ),
                    const SizedBox(height: 14),
                    Text(
                      view.name,
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontFamily: 'SpaceGrotesk',
                        fontSize: 20,
                        fontWeight: FontWeight.w800,
                        color: theme.colorScheme.onSurface,
                      ),
                    ),
                    if (view.description != null && view.description!.isNotEmpty) ...[
                      const SizedBox(height: 8),
                      Text(
                        view.description!,
                        textAlign: TextAlign.center,
                        maxLines: 3,
                        overflow: TextOverflow.ellipsis,
                        style: TextStyle(
                          fontFamily: 'Manrope',
                          fontSize: 13,
                          color: palette.textSecondary,
                          height: 1.5,
                        ),
                      ),
                    ],
                    const SizedBox(height: 8),
                    Text(
                      l10n.channelSubscribersCount(view.subscriberCount),
                      style: TextStyle(fontFamily: 'Manrope', fontSize: 13, color: palette.textSecondary),
                    ),
                    const SizedBox(height: 24),
                    Icon(Icons.lock_outline_rounded, size: 40, color: palette.textSecondary),
                    const SizedBox(height: 12),
                    Text(
                      l10n.channelPrivateBanner,
                      textAlign: TextAlign.center,
                      style: TextStyle(fontFamily: 'Manrope', fontSize: 15, color: palette.textSecondary),
                    ),
                  ],
                ),
              ),
            ),
          ),
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 0, 16, 16),
            child: isPendingRequest
                ? Column(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Text(
                        l10n.channelRequestSent,
                        style: TextStyle(fontFamily: 'Manrope', fontSize: 14,
                            color: palette.textSecondary),
                      ),
                      const SizedBox(height: 8),
                      SudaButton(
                        label: l10n.channelCancelRequest,
                        variant: SudaButtonVariant.dangerOutline,
                        size: SudaButtonSize.lg,
                        fullWidth: true,
                        onPressed: onCancelRequest,
                      ),
                    ],
                  )
                : SudaButton(
                    label: l10n.channelJoinRequest,
                    variant: SudaButtonVariant.primary,
                    size: SudaButtonSize.lg,
                    fullWidth: true,
                    onPressed: onJoinRequest,
                  ),
          ),
        ],
      ),
    );
  }
}

// ── Pending Invite body ────────────────────────────────────────

class _PendingInviteScaffoldBody extends StatelessWidget {
  final ChannelViewModel view;
  final VoidCallback onAccept;
  final VoidCallback onDecline;

  const _PendingInviteScaffoldBody({
    required this.view,
    required this.onAccept,
    required this.onDecline,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    return SafeArea(
      child: Column(
        children: [
          Align(
            alignment: Alignment.centerLeft,
            child: IconButton(
              icon: Icon(Icons.arrow_back_rounded, color: theme.colorScheme.onSurface),
              onPressed: () => context.pop(),
            ),
          ),
          Expanded(
            child: Center(
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 32),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    SudaAvatar(
                      mediaId: view.avatarMediaId,
                      initials: view.name,
                      size: 72,
                      ring: true,
                    ),
                    const SizedBox(height: 14),
                    Text(
                      view.name,
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontFamily: 'SpaceGrotesk',
                        fontSize: 20,
                        fontWeight: FontWeight.w800,
                        color: theme.colorScheme.onSurface,
                      ),
                    ),
                    if (view.description != null && view.description!.isNotEmpty) ...[
                      const SizedBox(height: 8),
                      Text(
                        view.description!,
                        textAlign: TextAlign.center,
                        maxLines: 3,
                        overflow: TextOverflow.ellipsis,
                        style: TextStyle(
                          fontFamily: 'Manrope',
                          fontSize: 13,
                          color: palette.textSecondary,
                          height: 1.5,
                        ),
                      ),
                    ],
                    const SizedBox(height: 8),
                    Text(
                      l10n.channelSubscribersCount(view.subscriberCount),
                      style: TextStyle(fontFamily: 'Manrope', fontSize: 13, color: palette.textSecondary),
                    ),
                    const SizedBox(height: 24),
                    Icon(Icons.mail_outline_rounded, size: 40, color: theme.colorScheme.tertiary),
                    const SizedBox(height: 12),
                    Text(
                      l10n.channelAcceptInvite,
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontFamily: 'Manrope',
                        fontSize: 15,
                        color: palette.textSecondary,
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ),
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 0, 16, 16),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                SudaButton(
                  label: l10n.channelAcceptInvite,
                  variant: SudaButtonVariant.primary,
                  size: SudaButtonSize.lg,
                  fullWidth: true,
                  onPressed: onAccept,
                ),
                const SizedBox(height: 8),
                SudaButton(
                  label: l10n.channelDeclineInvite,
                  variant: SudaButtonVariant.dangerOutline,
                  size: SudaButtonSize.lg,
                  fullWidth: true,
                  onPressed: onDecline,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

// ── Error body ─────────────────────────────────────────────────

class _ErrorScaffoldBody extends StatelessWidget {
  final String message;
  final Color danger;
  const _ErrorScaffoldBody({required this.message, required this.danger});

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: Column(children: [
        Align(
          alignment: Alignment.centerLeft,
          child: IconButton(
            icon: Icon(Icons.arrow_back_rounded, color: Theme.of(context).colorScheme.onSurface),
            onPressed: () => context.pop(),
          ),
        ),
        Expanded(child: Center(child: Text(message, style: TextStyle(color: danger)))),
      ]),
    );
  }
}

// ── Hero ───────────────────────────────────────────────────────

class _ChannelHero extends StatelessWidget {
  final ChannelViewModel view;
  const _ChannelHero({required this.view});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    final handle = (view.username?.isNotEmpty == true) ? formatHandle(view.username) : null;
    final subscribers = l10n.channelSubscribersCount(view.subscriberCount);
    final subtitle = handle != null ? '$handle · $subscribers' : subscribers;

    return Container(
      constraints: const BoxConstraints(minHeight: 220),
      child: Stack(children: [
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
        Align(
          alignment: Alignment.bottomCenter,
          child: Padding(
            padding: const EdgeInsets.fromLTRB(18, 80, 18, 20),
            child: Column(mainAxisSize: MainAxisSize.min, children: [
              SudaAvatar(
                mediaId: view.avatarMediaId,
                initials: view.name.isNotEmpty ? view.name : '?',
                size: 96,
                ring: true,
              ),
              const SizedBox(height: 10),
              Text(
                view.name,
                style: const TextStyle(
                    fontFamily: 'SpaceGrotesk', fontSize: 22, fontWeight: FontWeight.w800),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 4),
              Text(
                subtitle,
                style: TextStyle(fontFamily: 'Manrope', fontSize: 13, color: palette.textSecondary),
                textAlign: TextAlign.center,
              ),
            ]),
          ),
        ),
      ]),
    );
  }
}

// ── Action bars ─────────────────────────────────────────────────

class _PublicVisitorActionsBar extends StatelessWidget {
  final VoidCallback onSubscribe;
  // Paid channels hide the subscribe action — payment goes through the opening paywall.
  final bool isPaidGated;
  final VoidCallback onDonate;
  final bool isMuted;
  final VoidCallback onMute;
  const _PublicVisitorActionsBar({
    required this.onSubscribe,
    required this.isPaidGated,
    required this.onDonate,
    required this.isMuted,
    required this.onMute,
  });

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
      child: Row(children: [
        if (!isPaidGated) ...[
          _ChannelAction(icon: Icons.add_circle_outline_rounded, label: l10n.channelSubscribe, onTap: onSubscribe),
          const SizedBox(width: 8),
        ],
        _ChannelAction(icon: Icons.favorite_border_rounded, label: l10n.userProfileDonate, onTap: onDonate),
        const SizedBox(width: 8),
        _ChannelAction(
          icon: isMuted ? Icons.notifications_rounded : Icons.notifications_none_rounded,
          label: isMuted ? l10n.chatUnmute : l10n.chatMute,
          onTap: onMute,
        ),
      ]),
    );
  }
}

class _SubscriberActionsBar extends StatelessWidget {
  final VoidCallback onDonate;
  final bool isMuted;
  final VoidCallback onMute;
  const _SubscriberActionsBar({
    required this.onDonate,
    required this.isMuted,
    required this.onMute,
  });

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
      child: Row(children: [
        _ChannelAction(icon: Icons.favorite_border_rounded, label: l10n.userProfileDonate, onTap: onDonate),
        const SizedBox(width: 8),
        _ChannelAction(
          icon: isMuted ? Icons.notifications_rounded : Icons.notifications_none_rounded,
          label: isMuted ? l10n.chatUnmute : l10n.chatMute,
          onTap: onMute,
        ),
      ]),
    );
  }
}

class _OwnerActionsBar extends StatelessWidget {
  final VoidCallback onInvite;
  final VoidCallback onSettings;
  final VoidCallback onTreasury;
  const _OwnerActionsBar({
    required this.onInvite,
    required this.onSettings,
    required this.onTreasury,
  });

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
      child: Row(children: [
        _ChannelAction(icon: Icons.person_add_alt_rounded, label: l10n.channelInvite, onTap: onInvite),
        const SizedBox(width: 8),
        _ChannelAction(icon: Icons.account_balance_rounded, label: l10n.channelTreasury, onTap: onTreasury),
        const SizedBox(width: 8),
        _ChannelAction(icon: Icons.settings_outlined, label: l10n.channelSettings, onTap: onSettings),
      ]),
    );
  }
}

class _ChannelAction extends StatelessWidget {
  final IconData icon;
  final String label;
  final VoidCallback onTap;
  const _ChannelAction({required this.icon, required this.label, required this.onTap});

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
            child: Column(mainAxisSize: MainAxisSize.min, children: [
              Icon(icon, size: 18, color: theme.colorScheme.tertiary),
              const SizedBox(height: 6),
              Text(
                label,
                style: TextStyle(
                    fontFamily: 'Manrope', fontSize: 12, fontWeight: FontWeight.w600,
                    color: palette.textSecondary),
                textAlign: TextAlign.center,
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
              ),
            ]),
          ),
        ),
      ),
    );
  }
}

// ── Content ────────────────────────────────────────────────────

class _ChannelContent extends StatelessWidget {
  final ChannelViewModel view;
  final String chatId;
  final String chatName;
  final VoidCallback onLeave;
  final VoidCallback onDelete;

  const _ChannelContent({
    required this.view,
    required this.chatId,
    required this.chatName,
    required this.onLeave,
    required this.onDelete,
  });

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final isOwner = view.isOwner;
    final isMember = view.isMember;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        if (view.description != null && view.description!.isNotEmpty) ...[
          Padding(
            padding: const EdgeInsets.fromLTRB(10, 8, 10, 0),
            child: _InfoCard(
              label: l10n.channelDescription,
              child: Text(
                view.description!,
                style: const TextStyle(fontFamily: 'Manrope', fontSize: 15, height: 1.6),
              ),
            ),
          ),
          const SizedBox(height: 12),
        ],

        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 10),
          child: _InfoListCard(children: [
            _InfoListRow(
              icon: Icons.image_outlined,
              label: l10n.userProfileSharedMedia,
              onTap: () => context.pushNamed(
                AppRoutes.sharedMedia,
                pathParameters: {'id': chatId},
                extra: {'chatName': chatName},
              ),
            ),
            _InfoListRow(
              icon: Icons.search_rounded,
              label: l10n.channelSearchIn,
              onTap: () => context.pushNamed(
                AppRoutes.chatSearch,
                pathParameters: {'id': chatId},
                extra: {'chatName': chatName},
              ),
            ),
            if (isMember && !isOwner)
              _InfoListRow(
                icon: Icons.exit_to_app_rounded,
                label: l10n.channelLeave,
                isDanger: true,
                showChevron: false,
                onTap: onLeave,
              ),
            if (isOwner)
              _InfoListRow(
                icon: Icons.delete_outline_rounded,
                label: l10n.channelDelete,
                isDanger: true,
                showChevron: false,
                onTap: onDelete,
              ),
          ]),
        ),
        const SizedBox(height: 8),
      ],
    );
  }
}

// ── Reusable local card widgets ─────────────────────────────────

class _InfoCard extends StatelessWidget {
  final String label;
  final Widget child;
  const _InfoCard({required this.label, required this.child});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;

    return Container(
      width: double.infinity,
      padding: const EdgeInsets.fromLTRB(16, 14, 16, 18),
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        border: Border.all(color: palette.divider),
        borderRadius: BorderRadius.circular(16),
      ),
      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Text(
          label.toUpperCase(),
          style: TextStyle(
              fontFamily: 'Manrope', fontSize: 11, fontWeight: FontWeight.w700,
              letterSpacing: 0.1 * 11, color: palette.textSecondary),
        ),
        const SizedBox(height: 10),
        child,
      ]),
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
  final bool isDanger;
  final bool showChevron;
  final VoidCallback? onTap;

  const _InfoListRow({
    required this.icon,
    required this.label,
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
            border: Border(bottom: BorderSide(color: palette.divider, width: 0.5))),
        child: Row(children: [
          Icon(icon, size: 20, color: iconColor),
          const SizedBox(width: 12),
          Expanded(
            child: Text(label,
                style: TextStyle(
                    fontFamily: 'Manrope', fontSize: 14, fontWeight: FontWeight.w500,
                    color: textColor)),
          ),
          if (showChevron)
            Icon(Icons.chevron_right_rounded, size: 18, color: palette.textSecondary),
        ]),
      ),
    );
  }
}
