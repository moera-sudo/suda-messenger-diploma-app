import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import 'package:go_router/go_router.dart';
import 'package:scrollable_positioned_list/scrollable_positioned_list.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import 'package:messenger_app_v2/app/config/theme/app_theme.dart';
import 'package:messenger_app_v2/app/navigation/app_routes.dart';
import '../../../../features/channels/data/models/channel_view_model.dart';
import '../../../../features/channels/data/models/gating_rule.dart';
import '../../../../features/channels/domain/repositories/channel_repository.dart';
import '../../../../features/media/domain/repositories/media_repository.dart';
import '../../../../shared/data/logger/app_logger.dart';
import '../../../../shared/domain/models/app_failure.dart';
import '../../../../shared/presentation/widgets/suda_button.dart';
import '../../../../shared/presentation/feedback/app_feedback.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/placeholder/placeholder_page.dart';
import '../../data/models/chat_models.dart';
import '../../domain/repositories/chat_repository.dart';
import '../bloc/chat_detail_bloc.dart';
import '../bloc/chat_list_bloc.dart';
import '../widgets/chat_app_bar.dart';
import '../../../../features/channels/presentation/widgets/donate_sheet.dart';
import '../../../../features/user/domain/repositories/user_repository.dart';
import '../widgets/chat_input_field.dart';
import '../widgets/message_bubble.dart';
import '../widgets/message_readers_sheet.dart';
import '../widgets/transfer_sheet.dart';

class ChatDetailPage extends StatefulWidget {
  final String chatId;
  final String chatName;
  final String interlocutorId;
  final String chatType; // DIRECT, GROUP, CHANNEL, SAVED
  final String myRole; // OWNER, ADMIN, MEMBER

  const ChatDetailPage({
    super.key,
    required this.chatId,
    this.chatName = 'Unknown',
    this.interlocutorId = '',
    this.chatType = 'DIRECT',
    this.myRole = 'MEMBER',
  });

  @override
  State<ChatDetailPage> createState() => _ChatDetailPageState();
}

class _ChatDetailPageState extends State<ChatDetailPage> {
  Message? _editingMessage;

  // Channel-only: viewer relationship (drives input/banner/menu) and mute state.
  ChannelViewModel? _channelView;
  // Channel-only: token-gating / paid-subscription rule (drives the paywall + price).
  GatingRule? _gatingRule;
  // Optimistically flips the UI to "subscribed" right after a successful subscribe,
  // before the eventually-consistent channel-view read catches up.
  bool _subscribedOptimistic = false;
  bool _isMuted = false;

  // Indexed scrolling so search results can jump to a specific message.
  final ItemScrollController _itemScrollController = ItemScrollController();
  final ItemPositionsListener _itemPositionsListener =
      ItemPositionsListener.create();
  // Shows the "jump to latest" button once the user scrolls far up the (reversed)
  // list — index 0 is the newest message at the bottom.
  final ValueNotifier<bool> _showScrollToBottom = ValueNotifier(false);
  static const _scrollToBottomThreshold = 8;

  @override
  void initState() {
    super.initState();
    _itemPositionsListener.itemPositions.addListener(_onPositionsChanged);
    if (widget.chatType == 'CHANNEL') {
      _isMuted = _readMutedFromList();
      _loadChannelView();
      _loadGatingRule();
    }
  }

  @override
  void dispose() {
    _itemPositionsListener.itemPositions.removeListener(_onPositionsChanged);
    _showScrollToBottom.dispose();
    super.dispose();
  }

  /// Toggles the jump-to-latest button based on how far the user scrolled up.
  void _onPositionsChanged() {
    final positions = _itemPositionsListener.itemPositions.value;
    if (positions.isEmpty) return;
    // Reversed list → smallest index is nearest the bottom (newest message).
    final minIndex =
        positions.map((p) => p.index).reduce((a, b) => a < b ? a : b);
    final show = minIndex > _scrollToBottomThreshold;
    if (show != _showScrollToBottom.value) _showScrollToBottom.value = show;
  }

  void _scrollToBottom() {
    if (!_itemScrollController.isAttached) return;
    _itemScrollController.scrollTo(
      index: 0,
      duration: const Duration(milliseconds: 300),
      curve: Curves.easeOut,
    );
  }

  bool _readMutedFromList() {
    // ChatListBloc is scoped to the chats tab, not this pushed route — read it
    // defensively so opening a channel never crashes if it isn't an ancestor.
    try {
      final chats = context.read<ChatListBloc>().state.chats;
      for (final c in chats) {
        if (c.id == widget.chatId) return c.isMuted;
      }
    } catch (_) {
      // Provider not found above this route — fall back to "not muted".
    }
    return false;
  }

  Future<void> _loadChannelView() async {
    final result = await sl<ChannelRepository>().getChannelView(widget.chatId);
    if (!mounted) return;
    result.fold((_) {}, (view) => setState(() => _channelView = view));
  }

  Future<void> _loadGatingRule() async {
    final result = await sl<ChannelRepository>().getGatingRule(widget.chatId);
    if (!mounted) return;
    result.fold((_) {}, (rule) => setState(() => _gatingRule = rule));
  }

