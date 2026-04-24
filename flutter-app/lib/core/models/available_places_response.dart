import 'package:coworking_app/core/models/place.dart';

class AvailablePlacesResponse {
  final List<Place> places;

  AvailablePlacesResponse({required this.places});

  factory AvailablePlacesResponse.fromJson(Map<String, dynamic> json) {
    final list = (json['places'] as List).map((e) => Place.fromJson(e)).toList();
    return AvailablePlacesResponse(places: list);
  }
}