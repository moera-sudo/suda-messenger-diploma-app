import 'package:flutter/material.dart';

import '../../../app/config/theme/app_theme.dart';

class SudaChip extends StatelessWidget {
  final String label;
  final bool active;
  final VoidCallback? onTap;
  final int? count;

  const SudaChip({
    super.key,
    required this.label,
    this.active = false,
    this.onTap,
    this.count,
  });

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);

    return GestureDetector(
      onTap: onTap,
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 150),
        padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 6),
        decoration: BoxDecoration(
          color: active
              ? Color.lerp(theme.colorScheme.tertiary, theme.colorScheme.surface, 0.88)
              : theme.colorScheme.surface,
          borderRadius: BorderRadius.circular(100),
          border: Border.all(
            color: active ? theme.colorScheme.tertiary : palette.divider,
          ),
        ),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(
              label,
              style: TextStyle(
                fontFamily: 'Manrope',
                fontSize: 13,
                fontWeight: FontWeight.w600,
                color: active ? theme.colorScheme.tertiary : palette.textSecondary,
              ),
            ),
            if (count != null && count! > 0) ...[
              const SizedBox(width: 5),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 5),
                constraints: const BoxConstraints(minWidth: 18),
                height: 18,
                decoration: BoxDecoration(
                  color: theme.colorScheme.primary,
                  borderRadius: BorderRadius.circular(100),
                ),
                alignment: Alignment.center,
                child: Text(
                  '$count',
                  style: const TextStyle(
                    fontFamily: 'Manrope',
                    fontSize: 11,
                    fontWeight: FontWeight.w700,
                    color: Colors.white,
                  ),
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }
}
