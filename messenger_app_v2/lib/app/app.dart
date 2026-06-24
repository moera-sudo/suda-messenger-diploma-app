import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import 'package:messenger_app_v2/app/config/theme/app_theme.dart';
import 'package:messenger_app_v2/app/navigation/app_router.dart';
import 'package:messenger_app_v2/app/navigation/app_routes.dart';
import 'package:messenger_app_v2/features/auth/presentation/bloc/auth_bloc.dart';
import 'package:messenger_app_v2/features/chat/presentation/bloc/socket_bloc.dart';
import 'package:messenger_app_v2/features/preferences/domain/repositories/preferences_repository.dart';
import 'package:messenger_app_v2/l10n/arb/app_localizations.dart';
import 'package:messenger_app_v2/shared/data/logger/app_logger.dart';
import 'package:messenger_app_v2/shared/data/notifications/notification_service.dart';
import 'package:messenger_app_v2/shared/presentation/bloc/settings_cubit.dart';
import 'package:messenger_app_v2/shared/presentation/feedback/app_feedback.dart';

class MainApp extends StatelessWidget {
  const MainApp({super.key});

  @override
  Widget build(BuildContext context) {
    final appRouter = sl<AppRouter>();

    return MultiBlocProvider(
      providers: [
        BlocProvider(create: (_) => sl<SettingsCubit>()),
        BlocProvider(create: (_) => sl<AuthBloc>()),
        BlocProvider(create: (_) => sl<SocketBloc>()),
      ],
      child: BlocListener<AuthBloc, AuthState>(
        listenWhen: (prev, curr) => prev.status != curr.status,
        listener: (context, state) {
          // On successful auth — pull preferences and apply theme + locale.
          if (state.status == AuthStatus.authenticated) {
            _applyServerPreferences(context);
            // Register FCM token now that the user is authenticated.
            NotificationService.registerToken();
            // Navigate to a specific chat if the user tapped a notification.
            final chatId = NotificationService.consumePendingChatId();
            if (chatId != null) {
              WidgetsBinding.instance.addPostFrameCallback((_) {
                sl<AppRouter>().router.pushNamed(
                  AppRoutes.chatDetail,
                  pathParameters: {'id': chatId},
                  extra: {'name': '', 'chatType': 'DIRECT'},
                );
              });
            }
            return;
          }
          // Forced logout while using the app (invalid refresh token / revoked
          // session) — tear down the socket and redirect to Welcome.
          if (state.status == AuthStatus.unauthenticated) {
            context.read<SocketBloc>().add(SocketDisconnect());
            sl<AppRouter>().router.goNamed(AppRoutes.welcome);
          }
        },
        child: BlocBuilder<SettingsCubit, SettingsState>(
          builder: (context, state) {
            return MaterialApp.router(
              title: 'Suda Messenger',
              debugShowCheckedModeBanner: false,
              scaffoldMessengerKey: AppFeedback.messengerKey,
              theme: state.themeData,
              locale: state.locale,
              supportedLocales: AppLocalizations.supportedLocales,
              localizationsDelegates: AppLocalizations.localizationsDelegates,
              routerConfig: appRouter.router,
            );
          },
        ),
      ),
    );
  }

  /// Loads preferences from server and applies theme + locale locally.
  /// Fire-and-forget — errors are logged but not surfaced to the user.
  Future<void> _applyServerPreferences(BuildContext context) async {
    final logger = sl<AppLogger>();
    final result = await sl<PreferencesRepository>().getPreferences();
    result.fold(
      (failure) => logger.warning('Could not load server preferences: ${failure.message}'),
      (prefs) {
        if (!context.mounted) return;
        final settings = context.read<SettingsCubit>();

        // Apply theme — stored as AppThemeType.name (e.g. 'suda', 'teaChatsDark').
        // Ignore legacy values like 'LIGHT'/'DARK'/'AUTO' that have no AppThemeType match.
        try {
          final themeType = AppThemeType.values.byName(prefs.theme);
          settings.changeTheme(themeType);
          logger.debug('Applied server theme: ${prefs.theme}');
        } catch (_) {
          logger.debug('Unknown theme value from server: ${prefs.theme} — keeping local');
        }

        // Apply language.
        if (prefs.language.isNotEmpty) {
          settings.changeLocale(prefs.language);
          logger.debug('Applied server language: ${prefs.language}');
        }
      },
    );
  }
}
