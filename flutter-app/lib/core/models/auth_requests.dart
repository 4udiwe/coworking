class RegisterRequest {
  final String name;
  final String lastName;
  final String email;
  final String password;
  final String roleCode; // student, teacher, admin

  RegisterRequest({
    required this.name,
    required this.lastName,
    required this.email,
    required this.password,
    required this.roleCode,
  });

  Map<String, dynamic> toJson() => {
    'first_name': name,
    'last_name': lastName,
    'email': email,
    'password': password,
    'roleCode': roleCode,
  };
}

class LoginRequest {
  final String email;
  final String password;

  LoginRequest({required this.email, required this.password});

  Map<String, dynamic> toJson() => {'email': email, 'password': password};
}

class RefreshTokenRequest {
  final String refreshToken;

  RefreshTokenRequest({required this.refreshToken});

  Map<String, dynamic> toJson() => {'refreshToken': refreshToken};
}

class RevokeSessionRequest {
  final String sessionId;

  RevokeSessionRequest({required this.sessionId});

  Map<String, dynamic> toJson() => {'sessionId': sessionId};
}
