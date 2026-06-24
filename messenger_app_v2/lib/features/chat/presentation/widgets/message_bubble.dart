import 'dart:convert';
import 'dart:io';

import 'package:cached_network_image/cached_network_image.dart';
import 'package:chewie/chewie.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:open_file/open_file.dart';
import 'package:photo_manager/photo_manager.dart';
import 'package:photo_view/photo_view.dart';
import 'package:shimmer/shimmer.dart';
import 'package:video_player/video_player.dart';

import '../../../../app/config/theme/app_theme.dart';
import '../../../../features/media/domain/repositories/media_repository.dart';
import '../../../../features/wallet/data/models/wallet_models.dart';
import '../../../../shared/presentation/feedback/app_feedback.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/widgets/suda_avatar.dart';
import '../../data/models/chat_models.dart';
import 'voice_message_player.dart';

class MessageBubble extends StatelessWidget {
  final Message message;
  final bool isMe;
  final VoidCallback? onLongPress;
  final VoidCallback? onReadReceiptTap;
  final MediaRepository? mediaRepo;
  // Group chats: show sender avatar (newest in run) + name (oldest in run).
  final bool isGroupIncoming;
  final String? senderName;
  final String? senderAvatarUrl;
  // Sender id — lets the avatar resolve via the user's profile when no
  // pre-resolved avatar URL is available.
  final String? senderId;
  final bool showAvatar;
  final bool showName;
  // Channel posts: comments affordance below the bubble.
  final bool showComments;
  final int? commentCount;
  final VoidCallback? onOpenComments;
  // Search jump-to: briefly highlight this bubble when navigated to.
  final bool isHighlighted;

  const MessageBubble({
    super.key,
    required this.message,
    required this.isMe,
    this.onLongPress,
    this.onReadReceiptTap,
    this.mediaRepo,
    this.isGroupIncoming = false,
    this.senderName,
    this.senderAvatarUrl,
    this.senderId,
    this.showAvatar = false,
    this.showName = false,
    this.showComments = false,
    this.commentCount,
    this.onOpenComments,
    this.isHighlighted = false,
  });

  // Server may echo the type in lower-case (e.g. 'image'); normalise so media
  // bubbles render regardless of casing.
  String get _type => message.type.toUpperCase();

  bool get _isMediaMessage =>
      _type == 'IMAGE' ||
      _type == 'VIDEO' ||
      _type == 'VOICE' ||
      _type == 'FILE';

  @override
  Widget build(BuildContext context) {
    // System transaction messages render as centred labels, not chat bubbles.
    if (_type == 'SUDA_TRANSFER' || _type == 'DONATION') {
      return _SystemTransactionBubble(message: message);
    }

    // Generic system events (group created, members joined, etc.) render as a
    // centred label with localized text instead of a sender bubble.
    if (_type == 'SYSTEM') {
      return _SystemMessageBubble(message: message);
    }

    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);
    final maxWidth = MediaQuery.of(context).size.width * (isMe ? 0.80 : 0.72);

    // IMAGE bubbles are compact — no padding around the image
    final isImageType = _type == 'IMAGE' || _type == 'VIDEO';
    final innerPadding = isImageType && message.attachmentMediaId != null
        ? EdgeInsets.zero
        : const EdgeInsets.symmetric(horizontal: 14, vertical: 10);

