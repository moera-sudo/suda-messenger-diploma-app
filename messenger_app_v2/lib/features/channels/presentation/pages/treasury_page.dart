import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import '../../../../app/config/theme/app_theme.dart';
import '../../../../shared/presentation/feedback/app_feedback.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/widgets/suda_button.dart';
import '../../data/models/treasury_models.dart';
import '../bloc/treasury_cubit.dart';

class TreasuryPage extends StatelessWidget {
  final String channelId;
  final bool isOwner;

  const TreasuryPage({super.key, required this.channelId, this.isOwner = false});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) => sl<TreasuryCubit>()..load(channelId),
      child: _TreasuryView(channelId: channelId, isOwner: isOwner),
    );
  }
}

class _TreasuryView extends StatelessWidget {
  final String channelId;
  final bool isOwner;

  const _TreasuryView({required this.channelId, required this.isOwner});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    return BlocListener<TreasuryCubit, TreasuryState>(
      listenWhen: (p, c) => p.status != c.status,
      listener: (context, state) {
        if (state.status == TreasuryStatus.withdrawSuccess) {
          AppFeedback.showSuccess(l10n.treasuryWithdrawSuccess);
          context.read<TreasuryCubit>().load(channelId);
        } else if (state.status == TreasuryStatus.success && state.error != null) {
          AppFeedback.showError(l10n.treasuryWithdrawFailed);
        }
      },
      child: Scaffold(
        backgroundColor: theme.scaffoldBackgroundColor,
        appBar: AppBar(
          backgroundColor: theme.scaffoldBackgroundColor,
          elevation: 0,
          leading: IconButton(
            icon: Icon(Icons.arrow_back_rounded, color: theme.colorScheme.onSurface),
            onPressed: () => context.pop(),
          ),
          title: Text(
            l10n.channelTreasury,
            style: TextStyle(
              fontFamily: 'Manrope',
              fontWeight: FontWeight.w700,
              color: theme.colorScheme.onSurface,
            ),
          ),
        ),
        body: BlocBuilder<TreasuryCubit, TreasuryState>(
          builder: (context, state) {
            if (state.status == TreasuryStatus.loading ||
                state.status == TreasuryStatus.initial) {
              return const Center(child: CircularProgressIndicator());
            }

            if (state.status == TreasuryStatus.failure) {
              return Center(
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Icon(Icons.error_outline_rounded, size: 48, color: palette.danger),
                    const SizedBox(height: 12),
                    Text(
                      state.error ?? '',
                      style: TextStyle(fontFamily: 'Manrope', color: palette.textSecondary),
                      textAlign: TextAlign.center,
                    ),
                    const SizedBox(height: 16),
                    SudaButton(
                      label: l10n.buttonRetry,
                      onPressed: () => context.read<TreasuryCubit>().load(channelId),
                    ),
                  ],
                ),
              );
            }

            final stats = state.stats!;
            return ListView(
              padding: const EdgeInsets.all(16),
              children: [
                // Balance card
                _BalanceCard(stats: stats, isOwner: isOwner, channelId: channelId),
                const SizedBox(height: 16),

                // Top donors
                if (stats.topDonors.isNotEmpty) ...[
                  _SectionHeader(label: l10n.treasuryTopDonors),
                  const SizedBox(height: 8),
                  ...stats.topDonors.map((d) => _DonorTile(donor: d)),
                  const SizedBox(height: 16),
                ],

                // Recent donations
                _SectionHeader(label: l10n.treasuryRecentDonations),
                const SizedBox(height: 8),
                if (state.donations.isEmpty)
                  Padding(
                    padding: const EdgeInsets.symmetric(vertical: 24),
                    child: Center(
                      child: Text(
                        l10n.treasuryEmpty,
                        style: TextStyle(fontFamily: 'Manrope', color: palette.textSecondary),
                      ),
                    ),
                  )
                else
                  ...state.donations.map((d) => _DonationTile(item: d)),
              ],
            );
          },
        ),
      ),
    );
  }
}

