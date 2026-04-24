import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../models/auth_tokens.dart';

class TokenStorage {
  final _storage = FlutterSecureStorage();
  final _accessKey = 'access_token';
  final _refreshKey = 'refresh_token';

  Future<void> saveTokens(AuthTokens tokens) async {
    await _storage.write(key: _accessKey, value: tokens.accessToken);
    await _storage.write(key: _refreshKey, value: tokens.refreshToken);
  }

  Future<AuthTokens?> readTokens() async {
    final access = await _storage.read(key: _accessKey);
    final refresh = await _storage.read(key: _refreshKey);
    if (access != null && refresh != null) {
      return AuthTokens(accessToken: access, refreshToken: refresh, expiresIn: 900);
    }
    return null;
  }

  Future<void> clearTokens() async {
    await _storage.delete(key: _accessKey);
    await _storage.delete(key: _refreshKey);
  }
}