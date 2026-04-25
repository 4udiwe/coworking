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

class PollNotifications extends NotificationEvent {}

class FetchUnreadCount extends NotificationEvent {}

class MarkNotificationRead extends NotificationEvent {
  final String id;
  MarkNotificationRead(this.id);
}

class MarkAllNotificationsRead extends NotificationEvent {}

class ClearAction extends NotificationEvent {}

class FCMMessageReceived extends NotificationEvent {
  final String? notificationId;
  FCMMessageReceived({this.notificationId});
}
