import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:image_picker/image_picker.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import 'package:messenger_app_v2/app/navigation/app_routes.dart';
import '../../../../app/config/theme/app_theme.dart';
import '../../../../features/friends/data/models/friend_models.dart';
import '../../../../features/friends/presentation/bloc/friends_cubit.dart';
import '../../../../features/media/domain/repositories/media_repository.dart';
import '../../../../shared/data/auth/current_user.dart';
import '../../../../shared/presentation/feedback/app_feedback.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/widgets/suda_avatar.dart';
import '../../../../shared/presentation/widgets/suda_button.dart';
import '../../data/models/chat_models.dart';
import '../../domain/repositories/chat_repository.dart';

/// Loaded chat info plus the current user's resolved role.
typedef _ChatInfoData = ({ChatInfo info, String role});

class ChatInfoPage extends StatefulWidget {
  final String chatId;
  final String chatType;
  final String myRole;

  const ChatInfoPage({
    super.key,
    required this.chatId,
    this.chatType = 'GROUP',
    this.myRole = 'MEMBER',
  });

  @override
  State<ChatInfoPage> createState() => _ChatInfoPageState();
}

class _ChatInfoPageState extends State<ChatInfoPage> {
  late Future<_ChatInfoData?> _dataFuture;

  @override
  void initState() {
    super.initState();
    _dataFuture = _loadData();
  }

  Future<_ChatInfoData?> _loadData() async {
    final result = await sl<ChatRepository>().getChatInfo(widget.chatId);
    final info = result.fold((_) => null, (i) => i);
    if (info == null) return null;

    // Resolve current user's role from members (route extra never carries it).
    String role = widget.myRole;
    final myId = await sl<CurrentUser>().id();
    if (myId != null) {
      for (final m in info.members) {
        if (m.userId == myId) {
          role = m.role;
          break;
        }
      }
    }
    return (info: info, role: role);
  }

