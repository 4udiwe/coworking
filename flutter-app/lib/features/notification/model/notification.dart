import 'dart:convert';

import 'package:intl/intl.dart';

class NotificationPayload {
  final String? type;
  final String? bookingId;
  final String? placeId;
  final String? placeLabel;
  final DateTime? startTime;
  final DateTime? endTime;
  final Map<String, dynamic> raw;

  NotificationPayload({
    this.type,
    this.bookingId,
    this.placeId,
    this.placeLabel,
    this.startTime,
    this.endTime,
    required this.raw,
  });

  factory NotificationPayload.fromJson(Map<String, dynamic> json) {
    return NotificationPayload(
      type: json['type'] as String?,
      bookingId: json['bookingId'] as String?,
      placeId: json['placeId'] as String?,
      placeLabel: json['placeLabel'] as String?,
      startTime: json['startTime'] != null ?  DateFormat("yyyy-MM-dd HH:mm:ss ZZZZ", "ru").parse(json['startTime'], true).toLocal() : null,
      endTime: json['endTime'] != null ? DateFormat("yyyy-MM-dd HH:mm:ss ZZZZ", "ru").parse(json['endTime'], true).toLocal() : null,
      raw: json,
    );
  }
}

final _timeFormat = DateFormat('HH:mm', 'ru');


class NotificationModel {
  final String id;
  final String userId;
  final String type;
  final String title;
  final String body;
  final NotificationPayload? payload;
  final String? actionUrl;
  final bool isRead;
  final DateTime createdAt;
  final DateTime? readAt;

  String get displayBody {
    if (payload == null) return body;

    final p = payload!;
    final label = p.placeLabel ?? 'Place label';

    switch (type) {
      case 'booking_created':
        return 'Место $label забронировано на ${p.startTime != null ? _timeFormat.format(p.startTime!.toLocal()) : "time"}';
      case 'booking_cancelled':
        return 'Бронирование места $label отменено';
      case 'booking_reminder':
        return 'Время брони места $label начинается в ${p.startTime != null ? _timeFormat.format(p.startTime!.toLocal()) : "time"}';
      case 'booking_expired':
        return 'Бронирование места $label завершено в ${p.endTime != null ? _timeFormat.format(p.endTime!.toLocal()) : "time"}';
      default:
        return body;
    }
  }

  NotificationModel({
    required this.id,
    required this.userId,
    required this.type,
    required this.title,
    required this.body,
    this.payload,
    this.actionUrl,
    required this.isRead,
    required this.createdAt,
    this.readAt,
  });

  factory NotificationModel.fromJson(Map<String, dynamic> json) {
    NotificationPayload? parsedPayload;
    if (json['payload'] != null) {
      try {
        final decodedString = utf8.decode(base64Decode(json['payload']));
        final payloadMap = jsonDecode(decodedString) as Map<String, dynamic>;
        parsedPayload = NotificationPayload.fromJson(payloadMap);
      } catch (e) {
        // Fallback or log error
      }
    }

    return NotificationModel(
      id: json['id'],
      userId: json['userId'],
      type: json['type'],
      title: json['title'],
      body: json['body'],
      payload: parsedPayload,
      actionUrl: json['actionUrl'],
      isRead: json['isRead'],
      createdAt: DateTime.parse(json['createdAt']),
      readAt: json['readAt'] != null ? DateTime.parse(json['readAt']) : null,
    );
  }

  NotificationModel copyWith({
    String? id,
    String? userId,
    String? type,
    String? title,
    String? body,
    NotificationPayload? Function()? payload,
    String? Function()? actionUrl,
    bool? isRead,
    DateTime? createdAt,
    DateTime? Function()? readAt,
  }) {
    return NotificationModel(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      type: type ?? this.type,
      title: title ?? this.title,
      body: body ?? this.body,
      payload: payload != null ? payload() : this.payload,
      actionUrl: actionUrl != null ? actionUrl() : this.actionUrl,
      isRead: isRead ?? this.isRead,
      createdAt: createdAt ?? this.createdAt,
      readAt: readAt != null ? readAt() : this.readAt,
    );
  }
}

class RegisterDeviceTokenRequest {
  final String deviceToken;
  final String platform; // ios, android

  RegisterDeviceTokenRequest({required this.deviceToken, required this.platform});

  Map<String, dynamic> toJson() => {
    'deviceToken': deviceToken,
    'platform': platform,
  };
}
