import 'dart:typed_data';

import 'package:coworking_app/core/di/service_locator.dart';
import 'package:coworking_app/features/admin/bloc/admin_bloc.dart';
import 'package:coworking_app/features/admin/bloc/create_coworking_bloc.dart';
import 'package:coworking_app/features/admin/bloc/create_coworking_event.dart';
import 'package:coworking_app/features/admin/bloc/create_coworking_state.dart';
import 'package:coworking_app/features/admin/data/media_reposotory.dart';
import 'package:cross_file/cross_file.dart';
import 'package:desktop_drop/desktop_drop.dart';
import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';


class CreateCoworkingDialog extends StatelessWidget {
  const CreateCoworkingDialog({super.key});

  static Future<void> show(BuildContext context) {
    return showDialog(
      context: context,
      barrierDismissible: false,
      builder: (_) => BlocProvider(
        // Bloc живёт ровно столько, сколько открыт диалог
        create: (_) => CreateCoworkingBloc(
          mediaRepo: sl<MediaRepository>(),
          adminBloc: context.read<AdminBloc>(),
        ),
        child: const CreateCoworkingDialog(),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return const Dialog(
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.all(Radius.circular(16)),
      ),
      child: _DialogContent(),
    );
  }
}

// ──────────────────────────────────────────────
// Внутренний stateful виджет с формой
// ──────────────────────────────────────────────

class _DialogContent extends StatefulWidget {
  const _DialogContent();

  @override
  State<_DialogContent> createState() => _DialogContentState();
}

class _DialogContentState extends State<_DialogContent> {
  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();
  final _addressController = TextEditingController();
  bool _isDragOver = false;

  static const _allowedExtensions = ['jpg', 'jpeg', 'png', 'webp', 'gif'];
  static const _maxImages = 10;

  @override
  void dispose() {
    _nameController.dispose();
    _addressController.dispose();
    super.dispose();
  }

  // ──────────────────────────────────────────────
  // Файлы
  // ──────────────────────────────────────────────

  Future<void> _pickFiles() async {
    final result = await FilePicker.platform.pickFiles(
      type: FileType.image,
      allowMultiple: true,
      withData: true,
    );
    if (result == null) return;
    await _addPlatformFiles(result.files);
  }

  Future<void> _addPlatformFiles(List<PlatformFile> files) async {
    print("adding platrofm files");
    final bloc = context.read<CreateCoworkingBloc>();
    final remaining = _maxImages - bloc.state.mediaItems.length;
    if (remaining <= 0) return;

    final valid = files
        .where((f) => f.bytes != null && _isImage(f.name))
        .take(remaining)
        .map((f) => (bytes: f.bytes!, name: f.name))
        .toList();

    if (valid.isEmpty) return;
    bloc.add(AddMediaFilesEvent(valid));
  }

  Future<void> _addXFiles(List<XFile> xFiles) async {
    print("adding xfiles");
    final bloc = context.read<CreateCoworkingBloc>();
    final remaining = _maxImages - bloc.state.mediaItems.length;
    if (remaining <= 0) return;

    final validEntries = <({Uint8List bytes, String name})>[];
    for (final xf in xFiles.take(remaining)) {
      if (!_isImage(xf.name)) continue;
      final bytes = await xf.readAsBytes();
      validEntries.add((bytes: bytes, name: xf.name));
    }

    if (validEntries.isEmpty) return;
    bloc.add(AddMediaFilesEvent(validEntries));
  }

  bool _isImage(String name) {
    final ext = name.toLowerCase().split('.').last;
    return _allowedExtensions.contains(ext);
  }

  void _submit() {
    if (!_formKey.currentState!.validate()) return;
    context.read<CreateCoworkingBloc>().add(
      SubmitCreateCoworkingEvent(
        _nameController.text.trim(),
        _addressController.text.trim(),
      ),
    );
    Navigator.of(context).pop();
  }

