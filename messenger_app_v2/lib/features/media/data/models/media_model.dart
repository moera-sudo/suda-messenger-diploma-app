import 'package:json_annotation/json_annotation.dart';

part 'media_model.g.dart';

/// Response from POST /media/upload/init
@JsonSerializable()
class MediaUploadInit {
  @JsonKey(name: 'media_id')
  final String mediaId;
  @JsonKey(name: 'upload_url')
  final String uploadUrl;

  const MediaUploadInit({required this.mediaId, required this.uploadUrl});

  factory MediaUploadInit.fromJson(Map<String, dynamic> json) =>
      _$MediaUploadInitFromJson(json);
}

/// Response from POST /media/upload/confirm and GET /media/{id}
@JsonSerializable()
class MediaResult {
  @JsonKey(name: 'media_id')
  final String mediaId;
  final String url;
  final String? type;
  final String? filename;
  final int? size;

  const MediaResult({
    required this.mediaId,
    required this.url,
    this.type,
    this.filename,
    this.size,
  });

  factory MediaResult.fromJson(Map<String, dynamic> json) =>
      _$MediaResultFromJson(json);
}
