import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';

import 'package:messenger_app_v2/app/navigation/app_routes.dart';
import '../../../../app/config/theme/app_theme.dart';
import '../../../../features/channels/data/models/channel_view_model.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/widgets/suda_avatar.dart';
import '../bloc/chat_detail_bloc.dart';

class ChatAppBar extends StatelessWidget implements PreferredSizeWidget {
  final String title;
  final String chatId;
  final String interlocutorId;
  final String chatType;
  final String myRole;
  // Channel-only: viewer relationship + overflow-menu actions.
  final ChannelViewModel? channelView;
  final bool isMuted;
  final VoidCallback? onSubscribe;
  final VoidCallback? onUnsubscribe;
  final VoidCallback? onDonate;
  final VoidCallback? onToggleMute;
  // Called after returning from chat/channel info page so the parent can refresh state.
  final VoidCallback? onInfoReturn;

  const ChatAppBar({
    super.key,
    required this.title,
    required this.chatId,
    this.interlocutorId = '',
    this.chatType = 'DIRECT',
    this.myRole = 'MEMBER',
    this.channelView,
    this.isMuted = false,
    this.onSubscribe,
    this.onUnsubscribe,
    this.onDonate,
    this.onToggleMute,
    this.onInfoReturn,
  });

  String _formatLastSeen(BuildContext context, String? iso) {
    if (iso == null) return context.l10n.chatStatusOffline;
    try {
      final d = DateTime.parse(iso).toLocal();
      final diff = DateTime.now().difference(d);
      // Guard against unrealistic values (clock skew → future timestamp, or a
      // zero/epoch timestamp → absurdly large diff). Don't render nonsense.
      if (diff.isNegative) {
        return '${context.l10n.chatLastSeenPrefix} ${context.l10n.chatLastSeenJustNow}';
      }
      if (diff.inDays > 3650) {
        return context.l10n.chatStatusOffline;
      }
      if (diff.inSeconds < 60) {
        return '${context.l10n.chatLastSeenPrefix} ${context.l10n.chatLastSeenJustNow}';
      } else if (diff.inMinutes < 60) {
        return '${context.l10n.chatLastSeenPrefix} ${diff.inMinutes} ${context.l10n.chatLastSeenMinutes}';
      } else if (diff.inHours < 24) {
        return '${context.l10n.chatLastSeenPrefix} ${diff.inHours} ${context.l10n.chatLastSeenHours}';
      } else {
        final dateStr =
            '${d.day.toString().padLeft(2, '0')}.${d.month.toString().padLeft(2, '0')}.${d.year}';
        return '${context.l10n.chatLastSeenPrefix} $dateStr';
      }
    } catch (_) {
      return context.l10n.chatStatusOffline;
    }
  }

  void _onTitleTap(BuildContext context) {
    if (chatType == 'DIRECT' && interlocutorId.trim().isNotEmpty) {
      // DIRECT — open user profile, passing the chat so it can offer shared media.
      context.pushNamed(AppRoutes.chatProfile,
          pathParameters: {'id': interlocutorId.trim()},
          extra: {'chatId': chatId, 'chatName': title});
    } else if (chatType == 'GROUP' || chatType == 'CHANNEL') {
      // GROUP / CHANNEL — open group/channel info
      context.pushNamed(AppRoutes.chatInfo,
          pathParameters: {'id': chatId},
          extra: {'chatType': chatType, 'myRole': myRole, 'isMuted': isMuted},
      ).then((_) => onInfoReturn?.call());
    }
  }