  // ──────────────────────────────────────────────
  // Build
  // ──────────────────────────────────────────────

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return BlocBuilder<CreateCoworkingBloc, CreateCoworkingState>(
      builder: (context, state) {
        final atLimit = state.mediaItems.length >= _maxImages;

        return ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 560),
          child: Padding(
            padding: const EdgeInsets.all(24),
            child: Form(
              key: _formKey,
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // ── Header ──
                  _Header(onClose: () => Navigator.of(context).pop()),
                  const SizedBox(height: 20),

                  // ── Name ──
                  TextFormField(
                    controller: _nameController,
                    decoration: const InputDecoration(
                      labelText: 'Name',
                      border: OutlineInputBorder(),
                      prefixIcon: Icon(Icons.business),
                    ),
                    validator: (v) =>
                        (v == null || v.trim().isEmpty) ? 'Required' : null,
                  ),
                  const SizedBox(height: 12),

                  // ── Address ──
                  TextFormField(
                    controller: _addressController,
                    decoration: const InputDecoration(
                      labelText: 'Address',
                      border: OutlineInputBorder(),
                      prefixIcon: Icon(Icons.location_on),
                    ),
                    validator: (v) =>
                        (v == null || v.trim().isEmpty) ? 'Required' : null,
                  ),
                  const SizedBox(height: 20),

                  // ── Images label ──
                  Row(
                    children: [
                      Text('Images', style: theme.textTheme.labelLarge),
                      const SizedBox(width: 8),
                      Text(
                        '${state.mediaItems.length}/$_maxImages',
                        style: theme.textTheme.labelSmall?.copyWith(
                          color: theme.colorScheme.onSurfaceVariant,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 8),

                  // ── Drag & Drop zone ──
                  if (!atLimit)
                    _DropZone(
                      isDragOver: _isDragOver,
                      onDragEntered: () => setState(() => _isDragOver = true),
                      onDragExited: () => setState(() => _isDragOver = false),
                      onDrop: (xFiles) async {
                        setState(() => _isDragOver = false);
                        await _addXFiles(xFiles);
                      },
                      onTap: _pickFiles,
                    ),

                  // ── Thumbnails ──
                  if (state.mediaItems.isNotEmpty) ...[
                    const SizedBox(height: 12),
                    _ThumbnailList(
                      items: state.mediaItems,
                      onRemove: (localId) => context
                          .read<CreateCoworkingBloc>()
                          .add(RemoveMediaItemEvent(localId)),
                    ),
                  ],

                  // ── Upload error hint ──
                  if (state.hasErrors)
                    Padding(
                      padding: const EdgeInsets.only(top: 8),
                      child: Text(
                        'Some images failed to upload. Remove them and try again.',
                        style: theme.textTheme.bodySmall?.copyWith(
                          color: theme.colorScheme.error,
                        ),
                      ),
                    ),

                  const SizedBox(height: 24),

                  // ── Actions ──
                  _Actions(
                    isSubmitting: state.isSubmitting,
                    canSubmit: !state.hasUploading && !state.hasErrors,
                    onCancel: () => Navigator.of(context).pop(),
                    onSubmit: _submit,
                  ),
                ],
              ),
            ),
          ),
        );
      },
    );
  }
}

// ──────────────────────────────────────────────
// Мелкие sub-виджеты
// ──────────────────────────────────────────────

class _Header extends StatelessWidget {
  final VoidCallback onClose;
  const _Header({required this.onClose});

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Text('New Coworking', style: Theme.of(context).textTheme.titleLarge),
        const Spacer(),
        IconButton(icon: const Icon(Icons.close), onPressed: onClose),
      ],
    );
  }
}

class _DropZone extends StatelessWidget {
  final bool isDragOver;
  final VoidCallback onDragEntered;
  final VoidCallback onDragExited;
  final Future<void> Function(List<XFile>) onDrop;
  final VoidCallback onTap;

