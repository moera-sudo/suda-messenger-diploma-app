import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import '../../../../app/config/theme/app_theme.dart';
import '../../../../app/navigation/app_routes.dart';
import '../../../../features/wallet/data/models/wallet_models.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/placeholder/placeholder_page.dart';
import '../../../../shared/presentation/widgets/suda_avatar.dart';
import '../../data/models/chat_models.dart';
import '../bloc/chat_list_bloc.dart';

class ChatsPage extends StatefulWidget {
  const ChatsPage({super.key});

  @override
  State<ChatsPage> createState() => _ChatsPageState();
}

class _ChatsPageState extends State<ChatsPage> {
  // ignore: prefer_final_fields — intentionally mutable for future filter-chip feature
  String _activeFilter = 'all';

  List<ChatListItem> _filterChats(List<ChatListItem> chats) {
    return switch (_activeFilter) {
      'unread' => chats.where((c) => c.unreadCount > 0).toList(),
      'personal' => chats.where((c) => c.type == 'DIRECT' || c.type == 'SAVED').toList(),
      'groups' => chats.where((c) => c.type == 'GROUP').toList(),
      'channels' => chats.where((c) => c.type == 'CHANNEL').toList(),
      _ => chats,
    };
  }

  String _formatTime(String isoTime) {
    if (isoTime.isEmpty) return '';
    try {
      final date = DateTime.parse(isoTime).toLocal();
      final now = DateTime.now();
      if (date.year == now.year && date.month == now.month && date.day == now.day) {
        return '${date.hour.toString().padLeft(2, '0')}:${date.minute.toString().padLeft(2, '0')}';
      }
      return '${date.day.toString().padLeft(2, '0')}.${date.month.toString().padLeft(2, '0')}';
    } catch (_) {
      return '';
    }
  }

  void _showNewChatSheet(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);
    final l10n = context.l10n;

