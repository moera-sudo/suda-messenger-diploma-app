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

class ForgotPasswordPage extends StatefulWidget {
  final String initialEmail;
  const ForgotPasswordPage({super.key, this.initialEmail = ''});

  @override
  State<ForgotPasswordPage> createState() => _ForgotPasswordPageState();
}

class _ForgotPasswordPageState extends State<ForgotPasswordPage> {
  late final TextEditingController _emailCtrl;

  @override
  void initState() {
    super.initState();
    _emailCtrl = TextEditingController(text: widget.initialEmail);
    _emailCtrl.addListener(() => setState(() {}));
  }

  @override
  void dispose() {
    _emailCtrl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);

    return BlocProvider(
      create: (_) => sl<AuthBloc>(),
      child: BlocConsumer<AuthBloc, AuthState>(
        listener: (context, state) {
          if (state.status == AuthStatus.forgotCodeSent) {
            AppFeedback.showSuccess(l10n.authVerificationCodeSent);
            context.pushNamed(AppRoutes.resetPassword, extra: _emailCtrl.text.trim().toLowerCase());
          }
          if (state.status == AuthStatus.failure) {
            AppFeedback.showError(
              AppFeedback.failureMessage(state.failure, fallback: l10n.errorGeneric),
            );
          }
        },
        builder: (context, state) {
          final isLoading = state.status == AuthStatus.loading;
          final canSubmit = _emailCtrl.text.contains('@');

          return Scaffold(
            backgroundColor: theme.scaffoldBackgroundColor,
            body: SafeArea(
              child: SingleChildScrollView(
                padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 8),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    IconButton(
                      icon: Icon(Icons.arrow_back_rounded, color: theme.colorScheme.onSurface),
                      onPressed: () => context.pop(),
                      padding: EdgeInsets.zero,
                    ),

                    const SizedBox(height: 32),

                    // Mail icon tile
                    Container(
                      width: 72,
                      height: 72,
                      decoration: BoxDecoration(
                        color: theme.colorScheme.surface,
                        borderRadius: BorderRadius.circular(22),
                        border: Border.all(color: palette.divider),
                        boxShadow: [
                          BoxShadow(color: palette.glow, blurRadius: 32, spreadRadius: -8),
                        ],
                      ),
                      child: Icon(
                        Icons.mail_outline_rounded,
                        size: 28,
                        color: theme.colorScheme.tertiary,
                      ),
                    ),

                    const SizedBox(height: 24),

                    Text(
                      l10n.authForgotPasswordTitle,
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
                      l10n.authForgotPasswordSubtitle,
                      style: TextStyle(
                        fontFamily: 'Manrope',
                        fontSize: 14,
                        color: palette.textSecondary,
                        height: 1.5,
                      ),
                    ),

                    const SizedBox(height: 28),

                    SudaTextField(
                      label: l10n.authEmail,
                      hint: 'you@example.com',
                      controller: _emailCtrl,
                      leftIcon: Icons.mail_outline_rounded,
                      keyboardType: TextInputType.emailAddress,
                      textInputAction: TextInputAction.done,
                    ),

                    const SizedBox(height: 24),

                    SudaButton(
                      label: l10n.authSendVerificationCode,
                      variant: SudaButtonVariant.primary,
                      size: SudaButtonSize.lg,
                      fullWidth: true,
                      loading: isLoading,
                      onPressed: canSubmit && !isLoading
                          ? () {
                              final email = _emailCtrl.text.trim();
                              if (email.isEmpty) {
                                AppFeedback.showError(l10n.authEmailRequired);
                                return;
                              }
                              context.read<AuthBloc>().add(AuthForgotPasswordRequested(email));
                            }
                          : null,
                    ),

                    const SizedBox(height: 16),

                    Center(
                      child: TextButton(
                        onPressed: () => context.pop(),
                        style: TextButton.styleFrom(
                          foregroundColor: theme.colorScheme.tertiary,
                          padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 8),
                        ),
                        child: Text(
                          l10n.authBackToLogin,
                          style: const TextStyle(
                            fontFamily: 'Manrope',
                            fontSize: 14,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                      ),
                    ),
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
