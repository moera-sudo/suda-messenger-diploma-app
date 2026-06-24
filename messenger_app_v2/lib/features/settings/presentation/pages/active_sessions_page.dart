import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import '../../../../app/config/theme/app_theme.dart';
import '../../../../features/profile/data/models/session_model.dart';
import '../../../../features/profile/domain/repositories/profile_repository.dart';
import '../../../../features/auth/presentation/bloc/auth_bloc.dart';
import '../../../../shared/data/storage/secure_storage_client.dart';
import '../../../../shared/presentation/feedback/app_feedback.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/widgets/suda_button.dart';

class ActiveSessionsPage extends StatefulWidget {
  const ActiveSessionsPage({super.key});

  @override
  State<ActiveSessionsPage> createState() => _ActiveSessionsPageState();
}

class _ActiveSessionsPageState extends State<ActiveSessionsPage> {
  List<SessionModel>? _sessions;
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    if (mounted) setState(() => _isLoading = true);
    final result = await sl<ProfileRepository>().getSessions();
    if (!mounted) return;
    setState(() {
      _isLoading = false;
      _sessions = result.fold((_) => [], (list) => list);
    });
  }

  Future<void> _terminate(SessionModel session) async {
    // Optimistic removal — hide immediately so double-tap is impossible.
    setState(() => _sessions?.removeWhere((s) => s.id == session.id));

    final result = await sl<ProfileRepository>().deleteSession(session.id);
    if (!mounted) return;
    result.fold(
      (f) {
        // Rollback: restore the session and show the error.
        AppFeedback.showError(f.message);
        _load();
      },
      (_) async {
        // Terminating the current session ends our own access — force logout.
        if (session.isCurrent) {
          await sl<SecureStorageClient>().clearTokens();
          if (mounted) {
            context.read<AuthBloc>().add(AuthForcedLogout());
          }
        }
      },
    );
  }

  Future<void> _terminateAll() async {
    final l10n = context.l10n;
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (ctx) {
        final theme = Theme.of(ctx);
        final palette = theme.extension<AppPaletteTheme>()!;
        return AlertDialog(
          backgroundColor: theme.colorScheme.surface,
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(18)),
          content: Text(l10n.activeSessionsSignOutAllConfirm,
              style: TextStyle(fontFamily: 'Manrope', color: theme.colorScheme.onSurface)),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(ctx, false),
              child: Text(l10n.buttonRetry,
                  style: TextStyle(fontFamily: 'Manrope', color: palette.textSecondary)),
            ),
            TextButton(
              onPressed: () => Navigator.pop(ctx, true),
              child: Text(l10n.activeSessionsSignOutAll,
                  style: TextStyle(
                      fontFamily: 'Manrope',
                      color: palette.danger,
                      fontWeight: FontWeight.w700)),
            ),
          ],
        );
      },
    );
    if (confirmed != true || !mounted) return;

    final refreshToken = await sl<SecureStorageClient>().getRefreshToken() ?? '';
    final result = await sl<ProfileRepository>().terminateOtherSessions(refreshToken);
    if (!mounted) return;
    result.fold(
      (_) {},
      (_) {
        AppFeedback.showSuccess(l10n.activeSessionsSignOutAll);
        _load();
      },
    );
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
          l10n.activeSessionsTitle,
          style: TextStyle(
              fontFamily: 'Manrope', fontWeight: FontWeight.w700, color: theme.colorScheme.onSurface),
        ),
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : ListView(
              padding: const EdgeInsets.only(bottom: 32),
              children: [
                ...(_sessions ?? []).map((s) => _SessionRow(
                      session: s,
                      palette: palette,
                      theme: theme,
                      onTerminate: () => _terminate(s),
                      terminateLabel: l10n.activeSessionsTerminate,
                      currentLabel: l10n.activeSessionsCurrent,
                    )),
                if ((_sessions ?? []).length > 1) ...[
                  const SizedBox(height: 16),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 14),
                    child: SudaButton(
                      label: l10n.activeSessionsSignOutAll,
                      variant: SudaButtonVariant.dangerOutline,
                      size: SudaButtonSize.lg,
                      fullWidth: true,
                      onPressed: _terminateAll,
                    ),
                  ),
                ],
              ],
            ),
    );
  }
}

class _SessionRow extends StatelessWidget {
  final SessionModel session;
  final AppPaletteTheme palette;
  final ThemeData theme;
  final VoidCallback onTerminate;
  final String terminateLabel;
  final String currentLabel;

  const _SessionRow({
    required this.session,
    required this.palette,
    required this.theme,
    required this.onTerminate,
    required this.terminateLabel,
    required this.currentLabel,
  });

  String _formatDate(String iso) {
    try {
      final dt = DateTime.parse(iso).toLocal();
      return '${dt.day.toString().padLeft(2, '0')}.${dt.month.toString().padLeft(2, '0')}.${dt.year}';
    } catch (_) {
      return iso;
    }
  }

  @override
  Widget build(BuildContext context) {
    final label = session.deviceName.isNotEmpty
        ? session.deviceName
        : session.userAgent;

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 12),
      decoration: BoxDecoration(
          border: Border(bottom: BorderSide(color: palette.divider, width: 0.5))),
      child: Row(children: [
        Icon(Icons.devices_rounded, color: palette.textAccent, size: 28),
        const SizedBox(width: 12),
        Expanded(
          child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            Row(children: [
              Flexible(
                child: Text(
                  label,
                  style: TextStyle(
                      fontFamily: 'Manrope',
                      fontSize: 14,
                      fontWeight: FontWeight.w600,
                      color: theme.colorScheme.onSurface),
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
              ),
              if (session.isCurrent) ...[
                const SizedBox(width: 8),
                Container(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                  decoration: BoxDecoration(
                    color: palette.success.withValues(alpha: 0.15),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Text(
                    currentLabel,
                    style: TextStyle(
                        fontFamily: 'Manrope',
                        fontSize: 11,
                        fontWeight: FontWeight.w700,
                        color: palette.success),
                  ),
                ),
              ],
            ]),
            if (session.clientIp.isNotEmpty)
              Text(
                session.clientIp,
                style: TextStyle(fontFamily: 'Manrope', fontSize: 12, color: palette.textSecondary),
              ),
            Text(
              _formatDate(session.createdAt),
              style: TextStyle(fontFamily: 'Manrope', fontSize: 12, color: palette.textSecondary),
            ),
          ]),
        ),
        TextButton(
          onPressed: onTerminate,
          child: Text(terminateLabel,
              style: TextStyle(fontFamily: 'Manrope', color: palette.danger, fontWeight: FontWeight.w600)),
        ),
      ]),
    );
  }
}
