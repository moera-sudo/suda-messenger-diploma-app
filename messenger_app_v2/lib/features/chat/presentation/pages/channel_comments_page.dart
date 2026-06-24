import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:go_router/go_router.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import '../../../../app/config/theme/app_theme.dart';
import '../../../../features/channels/data/models/channel_comment_model.dart';
import '../../../../features/channels/domain/repositories/channel_repository.dart';
import '../../../../features/media/domain/repositories/media_repository.dart';
import '../../../../features/user/domain/repositories/user_repository.dart';
import '../../../../shared/data/api/socket_client.dart';
import '../../../../shared/domain/models/socket_event.dart';
import '../../../../shared/presentation/feedback/app_feedback.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/widgets/suda_avatar.dart';
import '../../data/models/chat_models.dart';

/// Comments under a single channel post (Channels-Groups-Frontend §7).
class ChannelCommentsPage extends StatefulWidget {
  final String channelId;
  final int postId;
  final Message? post;
  final bool commentsEnabled;
  final bool isSubscriber;
  final bool isAdmin;
  final String currentUserId;

  const ChannelCommentsPage({
    super.key,
    required this.channelId,
    required this.postId,
    this.post,
    this.commentsEnabled = false,
    this.isSubscriber = false,
    this.isAdmin = false,
    this.currentUserId = '',
  });

  @override
  State<ChannelCommentsPage> createState() => _ChannelCommentsPageState();
}

class _ChannelCommentsPageState extends State<ChannelCommentsPage> {
  final _inputCtrl = TextEditingController();
  final List<ChannelComment> _comments = [];
  // Resolved avatar URL per senderId (null = no avatar / still pending). Cached
  // so each author's profile is only fetched once.
  final Map<String, String?> _avatarUrls = {};
  StreamSubscription? _socketSub;
  bool _loading = true;
  bool _sending = false;
  ChannelComment? _replyTo;
  // Fresh commentsEnabled verified from the view API on mount.
  late bool _commentsEnabledLocally;

  bool get _canComment => _commentsEnabledLocally && widget.isSubscriber;

  @override
  void initState() {
    super.initState();
    _commentsEnabledLocally = widget.commentsEnabled;
    _load();
    _socketSub = sl<SocketClient>().events.listen(_onSocketEvent);
    // Verify commentsEnabled freshly — widget.commentsEnabled may be stale (cached _channelView).
    if (widget.commentsEnabled) _refreshCommentsState();
  }

  Future<void> _refreshCommentsState() async {
    final result = await sl<ChannelRepository>().getChannelView(widget.channelId);
    if (!mounted) return;
    result.fold((_) {}, (view) {
      if (_commentsEnabledLocally != view.commentsEnabled) {
        setState(() => _commentsEnabledLocally = view.commentsEnabled);
      }
    });
  }

  @override
  void dispose() {
    _socketSub?.cancel();
    _inputCtrl.dispose();
    super.dispose();
  }

  void _onSocketEvent(SocketEvent event) {
    if (event.type != SocketEventType.channelNewComment || event.payload == null) return;
    final comment = ChannelComment.fromJson(event.payload!);
    if (comment.postId != widget.postId) return;
    if (_comments.any((c) => c.id == comment.id)) return;
    setState(() => _comments.add(comment));
    _resolveAvatars([comment]);
  }

  /// Resolves avatar URLs for the authors of [comments] (own messages skipped),
  /// reusing the cache so each profile is fetched at most once.
  Future<void> _resolveAvatars(List<ChannelComment> comments) async {
    final ids = comments
        .map((c) => c.senderId)
        .where((id) => id.isNotEmpty && id != widget.currentUserId && !_avatarUrls.containsKey(id))
        .toSet();
    if (ids.isEmpty) return;
    // Mark as in-flight so concurrent calls don't refetch the same author.
    for (final id in ids) {
      _avatarUrls[id] = null;
    }
    for (final id in ids) {
      final profileRes = await sl<UserRepository>().getUserProfile(id);
      final mediaId = profileRes.fold((_) => null, (p) => p.avatarMediaId);
      String? url;
      if (mediaId != null) {
        final urlRes = await sl<MediaRepository>().getMediaUrl(mediaId);
        url = urlRes.fold((_) => null, (u) => u);
      }
      if (!mounted) return;
      setState(() => _avatarUrls[id] = url);
    }
  }

