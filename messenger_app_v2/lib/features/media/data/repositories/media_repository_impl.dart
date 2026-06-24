import 'dart:io';

import 'package:dartz/dartz.dart';
import 'package:dio/dio.dart';
import 'package:injectable/injectable.dart';
import 'package:path_provider/path_provider.dart';

import '../../../../shared/data/api/api_client.dart';
import '../../../../shared/data/api/server_exception.dart';
import '../../../../shared/data/logger/app_logger.dart';
import '../../../../shared/domain/models/app_failure.dart';
import '../../domain/repositories/media_repository.dart';
import '../models/media_model.dart';

@LazySingleton(as: MediaRepository)
class MediaRepositoryImpl implements MediaRepository {
  final ApiClient _api;
  final AppLogger _logger;

  // Bare Dio for direct S3 PUT — no auth interceptors
  final Dio _s3Dio = Dio();

  MediaRepositoryImpl(this._api, this._logger);

  @override
  Future<Either<AppFailure, String>> uploadFile({
    required String filePath,
    required String mediaType,
    String? filename,
    bool? isPrivate,
  }) async {
    try {
      final file = File(filePath);
      final name = filename ?? Uri.file(filePath).pathSegments.last;
      final size = await file.length();
      final kind = _kindFor(mediaType);
      final contentType = _contentType(mediaType, name);
      final private = isPrivate ?? _defaultPrivate(mediaType);

      _logger.info('Media upload start: $name (kind=$kind, $contentType, ${size}B)');

      // Step 1 — init (Frontend.md §6.2)
      final initResp = await _api.post(
        '/api/v1/media/upload/init',
        data: {
          'kind': kind,
          'content_type': contentType,
          'original_name': name,
          'is_private': private,
          'size_bytes': size,
        },
      );
      final init = MediaUploadInit.fromJson(
        Map<String, dynamic>.from(initResp as Map),
      );
      _logger.debug('Media init OK: mediaId=${init.mediaId}');

      // Step 2 — PUT raw bytes directly to S3 (no auth headers).
      // Content-Type MUST match the init content_type or the presigned
      // signature won't match (S3 → 403 SignatureDoesNotMatch).
      final bytes = await file.readAsBytes();
      await _s3Dio.put(
        init.uploadUrl,
        data: Stream.fromIterable([bytes]),
        options: Options(
          headers: {
            'Content-Length': size,
            'Content-Type': contentType,
          },
          sendTimeout: const Duration(minutes: 5),
          receiveTimeout: const Duration(minutes: 5),
        ),
      );
      _logger.debug('Media S3 upload OK');

      // Step 3 — confirm (POST /media/{media_id}/confirm, no body).
      // Response is {media_id, status, size_bytes} — read media_id directly.
      final confirmResp = await _api.post(
        '/api/v1/media/${init.mediaId}/confirm',
      );
      final confirmMap = Map<String, dynamic>.from(confirmResp as Map);
      final mediaId = confirmMap['media_id'] as String? ?? init.mediaId;
      _logger.info('Media upload confirmed: $mediaId');
      return Right(mediaId);
    } catch (e) {
      _logger.error('uploadFile failed ($mediaType)', e);
      if (e is ServerException) return Left(ServerFailure(message: e.message, code: e.errorCode));
      return const Left(UnknownFailure());
    }
  }

  @override
  Future<Either<AppFailure, String>> getMediaUrl(String mediaId) async {
    try {
      // View URL for inline display (Frontend.md §6.3).
      final resp = await _api.get('/api/v1/media/$mediaId/url');
      final map = Map<String, dynamic>.from(resp as Map);
      final url = map['url'] as String? ?? '';
      return Right(url);
    } catch (e) {
      _logger.error('getMediaUrl failed ($mediaId)', e);
      if (e is ServerException) return Left(ServerFailure(message: e.message, code: e.errorCode));
      return const Left(UnknownFailure());
    }
  }

  @override
  Future<Either<AppFailure, String>> downloadToCache(
    String mediaId,
    String filename,
  ) async {
    try {
      // Resolve the presigned URL, then stream the bytes to a cache file.
      final urlResp = await _api.get('/api/v1/media/$mediaId/url');
      final url = Map<String, dynamic>.from(urlResp as Map)['url'] as String? ?? '';
      if (url.isEmpty) {
        _logger.error('downloadToCache got empty url ($mediaId)');
        return const Left(UnknownFailure());
      }

      // Sanitize the filename so it is a safe single path segment.
      final safeName = filename.replaceAll(RegExp(r'[\\/:*?"<>|]'), '_');
      final dir = await getTemporaryDirectory();
      final localPath = '${dir.path}/$safeName';

      // Presigned URL is public — use the bare S3 Dio (no auth interceptors).
      await _s3Dio.download(url, localPath);
      _logger.info('Downloaded media $mediaId to $localPath');
      return Right(localPath);
    } catch (e) {
      _logger.error('downloadToCache failed ($mediaId)', e);
      if (e is ServerException) return Left(ServerFailure(message: e.message, code: e.errorCode));
      return const Left(UnknownFailure());
    }
  }

  /// Maps a message/media type to a media-service `kind` (Frontend.md §6.1).
  String _kindFor(String mediaType) => switch (mediaType) {
        'AVATAR' => 'AVATAR',
        'CHAT_AVATAR' => 'CHAT_AVATAR',
        'VOICE' => 'VOICE_MESSAGE',
        'IMAGE' || 'VIDEO' || 'FILE' => 'ATTACHMENT',
        _ => 'ATTACHMENT',
      };

  /// Avatars are public; everything else defaults to private.
  bool _defaultPrivate(String mediaType) =>
      !(mediaType == 'AVATAR' || mediaType == 'CHAT_AVATAR');

  String _contentType(String mediaType, String filename) {
    final ext = filename.contains('.') ? '.${filename.split('.').last.toLowerCase()}' : '';
    return switch (mediaType) {
      'AVATAR' || 'CHAT_AVATAR' || 'IMAGE' =>
        ext == '.png' ? 'image/png' : 'image/jpeg',
      'VIDEO' => 'video/mp4',
      // audio_waveforms records AAC in an MPEG-4 container (.m4a).
      'VOICE' => 'audio/mp4',
      'FILE' => 'application/octet-stream',
      _ => 'application/octet-stream',
    };
  }
}