class _BalanceCard extends StatelessWidget {
  final TreasuryStats stats;
  final bool isOwner;
  final String channelId;

  const _BalanceCard({required this.stats, required this.isOwner, required this.channelId});

  String _weiToSuda(String wei) {
    try {
      final big = BigInt.parse(wei);
      final divisor = BigInt.from(10).pow(18);
      final intPart = big ~/ divisor;
      final frac = big.remainder(divisor).abs();
      if (frac == BigInt.zero) return intPart.toString();
      final fracStr = frac.toString().padLeft(18, '0').replaceAll(RegExp(r'0+$'), '');
      return '$intPart.$fracStr';
    } catch (_) {
      return wei;
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    final balance = _weiToSuda(stats.sudaBalanceWei);
    final totalDonations = _weiToSuda(stats.totalDonationsWei);

    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: theme.colorScheme.tertiary.withValues(alpha: 0.3)),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Icon(Icons.account_balance_rounded, size: 20, color: theme.colorScheme.tertiary),
              const SizedBox(width: 8),
              Text(
                l10n.treasuryBalance,
                style: TextStyle(
                  fontFamily: 'Manrope',
                  fontSize: 14,
                  color: palette.textSecondary,
                ),
              ),
            ],
          ),
          const SizedBox(height: 8),
          Text(
            '$balance SUDA',
            style: TextStyle(
              fontFamily: 'SpaceGrotesk',
              fontSize: 28,
              fontWeight: FontWeight.w800,
              color: theme.colorScheme.tertiary,
            ),
          ),
          const SizedBox(height: 12),
          Row(
            children: [
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      l10n.treasuryTotalDonations,
                      style: TextStyle(fontFamily: 'Manrope', fontSize: 12, color: palette.textSecondary),
                    ),
                    const SizedBox(height: 2),
                    Text(
                      '$totalDonations SUDA (${stats.totalDonationsCount})',
                      style: TextStyle(
                        fontFamily: 'Manrope',
                        fontSize: 13,
                        fontWeight: FontWeight.w600,
                        color: theme.colorScheme.onSurface,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
          if (isOwner) ...[
            const SizedBox(height: 16),
            SudaButton(
              label: l10n.treasuryWithdraw,
              variant: SudaButtonVariant.primary,
              fullWidth: true,
              onPressed: () => _showWithdrawSheet(context, balance),
            ),
          ],
        ],
      ),
    );
  }

  void _showWithdrawSheet(BuildContext context, String currentBalance) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;
    final ctrl = TextEditingController();
    final cubit = context.read<TreasuryCubit>();

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: theme.colorScheme.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      ),
      builder: (ctx) => Padding(
        padding: EdgeInsets.fromLTRB(
          20, 20, 20, MediaQuery.of(ctx).viewInsets.bottom + 24,
        ),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              l10n.treasuryWithdraw,
              style: TextStyle(
                fontFamily: 'Manrope',
                fontSize: 16,
                fontWeight: FontWeight.w700,
                color: theme.colorScheme.onSurface,
              ),
            ),
            const SizedBox(height: 6),
            Text(
              '$currentBalance SUDA ${l10n.treasuryBalance.toLowerCase()}',
              style: TextStyle(fontFamily: 'Manrope', fontSize: 13, color: palette.textSecondary),
            ),
            const SizedBox(height: 16),
            TextField(
              controller: ctrl,
              autofocus: true,
              keyboardType: const TextInputType.numberWithOptions(decimal: true),
              style: TextStyle(fontFamily: 'Manrope', color: theme.colorScheme.onSurface),
              decoration: InputDecoration(
                hintText: l10n.treasuryWithdrawHint,
                hintStyle: TextStyle(color: palette.textSecondary),
                border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
              ),
            ),
            const SizedBox(height: 16),
            SudaButton(
              label: l10n.treasuryWithdraw,
              variant: SudaButtonVariant.primary,
              fullWidth: true,
              onPressed: () async {
                final text = ctrl.text.trim();
                if (text.isEmpty) return;
                final wei = _sudaToWei(text);
                Navigator.pop(ctx);
                await cubit.withdraw(channelId, wei.toString());
              },
            ),
          ],
        ),
      ),
    );
  }

  BigInt _sudaToWei(String suda) {
    const decimals = 18;
    final parts = suda.split('.');
    final intPart = BigInt.tryParse(parts[0]) ?? BigInt.zero;
    var wei = intPart * BigInt.from(10).pow(decimals);
    if (parts.length == 2 && parts[1].isNotEmpty) {
      var frac = parts[1];
      frac = frac.length > decimals ? frac.substring(0, decimals) : frac.padRight(decimals, '0');
      wei += BigInt.tryParse(frac) ?? BigInt.zero;
    }
    return wei;
  }
}

