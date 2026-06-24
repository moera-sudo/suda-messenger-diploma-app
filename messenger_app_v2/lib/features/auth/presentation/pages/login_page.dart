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
import '../../../chat/presentation/bloc/socket_bloc.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final _emailCtrl = TextEditingController();
  final _passCtrl = TextEditingController();
  bool _obscurePass = true;

  @override
  void dispose() {
    _emailCtrl.dispose();
    _passCtrl.dispose();
    super.dispose();
  }

  bool get _canSubmit => _emailCtrl.text.length > 1 && _passCtrl.text.length >= 6;

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
            context.read<SocketBloc>().add(SocketConnect());
            context.goNamed(AppRoutes.chats);
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
                    IconButton(
                      icon: Icon(Icons.arrow_back_rounded, color: theme.colorScheme.onSurface),
                      onPressed: () => context.pop(),
                      padding: EdgeInsets.zero,
                    ),

                    const SizedBox(height: 20),

                    Text(
                      l10n.authWelcomeBack,
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
                      l10n.authExcited,
                      style: TextStyle(
                        fontFamily: 'Manrope',
                        fontSize: 14,
                        color: palette.textSecondary,
                      ),
                    ),

                    const SizedBox(height: 32),

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
                      listenable: Listenable.merge([_emailCtrl, _passCtrl]),
                      builder: (context, _) => SudaTextField(
                        label: l10n.authPassword,
                        hint: 'Your password',
                        controller: _passCtrl,
                        leftIcon: Icons.lock_outline_rounded,
                        obscureText: _obscurePass,
                        rightIcon: _obscurePass ? Icons.visibility_off_outlined : Icons.visibility_outlined,
                        onRightTap: () => setState(() => _obscurePass = !_obscurePass),
                        textInputAction: TextInputAction.done,
                      ),
                    ),

                    // Forgot password link
                    Align(
                      alignment: Alignment.centerLeft,
                      child: TextButton(
                        onPressed: () => context.pushNamed(
                          AppRoutes.forgotPassword,
                          extra: _emailCtrl.text.trim(),
                        ),
                        style: TextButton.styleFrom(
                          foregroundColor: theme.colorScheme.tertiary,
                          padding: const EdgeInsets.symmetric(vertical: 8, horizontal: 4),
                        ),
                        child: Text(
                          l10n.authForgotPassword,
                          style: const TextStyle(
                            fontFamily: 'Manrope',
                            fontSize: 14,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                      ),
                    ),

                    const SizedBox(height: 16),

                    ListenableBuilder(
                      listenable: Listenable.merge([_emailCtrl, _passCtrl]),
                      builder: (context, _) => SudaButton(
                        label: l10n.authLogin,
                        variant: SudaButtonVariant.primary,
                        size: SudaButtonSize.lg,
                        fullWidth: true,
                        loading: isLoading,
                        onPressed: _canSubmit && !isLoading
                            ? () => context.read<AuthBloc>().add(
                                  AuthLoginRequested(_emailCtrl.text.trim(), _passCtrl.text),
                                )
                            : null,
                      ),
                    ),

                    const SizedBox(height: 24),

                    Center(
                      child: Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Text(
                            '${l10n.authNeedAccount} ',
                            style: TextStyle(
                              fontFamily: 'Manrope',
                              fontSize: 14,
                              color: palette.textSecondary,
                            ),
                          ),
                          GestureDetector(
                            onTap: () => context.pushReplacementNamed(AppRoutes.register),
                            child: Text(
                              l10n.authRegister,
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
