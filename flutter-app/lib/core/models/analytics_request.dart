class CoworkingAnalyticsRequest {
  final String coworkingId;

  CoworkingAnalyticsRequest({required this.coworkingId});

  Map<String, dynamic> toJson() => {
    'coworkingId': coworkingId,
  };
}

class PlaceAnalyticsRequest {
  final String placeId;

  PlaceAnalyticsRequest({required this.placeId});

  Map<String, dynamic> toJson() => {
    'placeId': placeId,
  };
}