    showModalBottomSheet(
      context: context,
      backgroundColor: theme.colorScheme.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      ),
      builder: (ctx) => SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              margin: const EdgeInsets.only(top: 8, bottom: 8),
              width: 38,
              height: 4,
              decoration: BoxDecoration(
                color: palette.textSecondary.withValues(alpha: 0.4),
                borderRadius: BorderRadius.circular(2),
              ),
            ),
            _NewChatAction(
              icon: Icons.group_rounded,
              title: l10n.newGroup,
              onTap: () {
                Navigator.pop(ctx);
                context.pushNamed(AppRoutes.chatNew);
              },
            ),
            _NewChatAction(
              icon: Icons.campaign_rounded,
              title: l10n.newChannel,
              onTap: () {
                Navigator.pop(ctx);
                context.pushNamed(AppRoutes.channelNew);
              },
            ),
            _NewChatAction(
              icon: Icons.person_search_rounded,
              title: l10n.chatFindOrStart,
              onTap: () {
                Navigator.pop(ctx);
                context.push('/search');
              },
            ),
            const SizedBox(height: 8),
          ],
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);
    final l10n = context.l10n;

    return BlocProvider(
      create: (context) => sl<ChatListBloc>()..add(LoadChats()),
      child: Scaffold(
        backgroundColor: theme.scaffoldBackgroundColor,
        body: BlocBuilder<ChatListBloc, ChatListState>(
          buildWhen: (prev, next) =>
              prev.status != next.status || prev.chats != next.chats,
          builder: (context, state) {
            final allChats = state.chats;
            final filtered = _filterChats(allChats);

            return CustomScrollView(
              slivers: [
                // ── AppBar ───────────────────────────────────
                SliverAppBar(
                  pinned: true,
                  backgroundColor: theme.scaffoldBackgroundColor,
                  elevation: 0,
                  scrolledUnderElevation: 0,
                  title: Text(
                    l10n.chatTitle,
                    style: TextStyle(
                      fontFamily: 'Manrope',
                      fontSize: 18,
                      fontWeight: FontWeight.w700,
                      color: theme.colorScheme.onSurface,
                    ),
                  ),
                  actions: [
                    IconButton(
                      icon: Icon(Icons.account_balance_wallet_outlined, color: palette.textAccent),
                      onPressed: () => context.pushNamed(AppRoutes.wallet),
                      tooltip: l10n.walletOpenWallet,
                    ),

                    IconButton(
                      icon: Icon(Icons.edit_outlined, color: theme.colorScheme.onSurface),
                      onPressed: () => _showNewChatSheet(context),
                    ),
                  ],
                ),

                // ── Search bar ───────────────────────────────
                SliverToBoxAdapter(
                  child: Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 4),
                    child: GestureDetector(
                      onTap: () {
                        final recent = context
                            .read<ChatListBloc>()
                            .state
                            .chats
                            .take(4)
                            .toList();
                        context.push('/search', extra: recent);
                      },
                      child: Container(
                        height: 44,
                        padding: const EdgeInsets.symmetric(horizontal: 14),
                        decoration: BoxDecoration(
                          color: theme.colorScheme.surface,
                          borderRadius: BorderRadius.circular(14),
                          border: Border.all(color: palette.divider),
                        ),
                        child: Row(
                          children: [
                            Icon(Icons.search_rounded, size: 18, color: palette.textSecondary),
                            const SizedBox(width: 10),
                            Text(
                              l10n.chatFindOrStart,
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
                ),

                // ── Filter chips (hidden — foundation for future folder tabs) ──
                // TODO: Uncomment when folder/tab feature is implemented
                // SliverToBoxAdapter(child: _FilterChipsRow(...)),

                // ── Chat list ────────────────────────────────
                if (state.status == ChatListStatus.loading)
                  const SliverFillRemaining(
                    child: Center(child: CircularProgressIndicator()),
                  )
                else if (state.status == ChatListStatus.failure)
                  SliverFillRemaining(
                    child: Center(
                      child: Text(
                        state.error ?? l10n.errorGeneric,
                        style: TextStyle(color: palette.danger),
                      ),
                    ),
                  )
                else if (filtered.isEmpty)
                  const SliverFillRemaining(
                    child: PlaceholderPage(type: PlaceholderType.noContent, scaffold: false),
                  )
                else
                  SliverList.separated(
                    itemCount: filtered.length,
                    separatorBuilder: (_, __) => const SizedBox.shrink(),
                    itemBuilder: (context, index) {
                      final chat = filtered[index];
                      final title = switch (chat.type) {
                        'SAVED' => l10n.chatSavedMessages,
                        'GROUP' || 'CHANNEL' => chat.name ?? l10n.chatUnknownUser,
                        _ => chat.interlocutorName ?? l10n.chatUnknownUser,
                      };
                      return _ChatRow(
                        chat: chat,
                        title: title,
                        timeStr: _formatTime(chat.lastMessageTime),
                        onTap: () {
                          // Clear the unread badge locally on open.
                          context.read<ChatListBloc>().add(MarkChatRead(chat.id));
                          context.pushNamed(
                            AppRoutes.chatDetail,
                            pathParameters: {'id': chat.id},
                            extra: {
                              'name': title,
                              'interlocutorId': chat.interlocutorId ?? '',
                              'chatType': chat.type,
                            },
                          );
                        },
                        onLongPress: () => _showChatOptions(context, chat),
                      );
                    },
                  ),

                // Bottom padding for FAB
                const SliverToBoxAdapter(child: SizedBox(height: 80)),
              ],
            );
          },
        ),

        // ── FAB ──────────────────────────────────────────────
        floatingActionButton: Padding(
          padding: const EdgeInsets.only(bottom: 8),
          child: SizedBox(
            width: 56,
            height: 56,
            child: FloatingActionButton(
              onPressed: () => _showNewChatSheet(context),
              backgroundColor: theme.colorScheme.tertiary,
              foregroundColor: theme.scaffoldBackgroundColor,
              elevation: 8,
              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(18)),
              child: const Icon(Icons.edit_rounded, size: 22),
            ),
          ),
        ),
      ),
    );
  }

  void _showChatOptions(BuildContext context, ChatListItem chat) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);
    final l10n = context.l10n;
    final bloc = context.read<ChatListBloc>();
    final isOwner = chat.memberRole == 'OWNER';
    final isDirect = chat.type == 'DIRECT';

    showModalBottomSheet(
      context: context,
      backgroundColor: theme.colorScheme.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      ),
      builder: (ctx) => SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              margin: const EdgeInsets.only(top: 8, bottom: 8),
              width: 38,
              height: 4,
              decoration: BoxDecoration(
                color: palette.textSecondary.withValues(alpha: 0.4),
                borderRadius: BorderRadius.circular(2),
              ),
            ),
            // Pin / Unpin
            ListTile(
              leading: Icon(
                chat.isPinned ? Icons.push_pin_outlined : Icons.push_pin_rounded,
                color: palette.textAccent,
              ),
              title: Text(
                chat.isPinned ? l10n.chatUnpin : l10n.chatPin,
                style: TextStyle(fontFamily: 'Manrope', fontWeight: FontWeight.w600, color: theme.colorScheme.onSurface),
              ),
              onTap: () {
                Navigator.of(ctx).pop();
                if (chat.isPinned) {
                  bloc.add(UnpinChatEvent(chat.id));
                } else {
                  bloc.add(PinChatEvent(chat.id));
                }
              },
            ),
            // Mute / Unmute
            ListTile(
              leading: Icon(
                chat.isMuted ? Icons.notifications_rounded : Icons.notifications_off_rounded,
                color: palette.textAccent,
              ),
              title: Text(
                chat.isMuted ? l10n.chatUnmute : l10n.chatMute,
                style: TextStyle(fontFamily: 'Manrope', fontWeight: FontWeight.w600, color: theme.colorScheme.onSurface),
              ),
              onTap: () {
                Navigator.of(ctx).pop();
                if (chat.isMuted) {
                  bloc.add(UnmuteChatEvent(chat.id));
                } else {
                  bloc.add(MuteChatEvent(chat.id));
                }
              },
            ),
            // Delete options — Saved Messages cannot be deleted.
            if (chat.type != 'SAVED') ...[
              if (isDirect) ...[
                _DeleteTile(
                  label: l10n.chatDeleteForMe,
                  palette: palette,
                  theme: theme,
                  onTap: () {
                    Navigator.of(ctx).pop();
                    _confirmDelete(context, chat, scope: 'me', message: l10n.chatDeleteConfirmForMe);
                  },
                ),
                _DeleteTile(
                  label: l10n.chatDeleteForEveryone,
                  palette: palette,
                  theme: theme,
                  onTap: () {
                    Navigator.of(ctx).pop();
                    _confirmDelete(context, chat, scope: 'everyone', message: l10n.chatDeleteConfirmForAll);
                  },
                ),
              ] else if (isOwner) ...[
                _DeleteTile(
                  label: chat.type == 'CHANNEL' ? l10n.chatDeleteChannel : l10n.chatDeleteGroup,
                  palette: palette,
                  theme: theme,
                  onTap: () {
                    Navigator.of(ctx).pop();
                    _confirmDelete(context, chat, scope: 'everyone', message: l10n.chatDeleteConfirmForAll);
                  },
                ),
                _DeleteTile(
                  label: l10n.chatDeleteForMe,
                  palette: palette,
                  theme: theme,
                  onTap: () {
                    Navigator.of(ctx).pop();
                    _confirmDelete(context, chat, scope: 'me', message: l10n.chatDeleteConfirmForMe);
                  },
                ),
              ] else ...[
                _DeleteTile(
                  label: l10n.chatDeleteForMe,
                  palette: palette,
                  theme: theme,
                  onTap: () {
                    Navigator.of(ctx).pop();
                    _confirmDelete(context, chat, scope: 'me', message: l10n.chatDeleteConfirmForMe);
                  },
                ),
              ],
            ],
            const SizedBox(height: 8),
          ],
        ),
      ),
    );
  }

  void _confirmDelete(BuildContext context, ChatListItem chat, {required String scope, required String message}) {
    final theme = Theme.of(context);
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final l10n = context.l10n;
    final bloc = context.read<ChatListBloc>();

    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        backgroundColor: theme.colorScheme.surface,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(18)),
        content: Text(message, style: TextStyle(fontFamily: 'Manrope', color: theme.colorScheme.onSurface)),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: Text(l10n.cancel, style: TextStyle(fontFamily: 'Manrope', color: palette.textSecondary)),
          ),
          TextButton(
            onPressed: () {
              Navigator.pop(ctx);
              bloc.add(DeleteChatEvent(chat.id, scope: scope));
            },
            child: Text(l10n.confirm, style: TextStyle(fontFamily: 'Manrope', color: palette.danger, fontWeight: FontWeight.w700)),
          ),
        ],
      ),
    );
  }

}

