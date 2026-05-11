import 'package:coworking_app/core/api/api_client.dart';
import 'package:coworking_app/core/models/upload_media.dart';
import 'package:coworking_app/core/utils/check_response_status.dart';
import 'package:cross_file/cross_file.dart';
import 'dart:convert';

class MediaRepository {
  final ApiClient apiClient;

  MediaRepository({required this.apiClient});

  // Загрузка изображений
  Future<UploadedMedia> uploadMedia(XFile file) async {
    final response = await apiClient.multipart(
      '/admin/media/upload',
      file: file,
    );

    checkStatus(response, validCodes: [201]);

    final data = jsonDecode(response.body);

    if (data is! Map<String, dynamic>) {
      throw Exception('Invalid upload response');
    }

    return UploadedMedia.fromJson(data);
  }

  /// Удаляет медиа по id
  Future<void> deleteMedia(String mediaId) async {
    final response = await apiClient.delete('/admin/media/delete/$mediaId');
    // 204 = success
    if (response.statusCode != 204) {
      throw Exception('Delete failed: ${response.statusCode}');
    }
  }
}
