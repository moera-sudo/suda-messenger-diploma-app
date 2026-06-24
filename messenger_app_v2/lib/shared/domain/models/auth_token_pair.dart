import 'package:json_annotation/json_annotation.dart';

part 'auth_token_pair.g.dart';  // Связывание с автоматически генерируемым файлом 

@JsonSerializable()
class AuthTokenPair {
  @JsonKey(name: 'access_token')
  final String accessToken;

  @JsonKey(name: 'refresh_token')
  final String refreshToken;

  AuthTokenPair({required this.accessToken, required this.refreshToken});

  factory AuthTokenPair.fromJson(Map<String, dynamic> json) => _$AuthTokenPairFromJson(json);
  Map<String, dynamic> toJson() => _$AuthTokenPairToJson(this);
}