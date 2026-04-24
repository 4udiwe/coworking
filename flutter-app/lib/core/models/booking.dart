import 'place.dart';

enum BookingStatus {
  active,
  cancelled,
  completed;

  static BookingStatus fromString(String value) {
    return BookingStatus.values.firstWhere(
      (e) => e.name == value.toLowerCase(),
      orElse: () => BookingStatus.active,
    );
  }
}

class Booking {
  final String id;
  final String userId;
  final String userName;
  final Place place;
  final DateTime startTime;
  final DateTime endTime;
  final BookingStatus status;
  final String? cancelReason;
  final DateTime createdAt;
  final DateTime updatedAt;
  final DateTime? cancelledAt;

  Booking({
    required this.id,
    this.userName = '',
    required this.userId,
    required this.place,
    required this.startTime,
    required this.endTime,
    required this.status,
    this.cancelReason,
    required this.createdAt,
    required this.updatedAt,
    this.cancelledAt,
  });

  factory Booking.fromJson(Map<String, dynamic> json) => Booking(
    id: json['id'] ?? '',
    userId: json['userId'] ?? '',
    userName: json['user_name'] ?? json['userName'] ?? '',
    place: Place.fromJson(json['place']),
    startTime: DateTime.parse(json['startTime']),
    endTime: DateTime.parse(json['endTime']),
    status: BookingStatus.fromString(json['status'] ?? ''),
    cancelReason: json['cancelReason'] as String?,
    createdAt: DateTime.parse(json['createdAt']),
    updatedAt: DateTime.parse(json['updatedAt']),
    cancelledAt: json['cancelledAt'] != null ? DateTime.parse(json['cancelledAt']) : null,
  );

  Map<String, dynamic> toJson() => {
    'id': id,
    'userId': userId,
    'user_name': userName,
    'place': place.toJson(),
    'startTime': startTime.toIso8601String(),
    'endTime': endTime.toIso8601String(),
    'status': status.name,
    'cancelReason': cancelReason,
    'createdAt': createdAt.toIso8601String(),
    'updatedAt': updatedAt.toIso8601String(),
    'cancelledAt': cancelledAt?.toIso8601String(),
  };
}

class PaginatedBookings {
  final List<Booking> items;
  final int totalItems;
  final int totalPages;
  final int page;
  final int pageSize;

  PaginatedBookings({
    required this.items,
    required this.totalItems,
    required this.totalPages,
    required this.page,
    required this.pageSize,
  });
}
