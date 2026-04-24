import 'package:jwt_decode/jwt_decode.dart';

class UserAccessClaims {
  final String userId;
  final String userName;
  final String email;
  final List<String> roles;

  UserAccessClaims({
    required this.userId,
    required this.userName,
    required this.email,
    required this.roles,
  });
}

class JwtClaimsParser {
  static bool isExpired(String token) {
    return Jwt.isExpired(token);
  }

  static UserAccessClaims parse(String token) {
    final payload = Jwt.parseJwt(token);

    return UserAccessClaims(
      userId: payload['sub']?.toString() ?? '',
      userName: payload['name']?.toString() ?? '',
      email: payload['email']?.toString() ?? '',
      roles: _parseRoles(payload),
    );
  }

  static List<String> _parseRoles(Map<String, dynamic> payload) {
    final rawRoles = payload['roles'];

    if (rawRoles == null) return [];

    if (rawRoles is List) {
      return rawRoles.map((e) => e.toString()).toList();
    }

    // иногда приходит строка
    if (rawRoles is String) {
      return [rawRoles];
    }

    return [];
  }
}