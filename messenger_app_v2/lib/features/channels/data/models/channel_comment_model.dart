/// One comment under a channel post — GET/POST /channels/{id}/posts/{msg_id}/comments
/// and the CHANNEL_NEW_COMMENT WS payload.
class ChannelComment {
  final int id;
  final String channelId;
  final int postId;
  final String senderId;
  final String content;
  final int? replyToCommentId;
  final String? editedAt;
  final String createdAt;
  final String senderUsername;
  final String senderDisplayName;

  const ChannelComment({
    required this.id,
    required this.channelId,
    required this.postId,
    required this.senderId,
    required this.content,
    this.replyToCommentId,
    this.editedAt,
    required this.createdAt,
    required this.senderUsername,
    required this.senderDisplayName,
  });

  bool get isEdited => editedAt != null;
  bool get isReply => replyToCommentId != null;

  factory ChannelComment.fromJson(Map<String, dynamic> json) => ChannelComment(
        id: (json['id'] as num?)?.toInt() ?? 0,
        channelId: json['channel_id'] as String? ?? '',
        postId: (json['post_id'] as num?)?.toInt() ?? 0,
        senderId: json['sender_id'] as String? ?? '',
        content: json['content'] as String? ?? '',
        replyToCommentId: (json['reply_to_comment_id'] as num?)?.toInt(),
        editedAt: json['edited_at'] as String?,
        createdAt: json['created_at'] as String? ?? '',
        senderUsername: json['sender_username'] as String? ?? '',
        senderDisplayName: json['sender_display_name'] as String? ?? '',
      );

  ChannelComment copyWith({String? content, String? editedAt}) => ChannelComment(
        id: id,
        channelId: channelId,
        postId: postId,
        senderId: senderId,
        content: content ?? this.content,
        replyToCommentId: replyToCommentId,
        editedAt: editedAt ?? this.editedAt,
        createdAt: createdAt,
        senderUsername: senderUsername,
        senderDisplayName: senderDisplayName,
      );
}