    final bubble = GestureDetector(
      onLongPress: () {
        HapticFeedback.mediumImpact();
        onLongPress?.call();
      },
      child: ConstrainedBox(
        constraints: BoxConstraints(maxWidth: maxWidth),
        child: Container(
          padding: innerPadding,
          decoration: BoxDecoration(
            color: isMe ? palette.messageMeBg : palette.messageOtherBg,
            borderRadius: BorderRadius.only(
              topLeft: const Radius.circular(18),
              topRight: const Radius.circular(18),
              bottomLeft: Radius.circular(isMe ? 18 : 4),
              bottomRight: Radius.circular(isMe ? 4 : 18),
            ),
            border: isMe ? null : Border.all(color: palette.divider),
          ),
          child: (showName && senderName?.isNotEmpty == true)
              ? Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Padding(
                      padding: const EdgeInsets.only(bottom: 2),
                      child: Text(
                        senderName!,
                        style: TextStyle(
                          fontFamily: 'Manrope',
                          fontSize: 12,
                          fontWeight: FontWeight.w700,
                          color: theme.colorScheme.tertiary,
                        ),
                      ),
                    ),
                    _buildContent(context, palette, theme),
                  ],
                )
              : _buildContent(context, palette, theme),
        ),
      ),
    );

    // Group incoming: leading avatar gutter (avatar on newest of a run).
    final Widget leading = isGroupIncoming
        ? Padding(
            padding: const EdgeInsets.only(right: 6),
            child: SizedBox(
              width: 28,
              child: showAvatar
                  ? SudaAvatar(
                      mediaUrl: senderAvatarUrl,
                      userId: senderId,
                      initials: (senderName?.isNotEmpty == true) ? senderName! : '?',
                      size: 28,
                    )
                  : null,
            ),
          )
        : SizedBox(width: isMe ? 4 : 4);

    final row = Row(
      mainAxisAlignment: isMe ? MainAxisAlignment.end : MainAxisAlignment.start,
      crossAxisAlignment: CrossAxisAlignment.end,
      children: [
        if (!isMe) leading,
        Flexible(
          child: Column(
            crossAxisAlignment:
                isMe ? CrossAxisAlignment.end : CrossAxisAlignment.start,
            children: [
              bubble,
              if (showComments) _CommentsButton(
                count: commentCount,
                palette: palette,
                onTap: onOpenComments,
              ),
            ],
          ),
        ),
        if (isMe) const SizedBox(width: 4),
      ],
    );

    // Search jump-to: fade a subtle highlight behind the whole row while the
    // bubble is the active jump target, then fade out when it clears.
    return AnimatedContainer(
      duration: const Duration(milliseconds: 350),
      curve: Curves.easeOut,
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 2),
      color: isHighlighted
          ? theme.colorScheme.tertiary.withValues(alpha: 0.18)
          : Colors.transparent,
      child: row,
    );
  }

  Widget _buildContent(
    BuildContext context,
    AppPaletteTheme palette,
    ThemeData theme,
  ) {
    if (_isMediaMessage && message.attachmentMediaId != null) {
      return switch (_type) {
        'IMAGE' => _ImageContent(
            mediaId: message.attachmentMediaId!,
            mediaRepo: mediaRepo,
            isMe: isMe,
            message: message,
            onReadReceiptTap: onReadReceiptTap,
          ),
        'VIDEO' => _VideoContent(
            mediaId: message.attachmentMediaId!,
            mediaRepo: mediaRepo,
            isMe: isMe,
            message: message,
            onReadReceiptTap: onReadReceiptTap,
          ),
        'VOICE' => Padding(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                if (message.isForwarded) _ForwardedLabel(palette: palette, message: message),
                VoiceMessagePlayer(
                  mediaId: message.attachmentMediaId,
                  durationText: message.content.isEmpty ? null : message.content,
                  isMe: isMe,
                  mediaRepo: mediaRepo,
                ),
                const SizedBox(height: 4),
                _MessageFooter(message: message, isMe: isMe, onReadReceiptTap: onReadReceiptTap),
              ],
            ),
          ),
        'FILE' => Padding(
            padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                if (message.isForwarded) _ForwardedLabel(palette: palette, message: message),
                _FileContent(
                  mediaId: message.attachmentMediaId!,
                  filename: message.content.isEmpty
                      ? context.l10n.messageTypeFile
                      : message.content,
                  mediaRepo: mediaRepo,
                  palette: palette,
                ),
                const SizedBox(height: 4),
                _MessageFooter(message: message, isMe: isMe, onReadReceiptTap: onReadReceiptTap),
              ],
            ),
          ),
        _ => _TextContent(message: message, isMe: isMe, palette: palette, onReadReceiptTap: onReadReceiptTap),
      };
    }
    // Media message whose attachment id is missing — show a localized label with
    // an icon instead of an empty text bubble.
    if (_isMediaMessage) {
      return _MediaPlaceholder(
        type: _type,
        message: message,
        isMe: isMe,
        palette: palette,
        onReadReceiptTap: onReadReceiptTap,
      );
    }
    return _TextContent(
      message: message,
      isMe: isMe,
      palette: palette,
      onReadReceiptTap: onReadReceiptTap,
    );
  }
}

// ─── Media placeholder (attachment id missing) ────────────────

class _MediaPlaceholder extends StatelessWidget {
  final String type;
  final Message message;
  final bool isMe;
  final AppPaletteTheme palette;
  final VoidCallback? onReadReceiptTap;