  void _reload() => setState(() => _dataFuture = _loadData());

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
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
          l10n.chatInfo,
          style: TextStyle(fontFamily: 'Manrope', fontWeight: FontWeight.w700, color: theme.colorScheme.onSurface),
        ),
        actions: [
          FutureBuilder<_ChatInfoData?>(
            future: _dataFuture,
            builder: (context, snapshot) {
              final isAdmin = snapshot.data?.role == 'OWNER' ||
                  snapshot.data?.role == 'ADMIN';
              if (!isAdmin) return const SizedBox.shrink();
              return IconButton(
                icon: Icon(Icons.edit_rounded, color: theme.colorScheme.onSurface),
                tooltip: l10n.chatEditGroup,
                onPressed: () => _openEditSheet(context, snapshot.data!.info),
              );
            },
          ),
        ],
      ),
      body: FutureBuilder<_ChatInfoData?>(
        future: _dataFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          }
          final data = snapshot.data;
          if (data == null) {
            return Center(
              child: Text(l10n.errorGeneric, style: TextStyle(color: palette.danger)),
            );
          }
          final info = data.info;
          final isOwner = data.role == 'OWNER';
          final isAdmin = data.role == 'OWNER' || data.role == 'ADMIN';

          return ListView(
            padding: const EdgeInsets.only(bottom: 32),
            children: [
              // ── Hero ──────────────────────────────────────────
              _HeroSection(info: info, palette: palette, theme: theme),

              const SizedBox(height: 16),

              // ── Shared media ──────────────────────────────────
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 14),
                child: Material(
                  color: theme.colorScheme.surface,
                  borderRadius: BorderRadius.circular(16),
                  clipBehavior: Clip.hardEdge,
                  child: ListTile(
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(16),
                      side: BorderSide(color: palette.divider),
                    ),
                    leading: Icon(Icons.image_outlined, color: palette.textAccent),
                    title: Text(
                      l10n.userProfileSharedMedia,
                      style: TextStyle(
                        fontFamily: 'Manrope',
                        fontSize: 14,
                        fontWeight: FontWeight.w500,
                        color: theme.colorScheme.onSurface,
                      ),
                    ),
                    trailing: Icon(Icons.chevron_right_rounded,
                        size: 18, color: palette.textSecondary),
                    onTap: () => context.pushNamed(
                      AppRoutes.sharedMedia,
                      pathParameters: {'id': widget.chatId},
                      extra: {'chatName': info.name ?? ''},
                    ),
                  ),
                ),
              ),

              const SizedBox(height: 16),

              // ── Members section ───────────────────────────────
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 14),
                child: Text(
                  l10n.chatMembers.toUpperCase(),
                  style: TextStyle(
                    fontFamily: 'Manrope',
                    fontSize: 11,
                    fontWeight: FontWeight.w700,
                    color: palette.textSecondary,
                    letterSpacing: 0.10,
                  ),
                ),
              ),
              const SizedBox(height: 8),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 14),
                child: Container(
                  decoration: BoxDecoration(
                    color: theme.colorScheme.surface,
                    borderRadius: BorderRadius.circular(16),
                    border: Border.all(color: palette.divider),
                  ),
                  clipBehavior: Clip.hardEdge,
                  child: Column(
                    children: [
                      // Add member — only admins
                      if (isAdmin)
                        _AddMemberTile(
                          palette: palette,
                          theme: theme,
                          l10n: l10n,
                          onTap: () => _openAddMembersSheet(context, info.members),
                        ),

                      ...info.members.asMap().entries.map((entry) {
                        final i = entry.key;
                        final member = entry.value;
                        final isLast = i == info.members.length - 1 && !isAdmin;
                        return Column(
                          children: [
                            _MemberTile(
                              member: member,
                              palette: palette,
                              theme: theme,
                              canManage: isAdmin && member.role != 'OWNER',
                              onMakeAdmin: () => _makeAdmin(context, member),
                              onRemove: () => _removeMember(context, member),
                            ),
                            if (!isLast)
                              Divider(height: 1, indent: 62, color: palette.divider),
                          ],
                        );
                      }),
                    ],
                  ),
                ),
              ),

              const SizedBox(height: 24),

              // ── Danger zone ───────────────────────────────────
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 14),
                child: Column(
                  children: [
                    if (!isOwner)
                      SudaButton(
                        label: l10n.chatLeave,
                        variant: SudaButtonVariant.dangerOutline,
                        size: SudaButtonSize.lg,
                        fullWidth: true,
                        icon: Icons.exit_to_app_rounded,
                        onPressed: () => _confirmLeave(context),
                      ),
                    if (isOwner) ...[
                      SudaButton(
                        label: l10n.chatDelete,
                        variant: SudaButtonVariant.dangerOutline,
                        size: SudaButtonSize.lg,
                        fullWidth: true,
                        icon: Icons.delete_outline_rounded,
                        onPressed: () => _confirmDelete(context),
                      ),
                    ],
                  ],
                ),
              ),
            ],
          );
        },
      ),
    );
  }

  Future<void> _makeAdmin(BuildContext context, ChatMember member) async {
    final result = await sl<ChatRepository>().changeMemberRole(widget.chatId, member.userId, 'ADMIN');
    if (!context.mounted) return;
    result.fold(
      (f) => AppFeedback.showError(f.message),
      (_) => _reload(),
    );
  }

  Future<void> _removeMember(BuildContext context, ChatMember member) async {
    final result = await sl<ChatRepository>().removeMember(widget.chatId, member.userId);
    if (!context.mounted) return;
    result.fold(
      (f) => AppFeedback.showError(f.message),
      (_) => _reload(),
    );
  }

  void _confirmLeave(BuildContext context) {
    final l10n = context.l10n;
    _showConfirm(
      context,
      message: l10n.chatLeaveConfirm,
      onConfirm: () async {
        final result = await sl<ChatRepository>().leaveChat(widget.chatId);
        if (!context.mounted) return;
        result.fold(
          (f) => AppFeedback.showError(f.message),
          (_) => context.go('/chats'),
        );
      },
    );
  }

  void _confirmDelete(BuildContext context) {
    final l10n = context.l10n;
    _showConfirm(
      context,
      message: l10n.chatDeleteConfirm,
      onConfirm: () async {
        final result = await sl<ChatRepository>().deleteChat(widget.chatId);
        if (!context.mounted) return;
        result.fold(
          (f) => AppFeedback.showError(f.message),
          (_) => context.go('/chats'),
        );
      },
    );
  }

  void _showConfirm(BuildContext context, {required String message, required VoidCallback onConfirm}) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);
    final l10n = context.l10n;
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        backgroundColor: theme.colorScheme.surface,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(18)),
        content: Text(message, style: TextStyle(fontFamily: 'Manrope', color: theme.colorScheme.onSurface)),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: Text(l10n.cancel, style: TextStyle(fontFamily: 'Manrope', color: palette.textSecondary)),
          ),
          TextButton(
            onPressed: () {
              Navigator.pop(ctx);
              onConfirm();
            },
            child: Text(l10n.confirm, style: TextStyle(fontFamily: 'Manrope', color: palette.danger, fontWeight: FontWeight.w700)),
          ),
        ],
      ),
    );
  }

  void _openEditSheet(BuildContext context, ChatInfo info) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Theme.of(context).colorScheme.surface,
      shape: const RoundedRectangleBorder(borderRadius: BorderRadius.vertical(top: Radius.circular(24))),
      builder: (_) => _EditGroupSheet(
        chatId: widget.chatId,
        currentName: info.name ?? '',
        currentDescription: info.description ?? '',
        currentAvatarMediaId: info.avatarMediaId,
        onSaved: _reload,
      ),
    );
  }

  void _openAddMembersSheet(BuildContext context, List<ChatMember> members) {
    final existingIds = members.map((m) => m.userId).toSet();
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Theme.of(context).colorScheme.surface,
      shape: const RoundedRectangleBorder(borderRadius: BorderRadius.vertical(top: Radius.circular(24))),
      builder: (_) => _AddMembersSheet(
        chatId: widget.chatId,
        existingMemberIds: existingIds,
        onAdded: _reload,
      ),
    );
  }
}

