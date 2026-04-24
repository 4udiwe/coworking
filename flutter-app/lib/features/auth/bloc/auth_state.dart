import '../../../core/models/user.dart';
import '../data/jwt_parser.dart';

enum AuthStatus {
  initial,
  loading,
  authenticated,
  unauthenticated,
  failure,
  sessionExpired,
}

class AuthState {
  final AuthStatus status;
  final UserAccessClaims? userClaims;
  //final LoadState<User>? user;
  final String? error;

  const AuthState({
    this.status = AuthStatus.initial,
    this.userClaims,
    //this.user = const LoadState(),
    this.error,
  });

  bool get isAuthenticated => status == AuthStatus.authenticated;

  bool get isAdmin => userClaims?.roles.contains(UserRole.admin) ?? false;

  AuthState copyWith({
    AuthStatus? status,
    UserAccessClaims? userClaims,
    //LoadState<User>? user,
    String? error,
  }) {
    return AuthState(
      status: status ?? this.status,
      userClaims: userClaims ?? this.userClaims,
      //user: user,
      error: error,
    );
  }
}
