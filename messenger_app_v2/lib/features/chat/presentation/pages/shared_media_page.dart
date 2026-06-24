import 'package:cached_network_image/cached_network_image.dart';
import 'package:chewie/chewie.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:open_file/open_file.dart';
import 'package:photo_view/photo_view.dart';
import 'package:video_player/video_player.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import '../../../../app/config/theme/app_theme.dart';
import '../../../../features/media/domain/repositories/media_repository.dart';
import '../../../../shared/presentation/feedback/app_feedback.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../data/models/chat_models.dart';
import '../../domain/repositories/chat_repository.dart';

/// Shared-media gallery for a chat — a grid of all images/videos plus lists of
/// files and audio. Works for DIRECT / GROUP / CHANNEL via GET /chats/{id}/media.
class SharedMediaPage extends StatefulWidget {
  final String chatId;
  final String chatName;

  const SharedMediaPage({super.key, required this.chatId, this.chatName = ''});

  @override
  State<SharedMediaPage> createState() => _SharedMediaPageState();
}

class _SharedMediaPageState extends State<SharedMediaPage> {
  ChatMedia? _media;
  bool _loading = true;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    final result = await sl<ChatRepository>().getChatMedia(widget.chatId);
    if (!mounted) return;
    result.fold(
      (f) {
        AppFeedback.showError(f.message);
        setState(() => _loading = false);
      },
      (media) => setState(() {
        _media = media;
        _loading = false;
      }),
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;
    final media = _media;

    // Media tab combines images + videos (per Shared-Media-Frontend.md §4).
    final mediaItems = <({MediaItem item, bool isVideo})>[
      ...(media?.images ?? const []).map((m) => (item: m, isVideo: false)),
      ...(media?.videos ?? const []).map((m) => (item: m, isVideo: true)),
    ];

    return DefaultTabController(
      length: 3,
      child: Scaffold(
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
            l10n.sharedMediaTitle,
            style: TextStyle(
                fontFamily: 'Manrope',
                fontWeight: FontWeight.w700,
                color: theme.colorScheme.onSurface),
          ),
          bottom: TabBar(
            labelColor: theme.colorScheme.tertiary,
            unselectedLabelColor: palette.textSecondary,
            indicatorColor: theme.colorScheme.tertiary,
            tabs: [
              Tab(text: l10n.sharedMediaTabMedia),
              Tab(text: l10n.sharedMediaFiles),
              Tab(text: l10n.sharedMediaAudio),
            ],
          ),
        ),
        body: _loading
            ? const Center(child: CircularProgressIndicator())
            : TabBarView(
                children: [
                  _MediaGrid(items: mediaItems),
                  _MediaList(items: media?.documents ?? const [], audio: false),
                  _MediaList(items: media?.audio ?? const [], audio: true),
                ],
              ),
      ),
    );
  }
}

// ─── Image / video grid ───────────────────────────────────────

class _MediaGrid extends StatelessWidget {
  final List<({MediaItem item, bool isVideo})> items;
  const _MediaGrid({required this.items});

  @override
  Widget build(BuildContext context) {
    if (items.isEmpty) return const _EmptyState();
    return GridView.builder(
      padding: const EdgeInsets.all(2),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 3,
        mainAxisSpacing: 2,
        crossAxisSpacing: 2,
      ),
      itemCount: items.length,
      itemBuilder: (_, i) => RepaintBoundary(
        child: _MediaTile(item: items[i].item, isVideo: items[i].isVideo),
      ),
    );
  }
}

class _MediaTile extends StatefulWidget {
  final MediaItem item;
  final bool isVideo;
  const _MediaTile({required this.item, required this.isVideo});

  @override
  State<_MediaTile> createState() => _MediaTileState();
}

class _MediaTileState extends State<_MediaTile> {
  String? _url;

  @override
  void initState() {
    super.initState();
    _resolve();
  }

  Future<void> _resolve() async {
    final res = await sl<MediaRepository>().getMediaUrl(widget.item.mediaId);
    final url = res.fold((_) => null, (u) => u);
    if (mounted) setState(() => _url = url);
  }

  void _open() {
    final url = _url;
    if (url == null) return;
    Navigator.of(context).push(MaterialPageRoute(
      fullscreenDialog: true,
      builder: (_) => widget.isVideo
          ? _VideoScreen(url: url)
          : _ImageScreen(url: url),
    ));
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final url = _url;
    return GestureDetector(
      onTap: url == null ? null : _open,
      child: Container(
        color: theme.colorScheme.surface,
        child: Stack(
          fit: StackFit.expand,
          children: [
            // Images render a thumbnail; videos use a dark tile (the backend
            // sends no preview — first frame would need per-tile decoding).
            if (widget.isVideo)
              Container(color: Colors.black)
            else if (url != null)
              CachedNetworkImage(
                imageUrl: url,
                fit: BoxFit.cover,
                errorWidget: (_, __, ___) =>
                    const Icon(Icons.broken_image_outlined, color: Colors.white54),
              )
            else
              const Center(
                child: SizedBox(
                  width: 20,
                  height: 20,
                  child: CircularProgressIndicator(strokeWidth: 2),
                ),
              ),
            if (widget.isVideo)
              const Center(
                child: Icon(Icons.play_circle_fill_rounded,
                    color: Colors.white, size: 36),
              ),
          ],
        ),
      ),
    );
  }
}

