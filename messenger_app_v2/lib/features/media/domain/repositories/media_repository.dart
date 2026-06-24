import 'package:dartz/dartz.dart';

import '../../../../shared/domain/models/app_failure.dart';

abstract class MediaRepository {
  /// 3-step upload: init → PUT S3 → confirm → returns media_id.
  /// [mediaType] maps to a media-service `kind` (IMAGE/VIDEO/FILE/VOICE/AVATAR/
  /// CHAT_AVATAR). [isPrivate] overrides the default (avatars public, rest private).
  Future<Either<AppFailure, String>> uploadFile({
    required String filePath,
    required String mediaType,
    String? filename,
    bool? isPrivate,
  });

  /// GET /api/v1/media/{mediaId}/url → presigned view URL
  Future<Either<AppFailure, String>> getMediaUrl(String mediaId);

  /// Resolves the media URL and downloads the bytes into the app cache,
  /// returning the local file path. Used for opening files and saving media
  /// to the device gallery. [filename] names the cached file.
  Future<Either<AppFailure, String>> downloadToCache(
    String mediaId,
    String filename,
  );
}
