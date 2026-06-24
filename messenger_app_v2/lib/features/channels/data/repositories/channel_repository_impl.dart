import 'package:dartz/dartz.dart';
import 'package:injectable/injectable.dart';

import '../../../../shared/data/api/api_client.dart';
import '../../../../shared/data/api/server_exception.dart';
import '../../../../shared/data/logger/app_logger.dart';
import '../../../../shared/domain/models/app_failure.dart';
import '../../domain/repositories/channel_repository.dart';
import '../models/channel_comment_model.dart';
import '../models/channel_settings_model.dart';
import '../models/channel_view_model.dart';
import '../models/gating_rule.dart';
import '../models/treasury_models.dart';

@LazySingleton(as: ChannelRepository)
class ChannelRepositoryImpl implements ChannelRepository {
  final ApiClient _api;
  final AppLogger _logger;

  static const _channelsBase = '/api/v1/messenger/channels';
  static const _gatingBase = '/api/v1/tx/gating/rule';

  ChannelRepositoryImpl(this._api, this._logger);

  @override
  Future<Either<AppFailure, Unit>> subscribe(String chatId) async {
    try {
      await _api.post('$_channelsBase/$chatId/subscribe', data: {});
      _logger.info('Subscribed to channel $chatId');
      return const Right(unit);
    } catch (e) {
      _logger.error('subscribe failed (chatId=$chatId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> unsubscribe(String chatId) async {
    try {
      await _api.post('$_channelsBase/$chatId/unsubscribe', data: {});
      _logger.info('Unsubscribed from channel $chatId');
      return const Right(unit);
    } catch (e) {
      _logger.error('unsubscribe failed (chatId=$chatId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, GatingRule?>> getGatingRule(String chatId) async {
    try {
      final response = await _api.get('$_gatingBase/$chatId');
      if (response == null) return const Right(null);
      final rule = GatingRule.fromJson(Map<String, dynamic>.from(response as Map));
      // has_rule=false → channel is open for everyone (no gating)
      if (!rule.hasRule) {
        _logger.debug('No active gating rule for $chatId');
        return const Right(null);
      }
      _logger.debug('Gating rule for $chatId: ${rule.subscriptionPrice} SUDA (subscription price)');
      return Right(rule);
    } catch (e) {
      // 404 = no gating rule → not an error (defensive fallback)
      if (e is ServerException && e.statusCode == 404) {
        _logger.debug('No gating rule for $chatId');
        return const Right(null);
      }
      _logger.error('getGatingRule failed (chatId=$chatId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> setGatingRule(String chatId, BigInt minSudaBalanceWei) async {
    try {
      await _api.post(_gatingBase, data: {
        'chat_id': chatId,
        'subscription_price_wei': minSudaBalanceWei.toString(),
      });
      _logger.info('Gating rule set for $chatId: $minSudaBalanceWei wei (subscription_price_wei)');
      return const Right(unit);
    } catch (e) {
      _logger.error('setGatingRule failed ($chatId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> deleteGatingRule(String chatId) async {
    try {
      await _api.delete('$_gatingBase/$chatId');
      _logger.info('Gating rule deleted for $chatId');
      return const Right(unit);
    } catch (e) {
      _logger.error('deleteGatingRule failed ($chatId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, ChannelViewModel>> getChannelView(String chatId) async {
    try {
      final resp = await _api.get('$_channelsBase/$chatId/view');
      return Right(ChannelViewModel.fromJson(Map<String, dynamic>.from(resp as Map)));
    } catch (e) {
      _logger.error('getChannelView failed ($chatId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> sendJoinRequest(String chatId) async {
    try {
      await _api.post('$_channelsBase/$chatId/join-request', data: {});
      _logger.info('Join request sent for channel $chatId');
      return const Right(unit);
    } catch (e) {
      _logger.error('sendJoinRequest failed ($chatId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> cancelJoinRequest(String chatId) async {
    try {
      await _api.delete('$_channelsBase/$chatId/join-request');
      _logger.info('Join request cancelled for channel $chatId');
      return const Right(unit);
    } catch (e) {
      _logger.error('cancelJoinRequest failed ($chatId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> acceptInvite(String chatId) async {
    try {
      await _api.post('$_channelsBase/$chatId/invites/accept', data: {});
      _logger.info('Invite accepted for channel $chatId');
      return const Right(unit);
    } catch (e) {
      _logger.error('acceptInvite failed ($chatId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> declineInvite(String chatId) async {
    try {
      await _api.post('$_channelsBase/$chatId/invites/decline', data: {});
      _logger.info('Invite declined for channel $chatId');
      return const Right(unit);
    } catch (e) {
      _logger.error('declineInvite failed ($chatId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, List<ChannelInviteItem>>> getMyInvites() async {
    try {
      final resp = await _api.get('$_channelsBase/invites');
      final List raw = resp as List? ?? [];
      final items = raw
          .whereType<Map<String, dynamic>>()
          .map(ChannelInviteItem.fromJson)
          .toList();
      return Right(items);
    } catch (e) {
      _logger.error('getMyInvites failed', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> inviteUser(String chatId, String username) async {
    try {
      await _api.post('$_channelsBase/$chatId/invites', data: {'username': username});
      _logger.info('User @$username invited to channel $chatId');
      return const Right(unit);
    } catch (e) {
      _logger.error('inviteUser failed ($chatId, @$username)', e);
      return Left(_handleError(e));
    }
  }

  // ─── Settings ──────────────────────────────────────────────

  @override
  Future<Either<AppFailure, ChannelSettings>> getChannelSettings(String chatId) async {
    try {
      final resp = await _api.get('$_channelsBase/$chatId/settings');
      return Right(ChannelSettings.fromJson(Map<String, dynamic>.from(resp as Map)));
    } catch (e) {
      _logger.error('getChannelSettings failed ($chatId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> updateChannelSettings(
    String chatId, {
    bool? commentsEnabled,
    String? visibility,
    String? username,
  }) async {
    try {
      final body = <String, dynamic>{};
      if (commentsEnabled != null) body['comments_enabled'] = commentsEnabled;
      if (visibility != null) body['visibility'] = visibility;
      if (username != null) body['username'] = username;
      await _api.put('$_channelsBase/$chatId/settings', data: body);
      _logger.info('Channel settings updated ($chatId): ${body.keys.toList()}');
      return const Right(unit);
    } catch (e) {
      _logger.error('updateChannelSettings failed ($chatId)', e);
      return Left(_handleError(e));
    }
  }

  // ─── Comments ──────────────────────────────────────────────

  @override
  Future<Either<AppFailure, (List<ChannelComment>, int)>> getComments(
    String chatId,
    int postId, {
    int limit = 50,
    int offset = 0,
  }) async {
    try {
      final resp = await _api.get(
        '$_channelsBase/$chatId/posts/$postId/comments',
        query: {'limit': limit, 'offset': offset},
      );
      final map = Map<String, dynamic>.from(resp as Map);
      final List raw = map['comments'] as List? ?? [];
      final comments = raw
          .whereType<Map>()
          .map((e) => ChannelComment.fromJson(Map<String, dynamic>.from(e)))
          .toList();
      final total = (map['total'] as num?)?.toInt() ?? comments.length;
      return Right((comments, total));
    } catch (e) {
      _logger.error('getComments failed ($chatId, post=$postId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, ChannelComment>> addComment(
    String chatId,
    int postId,
    String content, {
    int? replyToCommentId,
  }) async {
    try {
      final body = <String, dynamic>{'content': content};
      if (replyToCommentId != null) body['reply_to_comment_id'] = replyToCommentId;
      final resp = await _api.post('$_channelsBase/$chatId/posts/$postId/comments', data: body);
      _logger.info('Comment added to post $postId in channel $chatId');
      return Right(ChannelComment.fromJson(Map<String, dynamic>.from(resp as Map)));
    } catch (e) {
      _logger.error('addComment failed ($chatId, post=$postId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> editComment(int commentId, String content) async {
    try {
      await _api.put('$_channelsBase/comments/$commentId', data: {'content': content});
      _logger.info('Comment $commentId edited');
      return const Right(unit);
    } catch (e) {
      _logger.error('editComment failed ($commentId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> deleteComment(int commentId) async {
    try {
      await _api.delete('$_channelsBase/comments/$commentId');
      _logger.info('Comment $commentId deleted');
      return const Right(unit);
    } catch (e) {
      _logger.error('deleteComment failed ($commentId)', e);
      return Left(_handleError(e));
    }
  }

  // ─── Treasury ──────────────────────────────────────────────────

  static const _treasuryBase = '/api/v1/tx/wallet/channel';

  @override
  Future<Either<AppFailure, TreasuryStats>> getTreasury(String channelId) async {
    try {
      final resp = await _api.get('$_treasuryBase/$channelId/treasury');
      return Right(TreasuryStats.fromJson(Map<String, dynamic>.from(resp as Map)));
    } catch (e) {
      _logger.error('getTreasury failed ($channelId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, List<DonationItem>>> getDonations(
    String channelId, {
    int limit = 50,
    int offset = 0,
  }) async {
    try {
      final resp = await _api.get(
        '$_treasuryBase/$channelId/donations',
        query: {'limit': limit, 'offset': offset},
      );
      final map = Map<String, dynamic>.from(resp as Map);
      final items = (map['items'] as List? ?? [])
          .whereType<Map>()
          .map((e) => DonationItem.fromJson(Map<String, dynamic>.from(e)))
          .toList();
      return Right(items);
    } catch (e) {
      _logger.error('getDonations failed ($channelId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, WithdrawResponse>> withdrawFromTreasury(
    String channelId,
    String amountWei,
  ) async {
    try {
      final resp = await _api.post(
        '$_treasuryBase/$channelId/withdraw',
        data: {'amount_wei': amountWei},
      );
      _logger.info('Treasury withdrawal submitted for $channelId: $amountWei wei');
      return Right(WithdrawResponse.fromJson(Map<String, dynamic>.from(resp as Map)));
    } catch (e) {
      _logger.error('withdrawFromTreasury failed ($channelId)', e);
      return Left(_handleError(e));
    }
  }

  AppFailure _handleError(dynamic e) {
    if (e is ServerException) return ServerFailure(message: e.message, code: e.errorCode);
    return const UnknownFailure();
  }
}