  /// Scrolls the list so the jumped-to search result is centred. Index maps
  /// 1:1 to the reversed message list used by the builder.
  void _scrollToMessage(int messageId, ChatDetailState state) {
    final index = state.messages.indexWhere((m) => m.id == messageId);
    if (index < 0 || !_itemScrollController.isAttached) return;
    _itemScrollController.scrollTo(
      index: index,
      alignment: 0.4,
      duration: const Duration(milliseconds: 400),
      curve: Curves.easeInOut,
    );
  }

  void _startEditing(Message message) => setState(() => _editingMessage = message);
  void _cancelEditing() => setState(() => _editingMessage = null);

  void _submitEdit(BuildContext context, String newContent) {
    final msg = _editingMessage;
    if (msg == null) return;
    final trimmed = newContent.trim();
    if (trimmed.isNotEmpty && trimmed != msg.content) {
      context.read<ChatDetailBloc>().add(EditMessageEvent(msg.id, trimmed));
    }
    _cancelEditing();
  }

  // ── Channel actions (overflow menu / subscribe banner) ──────────

  /// Re-fetch the channel view and reload posts after a membership change.
  Future<void> _reloadChannel(BuildContext blocContext) async {
    blocContext.read<ChatDetailBloc>().add(LoadMessages());
    await _loadChannelView();
    await _loadGatingRule();
  }

  /// Optimistically flip to "subscribed" so the bottom banner hides immediately,
  /// then sync the real state (channel-view is eventually consistent).
  void _afterSubscribed() {
    if (mounted) setState(() => _subscribedOptimistic = true);
    _reloadChannel(context);
  }

  /// Pay the subscription price and join. Shared by the opening paywall dialog
  /// and the bottom paywall banner. Backend charges only for paid channels.
  Future<void> _onPaySubscribe({BuildContext? dialogCtx}) async {
    final result = await sl<ChannelRepository>().subscribe(widget.chatId);
    if (!mounted) return;
    result.fold(
      (f) {
        final code = f is ServerFailure ? f.code : null;
        if (code == 'INSUFFICIENT_BALANCE') {
          AppFeedback.showError(context.l10n.channelInsufficientBalance);
        } else if (code == 'NO_WALLET') {
          AppFeedback.showError(context.l10n.channelNoWallet);
        } else {
          AppFeedback.showError(context.l10n.channelSubscribeFailed);
        }
        // Keep the dialog open so the user can retry or press Back.
      },
      (_) {
        if (dialogCtx != null && Navigator.canPop(dialogCtx)) {
          Navigator.pop(dialogCtx);
        }
        AppFeedback.showSuccess(context.l10n.channelSubscribeSuccess);
        _afterSubscribed(); // hide paywall now + sync view/messages
      },
    );
  }

  Future<void> _onSubscribe(BuildContext blocContext) async {
    final result = await sl<ChannelRepository>().subscribe(widget.chatId);
    if (!mounted) return;
    result.fold(
      (f) => AppFeedback.showError(f.message),
      (_) {
        AppFeedback.showSuccess(blocContext.l10n.channelSubscribeSuccess);
        _afterSubscribed();
      },
    );
  }

