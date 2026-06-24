import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:image_picker/image_picker.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import '../../../../app/config/theme/app_theme.dart';
import '../../../../shared/presentation/feedback/app_feedback.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/widgets/suda_avatar.dart';
import '../../../../shared/presentation/widgets/suda_text_field.dart';
import '../bloc/profile_bloc.dart';

class EditProfilePage extends StatelessWidget {
  const EditProfilePage({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) => sl<ProfileBloc>()..add(const ProfileLoadRequested()),
      child: const _EditProfileView(),
    );
  }
}

class _EditProfileView extends StatefulWidget {
  const _EditProfileView();

  @override
  State<_EditProfileView> createState() => _EditProfileViewState();
}

class _EditProfileViewState extends State<_EditProfileView> {
  final _displayCtrl = TextEditingController();
  final _firstCtrl = TextEditingController();
  final _lastCtrl = TextEditingController();
  final _bioCtrl = TextEditingController();
  final _picker = ImagePicker();
  bool _seeded = false;

  @override
  void dispose() {
    _displayCtrl.dispose();
    _firstCtrl.dispose();
    _lastCtrl.dispose();
    _bioCtrl.dispose();
    super.dispose();
  }

  void _save(BuildContext context) {
    context.read<ProfileBloc>().add(ProfileSaveRequested(
          displayName: _displayCtrl.text.trim(),
          firstName: _firstCtrl.text.trim(),
          lastName: _lastCtrl.text.trim(),
          bio: _bioCtrl.text.trim(),
        ));
  }

  Future<void> _pickAvatar(BuildContext context) async {
    final file = await _picker.pickImage(source: ImageSource.gallery, imageQuality: 88, maxWidth: 1200);
    if (!context.mounted || file == null) return;
    context.read<ProfileBloc>().add(ProfileAvatarUploadRequested(filePath: file.path, fileName: file.name));
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    return BlocConsumer<ProfileBloc, ProfileState>(
      listenWhen: (p, c) =>
          p.profile != c.profile || p.actionStatus != c.actionStatus,
      listener: (context, state) {
        // Seed controllers once the profile is loaded.
        if (!_seeded && state.profile != null) {
          final p = state.profile!;
          _displayCtrl.text = p.displayName;
          _firstCtrl.text = p.firstName;
          _lastCtrl.text = p.lastName;
          _bioCtrl.text = p.bio;
          _seeded = true;
        }
        if (state.actionStatus == ProfileActionStatus.saveSuccess) {
          AppFeedback.showSuccess(l10n.profileUpdated);
          context.pop();
        } else if (state.actionStatus == ProfileActionStatus.failure) {
          final msg = state.actionMessage;
          if (msg != null && msg.isNotEmpty) AppFeedback.showError(msg);
        }
      },
      builder: (context, state) {
        final profile = state.profile;
        final username = profile?.username ?? '';
        final handle = username.startsWith('@') ? username : '@$username';

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
              l10n.editProfileTitle,
              style: TextStyle(fontFamily: 'Manrope', fontWeight: FontWeight.w700, color: theme.colorScheme.onSurface),
            ),
            actions: [
              TextButton(
                onPressed: state.isSaving ? null : () => _save(context),
                child: Text(
                  l10n.editProfileSave,
                  style: TextStyle(
                    fontFamily: 'Manrope',
                    fontWeight: FontWeight.w700,
                    color: state.isSaving ? palette.textSecondary : theme.colorScheme.tertiary,
                  ),
                ),
              ),
            ],
          ),
          body: state.status == ProfileViewStatus.loading && profile == null
              ? const Center(child: CircularProgressIndicator())
              : ListView(
                  padding: const EdgeInsets.fromLTRB(20, 8, 20, 32),
                  children: [
                    // Avatar with camera
                    Center(
                      child: Column(
                        children: [
                          Stack(
                            children: [
                              SudaAvatar(
                                mediaUrl: _cacheBusted(profile?.avatarUrl, state.avatarVersion),
                                initials: profile?.displayName ?? '?',
                                size: 96,
                                ring: true,
                              ),
                              Positioned(
                                bottom: -2,
                                right: -2,
                                child: GestureDetector(
                                  onTap: state.isUploadingAvatar ? null : () => _pickAvatar(context),
                                  child: Container(
                                    width: 30,
                                    height: 30,
                                    decoration: BoxDecoration(
                                      color: theme.colorScheme.tertiary,
                                      shape: BoxShape.circle,
                                      border: Border.all(color: theme.scaffoldBackgroundColor, width: 2.5),
                                    ),
                                    child: state.isUploadingAvatar
                                        ? Padding(
                                            padding: const EdgeInsets.all(6),
                                            child: CircularProgressIndicator(
                                              strokeWidth: 2,
                                              color: theme.scaffoldBackgroundColor,
                                            ),
                                          )
                                        : Icon(Icons.camera_alt_rounded, size: 14, color: theme.scaffoldBackgroundColor),
                                  ),
                                ),
                              ),
                            ],
                          ),
                          const SizedBox(height: 8),
                          Text(
                            handle,
                            style: TextStyle(fontFamily: 'Manrope', fontSize: 13, color: palette.textSecondary),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: 24),

                    SudaTextField(label: l10n.editProfileDisplayName, controller: _displayCtrl, textInputAction: TextInputAction.next),
                    const SizedBox(height: 18),
                    SudaTextField(label: l10n.editProfileFirstName, hint: l10n.editProfileOptional, controller: _firstCtrl, textInputAction: TextInputAction.next),
                    const SizedBox(height: 18),
                    SudaTextField(label: l10n.editProfileLastName, hint: l10n.editProfileOptional, controller: _lastCtrl, textInputAction: TextInputAction.next),
                    const SizedBox(height: 18),
                    SudaTextField(label: l10n.editProfileBio, hint: l10n.editProfileBioHint, controller: _bioCtrl, maxLines: 3),
                    const SizedBox(height: 18),
                    // Username (read-only)
                    SudaTextField(
                      label: l10n.editProfileUsername,
                      controller: TextEditingController(text: handle),
                      enabled: false,
                      hint2: l10n.editProfileUsernameHint,
                    ),
                  ],
                ),
        );
      },
    );
  }

  String? _cacheBusted(String? url, int version) {
    // Presigned S3 URLs sign their query string — appending `&v=N` breaks the
    // signature (403). A fresh presigned URL is returned per upload/reload, so
    // the string changes on its own and busts the cache.
    if (url == null || url.trim().isEmpty) return null;
    return url.trim();
  }
}
