import 'dart:typed_data';

import 'package:coworking_app/features/admin/bloc/admin_bloc.dart';
import 'package:coworking_app/features/admin/bloc/admin_event.dart' as admin;
import 'package:coworking_app/features/admin/data/media_reposotory.dart';
import 'package:cross_file/cross_file.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import 'create_coworking_event.dart';
import 'create_coworking_state.dart';

class CreateCoworkingBloc
    extends Bloc<CreateCoworkingEvent, CreateCoworkingState> {
  final MediaRepository _mediaRepo;
  final AdminBloc _adminBloc; // только для передачи финального события

  CreateCoworkingBloc({
    required MediaRepository mediaRepo,
    required AdminBloc adminBloc,
  }) : _mediaRepo = mediaRepo,
       _adminBloc = adminBloc,
       super(const CreateCoworkingState()) {
    on<AddMediaFilesEvent>(_onAddMediaFiles);
    on<RemoveMediaItemEvent>(_onRemoveMediaItem);
    on<SubmitCreateCoworkingEvent>(_onSubmit);
  }

  Future<void> _onAddMediaFiles(
    AddMediaFilesEvent event,
    Emitter<CreateCoworkingState> emit,
  ) async {
    print("Add media file event");
    // Генерируем уникальный localId для каждого файла отдельно
    final indexed = event.files
        .map(
          (f) => (
            localId:
                '${f.name}_${DateTime.now().microsecondsSinceEpoch}_${f.bytes.hashCode}',
            file: f,
          ),
        )
        .toList();
    // Добавляем все файлы сразу в состояние со статусом uploading
    final newItems = indexed
        .map(
          (entry) => MediaItemState(
            localId: entry.localId,
            previewBytes: entry.file.bytes,
            status: MediaUploadStatus.uploading,
          ),
        )
        .toList();

    emit(state.copyWith(mediaItems: [...state.mediaItems, ...newItems]));

    // Каждый item теперь точно связан со своим файлом
    await Future.wait(
      indexed.map((entry) => _uploadSingle(entry.localId, entry.file)),
    );
  }

  // Загружаем один файл и обновляем его статус в состоянии
  Future<void> _uploadSingle(
    String localId,
    ({Uint8List bytes, String name}) file,
  ) async {
    print("Upload single");
    try {
      final media = await _mediaRepo.uploadMedia(
        XFile.fromData(file.bytes, name: file.name),
      );

      final updated = state.mediaItems.map((m) {
        if (m.localId != localId) return m;
        return m.copyWith(status: MediaUploadStatus.success, mediaId: media.id);
      }).toList();

      emit(state.copyWith(mediaItems: updated));
    } catch (e) {
      print("upload error: ${e.toString()}");
      final updated = state.mediaItems.map((m) {
        if (m.localId != localId) return m;
        return m.copyWith(status: MediaUploadStatus.error, error: e.toString());
      }).toList();

      emit(state.copyWith(mediaItems: updated));
    }
  }

  Future<void> _onRemoveMediaItem(
    RemoveMediaItemEvent event,
    Emitter<CreateCoworkingState> emit,
  ) async {
    final item = state.mediaItems.firstWhere(
      (m) => m.localId == event.localId,
      orElse: () => throw StateError('Item not found'),
    );

    // Убираем из состояния сразу — не ждём удаления с сервера
    emit(
      state.copyWith(
        mediaItems: state.mediaItems
            .where((m) => m.localId != event.localId)
            .toList(),
      ),
    );

    // Удаляем с сервера в фоне, только если уже загружен
    if (item.mediaId != null) {
      try {
        await _mediaRepo.deleteMedia(item.mediaId!);
      } catch (_) {
        // Некритично — можно добавить логирование
      }
    }
  }

  Future<void> _onSubmit(
    SubmitCreateCoworkingEvent event,
    Emitter<CreateCoworkingState> emit,
  ) async {
    if (state.hasUploading || state.hasErrors) return;

    emit(state.copyWith(isSubmitting: true));

    // Передаём событие в AdminBloc — он не знает ничего про медиа загрузку
    _adminBloc.add(
      admin.CreateCoworkingEvent(
        event.name,
        event.address,
        state.uploadedMediaIds,
      ),
    );

    emit(state.copyWith(submitted: true, isSubmitting: false));
  }

  @override
  Future<void> close() {
    // При закрытии удаляем с сервера всё загруженное, если форма не была отправлена
    // (isSubmitting == false означает что submit не был вызван)
    for (final item in state.mediaItems) {
      if (item.mediaId != null && !state.submitted) {
        _mediaRepo.deleteMedia(item.mediaId!).ignore();
      }
    }
    return super.close();
  }
}
