import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../../../shared/presentation/bloc/settings_cubit.dart';
import '../../../../shared/presentation/l10n_extensions.dart';

class LanguageSelector extends StatelessWidget {
  const LanguageSelector({super.key});

  @override
  Widget build(BuildContext context) {
    // Получаем переводы для названий языков
    final l10n = context.l10n;

    return BlocBuilder<SettingsCubit, SettingsState>(
      builder: (context, state) {
        return PopupMenuButton<String>(
          icon: Container(
            padding: const EdgeInsets.all(8),
            decoration: BoxDecoration(
              color: Colors.white.withValues(alpha: 0.12),
              shape: BoxShape.circle,
            ),
            child: const Icon(Icons.language, color: Colors.white),
          ),
          tooltip: l10n.languageTitle,
          initialValue: state.locale.languageCode,
          onSelected: (code) {
            context.read<SettingsCubit>().changeLocale(code);
          },
          itemBuilder: (context) => [
            PopupMenuItem(
              value: 'en',
              child: Row(
                children: [
                  const Text("🇺🇸 ", style: TextStyle(fontSize: 18)),
                  Text(l10n.languageEnglish), // "English"
                ],
              ),
            ),
            PopupMenuItem(
              value: 'ru',
              child: Row(
                children: [
                  const Text("🇷🇺 ", style: TextStyle(fontSize: 18)),
                  Text(l10n.languageRussian), // "Русский"
                ],
              ),
            ),
            PopupMenuItem(
              value: 'kk',
              child: Row(
                children: [
                  const Text("🇰🇿 ", style: TextStyle(fontSize: 18)),
                  Text(l10n.languageKazakh), // "Қазақша"
                ],
              ),
            ),
          ],
        );
      },
    );
  }
}