// ─── Files / audio list ───────────────────────────────────────

class _MediaList extends StatelessWidget {
  final List<MediaItem> items;
  final bool audio;
  const _MediaList({required this.items, required this.audio});

  @override
  Widget build(BuildContext context) {
    if (items.isEmpty) return const _EmptyState();
    return ListView.builder(
      padding: const EdgeInsets.symmetric(vertical: 8),
      itemCount: items.length,
      itemBuilder: (_, i) => _MediaFileRow(item: items[i], audio: audio),
    );
  }
}

class _MediaFileRow extends StatefulWidget {
  final MediaItem item;
  final bool audio;
  const _MediaFileRow({required this.item, required this.audio});

  @override
  State<_MediaFileRow> createState() => _MediaFileRowState();
}

class _MediaFileRowState extends State<_MediaFileRow> {
  bool _busy = false;

  Future<void> _open() async {
    if (_busy) return;
    setState(() => _busy = true);
    final l10n = context.l10n;
    final name = widget.item.name?.trim().isNotEmpty == true
        ? widget.item.name!
        : '${widget.item.mediaId}${widget.audio ? '.m4a' : ''}';
    final result =
        await sl<MediaRepository>().downloadToCache(widget.item.mediaId, name);
    if (!mounted) return;
    setState(() => _busy = false);
    await result.fold(
      (_) async => AppFeedback.showError(l10n.fileOpenError),
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
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;
    final title = widget.item.name?.trim().isNotEmpty == true
        ? widget.item.name!
        : (widget.audio ? l10n.messageTypeVoice : l10n.messageTypeFile);

    return ListTile(
      onTap: _open,
      leading: Container(
        width: 40,
        height: 40,
        decoration: BoxDecoration(
          color: palette.surface2,
          borderRadius: BorderRadius.circular(10),
        ),
        child: Icon(
          widget.audio ? Icons.audiotrack_rounded : Icons.insert_drive_file_outlined,
          size: 22,
          color: theme.colorScheme.tertiary,
        ),
      ),
      title: Text(
        title,
        maxLines: 1,
        overflow: TextOverflow.ellipsis,
        style: TextStyle(
            fontFamily: 'Manrope', fontSize: 14, color: theme.colorScheme.onSurface),
      ),
      trailing: _busy
          ? const SizedBox(
              width: 18, height: 18, child: CircularProgressIndicator(strokeWidth: 2))
          : Icon(Icons.download_rounded, color: palette.textSecondary),
    );
  }
}

// ─── Fullscreen viewers ───────────────────────────────────────

class _ImageScreen extends StatelessWidget {
  final String url;
  const _ImageScreen({required this.url});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.black,
      appBar: AppBar(
        backgroundColor: Colors.black,
        iconTheme: const IconThemeData(color: Colors.white),
      ),
      body: PhotoView(imageProvider: CachedNetworkImageProvider(url)),
    );
  }
}

class _VideoScreen extends StatefulWidget {
  final String url;
  const _VideoScreen({required this.url});

  @override
  State<_VideoScreen> createState() => _VideoScreenState();
}

class _VideoScreenState extends State<_VideoScreen> {
  VideoPlayerController? _video;
  ChewieController? _chewie;

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
      _video = ctrl;
      _chewie = ChewieController(
        videoPlayerController: ctrl,
        autoPlay: true,
        looping: false,
      );
    });
  }

  @override
  void dispose() {
    _chewie?.dispose();
    _video?.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.black,
      appBar: AppBar(
        backgroundColor: Colors.black,
        iconTheme: const IconThemeData(color: Colors.white),
      ),
      body: SafeArea(
        child: Center(
          child: _chewie != null
              ? Chewie(controller: _chewie!)
              : const CircularProgressIndicator(),
        ),
      ),
    );
  }
}

// ─── Empty state ──────────────────────────────────────────────

class _EmptyState extends StatelessWidget {
  const _EmptyState();

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    return Center(
      child: Text(
        context.l10n.sharedMediaEmpty,
        style: TextStyle(fontFamily: 'Manrope', color: palette.textSecondary),
      ),
    );
  }
}
