class Coworking {
  final String id;
  final String name;
  final String address;
  final bool isActive;
  final DateTime createdAt;
  final DateTime updatedAt;

  Coworking({
    required this.id,
    required this.name,
    required this.address,
    required this.isActive,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Coworking.fromJson(Map<String, dynamic> json) => Coworking(
        id: json['id'],
        name: json['name'],
        address: json['address'],
        isActive: json['isActive'],
        createdAt: DateTime.parse(json['createdAt']).toLocal(),
        updatedAt: DateTime.parse(json['updatedAt']).toLocal(),
      );

  Map<String, dynamic> toJson() => {
        'id': id,
        'name': name,
        'address': address,
        'isActive': isActive,
        'createdAt': createdAt.toIso8601String(),
        'updatedAt': updatedAt.toIso8601String(),
      };
}
