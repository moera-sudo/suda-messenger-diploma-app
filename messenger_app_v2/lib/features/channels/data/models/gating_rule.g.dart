// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'gating_rule.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

GatingRule _$GatingRuleFromJson(Map<String, dynamic> json) => GatingRule(
  chatId: json['chat_id'] as String,
  subscriptionPriceWei: json['subscription_price_wei'] as String? ?? '0',
  minSudaBalanceWei: json['min_suda_balance_wei'] as String? ?? '0',
  requiredNftCollectionId: json['required_nft_collection_id'] as String?,
  hasRule: json['has_rule'] as bool? ?? false,
);

Map<String, dynamic> _$GatingRuleToJson(GatingRule instance) =>
    <String, dynamic>{
      'chat_id': instance.chatId,
      'subscription_price_wei': instance.subscriptionPriceWei,
      'min_suda_balance_wei': instance.minSudaBalanceWei,
      'required_nft_collection_id': instance.requiredNftCollectionId,
      'has_rule': instance.hasRule,
    };
