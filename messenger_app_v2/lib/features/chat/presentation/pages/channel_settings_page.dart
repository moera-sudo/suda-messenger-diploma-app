import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:go_router/go_router.dart';
import 'package:image_picker/image_picker.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import '../../../../app/config/theme/app_theme.dart';
import '../../../../features/channels/data/models/channel_settings_model.dart';
import '../../../../features/channels/domain/repositories/channel_repository.dart';
import '../../../../features/media/domain/repositories/media_repository.dart';
import '../../../../features/wallet/data/models/wallet_models.dart';
import '../../../../shared/presentation/feedback/app_feedback.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/widgets/suda_avatar.dart';
import '../../../../shared/presentation/widgets/suda_button.dart';
import '../../../../shared/presentation/widgets/suda_text_field.dart';
import '../../domain/repositories/chat_repository.dart';

/// Channel settings (OWNER/ADMIN). Saves to two endpoints:
/// profile fields → PUT /chats/{id}, channel fields → PUT /channels/{id}/settings.
class ChannelSettingsPage extends StatefulWidget {
  final String chatId;
  const ChannelSettingsPage({super.key, required this.chatId});

  @override
  State<ChannelSettingsPage> createState() => _ChannelSettingsPageState();
}

class _ChannelSettingsPageState extends State<ChannelSettingsPage> {
  final _nameCtrl = TextEditingController();
  final _handleCtrl = TextEditingController();
  final _descCtrl = TextEditingController();
  final _gatingAmountCtrl = TextEditingController();
  final _picker = ImagePicker();

  ChannelSettings? _original;
  bool _isPublic = true;
  bool _commentsEnabled = false;
  String? _avatarMediaId;
  String? _avatarUrl;
  bool _loading = true;
  bool _saving = false;

  // Token gating (public channels only).
  bool _gatingEnabled = false;
  bool _originalGatingEnabled = false;
  double _originalGatingAmount = 0;

  @override
  void initState() {
    super.initState();
    _load();
  }

  @override
  void dispose() {
    _nameCtrl.dispose();
    _handleCtrl.dispose();
    _descCtrl.dispose();
    _gatingAmountCtrl.dispose();
    super.dispose();
  }

  Future<void> _load() async {
    final result = await sl<ChannelRepository>().getChannelSettings(widget.chatId);
    if (!mounted) return;
    await result.fold(
      (f) async {
        AppFeedback.showError(f.message);
        setState(() => _loading = false);
      },
      (s) async {
        _original = s;
        _nameCtrl.text = s.name;
        _handleCtrl.text = s.username ?? '';
        _descCtrl.text = s.description ?? '';
        _isPublic = s.isPublic;
        _commentsEnabled = s.commentsEnabled;
        _avatarMediaId = s.avatarMediaId;
        if (s.avatarMediaId != null) {
          final urlRes = await sl<MediaRepository>().getMediaUrl(s.avatarMediaId!);
          if (mounted) _avatarUrl = urlRes.fold((_) => null, (u) => u);
        }

        // Load the current token-gating rule (public channels only).
        final ruleRes = await sl<ChannelRepository>().getGatingRule(widget.chatId);
        ruleRes.fold((_) {}, (rule) {
          if (rule != null) {
            _gatingEnabled = true;
            _originalGatingEnabled = true;
            _originalGatingAmount = rule.subscriptionPrice;
            _gatingAmountCtrl.text = rule.subscriptionPrice
                .toStringAsFixed(rule.subscriptionPrice.truncateToDouble() == rule.subscriptionPrice ? 0 : 2);
          }
        });

        if (mounted) setState(() => _loading = false);
      },
    );
  }

  Future<void> _pickAvatar() async {
    final file = await _picker.pickImage(
        source: ImageSource.gallery, imageQuality: 88, maxWidth: 1200);
    if (!mounted || file == null) return;
    setState(() => _saving = true);
    final upload = await sl<MediaRepository>().uploadFile(
      filePath: file.path,
      mediaType: 'CHAT_AVATAR',
      filename: file.name,
    );
    if (!mounted) return;
    await upload.fold(
      (f) async => AppFeedback.showError(f.message),
      (mediaId) async {
        final urlRes = await sl<MediaRepository>().getMediaUrl(mediaId);
        if (!mounted) return;
        setState(() {
          _avatarMediaId = mediaId;
          _avatarUrl = urlRes.fold((_) => null, (u) => u);
        });
      },
    );
    if (mounted) setState(() => _saving = false);
  }

