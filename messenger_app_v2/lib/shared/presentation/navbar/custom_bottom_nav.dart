import 'dart:ui';

import 'package:flutter/material.dart';

import '../../../app/config/theme/app_theme.dart';
import '../l10n_extensions.dart';

class CustomBottomNav extends StatelessWidget {
  final int currentIndex;
  final ValueChanged<int> onTap;

  const CustomBottomNav({
    super.key,
    required this.currentIndex,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);

    return ClipRect(
      child: BackdropFilter(
        filter: ImageFilter.blur(sigmaX: 20, sigmaY: 20),
        child: Container(
          decoration: BoxDecoration(
            color: theme.scaffoldBackgroundColor.withValues(alpha: 0.92),
            border: Border(top: BorderSide(color: palette.divider)),
          ),
          child: SafeArea(
            top: false,
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 6),
              child: Row(
                children: [
                  _TabItem(
                    icon: Icons.chat_bubble_outline_rounded,
                    activeIcon: Icons.chat_bubble_rounded,
                    label: context.l10n.navChat,
                    active: currentIndex == 0,
                    onTap: () => onTap(0),
                  ),
                  _TabItem(
                    icon: Icons.people_outline_rounded,
                    activeIcon: Icons.people_rounded,
                    label: context.l10n.navFriends,
                    active: currentIndex == 1,
                    onTap: () => onTap(1),
                  ),
                  _TabItem(
                    icon: Icons.person_outline_rounded,
                    activeIcon: Icons.person_rounded,
                    label: context.l10n.navProfile,
                    active: currentIndex == 2,
                    onTap: () => onTap(2),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class _TabItem extends StatelessWidget {
  final IconData icon;
  final IconData activeIcon;
  final String label;
  final bool active;
  final VoidCallback onTap;

  const _TabItem({
    required this.icon,
    required this.activeIcon,
    required this.label,
    required this.active,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);
    final color = active ? theme.colorScheme.tertiary : palette.textSecondary;

    return Expanded(
      child: GestureDetector(
        onTap: onTap,
        behavior: HitTestBehavior.opaque,
        child: AnimatedContainer(
          duration: const Duration(milliseconds: 150),
          padding: const EdgeInsets.symmetric(vertical: 8),
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(12),
            color: active
                ? theme.colorScheme.tertiary.withValues(alpha: 0.06)
                : Colors.transparent,
          ),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(active ? activeIcon : icon, size: 24, color: color),
              const SizedBox(height: 4),
              Text(
                label,
                style: TextStyle(
                  fontFamily: 'Manrope',
                  fontSize: 11,
                  fontWeight: FontWeight.w600,
                  color: color,
                  letterSpacing: 0.02,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
