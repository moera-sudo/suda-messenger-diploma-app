import 'package:equatable/equatable.dart';

/// Базовый класс для всех ошибок приложения
abstract class AppFailure extends Equatable {
  final String message; // Человекочитаемое сообщение (дефолтное)
  final String? code;   // Код ошибки (например "USER_NOT_FOUND")

  const AppFailure({required this.message, this.code});

  @override
  List<Object?> get props => [message, code];
}

/// Ошибка с сервера (например, 400 Bad Request, 500 Internal Error)
class ServerFailure extends AppFailure {
  const ServerFailure({required super.message, super.code});
}

/// Ошибка сети (нет интернета, таймаут)
class NetworkFailure extends AppFailure {
  const NetworkFailure() : super(message: 'No Internet Connection', code: 'NETWORK_ERROR');
}

/// Неизвестная ошибка (крэш, null pointer и т.д.)
class UnknownFailure extends AppFailure {
  const UnknownFailure({super.message = 'Unknown Error', super.code = 'UNKNOWN'});
}