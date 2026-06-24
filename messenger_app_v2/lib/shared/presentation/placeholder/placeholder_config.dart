part of 'placeholder_page.dart';

// Приватный класс конфига, доступен только внутри этой "библиотеки"
class _PlaceholderConfig {
  final IconData icon;
  final Color color;
  final String title;
  final String message;

  _PlaceholderConfig({
    required this.icon,
    required this.color,
    required this.title,
    required this.message,
  });
}

// Функция получения конфига
_PlaceholderConfig _getConfig(BuildContext context, PlaceholderType type) {
  final colors = Theme.of(context).colorScheme;
  // Мы можем использовать context.l10n, так как импорт extension есть в родительском файле
  final l10n = context.l10n; 

  switch (type) {
    case PlaceholderType.error:
      return _PlaceholderConfig(
        icon: Icons.error_outline_rounded,
        color: colors.error,
        title: l10n.placeholderErrorTitle,
        message: l10n.placeholderErrorMessage,
      );
    case PlaceholderType.inProgress:
      return _PlaceholderConfig(
        icon: Icons.construction_rounded,
        color: Colors.orangeAccent,
        title: l10n.placeholderInProgressTitle,
        message: l10n.placeholderInProgressMessage,
      );
    case PlaceholderType.test:
      return _PlaceholderConfig(
        icon: Icons.science_rounded,
        color: Colors.greenAccent,
        title: l10n.placeholderTestTitle,
        message: l10n.placeholderTestMessage,
      );
    case PlaceholderType.noContent:
      return _PlaceholderConfig(
        icon: Icons.inbox_rounded,
        color: colors.onSurface.withValues(alpha:0.5),
        title: l10n.placeholderNoContentTitle,
        message: l10n.placeholderNoContentMessage,
      );
  }
}