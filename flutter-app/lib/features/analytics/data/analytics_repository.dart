import 'dart:convert';

import 'package:coworking_app/core/models/heatmap_response.dart';
import 'package:coworking_app/core/models/load_response.dart';

import '../../../core/api/api_client.dart';

class AnalyticsRepository {
  final ApiClient apiClient;

  AnalyticsRepository({required this.apiClient});

  // ---------- Hourly ----------

  Future<HourlyLoad> getHourlyLoad(String coworkingId, int? weekday) async {
    final response = await apiClient.get(
      '/analytics/hourly/$coworkingId',
        queryParameters: {
          if (weekday != null) 'weekday': weekday,
        }
    );

    if (response.statusCode == 200) {
      return HourlyLoad.fromJson(jsonDecode(response.body));
    }

    throw Exception('Failed to fetch hourly analytics');
  }

  // ---------- Weekday ----------

  Future<WeekdayLoad> getWeekdayLoad(String coworkingId) async {
    final response = await apiClient.get(
      '/analytics/weekday/$coworkingId'
    );

    if (response.statusCode == 200) {
      return WeekdayLoad.fromJson(jsonDecode(response.body));
    }

    throw Exception('Failed to fetch weekday analytics');
  }

  // ---------- Coworking Heatmap ----------

  Future<CoworkingHeatmap> getCoworkingHeatmap(String coworkingId) async {
    final response = await apiClient.get(
      '/analytics/coworking_heatmap/$coworkingId',
    );

    if (response.statusCode == 200) {
      return CoworkingHeatmap.fromJson(jsonDecode(response.body));
    }

    throw Exception('Failed to fetch coworking heatmap');
  }

  // ---------- Place Heatmap ----------

  Future<PlaceHeatmap> getPlaceHeatmap(String placeId) async {
    final response = await apiClient.get(
      '/analytics/place_heatmap/$placeId',
    );

    if (response.statusCode == 200) {
      return PlaceHeatmap.fromJson(jsonDecode(response.body));
    }

    throw Exception('Failed to fetch place heatmap');
  }
}