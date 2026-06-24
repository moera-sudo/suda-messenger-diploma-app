import 'package:flutter/material.dart';

import 'app_colors.dart';

enum AppThemeType {
  suda,
  sudaEnlightened,
  teaChatsLight,
  teaChatsDark,
  etherealSky,
  etherealAbyss,
}

/// ThemeExtension exposing palette-specific colors not covered by ColorScheme.
class AppPaletteTheme extends ThemeExtension<AppPaletteTheme> {
  final Color surfaceContainer;
  final Color surface2;
  final Color textSecondary;
  final Color divider;
  final Color messageMeBg;
  final Color messageOtherBg;
  final Color textAccent;
  final Color danger;
  final Color success;
  final Color glow;

  const AppPaletteTheme({
    required this.surfaceContainer,
    required this.surface2,
    required this.textSecondary,
    required this.divider,
    required this.messageMeBg,
    required this.messageOtherBg,
    required this.textAccent,
    required this.danger,
    required this.success,
    required this.glow,
  });

  @override
  AppPaletteTheme copyWith({
    Color? surfaceContainer,
    Color? surface2,
    Color? textSecondary,
    Color? divider,
    Color? messageMeBg,
    Color? messageOtherBg,
    Color? textAccent,
    Color? danger,
    Color? success,
    Color? glow,
  }) {
    return AppPaletteTheme(
      surfaceContainer: surfaceContainer ?? this.surfaceContainer,
      surface2: surface2 ?? this.surface2,
      textSecondary: textSecondary ?? this.textSecondary,
      divider: divider ?? this.divider,
      messageMeBg: messageMeBg ?? this.messageMeBg,
      messageOtherBg: messageOtherBg ?? this.messageOtherBg,
      textAccent: textAccent ?? this.textAccent,
      danger: danger ?? this.danger,
      success: success ?? this.success,
      glow: glow ?? this.glow,
    );
  }

  @override
  AppPaletteTheme lerp(covariant ThemeExtension<AppPaletteTheme>? other, double t) {
    if (other is! AppPaletteTheme) return this;
    return AppPaletteTheme(
      surfaceContainer: Color.lerp(surfaceContainer, other.surfaceContainer, t)!,
      surface2: Color.lerp(surface2, other.surface2, t)!,
      textSecondary: Color.lerp(textSecondary, other.textSecondary, t)!,
      divider: Color.lerp(divider, other.divider, t)!,
      messageMeBg: Color.lerp(messageMeBg, other.messageMeBg, t)!,
      messageOtherBg: Color.lerp(messageOtherBg, other.messageOtherBg, t)!,
      textAccent: Color.lerp(textAccent, other.textAccent, t)!,
      danger: Color.lerp(danger, other.danger, t)!,
      success: Color.lerp(success, other.success, t)!,
      glow: Color.lerp(glow, other.glow, t)!,
    );
  }
}

