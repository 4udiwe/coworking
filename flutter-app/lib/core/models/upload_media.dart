class UploadedMedia {
  final String id;
  final String status;
  final Map<String, String> urls;

  UploadedMedia({required this.id, required this.status, required this.urls});

  factory UploadedMedia.fromJson(Map<String, dynamic> json) {
    return UploadedMedia(
      id: json['id'],
      status: json['status'],
      urls: Map<String, String>.from(json['urls'] ?? {}),
    );
  }
}
