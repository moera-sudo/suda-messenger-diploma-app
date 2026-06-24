import 'package:injectable/injectable.dart';
import 'package:shared_preferences/shared_preferences.dart';

@lazySingleton
class PreferencesLocalSource {
  final SharedPreferences _prefs;

  PreferencesLocalSource(this._prefs);

  static const _themeKey = 'app_theme';
  static const _localKey = 'app_locale';

  Future<void> saveTheme(String themeName) => _prefs.setString(_themeKey, themeName);

  String? getTheme() => _prefs.getString(_themeKey);

  Future<void> saveLocale(String languageCode) => _prefs.setString(_localKey, languageCode);

  String? getLocale() => _prefs.getString(_localKey);
}