// ── Sub-widgets ───────────────────────────────────────────────

class _HeroSection extends StatelessWidget {
  final ChatInfo info;
  final AppPaletteTheme palette;
  final ThemeData theme;

  const _HeroSection({required this.info, required this.palette, required this.theme});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(vertical: 24),
      decoration: BoxDecoration(
        gradient: RadialGradient(
          center: Alignment.center,
          radius: 1.2,
          colors: [theme.colorScheme.primary.withValues(alpha: 0.3), Colors.transparent],
        ),
      ),
      child: Column(
        children: [
          SudaAvatar(mediaId: info.avatarMediaId, initials: info.name ?? '?', size: 80, ring: true),
          const SizedBox(height: 12),
          Text(
            info.name ?? '?',
            style: TextStyle(fontFamily: 'SpaceGrotesk', fontSize: 22, fontWeight: FontWeight.w800, color: theme.colorScheme.onSurface),
          ),
          const SizedBox(height: 4),
          Text(
            '${info.memberCount} ${context.l10n.chatMembers.toLowerCase()}',
            style: TextStyle(fontFamily: 'Manrope', fontSize: 14, color: palette.textSecondary),
          ),
          if (info.description?.isNotEmpty == true) ...[
            const SizedBox(height: 8),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 32),
              child: Text(
                info.description!,
                textAlign: TextAlign.center,
                style: TextStyle(fontFamily: 'Manrope', fontSize: 13, color: palette.textSecondary),
              ),
            ),
          ],
        ],
      ),
    );
  }
}

class _AddMemberTile extends StatelessWidget {
  final AppPaletteTheme palette;
  final ThemeData theme;
  final dynamic l10n;
  final VoidCallback onTap;

  const _AddMemberTile({
    required this.palette,
    required this.theme,
    required this.l10n,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        ListTile(
          leading: Container(
            width: 40,
            height: 40,
            decoration: BoxDecoration(
              color: theme.colorScheme.tertiary.withValues(alpha: 0.12),
              shape: BoxShape.circle,
            ),
            child: Icon(Icons.person_add_rounded, color: theme.colorScheme.tertiary, size: 20),
          ),
          title: Text(
            l10n.chatAddMember,
            style: TextStyle(fontFamily: 'Manrope', fontWeight: FontWeight.w600, color: theme.colorScheme.tertiary),
          ),
          onTap: onTap,
        ),
        Divider(height: 1, indent: 62, color: palette.divider),
      ],
    );
  }
}

class _MemberTile extends StatelessWidget {
  final ChatMember member;
  final AppPaletteTheme palette;
  final ThemeData theme;
  final bool canManage;
  final VoidCallback onMakeAdmin;
  final VoidCallback onRemove;

