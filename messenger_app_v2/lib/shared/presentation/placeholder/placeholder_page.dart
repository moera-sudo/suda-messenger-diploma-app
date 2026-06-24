import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import '../l10n_extensions.dart';

part 'placeholder_config.dart';

enum PlaceholderType {
  error,
  inProgress,
  test,
  noContent,
}

/// Generic placeholder widget.
///
/// When [scaffold] is true (default) the widget is a full [Scaffold] with its
/// own [AppBar].  Set [scaffold] to **false** when embedding inside a parent
/// [Scaffold] body — otherwise the parent's back-button and this widget's
/// AppBar stack on top of each other (double arrow).
class PlaceholderPage extends StatelessWidget {
  final PlaceholderType type;
  final String? title;
  final String? message;
  final VoidCallback? onRetry;
  final bool scaffold;

  const PlaceholderPage({
    super.key,
    required this.type,
    this.title,
    this.message,
    this.onRetry,
    this.scaffold = true,
  });

  @override
  Widget build(BuildContext context) {
    final config = _getConfig(context, type);
    final theme = Theme.of(context);
    final body = Center(
      child: Padding(
        padding: const EdgeInsets.all(32.0),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(config.icon, size: 80, color: config.color)
                .animate()
                .scale(duration: 600.ms, curve: Curves.easeOutBack)
                .fadeIn(),
            const SizedBox(height: 32),
            Text(
              message ?? config.message,
              textAlign: TextAlign.center,
              style: theme.textTheme.headlineSmall?.copyWith(
                fontWeight: FontWeight.bold,
                color: theme.colorScheme.onSurface,
              ),
            ).animate().fadeIn(delay: 200.ms).slideY(begin: 0.2, end: 0),
            const SizedBox(height: 16),
            if (type == PlaceholderType.inProgress)
              Padding(
                padding: const EdgeInsets.only(top: 24.0),
                child: CircularProgressIndicator(color: config.color),
              ).animate().fadeIn(delay: 400.ms),
            if (type == PlaceholderType.error && onRetry != null)
              Padding(
                padding: const EdgeInsets.only(top: 32.0),
                child: FilledButton.icon(
                  onPressed: onRetry,
                  style: FilledButton.styleFrom(
                    backgroundColor: theme.colorScheme.error,
                    foregroundColor: theme.colorScheme.onError,
                  ),
                  icon: const Icon(Icons.refresh),
                  label: Text(context.l10n.buttonRetry),
                ),
              ).animate().fadeIn(delay: 400.ms),
            const SizedBox(height: 64),
          ],
        ),
      ),
    );

    if (!scaffold) return body;

    return Scaffold(
      appBar: AppBar(
        title: Text(title ?? config.title),
        centerTitle: true,
        backgroundColor: Colors.transparent,
      ),
      body: body,
    );
  }
}
