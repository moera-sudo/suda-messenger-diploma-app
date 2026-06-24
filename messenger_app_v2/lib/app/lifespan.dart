import 'package:flutter/material.dart';
import '../shared/data/logger/app_logger.dart';

class Lifespan extends WidgetsBindingObserver {
  final AppLogger _logger;

  Lifespan(this._logger);

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    _logger.info("App Lifecycle state: $state");

    if (state == AppLifecycleState.detached) {
      _logger.warning("Logger connection is closing..");
      _logger.close();
    }
  }
}