  Future<void> _onUnsubscribe(BuildContext blocContext) async {
    final result = await sl<ChannelRepository>().unsubscribe(widget.chatId);
    if (!mounted) return;
    result.fold(
      (f) => AppFeedback.showError(f.message),
      (_) {
        AppFeedback.showSuccess(blocContext.l10n.channelUnsubscribeSuccess);
        _reloadChannel(blocContext);
      },
    );
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

  /// Opens the comments screen for a channel post.
  void _openComments(BuildContext context, Message post, ChatDetailState state) {
    final view = _channelView;
    context.pushNamed(
      AppRoutes.channelComments,
      pathParameters: {'id': widget.chatId, 'postId': post.id.toString()},
      extra: {
        'post': post,
        'commentsEnabled': view?.commentsEnabled ?? false,
        'isSubscriber': view?.isMember ?? false,
        'isAdmin': view?.isAdmin ?? false,
        'currentUserId': state.currentUserId ?? '',
      },
    );
  }

  static bool _isSystemMessage(String type) {
    final t = type.toUpperCase();
    return t == 'SYSTEM' || t == 'SUDA_TRANSFER' || t == 'DONATION';
  }

  void _onMessageLongPress(BuildContext context, Message message, bool isMe) {
    final l10n = context.l10n;
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final bloc = context.read<ChatDetailBloc>();

    // Reply only where the user can actually post: any non-channel chat, or a
    // channel where the viewer is OWNER/ADMIN. Read-only subscribers get no Reply.
    final canReply = !_isSystemMessage(message.type) &&
        (widget.chatType != 'CHANNEL' || (_channelView?.isAdmin ?? false));

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
            // Handle pill
            Container(
              margin: const EdgeInsets.only(top: 8, bottom: 8),
              width: 38,
              height: 4,
              decoration: BoxDecoration(
                color: palette.textSecondary.withValues(alpha: 0.4),
                borderRadius: BorderRadius.circular(2),
              ),
            ),

            // Reply — hidden for system messages and read-only channel subscribers
            if (canReply)
              ListTile(
                leading: Icon(Icons.reply_rounded, color: theme.colorScheme.tertiary),
                title: Text(l10n.messageReply, style: TextStyle(color: theme.colorScheme.onSurface)),
                onTap: () {
                  Navigator.pop(ctx);
                  bloc.add(SetReplyTarget(message));
                },
              ),

            // Forward — hidden for system messages
            if (!_isSystemMessage(message.type))
              ListTile(
                leading: Icon(Icons.forward_rounded, color: theme.colorScheme.onSurface),
                title: Text(l10n.messageForward, style: TextStyle(color: theme.colorScheme.onSurface)),
                onTap: () {
                  Navigator.pop(ctx);
                  _showForwardSheet(context, message);
                },
              ),

            // Copy text
            ListTile(
              leading: Icon(Icons.copy_outlined, color: theme.colorScheme.onSurface),
              title: Text(l10n.messageActionCopy, style: TextStyle(color: theme.colorScheme.onSurface)),
              onTap: () {
                Navigator.pop(ctx);
                Clipboard.setData(ClipboardData(text: message.content));
              },
            ),

            // Edit — own text messages only
            if (isMe && message.type.toUpperCase() == 'TEXT')
              ListTile(
                leading: Icon(Icons.edit_outlined, color: theme.colorScheme.onSurface),
                title: Text(l10n.messageActionEdit, style: TextStyle(color: theme.colorScheme.onSurface)),
                onTap: () {
                  Navigator.pop(ctx);
                  _startEditing(message);
                },
              ),

            // Delete for me
            ListTile(
              leading: Icon(Icons.delete_outline, color: palette.danger),
              title: Text(l10n.messageActionDeleteForMe, style: TextStyle(color: palette.danger)),
              onTap: () {
                Navigator.pop(ctx);
                bloc.add(DeleteMessageEvent(message.id, forEveryone: false));
              },
            ),

            // Delete for everyone — own messages only
            if (isMe)
              ListTile(
                leading: Icon(Icons.delete_sweep_outlined, color: palette.danger),
                title: Text(l10n.messageActionDeleteForEveryone, style: TextStyle(color: palette.danger)),
                onTap: () {
                  Navigator.pop(ctx);
                  bloc.add(DeleteMessageEvent(message.id, forEveryone: true));
                },
              ),

            const SizedBox(height: 8),
          ],
        ),
      ),
    );
  }

  /// Only chats the user can post to are valid forward targets (forward = repost):
  /// DIRECT/GROUP/SAVED always, channels only when OWNER/ADMIN (read-only
  /// subscriptions can't receive posts).
  static bool _canForwardTo(ChatListItem chat) {
    switch (chat.type) {
      case 'DIRECT':
      case 'GROUP':
      case 'SAVED':
        return true;
      case 'CHANNEL':
        return chat.memberRole == 'OWNER' || chat.memberRole == 'ADMIN';
      default:
        return false;
    }
  }

  /// Loads writable chats for the forward sheet. ChatListBloc is scoped to the
  /// chats tab and isn't an ancestor of this pushed route, so we fetch fresh
  /// from the repository instead of reading the bloc.
  Future<List<ChatListItem>> _loadForwardTargets() async {
    final result = await sl<ChatRepository>().getChats();
    return result.fold(
      (f) {
        sl<AppLogger>().error('Forward: failed to load chats', f);
        return <ChatListItem>[];
      },
      (chats) => chats.where(_canForwardTo).toList(),
    );
  }

  /// Display title for a forward-target chat — mirrors the chat list logic.
  String _forwardTargetTitle(BuildContext context, ChatListItem chat) {
    final l10n = context.l10n;
    return switch (chat.type) {
      'SAVED' => l10n.chatSavedMessages,
      'GROUP' || 'CHANNEL' => chat.name ?? l10n.chatUnknownUser,
      _ => chat.interlocutorName ?? l10n.chatUnknownUser,
    };
  }

  /// Forwards [message] to [chat], then opens that chat on success so the user
  /// lands where the message was sent (Telegram-style). Awaiting the forward
  /// guarantees the message is persisted before the destination chat loads.
  Future<void> _forwardTo(
    BuildContext sheetContext,
    ChatListItem chat,
    String title,
    Message message,
  ) async {
    // Capture the router before the async gap — context may be defunct after.
    final router = GoRouter.of(context);
    Navigator.pop(sheetContext); // close the forward sheet

    final result = await sl<ChatRepository>().forwardMessage(
      fromChatId: widget.chatId,
      toChatId: chat.id,
      messageId: message.id,
    );
    if (!mounted) return;
    result.fold(
      (f) {
        sl<AppLogger>().error('Forward to ${chat.id} failed', f);
        AppFeedback.showError(f.message);
      },
      (_) {
        sl<AppLogger>().info('Message ${message.id} forwarded to ${chat.id} — opening chat');
        router.pushNamed(
          AppRoutes.chatDetail,
          pathParameters: {'id': chat.id},
          extra: {
            'name': title,
            'interlocutorId': chat.interlocutorId ?? '',
            'chatType': chat.type,
          },
        );
      },
    );
  }

  void _showForwardSheet(BuildContext context, Message message) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;

    showModalBottomSheet(
      context: context,
      backgroundColor: theme.colorScheme.surface,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      ),
      builder: (ctx) {
        return DraggableScrollableSheet(
          initialChildSize: 0.6,
          maxChildSize: 0.9,
          minChildSize: 0.4,
          expand: false,
          builder: (_, scrollController) => Column(
            children: [
              Container(
                margin: const EdgeInsets.only(top: 8, bottom: 12),
                width: 38,
                height: 4,
                decoration: BoxDecoration(
                  color: palette.textSecondary.withValues(alpha: 0.4),
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
              Padding(
                padding: const EdgeInsets.only(left: 16, bottom: 12),
                child: Align(
                  alignment: Alignment.centerLeft,
                  child: Text(
                    context.l10n.forwardTo,
                    style: TextStyle(
                      fontFamily: 'Manrope',
                      fontSize: 17,
                      fontWeight: FontWeight.w700,
                      color: theme.colorScheme.onSurface,
                    ),
                  ),
                ),
              ),
              Expanded(
                child: FutureBuilder<List<ChatListItem>>(
                  future: _loadForwardTargets(),
                  builder: (fctx, snap) {
                    if (snap.connectionState == ConnectionState.waiting) {
                      return const Center(child: CircularProgressIndicator());
                    }
                    final chats = snap.data ?? const <ChatListItem>[];
                    if (chats.isEmpty) {
                      return Center(
                        child: Text(
                          fctx.l10n.forwardNoTargets,
                          style: TextStyle(
                            fontFamily: 'Manrope',
                            color: palette.textSecondary,
                          ),
                        ),
                      );
                    }
                    return ListView.builder(
                      controller: scrollController,
                      itemCount: chats.length,
                      itemBuilder: (_, i) {
                        final chat = chats[i];
                        final title = _forwardTargetTitle(fctx, chat);
                        return ListTile(
                          leading: CircleAvatar(
                            radius: 20,
                            backgroundColor: theme.colorScheme.primary,
                            child: Text(
                              title.isNotEmpty ? title[0].toUpperCase() : '?',
                              style: const TextStyle(color: Colors.white, fontFamily: 'SpaceGrotesk'),
                            ),
                          ),
                          title: Text(title, style: TextStyle(color: theme.colorScheme.onSurface)),
                          onTap: () => _forwardTo(ctx, chat, title, message),
                        );
                      },
                    );
                  },
                ),
              ),
            ],
          ),
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (context) => sl<ChatDetailBloc>()
        ..init(widget.chatId, interlocutorId: widget.interlocutorId, chatType: widget.chatType),
      child: Builder(
        builder: (context) {
          final theme = Theme.of(context);
          final palette = theme.extension<AppPaletteTheme>()!;

          return BlocListener<ChatDetailBloc, ChatDetailState>(
            listenWhen: (p, c) => !p.isGatingRequired && c.isGatingRequired,
            listener: (ctx, state) => _showGatingAlert(ctx),
            child: Scaffold(
            backgroundColor: theme.scaffoldBackgroundColor,
            appBar: ChatAppBar(
              title: widget.chatName,
              chatId: widget.chatId,
              interlocutorId: widget.interlocutorId,
              chatType: widget.chatType,
              myRole: widget.myRole,
              channelView: _channelView,
              isMuted: _isMuted,
              onSubscribe: () => _onSubscribe(context),
              onUnsubscribe: () => _onUnsubscribe(context),
              onDonate: () => DonateSheet.show(context, channelId: widget.chatId, chatId: widget.chatId),
              onInfoReturn: _loadChannelView,
              onToggleMute: _onToggleMute,
            ),
            body: Column(
              children: [
                // Messages list with pagination
                Expanded(
                  child: Stack(
                    children: [
                      BlocListener<ChatDetailBloc, ChatDetailState>(
                    listenWhen: (p, c) =>
                        p.highlightedMessageId != c.highlightedMessageId &&
                        c.highlightedMessageId != null,
                    listener: (ctx, state) =>
                        _scrollToMessage(state.highlightedMessageId!, state),
                    child: BlocBuilder<ChatDetailBloc, ChatDetailState>(
                    buildWhen: (p, c) =>
                        p.messages != c.messages ||
                        p.status != c.status ||
                        p.isLoadingMore != c.isLoadingMore ||
                        p.hasMoreMessages != c.hasMoreMessages ||
                        p.highlightedMessageId != c.highlightedMessageId,
                    builder: (context, state) {
                      if (state.status == ChatStatus.loading) {
                        return const Center(child: CircularProgressIndicator());
                      }
                      if (state.status == ChatStatus.failure) {
                        return const PlaceholderPage(type: PlaceholderType.error, scaffold: false);
                      }
                      if (state.messages.isEmpty && state.status == ChatStatus.success) {
                        return Center(
                          child: Text(
                            context.l10n.chatNoMessagesYet,
                            style: TextStyle(color: palette.textSecondary),
                          ),
                        );
                      }

                      final isGroup = widget.chatType == 'GROUP';

                      return NotificationListener<ScrollNotification>(
                        onNotification: (n) {
                          // Load more when near the bottom of reversed list (= older messages)
                          if (n is ScrollEndNotification &&
                              n.metrics.pixels >= n.metrics.maxScrollExtent - 300 &&
                              state.hasMoreMessages &&
                              !state.isLoadingMore &&
                              state.messages.isNotEmpty) {
                            context.read<ChatDetailBloc>().add(
                                  LoadMoreMessages(state.messages.last.id),
                                );
                          }
                          return false;
                        },
                        child: ScrollablePositionedList.builder(
                          reverse: true,
                          itemScrollController: _itemScrollController,
                          itemPositionsListener: _itemPositionsListener,
                          padding: const EdgeInsets.symmetric(vertical: 10),
                          itemCount: state.messages.length + (state.isLoadingMore ? 1 : 0),
                          itemBuilder: (context, index) {
                            // Loading indicator at bottom (top visually, since reversed)
                            if (index == state.messages.length) {
                              return const Padding(
                                padding: EdgeInsets.all(16),
                                child: Center(child: CircularProgressIndicator()),
                              );
                            }
                            final message = state.messages[index];
                            final isMe = message.senderId == state.currentUserId;
                            final sid = message.senderId;
                            final isSystem = message.type == 'SYSTEM' ||
                                message.type == 'SUDA_TRANSFER' ||
                                message.type == 'DONATION';

                            // Group sender display (Telegram-style runs): the list
                            // is reversed, so index-1 is visually below (newer),
                            // index+1 is visually above (older).
                            final showSender = isGroup && !isMe && !isSystem && sid != null;
                            final newer = index > 0 ? state.messages[index - 1] : null;
                            final older = index < state.messages.length - 1
                                ? state.messages[index + 1]
                                : null;
                            final showAvatar = showSender &&
                                (newer == null || newer.senderId != sid);
                            final showName = showSender &&
                                (older == null || older.senderId != sid);
                            final senderName =
                                showSender ? state.members[sid]?.displayName : null;

                            // Channel post comments affordance.
                            final isChannelPost = widget.chatType == 'CHANNEL' && !isSystem;
                            final commentsEnabled = _channelView?.commentsEnabled == true;

                            return RepaintBoundary(
                              child: MessageBubble(
                                message: message,
                                isMe: isMe,
                                mediaRepo: sl<MediaRepository>(),
                                onLongPress: () => _onMessageLongPress(context, message, isMe),
                                onReadReceiptTap: (isMe && isGroup && message.status == 'READ')
                                    ? () => MessageReadersSheet.show(context, message.id)
                                    : null,
                                isGroupIncoming: showSender,
                                senderName: senderName,
                                senderId: sid,
                                senderAvatarUrl: showAvatar ? state.memberAvatarUrls[sid] : null,
                                showAvatar: showAvatar,
                                showName: showName,
                                showComments: isChannelPost && commentsEnabled,
                                commentCount: message.commentCount,
                                onOpenComments: () => _openComments(context, message, state),
                                isHighlighted: message.id == state.highlightedMessageId,
                              ),
                            );
                          },
                        ),
                      );
                    },
                  ),
                  ),
                      // Jump-to-latest floating button (long chats only).
                      Positioned(
                        right: 12,
                        bottom: 12,
                        child: ValueListenableBuilder<bool>(
                          valueListenable: _showScrollToBottom,
                          builder: (context, show, _) => AnimatedScale(
                            scale: show ? 1 : 0,
                            duration: const Duration(milliseconds: 150),
                            child: _ScrollToBottomButton(onTap: _scrollToBottom),
                          ),
                        ),
                      ),
                    ],
                  ),
                ),

                // Typing indicator
                BlocBuilder<ChatDetailBloc, ChatDetailState>(
                  buildWhen: (p, c) => p.isInterlocutorTyping != c.isInterlocutorTyping,
                  builder: (context, state) {
                    if (!state.isInterlocutorTyping) return const SizedBox.shrink();
                    return Padding(
                      padding: const EdgeInsets.only(left: 16, bottom: 4),
                      child: Align(
                        alignment: Alignment.centerLeft,
                        child: Text(
                          context.l10n.chatStatusTyping,
                          style: TextStyle(
                            fontFamily: 'Manrope',
                            color: theme.colorScheme.tertiary,
                            fontSize: 12,
                            fontStyle: FontStyle.italic,
                          ),
                        ),
                      ),
                    );
                  },
                ),

                // Edit bar — shown when editing a message
                if (_editingMessage != null)
                  _EditBar(
                    message: _editingMessage!,
                    onCancel: _cancelEditing,
                  ),

                // Input field or blocked banner
                BlocBuilder<ChatDetailBloc, ChatDetailState>(
                  buildWhen: (p, c) =>
                      p.replyTarget != c.replyTarget ||
                      p.blockedMe != c.blockedMe ||
                      p.blockedByMe != c.blockedByMe,
                  builder: (context, state) {
                    // Channel: bottom area is driven by the viewer's relationship
                    // (from GET /channels/{id}/view), not by myRole or block state.
                    if (widget.chatType == 'CHANNEL') {
                      final view = _channelView;
                      if (view == null) return const SizedBox.shrink(); // still loading
                      // OWNER/ADMIN → fall through to the post composer below.
                      if (!view.isAdmin) {
                        // Subscriber (or just-subscribed): read-only, no input/buttons.
                        if (view.isMember || _subscribedOptimistic) {
                          return const SizedBox.shrink();
                        }
                        final price = _gatingRule?.subscriptionPrice ?? 0;
                        // Public paid channel → single paywall (price + Pay).
                        if (view.isPublic && price > 0) {
                          return _ChannelPaywallBanner(
                            price: price,
                            onPay: () => _onPaySubscribe(),
                          );
                        }
                        // Free public channel → plain subscribe bar.
                        if (view.isPublic) {
                          return _ChannelSubscribeBanner(
                            chatId: widget.chatId,
                            onSubscribed: _afterSubscribed,
                          );
                        }
                        // Private channel, not a member → request-to-join banner.
                        return _ChannelJoinBanner(chatId: widget.chatId);
                      }
                    } else if (state.blockedMe) {
                      return const _BlockedMeBanner();
                    } else if (state.blockedByMe) {
                      return _BlockedByMeBanner(
                        interlocutorId: widget.interlocutorId,
                        onUnblocked: () {
                          context.read<ChatDetailBloc>().add(LoadMessages());
                        },
                      );
                    }
                    // Resolve human-readable reply author name.
                    final replyTarget = state.replyTarget;
                    String? replyAuthorName;
                    if (replyTarget != null) {
                      final sid = replyTarget.senderId;
                      if (sid == state.currentUserId) {
                        replyAuthorName = context.l10n.messageSentByMe;
                      } else if (widget.chatType == 'GROUP' && sid != null) {
                        replyAuthorName = state.members[sid]?.displayName ?? sid;
                      } else {
                        replyAuthorName = widget.chatName;
                      }
                    }
                    return ChatInputField(
                      chatType: widget.chatType,
                      initialText: _editingMessage?.content,
                      replyToAuthor: replyAuthorName,
                      replyToText: replyTarget?.content,
                      onClearReply: () => context.read<ChatDetailBloc>().add(ClearReplyTarget()),
                      onSend: (text) {
                        if (_editingMessage != null) {
                          _submitEdit(context, text);
                        } else {
                          context.read<ChatDetailBloc>().add(SendMessage(text));
                        }
                      },
                      onTyping: (isTyping) {
                        context.read<ChatDetailBloc>().add(OnTypingInput(isTyping));
                      },
                      onAttachMedia: (path, type) {
                        if (type == 'TRANSFER') {
                          _showTransferSheet(context);
                          return;
                        }
                        context.read<ChatDetailBloc>().add(
                          SendMediaMessage(filePath: path, messageType: type),
                        );
                      },
                      onVoiceSend: (path, durationSec) {
                        context.read<ChatDetailBloc>().add(
                          SendMediaMessage(
                            filePath: path,
                            messageType: 'VOICE',
                            content: '$durationSec',
                          ),
                        );
                      },
                    );
                  },
                ),
              ],
            ),
          )); // Scaffold + BlocListener
        },
      ),
    );
  }

  // ── Gating alert ──────────────────────────────────────────────

  Future<void> _showGatingAlert(BuildContext ctx) async {
    // Capture context-dependent values BEFORE the async call to satisfy the linter.
    final l10n    = ctx.l10n;
    final theme   = Theme.of(ctx);

    // Ensure the rule (and its price) is loaded before showing the paywall.
    if (_gatingRule == null) await _loadGatingRule();
    if (!mounted) return;

    final price = _gatingRule?.subscriptionPrice ?? 0;
    final amountStr = price == price.truncateToDouble()
        ? price.toInt().toString()
        : price.toStringAsFixed(2);

    // Use State.context after mounted guard — ctx may be stale across the async gap.
    await showDialog<void>(
      context: context,
      barrierDismissible: false,
      builder: (dCtx) => AlertDialog(
        backgroundColor: theme.colorScheme.surface,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(18)),
        title: Text(l10n.gatingAlertTitle,
            style: TextStyle(fontFamily: 'Manrope', fontWeight: FontWeight.w700,
                color: theme.colorScheme.onSurface)),
        content: Text(l10n.gatingAlertBody(amountStr),
            style: TextStyle(fontFamily: 'Manrope', color: theme.colorScheme.onSurface)),
        actions: [
          TextButton(
            onPressed: () {
              Navigator.pop(dCtx);
              context.pop(); // leave the channel
            },
            child: Text(l10n.gatingBack,
                style: const TextStyle(fontFamily: 'Manrope')),
          ),
          TextButton(
            onPressed: () => _onPaySubscribe(dialogCtx: dCtx),
            child: Text(l10n.gatingPay,
                style: TextStyle(fontFamily: 'Manrope',
                    color: theme.colorScheme.tertiary, fontWeight: FontWeight.w700)),
          ),
        ],
      ),
    );
  }

  // ── Transfer sheet ────────────────────────────────────────────

  Future<void> _showTransferSheet(BuildContext ctx) async {
    // Capture bloc state before async gap.
    final state  = ctx.read<ChatDetailBloc>().state;
    final chatId  = widget.chatId;
    final chatType = widget.chatType;
    final interlocutorId = widget.interlocutorId;
    String? interlocutorUsername;

    // For DIRECT chats resolve the interlocutor's username from their profile.
    if (chatType == 'DIRECT' && interlocutorId.isNotEmpty) {
      final result = await sl<UserRepository>().getUserProfile(interlocutorId);
      if (!mounted) return;
      interlocutorUsername = result.fold((_) => null, (p) {
        final u = p.username;
        return u.startsWith('@') ? u.substring(1) : u;
      });
    }

    if (!mounted) return;
    await TransferSheet.show(
      context, // use State.context after mounted guard
      chatId: chatId,
      chatType: chatType,
      interlocutorUsername: interlocutorUsername,
      members: state.members,
      currentUserId: state.currentUserId ?? '',
    );
  }
}

