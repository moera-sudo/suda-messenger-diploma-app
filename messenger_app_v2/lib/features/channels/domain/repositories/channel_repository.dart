import 'package:dartz/dartz.dart';

import '../../../../shared/domain/models/app_failure.dart';
import '../../data/models/channel_comment_model.dart';
import '../../data/models/channel_settings_model.dart';
import '../../data/models/channel_view_model.dart';
import '../../data/models/gating_rule.dart';
import '../../data/models/treasury_models.dart';

abstract class ChannelRepository {
  /// POST /api/v1/messenger/channels/{chatId}/subscribe
  Future<Either<AppFailure, Unit>> subscribe(String chatId);

  /// POST /api/v1/messenger/channels/{chatId}/unsubscribe
  Future<Either<AppFailure, Unit>> unsubscribe(String chatId);

  /// GET /api/v1/tx/gating/rule/{chatId} → null if channel is open (has_rule=false)
  Future<Either<AppFailure, GatingRule?>> getGatingRule(String chatId);

  /// POST /api/v1/tx/gating/rule — create/update a token-gating rule (OWNER).
  Future<Either<AppFailure, Unit>> setGatingRule(String chatId, BigInt minSudaBalanceWei);

  /// DELETE /api/v1/tx/gating/rule/{chatId} — remove the rule (OWNER).
  Future<Either<AppFailure, Unit>> deleteGatingRule(String chatId);

  /// GET /api/v1/messenger/channels/{chatId}/view
  Future<Either<AppFailure, ChannelViewModel>> getChannelView(String chatId);

  /// POST /api/v1/messenger/channels/{chatId}/join-request
  Future<Either<AppFailure, Unit>> sendJoinRequest(String chatId);

  /// DELETE /api/v1/messenger/channels/{chatId}/join-request
  Future<Either<AppFailure, Unit>> cancelJoinRequest(String chatId);

  /// POST /api/v1/messenger/channels/{chatId}/invites/accept
  Future<Either<AppFailure, Unit>> acceptInvite(String chatId);

  /// POST /api/v1/messenger/channels/{chatId}/invites/decline
  Future<Either<AppFailure, Unit>> declineInvite(String chatId);

  /// GET /api/v1/messenger/channels/invites
  Future<Either<AppFailure, List<ChannelInviteItem>>> getMyInvites();

  /// POST /api/v1/messenger/channels/{chatId}/invites   body: {username}
  Future<Either<AppFailure, Unit>> inviteUser(String chatId, String username);

  // ─── Settings (OWNER/ADMIN) ────────────────────────────────

  /// GET /api/v1/messenger/channels/{chatId}/settings
  Future<Either<AppFailure, ChannelSettings>> getChannelSettings(String chatId);

  /// PUT /api/v1/messenger/channels/{chatId}/settings — partial update.
  Future<Either<AppFailure, Unit>> updateChannelSettings(
    String chatId, {
    bool? commentsEnabled,
    String? visibility,
    String? username,
  });

  // ─── Comments ──────────────────────────────────────────────

  /// GET /api/v1/messenger/channels/{chatId}/posts/{postId}/comments
  /// Returns the comment list and the total count.
  Future<Either<AppFailure, (List<ChannelComment>, int)>> getComments(
    String chatId,
    int postId, {
    int limit = 50,
    int offset = 0,
  });

  /// POST /api/v1/messenger/channels/{chatId}/posts/{postId}/comments
  Future<Either<AppFailure, ChannelComment>> addComment(
    String chatId,
    int postId,
    String content, {
    int? replyToCommentId,
  });

  /// PUT /api/v1/messenger/channels/comments/{commentId}
  Future<Either<AppFailure, Unit>> editComment(int commentId, String content);

  /// DELETE /api/v1/messenger/channels/comments/{commentId}
  Future<Either<AppFailure, Unit>> deleteComment(int commentId);

  // ─── Treasury (OWNER/ADMIN) ────────────────────────────────────

  /// GET /api/v1/tx/wallet/channel/{channelId}/treasury
  Future<Either<AppFailure, TreasuryStats>> getTreasury(String channelId);

  /// GET /api/v1/tx/wallet/channel/{channelId}/donations?limit=&offset=
  Future<Either<AppFailure, List<DonationItem>>> getDonations(
    String channelId, {
    int limit = 50,
    int offset = 0,
  });

  /// POST /api/v1/tx/wallet/channel/{channelId}/withdraw  (OWNER only)
  Future<Either<AppFailure, WithdrawResponse>> withdrawFromTreasury(
    String channelId,
    String amountWei,
  );
}
