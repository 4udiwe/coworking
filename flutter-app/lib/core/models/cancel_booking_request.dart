class CreateBookingRequest {
  final String placeId;
  final DateTime startTime;
  final DateTime endTime;

  CreateBookingRequest({
    required this.placeId,
    required this.startTime,
    required this.endTime,
  });

  Map<String, dynamic> toJson() => {
    'placeId': placeId,
    'startTime': startTime.toIso8601String(),
    'endTime': endTime.toIso8601String(),
  };
}