import 'package:coworking_app/core/utils/media_url_builder.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:coworking_app/core/l10n/localization_helper.dart';

import '../../../../core/models/coworking.dart';
import '../../bloc/coworking_bloc.dart';
import '../../bloc/coworking_event.dart';
import '../screens/coworking_details_screen.dart';

class CoworkingCard extends StatelessWidget {
  final Coworking coworking;

  const CoworkingCard({super.key, required this.coworking});

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: () {
        if (!coworking.isActive) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              duration: const Duration(seconds: 1),
              content: Text(context.l10n.coworkingUnavailable),
            ),
          );
          return;
        }
        final bloc = context.read<CoworkingBloc>();
        bloc.add(SelectCoworking(coworking.id));

        Navigator.push(
          context,
          MaterialPageRoute(
            builder: (_) => BlocProvider.value(
              value: bloc,
              child: CoworkingDetailsScreen(coworkingId: coworking.id),
            ),
          ),
        );
      },
      child: Card(
        shadowColor: Colors.blue,
        surfaceTintColor: coworking.isActive ? Colors.white : Colors.black,
        elevation: 12,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            /// 🔷 Gallery
            _CoworkingImageGallery(
              imageIds: coworking.imageIDs,
              isActive: coworking.isActive,
            ),

            /// 🔷 Info
            Padding(
              padding: const EdgeInsets.all(12),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    coworking.name,
                    style: const TextStyle(
                      fontWeight: FontWeight.bold,
                      fontSize: 16,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    coworking.address,
                    style: TextStyle(color: Colors.grey[600], fontSize: 12),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

// ──────────────────────────────────────────────
// Gallery widget
// ──────────────────────────────────────────────

class _CoworkingImageGallery extends StatefulWidget {
  final List<String> imageIds;
  final bool isActive;

  const _CoworkingImageGallery({
    required this.imageIds,
    required this.isActive,
  });

  @override
  State<_CoworkingImageGallery> createState() => _CoworkingImageGalleryState();
}

class _CoworkingImageGalleryState extends State<_CoworkingImageGallery> {
  final PageController _pageController = PageController();
  int _currentPage = 0;

  @override
  void dispose() {
    _pageController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final hasImages = widget.imageIds.isNotEmpty;
    final count = widget.imageIds.length;

    print("Building gallery for coworking with ${widget.imageIds.length} images, active: ${widget.isActive}");

    return ClipRRect(
      borderRadius: const BorderRadius.vertical(top: Radius.circular(12)),
      child: SizedBox(
        height: 250,
        child: Stack(
          fit: StackFit.expand,
          children: [
            // ── Images or placeholder ──
            if (!hasImages)
              _Placeholder(isActive: widget.isActive)
            else
              // Останавливаем свайп карточки при скролле галереи
              GestureDetector(
                onHorizontalDragUpdate: (_) {},
                child: PageView.builder(
                  controller: _pageController,
                  itemCount: count,
                  onPageChanged: (i) => setState(() => _currentPage = i),
                  itemBuilder: (context, index) {
                    return Image.network(
                      MediaUrlBuilder.build(widget.imageIds[index]),
                      fit: BoxFit.cover,
                      // Пока грузится — показываем shimmer-заглушку
                      loadingBuilder: (_, child, progress) {
                        if (progress == null) return child;
                        return Container(
                          color: Colors.blue.withOpacity(0.15),
                          child: const Center(
                            child: CircularProgressIndicator(strokeWidth: 2),
                          ),
                        );
                      },
                      errorBuilder: (_, __, ___) => Container(
                        color: Colors.blue.withOpacity(0.15),
                        child: const Center(
                          child: Icon(Icons.broken_image_outlined, size: 48),
                        ),
                      ),
                    );
                  },
                ),
              ),

            // ── Inactive overlay ──
            if (!widget.isActive)
              Center(
                child: Icon(
                  Icons.block,
                  size: 200,
                  color: Colors.grey.withAlpha(150),
                ),
              ),

            // ── Page indicator (только если больше одной картинки) ──
            if (hasImages && count > 1)
              Positioned(
                bottom: 10,
                left: 0,
                right: 0,
                child: _PageIndicator(count: count, currentIndex: _currentPage),
              ),
          ],
        ),
      ),
    );
  }
}

// ──────────────────────────────────────────────
// Placeholder — когда картинок нет
// ──────────────────────────────────────────────

class _Placeholder extends StatelessWidget {
  final bool isActive;
  const _Placeholder({required this.isActive});

  @override
  Widget build(BuildContext context) {
    return Container(
      color: Colors.blue.withOpacity(0.15),
      child: const Center(child: Icon(Icons.business, size: 100)),
    );
  }
}

// ──────────────────────────────────────────────
// Dot indicator
// ──────────────────────────────────────────────

class _PageIndicator extends StatelessWidget {
  final int count;
  final int currentIndex;

  const _PageIndicator({required this.count, required this.currentIndex});

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: List.generate(count, (i) {
        final isActive = i == currentIndex;
        return AnimatedContainer(
          duration: const Duration(milliseconds: 250),
          curve: Curves.easeOutCubic,
          margin: const EdgeInsets.symmetric(horizontal: 3),
          height: 3,
          width: isActive ? 20 : 8,
          decoration: BoxDecoration(
            color: isActive ? Colors.white : Colors.white.withOpacity(0.5),
            borderRadius: BorderRadius.circular(2),
          ),
        );
      }),
    );
  }
}
