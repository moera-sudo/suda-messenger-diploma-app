/// Response for POST /api/v1/tx/wallet/transfer (202 Accepted).
class TransferResponse {
  final String txHash;
  final String fromAddress;
  final String toAddress;
  final String toUserId;
  final String amountWei;
  final String chatId;
  final String submittedAt;

  const TransferResponse({
    required this.txHash,
    required this.fromAddress,
    required this.toAddress,
    required this.toUserId,
    required this.amountWei,
    required this.chatId,
    required this.submittedAt,
  });

  factory TransferResponse.fromJson(Map<String, dynamic> json) => TransferResponse(
        txHash: json['tx_hash'] as String? ?? '',
        fromAddress: json['from_address'] as String? ?? '',
        toAddress: json['to_address'] as String? ?? '',
        toUserId: json['to_user_id'] as String? ?? '',
        amountWei: json['amount_wei'] as String? ?? '0',
        chatId: json['chat_id'] as String? ?? '',
        submittedAt: json['submitted_at'] as String? ?? '',
      );
}

/// Response for GET /api/v1/tx/wallet/me — the current user's wallet.
class MyWallet {
  final String address;
  final String sudaBalanceWei;
  final int decimals;

  const MyWallet({
    required this.address,
    required this.sudaBalanceWei,
    this.decimals = 18,
  });

  factory MyWallet.fromJson(Map<String, dynamic> json) => MyWallet(
        address: json['address'] as String? ?? '',
        sudaBalanceWei: json['suda_balance_wei']?.toString() ?? '0',
        decimals: (json['decimals'] as num?)?.toInt() ?? 18,
      );
}

/// Converts wei string to a human-readable SUDA display string.
/// e.g. "1500000000000000000" → "1.5", "1000000000000000000" → "1"
/// Trailing zeros after the decimal point are stripped.
String weiToSudaDisplay(String wei) {
  try {
    final bigWei = BigInt.parse(wei);
    const decimals = 18;
    final divisor = BigInt.from(10).pow(decimals);
    final intPart = bigWei ~/ divisor;
    final fracBig = bigWei.remainder(divisor).abs();
    if (fracBig == BigInt.zero) return intPart.toString();
    // Pad fractional part to 18 digits, then strip trailing zeros.
    final fracStr = fracBig.toString().padLeft(decimals, '0').replaceAll(RegExp(r'0+$'), '');
    return '$intPart.$fracStr';
  } catch (_) {
    return wei;
  }
}

/// Converts a human-readable SUDA amount string to wei (BigInt).
/// e.g. "1.5" → 1500000000000000000
///
/// Uses pure BigInt arithmetic (no double) so large amounts never overflow
/// int64 — `(amount * 1e18).round()` silently clamped to 9223372036854775807
/// for any amount ≳ 9.22 SUDA, which is the bug this replaces.
BigInt sudaToWei(String suda) {
  const decimals = 18;
  final trimmed = suda.trim();
  if (trimmed.isEmpty) return BigInt.zero;

  final parts = trimmed.split('.');
  if (parts.length > 2) return BigInt.zero; // malformed

  final intPart = parts[0].isEmpty ? BigInt.zero : (BigInt.tryParse(parts[0]) ?? BigInt.zero);
  var wei = intPart * BigInt.from(10).pow(decimals);

  if (parts.length == 2 && parts[1].isNotEmpty) {
    // Truncate/pad the fractional part to exactly `decimals` digits.
    var frac = parts[1];
    frac = frac.length > decimals ? frac.substring(0, decimals) : frac.padRight(decimals, '0');
    wei += BigInt.tryParse(frac) ?? BigInt.zero;
  }
  return wei;
}
