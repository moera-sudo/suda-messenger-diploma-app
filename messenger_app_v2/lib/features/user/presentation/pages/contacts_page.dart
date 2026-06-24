import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_contacts/flutter_contacts.dart';
import 'package:go_router/go_router.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';
import '../../../../app/config/theme/app_theme.dart';
import '../../../../app/navigation/app_routes.dart';
import '../../../../features/chat/domain/repositories/chat_repository.dart';
import '../../../../shared/data/contacts/device_contacts_service.dart';
import '../../../../shared/presentation/feedback/app_feedback.dart';
import '../../../../shared/presentation/format_utils.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/placeholder/placeholder_page.dart';
import '../../../../shared/presentation/widgets/suda_avatar.dart';
import '../bloc/contacts_bloc.dart';

class ContactsPage extends StatelessWidget {
  const ContactsPage({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) => sl<ContactsBloc>()..add(LoadContacts()),
      child: const _ContactsView(),
    );
  }
}

class _ContactsView extends StatefulWidget {
  const _ContactsView();

  @override
  State<_ContactsView> createState() => _ContactsViewState();
}

class _ContactsViewState extends State<_ContactsView> {
  List<Contact>? _deviceContacts;
  bool _loadingDevice = false;
  bool _permissionDenied = false;

  @override
  void initState() {
    super.initState();
    _loadDeviceContacts();
  }

  Future<void> _loadDeviceContacts() async {
    setState(() => _loadingDevice = true);
    // Cached per session — avoids re-reading the address book on every tab visit.
    final result = await sl<DeviceContactsService>().load();
    if (!mounted) return;
    setState(() {
      _loadingDevice = false;
      _permissionDenied = !result.granted;
      _deviceContacts = result.contacts;
    });
  }

  void _openChatWith(BuildContext context, String userId, String displayName) {
    sl<ChatRepository>()
        .createChat(type: 'DIRECT', targetId: userId)
        .then((result) {
      result.fold(
        (_) {},
        (chat) {
          if (context.mounted) {
            context.pushNamed(
              AppRoutes.chatDetail,
              pathParameters: {'id': chat.id},
              extra: {'name': displayName, 'interlocutorId': userId, 'chatType': 'DIRECT'},
            );
          }
        },
      );
    });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      appBar: AppBar(
        backgroundColor: theme.scaffoldBackgroundColor,
        elevation: 0,
        scrolledUnderElevation: 0,
        automaticallyImplyLeading: false,
        title: Text(
          l10n.contactsTitle,
          style: TextStyle(
            fontFamily: 'Manrope',
            fontSize: 18,
            fontWeight: FontWeight.w700,
            color: theme.colorScheme.onSurface,
          ),
        ),
        bottom: PreferredSize(
          preferredSize: const Size.fromHeight(1),
          child: Divider(height: 1, color: palette.divider),
        ),
      ),
      body: BlocBuilder<ContactsBloc, ContactsState>(
        buildWhen: (p, c) => p.status != c.status || p.contacts != c.contacts,
        builder: (context, state) {
          return ListView(
            children: [
              // ── In Suda ───────────────────────────────────────
              _SectionHeader(label: l10n.contactsInSuda),
              if (state.status == ContactsStatus.loading ||
                  state.status == ContactsStatus.initial)
                const Padding(
                  padding: EdgeInsets.all(24),
                  child: Center(child: CircularProgressIndicator()),
                )
              else if (state.status == ContactsStatus.failure)
                Padding(
                  padding: const EdgeInsets.all(16),
                  child: Text(state.error ?? l10n.errorGeneric,
                      style: TextStyle(color: palette.danger)),
                )
              else if (state.contacts.isEmpty)
                const PlaceholderPage(
                    type: PlaceholderType.noContent, scaffold: false)
              else
                ...state.contacts.map((contact) {
                  final displayName = contact.displayLabel;
                  return _ServerContactRow(
                    displayName: displayName,
                    username: contact.profile?.username ?? '',
                    avatarMediaId: contact.profile?.avatarMediaId,
                    isOnline: contact.profile?.isOnline ?? false,
                    onTap: () => _openChatWith(context, contact.contactId, displayName),
                    onRemove: () => context.read<ContactsBloc>().add(RemoveContact(contact.contactId)),
                  );
                }),

              const SizedBox(height: 16),

              // ── From phone ────────────────────────────────────
              _SectionHeader(label: l10n.contactsOnPhone),
              if (_loadingDevice)
                const Padding(
                  padding: EdgeInsets.all(24),
                  child: Center(child: CircularProgressIndicator()),
                )
              else if (_permissionDenied)
                Padding(
                  padding: const EdgeInsets.all(16),
                  child: Text(l10n.contactsPermissionDenied,
                      style: TextStyle(color: palette.textSecondary)),
                )
              else if (_deviceContacts == null || _deviceContacts!.isEmpty)
                const PlaceholderPage(
                    type: PlaceholderType.noContent, scaffold: false)
              else
                ..._deviceContacts!.map((c) {
                  final phones = c.phones.map((p) => p.number).join(', ');
                  return _DeviceContactRow(
                    displayName: c.displayName,
                    phone: phones,
                  );
                }),

              const SizedBox(height: 32),
            ],
          );
        },
      ),
    );
  }
}

