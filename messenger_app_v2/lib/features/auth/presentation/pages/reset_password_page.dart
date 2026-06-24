import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import '../../../../app/config/theme/app_theme.dart';
import '../../../../app/navigation/app_routes.dart';
import '../../../../shared/presentation/feedback/app_feedback.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/widgets/suda_button.dart';
import '../../../../shared/presentation/widgets/suda_otp_cell.dart';
import '../../../../shared/presentation/widgets/suda_text_field.dart';
import '../bloc/auth_bloc.dart';

class ResetPasswordPage extends StatefulWidget {
  final String email;
  const ResetPasswordPage({super.key, required this.email});

  @override
  State<ResetPasswordPage> createState() => _ResetPasswordPageState();
}

class _ResetPasswordPageState extends State<ResetPasswordPage> {
  final List<TextEditingController> _otpControllers = List.generate(6, (_) => TextEditingController());
  final List<FocusNode> _otpFocusNodes = List.generate(6, (_) => FocusNode());
  final _newPassCtrl = TextEditingController();
  final _confirmPassCtrl = TextEditingController();
  bool _obscureNew = true;
  bool _obscureConfirm = true;

  @override
  void initState() {
    super.initState();
    for (final c in _otpControllers) { c.addListener(() => setState(() {})); }
    _newPassCtrl.addListener(() => setState(() {}));
    _confirmPassCtrl.addListener(() => setState(() {}));
  }

  @override
  void dispose() {
    for (final c in _otpControllers) { c.dispose(); }
    for (final f in _otpFocusNodes) { f.dispose(); }
    _newPassCtrl.dispose();
    _confirmPassCtrl.dispose();
    super.dispose();
  }

  bool get _canSubmit =>
      _otpControllers.every((c) => c.text.isNotEmpty) &&
      _newPassCtrl.text.length >= 6 &&
      _newPassCtrl.text == _confirmPassCtrl.text;

  void _onDigitChanged(int index, String value) {
    final digit = value.replaceAll(RegExp(r'\D'), '');
    if (digit.length > 1) {
      for (int i = 0; i < digit.length && i + index < 6; i++) {
        _otpControllers[index + i].text = digit[i];
      }
      final next = (index + digit.length).clamp(0, 5);
      _otpFocusNodes[next].requestFocus();
      return;
    }
    if (digit.isNotEmpty && index < 5) {
      _otpFocusNodes[index + 1].requestFocus();
    }
  }

  void _submit(BuildContext context) {
    final l10n = context.l10n;
    if (widget.email.trim().isEmpty) {
      AppFeedback.showError(l10n.authEmailMissing);
      return;
    }
    final code = _otpControllers.map((c) => c.text).join();
    if (code.length != 6) { AppFeedback.showError(l10n.authEnterFullCode); return; }
    if (_newPassCtrl.text.length < 6) { AppFeedback.showError(l10n.authPasswordMinLength); return; }
    if (_newPassCtrl.text != _confirmPassCtrl.text) { AppFeedback.showError(l10n.authPasswordsDoNotMatch); return; }

    context.read<AuthBloc>().add(AuthResetPasswordRequested(
      email: widget.email.trim().toLowerCase(),
      code: code,
      newPassword: _newPassCtrl.text,
    ));
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
          if (state.status == AuthStatus.passwordReset) {
            AppFeedback.showSuccess(l10n.authPasswordChangedSuccess);
            context.goNamed(AppRoutes.login);
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

                    const SizedBox(height: 32),

                    // Lock icon tile
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
                        Icons.lock_outline_rounded,
                        size: 28,
                        color: theme.colorScheme.tertiary,
                      ),
                    ),

                    const SizedBox(height: 24),

                    Text(
                      l10n.authResetPasswordTitle,
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
                      l10n.authEnterCodeSentTo(widget.email),
                      style: TextStyle(
                        fontFamily: 'Manrope',
                        fontSize: 14,
                        color: palette.textSecondary,
                        height: 1.5,
                      ),
                    ),

                    const SizedBox(height: 28),

                    // OTP cells (reuse widget from verify_page)
                    Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: List.generate(6, (i) {
                        final filled = _otpControllers[i].text.isNotEmpty;
                        return Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 5),
                          child: SudaOtpCell(
                            controller: _otpControllers[i],
                            focusNode: _otpFocusNodes[i],
                            filled: filled,
                            onChanged: (v) => _onDigitChanged(i, v),
                            onBackspace: () {
                              if (_otpControllers[i].text.isEmpty && i > 0) {
                                _otpControllers[i - 1].clear();
                                _otpFocusNodes[i - 1].requestFocus();
                              }
                            },
                          ),
                        );
                      }),
                    ),

                    const SizedBox(height: 24),

                    SudaTextField(
                      label: l10n.authNewPassword,
                      hint: '••••••••',
                      controller: _newPassCtrl,
                      leftIcon: Icons.lock_outline_rounded,
                      obscureText: _obscureNew,
                      rightIcon: _obscureNew ? Icons.visibility_off_outlined : Icons.visibility_outlined,
                      onRightTap: () => setState(() => _obscureNew = !_obscureNew),
                      textInputAction: TextInputAction.next,
                    ),

                    const SizedBox(height: 16),

                    SudaTextField(
                      label: l10n.authConfirmPassword,
                      hint: '••••••••',
                      controller: _confirmPassCtrl,
                      obscureText: _obscureConfirm,
                      rightIcon: _obscureConfirm ? Icons.visibility_off_outlined : Icons.visibility_outlined,
                      onRightTap: () => setState(() => _obscureConfirm = !_obscureConfirm),
                      textInputAction: TextInputAction.done,
                    ),

                    const SizedBox(height: 28),

                    SudaButton(
                      label: l10n.authResetPasswordTitle,
                      variant: SudaButtonVariant.primary,
                      size: SudaButtonSize.lg,
                      fullWidth: true,
                      loading: isLoading,
                      onPressed: _canSubmit && !isLoading ? () => _submit(context) : null,
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