String _withAt(String u) => u.startsWith('@') ? u : '@$u';

class _SectionHeader extends StatelessWidget {
  final String label;
  const _SectionHeader({required this.label});

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    return Text(
      label,
      style: TextStyle(
        fontFamily: 'Manrope',
        fontSize: 13,
        fontWeight: FontWeight.w700,
        color: palette.textSecondary,
        letterSpacing: 0.5,
      ),
    );
  }
}

class _DonorTile extends StatelessWidget {
  final TopDonor donor;
  const _DonorTile({required this.donor});

  String _weiToSuda(String wei) {
    try {
      final big = BigInt.parse(wei);
      final divisor = BigInt.from(10).pow(18);
      final intPart = big ~/ divisor;
      return intPart.toString();
    } catch (_) {
      return wei;
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    return ListTile(
      contentPadding: const EdgeInsets.symmetric(horizontal: 0, vertical: 2),
      leading: CircleAvatar(
        backgroundColor: theme.colorScheme.tertiary.withValues(alpha: 0.15),
        child: Text(
          donor.displayName.isNotEmpty ? donor.displayName[0].toUpperCase() : '?',
          style: TextStyle(
            fontFamily: 'Manrope',
            fontWeight: FontWeight.w700,
            color: theme.colorScheme.tertiary,
          ),
        ),
      ),
      title: Text(
        donor.displayName.isNotEmpty ? donor.displayName : _withAt(donor.username),
        style: TextStyle(
          fontFamily: 'Manrope',
          fontWeight: FontWeight.w600,
          color: theme.colorScheme.onSurface,
        ),
      ),
      subtitle: Text(
        _withAt(donor.username),
        style: TextStyle(fontFamily: 'Manrope', fontSize: 12, color: palette.textSecondary),
      ),
      trailing: Text(
        '${_weiToSuda(donor.amountWei)} SUDA',
        style: TextStyle(
          fontFamily: 'Manrope',
          fontWeight: FontWeight.w700,
          color: theme.colorScheme.tertiary,
        ),
      ),
    );
  }
}

class _DonationTile extends StatelessWidget {
  final DonationItem item;
  const _DonationTile({required this.item});

  String _weiToSuda(String wei) {
    try {
      final big = BigInt.parse(wei);
      final divisor = BigInt.from(10).pow(18);
      final intPart = big ~/ divisor;
      return intPart.toString();
    } catch (_) {
      return wei;
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;

    return Container(
      margin: const EdgeInsets.only(bottom: 8),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  _withAt(item.fromUsername),
                  style: TextStyle(
                    fontFamily: 'Manrope',
                    fontWeight: FontWeight.w600,
                    color: theme.colorScheme.onSurface,
                  ),
                ),
                if (item.message != null && item.message!.isNotEmpty) ...[
                  const SizedBox(height: 2),
                  Text(
                    item.message!,
                    style: TextStyle(fontFamily: 'Manrope', fontSize: 13, color: palette.textSecondary),
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                  ),
                ],
              ],
            ),
          ),
          Text(
            '${_weiToSuda(item.amountWei)} SUDA',
            style: TextStyle(
              fontFamily: 'Manrope',
              fontWeight: FontWeight.w700,
              color: theme.colorScheme.tertiary,
            ),
          ),
        ],
      ),
    );
  }
}
