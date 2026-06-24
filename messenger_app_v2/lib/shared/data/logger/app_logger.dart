import 'package:logger/logger.dart';
import 'package:injectable/injectable.dart';

@lazySingleton
class AppLogger {
  final Logger _logger = Logger(
    printer: PrettyPrinter(
      methodCount: 1, 
      errorMethodCount: 8,
      lineLength: 80,
      colors: true,
      printEmojis: true,
      dateTimeFormat: DateTimeFormat.onlyTimeAndSinceStart,
    ),
  );

  void trace(String message) => _logger.t(message);
  void debug(String message) => _logger.d(message);
  void info(String message) => _logger.i(message);
  void warning(String message) => _logger.w(message);
  void error(String message, [dynamic error, StackTrace? stackTrace]) =>
      _logger.e(message, error: error, stackTrace: stackTrace);
  void fatal(String message) => _logger.f(message);

  Future<void> close() => _logger.close();
}