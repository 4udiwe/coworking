import 'dart:convert';
import 'package:coworking_app/core/models/booking.dart';
import 'package:coworking_app/core/models/coworking.dart';
import 'package:coworking_app/core/models/layout.dart';
import 'package:coworking_app/core/models/place.dart';
import 'package:coworking_app/core/models/upload_media.dart';
import '../../../../core/api/api_client.dart';
import '../../../core/models/user.dart';
import '../../../core/utils/check_response_status.dart';
import 'package:cross_file/cross_file.dart';

class AdminRepository {
  final ApiClient apiClient;

  AdminRepository({required this.apiClient});

  // 🔧 Универсальные helpers

  Map<String, dynamic> _decodeMap(String body) {
    final decoded = jsonDecode(body);
    if (decoded is! Map<String, dynamic>) {
      throw Exception('Expected JSON object but got ${decoded.runtimeType}');
    }
    return decoded;
  }

  // 📦 Coworkings

  Future<List<Coworking>> getCoworkings() async {
    final response = await apiClient.get('/coworkings');
    checkStatus(response);

    final data = _decodeMap(response.body);
    final list = (data['coworkings'] as List?) ?? [];

    return list.map((e) {
      if (e is! Map<String, dynamic>) {
        throw Exception('Invalid coworking format');
      }
      return Coworking.fromJson(e);
    }).toList();
  }

  Future<Coworking> getCoworkingById(String id) async {
    final response = await apiClient.get('/coworkings/$id');
    checkStatus(response);

    final data = _decodeMap(response.body);
    return Coworking.fromJson(data);
  }

  // 📍 Places

  Future<List<Place>> getPlaces(String coworkingId) async {
    final response = await apiClient.get('/coworkings/$coworkingId/places');
    checkStatus(response);

    final data = _decodeMap(response.body);

    final list = (data['places'] as List?) ?? [];

    return list.map((e) => Place.fromJson(e)).toList();
  }

  Future<void> addPlaces(String coworkingId, List<Place> places) async {
    final response = await apiClient.post(
      '/admin/places',
      body: {
        'coworkingId': coworkingId,
        'places': places.map((p) => p.toJson()).toList(),
      },
    );

    checkStatus(response, validCodes: [201]);
  }

  Future<void> setPlaceActive(String placeId, bool active) async {
    final response = await apiClient.patch(
      '/admin/places/$placeId/set_active',
      body: {'active': active},
    );

    checkStatus(response, validCodes: [202]);
  }

  // 🧩 Layouts

