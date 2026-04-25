import 'dart:convert';

import '../../../core/api/api_client.dart';
import '../model/notification.dart';

class NotificationsResponse {
  final List<NotificationModel> notifications;

  NotificationsResponse({required this.notifications});

  factory NotificationsResponse.fromJson(Map<String, dynamic> json) {
    final list = json['notifications'] as List;

    return NotificationsResponse(
      notifications: list.map((e) => NotificationModel.fromJson(e)).toList(),
    );
  }
}

class NotificationRepository {
  final ApiClient apiClient;

  NotificationRepository({required this.apiClient});

  /// =======================
  /// GET NOTIFICATIONS
  /// =======================
  Future<NotificationsResponse> getNotifications({
    int? limit,
    int? offset,
    bool? isRead,
    DateTime? since,
  }) async {
    final response = await apiClient.get(
      '/notifications',
      queryParameters: {
        if (limit != null) 'limit': limit,
        if (offset != null) 'offset': offset,
        if (isRead != null) 'isRead': isRead,
        if (since != null) 'since': since.toUtc().toIso8601String(),
      },
    );

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      return NotificationsResponse.fromJson(data);
    }

    throw Exception('Failed to fetch notifications');
  }

  /// =======================
  /// UNREAD COUNT
  /// =======================
  Future<int> getUnreadCount() async {
    final response = await apiClient.get('/notifications/unread-count');

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      return data['unreadCount'] as int;
    }

    throw Exception('Failed to fetch unread count');
  }

  /// =======================
  /// MARK READ
  /// =======================
  Future<void> markRead(String notificationId) async {
    final response = await apiClient.patch('/notifications/$notificationId');

    if (response.statusCode != 202) {
      throw Exception('Failed to mark notification as read');
    }
  }

  /// =======================
  /// READ ALL
  /// =======================
  Future<void> readAll() async {
    final response = await apiClient.patch('/notifications/read-all');

    if (response.statusCode != 202) {
      throw Exception('Failed to mark all notifications as read');
    }
  }

  /// =======================
  /// REGISTER DEVICE
  /// =======================
  Future<void> registerDeviceToken({
    required String token,
    required String platform,
  }) async {
    final response = await apiClient.post(
      '/notifications/device',
      body: {'deviceToken': token, 'platform': platform},
    );

    if (response.statusCode != 201) {
      throw Exception('Failed to register device');
    }
  }
}
