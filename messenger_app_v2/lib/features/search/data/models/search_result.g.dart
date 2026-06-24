// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'search_result.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

SearchResult _$SearchResultFromJson(Map<String, dynamic> json) => SearchResult(
  type: json['type'] as String,
  id: json['id'] as String,
  title: json['title'] as String,
  description: json['description'] as String,
  imageUrl: json['image_url'] as String? ?? '',
  chatType: json['chat_type'] as String?,
  chatId: json['chat_id'] as String?,
  messageId: (json['message_id'] as num?)?.toInt(),
);

Map<String, dynamic> _$SearchResultToJson(SearchResult instance) =>
    <String, dynamic>{
      'type': instance.type,
      'id': instance.id,
      'title': instance.title,
      'description': instance.description,
      'image_url': instance.imageUrl,
      'chat_type': instance.chatType,
      'chat_id': instance.chatId,
      'message_id': instance.messageId,
    };
