import 'dart:convert';

import '../../../core/api/api_client.dart';
import '../../../core/models/auth_requests.dart';
import '../../../core/models/user.dart';
import '../../../core/models/user_session.dart';

class UserRepository {
  final ApiClient apiClient;

  UserRepository({required this.apiClient});

  Future<User> getProfile() async {
    final response = await apiClient.get(
      '/users/me',
    );

    if (response.statusCode == 200) {
      return User.fromJson(jsonDecode(response.body));
    } else if (response.statusCode == 401) {
      throw Exception('Unauthorized');
    } else {
      throw Exception('Failed to fetch profile: ${response.statusCode}');
    }
  }

  Future<List<UserSession>> getActiveSessions() async {
    final response = await apiClient.get(
      '/users/sessions/active',
    );

    if (response.statusCode == 200) {
      final list = jsonDecode(response.body) as List;
      return list.map((e) => UserSession.fromJson(e)).toList();
    } else {
      throw Exception('Failed to fetch active sessions: ${response.statusCode}');
    }
  }

  Future<List<UserSession>> getAllSessions() async {
    final response = await apiClient.get(
      '/users/sessions/all',
    );

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

    if (response.statusCode != 202) {
      throw Exception('Failed to revoke session: ${response.statusCode}');
    }
  }
}

