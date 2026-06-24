import 'package:flutter/material.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import '../../../app/config/theme/app_theme.dart';
import '../../../features/media/domain/repositories/media_repository.dart';
import '../../../features/user/domain/repositories/user_repository.dart';

class SudaAvatar extends StatefulWidget {
  /// Pre-resolved image URL. Highest precedence.
  final String? mediaUrl;
  /// Media id to resolve into a view URL (cached process-wide). Used for
  /// group/channel/chat avatars when no pre-resolved [mediaUrl] is available.
  final String? mediaId;
  /// User id to resolve an avatar for (fetches the profile → avatar media id →
  /// url, cached). Lowest-precedence source; used for chat-list interlocutors
  /// and group message senders when no url/media id is on hand.
  final String? userId;
  final String initials;
  final double size;
  final bool showOnline;
  final bool ring;

  const SudaAvatar({
    super.key,
    this.mediaUrl,
    this.mediaId,
    this.userId,
    this.initials = '?',
    this.size = 40,
    this.showOnline = false,
    this.ring = false,
  });

  // Process-wide caches: media_id → url and user_id → url (null = resolved to
  // nothing). Avoids refetching the same avatar repeatedly within a session.
  static final Map<String, String?> _urlCache = {};
  static final Map<String, String?> _userUrlCache = {};

  @override
  State<SudaAvatar> createState() => _SudaAvatarState();
}

class _SudaAvatarState extends State<SudaAvatar> {
  String? _resolvedUrl;

  @override
  void initState() {
    super.initState();
    _resolvedUrl = widget.mediaUrl;
    _resolve();
  }

  @override
  void didUpdateWidget(SudaAvatar old) {
    super.didUpdateWidget(old);
    if (widget.mediaUrl != old.mediaUrl ||
        widget.mediaId != old.mediaId ||
        widget.userId != old.userId) {
      _resolvedUrl = widget.mediaUrl;
      _resolve();
    }
  }

  /// Resolves an avatar URL by precedence: mediaUrl → mediaId → userId.
  Future<void> _resolve() async {
    if (widget.mediaUrl?.isNotEmpty ?? false) return;

    final id = widget.mediaId;
    if (id != null && id.isNotEmpty) {
      if (SudaAvatar._urlCache.containsKey(id)) {
        _applyUrl(SudaAvatar._urlCache[id]);
        return;
      }
      final res = await sl<MediaRepository>().getMediaUrl(id);
      final url = res.fold((_) => null, (u) => u);
      SudaAvatar._urlCache[id] = url;
      _applyUrl(url);
      return;
    }

    final uid = widget.userId;
    if (uid != null && uid.isNotEmpty) {
      if (SudaAvatar._userUrlCache.containsKey(uid)) {
        _applyUrl(SudaAvatar._userUrlCache[uid]);
        return;
      }
      // Resolve via the user's profile: userId → avatar media id → url.
      final profileRes = await sl<UserRepository>().getUserProfile(uid);
      final mediaId = profileRes.fold((_) => null, (p) => p.avatarMediaId);
      String? url;
      if (mediaId != null && mediaId.isNotEmpty) {
        final urlRes = await sl<MediaRepository>().getMediaUrl(mediaId);
        url = urlRes.fold((_) => null, (u) => u);
      }
      SudaAvatar._userUrlCache[uid] = url;
      _applyUrl(url);
    }
  }

  void _applyUrl(String? url) {
    if (url != _resolvedUrl && mounted) setState(() => _resolvedUrl = url);
  }

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);
    final size = widget.size;
    final initials = widget.initials;
    final url = _resolvedUrl;

    Widget avatar = SizedBox(
      width: size,
      height: size,
      child: Stack(
        children: [
          // Circle avatar with clipping
          ClipOval(
            child: SizedBox(
              width: size,
              height: size,
              child: url != null && url.isNotEmpty
                  ? Image.network(
                      url,
                      fit: BoxFit.cover,
                      // Show initials while the network image loads, so there is
                      // no blank frame; fall back to initials on error too.
                      loadingBuilder: (context, child, progress) =>
                          progress == null
                              ? child
                              : _Initials(
                                  initials: initials,
                                  size: size,
                                  palette: palette,
                                  theme: theme,
                                ),
                      errorBuilder: (_, __, ___) => _Initials(
                        initials: initials,
                        size: size,
                        palette: palette,
                        theme: theme,
                      ),
                    )
                  : _Initials(
                      initials: initials,
                      size: size,
                      palette: palette,
                      theme: theme,
                    ),
            ),
          ),
          // Online dot
          if (widget.showOnline)
            Positioned(
              bottom: 0,
              right: 0,
              child: Container(
                width: size * 0.28,
                height: size * 0.28,
                decoration: BoxDecoration(
                  color: palette.success,
                  shape: BoxShape.circle,
                  border: Border.all(
                    color: theme.scaffoldBackgroundColor,
                    width: 2,
                  ),
                ),
              ),
            ),
        ],
      ),
    );

    if (widget.ring) {
      avatar = Container(
        decoration: BoxDecoration(
          shape: BoxShape.circle,
          boxShadow: [
            BoxShadow(
              color: theme.scaffoldBackgroundColor,
              spreadRadius: 3,
              blurRadius: 0,
            ),
            BoxShadow(
              color: palette.surface2,
              spreadRadius: 5,
              blurRadius: 0,
            ),
          ],
        ),
        child: avatar,
      );
    }

    return avatar;
  }
}

class _Initials extends StatelessWidget {
  final String initials;
  final double size;
  final AppPaletteTheme palette;
  final ThemeData theme;

  const _Initials({
    required this.initials,
    required this.size,
    required this.palette,
    required this.theme,
  });

  @override
  Widget build(BuildContext context) {
    final displayInit = initials.isNotEmpty ? initials[0].toUpperCase() : '?';
    return Container(
      width: size,
      height: size,
      decoration: BoxDecoration(
        gradient: LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [
            theme.colorScheme.primary,
            Color.lerp(theme.colorScheme.primary, theme.colorScheme.tertiary, 0.3)!,
          ],
        ),
      ),
      alignment: Alignment.center,
      child: Text(
        displayInit,
        style: TextStyle(
          fontFamily: 'SpaceGrotesk',
          color: Colors.white,
          fontWeight: FontWeight.w700,
          fontSize: size * 0.38,
          letterSpacing: -0.02,
        ),
      ),
    );
  }
}