  const _MediaPlaceholder({
    required this.type,
    required this.message,
    required this.isMe,
    required this.palette,
    this.onReadReceiptTap,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final l10n = context.l10n;
    final (IconData icon, String label) = switch (type) {
      'IMAGE' => (Icons.photo_outlined, l10n.messageTypeImage),
      'VIDEO' => (Icons.videocam_outlined, l10n.messageTypeVideo),
      'VOICE' => (Icons.mic_none_rounded, l10n.messageTypeVoice),
      _ => (Icons.insert_drive_file_outlined, l10n.messageTypeFile),
    };
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(icon, size: 18, color: theme.colorScheme.onSurface),
            const SizedBox(width: 6),
            Text(
              label,
              style: TextStyle(
                fontFamily: 'Manrope',
                fontSize: 15,
                color: theme.colorScheme.onSurface,
                height: 1.45,
              ),
            ),
          ],
        ),
        const SizedBox(height: 4),
        _MessageFooter(message: message, isMe: isMe, onReadReceiptTap: onReadReceiptTap),
      ],
    );
  }
}

// ─── Text content ─────────────────────────────────────────────

class _TextContent extends StatelessWidget {
  final Message message;
  final bool isMe;
  final AppPaletteTheme palette;
  final VoidCallback? onReadReceiptTap;

  const _TextContent({
    required this.message,
    required this.isMe,
    required this.palette,
    this.onReadReceiptTap,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        if (message.isReply) _ReplyQuote(palette: palette, message: message),
        if (message.isForwarded) _ForwardedLabel(palette: palette, message: message),
        Text(
          message.content,
          style: TextStyle(
            fontFamily: 'Manrope',
            fontSize: 15,
            color: theme.colorScheme.onSurface,
            height: 1.45,
          ),
        ),
        const SizedBox(height: 4),
        _MessageFooter(message: message, isMe: isMe, onReadReceiptTap: onReadReceiptTap),
      ],
    );
  }
}

// ─── Image content ────────────────────────────────────────────

class _ImageContent extends StatelessWidget {
  final String mediaId;
  final MediaRepository? mediaRepo;
  final bool isMe;
  final Message message;
  final VoidCallback? onReadReceiptTap;

  const _ImageContent({
    required this.mediaId,
    required this.mediaRepo,
    required this.isMe,
    required this.message,
    this.onReadReceiptTap,
  });

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<String?>(
      future: mediaRepo?.getMediaUrl(mediaId).then((r) => r.fold((_) => null, (u) => u)),
      builder: (ctx, snap) {
        final url = snap.data;
        return ClipRRect(
          borderRadius: const BorderRadius.only(
            topLeft: Radius.circular(18),
            topRight: Radius.circular(18),
            bottomLeft: Radius.circular(18),
            bottomRight: Radius.circular(18),
          ),
          child: Stack(
            children: [
              // Image or shimmer
              url != null
                  ? GestureDetector(
                      onTap: () => _openFullscreen(ctx, url),
                      child: CachedNetworkImage(
                        imageUrl: url,
                        width: 200,
                        height: 200,
                        fit: BoxFit.cover,
                        placeholder: (_, __) => _shimmer(),
                        errorWidget: (_, __, ___) => _errorPlaceholder(ctx),
                      ),
                    )
                  : _shimmer(),

              // Footer overlay at bottom of image
              Positioned(
                bottom: 0,
                right: 0,
                child: Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    gradient: LinearGradient(
                      begin: Alignment.topCenter,
                      end: Alignment.bottomCenter,
                      colors: [Colors.transparent, Colors.black.withValues(alpha: 0.5)],
                    ),
                  ),
                  child: _MessageFooter(
                    message: message,
                    isMe: isMe,
                    onReadReceiptTap: onReadReceiptTap,
                  ),
                ),
              ),
            ],
          ),
        );
      },
    );
  }

  void _openFullscreen(BuildContext context, String url) {
    Navigator.of(context).push(MaterialPageRoute(
      fullscreenDialog: true,
      builder: (_) => _FullscreenImageScreen(
        url: url,
        mediaId: mediaId,
        mediaRepo: mediaRepo,
      ),
    ));
  }

  Widget _shimmer() => Shimmer.fromColors(
        baseColor: Colors.grey.shade800,
        highlightColor: Colors.grey.shade600,
        child: Container(width: 200, height: 200, color: Colors.white),
      );

  Widget _errorPlaceholder(BuildContext context) => Container(
        width: 200,
        height: 200,
        color: Theme.of(context).colorScheme.surface,
        child: const Icon(Icons.broken_image_outlined, size: 48, color: Colors.white54),
      );
}

