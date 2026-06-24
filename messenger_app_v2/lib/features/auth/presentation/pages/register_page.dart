import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import '../../../../app/config/theme/app_theme.dart';
import '../../../../app/navigation/app_routes.dart';
import '../../../../shared/presentation/feedback/app_feedback.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/widgets/suda_button.dart';
import '../../../../shared/presentation/widgets/suda_text_field.dart';
import '../bloc/auth_bloc.dart';

class RegisterPage extends StatefulWidget {
  const RegisterPage({super.key});

  @override
  State<RegisterPage> createState() => _RegisterPageState();
}

class _RegisterPageState extends State<RegisterPage> {
  final _usernameCtrl = TextEditingController();
  final _displayNameCtrl = TextEditingController();
  final _emailCtrl = TextEditingController();
  final _passCtrl = TextEditingController();
  bool _obscurePass = true;

  @override
  void dispose() {
    _usernameCtrl.dispose();
    _displayNameCtrl.dispose();
    _emailCtrl.dispose();
    _passCtrl.dispose();
    super.dispose();
  }

  // @ prefix is added when sending to API, not shown to the user
  String get _usernameForApi {
    final raw = _usernameCtrl.text.trim();
    return raw.startsWith('@') ? raw : '@$raw';
  }

  bool get _canSubmit =>
      _usernameCtrl.text.trim().length >= 2 &&
      _emailCtrl.text.contains('@') &&
      _passCtrl.text.length >= 6;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);

    return BlocProvider(
      create: (context) => sl<AuthBloc>(),
      child: BlocConsumer<AuthBloc, AuthState>(
        listener: (context, state) {
          if (state.status == AuthStatus.success) {
            context.pushNamed(AppRoutes.verify, extra: _emailCtrl.text.trim());
          }
          if (state.status == AuthStatus.failure) {
            AppFeedback.showError(
              AppFeedback.failureMessage(state.failure, fallback: l10n.errorGeneric),
            );
          }
        },
        builder: (context, state) {
          final isLoading = state.status == AuthStatus.loading;

          return Scaffold(
            backgroundColor: theme.scaffoldBackgroundColor,
            body: SafeArea(
              child: SingleChildScrollView(
                padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 8),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // Back button
                    IconButton(
                      icon: Icon(Icons.arrow_back_rounded, color: theme.colorScheme.onSurface),
                      onPressed: () => context.pop(),
                      padding: EdgeInsets.zero,
                    ),

                    const SizedBox(height: 20),

                    // Headline
                    Text(
                      l10n.authCreateAccount,
                      style: TextStyle(
                        fontFamily: 'SpaceGrotesk',
                        fontSize: 26,
                        fontWeight: FontWeight.w800,
                        letterSpacing: -0.03,
                        color: theme.colorScheme.onSurface,
                      ),
                    ),
                    const SizedBox(height: 6),
                    Text(
                      l10n.authJoinFuture,
                      style: TextStyle(
                        fontFamily: 'Manrope',
                        fontSize: 14,
                        color: palette.textSecondary,
                        height: 1.5,
                      ),
                    ),

                    const SizedBox(height: 32),

                    // Form fields
                    SudaTextField(
                      label: l10n.authUsername,
                      hint: 'yourhandle',
                      controller: _usernameCtrl,
                      leftIcon: Icons.alternate_email_rounded,
                      hint2: l10n.authUsernameHint,
                      textInputAction: TextInputAction.next,
                    ),
                    const SizedBox(height: 16),
                    SudaTextField(
                      label: l10n.authDisplayName,
                      hint: 'Your name',
                      controller: _displayNameCtrl,
                      textInputAction: TextInputAction.next,
                    ),
                    const SizedBox(height: 16),
                    SudaTextField(
                      label: l10n.authEmail,
                      hint: 'you@example.com',
                      controller: _emailCtrl,
                      leftIcon: Icons.mail_outline_rounded,
                      keyboardType: TextInputType.emailAddress,
                      textInputAction: TextInputAction.next,
                    ),
                    const SizedBox(height: 16),
                    ListenableBuilder(
                      listenable: Listenable.merge([_usernameCtrl, _emailCtrl, _passCtrl]),
                      builder: (context, _) => SudaTextField(
                        label: l10n.authPassword,
                        hint: 'Minimum 6 characters',
                        controller: _passCtrl,
                        leftIcon: Icons.lock_outline_rounded,
                        obscureText: _obscurePass,
                        rightIcon: _obscurePass ? Icons.visibility_off_outlined : Icons.visibility_outlined,
                        onRightTap: () => setState(() => _obscurePass = !_obscurePass),
                        textInputAction: TextInputAction.done,
                      ),
                    ),

                    const SizedBox(height: 28),

                    ListenableBuilder(
                      listenable: Listenable.merge([_usernameCtrl, _emailCtrl, _passCtrl]),
                      builder: (context, _) => SudaButton(
                        label: l10n.authRegister,
                        variant: SudaButtonVariant.primary,
                        size: SudaButtonSize.lg,
                        fullWidth: true,
                        loading: isLoading,
                        onPressed: _canSubmit && !isLoading
                            ? () => context.read<AuthBloc>().add(
                                  AuthRegisterRequested(
                                    username: _usernameForApi,
                                    displayName: _displayNameCtrl.text.trim(),
                                    email: _emailCtrl.text.trim(),
                                    password: _passCtrl.text,
                                  ),
                                )
                            : null,
                      ),
                    ),

                    const SizedBox(height: 24),

                    // Already have account
                    Center(
                      child: Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Text(
                            '${l10n.authAlreadyHaveAccount} ',
                            style: TextStyle(
                              fontFamily: 'Manrope',
                              fontSize: 14,
                              color: palette.textSecondary,
                            ),
                          ),
                          GestureDetector(
                            onTap: () => context.pushReplacementNamed(AppRoutes.login),
                            child: Text(
                              l10n.authLogin,
                              style: TextStyle(
                                fontFamily: 'Manrope',
                                fontSize: 14,
                                fontWeight: FontWeight.w600,
                                color: theme.colorScheme.tertiary,
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),

                    const SizedBox(height: 24),
                  ],
                ),
              ),
            ),
          );
        },
      ),
    );
  }
}
