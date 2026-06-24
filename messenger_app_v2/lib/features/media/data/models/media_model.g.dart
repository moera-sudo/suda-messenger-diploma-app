// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'media_model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

MediaUploadInit _$MediaUploadInitFromJson(Map<String, dynamic> json) =>
    MediaUploadInit(
      mediaId: json['media_id'] as String,
      uploadUrl: json['upload_url'] as String,
    );

Map<String, dynamic> _$MediaUploadInitToJson(MediaUploadInit instance) =>
    <String, dynamic>{
      'media_id': instance.mediaId,
      'upload_url': instance.uploadUrl,
    };

MediaResult _$MediaResultFromJson(Map<String, dynamic> json) => MediaResult(
  mediaId: json['media_id'] as String,
  url: json['url'] as String,
  type: json['type'] as String?,
  filename: json['filename'] as String?,
  size: (json['size'] as num?)?.toInt(),
);

Map<String, dynamic> _$MediaResultToJson(MediaResult instance) =>
    <String, dynamic>{
      'media_id': instance.mediaId,
      'url': instance.url,
      'type': instance.type,
      'filename': instance.filename,
      'size': instance.size,
    };
