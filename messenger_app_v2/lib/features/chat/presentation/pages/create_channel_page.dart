import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import '../../../../app/config/theme/app_theme.dart';
import '../../../../app/navigation/app_routes.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/widgets/suda_button.dart';
import '../../../../shared/presentation/widgets/suda_text_field.dart';
import '../../domain/repositories/chat_repository.dart';

/// Single-step channel creation form (simpler than group: no member picker).
class CreateChannelPage extends StatefulWidget {
  const CreateChannelPage({super.key});

  @override
  State<CreateChannelPage> createState() => _CreateChannelPageState();
}

class _CreateChannelPageState extends State<CreateChannelPage> {
  final _nameCtrl = TextEditingController();
  final _handleCtrl = TextEditingController();
  final _descCtrl = TextEditingController();
  bool _isPublic = true;
  bool _isCreating = false;

  @override
  void dispose() {
    _nameCtrl.dispose();
    _handleCtrl.dispose();
    _descCtrl.dispose();
    super.dispose();
  }

  Future<void> _createChannel(BuildContext context) async {
    final name = _nameCtrl.text.trim();
    if (name.isEmpty) return;

    final handle = _handleCtrl.text.trim().replaceFirst('@', '');
    final desc = _descCtrl.text.trim();

    setState(() => _isCreating = true);
    final result = await sl<ChatRepository>().createChat(
      type: 'CHANNEL',
      name: name,
      username: handle.isEmpty ? null : handle,
      description: desc.isEmpty ? null : desc,
      visibility: _isPublic ? 'PUBLIC' : 'PRIVATE',
    );
    if (!context.mounted) return;
    setState(() => _isCreating = false);

    // ApiClient already surfaces failures via AppFeedback — only handle success here.
    result.fold(
      (_) {},
      (chat) => context.pushReplacementNamed(
        AppRoutes.chatDetail,
        pathParameters: {'id': chat.id},
        extra: {'name': name, 'interlocutorId': '', 'chatType': 'CHANNEL'},
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final l10n = context.l10n;

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      appBar: AppBar(
        backgroundColor: theme.scaffoldBackgroundColor,
        elevation: 0,
        leading: IconButton(
          icon: Icon(Icons.arrow_back_rounded, color: theme.colorScheme.onSurface),
          onPressed: () => context.pop(),
        ),
        title: Text(
          l10n.channelNewTitle,
          style: TextStyle(
            fontFamily: 'Manrope',
            fontWeight: FontWeight.w700,
            color: theme.colorScheme.onSurface,
          ),
        ),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            SudaTextField(
              label: l10n.channelName,
              controller: _nameCtrl,
              textInputAction: TextInputAction.next,
            ),
            const SizedBox(height: 20),
            SudaTextField(
              label: l10n.channelHandle,
              hint: l10n.channelHandleHint,
              controller: _handleCtrl,
              textInputAction: TextInputAction.next,
            ),
            const SizedBox(height: 20),
            SudaTextField(
              label: l10n.channelDescription,
              hint: l10n.channelDescriptionHint,
              controller: _descCtrl,
              maxLines: 3,
            ),
            const SizedBox(height: 24),
            _VisibilitySelector(
              isPublic: _isPublic,
              onChanged: (v) => setState(() => _isPublic = v),
            ),
            const SizedBox(height: 28),
            ListenableBuilder(
              listenable: _nameCtrl,
              builder: (_, __) => SudaButton(
                label: l10n.channelCreate,
                variant: SudaButtonVariant.primary,
                size: SudaButtonSize.lg,
                fullWidth: true,
                loading: _isCreating,
                onPressed: (_nameCtrl.text.trim().isNotEmpty && !_isCreating)
                    ? () => _createChannel(context)
                    : null,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

/// PUBLIC / PRIVATE radio-style selector built from two tappable cards.
class _VisibilitySelector extends StatelessWidget {
  final bool isPublic;
  final ValueChanged<bool> onChanged;

  const _VisibilitySelector({required this.isPublic, required this.onChanged});

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          l10n.channelVisibilityLabel.toUpperCase(),
          style: TextStyle(
            fontFamily: 'Manrope',
            fontSize: 11,
            fontWeight: FontWeight.w700,
            color: palette.textAccent,
            letterSpacing: 0.12,
          ),
        ),
        const SizedBox(height: 8),
        Row(
          children: [
            Expanded(
              child: _VisibilityOption(
                icon: Icons.public_rounded,
                label: l10n.channelPublic,
                selected: isPublic,
                onTap: () => onChanged(true),
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: _VisibilityOption(
                icon: Icons.lock_outline_rounded,
                label: l10n.channelPrivate,
                selected: !isPublic,
                onTap: () => onChanged(false),
              ),
            ),
          ],
        ),
      ],
    );
  }
}

class _VisibilityOption extends StatelessWidget {
  final IconData icon;
  final String label;
  final bool selected;
  final VoidCallback onTap;

  const _VisibilityOption({
    required this.icon,
    required this.label,
    required this.selected,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final accent = theme.colorScheme.tertiary;

    return Material(
      color: selected
          ? accent.withValues(alpha: 0.12)
          : theme.colorScheme.surface,
      borderRadius: BorderRadius.circular(14),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(14),
        child: Container(
          padding: const EdgeInsets.symmetric(vertical: 14, horizontal: 12),
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(14),
            border: Border.all(
              color: selected ? accent : palette.divider,
              width: selected ? 1.5 : 1,
            ),
          ),
          child: Row(
            children: [
              Icon(
                icon,
                size: 18,
                color: selected ? accent : palette.textSecondary,
              ),
              const SizedBox(width: 8),
              Text(
                label,
                style: TextStyle(
                  fontFamily: 'Manrope',
                  fontSize: 14,
                  fontWeight: FontWeight.w600,
                  color: selected ? accent : theme.colorScheme.onSurface,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