// ── New chat sheet action ─────────────────────────────────────
class _NewChatAction extends StatelessWidget {
  final IconData icon;
  final String title;
  final VoidCallback onTap;

  const _NewChatAction({required this.icon, required this.title, required this.onTap});

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);

    return ListTile(
      leading: Container(
        width: 44,
        height: 44,
        decoration: BoxDecoration(
          color: theme.colorScheme.tertiary.withValues(alpha: 0.12),
          borderRadius: BorderRadius.circular(12),
        ),
        child: Icon(icon, color: theme.colorScheme.tertiary, size: 22),
      ),
      title: Text(
        title,
        style: TextStyle(
          fontFamily: 'Manrope',
          fontSize: 15,
          fontWeight: FontWeight.w600,
          color: theme.colorScheme.onSurface,
        ),
      ),
      trailing: Icon(Icons.chevron_right_rounded, color: palette.textSecondary),
      onTap: onTap,
    );
  }
}

// ── Single chat row ───────────────────────────────────────────
class _ChatRow extends StatelessWidget {
  final ChatListItem chat;
  final String title;
  final String timeStr;
  final VoidCallback onTap;
  final VoidCallback? onLongPress;

  const _ChatRow({
    required this.chat,
    required this.title,
    required this.timeStr,
    required this.onTap,
    this.onLongPress,
  });