/// Floating round button that jumps the chat back to the newest message.
class _ScrollToBottomButton extends StatelessWidget {
  final VoidCallback onTap;
  const _ScrollToBottomButton({required this.onTap});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    return Material(
      color: theme.colorScheme.surface,
      elevation: 3,
      shape: const CircleBorder(),
      child: InkWell(
        customBorder: const CircleBorder(),
        onTap: onTap,
        child: Container(
          width: 44,
          height: 44,
          decoration: BoxDecoration(
            shape: BoxShape.circle,
            border: Border.all(color: palette.divider),
          ),
          child: Icon(
            Icons.keyboard_arrow_down_rounded,
            color: theme.colorScheme.onSurface,
          ),
        ),
      ),
    );
  }
}

/// Banner shown when the interlocutor has blocked the current user.
class _BlockedMeBanner extends StatelessWidget {
  const _BlockedMeBanner();

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;

    return SafeArea(
      top: false,
      child: Container(
      height: 56,
      padding: const EdgeInsets.symmetric(horizontal: 16),
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        border: Border(top: BorderSide(color: palette.divider)),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(Icons.block_rounded, size: 18, color: palette.danger),
          const SizedBox(width: 8),
          Text(
            context.l10n.chatBlockedMeBanner,
            style: TextStyle(
              fontFamily: 'Manrope',
              fontSize: 13,
              fontWeight: FontWeight.w500,
              color: palette.textSecondary,
            ),
          ),
        ],
      ),
      ),
    );
  }
}

