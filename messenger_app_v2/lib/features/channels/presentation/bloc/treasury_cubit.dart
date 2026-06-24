import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:injectable/injectable.dart';

import '../../data/models/treasury_models.dart';
import '../../domain/repositories/channel_repository.dart';

// ── State ─────────────────────────────────────────────────────────

enum TreasuryStatus { initial, loading, success, failure, withdrawing, withdrawSuccess }

class TreasuryState extends Equatable {
  final TreasuryStatus status;
  final TreasuryStats? stats;
  final List<DonationItem> donations;
  final String? error;

  const TreasuryState({
    this.status = TreasuryStatus.initial,
    this.stats,
    this.donations = const [],
    this.error,
  });

  TreasuryState copyWith({
    TreasuryStatus? status,
    TreasuryStats? stats,
    List<DonationItem>? donations,
    String? error,
  }) =>
      TreasuryState(
        status: status ?? this.status,
        stats: stats ?? this.stats,
        donations: donations ?? this.donations,
        error: error ?? this.error,
      );

  @override
  List<Object?> get props => [status, stats, donations, error];
}

// ── Cubit ─────────────────────────────────────────────────────────

@injectable
class TreasuryCubit extends Cubit<TreasuryState> {
  final ChannelRepository _repo;

  TreasuryCubit(this._repo) : super(const TreasuryState());

  Future<void> load(String channelId) async {
    emit(state.copyWith(status: TreasuryStatus.loading));

    final statsResult = await _repo.getTreasury(channelId);
    final donationsResult = await _repo.getDonations(channelId);

    final stats = statsResult.fold((_) => null, (s) => s);
    final donations = donationsResult.fold((_) => <DonationItem>[], (d) => d);

    if (stats == null) {
      emit(state.copyWith(
        status: TreasuryStatus.failure,
        error: statsResult.fold((f) => f.message, (_) => null),
      ));
      return;
    }

    emit(state.copyWith(
      status: TreasuryStatus.success,
      stats: stats,
      donations: donations,
    ));
  }

  Future<bool> withdraw(String channelId, String amountWei) async {
    emit(state.copyWith(status: TreasuryStatus.withdrawing));
    final result = await _repo.withdrawFromTreasury(channelId, amountWei);
    return result.fold(
      (f) {
        emit(state.copyWith(status: TreasuryStatus.success, error: f.message));
        return false;
      },
      (_) {
        emit(state.copyWith(status: TreasuryStatus.withdrawSuccess));
        return true;
      },
    );
  }
}
