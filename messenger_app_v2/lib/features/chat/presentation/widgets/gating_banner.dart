import 'package:flutter/material.dart';

import '../../../../app/config/theme/app_theme.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/widgets/suda_button.dart';

/// Paywall banner shown to non-subscribers on paid channels.
/// Iteration-2: token-gating = one-time subscription price in SUDA.
class GatingBanner extends StatelessWidget {
  final double subscriptionPrice;
  final VoidCallback onSubscribeTap;

  /// Legacy compat — ignored, kept so existing call sites compile without change.
  // ignore: avoid_unused_constructor_parameters
  final double? userBalance;

  const GatingBanner({
    super.key,
    required this.onSubscribeTap,
    // New API
    double? subscriptionPrice,
    this.userBalance,
    // Legacy alias
    double? minBalance,
  }) : subscriptionPrice = subscriptionPrice ?? minBalance ?? 0.0;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    final priceStr = subscriptionPrice == subscriptionPrice.truncateToDouble()
        ? subscriptionPrice.toInt().toString()
        : subscriptionPrice.toStringAsFixed(2);

    return Container(
      padding: const EdgeInsets.all(14),
      decoration: BoxDecoration(
        color: theme.colorScheme.tertiary.withValues(alpha: 0.08),
        border: Border.all(color: theme.colorScheme.tertiary.withValues(alpha: 0.4)),
        borderRadius: BorderRadius.circular(14),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Icon(Icons.lock_rounded, size: 16, color: theme.colorScheme.tertiary),
              const SizedBox(width: 8),
              Text(
                l10n.channelTokenGated,
                style: TextStyle(
                  fontFamily: 'Manrope',
                  fontSize: 14,
                  fontWeight: FontWeight.w700,
                  color: theme.colorScheme.tertiary,
                ),
              ),
            ],
          ),
          const SizedBox(height: 8),
          Text(
            l10n.channelTokenGatedDesc(priceStr),
            style: TextStyle(fontFamily: 'Manrope', fontSize: 13, color: palette.textSecondary),
          ),
          const SizedBox(height: 14),
          SudaButton(
            label: l10n.channelSubscribeForPrice(priceStr),
            variant: SudaButtonVariant.primary,
            size: SudaButtonSize.md,
            fullWidth: true,
            onPressed: onSubscribeTap,
          ),
        ],
      ),
    );
  }
}
