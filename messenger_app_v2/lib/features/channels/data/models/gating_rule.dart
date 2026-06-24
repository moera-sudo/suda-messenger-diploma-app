import 'package:json_annotation/json_annotation.dart';

part 'gating_rule.g.dart';

/// Token-gating / paid-subscription rule for a channel.
/// Iteration-2: gating now means a one-time subscription price paid in SUDA.
@JsonSerializable()
class GatingRule {
  @JsonKey(name: 'chat_id')
  final String chatId;

  /// Subscription price in wei (Iteration-2 paid model). Primary field.
  @JsonKey(name: 'subscription_price_wei', defaultValue: '0')
  final String subscriptionPriceWei;

  /// Legacy min-balance field kept for backward compatibility.
  @JsonKey(name: 'min_suda_balance_wei', defaultValue: '0')
  final String minSudaBalanceWei;

  @JsonKey(name: 'required_nft_collection_id')
  final String? requiredNftCollectionId;

  /// If false — channel is open for everyone (no gating).
  @JsonKey(name: 'has_rule', defaultValue: false)
  final bool hasRule;

  const GatingRule({
    required this.chatId,
    this.subscriptionPriceWei = '0',
    this.minSudaBalanceWei = '0',
    this.requiredNftCollectionId,
    this.hasRule = false,
  });

  /// Subscription price converted from wei to SUDA (human-readable double).
  double get subscriptionPrice =>
      BigInt.parse(subscriptionPriceWei).toDouble() / 1e18;

  /// Legacy getter — use [subscriptionPrice] for new code.
  double get minSudaBalance =>
      BigInt.parse(minSudaBalanceWei).toDouble() / 1e18;

  factory GatingRule.fromJson(Map<String, dynamic> json) =>
      _$GatingRuleFromJson(json);

  Map<String, dynamic> toJson() => _$GatingRuleToJson(this);
}
