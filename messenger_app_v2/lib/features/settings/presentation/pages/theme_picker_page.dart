import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import '../../../../app/config/theme/app_theme.dart';
import '../../../../features/preferences/domain/repositories/preferences_repository.dart';
import '../../../../shared/presentation/bloc/settings_cubit.dart';
import '../../../../shared/presentation/l10n_extensions.dart';

/// Human-readable name for a theme. Shared between Settings and the picker.
String themeLabel(AppThemeType type) => switch (type) {
      AppThemeType.suda => 'Suda Dark',
      AppThemeType.sudaEnlightened => 'Suda Light',
      AppThemeType.teaChatsLight => 'TeaChats Light',
      AppThemeType.teaChatsDark => 'TeaChats Dark',
      AppThemeType.etherealSky => 'Ethereal Sky',
      AppThemeType.etherealAbyss => 'Ethereal Abyss',
    };

class ThemePickerPage extends StatelessWidget {
  const ThemePickerPage({super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      appBar: AppBar(
        backgroundColor: theme.scaffoldBackgroundColor,
        elevation: 0,
        leading: IconButton(
          icon: Icon(Icons.arrow_back_rounded, color: theme.colorScheme.onSurface),
          onPressed: () => context.pop(),
        ),
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(
              l10n.themeSelectTitle,
              style: TextStyle(fontFamily: 'Manrope', fontSize: 17, fontWeight: FontWeight.w700, color: theme.colorScheme.onSurface),
            ),
            Text(
              l10n.themePickerSubtitle,
              style: TextStyle(fontFamily: 'Manrope', fontSize: 12, fontWeight: FontWeight.w400, color: palette.textSecondary),
            ),
          ],
        ),
      ),
      body: BlocBuilder<SettingsCubit, SettingsState>(
        buildWhen: (p, c) => p.themeType != c.themeType,
        builder: (context, state) {
          return GridView.count(
            padding: const EdgeInsets.all(14),
            crossAxisCount: 2,
            crossAxisSpacing: 12,
            mainAxisSpacing: 12,
            childAspectRatio: 0.82,
            children: AppThemeType.values.map((type) {
              return _ThemeCard(
                type: type,
                selected: type == state.themeType,
                onTap: () {
                  context.read<SettingsCubit>().changeTheme(type);
                  // Sync to backend so the theme restores on new devices.
                  sl<PreferencesRepository>()
                      .updatePreferences({'theme': type.name});
                },
              );
            }).toList(),
          );
        },
      ),
    );
  }
}

class _ThemeCard extends StatelessWidget {
  final AppThemeType type;
  final bool selected;
  final VoidCallback onTap;

  const _ThemeCard({required this.type, required this.selected, required this.onTap});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final accent = theme.colorScheme.tertiary;

    return Material(
      color: theme.colorScheme.surface,
      borderRadius: BorderRadius.circular(16),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(16),
        child: Container(
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(16),
            border: Border.all(
              color: selected ? accent : palette.divider,
              width: selected ? 2 : 1,
            ),
          ),
          clipBehavior: Clip.hardEdge,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Expanded(child: RepaintBoundary(child: _ThemePreview(type: type))),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
                child: Row(
                  children: [
                    Expanded(
                      child: Text(
                        themeLabel(type),
                        style: TextStyle(
                          fontFamily: 'Manrope',
                          fontSize: 13,
                          fontWeight: FontWeight.w700,
                          color: theme.colorScheme.onSurface,
                        ),
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                      ),
                    ),
                    Container(
                      width: 18,
                      height: 18,
                      decoration: BoxDecoration(
                        shape: BoxShape.circle,
                        border: Border.all(
                          color: selected ? accent : palette.textSecondary,
                          width: 2,
                        ),
                        color: selected ? accent : Colors.transparent,
                      ),
                      child: selected
                          ? Icon(Icons.check_rounded, size: 12, color: theme.scaffoldBackgroundColor)
                          : null,
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

/// Miniature chat preview rendered with the target theme's palette.
class _ThemePreview extends StatelessWidget {
  final AppThemeType type;
  const _ThemePreview({required this.type});

  @override
  Widget build(BuildContext context) {
    final td = AppTheme.getTheme(type);
    final palette = td.extension<AppPaletteTheme>()!;

    return Container(
      color: td.scaffoldBackgroundColor,
      padding: const EdgeInsets.all(10),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          // Top bar
          Container(
            height: 14,
            decoration: BoxDecoration(
              color: td.colorScheme.surface,
              borderRadius: BorderRadius.circular(4),
            ),
          ),
          const Spacer(),
          // Incoming bubble
          Align(
            alignment: Alignment.centerLeft,
            child: Container(
              width: 60,
              height: 16,
              decoration: BoxDecoration(
                color: palette.messageOtherBg,
                borderRadius: BorderRadius.circular(8),
              ),
            ),
          ),
          const SizedBox(height: 6),
          // Outgoing bubble
          Align(
            alignment: Alignment.centerRight,
            child: Container(
              width: 70,
              height: 16,
              decoration: BoxDecoration(
                color: palette.messageMeBg,
                borderRadius: BorderRadius.circular(8),
              ),
            ),
          ),
          const Spacer(),
          // Accent dots
          Row(
            children: [
              _Dot(color: td.colorScheme.primary),
              const SizedBox(width: 5),
              _Dot(color: td.colorScheme.secondary),
              const SizedBox(width: 5),
              _Dot(color: td.colorScheme.tertiary),
            ],
          ),
        ],
      ),
    );
  }
}

class _Dot extends StatelessWidget {
  final Color color;
  const _Dot({required this.color});

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 10,
      height: 10,
      decoration: BoxDecoration(color: color, shape: BoxShape.circle),
    );
  }
}