// ── Section header ─────────────────────────────────────────────

class _SectionHeader extends StatelessWidget {
  final String label;
  const _SectionHeader({required this.label});

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    return Padding(
      padding: const EdgeInsets.fromLTRB(14, 14, 14, 6),
      child: Text(
        label.toUpperCase(),
        style: TextStyle(
          fontFamily: 'Manrope',
          fontSize: 11,
          fontWeight: FontWeight.w700,
          letterSpacing: 0.1 * 11,
          color: palette.textSecondary,
        ),
      ),
    );
  }
}

// ── Server contact row ─────────────────────────────────────────

class _ServerContactRow extends StatelessWidget {
  final String displayName;
  final String username;
  final String? avatarMediaId;
  final bool isOnline;
  final VoidCallback onTap;
  final VoidCallback onRemove;

  const _ServerContactRow({
    required this.displayName,
    required this.username,
    this.avatarMediaId,
    this.isOnline = false,
    required this.onTap,
    required this.onRemove,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    final initials = displayName.isNotEmpty ? displayName[0].toUpperCase() : '?';

    return InkWell(
      onTap: onTap,
      onLongPress: () => _showRemoveSheet(context, l10n, palette, theme),
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
        decoration: BoxDecoration(
          border: Border(bottom: BorderSide(color: palette.divider, width: 0.5)),
        ),
        child: Row(children: [
          SudaAvatar(
            mediaUrl: avatarMediaId,
            initials: initials,
            size: 48,
            showOnline: isOnline,
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(displayName,
                    style: TextStyle(
                        fontFamily: 'Manrope',
                        fontSize: 15,
                        fontWeight: FontWeight.w600,
                        color: theme.colorScheme.onSurface)),
                if (username.isNotEmpty)
                  Text(formatHandle(username),
                      style: TextStyle(
                          fontFamily: 'Manrope',
                          fontSize: 13,
                          color: palette.textSecondary)),
              ],
            ),
          ),
          IconButton(
            icon: Icon(Icons.chat_rounded, color: theme.colorScheme.tertiary),
            onPressed: onTap,
            tooltip: l10n.userProfileMessage,
          ),
        ]),
      ),
    );
  }

  void _showRemoveSheet(
      BuildContext context, dynamic l10n, AppPaletteTheme palette, ThemeData theme) {
    showModalBottomSheet(
      context: context,
      backgroundColor: theme.colorScheme.surface,
      shape: const RoundedRectangleBorder(
          borderRadius: BorderRadius.vertical(top: Radius.circular(24))),
      builder: (ctx) => SafeArea(
        child: Column(mainAxisSize: MainAxisSize.min, children: [
          Container(
            margin: const EdgeInsets.only(top: 8, bottom: 8),
            width: 38,
            height: 4,
            decoration: BoxDecoration(
              color: palette.textSecondary.withValues(alpha: 0.4),
              borderRadius: BorderRadius.circular(2),
            ),
          ),
          ListTile(
            leading: Icon(Icons.person_remove_outlined, color: palette.danger),
            title: Text(l10n.contactsRemove,
                style: TextStyle(
                    color: palette.danger,
                    fontFamily: 'Manrope',
                    fontWeight: FontWeight.w600)),
            onTap: () { Navigator.of(ctx).pop(); onRemove(); },
          ),
          const SizedBox(height: 8),
        ]),
      ),
    );
  }
}

// ── Device contact row (phone) ─────────────────────────────────

class _DeviceContactRow extends StatelessWidget {
  final String displayName;
  final String phone;

  const _DeviceContactRow({required this.displayName, required this.phone});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final palette = theme.extension<AppPaletteTheme>()!;
    final l10n = context.l10n;

    final initials = displayName.isNotEmpty ? displayName[0].toUpperCase() : '?';

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
      decoration: BoxDecoration(
        border: Border(bottom: BorderSide(color: palette.divider, width: 0.5)),
      ),
      child: Row(children: [
        SudaAvatar(initials: initials, size: 48),
        const SizedBox(width: 12),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(displayName,
                  style: TextStyle(
                      fontFamily: 'Manrope',
                      fontSize: 15,
                      fontWeight: FontWeight.w600,
                      color: theme.colorScheme.onSurface)),
              if (phone.isNotEmpty)
                Text(phone,
                    style: TextStyle(
                        fontFamily: 'Manrope',
                        fontSize: 13,
                        color: palette.textSecondary)),
            ],
          ),
        ),
        // Invite — coming soon (no backend phone→user resolution yet)
        Tooltip(
          message: l10n.comingSoon,
          child: IconButton(
            icon: Icon(Icons.person_add_outlined, color: palette.textSecondary),
            onPressed: () => AppFeedback.showSuccess(l10n.comingSoon),
          ),
        ),
      ]),
    );
  }
}
