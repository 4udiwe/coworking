class User {
  final String id;
  final String firstName;
  final String lastName;
  final String email;
  final bool isActive;
  final DateTime createdAt;
  final DateTime updatedAt;
  final List<UserRole> roles;

  User({
    required this.id,
    required this.firstName,
    required this.lastName,
    required this.email,
    required this.isActive,
    required this.createdAt,
    required this.updatedAt,
    required this.roles,
  });

  String get fullName => '$firstName $lastName';

  bool get isAdmin => roles.any((role) => role.roleCode == UserRole.admin);

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'] ?? '',
      firstName: json['first_name'] ?? '',
      lastName: json['last_name'] ?? '',
      email: json['email'] ?? '',
      isActive: json['isActive'] ?? false,
      createdAt: json['createdAt'] != null
          ? DateTime.parse(json['createdAt'])
          : DateTime.now(),
      updatedAt: json['updatedAt'] != null
          ? DateTime.parse(json['updatedAt'])
          : DateTime.now(),
      roles: (json['roles'] as List? ?? [])
          .map((role) => UserRole.fromJson(role))
          .toList(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'first_name': firstName,
      'last_name': lastName,
      'email': email,
      'isActive': isActive,
      'createdAt': createdAt.toIso8601String(),
      'updatedAt': updatedAt.toIso8601String(),
      'roles': roles.map((role) => role.toJson()).toList(),
    };
  }
}

class UserRole {
  static const String admin = 'admin';
  static const String student = 'student';
  static const String teacher = 'teacher';


  final String id;
  final String roleCode;
  final String name;

  UserRole({
    required this.id,
    required this.roleCode,
    required this.name,
  });

  factory UserRole.fromJson(Map<String, dynamic> json) {
    return UserRole(
      id: json['id'] ?? '',
      roleCode: json['roleCode'] ?? '',
      name: json['name'] ?? '',
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'roleCode': roleCode,
      'name': name,
    };
  }
}

class PaginatedUsers {
  final List<User> items;
  final int total;
  final int page;
  final int size;

  PaginatedUsers({
    required this.items,
    required this.total,
    required this.page,
    required this.size,
  });
}