/// Banner shown when the current user has blocked the interlocutor.
/// Provides an unblock button to restore the conversation.
class _BlockedByMeBanner extends StatefulWidget {
  final String interlocutorId;
  final VoidCallback onUnblocked;

  const _BlockedByMeBanner({
    required this.interlocutorId,
    required this.onUnblocked,
  });

  @override
  State<_BlockedByMeBanner> createState() => _BlockedByMeBannerState();
}

class _BlockedByMeBannerState extends State<_BlockedByMeBanner> {
  bool _loading = false;

  Future<void> _unblock() async {
    setState(() => _loading = true);
    final result = await sl<UserRepository>().unblockUser(widget.interlocutorId);
    if (!mounted) return;
    result.fold(
      (f) {
        setState(() => _loading = false);
        AppFeedback.showError(f.message);
      },
      (_) {
        AppFeedback.showSuccess(context.l10n.unblockAction);
        widget.onUnblocked();
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;

    return SafeArea(
      top: false,
      child: Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        border: Border(top: BorderSide(color: palette.divider)),
      ),
      child: Row(
        children: [
          Icon(Icons.block_rounded, size: 18, color: palette.danger),
          const SizedBox(width: 8),
          Expanded(
            child: Text(
              context.l10n.chatBlockedByMeBanner,
              style: TextStyle(
                fontFamily: 'Manrope',
                fontSize: 13,
                fontWeight: FontWeight.w500,
                color: palette.textSecondary,
              ),
            ),
          ),
          const SizedBox(width: 12),
          SudaButton(
            label: context.l10n.unblockAction,
            variant: SudaButtonVariant.dangerOutline,
            size: SudaButtonSize.sm,
            loading: _loading,
            onPressed: _loading ? null : _unblock,
          ),
        ],
      ),
      ),
    );
  }
}

