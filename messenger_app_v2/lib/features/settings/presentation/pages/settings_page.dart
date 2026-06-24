import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import '../../../../app/config/app_config.dart';
import '../../../../app/config/theme/app_theme.dart';
import '../../../../app/navigation/app_routes.dart';
import '../../../../features/auth/presentation/bloc/auth_bloc.dart';
import '../../../../features/chat/presentation/bloc/socket_bloc.dart';
import '../../../../features/preferences/presentation/bloc/preferences_cubit.dart';
import '../../../../features/profile/presentation/bloc/profile_bloc.dart';
import '../../../../shared/presentation/bloc/settings_cubit.dart';
import '../../../../shared/presentation/feedback/app_feedback.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/widgets/suda_info_row.dart';
import 'theme_picker_page.dart';

class SettingsPage extends StatelessWidget {
  const SettingsPage({super.key});

  @override
  Widget build(BuildContext context) {
    return MultiBlocProvider(
      providers: [
        BlocProvider(create: (_) => sl<ProfileBloc>()..add(const ProfileLoadRequested())),
        BlocProvider(create: (_) => sl<PreferencesCubit>()..load()),
      ],
      child: const _SettingsView(),
    );
  }
}

class _SettingsView extends StatelessWidget {
  const _SettingsView();

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    return BlocListener<ProfileBloc, ProfileState>(
      listenWhen: (p, c) => p.actionStatus != c.actionStatus,
      listener: (context, state) {
        if (state.actionStatus == ProfileActionStatus.logoutSuccess) {
          AppFeedback.showSuccess(l10n.loggedOut);
          context.read<SocketBloc>().add(SocketDisconnect());
          context.read<AuthBloc>().add(AuthCheckRequested());
          context.goNamed(AppRoutes.welcome);
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
            l10n.settingsTitle,
            style: TextStyle(fontFamily: 'Manrope', fontWeight: FontWeight.w700, color: theme.colorScheme.onSurface),
          ),
        ),
        body: ListView(
          padding: const EdgeInsets.only(top: 8, bottom: 24),
          children: [
            // ── Account ───────────────────────────────────────
            SudaInfoList(
              title: l10n.settingsAccount,
              children: [
                SudaInfoRow(
                  icon: Icons.person_outline_rounded,
                  label: l10n.profileEditProfile,
                  onTap: () => context.pushNamed(AppRoutes.editProfile),
                ),
                SudaInfoRow(
                  icon: Icons.lock_outline_rounded,
                  label: l10n.settingsChangePassword,
                  onTap: () => context.pushNamed(AppRoutes.changePassword),
                ),
                SudaInfoRow(
                  icon: Icons.verified_user_outlined,
                  label: l10n.settingsActiveSessions,
                  onTap: () => context.pushNamed(AppRoutes.activeSessions),
                ),
                SudaInfoRow(
                  icon: Icons.mail_outline_rounded,
                  label: l10n.settingsChannelInvitations,
                  onTap: () => context.pushNamed(AppRoutes.channelInvitations),
                ),
                SudaInfoRow(
                  icon: Icons.block_rounded,
                  label: l10n.settingsBlocked,
                  onTap: () => context.pushNamed(AppRoutes.blocked),
                ),
              ],
            ),
            const SizedBox(height: 16),

            // ── Appearance ────────────────────────────────────
            BlocBuilder<PreferencesCubit, PreferencesState>(
              buildWhen: (p, c) => p.prefs.language != c.prefs.language,
              builder: (context, state) {
                final prefs = state.prefs;
                return SudaInfoList(
                  title: l10n.settingsAppearance,
                  children: [
                    SudaInfoRow(
                      icon: Icons.language_rounded,
                      label: l10n.profileLanguage,
                      right: InfoMeta(_languageLabel(l10n, context.watch<SettingsCubit>().state.locale.languageCode)),
                      onTap: () => _showLanguageSheet(context, prefs),
                    ),
                    SudaInfoRow(
                      icon: Icons.palette_outlined,
                      label: l10n.settingsTheme,
                      right: InfoMeta(themeLabel(context.watch<SettingsCubit>().state.themeType)),
                      onTap: () => context.pushNamed(AppRoutes.themePicker),
                    ),
                  ],
                );
              },
            ),
            const SizedBox(height: 16),

            // ── Suda ──────────────────────────────────────────
            SudaInfoList(
              title: l10n.settingsSuda,
              children: [
                SudaInfoRow(
                  icon: Icons.account_balance_wallet_outlined,
                  label: l10n.settingsWallet,
                  onTap: () => context.pushNamed(AppRoutes.wallet),
                ),
                SudaInfoRow(
                  icon: Icons.info_outline_rounded,
                  label: l10n.settingsAboutSuda,
                  onTap: () => _showAboutSheet(context),
                ),
                SudaInfoRow(
                  icon: Icons.shield_outlined,
                  label: l10n.settingsPrivacyPolicy,
                  onTap: () => _comingSoon(context),
                ),
              ],
            ),
            const SizedBox(height: 16),

            // ── Sign out ──────────────────────────────────────
            BlocBuilder<ProfileBloc, ProfileState>(
              buildWhen: (p, c) => p.isLoggingOut != c.isLoggingOut,
              builder: (context, state) => SudaInfoList(
                children: [
                  SudaInfoRow(
                    icon: Icons.logout_rounded,
                    label: l10n.settingsSignOut,
                    danger: true,
                    showChevron: false,
                    onTap: state.isLoggingOut
                        ? null
                        : () => context.read<ProfileBloc>().add(const ProfileLogoutRequested()),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 16),

            // ── Footer ────────────────────────────────────────
            Center(
              child: Text(
                l10n.settingsVersion(AppConfig.appVersion),
                style: TextStyle(fontFamily: 'Manrope', fontSize: 12, color: palette.textSecondary),
              ),
            ),
          ],
        ),
      ),
    );
  }

  // ── Helpers ─────────────────────────────────────────────────

  String _languageLabel(dynamic l10n, String code) => switch (code) {
        'ru' => l10n.languageRussian,
        'kk' => l10n.languageKazakh,
        _ => l10n.languageEnglish,
      };

  void _showLanguageSheet(BuildContext context, prefs) {
    final l10n = context.l10n;
    final settings = context.read<SettingsCubit>();
    final prefsCubit = context.read<PreferencesCubit>();
    _showOptionSheet(
      context,
      title: l10n.languageTitle,
      current: settings.state.locale.languageCode,
      options: {
        'en': l10n.languageEnglish,
        'ru': l10n.languageRussian,
        'kk': l10n.languageKazakh,
      },
      onSelected: (code) {
        settings.changeLocale(code);
        prefsCubit.update(prefs.copyWith(language: code), {'language': code});
      },
    );
  }

  void _showOptionSheet(
    BuildContext context, {
    required String title,
    required String current,
    required Map<String, String> options,
    required ValueChanged<String> onSelected,
  }) {
    final theme = Theme.of(context);

    showModalBottomSheet(
      context: context,
      backgroundColor: theme.colorScheme.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      ),
      builder: (ctx) => SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(20, 18, 20, 8),
              child: Text(
                title,
                style: TextStyle(
                  fontFamily: 'Manrope',
                  fontSize: 16,
                  fontWeight: FontWeight.w700,
                  color: theme.colorScheme.onSurface,
                ),
              ),
            ),
            ...options.entries.map((e) {
              final selected = e.key == current;
              return ListTile(
                title: Text(
                  e.value,
                  style: TextStyle(fontFamily: 'Manrope', color: theme.colorScheme.onSurface),
                ),
                trailing: selected
                    ? Icon(Icons.check_rounded, color: theme.colorScheme.tertiary)
                    : null,
                onTap: () {
                  Navigator.pop(ctx);
                  if (!selected) onSelected(e.key);
                },
              );
            }),
            const SizedBox(height: 8),
          ],
        ),
      ),
    );
  }

  void _showAboutSheet(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    showModalBottomSheet(
      context: context,
      backgroundColor: theme.colorScheme.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      ),
      builder: (_) => SafeArea(
        child: Padding(
          padding: const EdgeInsets.fromLTRB(24, 24, 24, 32),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Text(
                'Suda',
                style: TextStyle(
                  fontFamily: 'SpaceGrotesk',
                  fontSize: 26,
                  fontWeight: FontWeight.w800,
                  color: theme.colorScheme.tertiary,
                ),
              ),
              const SizedBox(height: 6),
              Text(
                l10n.settingsVersion(AppConfig.appVersion),
                style: TextStyle(fontFamily: 'Manrope', fontSize: 13, color: palette.textSecondary),
              ),
              const SizedBox(height: 14),
              Text(
                l10n.settingsAboutDesc,
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontFamily: 'Manrope',
                  fontSize: 14,
                  height: 1.5,
                  color: theme.colorScheme.onSurface,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  void _comingSoon(BuildContext context) {
    AppFeedback.showSuccess(context.l10n.comingSoon);
  }
}
