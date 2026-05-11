import 'dart:convert';
import 'dart:ui';

import 'package:http/http.dart' as http;
import '../storage/token_storage.dart';
import 'auth_interceptor.dart';
import '../models/auth_tokens.dart';
import 'package:http_parser/http_parser.dart';
import 'package:cross_file/cross_file.dart';

class ApiClient {
  final String baseUrl;
  final TokenStorage tokenStorage;
  final AuthInterceptor _interceptor;

  Future<AuthTokens> Function(String refreshToken) onRefresh;

  ApiClient({
    required this.baseUrl,
    required this.tokenStorage,
    required this.onRefresh,
    required VoidCallback onSessionExpired,
  }) : _interceptor = AuthInterceptor(
         tokenStorage: tokenStorage,
         onSessionExpired: onSessionExpired,
       );

  Future<http.Response> get(
    String path, {
    Map<String, String>? headers,
    Map<String, dynamic>? queryParameters,
  }) async {
    var uri = Uri.parse('$baseUrl$path');

    if (queryParameters != null && queryParameters.isNotEmpty) {
      uri = uri.replace(
        queryParameters: queryParameters.map(
          (key, value) => MapEntry(key, value?.toString()),
        ),
      );
    }

    return _interceptor.execute(() async {
      final combinedHeaders = await _headers(headers); // ← каждый раз заново
      return http.get(uri, headers: combinedHeaders);
    }, onRefresh: onRefresh);
  }

  Future<http.Response> post(
    String path, {
    Map<String, String>? headers,
    Object? body,
  }) async {
    return _interceptor.execute(() async {
      final combinedHeaders = await _headers(headers);
      return http.post(
        Uri.parse('$baseUrl$path'),
        headers: combinedHeaders,
        body: body != null ? jsonEncode(body) : null,
      );
    }, onRefresh: onRefresh);
  }

  Future<http.Response> put(
    String path, {
    Map<String, String>? headers,
    Object? body,
  }) async {
    return _interceptor.execute(() async {
      final combinedHeaders = await _headers(headers);
      return http.put(
        Uri.parse('$baseUrl$path'),
        headers: combinedHeaders,
        body: body != null ? jsonEncode(body) : null,
      );
    }, onRefresh: onRefresh);
  }

  Future<http.Response> patch(
    String path, {
    Map<String, String>? headers,
    Object? body,
  }) async {
    return _interceptor.execute(() async {
      final combinedHeaders = await _headers(headers);
      return http.patch(
        Uri.parse('$baseUrl$path'),
        headers: combinedHeaders,
        body: body != null ? jsonEncode(body) : null,
      );
    }, onRefresh: onRefresh);
  }

  Future<http.Response> delete(
    String path, {
    Map<String, String>? headers,
    Object? body,
  }) async {
    return _interceptor.execute(() async {
      final combinedHeaders = await _headers(headers);
      return http.delete(
        Uri.parse('$baseUrl$path'),
        headers: combinedHeaders,
        body: body != null ? jsonEncode(body) : null,
      );
    }, onRefresh: onRefresh);
  }

  Future<http.Response> multipart(
    String path, {
    required XFile file,
    String field = 'file',
    Map<String, String>? headers,
  }) async {
    return _interceptor.execute(() async {
      final uri = Uri.parse('$baseUrl$path');

      final request = http.MultipartRequest('POST', uri);

      final authHeaders = await _headers(headers, json: false);

      request.headers.addAll(authHeaders);

      final bytes = await file.readAsBytes();

      request.files.add(
        http.MultipartFile.fromBytes(
          field,
          bytes,
          filename: file.name,
          contentType: MediaType('image', 'webp'),
        ),
      );

      final streamed = await request.send();

      return http.Response.fromStream(streamed);
    }, onRefresh: onRefresh);
  }

  Future<Map<String, String>> _headers(
    Map<String, String>? headers, {
    bool json = true,
  }) async {
    final tokens = await tokenStorage.readTokens();

    final Map<String, String> authHeader = {};

    if (tokens != null && tokens.accessToken.isNotEmpty) {
      authHeader['Authorization'] = 'Bearer ${tokens.accessToken}';
    }

    return {
      ...?headers,
      ...authHeader,
      if (json) 'Content-Type': 'application/json',
    };
  }
}