/// Paywall bar for a public PAID channel non-subscriber: shows the price and a
/// Pay button. Subscription (and the charge) is handled by [onPay].
class _ChannelPaywallBanner extends StatelessWidget {
  final double price;
  final VoidCallback onPay;
  const _ChannelPaywallBanner({required this.price, required this.onPay});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;
    final priceStr = price == price.truncateToDouble()
        ? price.toInt().toString()
        : price.toStringAsFixed(2);

    return SafeArea(
      top: false,
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
        decoration: BoxDecoration(
          color: theme.colorScheme.surface,
          border: Border(top: BorderSide(color: palette.divider)),
        ),
        child: Row(children: [
          Expanded(
            child: Text(
              l10n.channelTokenGatedDesc(priceStr),
              style: TextStyle(fontFamily: 'Manrope', fontSize: 13, color: palette.textSecondary),
            ),
          ),
          const SizedBox(width: 12),
          SudaButton(
            label: l10n.channelSubscribeForPrice(priceStr),
            variant: SudaButtonVariant.primary,
            size: SudaButtonSize.sm,
            onPressed: onPay,
          ),
        ]),
      ),
    );
  }
}

/// Shown for public channel visitors (can read, but not subscribed).
class _ChannelSubscribeBanner extends StatefulWidget {
  final String chatId;
  final VoidCallback? onSubscribed;
  const _ChannelSubscribeBanner({required this.chatId, this.onSubscribed});

