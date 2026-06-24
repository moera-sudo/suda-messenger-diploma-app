import 'package:dartz/dartz.dart';
import 'package:injectable/injectable.dart';

import '../../../../shared/data/api/api_client.dart';
import '../../../../shared/data/api/server_exception.dart';
import '../../../../shared/data/logger/app_logger.dart';
import '../../../../shared/domain/models/app_failure.dart';
import '../../domain/repositories/wallet_repository.dart';
import '../models/wallet_models.dart';

@LazySingleton(as: WalletRepository)
class WalletRepositoryImpl implements WalletRepository {
  final ApiClient _api;
  final AppLogger _logger;

  static const _walletPath   = '/api/v1/tx/wallet/me';
  static const _transferPath = '/api/v1/tx/wallet/transfer';
  static const _donatePath   = '/api/v1/tx/donate';

  WalletRepositoryImpl(this._api, this._logger);

  @override
  Future<Either<AppFailure, MyWallet>> getMyWallet() async {
    try {
      final resp = await _api.get(_walletPath);
      final wallet = MyWallet.fromJson(Map<String, dynamic>.from(resp as Map));
      _logger.debug('Loaded wallet: balance=${wallet.sudaBalanceWei} wei');
      return Right(wallet);
    } catch (e) {
      _logger.error('getMyWallet failed', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, TransferResponse>> transferInChat({
    required String toUsername,
    required BigInt amountWei,
    required String chatId,
    String? note,
  }) async {
    try {
      // Strip any leading '@' so we never send "@@name" / "@name" to the backend.
      final username = toUsername.trim().replaceFirst(RegExp(r'^@+'), '');
      final body = <String, dynamic>{
        'to_username': username,
        'amount_wei': amountWei.toString(),
        'chat_id': chatId,
      };
      if (note != null && note.isNotEmpty) body['note'] = note;

      _logger.info('Initiating transfer of $amountWei wei to @$username in chat $chatId');
      final resp = await _api.post(_transferPath, data: body);
      final result = TransferResponse.fromJson(Map<String, dynamic>.from(resp as Map));
      _logger.info('Transfer submitted — tx_hash=${result.txHash}');
      return Right(result);
    } catch (e) {
      _logger.error('transferInChat failed (to=@$toUsername, chat=$chatId)', e);
      return Left(_handleError(e));
    }
  }

  @override
  Future<Either<AppFailure, Unit>> donateToChannel({
    required String toChannelId,
    required BigInt amountWei,
    String? message,
    String? chatId,
  }) async {
    try {
      final body = <String, dynamic>{
        'to_channel_id': toChannelId,
        'amount_wei': amountWei.toString(),
      };
      if (message != null && message.isNotEmpty) body['message'] = message;
      if (chatId != null && chatId.isNotEmpty) body['chat_id'] = chatId;

      _logger.info('Initiating donation of $amountWei wei to channel $toChannelId');
      await _api.post(_donatePath, data: body);
      _logger.info('Donation submitted to channel $toChannelId');
      return const Right(unit);
    } catch (e) {
      _logger.error('donateToChannel failed (channel=$toChannelId)', e);
      return Left(_handleError(e));
    }
  }

  AppFailure _handleError(dynamic e) {
    if (e is ServerException) return ServerFailure(message: e.message, code: e.errorCode);
    return const UnknownFailure();
  }
}
