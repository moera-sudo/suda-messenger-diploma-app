import 'dart:async';
import 'dart:convert';

import 'package:injectable/injectable.dart';
import 'package:web_socket_channel/io.dart';
import 'package:web_socket_channel/status.dart' as ws_status;

import '../../../app/config/app_config.dart';
import '../logger/app_logger.dart';
import '../storage/secure_storage_client.dart';
import '../../domain/models/socket_event.dart';

@lazySingleton
class SocketClient {
  final SecureStorageClient _storage;
  final AppLogger _logger;
  final AppConfig _config;

  IOWebSocketChannel? _channel;
  final _eventController = StreamController<SocketEvent>.broadcast();

  // Exponential backoff: 1s → 2s → 5s → 10s → 30s → 30s (cap)
  static const _backoffSeconds = [1, 2, 5, 10, 30, 30];
  int _retryCount = 0;
  bool _isConnected = false;
  bool _isConnecting = false;
  bool _intentionalClose = false;
  Timer? _reconnectTimer;

  SocketClient(this._storage, this._logger, this._config);

  Stream<SocketEvent> get events => _eventController.stream;
  bool get isConnected => _isConnected;

  Future<void> connect() async {
    if (_isConnecting || _isConnected) return;
    // Allow reconnecting after max retries have been exhausted by an external caller
    if (_retryCount >= _backoffSeconds.length) _retryCount = 0;
    _isConnecting = true;
    _intentionalClose = false;

    final token = await _storage.getAccessToken();
    if (token == null) {
      _logger.warning('WS connect skipped: no access token');
      _isConnecting = false;
      return;
    }

    // Use Uri.replace to avoid empty-fragment artifact (#) from string concatenation
    final wsUrl = Uri.parse(_config.wsUrl).replace(
      queryParameters: {'token': token},
    );
    _logger.info('🔌 Connecting to WS: ${_config.wsUrl}');

    try {
      final channel = IOWebSocketChannel.connect(wsUrl);
      await channel.ready; // throws if handshake fails

      _channel = channel;
      _isConnected = true;
      _isConnecting = false;
      _retryCount = 0;
      _logger.info('✅ WS Connected');

      _channel!.stream.listen(
        _onMessage,
        onError: (e) {
          _logger.error('WS stream error', e);
          _handleDisconnect();
        },
        onDone: _handleDisconnect,
        cancelOnError: true,
      );
    } catch (e) {
      _logger.error('WS connection failed', e);
      _isConnecting = false;
      _channel = null;
      _handleDisconnect();
    }
  }

  void send(String type, Map<String, dynamic> payload) {
    if (_channel == null || !_isConnected) {
      _logger.warning('Cannot send "$type": WS not connected');
      return;
    }
    final event = {'type': type, 'payload': payload};
    _channel!.sink.add(jsonEncode(event));
    _logger.debug('📤 WS sent: $type');
  }

  void disconnect() {
    _intentionalClose = true;
    _reconnectTimer?.cancel();
    _reconnectTimer = null;
    _channel?.sink.close(ws_status.goingAway);
    _channel = null;
    _isConnected = false;
    _isConnecting = false;
    _logger.info('🔌 WS disconnected intentionally');
  }

  void _onMessage(dynamic raw) {
    try {
      _logger.debug('📩 WS raw: $raw');
      final jsonMap = jsonDecode(raw as String) as Map<String, dynamic>;
      final event = SocketEvent.fromJson(jsonMap);
      _eventController.add(event);
    } catch (e) {
      _logger.error('WS parse error', e);
    }
  }

  void _handleDisconnect() {
    if (_intentionalClose) return;
    _isConnected = false;
    _isConnecting = false;
    _channel = null;
    _scheduleReconnect();
  }

  void _scheduleReconnect() {
    _reconnectTimer?.cancel();

    if (_retryCount >= _backoffSeconds.length) {
      _logger.warning('WS max retries (${_backoffSeconds.length}) reached — giving up');
      return;
    }

    final delay = _backoffSeconds[_retryCount];
    _retryCount++;
    _logger.info('🔄 WS reconnect in ${delay}s (attempt $_retryCount/${_backoffSeconds.length})');
    _reconnectTimer = Timer(Duration(seconds: delay), connect);
  }
}