  const _MemberTile({
    required this.member,
    required this.palette,
    required this.theme,
    required this.canManage,
    required this.onMakeAdmin,
    required this.onRemove,
  });

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final roleName = switch (member.role) {
      'OWNER' => l10n.roleOwner,
      'ADMIN' => l10n.roleAdmin,
      _ => l10n.roleMember,
    };
    final roleColor = switch (member.role) {
      'OWNER' => theme.colorScheme.tertiary,
      'ADMIN' => theme.colorScheme.primary,
      _ => palette.textSecondary,
    };

    return ListTile(
      leading: SudaAvatar(
        mediaId: member.avatarMediaId,
        userId: member.userId,
        initials: member.displayName,
        size: 40,
      ),
      title: Text(
        member.displayName,
        style: TextStyle(fontFamily: 'Manrope', fontWeight: FontWeight.w600, color: theme.colorScheme.onSurface),
      ),
      subtitle: Text(
        member.username,
        style: TextStyle(fontFamily: 'Manrope', fontSize: 13, color: palette.textSecondary),
      ),
      trailing: Container(
        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
        decoration: BoxDecoration(
          color: roleColor.withValues(alpha: 0.12),
          borderRadius: BorderRadius.circular(100),
        ),
        child: Text(roleName, style: TextStyle(fontFamily: 'Manrope', fontSize: 11, fontWeight: FontWeight.w700, color: roleColor)),
      ),
      onLongPress: canManage
          ? () => showModalBottomSheet(
                context: context,
                backgroundColor: theme.colorScheme.surface,
                shape: const RoundedRectangleBorder(borderRadius: BorderRadius.vertical(top: Radius.circular(24))),
                builder: (_) => SafeArea(
                  child: Column(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      if (member.role == 'MEMBER')
                        ListTile(
                          leading: Icon(Icons.admin_panel_settings_outlined, color: theme.colorScheme.primary),
                          title: Text(l10n.chatMakeAdmin, style: TextStyle(color: theme.colorScheme.onSurface)),
                          onTap: () { Navigator.pop(context); onMakeAdmin(); },
                        ),
                      ListTile(
                        leading: Icon(Icons.person_remove_outlined, color: palette.danger),
                        title: Text(l10n.chatRemoveMember, style: TextStyle(color: palette.danger)),
                        onTap: () { Navigator.pop(context); onRemove(); },
                      ),
                    ],
                  ),
                ),
              )
          : null,
    );
  }
}

// ── _EditGroupSheet ───────────────────────────────────────────

class _EditGroupSheet extends StatefulWidget {
  final String chatId;
  final String currentName;
  final String currentDescription;
  final String? currentAvatarMediaId;
  final VoidCallback onSaved;

  const _EditGroupSheet({
    required this.chatId,
    required this.currentName,
    required this.currentDescription,
    this.currentAvatarMediaId,
    required this.onSaved,
  });

  @override
  State<_EditGroupSheet> createState() => _EditGroupSheetState();
}

class _EditGroupSheetState extends State<_EditGroupSheet> {
  late final TextEditingController _nameCtrl;
  late final TextEditingController _descCtrl;
  String? _pendingAvatarPath;
  bool _saving = false;

  @override
  void initState() {
    super.initState();
    _nameCtrl = TextEditingController(text: widget.currentName);
    _descCtrl = TextEditingController(text: widget.currentDescription);
  }

  @override
  void dispose() {
    _nameCtrl.dispose();
    _descCtrl.dispose();
    super.dispose();
  }

  Future<void> _pickAvatar() async {
    final picked = await ImagePicker().pickImage(source: ImageSource.gallery, imageQuality: 80);
    if (picked != null) setState(() => _pendingAvatarPath = picked.path);
  }

