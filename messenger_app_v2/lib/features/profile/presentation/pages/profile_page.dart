import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:image_picker/image_picker.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import '../../../../app/config/theme/app_theme.dart';
import '../../../../app/navigation/app_routes.dart';
import '../../../../features/auth/presentation/bloc/auth_bloc.dart';
import '../../../../features/chat/presentation/bloc/socket_bloc.dart';
import '../../../../features/wallet/data/models/wallet_models.dart';
import '../../../../shared/presentation/feedback/app_feedback.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/widgets/suda_avatar.dart';
import '../../../../shared/presentation/widgets/suda_info_row.dart';
import '../bloc/profile_bloc.dart';

class ProfilePage extends StatefulWidget {
  const ProfilePage({super.key});

  @override
  State<ProfilePage> createState() => _ProfilePageState();
}

class _ProfilePageState extends State<ProfilePage> {
  final _picker = ImagePicker();

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) => sl<ProfileBloc>()..add(const ProfileLoadRequested()),
      child: BlocConsumer<ProfileBloc, ProfileState>(
        listener: (context, state) {
          final l10n = context.l10n;
          if (state.actionStatus == ProfileActionStatus.failure) {
            final msg = state.actionMessage;
            if (msg != null && msg.isNotEmpty) AppFeedback.showError(msg);
          } else if (state.actionStatus == ProfileActionStatus.saveSuccess) {
            AppFeedback.showSuccess(l10n.profileUpdated);
          } else if (state.actionStatus == ProfileActionStatus.uploadSuccess) {
            AppFeedback.showSuccess(l10n.avatarUpdated);
          }

          if (state.actionStatus == ProfileActionStatus.logoutSuccess) {
            AppFeedback.showSuccess(l10n.loggedOut);
            context.read<SocketBloc>().add(SocketDisconnect());
            context.read<AuthBloc>().add(AuthCheckRequested());
            context.goNamed(AppRoutes.welcome);
          }
        },
        builder: (context, state) {
          final palette = Theme.of(context).extension<AppPaletteTheme>()!;
          final theme = Theme.of(context);
          final l10n = context.l10n;

          return Scaffold(
            backgroundColor: theme.scaffoldBackgroundColor,
            body: state.status == ProfileViewStatus.loading
                ? const Center(child: CircularProgressIndicator())
                : state.status == ProfileViewStatus.failure
                    ? _ErrorView(failure: state.failure)
                    : CustomScrollView(
                        slivers: [
                          // ── Large AppBar ──────────────────────────────
                          SliverAppBar(
                            pinned: true,
                            expandedHeight: 0,
                            backgroundColor: theme.scaffoldBackgroundColor,
                            elevation: 0,
                            scrolledUnderElevation: 0,
                            title: Text(
                              l10n.profileTitle,
                              style: TextStyle(
                                fontFamily: 'Manrope',
                                fontSize: 18,
                                fontWeight: FontWeight.w700,
                                color: theme.colorScheme.onSurface,
                              ),
                            ),
                          ),

                          SliverToBoxAdapter(
                            child: Column(
                              children: [
                                // ── Hero section ─────────────────────
                                Stack(
                                  children: [
                                    // Background radial glow
                                    Container(
                                      height: 200,
                                      decoration: BoxDecoration(
                                        gradient: RadialGradient(
                                          center: Alignment.center,
                                          radius: 1.0,
                                          colors: [
                                            theme.colorScheme.primary.withValues(alpha: 0.45),
                                            Colors.transparent,
                                          ],
                                        ),
                                      ),
                                    ),
                                    Center(
                                      child: Padding(
                                      padding: const EdgeInsets.only(bottom: 20, top: 12),
                                      child: Column(
                                        crossAxisAlignment: CrossAxisAlignment.center,
                                        children: [
                                          // Avatar with camera button
                                          Stack(
                                            children: [
                                              SudaAvatar(
                                                mediaUrl: _cacheBusted(state.profile?.avatarUrl, state.avatarVersion),
                                                initials: state.profile?.displayName ?? '?',
                                                size: 88,
                                                ring: true,
                                              ),
                                              Positioned(
                                                bottom: -2,
                                                right: -2,
                                                child: GestureDetector(
                                                  onTap: state.isUploadingAvatar
                                                      ? null
                                                      : () => _pickAvatar(context),
                                                  child: Container(
                                                    width: 30,
                                                    height: 30,
                                                    decoration: BoxDecoration(
                                                      color: theme.colorScheme.tertiary,
                                                      shape: BoxShape.circle,
                                                      border: Border.all(
                                                        color: theme.scaffoldBackgroundColor,
                                                        width: 2.5,
                                                      ),
                                                    ),
                                                    child: state.isUploadingAvatar
                                                        ? Padding(
                                                            padding: const EdgeInsets.all(6),
                                                            child: CircularProgressIndicator(
                                                              strokeWidth: 2,
                                                              color: theme.scaffoldBackgroundColor,
                                                            ),
                                                          )
                                                        : Icon(
                                                            Icons.camera_alt_rounded,
                                                            size: 14,
                                                            color: theme.scaffoldBackgroundColor,
                                                          ),
                                                  ),
                                                ),
                                              ),
                                            ],
                                          ),

                                          const SizedBox(height: 12),

                                          Text(
                                            state.profile?.displayName ?? '',
                                            style: TextStyle(
                                              fontFamily: 'SpaceGrotesk',
                                              fontSize: 22,
                                              fontWeight: FontWeight.w800,
                                              letterSpacing: -0.03,
                                              color: theme.colorScheme.onSurface,
                                            ),
                                          ),

                                          const SizedBox(height: 2),

                                          if (state.profile?.username.isNotEmpty == true)
                                            Text(
                                              state.profile!.username.startsWith('@')
                                                  ? state.profile!.username
                                                  : '@${state.profile!.username}',
                                              style: TextStyle(
                                                fontFamily: 'Manrope',
                                                fontSize: 14,
                                                fontWeight: FontWeight.w600,
                                                color: theme.colorScheme.tertiary,
                                              ),
                                            ),

                                          if (state.profile?.bio.isNotEmpty == true) ...[
                                            const SizedBox(height: 4),
                                            Padding(
                                              padding: const EdgeInsets.symmetric(horizontal: 32),
                                              child: Text(
                                                state.profile!.bio,
                                                textAlign: TextAlign.center,
                                                style: TextStyle(
                                                  fontFamily: 'Manrope',
                                                  fontSize: 13,
                                                  color: palette.textSecondary,
                                                  height: 1.4,
                                                ),
                                              ),
                                            ),
                                          ],
                                        ],
                                      ),
                                    ),
                                    ), // Center
                                  ],
                                ),

                                // ── Balance card ──────────────────────
                                Padding(
                                  padding: const EdgeInsets.fromLTRB(14, 0, 14, 10),
                                  child: Container(
                                    decoration: BoxDecoration(
                                      gradient: LinearGradient(
                                        begin: Alignment.topLeft,
                                        end: Alignment.bottomRight,
                                        colors: [
                                          Color.lerp(theme.colorScheme.primary, theme.colorScheme.surface, 0.5)!,
                                          palette.surface2,
                                        ],
                                      ),
                                      borderRadius: BorderRadius.circular(18),
                                      border: Border.all(
                                        color: theme.colorScheme.primary.withValues(alpha: 0.4),
                                      ),
                                    ),
                                    padding: const EdgeInsets.all(16),
                                    child: Row(
                                      children: [
                                        // Balance info
                                        Expanded(
                                          child: Column(
                                            crossAxisAlignment: CrossAxisAlignment.start,
                                            children: [
                                              Text(
                                                l10n.profileSudaBalance.toUpperCase(),
                                                style: TextStyle(
                                                  fontFamily: 'Manrope',
                                                  fontSize: 11,
                                                  fontWeight: FontWeight.w700,
                                                  letterSpacing: 0.10,
                                                  color: palette.textSecondary,
                                                ),
                                              ),
                                              const SizedBox(height: 4),
                                              Row(
                                                crossAxisAlignment: CrossAxisAlignment.baseline,
                                                textBaseline: TextBaseline.alphabetic,
                                                children: [
                                                  Text(
                                                    state.wallet != null
                                                        ? '${weiToSudaDisplay(state.wallet!.sudaBalanceWei)} '
                                                        : '— ',
                                                    style: TextStyle(
                                                      fontFamily: 'SpaceGrotesk',
                                                      fontSize: 26,
                                                      fontWeight: FontWeight.w800,
                                                      color: theme.colorScheme.onSurface,
                                                    ),
                                                  ),
                                                  Text(
                                                    l10n.walletSuda,
                                                    style: TextStyle(
                                                      fontFamily: 'Manrope',
                                                      fontSize: 14,
                                                      fontWeight: FontWeight.w600,
                                                      color: palette.textSecondary,
                                                    ),
                                                  ),
                                                ],
                                              ),
                                              const SizedBox(height: 4),
                                              Row(
                                                children: [
                                                  Icon(Icons.lock_outline_rounded, size: 11, color: palette.textSecondary),
                                                  const SizedBox(width: 4),
                                                  Text(
                                                    '0x…••••',
                                                    style: TextStyle(
                                                      fontFamily: 'JetBrainsMono',
                                                      fontSize: 12,
                                                      color: palette.textSecondary,
                                                    ),
                                                  ),
                                                ],
                                              ),
                                            ],
                                          ),
                                        ),
                                        // Action buttons
                                        Column(
                                          children: [
                                            _BalanceBtn(
                                              icon: Icons.account_balance_wallet_rounded,
                                              label: l10n.profileOpenWallet,
                                              accent: true,
                                              onTap: () => context.pushNamed(AppRoutes.wallet),
                                            ),
                                          ],
                                        ),
                                      ],
                                    ),
                                  ),
                                ),

                                // ── Quick action tiles ───────────────
                                Padding(
                                  padding: const EdgeInsets.fromLTRB(14, 0, 14, 8),
                                  child: Container(
                                    decoration: BoxDecoration(
                                      color: theme.colorScheme.surface,
                                      borderRadius: BorderRadius.circular(16),
                                      border: Border.all(color: palette.divider),
                                    ),
                                    child: IntrinsicHeight(
                                      child: Row(
                                        children: [
                                          _ProfileActionTile(
                                            icon: Icons.camera_alt_rounded,
                                            label: l10n.profileSetPhoto,
                                            onTap: state.isUploadingAvatar
                                                ? null
                                                : () => _pickAvatar(context),
                                          ),
                                          VerticalDivider(width: 1, color: palette.divider),
                                          _ProfileActionTile(
                                            icon: Icons.person_outline_rounded,
                                            label: l10n.profileEditProfile,
                                            onTap: () => context.pushNamed(AppRoutes.editProfile),
                                          ),
                                          VerticalDivider(width: 1, color: palette.divider),
                                          _ProfileActionTile(
                                            icon: Icons.settings_outlined,
                                            label: l10n.settingsTitle,
                                            onTap: () => context.pushNamed(AppRoutes.settings),
                                          ),
                                        ],
                                      ),
                                    ),
                                  ),
                                ),

                                // ── Quick actions ─────────────────────
                                // Full settings (theme/privacy/language/etc.) → gear icon.
                                const SizedBox(height: 4),
                                SudaInfoList(
                                  children: [
                                    SudaInfoRow(
                                      icon: Icons.logout_rounded,
                                      label: l10n.profileSignOut,
                                      danger: true,
                                      showChevron: false,
                                      onTap: state.isLoggingOut
                                          ? null
                                          : () => context.read<ProfileBloc>().add(const ProfileLogoutRequested()),
                                    ),
                                  ],
                                ),

                                const SizedBox(height: 32),
                              ],
                            ),
                          ),
                        ],
                      ),
          );
        },
      ),
    );
  }

  String? _cacheBusted(String? url, int version) {
    // Presigned S3 URLs sign their query string — appending `&v=N` breaks the
    // signature (403). Each upload/reload yields a fresh presigned URL, so the
    // string already changes and busts the image cache on its own.
    if (url == null || url.trim().isEmpty) return null;
    return url.trim();
  }

  Future<void> _pickAvatar(BuildContext context) async {
    final file = await _picker.pickImage(source: ImageSource.gallery, imageQuality: 88, maxWidth: 1200);
    if (!context.mounted || file == null) return;
    context.read<ProfileBloc>().add(ProfileAvatarUploadRequested(filePath: file.path, fileName: file.name));
  }
}