  Future<void> _load() async {
    final result = await sl<ChannelRepository>()
        .getComments(widget.channelId, widget.postId, limit: 100);
    if (!mounted) return;
    result.fold(
      (f) {
        AppFeedback.showError(f.message);
        setState(() => _loading = false);
      },
      (data) {
        final (list, _) = data;
        setState(() {
          _comments
            ..clear()
            ..addAll(list);
          _loading = false;
        });
        _resolveAvatars(list);
      },
    );
  }

  Future<void> _send() async {
    final text = _inputCtrl.text.trim();
    if (text.isEmpty) return;
    setState(() => _sending = true);
    final result = await sl<ChannelRepository>().addComment(
      widget.channelId,
      widget.postId,
      text,
      replyToCommentId: _replyTo?.id,
    );
    if (!mounted) return;
    result.fold(
      (f) {
        AppFeedback.showError(f.message);
        setState(() => _sending = false);
      },
      (comment) {
        setState(() {
          if (!_comments.any((c) => c.id == comment.id)) _comments.add(comment);
          _inputCtrl.clear();
          _replyTo = null;
          _sending = false;
        });
      },
    );
  }

  void _onCommentLongPress(ChannelComment comment) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;
    final isAuthor = comment.senderId == widget.currentUserId;

    showModalBottomSheet(
      context: context,
      backgroundColor: theme.colorScheme.surface,
      shape: const RoundedRectangleBorder(
          borderRadius: BorderRadius.vertical(top: Radius.circular(24))),
      builder: (ctx) => SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const SizedBox(height: 12),
            if (_canComment)
              ListTile(
                leading: Icon(Icons.reply_rounded, color: theme.colorScheme.tertiary),
                title: Text(l10n.commentActionReply,
                    style: TextStyle(color: theme.colorScheme.onSurface)),
                onTap: () {
                  Navigator.pop(ctx);
                  setState(() => _replyTo = comment);
                },
              ),
            if (isAuthor)
              ListTile(
                leading: Icon(Icons.edit_outlined, color: theme.colorScheme.onSurface),
                title: Text(l10n.commentActionEdit,
                    style: TextStyle(color: theme.colorScheme.onSurface)),
                onTap: () {
                  Navigator.pop(ctx);
                  _editComment(comment);
                },
              ),
            if (isAuthor || widget.isAdmin)
              ListTile(
                leading: Icon(Icons.delete_outline, color: palette.danger),
                title: Text(l10n.commentActionDelete, style: TextStyle(color: palette.danger)),
                onTap: () {
                  Navigator.pop(ctx);
                  _deleteComment(comment);
                },
              ),
            const SizedBox(height: 8),
          ],
        ),
      ),
    );
  }

  Future<void> _editComment(ChannelComment comment) async {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;
    final ctrl = TextEditingController(text: comment.content);

    final newText = await showDialog<String>(
      context: context,
      builder: (ctx) => AlertDialog(
        backgroundColor: theme.colorScheme.surface,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(18)),
        content: TextField(
          controller: ctrl,
          autofocus: true,
          maxLines: 4,
          minLines: 1,
          style: TextStyle(fontFamily: 'Manrope', color: theme.colorScheme.onSurface),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: Text(l10n.buttonRetry,
                style: TextStyle(color: palette.textSecondary)),
          ),
          TextButton(
            onPressed: () => Navigator.pop(ctx, ctrl.text.trim()),
            child: Text(l10n.commentActionEdit,
                style: TextStyle(color: theme.colorScheme.tertiary, fontWeight: FontWeight.w700)),
          ),
        ],
      ),
    );
    if (newText == null || newText.isEmpty || newText == comment.content) return;
    final result = await sl<ChannelRepository>().editComment(comment.id, newText);
    if (!mounted) return;
    result.fold(
      (f) => AppFeedback.showError(f.message),
      (_) {
        final i = _comments.indexWhere((c) => c.id == comment.id);
        if (i != -1) {
          setState(() => _comments[i] = _comments[i].copyWith(
                content: newText,
                editedAt: DateTime.now().toUtc().toIso8601String(),
              ));
        }
      },
    );
  }

  Future<void> _deleteComment(ChannelComment comment) async {
    final result = await sl<ChannelRepository>().deleteComment(comment.id);
    if (!mounted) return;
    result.fold(
      (f) => AppFeedback.showError(f.message),
      (_) => setState(() => _comments.removeWhere((c) => c.id == comment.id)),
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
        leading: IconButton(
          icon: Icon(Icons.arrow_back_rounded, color: theme.colorScheme.onSurface),
          onPressed: () => context.pop(),
        ),
        title: Text(
          l10n.channelComments,
          style: TextStyle(
              fontFamily: 'Manrope', fontWeight: FontWeight.w700, color: theme.colorScheme.onSurface),
        ),
      ),
      body: Column(
        children: [
          // The post this thread belongs to.
          if (widget.post != null) _PostHeader(post: widget.post!),
          Expanded(
            child: _loading
                ? const Center(child: CircularProgressIndicator())
                : _comments.isEmpty
                    ? Center(
                        child: Text(
                          l10n.channelCommentsEmpty,
                          style: TextStyle(fontFamily: 'Manrope', color: palette.textSecondary),
                        ),
                      )
                    : Builder(
                        builder: (_) {
                          // Lookup map so reply bubbles can show the quoted parent.
                          final byId = {for (final c in _comments) c.id: c};
                          return ListView.builder(
                            padding: const EdgeInsets.symmetric(vertical: 8),
                            itemCount: _comments.length,
                            itemBuilder: (_, i) {
                              final c = _comments[i];
                              final parent = c.replyToCommentId == null
                                  ? null
                                  : byId[c.replyToCommentId];
                              return RepaintBoundary(
                                child: _CommentBubble(
                                  comment: c,
                                  isMine: c.senderId == widget.currentUserId,
                                  avatarUrl: _avatarUrls[c.senderId],
                                  replyTo: parent,
                                  onLongPress: () => _onCommentLongPress(c),
                                ),
                              );
                            },
                          );
                        },
                      ),
          ),
          _buildInputArea(context, theme, palette, l10n),
        ],
      ),
    );
  }

  Widget _buildInputArea(
    BuildContext context,
    ThemeData theme,
    AppPaletteTheme palette,
    dynamic l10n,
  ) {
    if (!_canComment) {
      return Container(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
        decoration: BoxDecoration(
          color: theme.colorScheme.surface,
          border: Border(top: BorderSide(color: palette.divider)),
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.lock_outline_rounded, size: 18, color: palette.textSecondary),
            const SizedBox(width: 8),
            Text(
              _commentsEnabledLocally
                  ? l10n.channelCommentSubscribePrompt
                  : l10n.channelCommentsDisabled,
              style: TextStyle(fontFamily: 'Manrope', fontSize: 13, color: palette.textSecondary),
            ),
          ],
        ),
      );
    }

    return Container(
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        border: Border(top: BorderSide(color: palette.divider)),
      ),
      child: SafeArea(
        top: false,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            if (_replyTo != null)
              Padding(
                padding: const EdgeInsets.fromLTRB(16, 8, 8, 0),
                child: Row(
                  children: [
                    Icon(Icons.reply_rounded, size: 14, color: theme.colorScheme.tertiary),
                    const SizedBox(width: 6),
                    Expanded(
                      child: Text(
                        '${_replyTo!.senderDisplayName}: ${_replyTo!.content}',
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                        style: TextStyle(
                            fontFamily: 'Manrope', fontSize: 12, color: palette.textSecondary),
                      ),
                    ),
                    IconButton(
                      icon: Icon(Icons.close_rounded, size: 16, color: palette.textSecondary),
                      onPressed: () => setState(() => _replyTo = null),
                    ),
                  ],
                ),
              ),
            Padding(
              padding: const EdgeInsets.fromLTRB(16, 8, 8, 8),
              child: Row(
                children: [
                  Expanded(
                    child: TextField(
                      controller: _inputCtrl,
                      minLines: 1,
                      maxLines: 4,
                      style: TextStyle(fontFamily: 'Manrope', color: theme.colorScheme.onSurface),
                      decoration: InputDecoration(
                        hintText: l10n.channelCommentHint,
                        hintStyle: TextStyle(color: palette.textSecondary),
                        border: InputBorder.none,
                      ),
                    ),
                  ),
                  IconButton(
                    icon: _sending
                        ? const SizedBox(
                            width: 18, height: 18, child: CircularProgressIndicator(strokeWidth: 2))
                        : Icon(Icons.send_rounded, color: theme.colorScheme.tertiary),
                    onPressed: _sending ? null : _send,
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

class _PostHeader extends StatelessWidget {
  final Message post;
  const _PostHeader({required this.post});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        border: Border(bottom: BorderSide(color: palette.divider)),
      ),
      child: Text(
        post.content,
        style: TextStyle(
            fontFamily: 'Manrope', fontSize: 14, height: 1.5, color: theme.colorScheme.onSurface),
      ),
    );
  }
}

/// A single comment rendered as a chat bubble — own comments align right with
/// no avatar; others align left with the author's avatar + display name, mirror-
/// ing the message bubbles in the main chat view.
class _CommentBubble extends StatelessWidget {
  final ChannelComment comment;
  final bool isMine;
  final String? avatarUrl;
  // Parent comment this one replies to (resolved from the loaded list), or null.
  final ChannelComment? replyTo;
  final VoidCallback onLongPress;

  const _CommentBubble({
    required this.comment,
    required this.isMine,
    required this.avatarUrl,
    required this.replyTo,
    required this.onLongPress,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;
    final maxWidth = MediaQuery.of(context).size.width * (isMine ? 0.80 : 0.72);

    final bubble = GestureDetector(
      onLongPress: () {
        HapticFeedback.mediumImpact();
        onLongPress();
      },
      child: ConstrainedBox(
        constraints: BoxConstraints(maxWidth: maxWidth),
        child: Container(
          padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
          decoration: BoxDecoration(
            color: isMine ? palette.messageMeBg : palette.messageOtherBg,
            borderRadius: BorderRadius.only(
              topLeft: const Radius.circular(18),
              topRight: const Radius.circular(18),
              bottomLeft: Radius.circular(isMine ? 18 : 4),
              bottomRight: Radius.circular(isMine ? 4 : 18),
            ),
            border: isMine ? null : Border.all(color: palette.divider),
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            mainAxisSize: MainAxisSize.min,
            children: [
              // Author name only on incoming bubbles (own name is redundant).
              if (!isMine)
                Padding(
                  padding: const EdgeInsets.only(bottom: 2),
                  child: Text(
                    comment.senderDisplayName,
                    style: TextStyle(
                      fontFamily: 'Manrope',
                      fontSize: 12,
                      fontWeight: FontWeight.w700,
                      color: theme.colorScheme.tertiary,
                    ),
                  ),
                ),
              // Quoted parent comment for replies.
              if (comment.isReply)
                _ReplyQuote(parent: replyTo, palette: palette, theme: theme),
              Text(
                comment.content,
                style: TextStyle(
                  fontFamily: 'Manrope',
                  fontSize: 15,
                  height: 1.45,
                  color: theme.colorScheme.onSurface,
                ),
              ),
              if (comment.isEdited)
                Padding(
                  padding: const EdgeInsets.only(top: 2),
                  child: Text(
                    l10n.commentEdited,
                    style: TextStyle(
                      fontFamily: 'Manrope',
                      fontSize: 11,
                      fontStyle: FontStyle.italic,
                      color: palette.textSecondary,
                    ),
                  ),
                ),
            ],
          ),
        ),
      ),
    );

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 2),
      child: Row(
        mainAxisAlignment: isMine ? MainAxisAlignment.end : MainAxisAlignment.start,
        crossAxisAlignment: CrossAxisAlignment.end,
        children: [
          if (!isMine) ...[
            SudaAvatar(
              mediaUrl: avatarUrl,
              initials: comment.senderDisplayName.isNotEmpty
                  ? comment.senderDisplayName
                  : '?',
              size: 28,
            ),
            const SizedBox(width: 6),
          ],
          Flexible(child: bubble),
          if (isMine) const SizedBox(width: 4),
        ],
      ),
    );
  }
}

/// Compact quote of the parent comment shown above a reply bubble.
class _ReplyQuote extends StatelessWidget {
  final ChannelComment? parent;
  final AppPaletteTheme palette;
  final ThemeData theme;

  const _ReplyQuote({
    required this.parent,
    required this.palette,
    required this.theme,
  });

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    // Parent may be outside the loaded window → neutral fallback label.
    final author = parent?.senderDisplayName ?? l10n.commentReplyDeleted;
    final text = parent?.content ?? '';

    return Container(
      margin: const EdgeInsets.only(bottom: 6),
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        border: Border(
          left: BorderSide(color: theme.colorScheme.tertiary, width: 3),
        ),
        color: theme.colorScheme.onSurface.withValues(alpha: 0.05),
        borderRadius: const BorderRadius.only(
          topRight: Radius.circular(8),
          bottomRight: Radius.circular(8),
        ),
      ),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            author,
            style: TextStyle(
              fontFamily: 'Manrope',
              fontSize: 11,
              fontWeight: FontWeight.w700,
              color: theme.colorScheme.tertiary,
            ),
          ),
          if (text.isNotEmpty)
            Text(
              text,
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
              style: TextStyle(
                fontFamily: 'Manrope',
                fontSize: 12,
                color: palette.textSecondary,
              ),
            ),
        ],
      ),
    );
  }
}