  const _DropZone({
    required this.isDragOver,
    required this.onDragEntered,
    required this.onDragExited,
    required this.onDrop,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return DropTarget(
      onDragEntered: (_) => onDragEntered(),
      onDragExited: (_) => onDragExited(),
      onDragDone: (details) => onDrop(details.files),
      child: GestureDetector(
        onTap: onTap,
        child: AnimatedContainer(
          duration: const Duration(milliseconds: 200),
          width: double.infinity,
          height: 100,
          decoration: BoxDecoration(
            color: isDragOver
                ? theme.colorScheme.primaryContainer
                : theme.colorScheme.surfaceContainerHighest.withOpacity(0.4),
            border: Border.all(
              color: isDragOver
                  ? theme.colorScheme.primary
                  : theme.colorScheme.outline.withOpacity(0.5),
              width: 2,
            ),
            borderRadius: BorderRadius.circular(12),
          ),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(
                Icons.cloud_upload_outlined,
                size: 32,
                color: isDragOver
                    ? theme.colorScheme.primary
                    : theme.colorScheme.onSurfaceVariant,
              ),
              const SizedBox(height: 6),
              Text(
                isDragOver
                    ? 'Drop images here'
                    : 'Drag & drop or click to select',
                style: theme.textTheme.bodySmall?.copyWith(
                  color: theme.colorScheme.onSurfaceVariant,
                ),
              ),
              Text(
                'JPG, PNG, WEBP · max $_maxImages images', // _maxImages недоступен тут — см. ниже
                style: theme.textTheme.labelSmall?.copyWith(
                  color: theme.colorScheme.onSurfaceVariant.withOpacity(0.6),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  static const _maxImages = 10; // локальная копия константы
}

class _ThumbnailList extends StatelessWidget {
  final List<MediaItemState> items;
  final void Function(String localId) onRemove;

  const _ThumbnailList({required this.items, required this.onRemove});

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: 96,
      child: ListView.separated(
        scrollDirection: Axis.horizontal,
        itemCount: items.length,
        separatorBuilder: (_, __) => const SizedBox(width: 8),
        itemBuilder: (_, index) => _ThumbnailItem(
          item: items[index],
          onRemove: () => onRemove(items[index].localId),
        ),
      ),
    );
  }
}

class _ThumbnailItem extends StatefulWidget {
  final MediaItemState item;
  final VoidCallback onRemove;

  const _ThumbnailItem({required this.item, required this.onRemove});

  @override
  State<_ThumbnailItem> createState() => _ThumbnailItemState();
}

class _ThumbnailItemState extends State<_ThumbnailItem> {
  bool _hovered = false;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isUploading = widget.item.status == MediaUploadStatus.uploading;
    final isError = widget.item.status == MediaUploadStatus.error;

    return MouseRegion(
      onEnter: (_) => setState(() => _hovered = true),
      onExit: (_) => setState(() => _hovered = false),
      child: SizedBox(
        width: 96,
        height: 96,
        child: Stack(
          fit: StackFit.expand,
          children: [
            // ── Image ──
            ClipRRect(
              borderRadius: BorderRadius.circular(8),
              child: Image.memory(widget.item.previewBytes, fit: BoxFit.cover),
            ),

            // ── Loading overlay ──
            if (isUploading)
              Container(
                decoration: BoxDecoration(
                  color: Colors.black45,
                  borderRadius: BorderRadius.circular(8),
                ),
                child: const Center(
                  child: SizedBox(
                    width: 24,
                    height: 24,
                    child: CircularProgressIndicator(
                      strokeWidth: 2,
                      color: Colors.white,
                    ),
                  ),
                ),
              ),

            // ── Error overlay ──
            if (isError)
              Container(
                decoration: BoxDecoration(
                  color: theme.colorScheme.error.withOpacity(0.6),
                  borderRadius: BorderRadius.circular(8),
                ),
                child: const Center(
                  child: Icon(Icons.error_outline, color: Colors.white),
                ),
              ),

            // ── Remove button (visible on hover, hidden while uploading) ──
            if (!isUploading)
              AnimatedOpacity(
                duration: const Duration(milliseconds: 150),
                opacity: _hovered ? 1.0 : 0.0,
                child: Align(
                  alignment: Alignment.topRight,
                  child: Padding(
                    padding: const EdgeInsets.all(4),
                    child: GestureDetector(
                      onTap: widget.onRemove,
                      child: Container(
                        decoration: const BoxDecoration(
                          color: Colors.black54,
                          shape: BoxShape.circle,
                        ),
                        padding: const EdgeInsets.all(2),
                        child: const Icon(
                          Icons.close,
                          size: 14,
                          color: Colors.white,
                        ),
                      ),
                    ),
                  ),
                ),
              ),
          ],
        ),
      ),
    );
  }
}

class _Actions extends StatelessWidget {
  final bool isSubmitting;
  final bool canSubmit;
  final VoidCallback onCancel;
  final VoidCallback onSubmit;

  const _Actions({
    required this.isSubmitting,
    required this.canSubmit,
    required this.onCancel,
    required this.onSubmit,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.end,
      children: [
        TextButton(onPressed: onCancel, child: const Text('Cancel')),
        const SizedBox(width: 8),
        FilledButton.icon(
          onPressed: (isSubmitting || !canSubmit) ? null : onSubmit,
          icon: isSubmitting
              ? const SizedBox(
                  width: 16,
                  height: 16,
                  child: CircularProgressIndicator(strokeWidth: 2),
                )
              : const Icon(Icons.add),
          label: const Text('Create'),
        ),
      ],
    );
  }
}
