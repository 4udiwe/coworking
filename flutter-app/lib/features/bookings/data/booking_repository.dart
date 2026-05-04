import 'dart:convert';

import '../../../core/api/api_client.dart';
import '../../../core/models/booking.dart';
import '../../../core/utils/check_response_status.dart';

class BookingRepository {
  final ApiClient apiClient;

  BookingRepository({required this.apiClient});

  Future<PaginatedBookings> getActiveBookings({required int page, required int pageSize}) async {
    final response = await apiClient.get(
        '/bookings/active',
      queryParameters: {
          'page': page,
          'pageSize': pageSize,
      }
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

  Future<PaginatedBookings> getHistoryBookings({
    required int page,
    required int pageSize,
  }) async {
    final response = await apiClient.get(
      '/bookings/history',
      queryParameters: {
        'page': page,
        'pageSize': pageSize,
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


  Future<void> cancelBooking(String id, {String? reason}) async {
    final response = await apiClient.delete(
      '/bookings/$id',
      body: {
        'reason' : reason ?? 'cancelled_by_user',
      },
    );

    if (response.statusCode != 202) {
      throw Exception('Failed to cancel booking $id');
    }
  }
}