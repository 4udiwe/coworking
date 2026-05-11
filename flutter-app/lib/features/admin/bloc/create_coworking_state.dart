import 'package:equatable/equatable.dart';
import 'dart:typed_data';

enum MediaUploadStatus { uploading, success, error }

class MediaItemState extends Equatable {
  final String localId; // временный id до получения mediaId с сервера
  final Uint8List previewBytes;
  final MediaUploadStatus status;
  final String? mediaId; // получаем после успешной загрузки
  final String? error;

  const MediaItemState({
    required this.localId,
    required this.previewBytes,
    required this.status,
    this.mediaId,
    this.error,
  });

  MediaItemState copyWith({
    MediaUploadStatus? status,
    String? mediaId,
    String? error,
  }) {
    return MediaItemState(
      localId: localId,
      previewBytes: previewBytes,
      status: status ?? this.status,
      mediaId: mediaId ?? this.mediaId,
      error: error ?? this.error,
    );
  }

  @override
  List<Object?> get props => [localId, status, mediaId, error];
}

class CreateCoworkingState extends Equatable {
  final List<MediaItemState> mediaItems;
  final bool isSubmitting;
  final String? submitError;
  final bool submitted;

  const CreateCoworkingState({
    this.mediaItems = const [],
    this.isSubmitting = false,
    this.submitError,
    this.submitted = false,
  });

  bool get hasUploading =>
      mediaItems.any((m) => m.status == MediaUploadStatus.uploading);

  bool get hasErrors =>
      mediaItems.any((m) => m.status == MediaUploadStatus.error);

  List<String> get uploadedMediaIds => mediaItems
      .where((m) => m.status == MediaUploadStatus.success && m.mediaId != null)
      .map((m) => m.mediaId!)
      .toList();

  CreateCoworkingState copyWith({
    List<MediaItemState>? mediaItems,
    bool? isSubmitting,
    String? submitError,
    bool? submitted,
    bool clearSubmitError = false,
  }) {
    return CreateCoworkingState(
      mediaItems: mediaItems ?? this.mediaItems,
      isSubmitting: isSubmitting ?? this.isSubmitting,
      submitError: clearSubmitError ? null : (submitError ?? this.submitError),
      submitted: submitted ?? this.submitted,
    );
  }

  @override
  List<Object?> get props => [mediaItems, isSubmitting, submitError];
}