  Future<void> _save() async {
    final orig = _original;
    if (orig == null) return;
    final name = _nameCtrl.text.trim();
    final desc = _descCtrl.text.trim();
    final handle = _handleCtrl.text.trim().replaceFirst('@', '');
    final visibility = _isPublic ? 'PUBLIC' : 'PRIVATE';

    setState(() => _saving = true);

    // Profile fields → PUT /chats/{id}
    final nameChanged = name.isNotEmpty && name != orig.name;
    final descChanged = desc != (orig.description ?? '');
    final avatarChanged = _avatarMediaId != orig.avatarMediaId;
    if (nameChanged || descChanged || avatarChanged) {
      final res = await sl<ChatRepository>().updateChat(
        widget.chatId,
        name: nameChanged ? name : null,
        description: descChanged ? desc : null,
        avatarMediaId: avatarChanged ? _avatarMediaId : null,
      );
      if (!mounted) return;
      final failed = res.isLeft();
      if (failed) {
        res.leftMap((f) => AppFeedback.showError(f.message));
        setState(() => _saving = false);
        return;
      }
    }

    // Channel fields → PUT /channels/{id}/settings
    final commentsChanged = _commentsEnabled != orig.commentsEnabled;
    final visibilityChanged = visibility != orig.visibility;
    final usernameChanged = handle.isNotEmpty && handle != (orig.username ?? '');
    if (commentsChanged || visibilityChanged || usernameChanged) {
      final res = await sl<ChannelRepository>().updateChannelSettings(
        widget.chatId,
        commentsEnabled: commentsChanged ? _commentsEnabled : null,
        visibility: visibilityChanged ? visibility : null,
        username: usernameChanged ? handle : null,
      );
      if (!mounted) return;
      if (res.isLeft()) {
        res.leftMap((f) => AppFeedback.showError(f.message));
        setState(() => _saving = false);
        return;
      }
    }

    // Token gating → /tx/gating/rule (only meaningful for PUBLIC channels).
    if (_isPublic) {
      final amount = double.tryParse(_gatingAmountCtrl.text.trim()) ?? 0;
      final amountChanged = amount != _originalGatingAmount;
      if (_gatingEnabled && amount > 0 && (!_originalGatingEnabled || amountChanged)) {
        final res = await sl<ChannelRepository>()
            .setGatingRule(widget.chatId, sudaToWei(_gatingAmountCtrl.text.trim()));
        if (!mounted) return;
        if (res.isLeft()) {
          res.leftMap((f) => AppFeedback.showError(f.message));
          setState(() => _saving = false);
          return;
        }
      } else if (!_gatingEnabled && _originalGatingEnabled) {
        final res = await sl<ChannelRepository>().deleteGatingRule(widget.chatId);
        if (!mounted) return;
        if (res.isLeft()) {
          res.leftMap((f) => AppFeedback.showError(f.message));
          setState(() => _saving = false);
          return;
        }
      }
    }

    if (!mounted) return;
    setState(() => _saving = false);
    AppFeedback.showSuccess(context.l10n.channelSettingsSaved);
    context.pop();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
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
          l10n.channelSettingsTitle,
          style: TextStyle(
              fontFamily: 'Manrope', fontWeight: FontWeight.w700, color: theme.colorScheme.onSurface),
        ),
      ),
      body: _loading
          ? const Center(child: CircularProgressIndicator())
          : SingleChildScrollView(
              padding: const EdgeInsets.all(20),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Avatar
                  Center(
                    child: GestureDetector(
                      onTap: _saving ? null : _pickAvatar,
                      child: Stack(
                        alignment: Alignment.bottomRight,
                        children: [
                          SudaAvatar(
                            mediaUrl: _avatarUrl,
                            initials: _nameCtrl.text.isNotEmpty ? _nameCtrl.text : '?',
                            size: 88,
                            ring: true,
                          ),
                          Container(
                            padding: const EdgeInsets.all(6),
                            decoration: BoxDecoration(
                              color: theme.colorScheme.primary,
                              shape: BoxShape.circle,
                              border: Border.all(color: theme.scaffoldBackgroundColor, width: 2),
                            ),
                            child: const Icon(Icons.camera_alt_rounded, size: 16, color: Colors.white),
                          ),
                        ],
                      ),
                    ),
                  ),
                  const SizedBox(height: 24),
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
                  // Token gating — public channels only.
                  if (_isPublic) ...[
                    const SizedBox(height: 20),
                    _GatingSection(
                      enabled: _gatingEnabled,
                      amountCtrl: _gatingAmountCtrl,
                      onToggle: (v) => setState(() => _gatingEnabled = v),
                    ),
                  ],
                  const SizedBox(height: 20),
                  // Comments toggle
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 6),
                    decoration: BoxDecoration(
                      color: theme.colorScheme.surface,
                      border: Border.all(color: palette.divider),
                      borderRadius: BorderRadius.circular(14),
                    ),
                    child: Row(
                      children: [
                        Icon(Icons.mode_comment_outlined, size: 20, color: palette.textAccent),
                        const SizedBox(width: 12),
                        Expanded(
                          child: Text(
                            l10n.channelCommentsEnabledLabel,
                            style: TextStyle(
                                fontFamily: 'Manrope', fontSize: 14, fontWeight: FontWeight.w500,
                                color: theme.colorScheme.onSurface),
                          ),
                        ),
                        Switch(
                          value: _commentsEnabled,
                          onChanged: (v) => setState(() => _commentsEnabled = v),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: 28),
                  SudaButton(
                    label: l10n.channelSettingsSave,
                    variant: SudaButtonVariant.primary,
                    size: SudaButtonSize.lg,
                    fullWidth: true,
                    loading: _saving,
                    onPressed: _saving ? null : _save,
                  ),
                ],
              ),
            ),
    );
  }
}

