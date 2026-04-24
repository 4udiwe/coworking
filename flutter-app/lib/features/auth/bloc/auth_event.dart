abstract class AuthEvent {}

class AuthRegister extends AuthEvent {
  final String name,lastName, email, password, role;
  AuthRegister({required this.name, required this.lastName, required this.email, required this.password, required this.role});
}

class AuthLogin extends AuthEvent {
  final String email, password;
  AuthLogin({required this.email, required this.password});
}

// Пользователь сам нажал "Выйти"
class AuthLogout extends AuthEvent {}

// Сессия была убита извне (другое устройство, сервер)
class AuthSessionExpired extends AuthEvent {}

class AuthCheckSession extends AuthEvent {}

class AuthRefresh extends AuthEvent {}