  Future<void> _save() async {
    if (_nameCtrl.text.trim().isEmpty) return;
    setState(() => _saving = true);

    String? newAvatarMediaId;
    if (_pendingAvatarPath != null) {
      final uploadResult = await sl<MediaRepository>().uploadFile(
        filePath: _pendingAvatarPath!,
        mediaType: 'CHAT_AVATAR',
      );
      final id = uploadResult.fold((_) => null, (id) => id);
      if (id == null) {
        if (mounted) setState(() => _saving = false);
        return;
      }
      newAvatarMediaId = id;
    }

    final result = await sl<ChatRepository>().updateChat(
      widget.chatId,
      name: _nameCtrl.text.trim(),
      description: _descCtrl.text.trim(),
      avatarMediaId: newAvatarMediaId,
    );

    if (!mounted) return;
    setState(() => _saving = false);
    result.fold(
      (f) => AppFeedback.showError(f.message),
      (_) {
        Navigator.pop(context);
        widget.onSaved();
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);
    final l10n = context.l10n;

    return Padding(
      padding: EdgeInsets.only(bottom: MediaQuery.of(context).viewInsets.bottom),
      child: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.fromLTRB(20, 20, 20, 8),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              // Handle bar
              Center(
                child: Container(
                  width: 40,
                  height: 4,
                  margin: const EdgeInsets.only(bottom: 16),
                  decoration: BoxDecoration(
                    color: palette.divider,
                    borderRadius: BorderRadius.circular(100),
                  ),
                ),
              ),

              Text(
                l10n.chatEditGroup,
                style: TextStyle(fontFamily: 'SpaceGrotesk', fontSize: 18, fontWeight: FontWeight.w700, color: theme.colorScheme.onSurface),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 20),

              // Avatar picker
              Center(
                child: GestureDetector(
                  onTap: _pickAvatar,
                  child: Stack(
                    children: [
                      _pendingAvatarPath != null
                          ? ClipOval(
                              child: Image.file(
                                File(_pendingAvatarPath!),
                                width: 88,
                                height: 88,
                                fit: BoxFit.cover,
                              ),
                            )
                          : SudaAvatar(
                              mediaId: widget.currentAvatarMediaId,
                              initials: widget.currentName.isNotEmpty ? widget.currentName : '?',
                              size: 88,
                            ),
                      Positioned(
                        right: 0,
                        bottom: 0,
                        child: Container(
                          width: 28,
                          height: 28,
                          decoration: BoxDecoration(
                            color: theme.colorScheme.primary,
                            shape: BoxShape.circle,
                            border: Border.all(color: theme.colorScheme.surface, width: 2),
                          ),
                          child: const Icon(Icons.camera_alt_rounded, size: 14, color: Colors.white),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
              const SizedBox(height: 20),

              // Name field
              TextField(
                controller: _nameCtrl,
                style: TextStyle(fontFamily: 'Manrope', color: theme.colorScheme.onSurface),
                decoration: InputDecoration(
                  labelText: l10n.groupName,
                  labelStyle: TextStyle(fontFamily: 'Manrope', color: palette.textSecondary),
                  border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                  enabledBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                    borderSide: BorderSide(color: palette.divider),
                  ),
                ),
              ),
              const SizedBox(height: 12),

              // Description field
              TextField(
                controller: _descCtrl,
                maxLines: 3,
                style: TextStyle(fontFamily: 'Manrope', color: theme.colorScheme.onSurface),
                decoration: InputDecoration(
                  labelText: l10n.chatGroupDescription,
                  hintText: l10n.chatGroupDescriptionHint,
                  hintStyle: TextStyle(fontFamily: 'Manrope', color: palette.textSecondary),
                  labelStyle: TextStyle(fontFamily: 'Manrope', color: palette.textSecondary),
                  border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                  enabledBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                    borderSide: BorderSide(color: palette.divider),
                  ),
                ),
              ),
              const SizedBox(height: 20),

              SudaButton(
                label: l10n.chatSaveChanges,
                variant: SudaButtonVariant.primary,
                size: SudaButtonSize.lg,
                fullWidth: true,
                loading: _saving,
                onPressed: _saving ? null : _save,
              ),
              const SizedBox(height: 8),
            ],
          ),
        ),
      ),
    );
  }
}

// ── _AddMembersSheet ──────────────────────────────────────────

class _AddMembersSheet extends StatefulWidget {
  final String chatId;
  final Set<String> existingMemberIds;
  final VoidCallback onAdded;

  const _AddMembersSheet({
    required this.chatId,
    required this.existingMemberIds,
    required this.onAdded,
  });

  @override
  State<_AddMembersSheet> createState() => _AddMembersSheetState();
}

class _AddMembersSheetState extends State<_AddMembersSheet> {
  final Set<String> _selected = {};
  bool _adding = false;