  /// Channel overflow menu: Subscribe/Unsubscribe, Donate, Mute/Unmute.
  Widget _buildChannelMenu(BuildContext context, ThemeData theme, AppPaletteTheme palette) {
    final l10n = context.l10n;
    final view = channelView!;
    final textStyle = TextStyle(fontFamily: 'Manrope', color: theme.colorScheme.onSurface);

    return PopupMenuButton<String>(
      icon: Icon(Icons.more_vert_rounded, color: palette.textSecondary),
      color: theme.colorScheme.surface,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
      onSelected: (value) {
        switch (value) {
          case 'subscribe':
            onSubscribe?.call();
          case 'unsubscribe':
            onUnsubscribe?.call();
          case 'donate':
            onDonate?.call();
          case 'mute':
            onToggleMute?.call();
        }
      },
      itemBuilder: (_) => [
        if (!view.isMember)
          PopupMenuItem(value: 'subscribe', child: Text(l10n.channelSubscribe, style: textStyle)),
        if (view.isMember && !view.isOwner)
          PopupMenuItem(value: 'unsubscribe', child: Text(l10n.channelUnsubscribe, style: textStyle)),
        PopupMenuItem(value: 'donate', child: Text(l10n.userProfileDonate, style: textStyle)),
        PopupMenuItem(
          value: 'mute',
          child: Text(isMuted ? l10n.userProfileUnmute : l10n.userProfileMute, style: textStyle),
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);

    return AppBar(
      backgroundColor: theme.scaffoldBackgroundColor,
      elevation: 0,
      scrolledUnderElevation: 0,
      leading: IconButton(
        icon: Icon(Icons.arrow_back_rounded, color: theme.colorScheme.onSurface),
        onPressed: () => context.pop(),
      ),
      titleSpacing: 0,
      title: BlocBuilder<ChatDetailBloc, ChatDetailState>(
        buildWhen: (p, c) =>
            p.isInterlocutorOnline != c.isInterlocutorOnline ||
            p.isInterlocutorTyping != c.isInterlocutorTyping ||
            p.interlocutorLastSeenAt != c.interlocutorLastSeenAt ||
            p.interlocutorAvatarMediaId != c.interlocutorAvatarMediaId,
        builder: (context, state) {
          // Presence (online / typing / last seen) is only meaningful in DIRECT
          // chats. Groups, channels and Saved Messages have no single interlocutor.
          final hasPresence = chatType == 'DIRECT';
          final isTyping = hasPresence && state.isInterlocutorTyping;
          final isOnline = hasPresence && state.isInterlocutorOnline;

          final String statusText;
          final Color statusColor;
          if (chatType == 'CHANNEL' && channelView != null) {
            // Channels show subscriber count instead of presence.
            statusText = context.l10n.channelSubscribersCount(channelView!.subscriberCount);
            statusColor = palette.textSecondary;
          } else if (!hasPresence) {
            statusText = '';
            statusColor = palette.textSecondary;
          } else if (isTyping) {
            statusText = context.l10n.chatStatusTyping;
            statusColor = theme.colorScheme.tertiary;
          } else if (isOnline) {
            statusText = context.l10n.chatStatusOnline;
            statusColor = palette.success;
          } else {
            statusText = _formatLastSeen(context, state.interlocutorLastSeenAt);
            statusColor = palette.textSecondary;
          }

          // Tapping title/avatar navigates to profile or chat info (no "i" button needed).
          return GestureDetector(
            onTap: () => _onTitleTap(context),
            behavior: HitTestBehavior.opaque,
            child: Row(children: [
              SudaAvatar(
                // DIRECT: resolve the interlocutor's avatar; CHANNEL: the
                // channel avatar; GROUP has no single avatar in the header.
                mediaId: chatType == 'CHANNEL'
                    ? channelView?.avatarMediaId
                    : state.interlocutorAvatarMediaId,
                initials: title.isNotEmpty ? title[0] : '?',
                size: 36,
                showOnline: isOnline,
              ),
              const SizedBox(width: 10),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  mainAxisSize: MainAxisSize.min,
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
                    if (statusText.isNotEmpty)
                      Text(
                        statusText,
                        style: TextStyle(
                          fontFamily: 'Manrope',
                          fontSize: 12,
                          fontStyle:
                              isTyping ? FontStyle.italic : FontStyle.normal,
                          color: statusColor,
                        ),
                      ),
                  ],
                ),
              ),
            ]),
          );
        },
      ),
      actions: [
        IconButton(
          icon: Icon(Icons.search_rounded, color: palette.textSecondary),
          onPressed: () async {
            // Capture the bloc before awaiting so we don't touch a stale context.
            final bloc = context.read<ChatDetailBloc>();
            final target = await context.pushNamed(
              AppRoutes.chatSearch,
              pathParameters: {'id': chatId},
              extra: {'chatName': title},
            );
            // Search returns the tapped message id — jump to it in the chat.
            if (target is int) bloc.add(JumpToMessage(target));
          },
        ),
        if (chatType == 'CHANNEL' && channelView != null)
          _buildChannelMenu(context, theme, palette),
      ],
      bottom: PreferredSize(
        preferredSize: const Size.fromHeight(1),
        child: Divider(height: 1, color: palette.divider),
      ),
    );
  }

  @override
  Size get preferredSize => const Size.fromHeight(kToolbarHeight + 1);
}
