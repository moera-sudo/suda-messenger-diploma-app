import 'package:audio_waveforms/audio_waveforms.dart';
import 'package:flutter/material.dart';

import '../../../../features/media/domain/repositories/media_repository.dart';

/// Displays a voice message bubble with waveform and play/stop controls.
class VoiceMessagePlayer extends StatefulWidget {
  final String? mediaId;
  final String? durationText;
  final bool isMe;
  final MediaRepository? mediaRepo;

  const VoiceMessagePlayer({
    super.key,
    this.mediaId,
    this.durationText,
    this.isMe = false,
    this.mediaRepo,
  });

  @override
  State<VoiceMessagePlayer> createState() => _VoiceMessagePlayerState();
}

class _VoiceMessagePlayerState extends State<VoiceMessagePlayer> {
  PlayerController? _player;
  bool _isPlaying = false;
  bool _isLoading = false;
  // audio_waveforms plays local files only — cache the downloaded path.
  String? _localPath;
  bool _prepared = false;

  Future<void> _togglePlay() async {
    if (_isLoading) return; // guard against double taps while preparing
    if (_isPlaying) {
      await _player?.pausePlayer();
      if (mounted) setState(() => _isPlaying = false);
      return;
    }
    if (widget.mediaId == null || widget.mediaRepo == null) return;

    setState(() => _isLoading = true);
    try {
      // Download once, then reuse the cached file on subsequent plays.
      if (_localPath == null) {
        final res = await widget.mediaRepo!
            .downloadToCache(widget.mediaId!, 'voice_${widget.mediaId}.m4a');
        final path = res.fold((_) => null, (p) => p);
        if (path == null) {
          if (mounted) setState(() => _isLoading = false);
          return;
        }
        _localPath = path;
      }

      // Recreate a fresh controller for every playback: audio_waveforms cannot
      // restart a player that has finished (it sits in `stopped`), so a clean
      // controller is the reliable way to allow replays.
      await _disposePlayer();
      final p = PlayerController()..updateFrequency = UpdateFrequency.high;
      await p.preparePlayer(path: _localPath!, shouldExtractWaveform: true);
      p.onPlayerStateChanged.listen((s) {
        if (mounted) setState(() => _isPlaying = s == PlayerState.playing);
      });
      _player = p;
      _prepared = true;
      await p.startPlayer();
      if (mounted) setState(() { _isPlaying = true; _isLoading = false; });
    } catch (_) {
      if (mounted) setState(() { _isLoading = false; });
    }
  }

  /// Tears down the current controller so the next play starts clean.
  Future<void> _disposePlayer() async {
    final p = _player;
    _player = null;
    _prepared = false;
    if (p == null) return;
    try {
      await p.stopPlayer();
    } catch (_) {/* already stopped */}
    p.dispose();
  }

  @override
  void dispose() {
    _player?.dispose();
    super.dispose();
  }

  String get _durationLabel {
    if (widget.durationText == null) return '0:00';
    final secs = int.tryParse(widget.durationText!) ?? 0;
    final m = secs ~/ 60;
    final s = secs % 60;
    return '$m:${s.toString().padLeft(2, '0')}';
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return RepaintBoundary(
      child: SizedBox(
        width: 220,
        child: Row(
          children: [
            // Play / Stop / Loading button
            GestureDetector(
              onTap: _togglePlay,
              child: Container(
                width: 32,
                height: 32,
                decoration: const BoxDecoration(
                  color: Color(0x2DFFFFFF),
                  shape: BoxShape.circle,
                ),
                child: _isLoading
                    ? const SizedBox(
                        width: 16,
                        height: 16,
                        child: CircularProgressIndicator(
                          strokeWidth: 2,
                          color: Colors.white70,
                        ),
                      )
                    : Icon(
                        _isPlaying ? Icons.stop_rounded : Icons.play_arrow_rounded,
                        size: 18,
                        color: Colors.white,
                      ),
              ),
            ),
            const SizedBox(width: 8),

            // Waveform or placeholder bars
            Expanded(
              child: _prepared && _player != null
                  ? AudioFileWaveforms(
                      size: const Size(double.infinity, 28),
                      playerController: _player!,
                      waveformType: WaveformType.fitWidth,
                      playerWaveStyle: PlayerWaveStyle(
                        fixedWaveColor: Colors.white.withValues(alpha: 0.4),
                        liveWaveColor: theme.colorScheme.tertiary,
                        seekLineColor: Colors.transparent,
                        waveThickness: 2.5,
                      ),
                    )
                  : _StaticWaveform(isMe: widget.isMe),
            ),
            const SizedBox(width: 8),

            // Duration
            Text(
              _durationLabel,
              style: const TextStyle(
                fontFamily: 'JetBrainsMono',
                fontSize: 12,
                color: Color(0x8CFFFFFF),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

/// Static waveform placeholder shown before audio URL is loaded.
class _StaticWaveform extends StatelessWidget {
  final bool isMe;
  const _StaticWaveform({required this.isMe});

  static const _heights = [6, 10, 14, 8, 18, 12, 16, 6, 14, 10, 18, 8, 12, 16, 6];

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: 28,
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        crossAxisAlignment: CrossAxisAlignment.center,
        children: _heights.map((h) => Container(
          width: 2,
          height: h.toDouble(),
          decoration: BoxDecoration(
            color: Colors.white.withValues(alpha: 0.45),
            borderRadius: BorderRadius.circular(1),
          ),
        )).toList(),
      ),
    );
  }
}
