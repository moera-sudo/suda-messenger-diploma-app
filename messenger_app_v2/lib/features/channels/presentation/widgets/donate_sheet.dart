import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../../../../app/DI/get_it.dart';
import '../../../../app/config/theme/app_theme.dart';
import '../../../../features/wallet/data/models/wallet_models.dart';
import '../../../../features/wallet/domain/repositories/wallet_repository.dart';
import '../../../../shared/presentation/feedback/app_feedback.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/widgets/suda_button.dart';

/// Bottom sheet for donating SUDA to a channel treasury (§5 of
/// TransactionService-Implementation.md). Available ONLY for channels.
class DonateSheet extends StatefulWidget {
  final String channelId;
  final String? chatId; // optional: links donation to a chat system message

  const DonateSheet({super.key, required this.channelId, this.chatId});

  static Future<void> show(
    BuildContext context, {
    required String channelId,
    String? chatId,
  }) =>
      showModalBottomSheet<void>(
        context: context,
        isScrollControlled: true,
        backgroundColor: Theme.of(context).colorScheme.surface,
        shape: const RoundedRectangleBorder(
          borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
        ),
        builder: (_) => Padding(
          padding: MediaQuery.of(context).viewInsets,
          child: DonateSheet(channelId: channelId, chatId: chatId),
        ),
      );

  @override
  State<DonateSheet> createState() => _DonateSheetState();
}

class _DonateSheetState extends State<DonateSheet> {
  final _amountCtrl  = TextEditingController();
  final _messageCtrl = TextEditingController();
  bool _sending = false;

  static const _maxMessageLength = 200;

  @override
  void dispose() {
    _amountCtrl.dispose();
    _messageCtrl.dispose();
    super.dispose();
  }

  Future<void> _send() async {
    final l10n = context.l10n;
    final amountWei = sudaToWei(_amountCtrl.text.trim());
    if (amountWei <= BigInt.zero) {
      AppFeedback.showError(l10n.donateAmountHint);
      return;
    }
    final msg = _messageCtrl.text.trim();

    setState(() => _sending = true);
    final result = await sl<WalletRepository>().donateToChannel(
      toChannelId: widget.channelId,
      amountWei: amountWei,
      message: msg.isEmpty ? null : msg,
      chatId: widget.chatId,
    );
    if (!mounted) return;
    result.fold(
      (f) {
        setState(() => _sending = false);
        AppFeedback.showError(f.message);
      },
      (_) {
        Navigator.pop(context);
        AppFeedback.showSuccess(l10n.donateSent);
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme   = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n    = context.l10n;

    return SafeArea(
      top: false,
      child: Padding(
        padding: const EdgeInsets.fromLTRB(20, 0, 20, 20),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Center(
              child: Container(
                margin: const EdgeInsets.only(top: 12, bottom: 20),
                width: 38,
                height: 4,
                decoration: BoxDecoration(
                  color: palette.textSecondary.withValues(alpha: 0.4),
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
            ),
            Text(
              l10n.userProfileDonate,
              style: TextStyle(
                fontFamily: 'Manrope',
                fontSize: 18,
                fontWeight: FontWeight.w700,
                color: theme.colorScheme.onSurface,
              ),
            ),
            const SizedBox(height: 20),

            // Amount
            TextField(
              controller: _amountCtrl,
              keyboardType: const TextInputType.numberWithOptions(decimal: true),
              inputFormatters: [
                FilteringTextInputFormatter.allow(RegExp(r'[0-9.]')),
              ],
              style: TextStyle(fontFamily: 'Manrope', color: theme.colorScheme.onSurface),
              decoration: InputDecoration(
                hintText: l10n.donateAmountHint,
                hintStyle: TextStyle(color: palette.textSecondary),
                suffixText: 'SUDA',
                suffixStyle: TextStyle(
                  fontFamily: 'Manrope',
                  fontWeight: FontWeight.w600,
                  color: theme.colorScheme.tertiary,
                ),
                filled: true,
                fillColor: theme.colorScheme.surface,
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                  borderSide: BorderSide(color: palette.divider),
                ),
                enabledBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                  borderSide: BorderSide(color: palette.divider),
                ),
                contentPadding: const EdgeInsets.symmetric(horizontal: 14, vertical: 12),
              ),
            ),
            const SizedBox(height: 12),

            // Message with character counter
            ValueListenableBuilder<TextEditingValue>(
              valueListenable: _messageCtrl,
              builder: (_, v, __) {
                final remaining = _maxMessageLength - v.text.length;
                return TextField(
                  controller: _messageCtrl,
                  maxLines: 3,
                  maxLength: _maxMessageLength,
                  style: TextStyle(fontFamily: 'Manrope', color: theme.colorScheme.onSurface),
                  decoration: InputDecoration(
                    hintText: l10n.donateMessageHint,
                    hintStyle: TextStyle(color: palette.textSecondary),
                    counterText: '$remaining',
                    counterStyle: TextStyle(
                      fontFamily: 'Manrope',
                      fontSize: 11,
                      color: remaining < 20 ? palette.danger : palette.textSecondary,
                    ),
                    filled: true,
                    fillColor: theme.colorScheme.surface,
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(12),
                      borderSide: BorderSide(color: palette.divider),
                    ),
                    enabledBorder: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(12),
                      borderSide: BorderSide(color: palette.divider),
                    ),
                    contentPadding: const EdgeInsets.symmetric(horizontal: 14, vertical: 12),
                  ),
                );
              },
            ),
            const SizedBox(height: 20),

            SudaButton(
              label: l10n.userProfileDonate,
              variant: SudaButtonVariant.primary,
              size: SudaButtonSize.lg,
              fullWidth: true,
              loading: _sending,
              onPressed: _sending ? null : _send,
            ),
          ],
        ),
      ),
    );
  }
}
