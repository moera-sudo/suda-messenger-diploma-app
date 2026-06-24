import 'dart:async';
import 'dart:io';

import 'package:audio_waveforms/audio_waveforms.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:path_provider/path_provider.dart';
import 'package:permission_handler/permission_handler.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import '../../../../app/config/theme/app_theme.dart';
import '../../../../shared/data/logger/app_logger.dart';
import '../../../../shared/presentation/feedback/app_feedback.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import 'attach_sheet.dart';

class ChatInputField extends StatefulWidget {
  final ValueChanged<String> onSend;
  final ValueChanged<bool> onTyping;
  final String? initialText;
  final String? replyToAuthor;
  final String? replyToText;
  final VoidCallback? onClearReply;
  final void Function(String filePath, String messageType)? onAttachMedia;
  final void Function(String audioPath, int durationSec)? onVoiceSend;
  final String chatType;

  const ChatInputField({
    super.key,
    required this.onSend,
    required this.onTyping,
    this.initialText,
    this.replyToAuthor,
    this.replyToText,
    this.onClearReply,
    this.onAttachMedia,
    this.onVoiceSend,
    this.chatType = 'DIRECT',
  });

  @override
  State<ChatInputField> createState() => _ChatInputFieldState();
}

class _ChatInputFieldState extends State<ChatInputField> {
  late final TextEditingController _ctrl;
  bool _hasText = false;
  bool _isTypingSent = false;

  // Voice recording state
  bool _isRecording = false;
  int _recordSec = 0;
  Timer? _recordTimer;
  RecorderController? _recorder;

  @override
  void initState() {
    super.initState();
    _ctrl = TextEditingController(text: widget.initialText);
    _hasText = widget.initialText?.isNotEmpty ?? false;
    _ctrl.addListener(_onTextChanged);
  }

  void _onTextChanged() {
    final hasText = _ctrl.text.isNotEmpty;
    if (hasText != _hasText) setState(() => _hasText = hasText);
    if (hasText && !_isTypingSent) {
      widget.onTyping(true);
      _isTypingSent = true;
    } else if (!hasText && _isTypingSent) {
      widget.onTyping(false);
      _isTypingSent = false;
    }
  }

  @override
  void didUpdateWidget(ChatInputField old) {
    super.didUpdateWidget(old);
    if (widget.initialText != old.initialText) {
      final t = widget.initialText ?? '';
      _ctrl.value = TextEditingValue(
        text: t,
        selection: TextSelection.collapsed(offset: t.length),
      );
      setState(() => _hasText = t.isNotEmpty);
    }
  }

  @override
  void dispose() {
    _ctrl.dispose();
    _recordTimer?.cancel();
    _recorder?.dispose();
    super.dispose();
  }

  void _send() {
    final text = _ctrl.text.trim();
    if (text.isEmpty) return;
    widget.onSend(text);
    _ctrl.clear();
    widget.onTyping(false);
    _isTypingSent = false;
    HapticFeedback.lightImpact();
  }

  Future<void> _startRecording() async {
    final l10n = context.l10n;
    final logger = sl<AppLogger>();

    // Microphone permission is mandatory — without it the recorder produces an
    // empty file and nothing can be sent.
    final status = await Permission.microphone.request();
    if (!status.isGranted) {
      logger.error('Mic permission not granted for voice recording: $status');
      AppFeedback.showError(l10n.micPermissionDenied);
      return;
    }

    try {
      final dir = await getTemporaryDirectory();
      // audio_waveforms records AAC inside an MPEG-4 container on Android → .m4a.
      final path = '${dir.path}/voice_${DateTime.now().millisecondsSinceEpoch}.m4a';
      // Explicit bitrate (sampleRate already 44100) for consistent encoding and
      // accurate duration metadata, which keeps the waveform in sync on playback.
      _recorder = RecorderController()..bitRate = 128000;
      await _recorder!.record(path: path);
      logger.debug('Voice recording started: $path');
      _recordTimer = Timer.periodic(const Duration(seconds: 1), (_) {
        setState(() => _recordSec++);
      });
      setState(() {
        _isRecording = true;
        _recordSec = 0;
      });
      HapticFeedback.mediumImpact();
    } catch (e) {
      logger.error('Failed to start voice recording', e);
      _recorder?.dispose();
      _recorder = null;
      AppFeedback.showError(l10n.voiceRecordError);
    }
  }