// ─── Save media to the device gallery ─────────────────────────

/// Downloads [mediaId] to cache and saves it into the device gallery via
/// photo_manager. Feedback strings are read before any await so the global
/// AppFeedback messenger can be used safely afterwards.
Future<void> _saveMediaToGallery({
  required BuildContext context,
  required MediaRepository? repo,
  required String mediaId,
  required String filename,
  required bool isVideo,
}) async {
  final l10n = context.l10n;
  final errorMsg = l10n.saveToGalleryError;
  final successMsg = l10n.saveToGallerySuccess;
  if (repo == null) {
    AppFeedback.showError(errorMsg);
    return;
  }

  final perm = await PhotoManager.requestPermissionExtend();
  if (!perm.hasAccess) {
    AppFeedback.showError(errorMsg);
    return;
  }

  final dl = await repo.downloadToCache(mediaId, filename);
  await dl.fold(
    (failure) async => AppFeedback.showError(errorMsg),
    (path) async {
      try {
        if (isVideo) {
          await PhotoManager.editor.saveVideo(File(path), title: filename);
        } else {
          await PhotoManager.editor.saveImageWithPath(path, title: filename);
        }
        AppFeedback.showSuccess(successMsg);
      } catch (_) {
        AppFeedback.showError(errorMsg);
      }
    },
  );
}

// ─── Fullscreen image viewer (with save-to-gallery) ───────────

class _FullscreenImageScreen extends StatefulWidget {
  final String url;
  final String mediaId;
  final MediaRepository? mediaRepo;

  const _FullscreenImageScreen({
    required this.url,
    required this.mediaId,
    required this.mediaRepo,
  });

  @override
  State<_FullscreenImageScreen> createState() => _FullscreenImageScreenState();
}

class _FullscreenImageScreenState extends State<_FullscreenImageScreen> {
  bool _saving = false;

  Future<void> _save() async {
    if (_saving) return;
    setState(() => _saving = true);
    await _saveMediaToGallery(
      context: context,
      repo: widget.mediaRepo,
      mediaId: widget.mediaId,
      filename: 'IMG_${widget.mediaId}.jpg',
      isVideo: false,
    );
    if (mounted) setState(() => _saving = false);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.black,
      appBar: AppBar(
        backgroundColor: Colors.black,
        iconTheme: const IconThemeData(color: Colors.white),
        actions: [_SaveButton(saving: _saving, onTap: _save)],
      ),
      body: PhotoView(imageProvider: CachedNetworkImageProvider(widget.url)),
    );
  }
}

// ─── Video content ────────────────────────────────────────────

class _VideoContent extends StatefulWidget {
  final String mediaId;
  final MediaRepository? mediaRepo;
  final bool isMe;
  final Message message;
  final VoidCallback? onReadReceiptTap;

  const _VideoContent({
    required this.mediaId,
    required this.mediaRepo,
    required this.isMe,
    required this.message,
    this.onReadReceiptTap,
  });

  @override
  State<_VideoContent> createState() => _VideoContentState();
}

class _VideoContentState extends State<_VideoContent> {
  VideoPlayerController? _controller;
  String? _url;
  bool _ready = false;
  bool _error = false;

  @override
  void initState() {
    super.initState();
    _init();
  }

  Future<void> _init() async {
    final repo = widget.mediaRepo;
    if (repo == null) {
      setState(() => _error = true);
      return;
    }
    final res = await repo.getMediaUrl(widget.mediaId);
    final url = res.fold((_) => null, (u) => u);
    if (!mounted) return;
    if (url == null || url.isEmpty) {
      setState(() => _error = true);
      return;
    }
    _url = url;
    final ctrl = VideoPlayerController.networkUrl(Uri.parse(url));
    try {
      await ctrl.initialize();
      if (!mounted) {
        ctrl.dispose();
        return;
      }
      _controller = ctrl;
      setState(() => _ready = true);
    } catch (_) {
      ctrl.dispose();
      if (mounted) setState(() => _error = true);
    }
  }

  @override
  void dispose() {
    _controller?.dispose();
    super.dispose();
  }

