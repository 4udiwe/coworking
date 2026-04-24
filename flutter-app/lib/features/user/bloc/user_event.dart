abstract class UserEvent {}

/// 🔹 профиль
class LoadProfile extends UserEvent {}

/// 🔹 активные сессии
class LoadActiveSessions extends UserEvent {}

/// 🔹 все сессии (если нужен экран)
class LoadAllSessions extends UserEvent {}

/// 🔹 revoke
class RevokeSessionEvent extends UserEvent {
  final String sessionId;

  RevokeSessionEvent(this.sessionId);
}

/// 🔹 refresh всего (удобно для pull-to-refresh)
class RefreshUserData extends UserEvent {}