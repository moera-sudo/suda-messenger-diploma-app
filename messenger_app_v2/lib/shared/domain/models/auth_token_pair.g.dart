// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'auth_token_pair.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

AuthTokenPair _$AuthTokenPairFromJson(Map<String, dynamic> json) =>
    AuthTokenPair(
      accessToken: json['access_token'] as String,
      refreshToken: json['refresh_token'] as String,
    );

Map<String, dynamic> _$AuthTokenPairToJson(AuthTokenPair instance) =>
    <String, dynamic>{
      'access_token': instance.accessToken,
      'refresh_token': instance.refreshToken,
    };
