import '../../../core/utils/bloc_load_state.dart';
import '../model/notification.dart';

class NotificationState {
  final LoadState<List<NotificationModel>> notifications;

  final int unreadCount;

  final DateTime? lastFetchedAt;

  final String? actionMessage;
  final String? messageId;
  final bool isError;

  DateTime? get latestCreatedAt =>
      notifications.data!.isNotEmpty ? notifications.data!.first.createdAt : null;

  const NotificationState({
    this.notifications = const LoadState(),
    this.unreadCount = 0,
    this.lastFetchedAt,
    this.actionMessage,
    this.messageId,
    this.isError = false,
  });

  NotificationState copyWith({
    LoadState<List<NotificationModel>>? notifications,
    int? unreadCount,
    DateTime? Function()? lastFetchedAt,
    String? Function()? actionMessage,
    String? messageId,
    bool? isError,
  }) {
    return NotificationState(
      notifications: notifications ?? this.notifications,
      unreadCount: unreadCount ?? this.unreadCount,
      lastFetchedAt:
      lastFetchedAt != null ? lastFetchedAt() : this.lastFetchedAt,
      actionMessage:
      actionMessage != null ? actionMessage() : this.actionMessage,
      messageId: messageId ?? this.messageId,
      isError: isError ?? this.isError,
    );
  }
}