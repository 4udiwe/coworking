import 'dart:ui';

import 'package:http/http.dart' as http;
import '../storage/token_storage.dart';
import '../models/auth_tokens.dart';

typedef RequestBuilder = Future<http.Response> Function();

class AuthInterceptor {
  final TokenStorage tokenStorage;
  final VoidCallback onSessionExpired;
  Future<AuthTokens>? _refreshing;

  AuthInterceptor({required this.tokenStorage, required this.onSessionExpired});

  Future<http.Response> execute(
    RequestBuilder request, {
    required Future<AuthTokens> Function(String) onRefresh,
  }) async {
    var response = await request();

    if (response.statusCode != 401) return response;

    final tokens = await tokenStorage.readTokens();
    if (tokens == null) {
      onSessionExpired();
      throw Exception('No tokens available');
    }

    try {
      // Дедупликация параллельных refresh-запросов
      _refreshing ??= onRefresh(tokens.refreshToken)
          .then((newTokens) async {
            await tokenStorage.saveTokens(newTokens);
            return newTokens;
          })
          .whenComplete(() => _refreshing = null);

      await _refreshing;

      // Повторяем оригинальный запрос
      response = await request();

      if (response.statusCode == 401) {
        // refresh прошёл, но запрос всё равно 401 — сессия мертва
        await tokenStorage.clearTokens();
        onSessionExpired();
        throw Exception('Session expired after refresh');
      }

      return response;
    } catch (e) {
      _refreshing = null;
      await tokenStorage.clearTokens();
      onSessionExpired();
      rethrow;
    }
  }
}