  @override
  State<_ChannelSubscribeBanner> createState() => _ChannelSubscribeBannerState();
}

class _ChannelSubscribeBannerState extends State<_ChannelSubscribeBanner> {
  bool _loading = false;

  Future<void> _subscribe() async {
    setState(() => _loading = true);
    final result = await sl<ChannelRepository>().subscribe(widget.chatId);
    if (!mounted) return;
    result.fold(
      (f) {
        setState(() => _loading = false);
        AppFeedback.showError(f.message);
      },
      (_) {
        AppFeedback.showSuccess(context.l10n.channelSubscribeSuccess);
        widget.onSubscribed?.call();
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;

    return SafeArea(
      top: false,
      child: Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        border: Border(top: BorderSide(color: palette.divider)),
      ),
      child: Row(children: [
        Expanded(
          child: Text(
            context.l10n.channelSubscribeBanner,
            style: TextStyle(fontFamily: 'Manrope', fontSize: 13, color: palette.textSecondary),
          ),
        ),
        const SizedBox(width: 12),
        SudaButton(
          label: context.l10n.channelSubscribe,
          variant: SudaButtonVariant.primary,
          size: SudaButtonSize.sm,
          loading: _loading,
          onPressed: _loading ? null : _subscribe,
        ),
      ]),
      ),
    );
  }
}

