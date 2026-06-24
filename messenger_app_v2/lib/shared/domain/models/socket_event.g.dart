// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'socket_event.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

SocketEvent _$SocketEventFromJson(Map<String, dynamic> json) => SocketEvent(
  type: json['type'] as String,
  payload: json['payload'] as Map<String, dynamic>?,
);

Map<String, dynamic> _$SocketEventToJson(SocketEvent instance) =>
    <String, dynamic>{'type': instance.type, 'payload': instance.payload};
