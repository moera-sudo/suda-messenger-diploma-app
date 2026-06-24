import 'package:flutter/material.dart';
import 'package:webview_flutter/webview_flutter.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import '../../../../shared/data/api/api_client.dart';
import '../../../../shared/data/logger/app_logger.dart';
import '../../../../shared/presentation/l10n_extensions.dart';

/// Generic Mini App host screen.
/// Fetches initData via POST /user/init-data, then opens the SPA in a WebView.
/// Both Wallet and Marketplace reuse this screen — they differ only by [appUrl].
class MiniAppWebViewScreen extends StatefulWidget {
  final String appUrl;
  final String appName;

  const MiniAppWebViewScreen({
    super.key,
    required this.appUrl,
    required this.appName,
  });

  @override
  State<MiniAppWebViewScreen> createState() => _MiniAppWebViewScreenState();
}

class _MiniAppWebViewScreenState extends State<MiniAppWebViewScreen> {
  // Controller is created synchronously in initState so WebViewWidget can be
  // included in the widget tree from the very first build, avoiding the
  // ImageReader_JNI buffer race condition that occurs when the widget appears
  // mid-frame after an async gap.
  late final WebViewController _controller;
  bool _loading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _controller = WebViewController()
      ..setJavaScriptMode(JavaScriptMode.unrestricted)
      // setBackgroundColor is omitted intentionally: calling it forces an early
      // WebView redraw before the Chromium render pipeline is ready, which adds
      // an extra ImageReader_JNI buffer-acquisition warning in logcat.
      // The dark scaffold background behind the WebView serves as visual cover.
      ..setNavigationDelegate(NavigationDelegate(
        onNavigationRequest: (req) {
          // SPA may redirect to close://wallet to signal it wants to close.
          if (req.url.startsWith('close://')) {
            if (mounted) Navigator.of(context).maybePop();
            return NavigationDecision.prevent;
          }
          return NavigationDecision.navigate;
        },
        onPageStarted: (_) {
          if (mounted) setState(() => _loading = true);
        },
        onPageFinished: (_) {
          if (mounted) setState(() => _loading = false);
        },
        onWebResourceError: (err) {
          if ((err.isForMainFrame ?? false) && mounted) {
            setState(() {
              _error = err.description;
              _loading = false;
            });
          }
        },
      ));
    _bootstrap();
  }

  Future<void> _bootstrap() async {
    if (mounted) setState(() { _loading = true; _error = null; });

    try {
      final api    = sl<ApiClient>();
      final logger = sl<AppLogger>();

      // 1. Obtain HMAC-signed initData from backend.
      final res      = await api.post('/api/v1/messenger/user/init-data', data: {});
      final initData = res['initData'] as String;
      logger.debug('init-data obtained for ${widget.appName}');

      // 2. Build URL — trailing slash is required to avoid redirect loop in the SPA.
      final encoded = Uri.encodeQueryComponent(initData);
      final fullUrl = '${widget.appUrl}?initData=$encoded';

      // 3. Load into the already-initialised controller (no setState needed here —
      //    onPageStarted/onPageFinished callbacks will drive the loading indicator).
      await _controller.loadRequest(Uri.parse(fullUrl));
    } catch (e) {
      sl<AppLogger>().error('MiniAppWebView bootstrap failed', e);
      if (mounted) setState(() { _error = e.toString(); _loading = false; });
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;

    return Scaffold(
      backgroundColor: const Color(0xFF0F1020),
      body: SafeArea(
        child: Stack(
          children: [
            // WebViewWidget is always in the tree from the first build so the
            // Android SurfaceTexture is allocated once, preventing the
            // ImageReader_JNI "unable to acquire buffer" warning.
            WebViewWidget(controller: _controller),

            // Loading overlay
            if (_loading)
              const Center(
                child: CircularProgressIndicator(color: Colors.white),
              ),

            // Error overlay — shown on top of the WebView so the widget stays
            // mounted and Retry just calls _bootstrap() without rebuilding.
            if (_error != null)
              ColoredBox(
                color: const Color(0xFF0F1020),
                child: Center(
                  child: Padding(
                    padding: const EdgeInsets.all(24),
                    child: Column(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        const Icon(Icons.error_outline_rounded,
                            color: Colors.redAccent, size: 48),
                        const SizedBox(height: 16),
                        Text(
                          l10n.walletLoadingError,
                          textAlign: TextAlign.center,
                          style: const TextStyle(
                              fontFamily: 'Manrope',
                              color: Colors.white70,
                              fontSize: 15),
                        ),
                        const SizedBox(height: 20),
                        ElevatedButton(
                          onPressed: _bootstrap,
                          child: const Text('Retry'),
                        ),
                      ],
                    ),
                  ),
                ),
              ),

            // Back button (always visible so user can leave on error)
            Positioned(
              top: 0,
              left: 0,
              child: IconButton(
                icon: const Icon(Icons.arrow_back_rounded, color: Colors.white),
                onPressed: () => Navigator.of(context).maybePop(),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
