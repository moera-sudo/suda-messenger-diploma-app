/// GET /api/v1/tx/wallet/channel/{channel_id}/treasury
class TreasuryStats {
  final String channelId;
  final String address;
  final String sudaBalanceWei;
  final int decimals;
  final int totalDonationsCount;
  final String totalDonationsWei;
  final List<TopDonor> topDonors;

  const TreasuryStats({
    required this.channelId,
    required this.address,
    required this.sudaBalanceWei,
    this.decimals = 18,
    required this.totalDonationsCount,
    required this.totalDonationsWei,
    required this.topDonors,
  });

  factory TreasuryStats.fromJson(Map<String, dynamic> json) => TreasuryStats(
        channelId: json['channel_id'] as String? ?? '',
        address: json['address'] as String? ?? '',
        sudaBalanceWei: json['suda_balance_wei']?.toString() ?? '0',
        decimals: (json['decimals'] as num?)?.toInt() ?? 18,
        totalDonationsCount: (json['total_donations_count'] as num?)?.toInt() ?? 0,
        totalDonationsWei: json['total_donations_wei']?.toString() ?? '0',
        topDonors: (json['top_donors'] as List? ?? [])
            .whereType<Map>()
            .map((e) => TopDonor.fromJson(Map<String, dynamic>.from(e)))
            .toList(),
      );
}

class TopDonor {
  final String userId;
  final String username;
  final String displayName;
  final String amountWei;
  final int donationCount;

  const TopDonor({
    required this.userId,
    required this.username,
    required this.displayName,
    required this.amountWei,
    required this.donationCount,
  });

  factory TopDonor.fromJson(Map<String, dynamic> json) => TopDonor(
        userId: json['user_id'] as String? ?? '',
        username: json['username'] as String? ?? '',
        displayName: json['display_name'] as String? ?? '',
        amountWei: json['amount_wei']?.toString() ?? '0',
        donationCount: (json['donation_count'] as num?)?.toInt() ?? 0,
      );
}

/// Single item from GET /api/v1/tx/wallet/channel/{channel_id}/donations
class DonationItem {
  final String fromUsername;
  final String amountWei;
  final String? message;
  final String txHash;
  final String createdAt;

  const DonationItem({
    required this.fromUsername,
    required this.amountWei,
    this.message,
    required this.txHash,
    required this.createdAt,
  });

  factory DonationItem.fromJson(Map<String, dynamic> json) => DonationItem(
        fromUsername: json['from_username'] as String? ?? '',
        amountWei: json['amount_wei']?.toString() ?? '0',
        message: json['message'] as String?,
        txHash: json['tx_hash'] as String? ?? '',
        createdAt: json['created_at'] as String? ?? '',
      );
}

/// Response for POST /api/v1/tx/wallet/channel/{channel_id}/withdraw
class WithdrawResponse {
  final String txHash;
  final String fromAddress;
  final String toAddress;
  final String amountWei;
  final String submittedAt;

  const WithdrawResponse({
    required this.txHash,
    required this.fromAddress,
    required this.toAddress,
    required this.amountWei,
    required this.submittedAt,
  });

  factory WithdrawResponse.fromJson(Map<String, dynamic> json) => WithdrawResponse(
        txHash: json['tx_hash'] as String? ?? '',
        fromAddress: json['from_address'] as String? ?? '',
        toAddress: json['to_address'] as String? ?? '',
        amountWei: json['amount_wei']?.toString() ?? '0',
        submittedAt: json['submitted_at'] as String? ?? '',
      );
}