  void _openPlayer() {
    final url = _url;
    if (url == null) return;
    Navigator.of(context).push(MaterialPageRoute(
      fullscreenDialog: true,
      builder: (_) => _VideoPlayerScreen(
        url: url,
        mediaId: widget.mediaId,
        mediaRepo: widget.mediaRepo,
      ),
    ));
  }

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: _ready ? _openPlayer : null,
      child: ClipRRect(
        borderRadius: BorderRadius.circular(18),
        child: Stack(
          alignment: Alignment.center,
          children: [
            SizedBox(
              width: 220,
              height: 220,
              child: _ready && _controller != null
                  ? FittedBox(
                      fit: BoxFit.cover,
                      child: SizedBox(
                        width: _controller!.value.size.width,
                        height: _controller!.value.size.height,
                        child: VideoPlayer(_controller!),
                      ),
                    )
                  : _error
                      ? Container(
                          color: Theme.of(context).colorScheme.surface,
                          child: const Icon(Icons.videocam_off_outlined,
                              size: 48, color: Colors.white54),
                        )
                      : Shimmer.fromColors(
                          baseColor: Colors.grey.shade800,
                          highlightColor: Colors.grey.shade600,
                          child: Container(color: Colors.white),
                        ),
            ),
            // Centre play affordance.
            if (_ready)
              Container(
                width: 54,
                height: 54,
                decoration: BoxDecoration(
                  color: Colors.black.withValues(alpha: 0.45),
                  shape: BoxShape.circle,
                ),
                child: const Icon(Icons.play_arrow_rounded,
                    size: 34, color: Colors.white),
              ),
            // Footer overlay (timestamp / read receipt).
            Positioned(
              bottom: 0,
              right: 0,
              child: Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                decoration: BoxDecoration(
                  gradient: LinearGradient(
                    begin: Alignment.topCenter,
                    end: Alignment.bottomCenter,
                    colors: [Colors.transparent, Colors.black.withValues(alpha: 0.5)],
                  ),
                ),
                child: _MessageFooter(
                  message: widget.message,
                  isMe: widget.isMe,
                  onReadReceiptTap: widget.onReadReceiptTap,
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

// ─── Fullscreen video player (Chewie, with save-to-gallery) ───

class _VideoPlayerScreen extends StatefulWidget {
  final String url;
  final String mediaId;
  final MediaRepository? mediaRepo;

  const _VideoPlayerScreen({
    required this.url,
    required this.mediaId,
    required this.mediaRepo,
  });

  @override
  State<_VideoPlayerScreen> createState() => _VideoPlayerScreenState();
}

class _VideoPlayerScreenState extends State<_VideoPlayerScreen> {
  VideoPlayerController? _videoController;
  ChewieController? _chewieController;
  bool _saving = false;

  @override
  void initState() {
    super.initState();
    _init();
  }

  Future<void> _init() async {
    final ctrl = VideoPlayerController.networkUrl(Uri.parse(widget.url));
    await ctrl.initialize();
    if (!mounted) {
      ctrl.dispose();
      return;
    }
    setState(() {
      _videoController = ctrl;
      _chewieController = ChewieController(
        videoPlayerController: ctrl,
        autoPlay: true,
        looping: false,
        allowFullScreen: true,
      );
    });
  }

  @override
  void dispose() {
    _chewieController?.dispose();
    _videoController?.dispose();
    super.dispose();
  }

  Future<void> _save() async {
    if (_saving) return;
    setState(() => _saving = true);
    await _saveMediaToGallery(
      context: context,
      repo: widget.mediaRepo,
      mediaId: widget.mediaId,
      filename: 'VID_${widget.mediaId}.mp4',
      isVideo: true,
    );
    if (mounted) setState(() => _saving = false);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.black,
      appBar: AppBar(
        backgroundColor: Colors.black,
        iconTheme: const IconThemeData(color: Colors.white),
        actions: [_SaveButton(saving: _saving, onTap: _save)],
      ),
      // Keep the Chewie timeline/controls clear of the Android navigation bar.
      body: SafeArea(
        child: Center(
          child: _chewieController != null
              ? Chewie(controller: _chewieController!)
              : const CircularProgressIndicator(),
        ),
      ),
    );
  }
}

// ─── Save-to-gallery app-bar button ───────────────────────────

class _SaveButton extends StatelessWidget {
  final bool saving;
  final VoidCallback onTap;

  const _SaveButton({required this.saving, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return IconButton(
      tooltip: context.l10n.saveToGallery,
      onPressed: saving ? null : onTap,
      icon: saving
          ? const SizedBox(
              width: 20,
              height: 20,
              child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white),
            )
          : const Icon(Icons.download_rounded, color: Colors.white),
    );
  }
}

// ─── File content ─────────────────────────────────────────────

class _FileContent extends StatefulWidget {
  final String mediaId;
  final String filename;
  final MediaRepository? mediaRepo;
  final AppPaletteTheme palette;

  const _FileContent({
    required this.mediaId,
    required this.filename,
    required this.mediaRepo,
    required this.palette,
  });

  @override
  State<_FileContent> createState() => _FileContentState();
}

class _FileContentState extends State<_FileContent> {
  bool _busy = false;

  Future<void> _openFile() async {
    final repo = widget.mediaRepo;
    if (_busy || repo == null) return;
    setState(() => _busy = true);
    final l10n = context.l10n;
    // Download to cache (presigned URL) then hand off to the system viewer —
    // OpenFile needs a local path, not a remote URL.
    final result = await repo.downloadToCache(widget.mediaId, widget.filename);
    if (!mounted) return;
    setState(() => _busy = false);
    await result.fold(
      (failure) async => AppFeedback.showError(l10n.fileOpenError),
      (path) async {
        final opened = await OpenFile.open(path);
        if (opened.type != ResultType.done) {
          AppFeedback.showError(l10n.fileOpenError);
        }
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final l10n = context.l10n;

    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Container(
          width: 40,
          height: 40,
          decoration: BoxDecoration(
            color: widget.palette.surface2,
            borderRadius: BorderRadius.circular(10),
          ),
          child: Icon(Icons.insert_drive_file_outlined,
              size: 22, color: theme.colorScheme.tertiary),
        ),
        const SizedBox(width: 10),
        Flexible(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                widget.filename,
                style: TextStyle(
                  fontFamily: 'Manrope',
                  fontSize: 13,
                  fontWeight: FontWeight.w600,
                  color: theme.colorScheme.onSurface,
                ),
                maxLines: 2,
                overflow: TextOverflow.ellipsis,
              ),
              GestureDetector(
                onTap: _openFile,
                child: _busy
                    ? Padding(
                        padding: const EdgeInsets.only(top: 4),
                        child: SizedBox(
                          width: 12,
                          height: 12,
                          child: CircularProgressIndicator(
                            strokeWidth: 2,
                            color: theme.colorScheme.tertiary,
                          ),
                        ),
                      )
                    : Text(
                        l10n.openFile,
                        style: TextStyle(
                          fontFamily: 'Manrope',
                          fontSize: 12,
                          color: theme.colorScheme.tertiary,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
              ),
            ],
          ),
        ),
      ],
    );
  }
}

// ─── Forwarded label ──────────────────────────────────────────

class _ForwardedLabel extends StatelessWidget {
  final AppPaletteTheme palette;
  final Message message;
  const _ForwardedLabel({required this.palette, required this.message});

  @override
  Widget build(BuildContext context) {
    final senderName = message.forwardedFrom?.senderName;
    final label = (senderName != null && senderName.isNotEmpty)
        ? context.l10n.messageForwardedFrom(senderName)
        : context.l10n.messageForwardedLabel;

    return Padding(
      padding: const EdgeInsets.only(bottom: 4),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(Icons.forward_rounded, size: 12, color: Colors.white.withValues(alpha: 0.5)),
          const SizedBox(width: 4),
          Text(
            label,
            style: TextStyle(
              fontFamily: 'Manrope',
              fontSize: 11,
              fontStyle: FontStyle.italic,
              color: Colors.white.withValues(alpha: 0.5),
            ),
          ),
        ],
      ),
    );
  }
}

// ─── Message footer ───────────────────────────────────────────

class _MessageFooter extends StatelessWidget {
  final Message message;
  final bool isMe;
  final VoidCallback? onReadReceiptTap;

  const _MessageFooter({
    required this.message,
    required this.isMe,
    this.onReadReceiptTap,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      mainAxisAlignment: MainAxisAlignment.end,
      children: [
        if (message.isEdited) ...[
          Text(
            context.l10n.messageEdited,
            style: TextStyle(
              fontFamily: 'Manrope',
              fontSize: 11,
              fontStyle: FontStyle.italic,
              color: Colors.white.withValues(alpha: 0.35),
            ),
          ),
          const SizedBox(width: 4),
        ],
        Text(
          _formatTime(message.createdAt),
          style: TextStyle(
            fontFamily: 'Manrope',
            fontSize: 11,
            color: Colors.white.withValues(alpha: 0.45),
          ),
        ),
        if (isMe) ...[
          const SizedBox(width: 4),
          GestureDetector(
            onTap: onReadReceiptTap,
            child: _TickIcon(status: message.status),
          ),
        ],
      ],
    );
  }

  String _formatTime(String iso) {
    try {
      final d = DateTime.parse(iso).toLocal();
      return '${d.hour.toString().padLeft(2, '0')}:${d.minute.toString().padLeft(2, '0')}';
    } catch (_) {
      return '';
    }
  }
}

// ─── Reply quote ──────────────────────────────────────────────

class _ReplyQuote extends StatelessWidget {
  final AppPaletteTheme palette;
  final Message message;

  const _ReplyQuote({required this.palette, required this.message});

  // Returns the text to show as reply preview content.
  String _previewText(BuildContext context) {
    final preview = message.replyPreview;
    if (preview == null) {
      // Fallback for old messages that don't carry hydrated preview yet.
      return '#${message.replyToMessageId}';
    }
    if (preview.deleted) return context.l10n.replyDeleted;
    // System transaction messages carry JSON — show a short label instead.
    final t = preview.type.toUpperCase();
    if (t == 'SUDA_TRANSFER' || t == 'DONATION') {
      try {
        final payload = jsonDecode(preview.content) as Map<String, dynamic>;
        final wei = payload['amount_wei'] as String? ?? '0';
        final amount = weiToSudaDisplay(wei);
        return t == 'SUDA_TRANSFER'
            ? '💸 $amount SUDA'
            : '💝 $amount SUDA';
      } catch (_) {
        return preview.content;
      }
    }
    return preview.content;
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final l10n = context.l10n;
    final preview = message.replyPreview;
    final authorName = preview?.senderName ?? l10n.messageReply;

    return Container(
      margin: const EdgeInsets.only(bottom: 8),
      padding: const EdgeInsets.all(8),
      decoration: BoxDecoration(
        color: Colors.black.withValues(alpha: 0.18),
        borderRadius: BorderRadius.circular(10),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Fixed height avoids "infinite height" from CrossAxisAlignment.stretch
          // when the bubble is laid out inside the reversed message ListView.
          Container(
            width: 3,
            height: 32,
            decoration: BoxDecoration(
              color: theme.colorScheme.tertiary,
              borderRadius: BorderRadius.circular(2),
            ),
          ),
          const SizedBox(width: 8),
          Flexible(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Icon(Icons.reply_rounded, size: 12, color: theme.colorScheme.tertiary),
                    const SizedBox(width: 4),
                    Flexible(
                      child: Text(
                        authorName,
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                        style: TextStyle(
                          fontFamily: 'Manrope',
                          fontSize: 12,
                          fontWeight: FontWeight.w700,
                          color: theme.colorScheme.tertiary,
                        ),
                      ),
                    ),
                  ],
                ),
                Text(
                  _previewText(context),
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
          ),
        ],
      ),
    );
  }
}