  // Returns a human-readable last-message preview, handling SUDA_TRANSFER /
  // DONATION system messages so raw JSON never appears in the chat list (§2.3).
  static String _buildPreview(BuildContext context, ChatListItem chat) {
    final l10n = context.l10n;
    final msgType = chat.lastMessageType?.toUpperCase();

    // Media messages carry no text content — show a localized type label with
    // an icon hint instead of falling through to "No messages yet".
    switch (msgType) {
      case 'IMAGE':
        return '📷 ${l10n.messageTypeImage}';
      case 'VIDEO':
        return '🎥 ${l10n.messageTypeVideo}';
      case 'VOICE':
        return '🎤 ${l10n.messageTypeVoice}';
      case 'FILE':
        return '📎 ${l10n.messageTypeFile}';
    }

    if (chat.lastMessageContent.isEmpty) return l10n.chatNoMessagesYet;

    if (msgType == 'SUDA_TRANSFER' || msgType == 'DONATION') {
      try {
        final payload = jsonDecode(chat.lastMessageContent) as Map<String, dynamic>;
        final amountWei = payload['amount_wei'] as String? ?? '0';
        final amount = weiToSudaDisplay(amountWei);

        final fromDisplay = payload['from_display_name'] as String?;
        final fromAddr = payload['from_address'] as String?;
        final fromId = payload['from_user_id'] as String?;
        final from = (fromDisplay?.isNotEmpty == true)
            ? fromDisplay!
            : (fromAddr != null && fromAddr.length > 6)
                ? '${fromAddr.substring(0, 6)}…'
                : (fromId != null && fromId.length > 8)
                    ? '${fromId.substring(0, 8)}…'
                    : '?';

        if (msgType == 'SUDA_TRANSFER') {
          final toDisplay = payload['to_display_name'] as String?;
          final toAddr = payload['to_address'] as String?;
          final toId = payload['to_user_id'] as String?;
          final to = (toDisplay?.isNotEmpty == true)
              ? toDisplay!
              : (toAddr != null && toAddr.length > 6)
                  ? '${toAddr.substring(0, 6)}…'
                  : (toId != null && toId.length > 8)
                      ? '${toId.substring(0, 8)}…'
                      : '?';
          return l10n.chatPreviewTransfer(from, to, amount);
        } else {
          return l10n.chatPreviewDonation(from, amount);
        }
      } catch (_) {
        // Fallback: show generic label instead of raw JSON.
        return msgType == 'SUDA_TRANSFER'
            ? '💸 SUDA transfer'
            : '💝 Donation';
      }
    }

    if (msgType == 'SYSTEM') return l10n.sysGenericEvent;

    return chat.lastMessageContent;
  }

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);

    final initials = title.isNotEmpty ? title[0].toUpperCase() : '?';
    final preview = _buildPreview(context, chat);

    return InkWell(
      onTap: onTap,
      onLongPress: onLongPress,
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
        decoration: BoxDecoration(
          border: Border(
            bottom: BorderSide(
              color: theme.colorScheme.onSurface.withValues(alpha: 0.06),
            ),
          ),
        ),
        child: Row(
          children: [
            // Avatar — DIRECT: interlocutor's avatar media id (userId as fallback
            // for older data); GROUP/CHANNEL: the chat's avatar media id.
            SudaAvatar(
              mediaId: chat.type == 'DIRECT'
                  ? chat.interlocutorAvatarMediaId
                  : chat.avatarMediaId,
              userId: chat.type == 'DIRECT' ? chat.interlocutorId : null,
              initials: initials,
              size: 48,
            ),
            const SizedBox(width: 12),

            // Content
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Top row: name + kind pill + time
                  Row(
                    children: [
                      Expanded(
                        child: Row(
                          children: [
                            Flexible(
                              child: Text(
                                title,
                                style: TextStyle(
                                  fontFamily: 'Manrope',
                                  fontSize: 15,
                                  fontWeight: FontWeight.w600,
                                  color: theme.colorScheme.onSurface,
                                ),
                                overflow: TextOverflow.ellipsis,
                              ),
                            ),
                            if (chat.type == 'CHANNEL')
                              Padding(
                                padding: const EdgeInsets.only(left: 5),
                                child: _KindPill(icon: Icons.campaign_rounded),
                              ),
                            if (chat.type == 'GROUP')
                              Padding(
                                padding: const EdgeInsets.only(left: 5),
                                child: _KindPill(icon: Icons.group_rounded),
                              ),
                          ],
                        ),
                      ),
                      const SizedBox(width: 8),
                      if (chat.isPinned)
                        Padding(
                          padding: const EdgeInsets.only(right: 4),
                          child: Icon(Icons.push_pin_rounded, size: 13, color: palette.textSecondary),
                        ),
                      Text(
                        timeStr,
                        style: TextStyle(
                          fontFamily: 'Manrope',
                          fontSize: 12,
                          color: palette.textSecondary,
                        ),
                      ),
                    ],
                  ),

                  const SizedBox(height: 4),

                  // Bottom row: preview + badges
                  Row(
                    children: [
                      Expanded(
                        child: Text(
                          preview,
                          style: TextStyle(
                            fontFamily: 'Manrope',
                            fontSize: 13,
                            color: palette.textSecondary,
                          ),
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ),
                      const SizedBox(width: 8),
                      // Badges: mute, pin, unread
                      if (chat.unreadCount > 0)
                        _UnreadBadge(count: chat.unreadCount, muted: chat.isMuted),
                    ],
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

class _KindPill extends StatelessWidget {
  final IconData icon;
  const _KindPill({required this.icon});

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 2),
      decoration: BoxDecoration(
        color: palette.surface2,
        borderRadius: BorderRadius.circular(6),
      ),
      child: Icon(icon, size: 11, color: palette.textSecondary),
    );
  }
}

class _DeleteTile extends StatelessWidget {
  final String label;
  final AppPaletteTheme palette;
  final ThemeData theme;
  final VoidCallback onTap;

  const _DeleteTile({required this.label, required this.palette, required this.theme, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return ListTile(
      leading: Icon(Icons.delete_outline_rounded, color: palette.danger),
      title: Text(
        label,
        style: TextStyle(fontFamily: 'Manrope', fontWeight: FontWeight.w600, color: palette.danger),
      ),
      onTap: onTap,
    );
  }
}

class _UnreadBadge extends StatelessWidget {
  final int count;
  final bool muted;
  const _UnreadBadge({required this.count, required this.muted});

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);

    return Container(
      constraints: const BoxConstraints(minWidth: 20),
      height: 20,
      padding: const EdgeInsets.symmetric(horizontal: 5),
      decoration: BoxDecoration(
        color: muted ? palette.textSecondary : theme.colorScheme.tertiary,
        borderRadius: BorderRadius.circular(100),
      ),
      alignment: Alignment.center,
      child: Text(
        count > 99 ? '99+' : '$count',
        style: TextStyle(
          fontFamily: 'Manrope',
          fontSize: 11,
          fontWeight: FontWeight.w700,
          color: muted ? theme.colorScheme.surface : theme.scaffoldBackgroundColor,
        ),
      ),
    );
  }
}
