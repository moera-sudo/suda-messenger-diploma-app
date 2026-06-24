import 'package:json_annotation/json_annotation.dart';

part 'search_result.g.dart';

/// Type constants for search result categories.
abstract class SearchResultType {
  static const user    = 'USER';
  static const chat    = 'CHAT';
  static const message = 'MESSAGE';
}

@JsonSerializable()
class SearchResult {
  final String type; // result category: USER, CHAT, MESSAGE
  final String id;
  final String title;
  final String description;
  @JsonKey(name: 'image_url', defaultValue: '')
  final String imageUrl;
  // Real chat type for CHAT results: DIRECT, GROUP, CHANNEL, SAVED.
  @JsonKey(name: 'chat_type')
  final String? chatType;
  // Only populated for MESSAGE type results
  @JsonKey(name: 'chat_id')
  final String? chatId;
  @JsonKey(name: 'message_id')
  final int? messageId;

  SearchResult({
    required this.type,
    required this.id,
    required this.title,
    required this.description,
    this.imageUrl = '',
    this.chatType,
    this.chatId,
    this.messageId,
  });

  factory SearchResult.fromJson(Map<String, dynamic> json) =>
      _$SearchResultFromJson(json);

  Map<String, dynamic> toJson() => _$SearchResultToJson(this);
}
