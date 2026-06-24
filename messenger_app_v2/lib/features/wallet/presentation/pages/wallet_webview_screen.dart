import 'package:flutter/material.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import 'package:messenger_app_v2/app/config/app_config.dart';
import 'mini_app_webview_screen.dart';

/// Opens the Wallet SPA in a WebView. The SPA lives at baseUrl/wallet/
/// and authenticates via initData (Frontend.md §5.2 / PreDefence-Fix.md §9).
class WalletWebViewScreen extends StatelessWidget {
  const WalletWebViewScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final walletUrl = '${sl<AppConfig>().baseUrl}/wallet/';
    return MiniAppWebViewScreen(appUrl: walletUrl, appName: 'Wallet');
  }
}
