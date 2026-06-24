import 'package:flutter/material.dart';

import '../../../app/config/theme/app_theme.dart';

enum SudaButtonVariant { primary, ghost, outline, dangerOutline, link }

enum SudaButtonSize { sm, md, lg }

class SudaButton extends StatelessWidget {
  final String label;
  final VoidCallback? onPressed;
  final SudaButtonVariant variant;
  final SudaButtonSize size;
  final IconData? icon;
  final IconData? iconRight;
  final bool fullWidth;
  final bool loading;

  const SudaButton({
    super.key,
    required this.label,
    this.onPressed,
    this.variant = SudaButtonVariant.primary,
    this.size = SudaButtonSize.md,
    this.icon,
    this.iconRight,
    this.fullWidth = false,
    this.loading = false,
  });

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);
    final isDisabled = onPressed == null || loading;

    final (double height, double radius, double fontSize) = switch (size) {
      SudaButtonSize.sm => (36.0, 10.0, 13.0),
      SudaButtonSize.md => (44.0, 14.0, 14.0),
      SudaButtonSize.lg => (52.0, 16.0, 15.0),
    };

    Widget child = Row(
      mainAxisSize: fullWidth ? MainAxisSize.max : MainAxisSize.min,
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        if (loading)
          SizedBox(
            width: 16,
            height: 16,
            child: CircularProgressIndicator(
              strokeWidth: 2,
              color: _foreground(palette, theme),
            ),
          )
        else ...[
          if (icon != null) ...[
            Icon(icon, size: fontSize + 2, color: _foreground(palette, theme)),
            const SizedBox(width: 6),
          ],
          Text(
            label,
            style: TextStyle(
              fontFamily: 'Manrope',
              fontSize: fontSize,
              fontWeight: FontWeight.w600,
              color: _foreground(palette, theme),
            ),
          ),
          if (iconRight != null) ...[
            const SizedBox(width: 6),
            Icon(iconRight, size: fontSize + 2, color: _foreground(palette, theme)),
          ],
        ],
      ],
    );

    if (variant == SudaButtonVariant.link) {
      return TextButton(
        onPressed: isDisabled ? null : onPressed,
        style: TextButton.styleFrom(
          foregroundColor: theme.colorScheme.tertiary,
          padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 8),
        ),
        child: child,
      );
    }

    final decoration = _decoration(palette, theme, radius, isDisabled);

    return Opacity(
      opacity: isDisabled ? 0.4 : 1.0,
      child: GestureDetector(
        onTap: isDisabled ? null : onPressed,
        child: AnimatedContainer(
          duration: const Duration(milliseconds: 80),
          height: height,
          width: fullWidth ? double.infinity : null,
          padding: EdgeInsets.symmetric(horizontal: size == SudaButtonSize.sm ? 12 : 18),
          decoration: decoration,
          alignment: Alignment.center,
          child: child,
        ),
      ),
    );
  }

  BoxDecoration _decoration(
    AppPaletteTheme palette,
    ThemeData theme,
    double radius,
    bool disabled,
  ) {
    return switch (variant) {
      SudaButtonVariant.primary => BoxDecoration(
          borderRadius: BorderRadius.circular(radius),
          gradient: LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [
              theme.colorScheme.primary,
              Color.lerp(theme.colorScheme.primary, theme.colorScheme.tertiary, 0.2)!,
            ],
          ),
          boxShadow: disabled
              ? null
              : [
                  BoxShadow(
                    color: theme.colorScheme.primary.withValues(alpha: 0.4),
                    blurRadius: 24,
                    offset: const Offset(0, 6),
                    spreadRadius: -8,
                  ),
                ],
        ),
      SudaButtonVariant.ghost => BoxDecoration(
          borderRadius: BorderRadius.circular(radius),
          color: theme.colorScheme.surface,
          border: Border.all(color: palette.divider),
        ),
      SudaButtonVariant.outline => BoxDecoration(
          borderRadius: BorderRadius.circular(radius),
          border: Border.all(color: theme.colorScheme.tertiary, width: 1.5),
        ),
      SudaButtonVariant.dangerOutline => BoxDecoration(
          borderRadius: BorderRadius.circular(radius),
          border: Border.all(
            color: palette.danger.withValues(alpha: 0.5),
            width: 1.5,
          ),
        ),
      SudaButtonVariant.link => const BoxDecoration(),
    };
  }

  Color _foreground(AppPaletteTheme palette, ThemeData theme) {
    return switch (variant) {
      SudaButtonVariant.primary => Colors.white,
      SudaButtonVariant.ghost => theme.colorScheme.onSurface,
      SudaButtonVariant.outline => theme.colorScheme.tertiary,
      SudaButtonVariant.dangerOutline => palette.danger,
      SudaButtonVariant.link => theme.colorScheme.tertiary,
    };
  }
}
