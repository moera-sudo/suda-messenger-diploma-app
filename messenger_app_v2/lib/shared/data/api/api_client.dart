import 'dart:developer' as developer;

import 'package:dio/dio.dart';
import 'package:injectable/injectable.dart';

import '../../domain/models/error_dictionary.dart';
import '../../presentation/feedback/app_feedback.dart';
import 'server_exception.dart';

@lazySingleton
class ApiClient {
  final Dio _dio;

  ApiClient(this._dio);

  Future<dynamic> get(String path, {Map<String, dynamic>? query}) async {
    try {
      final response = await _dio.get(path, queryParameters: query);
      return response.data;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<dynamic> post(String path, {dynamic data, Options? options}) async {
    try {
      final response = await _dio.post(path, data: data, options: options);
      return response.data;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<dynamic> put(String path, {dynamic data}) async {
    try {
      final response = await _dio.put(path, data: data);
      return response.data;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<dynamic> delete(String path, {Map<String, dynamic>? queryParameters}) async {
    try {
      final response = await _dio.delete(path, queryParameters: queryParameters);
      return response.data;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Exception _handleError(DioException e) {
    _logDioError(e);

    if (e.error is ServerException) {
      final serverError = e.error as ServerException;
      final userMessage = _shortMessage(
        statusCode: serverError.statusCode,
        errorCode: serverError.errorCode,
        backendMessage: serverError.message,
      );
      // Token-gating is handled entirely by the paywall UI — never toast it.
      if (serverError.errorCode != 'GATING_REQUIRED') {
        AppFeedback.showError(userMessage);
      }
      return ServerException(
        errorCode: serverError.errorCode,
        statusCode: serverError.statusCode,
        message: userMessage,
      );
    }

    if (e.type == DioExceptionType.connectionTimeout ||
        e.type == DioExceptionType.receiveTimeout ||
        e.type == DioExceptionType.connectionError ||
        e.type == DioExceptionType.sendTimeout) {
      const userMessage = 'No internet connection';
      AppFeedback.showError(userMessage);
      return ServerException(
        errorCode: 'NETWORK_ERROR',
        message: userMessage,
        statusCode: 0,
      );
    }

    final statusCode = e.response?.statusCode ?? 500;
    final errorCode = _extractErrorCode(e.response?.data);
    final backendMessage = _extractBackendMessage(e.response?.data);

    final userMessage = _shortMessage(
      statusCode: statusCode,
      errorCode: errorCode,
      backendMessage: backendMessage,
    );
    // Token-gating is handled entirely by the paywall UI — never toast it.
    if (errorCode != 'GATING_REQUIRED') {
      AppFeedback.showError(userMessage);
    }

    return ServerException(
      errorCode: errorCode ?? 'UNKNOWN',
      statusCode: statusCode,
      message: userMessage,
    );
  }

  String _shortMessage({
    required int statusCode,
    required String? errorCode,
    required String? backendMessage,
  }) {
    if (errorCode == 'NETWORK_ERROR' || statusCode == 0) {
      return 'No internet connection';
    }
    return ErrorDictionary.humanMessage(
      statusCode: statusCode,
      errorKey: errorCode,
    );
  }

  String? _extractErrorCode(dynamic responseData) {
    if (responseData is Map && responseData['error'] != null) {
      return responseData['error'].toString();
    }
    return null;
  }

  String? _extractBackendMessage(dynamic responseData) {
    if (responseData is Map && responseData['message'] != null) {
      return responseData['message'].toString();
    }
    if (responseData is String) {
      return responseData;
    }
    return null;
  }

  void _logDioError(DioException e) {
    final request = e.requestOptions;
    final response = e.response;

    developer.log(
      'HTTP ${request.method} ${request.uri} failed. '
      'status=${response?.statusCode}, response=${response?.data}, '
      'dioMessage=${e.message}, type=${e.type}',
      name: 'ApiClient',
      error: e,
      stackTrace: e.stackTrace,
    );
  }
}
