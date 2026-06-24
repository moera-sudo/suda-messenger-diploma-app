import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';
import 'package:injectable/injectable.dart';
import '../../../app/config/theme/app_theme.dart';
import '../../data/locale/shared_storage.dart';

class SettingsState extends Equatable{
  final AppThemeType themeType;
  final Locale locale;

  const SettingsState({required this.themeType, required this.locale});

  @override
  List<Object> get props => [themeType, locale];

  ThemeData get themeData => AppTheme.getTheme(themeType);
}

@injectable
class SettingsCubit extends Cubit<SettingsState> {
  final PreferencesLocalSource _source;

  SettingsCubit(this._source) : super(const SettingsState(
    themeType: AppThemeType.suda, 
    locale: Locale('en'),
  )) {
    _loadSettings();
  }

  void _loadSettings() {
    final savedTheme = _source.getTheme();
    AppThemeType theme = AppThemeType.suda;
    if (savedTheme != null) {
     try {
      theme = AppThemeType.values.byName(savedTheme);
     } catch(_) {
      // хз
     }
    }

    final savedLocale = _source.getLocale();
    Locale locale = const Locale('en');
    if (savedLocale != null) {
      locale = Locale(savedLocale);
    }

    emit(SettingsState(themeType: theme, locale: locale));
  }

  Future<void> changeTheme(AppThemeType type) async{
    await _source.saveTheme(type.name);
    emit(SettingsState(themeType: type, locale: state.locale));
  }

  Future<void> changeLocale(String languageCode) async {
    await _source.saveLocale(languageCode);
    emit(SettingsState(themeType: state.themeType, locale: Locale(languageCode)));
  }
}