abstract class NotificationEvent {}

class FetchNotifications extends NotificationEvent {
  final int limit;
  final int offset;
  final bool refresh;

  FetchNotifications({this.limit = 20, this.offset = 0, this.refresh = false});
}

class MarkAsRead extends NotificationEvent {
  final String notificationId;
  MarkAsRead(this.notificationId);
}

class MarkAllRead extends NotificationEvent {}

/// 🔁 polling
class PollNotifications extends NotificationEvent {}

/// 🔔 unread count
class FetchUnreadCount extends NotificationEvent {}

/// ✅ mark read
class MarkNotificationRead extends NotificationEvent {
  final String id;
  MarkNotificationRead(this.id);
}

/// ✅ read all
class MarkAllNotificationsRead extends NotificationEvent {}

/// ❌ очистка snackbar
class ClearAction extends NotificationEvent {}