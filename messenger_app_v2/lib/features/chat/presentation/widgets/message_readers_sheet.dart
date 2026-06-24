import 'package:flutter/material.dart';

import 'package:messenger_app_v2/app/DI/get_it.dart';

import '../../../../shared/presentation/l10n_extensions.dart';
import '../../data/models/chat_models.dart';
import '../../domain/repositories/chat_repository.dart';

/// Bottom sheet that shows who has read a given message.
/// Only available for GROUP chats (API returns 403 for DIRECT/SAVED).
class MessageReadersSheet extends StatefulWidget {
  final int messageId;

  const MessageReadersSheet({super.key, required this.messageId});

  static Future<void> show(BuildContext context, int messageId) {
    return showModalBottomSheet(
      context: context,
      backgroundColor: Theme.of(context).colorScheme.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
      ),
      builder: (_) => MessageReadersSheet(messageId: messageId),
    );
  }

  @override
  State<MessageReadersSheet> createState() => _MessageReadersSheetState();
}

class _MessageReadersSheetState extends State<MessageReadersSheet> {
  late final Future<List<MessageReader>> _future;

  @override
  void initState() {
    super.initState();
    _future = sl<ChatRepository>()
        .getMessageReaders(widget.messageId)
        .then((result) => result.fold((_) => <MessageReader>[], (r) => r));
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return SafeArea(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          // Handle
          Container(
            margin: const EdgeInsets.only(top: 8, bottom: 4),
            width: 40,
            height: 4,
            decoration: BoxDecoration(
              color: theme.colorScheme.onSurface.withValues(alpha:0.2),
              borderRadius: BorderRadius.circular(2),
            ),
          ),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            child: Text(
              context.l10n.messageReadersTitle,
              style: TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.bold,
                color: theme.colorScheme.onSurface,
              ),
            ),
          ),
          const Divider(height: 1),
          FutureBuilder<List<MessageReader>>(
            future: _future,
            builder: (context, snapshot) {
              if (snapshot.connectionState == ConnectionState.waiting) {
                return const Padding(
                  padding: EdgeInsets.all(24),
                  child: Center(child: CircularProgressIndicator()),
                );
              }

              final readers = snapshot.data ?? [];
              if (readers.isEmpty) {
                return Padding(
                  padding: const EdgeInsets.all(24),
                  child: Center(
                    child: Text(
                      context.l10n.messageReadersEmpty,
                      style: TextStyle(color: theme.colorScheme.onSurface.withValues(alpha:0.5)),
                    ),
                  ),
                );
              }

              return ListView.builder(
                shrinkWrap: true,
                physics: const NeverScrollableScrollPhysics(),
                itemCount: readers.length,
                itemBuilder: (context, index) {
                  final reader = readers[index];
                  final readTime = _formatReadAt(reader.readAt);
                  return ListTile(
                    leading: CircleAvatar(
                      backgroundColor: theme.colorScheme.surfaceContainerHighest,
                      child: Text(
                        reader.displayName.isNotEmpty ? reader.displayName[0].toUpperCase() : '?',
                        style: TextStyle(color: theme.colorScheme.onSurface),
                      ),
                    ),
                    title: Text(reader.displayName, style: TextStyle(color: theme.colorScheme.onSurface)),
                    subtitle: Text(reader.username, style: TextStyle(color: theme.colorScheme.onSurface.withValues(alpha:0.5))),
                    trailing: Text(
                      readTime,
                      style: TextStyle(fontSize: 11, color: theme.colorScheme.onSurface.withValues(alpha:0.4)),
                    ),
                  );
                },
              );
            },
          ),
          const SizedBox(height: 8),
        ],
      ),
    );
  }

  String _formatReadAt(String iso) {
    try {
      final d = DateTime.parse(iso).toLocal();
      return '${d.hour.toString().padLeft(2, '0')}:${d.minute.toString().padLeft(2, '0')}';
    } catch (_) {
      return '';
    }
  }
}