// ── Small helpers ──────────────────────────────────────────────

class _ErrorView extends StatelessWidget {
  final dynamic failure;
  const _ErrorView({this.failure});

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    return Center(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(Icons.error_outline_rounded, color: palette.danger, size: 44),
          const SizedBox(height: 12),
          Text(
            AppFeedback.failureMessage(failure, fallback: context.l10n.errorGeneric),
            textAlign: TextAlign.center,
            style: TextStyle(color: palette.textSecondary, height: 1.4),
          ),
          const SizedBox(height: 16),
          FilledButton.icon(
            onPressed: () => context.read<ProfileBloc>().add(const ProfileLoadRequested()),
            icon: const Icon(Icons.refresh_rounded),
            label: Text(context.l10n.buttonRetry),
          ),
        ],
      ),
    );
  }
}

class _ProfileActionTile extends StatelessWidget {
  final IconData icon;
  final String label;
  final VoidCallback? onTap;
  const _ProfileActionTile({required this.icon, required this.label, this.onTap});

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);
    return Expanded(
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(16),
        child: Padding(
          padding: const EdgeInsets.symmetric(vertical: 14, horizontal: 8),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(icon, size: 20, color: theme.colorScheme.tertiary),
              const SizedBox(height: 6),
              Text(
                label,
                style: TextStyle(
                  fontFamily: 'Manrope',
                  fontSize: 12,
                  color: palette.textSecondary,
                ),
                textAlign: TextAlign.center,
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _BalanceBtn extends StatelessWidget {
  final IconData icon;
  final String label;
  final bool accent;
  final VoidCallback onTap;
  const _BalanceBtn({required this.icon, required this.label, required this.accent, required this.onTap});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        decoration: BoxDecoration(
          color: accent ? theme.colorScheme.tertiary : Colors.white.withValues(alpha: 0.08),
          borderRadius: BorderRadius.circular(10),
        ),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              icon,
              size: 16,
              color: accent ? theme.scaffoldBackgroundColor : theme.colorScheme.onSurface,
            ),
            const SizedBox(width: 6),
            Text(
              label,
              style: TextStyle(
                fontFamily: 'Manrope',
                fontSize: 13,
                fontWeight: FontWeight.w600,
                color: accent ? theme.scaffoldBackgroundColor : theme.colorScheme.onSurface,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
