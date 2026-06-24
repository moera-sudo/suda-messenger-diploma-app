import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../../../../app/DI/get_it.dart';
import '../../../../app/config/theme/app_theme.dart';
import '../../../../features/wallet/data/models/wallet_models.dart';
import '../../../../features/wallet/domain/repositories/wallet_repository.dart';
import '../../../../shared/presentation/feedback/app_feedback.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/widgets/suda_button.dart';
import '../../data/models/chat_models.dart';

/// Bottom sheet for in-chat SUDA transfer (§4 of TransactionService-Implementation.md).
///
/// DIRECT chat: recipient username is pre-filled (loaded from the interlocutor profile).
/// GROUP chat:  user picks from a dropdown of chat members.
class TransferSheet extends StatefulWidget {
  final String chatId;
  final String chatType; // DIRECT | GROUP
  final String? interlocutorUsername; // pre-filled for DIRECT
  final Map<String, ChatMember> members;
  final String currentUserId;

  const TransferSheet({
    super.key,
    required this.chatId,
    required this.chatType,
    this.interlocutorUsername,
    required this.members,
    required this.currentUserId,
  });

  static Future<void> show(
    BuildContext context, {
    required String chatId,
    required String chatType,
    String? interlocutorUsername,
    required Map<String, ChatMember> members,
    required String currentUserId,
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
          child: TransferSheet(
            chatId: chatId,
            chatType: chatType,
            interlocutorUsername: interlocutorUsername,
            members: members,
            currentUserId: currentUserId,
          ),
        ),
      );

  @override
  State<TransferSheet> createState() => _TransferSheetState();
}

class _TransferSheetState extends State<TransferSheet> {
  final _amountCtrl = TextEditingController();
  final _noteCtrl   = TextEditingController();

  // GROUP: selected member id → derive username
  String? _selectedMemberId;
  bool _sending = false;

  bool get _isDirect => widget.chatType == 'DIRECT';

  String? get _recipientUsername {
    if (_isDirect) return widget.interlocutorUsername;
    if (_selectedMemberId == null) return null;
    return widget.members[_selectedMemberId]?.username;
  }

  @override
  void dispose() {
    _amountCtrl.dispose();
    _noteCtrl.dispose();
    super.dispose();
  }

  Future<void> _send() async {
    final l10n = context.l10n;
    final username = _recipientUsername;
    if (username == null || username.isEmpty) {
      AppFeedback.showError(l10n.transferSelectRecipient);
      return;
    }
    final amountWei = sudaToWei(_amountCtrl.text.trim());
    if (amountWei <= BigInt.zero) {
      AppFeedback.showError(l10n.transferAmountHint);
      return;
    }

    setState(() => _sending = true);
    final result = await sl<WalletRepository>().transferInChat(
      toUsername: username,
      amountWei: amountWei,
      chatId: widget.chatId,
      note: _noteCtrl.text.trim().isEmpty ? null : _noteCtrl.text.trim(),
    );
    if (!mounted) return;
    result.fold(
      (f) {
        setState(() => _sending = false);
        AppFeedback.showError(f.message);
      },
      (_) {
        Navigator.pop(context);
        AppFeedback.showSuccess(l10n.transferPending);
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme  = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n   = context.l10n;

    // GROUP members list (excluding self)
    final memberEntries = widget.members.entries
        .where((e) => e.key != widget.currentUserId)
        .toList();

    return SafeArea(
      top: false,
      child: Padding(
        padding: const EdgeInsets.fromLTRB(20, 0, 20, 20),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Handle
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
              l10n.attachTransfer,
              style: TextStyle(
                fontFamily: 'Manrope',
                fontSize: 18,
                fontWeight: FontWeight.w700,
                color: theme.colorScheme.onSurface,
              ),
            ),
            const SizedBox(height: 20),

            // Recipient
            if (_isDirect && widget.interlocutorUsername != null)
              _RecipientLabel(
                label: l10n.transferSendingTo(widget.interlocutorUsername!),
                palette: palette,
                theme: theme,
              )
            else if (!_isDirect) ...[
              Text(
                l10n.transferSelectRecipient,
                style: TextStyle(fontFamily: 'Manrope', fontSize: 13, color: palette.textSecondary),
              ),
              const SizedBox(height: 8),
              DropdownButtonFormField<String>(
                initialValue: _selectedMemberId,
                dropdownColor: theme.colorScheme.surface,
                decoration: InputDecoration(
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
                hint: Text(l10n.transferSelectRecipient,
                    style: TextStyle(fontFamily: 'Manrope', color: palette.textSecondary)),
                items: memberEntries.map((e) {
                  final m = e.value;
                  return DropdownMenuItem(
                    value: e.key,
                    child: Text(
                      m.displayName.isNotEmpty ? m.displayName : m.username,
                      style: TextStyle(fontFamily: 'Manrope', color: theme.colorScheme.onSurface),
                    ),
                  );
                }).toList(),
                onChanged: (v) => setState(() => _selectedMemberId = v),
              ),
            ],

            const SizedBox(height: 16),

            // Amount
            TextField(
              controller: _amountCtrl,
              keyboardType: const TextInputType.numberWithOptions(decimal: true),
              inputFormatters: [
                FilteringTextInputFormatter.allow(RegExp(r'[0-9.]')),
              ],
              style: TextStyle(fontFamily: 'Manrope', color: theme.colorScheme.onSurface),
              decoration: InputDecoration(
                hintText: l10n.transferAmountHint,
                hintStyle: TextStyle(color: palette.textSecondary),
                suffixText: 'SUDA',
                suffixStyle: TextStyle(
                  fontFamily: 'Manrope',
                  fontWeight: FontWeight.w600,
                  color: theme.colorScheme.tertiary,
                ),
                filled: true,
                fillColor: theme.colorScheme.surface,
                border: OutlineInputBorder(borderRadius: BorderRadius.circular(12),
                    borderSide: BorderSide(color: palette.divider)),
                enabledBorder: OutlineInputBorder(borderRadius: BorderRadius.circular(12),
                    borderSide: BorderSide(color: palette.divider)),
                contentPadding: const EdgeInsets.symmetric(horizontal: 14, vertical: 12),
              ),
            ),
            const SizedBox(height: 12),

            // Note
            TextField(
              controller: _noteCtrl,
              style: TextStyle(fontFamily: 'Manrope', color: theme.colorScheme.onSurface),
              decoration: InputDecoration(
                hintText: l10n.transferNoteHint,
                hintStyle: TextStyle(color: palette.textSecondary),
                filled: true,
                fillColor: theme.colorScheme.surface,
                border: OutlineInputBorder(borderRadius: BorderRadius.circular(12),
                    borderSide: BorderSide(color: palette.divider)),
                enabledBorder: OutlineInputBorder(borderRadius: BorderRadius.circular(12),
                    borderSide: BorderSide(color: palette.divider)),
                contentPadding: const EdgeInsets.symmetric(horizontal: 14, vertical: 12),
              ),
            ),
            const SizedBox(height: 24),

            SudaButton(
              label: l10n.attachTransfer,
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

class _RecipientLabel extends StatelessWidget {
  final String label;
  final AppPaletteTheme palette;
  final ThemeData theme;

  const _RecipientLabel({
    required this.label,
    required this.palette,
    required this.theme,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 12),
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        border: Border.all(color: palette.divider),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Text(
        label,
        style: TextStyle(
          fontFamily: 'Manrope',
          fontSize: 14,
          color: theme.colorScheme.onSurface,
        ),
      ),
    );
  }
}
