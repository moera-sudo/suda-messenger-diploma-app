import 'package:dartz/dartz.dart';

import '../../../../shared/domain/models/app_failure.dart';
import '../../data/models/wallet_models.dart';

abstract class WalletRepository {
  /// GET /api/v1/tx/wallet/me — current user's wallet (address + SUDA balance).
  /// May 404 for a few seconds after signup while the wallet is created.
  Future<Either<AppFailure, MyWallet>> getMyWallet();

  /// POST /api/v1/tx/wallet/transfer
  /// In-chat P2P transfer. chat_id links the tx to a chat system message.
  Future<Either<AppFailure, TransferResponse>> transferInChat({
    required String toUsername,
    required BigInt amountWei,
    required String chatId,
    String? note,
  });

  /// POST /api/v1/tx/donate
  /// Donate SUDA to a channel treasury. chat_id is optional.
  Future<Either<AppFailure, Unit>> donateToChannel({
    required String toChannelId,
    required BigInt amountWei,
    String? message,
    String? chatId,
  });
}
