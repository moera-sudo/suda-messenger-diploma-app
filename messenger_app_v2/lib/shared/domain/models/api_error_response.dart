import 'package:json_annotation/json_annotation.dart';

part 'api_error_response.g.dart';

@JsonSerializable()
class ApiErrorResponse {
  @JsonKey(name: 'error')
  final String error; // Код ошибки, например "BAD_REQUEST"

  @JsonKey(name: 'message')
  final String? message; // Сообщение, например "Invalid Input"

  ApiErrorResponse({required this.error, this.message});

  factory ApiErrorResponse.fromJson(Map<String, dynamic> json) => 
      _$ApiErrorResponseFromJson(json);
      
  Map<String, dynamic> toJson() => _$ApiErrorResponseToJson(this);
}