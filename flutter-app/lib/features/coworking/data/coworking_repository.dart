import 'dart:convert';
import '../../../core/api/api_client.dart';
import '../../../core/models/coworking.dart';
import '../../../core/models/place.dart';
import '../../../core/models/layout.dart';
import '../../../core/models/available_places_response.dart';

class CoworkingRepository {
  final ApiClient apiClient;

  CoworkingRepository({required this.apiClient});

  void _checkStatus(dynamic response, {List<int> validCodes = const [200]}) {
    if (!validCodes.contains(response.statusCode)) {
      throw Exception(
        'Request failed: ${response.statusCode}, body: ${response.body}',
      );
    }
  }

  Map<String, dynamic> _decodeMap(String body) {
    final decoded = jsonDecode(body);
    if (decoded is! Map<String, dynamic>) {
      throw Exception('Expected JSON object but got ${decoded.runtimeType}');
    }
    return decoded;
  }

  Future<List<Coworking>> getCoworkings() async {
    final response = await apiClient.get(
      '/coworkings',
    );
    if (response.statusCode == 200) {
      final jsonMap = jsonDecode(response.body) as Map<String, dynamic>;
      final jsonList = jsonMap['coworkings'] as List<dynamic>;
      return jsonList.map((e) => Coworking.fromJson(e as Map<String, dynamic>)).toList();
    }
    throw Exception('Failed to fetch coworkings: ${response.statusCode}');
  }

  Future<Coworking> getCoworkingById(String id) async {
    final response = await apiClient.get(
      '/coworkings/$id',
    );
    if (response.statusCode == 200) return Coworking.fromJson(jsonDecode(response.body));
    throw Exception('Failed to fetch coworking $id: ${response.statusCode}');
  }

  Future<List<Place>> getPlaces(String coworkingId) async {
    final response = await apiClient.get(
      '/coworkings/$coworkingId/places',
    );
    _checkStatus(response);

    final data = _decodeMap(response.body);

    final list = (data['places'] as List?) ?? [];

    return list.map((e) => Place.fromJson(e)).toList();
  }

  Future<AvailablePlacesResponse> getAvailablePlaces(String coworkingId, DateTime start, DateTime end) async {
    final query = '?startTime=${start.toUtc().toIso8601String()}&endTime=${end.toUtc().toIso8601String()}';
    final response = await apiClient.get(
      '/coworkings/$coworkingId/available-places$query',
    );
    if (response.statusCode == 200) return AvailablePlacesResponse.fromJson(jsonDecode(response.body));
    throw Exception('Failed to fetch available places');
  }

  Future<Layout> getLayout(String coworkingId) async {
    final response = await apiClient.get(
      '/coworkings/$coworkingId/layout',
    );
    if (response.statusCode == 200) return Layout.fromJson(jsonDecode(response.body));
    throw Exception('Failed to fetch layout for coworking $coworkingId');
  }

  Future<void> createBooking(String placeId, DateTime start, DateTime end) async {
    final response = await apiClient.post(
      '/bookings',
      body: {
        'placeId': placeId,
        'startTime': start.toUtc().toIso8601String(),
        'endTime': end.toUtc().toIso8601String(),
      },
    );

    if (response.statusCode == 201) {
      return;
    } else {
      throw Exception('Failed to create booking: ${response.statusCode}, ${response.body}');
    }
  }

}