class AppTheme {
  static ThemeData getTheme(AppThemeType type) {
    final palette = _getPalette(type);
    final isDark = _isDark(type);

    return ThemeData(
      useMaterial3: true,
      brightness: isDark ? Brightness.dark : Brightness.light,
      scaffoldBackgroundColor: palette.background,
      primaryColor: palette.primary,
      fontFamily: 'Manrope',
      colorScheme: ColorScheme(
        brightness: isDark ? Brightness.dark : Brightness.light,
        primary: palette.primary,
        onPrimary: Colors.white,
        secondary: palette.secondary,
        onSecondary: Colors.white,
        tertiary: palette.accent,
        onTertiary: palette.background,
        surface: palette.surface,
        onSurface: palette.textPrimary,
        surfaceContainerHighest: palette.surfaceContainer,
        error: palette.danger,
        onError: Colors.white,
      ),
      dividerColor: palette.divider,
      textTheme: TextTheme(
        displayLarge: TextStyle(
          fontFamily: 'SpaceGrotesk',
          color: palette.textPrimary,
          fontWeight: FontWeight.w800,
        ),
        displayMedium: TextStyle(
          fontFamily: 'SpaceGrotesk',
          color: palette.textPrimary,
          fontWeight: FontWeight.w700,
        ),
        displaySmall: TextStyle(
          fontFamily: 'SpaceGrotesk',
          color: palette.textPrimary,
          fontWeight: FontWeight.w700,
        ),
        headlineLarge: TextStyle(
          fontFamily: 'SpaceGrotesk',
          color: palette.textPrimary,
          fontWeight: FontWeight.w800,
          letterSpacing: -0.03,
        ),
        headlineMedium: TextStyle(
          fontFamily: 'SpaceGrotesk',
          color: palette.textPrimary,
          fontWeight: FontWeight.w700,
          letterSpacing: -0.02,
        ),
        titleLarge: TextStyle(
          fontFamily: 'Manrope',
          color: palette.textPrimary,
          fontWeight: FontWeight.w700,
          fontSize: 18,
        ),
        titleMedium: TextStyle(
          fontFamily: 'Manrope',
          color: palette.textPrimary,
          fontWeight: FontWeight.w600,
          fontSize: 15,
        ),
        titleSmall: TextStyle(
          fontFamily: 'Manrope',
          color: palette.textAccent,
          fontWeight: FontWeight.w600,
          fontSize: 13,
        ),
        bodyLarge: TextStyle(
          fontFamily: 'Manrope',
          color: palette.textPrimary,
          fontSize: 15,
        ),
        bodyMedium: TextStyle(
          fontFamily: 'Manrope',
          color: palette.textPrimary,
          fontSize: 14,
        ),
        bodySmall: TextStyle(
          fontFamily: 'Manrope',
          color: palette.textSecondary,
          fontSize: 13,
        ),
        labelSmall: TextStyle(
          fontFamily: 'Manrope',
          color: palette.textSecondary,
          fontSize: 11,
          fontWeight: FontWeight.w600,
        ),
      ),
      appBarTheme: AppBarTheme(
        backgroundColor: palette.background,
        foregroundColor: palette.textPrimary,
        elevation: 0,
        scrolledUnderElevation: 0,
        titleTextStyle: TextStyle(
          fontFamily: 'Manrope',
          color: palette.textPrimary,
          fontSize: 18,
          fontWeight: FontWeight.w700,
        ),
        iconTheme: IconThemeData(color: palette.textPrimary),
      ),
      inputDecorationTheme: InputDecorationTheme(
        filled: true,
        fillColor: palette.surface,
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(14),
          borderSide: BorderSide(color: palette.divider),
        ),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(14),
          borderSide: BorderSide(color: palette.divider),
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(14),
          borderSide: BorderSide(color: palette.accent, width: 1.5),
        ),
        hintStyle: TextStyle(color: palette.textSecondary),
        contentPadding: const EdgeInsets.symmetric(horizontal: 14, vertical: 14),
      ),
      extensions: [
        AppPaletteTheme(
          surfaceContainer: palette.surfaceContainer,
          surface2: palette.surface2,
          textSecondary: palette.textSecondary,
          divider: palette.divider,
          messageMeBg: palette.messageMeBg,
          messageOtherBg: palette.messageOtherBg,
          textAccent: palette.textAccent,
          danger: palette.danger,
          success: palette.success,
          glow: palette.glow,
        ),
      ],
    );
  }

  static AppPalette _getPalette(AppThemeType type) {
    switch (type) {
      case AppThemeType.suda:
        return sudaPalette;
      case AppThemeType.sudaEnlightened:
        return sudaEnlightenedPalette;
      case AppThemeType.teaChatsLight:
        return teaChatsLightPalette;
      case AppThemeType.teaChatsDark:
        return teaChatsDarkPalette;
      case AppThemeType.etherealSky:
        return etherealSkyPalette;
      case AppThemeType.etherealAbyss:
        return etherealAbyssPalette;
    }
  }

  static bool _isDark(AppThemeType type) {
    return [
      AppThemeType.suda,
      AppThemeType.teaChatsDark,
      AppThemeType.etherealAbyss,
    ].contains(type);
  }

  /// Returns a TextStyle using SpaceGrotesk (display/heading use cases).
  static TextStyle displayStyle({
    required BuildContext context,
    double fontSize = 22,
    FontWeight fontWeight = FontWeight.w800,
    Color? color,
  }) {
    final theme = Theme.of(context);
    return TextStyle(
      fontFamily: 'SpaceGrotesk',
      fontSize: fontSize,
      fontWeight: fontWeight,
      color: color ?? theme.colorScheme.onSurface,
      letterSpacing: -0.03,
    );
  }

  /// Returns a TextStyle using JetBrains Mono (wallet addresses, timers, amounts).
  static TextStyle monoStyle({
    required BuildContext context,
    double fontSize = 13,
    FontWeight fontWeight = FontWeight.w400,
    Color? color,
  }) {
    final theme = Theme.of(context);
    return TextStyle(
      fontFamily: 'JetBrainsMono',
      fontSize: fontSize,
      fontWeight: fontWeight,
      color: color ?? theme.colorScheme.onSurface,
    );
  }
}