  Future<void> _stopRecording({bool cancel = false}) async {
    _recordTimer?.cancel();
    _recordTimer = null;
    final durationSec = _recordSec;
    final logger = sl<AppLogger>();

    try {
      final path = await _recorder?.stop();
      _recorder?.dispose();
      _recorder = null;

      if (cancel || path == null || durationSec <= 0) {
        if (path != null) {
          final f = File(path);
          if (await f.exists()) await f.delete();
        }
      } else {
        // Guard against empty captures (e.g. permission silently revoked).
        final file = File(path);
        final length = await file.exists() ? await file.length() : 0;
        if (length > 0) {
          logger.info('Voice recorded: $path (${durationSec}s, ${length}B)');
          if (mounted) widget.onVoiceSend?.call(path, durationSec);
        } else {
          logger.error('Voice file empty/missing: $path (length=$length)');
          if (mounted) AppFeedback.showError(context.l10n.voiceRecordError);
        }
      }
    } catch (e) {
      logger.error('Failed to stop voice recording', e);
      _recorder?.dispose();
      _recorder = null;
    }

    if (mounted) setState(() { _isRecording = false; _recordSec = 0; });
  }

  Future<void> _openAttachSheet() async {
    final result = await showModalBottomSheet<AttachResult>(
      context: context,
      backgroundColor: Theme.of(context).colorScheme.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      ),
      builder: (_) => AttachSheet(chatType: widget.chatType),
    );
    if (result != null && mounted) {
      widget.onAttachMedia?.call(result.path, result.type);
    }
  }

  String _formatSecs(int sec) {
    final m = sec ~/ 60;
    final s = sec % 60;
    return '${m.toString().padLeft(2, '0')}:${s.toString().padLeft(2, '0')}';
  }

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);
    final l10n = context.l10n;

    return Container(
      color: theme.scaffoldBackgroundColor,
      // Keep the composer clear of the Android navigation bar / gesture inset.
      child: SafeArea(
        top: false,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
          Divider(height: 1, color: palette.divider),

          // Reply strip
          if (widget.replyToAuthor != null)
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 8),
              child: Row(
                children: [
                  Container(
                    width: 3,
                    height: 32,
                    decoration: BoxDecoration(
                      color: theme.colorScheme.tertiary,
                      borderRadius: BorderRadius.circular(2),
                    ),
                  ),
                  const SizedBox(width: 10),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          widget.replyToAuthor!,
                          style: TextStyle(
                            fontFamily: 'Manrope',
                            fontSize: 12,
                            fontWeight: FontWeight.w700,
                            color: theme.colorScheme.tertiary,
                          ),
                        ),
                        if (widget.replyToText != null)
                          Text(
                            widget.replyToText!,
                            style: TextStyle(
                              fontFamily: 'Manrope',
                              fontSize: 12,
                              color: palette.textSecondary,
                            ),
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                          ),
                      ],
                    ),
                  ),
                  IconButton(
                    icon: Icon(Icons.close_rounded, size: 18, color: palette.textSecondary),
                    onPressed: widget.onClearReply,
                    padding: EdgeInsets.zero,
                    constraints: const BoxConstraints(minWidth: 32, minHeight: 32),
                  ),
                ],
              ),
            ),

          // Input / Recording row
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 8),
            child: _isRecording
                ? _RecordingBar(
                    seconds: _recordSec,
                    formatSecs: _formatSecs,
                    onCancel: () => _stopRecording(cancel: true),
                    onStop: () => _stopRecording(),
                  )
                : Row(
                    crossAxisAlignment: CrossAxisAlignment.end,
                    children: [
                      _CircleIconBtn(
                        icon: Icons.attach_file_rounded,
                        color: palette.textSecondary,
                        size: 38,
                        onTap: _openAttachSheet,
                      ),
                      const SizedBox(width: 6),
                      Expanded(
                        child: Container(
                          constraints: const BoxConstraints(minHeight: 44),
                          decoration: BoxDecoration(
                            color: theme.colorScheme.surface,
                            borderRadius: BorderRadius.circular(20),
                            border: Border.all(color: palette.divider),
                          ),
                          child: TextField(
                            controller: _ctrl,
                            minLines: 1,
                            maxLines: 5,
                            style: TextStyle(
                              fontFamily: 'Manrope',
                              fontSize: 15,
                              color: theme.colorScheme.onSurface,
                            ),
                            cursorColor: theme.colorScheme.tertiary,
                            decoration: InputDecoration(
                              hintText: l10n.messageInputHint,
                              hintStyle: TextStyle(
                                fontFamily: 'Manrope',
                                color: palette.textSecondary,
                                fontSize: 15,
                              ),
                              border: InputBorder.none,
                              isDense: true,
                              contentPadding: const EdgeInsets.symmetric(
                                horizontal: 14,
                                vertical: 10,
                              ),
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(width: 6),
                      AnimatedSwitcher(
                        duration: const Duration(milliseconds: 150),
                        transitionBuilder: (child, anim) =>
                            ScaleTransition(scale: anim, child: child),
                        child: _hasText
                            ? _AccentCircleBtn(
                                key: const ValueKey('send'),
                                icon: Icons.send_rounded,
                                onTap: _send,
                              )
                            : _AccentCircleBtn(
                                key: const ValueKey('mic'),
                                icon: Icons.mic_rounded,
                                onTap: _startRecording,
                              ),
                      ),
                    ],
                  ),
          ),
        ],
        ),
      ),
    );
  }
}

