import 'dart:async';

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
import '../bloc/auth_bloc.dart';

class VerifyPage extends StatefulWidget {
  final String email;
  const VerifyPage({super.key, required this.email});

  @override
  State<VerifyPage> createState() => _VerifyPageState();
}

class _VerifyPageState extends State<VerifyPage> {
  final List<TextEditingController> _controllers = List.generate(6, (_) => TextEditingController());
  final List<FocusNode> _focusNodes = List.generate(6, (_) => FocusNode());

  int _remainingSeconds = 299;
  Timer? _timer;

  @override
  void initState() {
    super.initState();
    _startTimer();
    for (var c in _controllers) {
      c.addListener(() => setState(() {}));
    }
  }

  void _startTimer() {
    _timer = Timer.periodic(const Duration(seconds: 1), (t) {
      if (!mounted) { t.cancel(); return; }
      setState(() {
        if (_remainingSeconds > 0) {
          _remainingSeconds--;
        } else {
          t.cancel();
        }
      });
    });
  }

  @override
  void dispose() {
    _timer?.cancel();
    for (final c in _controllers) { c.dispose(); }
    for (final f in _focusNodes) { f.dispose(); }
    super.dispose();
  }

  void _onDigitChanged(int index, String value) {
    final digit = value.replaceAll(RegExp(r'\D'), '');
    if (digit.length > 1) {
      // Handle paste: distribute digits across cells
      for (int i = 0; i < digit.length && i + index < 6; i++) {
        _controllers[index + i].text = digit[i];
      }
      final next = (index + digit.length).clamp(0, 5);
      _focusNodes[next].requestFocus();
      return;
    }
    if (digit.isNotEmpty) {
      if (index < 5) _focusNodes[index + 1].requestFocus();
    }
  }

  bool get _isFilled => _controllers.every((c) => c.text.isNotEmpty);

  String get _timerText {
    final mm = (_remainingSeconds ~/ 60).toString().padLeft(2, '0');
    final ss = (_remainingSeconds % 60).toString().padLeft(2, '0');
    return '$mm:$ss';
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
          if (state.status == AuthStatus.success) {
            context.goNamed(AppRoutes.welcome);
          } else if (state.status == AuthStatus.failure) {
            AppFeedback.showError(AppFeedback.failureMessage(state.failure));
          }
        },
        builder: (context, state) {
          final isLoading = state.status == AuthStatus.loading;

          return Scaffold(
            backgroundColor: theme.scaffoldBackgroundColor,
            body: SafeArea(
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 8),
                child: Column(
                  children: [
                    Align(
                      alignment: Alignment.topLeft,
                      child: IconButton(
                        icon: Icon(Icons.arrow_back_rounded, color: theme.colorScheme.onSurface),
                        onPressed: () => context.pop(),
                        padding: EdgeInsets.zero,
                      ),
                    ),

                    const Spacer(flex: 2),

                    // Shield icon tile
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
                        Icons.shield_outlined,
                        size: 32,
                        color: theme.colorScheme.tertiary,
                      ),
                    ),

                    const SizedBox(height: 24),

                    Text(
                      l10n.authVerifyTitle,
                      style: TextStyle(
                        fontFamily: 'SpaceGrotesk',
                        fontSize: 26,
                        fontWeight: FontWeight.w800,
                        letterSpacing: -0.03,
                        color: theme.colorScheme.onSurface,
                      ),
                      textAlign: TextAlign.center,
                    ),
                    const SizedBox(height: 10),
                    Text(
                      l10n.authEnterCodeSentTo(widget.email),
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontFamily: 'Manrope',
                        fontSize: 14,
                        color: palette.textSecondary,
                        height: 1.5,
                      ),
                    ),

                    const SizedBox(height: 32),

                    // OTP cells
                    Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: List.generate(6, (i) {
                        final filled = _controllers[i].text.isNotEmpty;
                        return Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 5),
                          child: SudaOtpCell(
                            controller: _controllers[i],
                            focusNode: _focusNodes[i],
                            filled: filled,
                            onChanged: (v) => _onDigitChanged(i, v),
                            onBackspace: () {
                              if (_controllers[i].text.isEmpty && i > 0) {
                                _controllers[i - 1].clear();
                                _focusNodes[i - 1].requestFocus();
                              }
                            },
                          ),
                        );
                      }),
                    ),

                    const SizedBox(height: 20),

                    // Timer
                    Text(
                      '${l10n.authCodeExpires} ',
                      style: TextStyle(fontFamily: 'Manrope', fontSize: 14, color: palette.textSecondary),
                    ),
                    Text(
                      _timerText,
                      style: TextStyle(
                        fontFamily: 'JetBrainsMono',
                        fontSize: 14,
                        fontWeight: FontWeight.w700,
                        color: theme.colorScheme.tertiary,
                      ),
                    ),

                    const SizedBox(height: 8),

                    // Resend
                    TextButton.icon(
                      onPressed: _remainingSeconds == 0
                          ? () {
                              context.read<AuthBloc>().add(
                                    AuthResendCodeRequested(widget.email),
                                  );
                              setState(() => _remainingSeconds = 299);
                              _startTimer();
                            }
                          : null,
                      icon: Icon(Icons.refresh_rounded, size: 14, color: theme.colorScheme.tertiary),
                      label: Text(
                        l10n.authResendCode,
                        style: TextStyle(
                          fontFamily: 'Manrope',
                          fontSize: 14,
                          fontWeight: FontWeight.w600,
                          color: theme.colorScheme.tertiary,
                        ),
                      ),
                      style: TextButton.styleFrom(padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4)),
                    ),

                    const Spacer(flex: 3),

                    // Verify button
                    SudaButton(
                      label: l10n.authVerifyButton,
                      variant: SudaButtonVariant.primary,
                      size: SudaButtonSize.lg,
                      fullWidth: true,
                      iconRight: Icons.arrow_forward_rounded,
                      loading: isLoading,
                      onPressed: _isFilled && !isLoading
                          ? () {
                              final code = _controllers.map((c) => c.text).join();
                              context.read<AuthBloc>().add(AuthVerifyRequested(widget.email, code));
                            }
                          : null,
                    ),

                    const SizedBox(height: 28),
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

