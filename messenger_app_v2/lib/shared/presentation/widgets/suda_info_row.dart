import 'package:flutter/material.dart';

import '../../../app/config/theme/app_theme.dart';

/// A settings/info list row with icon, label and optional right widget.
/// Wrap multiple [SudaInfoRow]s in a [SudaInfoList] for proper card styling.
class SudaInfoRow extends StatelessWidget {
  final IconData icon;
  final String label;
  final Widget? right;
  final VoidCallback? onTap;
  final bool danger;
  final bool showChevron;

  const SudaInfoRow({
    super.key,
    required this.icon,
    required this.label,
    this.right,
    this.onTap,
    this.danger = false,
    this.showChevron = true,
  });

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);
    final color = danger ? palette.danger : theme.colorScheme.onSurface;

    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.zero,
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 14),
        child: Row(
          children: [
            Icon(
              icon,
              size: 20,
              color: danger ? palette.danger : palette.textAccent,
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Text(
                label,
                style: TextStyle(
                  fontFamily: 'Manrope',
                  fontSize: 14,
                  fontWeight: FontWeight.w500,
                  color: color,
                ),
              ),
            ),
            if (right != null) right!,
            if (onTap != null && showChevron && right == null)
              Icon(Icons.chevron_right, size: 18, color: palette.textSecondary),
            if (onTap != null && showChevron && right != null)
              const SizedBox(width: 4),
            if (onTap != null && showChevron && right != null)
              Icon(Icons.chevron_right, size: 18, color: palette.textSecondary),
          ],
        ),
      ),
    );
  }
}

/// Card container wrapping a list of [SudaInfoRow]s.
class SudaInfoList extends StatelessWidget {
  final List<Widget> children;
  final String? title;
  final EdgeInsets? padding;

  const SudaInfoList({
    super.key,
    required this.children,
    this.title,
    this.padding,
  });

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);

    return Padding(
      padding: padding ?? const EdgeInsets.symmetric(horizontal: 14),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          if (title != null) ...[
            Padding(
              padding: const EdgeInsets.only(bottom: 8, left: 2, top: 4),
              child: Text(
                title!.toUpperCase(),
                style: TextStyle(
                  fontFamily: 'Manrope',
                  fontSize: 12,
                  fontWeight: FontWeight.w700,
                  color: palette.textSecondary,
                  letterSpacing: 0.10,
                ),
              ),
            ),
          ],
          Container(
            decoration: BoxDecoration(
              color: theme.colorScheme.surface,
              borderRadius: BorderRadius.circular(16),
              border: Border.all(color: palette.divider),
            ),
            clipBehavior: Clip.hardEdge,
            child: Column(
              children: [
                for (int i = 0; i < children.length; i++) ...[
                  children[i],
                  if (i < children.length - 1)
                    Divider(height: 1, color: palette.divider, indent: 46),
                ],
              ],
            ),
          ),
        ],
      ),
    );
  }
}

/// Small muted metadata text for the right side of a [SudaInfoRow].
class InfoMeta extends StatelessWidget {
  final String text;

  const InfoMeta(this.text, {super.key});

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    return Text(
      text,
      style: TextStyle(
        fontFamily: 'Manrope',
        fontSize: 13,
        color: palette.textSecondary,
      ),
    );
  }
}
