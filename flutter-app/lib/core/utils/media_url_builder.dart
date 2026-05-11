import 'package:flutter/foundation.dart';

enum MediaSize { thumbnail, medium, large }

class MediaUrlBuilder {
  static const _baseUrl = 'http://localhost:8080';

  static String build(String mediaId, {MediaSize? size}) {
    final resolvedSize = size ?? defaultSize();

    return '$_baseUrl/media/$mediaId/${resolvedSize.name}.webp';
  }

  static String thumbnail(String mediaId) {
    return build(mediaId, size: MediaSize.thumbnail);
  }

  static MediaSize defaultSize() {
    if (kIsWeb) {
      return MediaSize.large;
    }

    return MediaSize.medium;
  }
}