  Future<List<LayoutVersion>> getLayoutVersions(String coworkingId) async {
    final response = await apiClient.get(
      '/admin/coworkings/$coworkingId/layouts',
    );
    checkStatus(response);

    final data = _decodeMap(response.body);
    final list = (data['layoutVersions'] as List?) ?? [];
    return list
        .map((e) => LayoutVersion.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<Layout> getLayoutVersion(String coworkingId, int version) async {
    final response = await apiClient.get(
      '/admin/coworkings/$coworkingId/layouts/$version',
    );
    checkStatus(response);

    final data = _decodeMap(response.body);
    return Layout.fromJson(data);
  }

  Future<Layout> getLatestLayout(String coworkingId) async {
    final response = await apiClient.get('/coworkings/$coworkingId/layout');
    checkStatus(response);

    final data = _decodeMap(response.body);
    return Layout.fromJson(data);
  }

  Future<void> createLayout(
    String coworkingId,
    LayoutSchema layout,
    int version,
  ) async {
    final response = await apiClient.post(
      '/admin/coworkings/$coworkingId/layouts',
      body: {
        'coworkingId': coworkingId,
        'version': version,
        'layout': layout.toJson(),
      },
    );

    checkStatus(response, validCodes: [201]);
  }

  Future<void> deleteLayout(String coworkingId, int layoutVersion) async {
    final response = await apiClient.delete(
      '/admin/coworkings/$coworkingId/layouts/$layoutVersion',
    );

    checkStatus(response, validCodes: [202]);
  }

  Future<void> setActiveLayout(String coworkingId, int layoutVersion) async {
    final response = await apiClient.patch(
      '/admin/coworkings/$coworkingId/layouts/$layoutVersion',
    );

    checkStatus(response, validCodes: [202]);
  }
  // 🏢 Coworking admin

  Future<void> createCoworking({
    required String name,
    required String address,
    required List<String> mediaIDs,
  }) async {
    final response = await apiClient.post(
      '/admin/coworkings',
      body: {'name': name, 'address': address, 'mediaIDs': mediaIDs},
    );

    checkStatus(response, validCodes: [201]);
  }

  Future<void> updateCoworking(String id, String name, String address) async {
    final response = await apiClient.put(
      '/admin/coworkings/$id',
      body: {'name': name, 'address': address},
    );

    checkStatus(response, validCodes: [202]);
  }

  Future<void> setCoworkingActive(String id, bool active) async {
    final response = await apiClient.patch(
      '/admin/coworkings/$id/set_active',
      body: {'active': active},
    );

    checkStatus(response, validCodes: [202]);
  }

  // 📅 Booking

  Future<void> adminCancelBooking(String bookingId, {String? reason}) async {
    final response = await apiClient.delete(
      '/admin/bookings/$bookingId',
      body: {'reason': reason ?? 'cancelled by admin'},
    );

    checkStatus(response, validCodes: [202]);
  }

  Future<PaginatedBookings> getActiveBookings({
    required String coworkingId,
    required int page,
    required int pageSize,
    DateTime? date,
    String? placeType,
    String? sortBy = 'desc',
  }) async {
    print("getActiveBookings page: $page, pageSize: $pageSize");

    DateTime? dateFrom;
    DateTime? dateTo;

    if (date != null) {
      dateFrom = DateTime(date.year, date.month, date.day);
      dateTo = dateFrom.add(const Duration(days: 1));
    }

    final response = await apiClient.get(
      '/admin/bookings',
      queryParameters: {
        'coworkingId': coworkingId,
        'page': page,
        'pageSize': pageSize,
        if (dateFrom != null) 'dateFrom': dateFrom.toUtc().toIso8601String(),
        if (dateTo != null) 'dateTo': dateTo.toUtc().toIso8601String(),
        if (placeType != null && placeType.isNotEmpty) 'placeType': placeType,
        if (sortBy != null) 'sortBy': sortBy, // asc / desc
      },
    );

    checkStatus(response, validCodes: [200]);

    final data = jsonDecode(response.body) as Map<String, dynamic>;
    final bookingsList = (data['bookings'] as List? ?? [])
        .map((e) => Booking.fromJson(e as Map<String, dynamic>))
        .toList();

    final pagination = data['pagination'] as Map<String, dynamic>? ?? {};

    return PaginatedBookings(
      items: bookingsList,
      totalItems: pagination['totalItems'] ?? 0,
      totalPages: pagination['totalPages'] ?? 0,
      page: pagination['page'] ?? page,
      pageSize: pagination['pageSize'] ?? pageSize,
    );
  }

  /// Получить список пользователей с фильтрацией, поиском и пагинацией
  Future<PaginatedUsers> getUsers({
    String? search,
    int page = 1,
    int size = 10,
    String? role,
    bool? isActive,
    String? sort,
  }) async {
    final queryParams = <String, dynamic>{
      if (search != null) 'search': search,
      'page': page,
      'size': size,
      if (role != null) 'role': role,
      if (isActive != null) 'isActive': isActive,
      if (sort != null) 'sort': sort,
    };

    final response = await apiClient.get(
      '/admin/users',
      queryParameters: queryParams,
    );

    checkStatus(response);

    final data = jsonDecode(response.body) as Map<String, dynamic>;

    final usersList = (data['users'] as List? ?? [])
        .map((e) => User.fromJson(e as Map<String, dynamic>))
        .toList();

    return PaginatedUsers(
      items: usersList,
      total: data['total'] ?? 0,
      page: data['page'] ?? page,
      size: data['size'] ?? size,
    );
  }

  /// Получить одного пользователя по ID
  Future<User> getUserById(String userId) async {
    final response = await apiClient.get('/admin/users/$userId');
    checkStatus(response, validCodes: [200]);

    final data = jsonDecode(response.body) as Map<String, dynamic>;
    return User.fromJson(data);
  }

  /// Обновить роли пользователя (полностью заменяет роли)
  Future<void> updateUserRoles(String userId, List<String> roles) async {
    final response = await apiClient.put(
      '/admin/users/$userId/roles',
      body: {'role_codes': roles},
    );

    checkStatus(response, validCodes: [200]);
  }

  /// Активировать или деактивировать пользователя
  Future<void> setUserActive(String userId, bool active) async {
    final response = await apiClient.patch(
      '/admin/users/$userId/set_active',
      body: {'active': active},
    );

    checkStatus(response, validCodes: [200]);
  }
}
