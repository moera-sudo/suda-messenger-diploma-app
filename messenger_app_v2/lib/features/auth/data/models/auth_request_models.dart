import 'package:json_annotation/json_annotation.dart';

part 'auth_request_models.g.dart';

@JsonSerializable()
class RegisterRequest {
  final String username;
  @JsonKey(name: 'name')
  final String displayName;
  final String email;
  final String password;

  RegisterRequest({
    required this.username,
    required this.displayName,
    required this.email,
    required this.password,
  });

  Map<String, dynamic> toJson() => _$RegisterRequestToJson(this);
}

@JsonSerializable()
class LoginRequest {
  final String email;
  final String password;
  @JsonKey(name: 'user_agent')
  final String userAgent;
  @JsonKey(name: 'client_ip')
  final String clientIp;

  LoginRequest({
    required this.email,
    required this.password,
    this.userAgent = 'FlutterApp',
    this.clientIp = '127.0.0.1', // Бэк сам подставит RealIP
  });

  Map<String, dynamic> toJson() => _$LoginRequestToJson(this);
}

@JsonSerializable()
class VerifyRequest {
  final String email;
  final String code;

  VerifyRequest({required this.email, required this.code});

  Map<String, dynamic> toJson() => _$VerifyRequestToJson(this);
}

class ForgotPasswordRequest {
  final String email;

  const ForgotPasswordRequest({required this.email});

  Map<String, dynamic> toJson() => {'email': email};
}

class ResetPasswordRequest {
  final String email;
  final String code;
  final String newPassword;

  const ResetPasswordRequest({
    required this.email,
    required this.code,
    required this.newPassword,
  });

  Map<String, dynamic> toJson() => {
    'email': email,
    'code': code,
    // Backend contract: json tag is exactly `Password`
    'Password': newPassword,
  };
}
