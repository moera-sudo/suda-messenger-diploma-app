import 'package:flutter/material.dart';
import 'package:file_picker/file_picker.dart';
import 'package:wechat_assets_picker/wechat_assets_picker.dart';

import '../../../../app/config/theme/app_theme.dart';
import '../../../../shared/presentation/l10n_extensions.dart';

typedef AttachResult = ({String path, String type});

/// Bottom sheet for selecting media attachments.
/// Returns [AttachResult] with file path and media type (IMAGE, FILE, VOICE).
class AttachSheet extends StatelessWidget {
  final String chatType;

  const AttachSheet({super.key, this.chatType = 'DIRECT'});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    final tiles = [
      _AttachTile(
        icon: Icons.image_rounded,
        label: l10n.attachPhoto,
        color: const Color(0xFFFF8FB3),
        onTap: () => _pickMedia(context),
      ),
      _AttachTile(
        icon: Icons.attach_file_rounded,
        label: l10n.attachFile,
        color: theme.colorScheme.tertiary,
        onTap: () => _pickFile(context),
      ),
      _AttachTile(
        icon: Icons.mic_rounded,
        label: l10n.attachVoice,
        color: palette.success,
        onTap: () => Navigator.of(context).pop(null), // voice handled in input
      ),
      if (chatType != 'SAVED' && chatType != 'CHANNEL')
        _AttachTile(
          icon: Icons.compare_arrows_rounded,
          label: l10n.attachTransfer,
          color: const Color(0xFFFFB347),
          onTap: () => Navigator.of(context).pop<AttachResult>(
            (path: '', type: 'TRANSFER'),
          ),
        ),
    ];

    return SafeArea(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            margin: const EdgeInsets.only(top: 8, bottom: 16),
            width: 38,
            height: 4,
            decoration: BoxDecoration(
              color: palette.textSecondary.withValues(alpha: 0.4),
              borderRadius: BorderRadius.circular(2),
            ),
          ),
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 0, 16, 20),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: tiles,
            ),
          ),
        ],
      ),
    );
  }

  Future<void> _pickMedia(BuildContext context) async {
    final theme = Theme.of(context);
    // In-app gallery grid (photos + videos) rendered inside the messenger,
    // themed to match the app — no external gallery/file-explorer app.
    final assets = await AssetPicker.pickAssets(
      context,
      pickerConfig: AssetPickerConfig(
        maxAssets: 1,
        requestType: RequestType.common,
        pickerTheme: AssetPicker.themeData(
          theme.colorScheme.tertiary,
          light: theme.brightness == Brightness.light,
        ),
        textDelegate: assetPickerTextDelegateFromLocale(
          Localizations.localeOf(context),
        ),
      ),
    );
    if (assets == null || assets.isEmpty || !context.mounted) return;

    final asset = assets.first;
    final file = await asset.file;
    if (file == null || !context.mounted) return;

    final type = asset.type == AssetType.video ? 'VIDEO' : 'IMAGE';
    Navigator.of(context).pop<AttachResult>((path: file.path, type: type));
  }

  Future<void> _pickFile(BuildContext context) async {
    final result = await FilePicker.platform.pickFiles();
    if (result == null || result.files.isEmpty || !context.mounted) return;
    final file = result.files.single;
    if (file.path == null) return;

    Navigator.of(context).pop<AttachResult>(
      (path: file.path!, type: 'FILE'),
    );
  }
}

class _AttachTile extends StatelessWidget {
  final IconData icon;
  final String label;
  final Color color;
  final VoidCallback onTap;

  const _AttachTile({
    required this.icon,
    required this.label,
    required this.color,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;

    return GestureDetector(
      onTap: onTap,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            width: 56,
            height: 56,
            decoration: BoxDecoration(
              color: color.withValues(alpha: 0.18),
              shape: BoxShape.circle,
            ),
            child: Icon(icon, size: 26, color: color),
          ),
          const SizedBox(height: 8),
          Text(
            label,
            style: TextStyle(
              fontFamily: 'Manrope',
              fontSize: 12,
              fontWeight: FontWeight.w600,
              color: palette.textSecondary,
            ),
          ),
        ],
      ),
    );
  }
}
