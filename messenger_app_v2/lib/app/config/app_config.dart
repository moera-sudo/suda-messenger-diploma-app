import 'package:injectable/injectable.dart';
// import 'const.dart';

@singleton
class AppConfig {
  /// App version shown in Settings footer / About.
  static const String appVersion = '1.0.0';

  String get baseUrl {
    return 'http://172.26.183.121:80'; // * VPN IP Address
  }

  String get PrefixMessenger {
    return 'api/v1/messenger';
  }

  // URL для веб-сокетов (автоматически меняет http на ws)
  String get wsUrl {
    // Frontend.md §4.1: WS hub lives at /ws (NOT /api/v1/messenger/ws)
    const wsPath = '/ws';
  
    final base = baseUrl;
    if (base.startsWith('https')) {
      return base.replaceFirst('https', 'wss') + wsPath;
    }
    return base.replaceFirst('http', 'ws') + wsPath;
  }

  // Таймауты
  Duration get connectTimeout => const Duration(seconds: 15);
  Duration get receiveTimeout => const Duration(seconds: 15);
}