/// Shown for private channel non-members (messages blocked, join-request needed).
class _ChannelJoinBanner extends StatefulWidget {
  final String chatId;
  const _ChannelJoinBanner({required this.chatId});

  @override
  State<_ChannelJoinBanner> createState() => _ChannelJoinBannerState();
}

class _ChannelJoinBannerState extends State<_ChannelJoinBanner> {
  bool _requested = false;
  bool _loading = false;

  Future<void> _sendRequest() async {
    setState(() => _loading = true);
    final result = await sl<ChannelRepository>().sendJoinRequest(widget.chatId);
    if (!mounted) return;
    result.fold(
      (f) {
        setState(() => _loading = false);
        AppFeedback.showError(f.message);
      },
      (_) => setState(() { _loading = false; _requested = true; }),
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    return SafeArea(
      top: false,
      child: Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        border: Border(top: BorderSide(color: palette.divider)),
      ),
      child: Row(children: [
        Icon(Icons.lock_outline_rounded, size: 18, color: palette.textSecondary),
        const SizedBox(width: 8),
        Expanded(
          child: Text(
            _requested ? l10n.channelRequestSent : l10n.channelPrivateBanner,
            style: TextStyle(fontFamily: 'Manrope', fontSize: 13, color: palette.textSecondary),
          ),
        ),
        if (!_requested) ...[
          const SizedBox(width: 12),
          SudaButton(
            label: l10n.channelJoinRequest,
            variant: SudaButtonVariant.primary,
            size: SudaButtonSize.sm,
            loading: _loading,
            onPressed: _loading ? null : _sendRequest,
          ),
        ],
      ]),
      ),
    );
  }
}

class _EditBar extends StatelessWidget {
  final Message message;
  final VoidCallback onCancel;

  const _EditBar({required this.message, required this.onCancel});

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 8),
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        border: Border(top: BorderSide(color: palette.divider)),
      ),
      child: Row(
        children: [
          Icon(Icons.edit_rounded, size: 18, color: theme.colorScheme.tertiary),
          const SizedBox(width: 10),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisSize: MainAxisSize.min,
              children: [
                Text(
                  context.l10n.messageEditingLabel,
                  style: TextStyle(
                    fontFamily: 'Manrope',
                    color: theme.colorScheme.tertiary,
                    fontSize: 12,
                    fontWeight: FontWeight.w700,
                  ),
                ),
                Text(
                  message.content,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                  style: TextStyle(
                    fontFamily: 'Manrope',
                    color: palette.textSecondary,
                    fontSize: 13,
                  ),
                ),
              ],
            ),
          ),
          IconButton(
            icon: Icon(Icons.close_rounded, color: palette.textSecondary),
            onPressed: onCancel,
            tooltip: context.l10n.messageEditCancel,
            padding: EdgeInsets.zero,
            constraints: const BoxConstraints(minWidth: 32, minHeight: 32),
          ),
        ],
      ),
    );
  }
}
