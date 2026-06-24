import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:injectable/injectable.dart';
import 'package:shared_preferences/shared_preferences.dart';

@module
abstract class RegisterModule {
  // --- Storage ---
  @lazySingleton
  FlutterSecureStorage get secureStorage => const FlutterSecureStorage();

  @preResolve
  Future<SharedPreferences> get prefs => SharedPreferences.getInstance();

  // ВАЖНО: Здесь НЕ должно быть ни Dio, ни BaseUrl.
  // Они теперь живут в lib/shared/data/dio/dio_module.dart
}