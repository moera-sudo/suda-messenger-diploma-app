// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'auth_request_models.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

RegisterRequest _$RegisterRequestFromJson(Map<String, dynamic> json) =>
    RegisterRequest(
      username: json['username'] as String,
      displayName: json['name'] as String,
      email: json['email'] as String,
      password: json['password'] as String,
    );

Map<String, dynamic> _$RegisterRequestToJson(RegisterRequest instance) =>
    <String, dynamic>{
      'username': instance.username,
      'name': instance.displayName,
      'email': instance.email,
      'password': instance.password,
    };

LoginRequest _$LoginRequestFromJson(Map<String, dynamic> json) => LoginRequest(
  email: json['email'] as String,
  password: json['password'] as String,
  userAgent: json['user_agent'] as String? ?? 'FlutterApp',
  clientIp: json['client_ip'] as String? ?? '127.0.0.1',
);

Map<String, dynamic> _$LoginRequestToJson(LoginRequest instance) =>
    <String, dynamic>{
      'email': instance.email,
      'password': instance.password,
      'user_agent': instance.userAgent,
      'client_ip': instance.clientIp,
    };

VerifyRequest _$VerifyRequestFromJson(Map<String, dynamic> json) =>
    VerifyRequest(email: json['email'] as String, code: json['code'] as String);

Map<String, dynamic> _$VerifyRequestToJson(VerifyRequest instance) =>
    <String, dynamic>{'email': instance.email, 'code': instance.code};