// ─── Comments button (channel posts) ──────────────────────────

class _CommentsButton extends StatelessWidget {
  final int? count;
  final AppPaletteTheme palette;
  final VoidCallback? onTap;
  const _CommentsButton({required this.count, required this.palette, this.onTap});

  @override
  Widget build(BuildContext context) {
    final label = (count == null)
        ? context.l10n.channelComments
        : context.l10n.channelCommentsCount(count!);
    return Padding(
      padding: const EdgeInsets.only(top: 4, left: 4),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(10),
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 4),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(Icons.mode_comment_outlined, size: 14, color: palette.textSecondary),
              const SizedBox(width: 6),
              Text(
                label,
                style: TextStyle(
                  fontFamily: 'Manrope',
                  fontSize: 12,
                  fontWeight: FontWeight.w600,
                  color: palette.textSecondary,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

// ─── Tick icon ────────────────────────────────────────────────

class _TickIcon extends StatelessWidget {
  final String status;
  const _TickIcon({required this.status});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return switch (status) {
      MessageStatus.read    => Icon(Icons.done_all_rounded, size: 12, color: theme.colorScheme.tertiary),
      MessageStatus.sent    => Icon(Icons.done_all_rounded, size: 12, color: Colors.white.withValues(alpha: 0.5)),
      MessageStatus.sending => Icon(Icons.schedule_rounded, size: 11, color: Colors.white.withValues(alpha: 0.5)),
      MessageStatus.failed  => Icon(Icons.error_outline_rounded, size: 12, color: Colors.red.shade400),
      _                     => Icon(Icons.done_rounded, size: 12, color: Colors.white.withValues(alpha: 0.5)),
    };
  }
}

/// Centred system bubble for SUDA_TRANSFER and DONATION messages.
/// content is a JSON string (§2.1): {from_display_name, to_display_name,
/// from_address, to_address, from_user_id, to_user_id, amount_wei, comment, ...}
class _SystemTransactionBubble extends StatelessWidget {
  final Message message;
  const _SystemTransactionBubble({required this.message});

  // Resolve the best human-readable name for a transfer participant.
  // Priority: display_name → shortened address → shortened user_id.
  static String _resolveName(
    Map<String, dynamic> payload,
    String displayKey,
    String addressKey,
    String idKey,
  ) {
    final display = payload[displayKey] as String?;
    if (display != null && display.isNotEmpty) return display;
    final addr = payload[addressKey] as String?;
    if (addr != null && addr.length > 8) return '${addr.substring(0, 6)}…';
    final id = payload[idKey] as String?;
    if (id != null && id.length > 8) return '${id.substring(0, 8)}…';
    return '?';
  }

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    String mainLabel;
    String? subLabel;
    try {
      final payload = jsonDecode(message.content) as Map<String, dynamic>;
      final amountWei = payload['amount_wei'] as String? ?? '0';
      final amount = weiToSudaDisplay(amountWei);

      if (message.type == 'SUDA_TRANSFER') {
        final from = _resolveName(payload, 'from_display_name', 'from_address', 'from_user_id');
        final to = _resolveName(payload, 'to_display_name', 'to_address', 'to_user_id');
        mainLabel = l10n.msgSudaTransfer(from, to, amount);
        final comment = payload['comment'] as String?;
        if (comment != null && comment.isNotEmpty) subLabel = comment;
      } else {
        // DONATION
        final from = _resolveName(payload, 'from_display_name', 'from_address', 'from_user_id');
        mainLabel = l10n.msgDonation(from, amount);
        final msg = payload['message'] as String?;
        if (msg != null && msg.isNotEmpty) subLabel = msg;
      }
    } catch (_) {
      mainLabel = message.content;
    }

    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 6, horizontal: 24),
      child: Center(
        child: Container(
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
          decoration: BoxDecoration(
            color: palette.surface2,
            borderRadius: BorderRadius.circular(12),
          ),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Text(
                mainLabel,
                style: TextStyle(
                  fontFamily: 'Manrope',
                  fontSize: 13,
                  color: palette.textSecondary,
                ),
                textAlign: TextAlign.center,
              ),
              if (subLabel != null)
                Padding(
                  padding: const EdgeInsets.only(top: 2),
                  child: Text(
                    subLabel,
                    style: TextStyle(
                      fontFamily: 'Manrope',
                      fontSize: 12,
                      fontStyle: FontStyle.italic,
                      color: palette.textSecondary,
                    ),
                    textAlign: TextAlign.center,
                  ),
                ),
            ],
          ),
        ),
      ),
    );
  }
}

/// Centred system bubble for generic SYSTEM events (group lifecycle, etc.).
/// `content` carries a server event code (e.g. "GROUP_CREATED") which is mapped
/// to a localized, user-facing label — the raw code is never shown.
class _SystemMessageBubble extends StatelessWidget {
  final Message message;
  const _SystemMessageBubble({required this.message});

  String _label(BuildContext context) {
    final l10n = context.l10n;
    // Extend this switch as the backend introduces more system event codes.
    return switch (message.content) {
      'GROUP_CREATED' => l10n.sysGroupCreated,
      _ => l10n.sysGenericEvent,
    };
  }

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;

    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 6, horizontal: 24),
      child: Center(
        child: Container(
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
          decoration: BoxDecoration(
            color: palette.surface2,
            borderRadius: BorderRadius.circular(12),
          ),
          child: Text(
            _label(context),
            style: TextStyle(
              fontFamily: 'Manrope',
              fontSize: 13,
              color: palette.textSecondary,
            ),
            textAlign: TextAlign.center,
          ),
        ),
      ),
    );
  }
}
