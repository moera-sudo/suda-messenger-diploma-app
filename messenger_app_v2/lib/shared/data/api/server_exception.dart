class ServerException implements Exception {
  final String errorCode; // "BAD_REQUEST"
  final String message;   // "Invalid Input"
  final int statusCode;   // 400

  ServerException({
    required this.errorCode,
    required this.message,
    required this.statusCode,
  });

  @override
  String toString() => 'ServerException(code: $statusCode, error: $errorCode, msg: $message)';
}