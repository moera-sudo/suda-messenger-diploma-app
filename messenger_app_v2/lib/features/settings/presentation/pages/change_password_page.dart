import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import '../../../../app/navigation/app_routes.dart';
import '../../../../features/profile/domain/repositories/profile_repository.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/widgets/suda_button.dart';
import '../../../../shared/presentation/widgets/suda_text_field.dart';

class ChangePasswordPage extends StatefulWidget {
  const ChangePasswordPage({super.key});

  @override
  State<ChangePasswordPage> createState() => _ChangePasswordPageState();
}

class _ChangePasswordPageState extends State<ChangePasswordPage> {
  final _oldCtrl = TextEditingController();
  final _newCtrl = TextEditingController();
  final _confirmCtrl = TextEditingController();
  bool _loading = false;
  String? _oldError;
  String? _newError;
  String? _confirmError;

  @override
  void dispose() {
    _oldCtrl.dispose();
    _newCtrl.dispose();
    _confirmCtrl.dispose();
    super.dispose();
  }

  bool _validate(BuildContext context) {
    final l10n = context.l10n;
    bool ok = true;

    setState(() {
      _oldError = null;
      _newError = null;
      _confirmError = null;
    });

    if (_newCtrl.text.length < 6) {
      setState(() => _newError = l10n.changePasswordMinLength);
      ok = false;
    }
    if (_confirmCtrl.text != _newCtrl.text) {
      setState(() => _confirmError = l10n.changePasswordMismatch);
      ok = false;
    }
    return ok;
  }

  Future<void> _submit(BuildContext context) async {
    if (!_validate(context)) return;
    setState(() => _loading = true);

    final result = await sl<ProfileRepository>().changePassword(
      oldPassword: _oldCtrl.text,
      newPassword: _newCtrl.text,
    );
    if (!context.mounted) return;
    setState(() => _loading = false);

    result.fold(
      (failure) {
        // 400 INVALID_OLD_PASSWORD → show inline error; other errors shown by ApiClient.
        if (failure.code == 'INVALID_OLD_PASSWORD') {
          setState(() => _oldError = context.l10n.changePasswordInvalid);
        }
      },
      (_) {
        showDialog(
          context: context,
          barrierDismissible: false,
          builder: (ctx) => AlertDialog(
            title: Text(
              context.l10n.changePasswordSuccess,
              style: const TextStyle(fontFamily: 'Manrope'),
            ),
            actions: [
              TextButton(
                onPressed: () {
                  Navigator.pop(ctx);
                  // Password change revokes all sessions → force re-login.
                  context.goNamed(AppRoutes.welcome);
                },
                child: const Text('OK'),
              ),
            ],
          ),
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
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
          l10n.settingsChangePassword,
          style: TextStyle(
              fontFamily: 'Manrope', fontWeight: FontWeight.w700, color: theme.colorScheme.onSurface),
        ),
      ),
      body: ListView(
        padding: const EdgeInsets.all(20),
        children: [
          SudaTextField(
            label: l10n.changePasswordOld,
            controller: _oldCtrl,
            obscureText: true,
            textInputAction: TextInputAction.next,
            error: _oldError,
          ),
          const SizedBox(height: 18),
          SudaTextField(
            label: l10n.changePasswordNew,
            controller: _newCtrl,
            obscureText: true,
            textInputAction: TextInputAction.next,
            error: _newError,
          ),
          const SizedBox(height: 18),
          SudaTextField(
            label: l10n.changePasswordConfirm,
            controller: _confirmCtrl,
            obscureText: true,
            error: _confirmError,
          ),
          const SizedBox(height: 28),
          ListenableBuilder(
            listenable: Listenable.merge([_oldCtrl, _newCtrl, _confirmCtrl]),
            builder: (_, __) => SudaButton(
              label: l10n.settingsChangePassword,
              variant: SudaButtonVariant.primary,
              size: SudaButtonSize.lg,
              fullWidth: true,
              loading: _loading,
              onPressed: (_oldCtrl.text.isNotEmpty &&
                      _newCtrl.text.isNotEmpty &&
                      _confirmCtrl.text.isNotEmpty &&
                      !_loading)
                  ? () => _submit(context)
                  : null,
            ),
          ),
        ],
      ),
    );
  }
}
