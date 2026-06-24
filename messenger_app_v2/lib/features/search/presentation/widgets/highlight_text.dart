import 'package:flutter/material.dart';

import '../../../../app/config/theme/app_theme.dart';

/// Renders [text] with occurrences of [query] highlighted in accent color.
/// Uses pure [TextSpan] (no WidgetSpan/Container) for correct baseline alignment.
class HighlightText extends StatelessWidget {
  final String text;
  final String query;
  final AppPaletteTheme palette;
  final int? maxLines;

  const HighlightText({
    super.key,
    required this.text,
    required this.query,
    required this.palette,
    this.maxLines = 1,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final accent = theme.colorScheme.tertiary;

    final baseStyle = TextStyle(
      fontFamily: 'Manrope',
      fontSize: 13,
      color: palette.textSecondary,
    );

    if (query.isEmpty) {
      return Text(
        text,
        style: baseStyle,
        maxLines: maxLines,
        overflow: TextOverflow.ellipsis,
      );
    }

    final pattern = RegExp(RegExp.escape(query), caseSensitive: false);
    final spans = <TextSpan>[];
    int lastEnd = 0;

    for (final match in pattern.allMatches(text)) {
      if (match.start > lastEnd) {
        spans.add(TextSpan(text: text.substring(lastEnd, match.start), style: baseStyle));
      }
      spans.add(TextSpan(
        text: match.group(0),
        style: TextStyle(
          fontFamily: 'Manrope',
          fontSize: 13,
          fontWeight: FontWeight.w700,
          color: accent,
          backgroundColor: accent.withValues(alpha: 0.15),
        ),
      ));
      lastEnd = match.end;
    }
    if (lastEnd < text.length) {
      spans.add(TextSpan(text: text.substring(lastEnd), style: baseStyle));
    }

    return Text.rich(
      TextSpan(children: spans),
      maxLines: maxLines,
      overflow: TextOverflow.ellipsis,
    );
  }
}