// ─── Recording bar ────────────────────────────────────────────

class _RecordingBar extends StatefulWidget {
  final int seconds;
  final String Function(int) formatSecs;
  final VoidCallback onCancel;
  final VoidCallback onStop;

  const _RecordingBar({
    required this.seconds,
    required this.formatSecs,
    required this.onCancel,
    required this.onStop,
  });

  @override
  State<_RecordingBar> createState() => _RecordingBarState();
}

class _RecordingBarState extends State<_RecordingBar>
    with SingleTickerProviderStateMixin {
  late final AnimationController _pulse;

  @override
  void initState() {
    super.initState();
    _pulse = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 800),
    )..repeat(reverse: true);
  }

  @override
  void dispose() {
    _pulse.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);
    final l10n = context.l10n;

    return RepaintBoundary(
      child: SizedBox(
        height: 52,
        child: Row(
          children: [
            // Pulsing dot
            FadeTransition(
              opacity: _pulse,
              child: Container(
                width: 10,
                height: 10,
                decoration: BoxDecoration(
                  color: palette.danger,
                  shape: BoxShape.circle,
                ),
              ),
            ),
            const SizedBox(width: 10),
            // Timer
            Text(
              widget.formatSecs(widget.seconds),
              style: TextStyle(
                fontFamily: 'JetBrainsMono',
                fontSize: 17,
                fontWeight: FontWeight.w700,
                color: theme.colorScheme.onSurface,
              ),
            ),
            const SizedBox(width: 10),
            // Hint
            Expanded(
              child: Text(
                l10n.voiceSlideCancelHint,
                style: TextStyle(
                  fontFamily: 'Manrope',
                  fontSize: 13,
                  color: palette.textSecondary,
                ),
              ),
            ),
            // Cancel button
            GestureDetector(
              onTap: widget.onCancel,
              child: Container(
                width: 44,
                height: 44,
                decoration: BoxDecoration(
                  color: theme.colorScheme.surface,
                  shape: BoxShape.circle,
                  border: Border.all(color: palette.divider),
                ),
                child: Icon(Icons.delete_outline_rounded,
                    size: 20, color: palette.danger),
              ),
            ),
            const SizedBox(width: 8),
            // Send recording
            _AccentCircleBtn(
              icon: Icons.send_rounded,
              onTap: widget.onStop,
            ),
          ],
        ),
      ),
    );
  }
}

// ─── Shared button widgets ────────────────────────────────────

class _CircleIconBtn extends StatelessWidget {
  final IconData icon;
  final Color color;
  final double size;
  final VoidCallback onTap;

  const _CircleIconBtn({
    required this.icon,
    required this.color,
    required this.size,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: SizedBox(
        width: size,
        height: size,
        child: Icon(icon, size: 22, color: color),
      ),
    );
  }
}

class _AccentCircleBtn extends StatelessWidget {
  final IconData icon;
  final VoidCallback onTap;

  const _AccentCircleBtn({super.key, required this.icon, required this.onTap});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return GestureDetector(
      onTap: onTap,
      child: Container(
        width: 44,
        height: 44,
        decoration: BoxDecoration(
          color: theme.colorScheme.tertiary,
          shape: BoxShape.circle,
          boxShadow: [
            BoxShadow(
              color: theme.colorScheme.tertiary.withValues(alpha: 0.5),
              blurRadius: 16,
              offset: const Offset(0, 4),
              spreadRadius: -4,
            ),
          ],
        ),
        child: Icon(icon, size: 20, color: theme.scaffoldBackgroundColor),
      ),
    );
  }
}
