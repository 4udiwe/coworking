class UserSession {
  final String id;
  final String userId;
  final String userAgent;
  final String device;
  final bool revoked;
  final bool current;
  final DateTime createdAt;
  final DateTime expiresAt;
  final String ipAddress;
  final DateTime lastUsedAt;

  UserSession({
    required this.id,
    required this.userId,
    required this.userAgent,
    required this.device,
    required this.revoked,
    required this.current,
    required this.createdAt,
    required this.expiresAt,
    required this.ipAddress,
    required this.lastUsedAt,
  });

  factory UserSession.fromJson(Map<String, dynamic> json) => UserSession(
    id: json['id'],
    userId: json['userId'],
    userAgent: json['userAgent'],
    device: json['device'],
    revoked: json['revoked'],
    current: json['current'],
    createdAt: DateTime.parse(json['createdAt']),
    expiresAt: DateTime.parse(json['expiresAt']),
    ipAddress: json['ipAddress'],
    lastUsedAt: DateTime.parse(json['lastUsedAt']),
  );
}