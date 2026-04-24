import 'dart:convert';
import 'package:coworking_app/core/storage/token_storage.dart';
import 'package:http/http.dart' as http;

import '../../../core/api/api_client.dart';
import '../../../core/models/auth_tokens.dart';
import '../../../core/models/auth_requests.dart';
import '../../../core/models/user.dart';
import '../../../core/models/user_session.dart';

class AuthRepository {
  final ApiClient apiClient;
  final TokenStorage tokenStorage;

  AuthRepository({required this.apiClient, required this.tokenStorage});

  // -------------------- Auth --------------------

  Future<AuthTokens> register(RegisterRequest request) async {
    final response = await apiClient.post(
      '/auth/register',
      body: request.toJson(),
    );

    if (response.statusCode == 201) {
      return AuthTokens.fromJson(jsonDecode(response.body));
    } else if (response.statusCode == 409) {
      throw Exception('User already exists');
    } else {
      throw Exception('Registration failed: ${response.statusCode}');
    }
  }

  Future<AuthTokens> login(LoginRequest request) async {
    final response = await apiClient.post(
      '/auth/login',
      body: request.toJson(),
    );

    if (response.statusCode == 200) {
      return AuthTokens.fromJson(jsonDecode(response.body));
    } else if (response.statusCode == 401) {
      throw Exception('Invalid credentials');
    } else {
      throw Exception('Login failed: ${response.statusCode}');
    }
  }

  Future<AuthTokens> refresh() async {
    final tokens = await tokenStorage.readTokens();

    if (tokens == null) {
      throw Exception('No refresh token');
    }

    final request = RefreshTokenRequest(refreshToken: tokens.refreshToken);

    final response = await http.post(
      Uri.parse('${apiClient.baseUrl}/auth/refresh'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode(request.toJson()),
    );

    if (response.statusCode == 200) {
      return AuthTokens.fromJson(jsonDecode(response.body));
    }
    throw Exception('Refresh token failed: ${response.statusCode}');
  }

  Future<void> logout(String refreshToken) async {
    final request = RefreshTokenRequest(refreshToken: refreshToken);
    final response = await apiClient.post(
      '/auth/logout',
      body: request.toJson(),
    );

    if (response.statusCode != 204) {
      throw Exception('Logout failed: ${response.statusCode}');
    }
  }

  // -------------------- Users --------------------

  Future<User> getProfile() async {
    final response = await apiClient.get('/users/me');

    if (response.statusCode == 200) {
      return User.fromJson(jsonDecode(response.body));
    } else if (response.statusCode == 401) {
      throw Exception('Unauthorized');
    } else {
      throw Exception('Failed to fetch profile: ${response.statusCode}');
    }
  }

  Future<List<UserSession>> getActiveSessions() async {
    final response = await apiClient.get('/users/sessions/active');

    if (response.statusCode == 200) {
      final list = jsonDecode(response.body) as List;
      return list.map((e) => UserSession.fromJson(e)).toList();
    } else {
      throw Exception(
        'Failed to fetch active sessions: ${response.statusCode}',
      );
    }
  }

  Future<List<UserSession>> getAllSessions() async {
    final response = await apiClient.get('/users/sessions/all');

    if (response.statusCode == 200) {
      final list = jsonDecode(response.body) as List;
      return list.map((e) => UserSession.fromJson(e)).toList();
    } else {
      throw Exception('Failed to fetch all sessions: ${response.statusCode}');
    }
  }

  Future<void> revokeSession(String sessionId) async {
    final request = RevokeSessionRequest(sessionId: sessionId);
    final response = await apiClient.post(
      '/users/sessions/revoke',
      body: request,
    );

    if (response.statusCode != 200) {
      throw Exception('Failed to revoke session: ${response.statusCode}');
    }
  }
}