  Future<void> _confirm(BuildContext ctx, List<FriendModel> eligible) async {
    if (_selected.isEmpty) return;
    setState(() => _adding = true);

    final repo = sl<ChatRepository>();
    for (final userId in _selected) {
      final result = await repo.addMember(widget.chatId, userId);
      result.fold((f) => AppFeedback.showError(f.message), (_) {});
    }

    if (!ctx.mounted) return;
    setState(() => _adding = false);
    Navigator.pop(ctx);
    widget.onAdded();
  }

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);
    final l10n = context.l10n;

    return BlocProvider(
      create: (_) => sl<FriendsCubit>()..load(),
      child: Builder(
        builder: (ctx) => SafeArea(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              // Handle bar
              Container(
                width: 40,
                height: 4,
                margin: const EdgeInsets.symmetric(vertical: 12),
                decoration: BoxDecoration(color: palette.divider, borderRadius: BorderRadius.circular(100)),
              ),

              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20),
                child: Row(
                  children: [
                    Expanded(
                      child: Text(
                        l10n.chatAddMembers,
                        style: TextStyle(fontFamily: 'SpaceGrotesk', fontSize: 18, fontWeight: FontWeight.w700, color: theme.colorScheme.onSurface),
                      ),
                    ),
                    BlocBuilder<FriendsCubit, FriendsState>(
                      builder: (ctx2, state) {
                        final eligible = state.friends
                            .where((f) => !widget.existingMemberIds.contains(f.userId))
                            .toList();
                        return TextButton(
                          onPressed: _adding || _selected.isEmpty
                              ? null
                              : () => _confirm(ctx2, eligible),
                          child: _adding
                              ? const SizedBox(width: 16, height: 16, child: CircularProgressIndicator(strokeWidth: 2))
                              : Text(
                                  '${l10n.chatAddMembers} (${_selected.length})',
                                  style: TextStyle(fontFamily: 'Manrope', color: theme.colorScheme.primary),
                                ),
                        );
                      },
                    ),
                  ],
                ),
              ),

              const Divider(height: 1),

              BlocBuilder<FriendsCubit, FriendsState>(
                builder: (_, state) {
                  if (state.status == FriendsStatus.loading || state.status == FriendsStatus.initial) {
                    return const Padding(
                      padding: EdgeInsets.symmetric(vertical: 32),
                      child: Center(child: CircularProgressIndicator()),
                    );
                  }

                  final eligible = state.friends
                      .where((f) => !widget.existingMemberIds.contains(f.userId))
                      .toList();

                  if (eligible.isEmpty) {
                    return Padding(
                      padding: const EdgeInsets.symmetric(vertical: 32),
                      child: Center(
                        child: Text(l10n.friendsEmptyFriends, style: TextStyle(fontFamily: 'Manrope', color: palette.textSecondary)),
                      ),
                    );
                  }

                  return ConstrainedBox(
                    constraints: BoxConstraints(maxHeight: MediaQuery.of(context).size.height * 0.5),
                    child: ListView.builder(
                      shrinkWrap: true,
                      itemCount: eligible.length,
                      itemBuilder: (_, i) {
                        final friend = eligible[i];
                        final checked = _selected.contains(friend.userId);
                        return CheckboxListTile(
                          value: checked,
                          onChanged: (_) => setState(() {
                            if (checked) {
                              _selected.remove(friend.userId);
                            } else {
                              _selected.add(friend.userId);
                            }
                          }),
                          secondary: SudaAvatar(
                            mediaId: friend.avatarMediaId,
                            initials: friend.displayName,
                            size: 40,
                          ),
                          title: Text(
                            friend.displayName,
                            style: TextStyle(fontFamily: 'Manrope', fontWeight: FontWeight.w600, color: theme.colorScheme.onSurface),
                          ),
                          subtitle: Text(
                            '@${friend.username}',
                            style: TextStyle(fontFamily: 'Manrope', fontSize: 13, color: palette.textSecondary),
                          ),
                          activeColor: theme.colorScheme.primary,
                          controlAffinity: ListTileControlAffinity.trailing,
                        );
                      },
                    ),
                  );
                },
              ),
              const SizedBox(height: 8),
            ],
          ),
        ),
      ),
    );
  }
}
