import 'package:dio/dio.dart';
import 'package:injectable/injectable.dart';
import 'package:pretty_dio_logger/pretty_dio_logger.dart';
import 'auth_interceptor.dart';
import 'package:messenger_app_v2/app/config/app_config.dart';

@module
abstract class DioModule {
  @lazySingleton
  Dio dio(
    AppConfig config,
    AuthInterceptor authInterceptor,
  ) {
    final dio = Dio(BaseOptions(
      baseUrl: config.baseUrl,
      connectTimeout: config.connectTimeout,
      receiveTimeout: config.receiveTimeout,
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
    ));

    dio.interceptors.add(authInterceptor);

    dio.interceptors.add(PrettyDioLogger(
      requestHeader: true,
      requestBody: true,
      responseBody: true,
      compact: true,
      maxWidth: 90,
    ));

    return dio;
  }
}