/// PUBLIC / PRIVATE selector (mirrors the one on create_channel_page).
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
              fontFamily: 'Manrope', fontSize: 11, fontWeight: FontWeight.w700,
              color: palette.textAccent, letterSpacing: 0.12),
        ),
        const SizedBox(height: 8),
        Row(children: [
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
        ]),
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
      color: selected ? accent.withValues(alpha: 0.12) : theme.colorScheme.surface,
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
          child: Row(children: [
            Icon(icon, size: 18, color: selected ? accent : palette.textSecondary),
            const SizedBox(width: 8),
            Text(
              label,
              style: TextStyle(
                  fontFamily: 'Manrope', fontSize: 14, fontWeight: FontWeight.w600,
                  color: selected ? accent : theme.colorScheme.onSurface),
            ),
          ]),
        ),
      ),
    );
  }
}

/// Token-gating section: toggle + min SUDA balance field (public channels only).
class _GatingSection extends StatelessWidget {
  final bool enabled;
  final TextEditingController amountCtrl;
  final ValueChanged<bool> onToggle;

  const _GatingSection({
    required this.enabled,
    required this.amountCtrl,
    required this.onToggle,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 6),
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        border: Border.all(color: palette.divider),
        borderRadius: BorderRadius.circular(14),
      ),
      child: Column(
        children: [
          Row(
            children: [
              Icon(Icons.lock_outline_rounded, size: 20, color: palette.textAccent),
              const SizedBox(width: 12),
              Expanded(
                child: Text(
                  l10n.gatingSettingsEnable,
                  style: TextStyle(
                      fontFamily: 'Manrope', fontSize: 14, fontWeight: FontWeight.w500,
                      color: theme.colorScheme.onSurface),
                ),
              ),
              Switch(value: enabled, onChanged: onToggle),
            ],
          ),
          if (enabled)
            Padding(
              padding: const EdgeInsets.fromLTRB(0, 4, 0, 10),
              child: TextField(
                controller: amountCtrl,
                keyboardType: const TextInputType.numberWithOptions(decimal: true),
                inputFormatters: [
                  FilteringTextInputFormatter.allow(RegExp(r'[0-9.]')),
                ],
                style: TextStyle(fontFamily: 'Manrope', color: theme.colorScheme.onSurface),
                decoration: InputDecoration(
                  labelText: l10n.gatingSettingsMinBalance,
                  labelStyle: TextStyle(color: palette.textSecondary),
                  suffixText: 'SUDA',
                  suffixStyle: TextStyle(
                    fontFamily: 'Manrope',
                    fontWeight: FontWeight.w600,
                    color: theme.colorScheme.tertiary,
                  ),
                  isDense: true,
                ),
              ),
            ),
        ],
      ),
    );
  }
}
