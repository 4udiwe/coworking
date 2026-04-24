import 'dart:math';
import 'package:flutter/material.dart';
import '../../../core/models/layout.dart';
import '../../../core/models/place.dart';

class LayoutPreview extends StatelessWidget {
  final Layout layout;
  final List<Place> places;

  /// 🆕 интерактив
  final Function(Place place)? onPlaceTap;

  /// 🆕 состояния
  final Set<String> selectedPlaceIds;
  final Set<String> unavailablePlaceIds;

  const LayoutPreview({
    super.key,
    required this.layout,
    required this.places,
    this.onPlaceTap,
    this.selectedPlaceIds = const {},
    this.unavailablePlaceIds = const {},
  });

  @override
  Widget build(BuildContext context) {
    final canvasWidth = layout.layout.canvas.width;
    final canvasHeight = layout.layout.canvas.height;

    final Map<String, Place> placesMap = {
      for (final p in places) p.id: p,
    };

    return Container(
      color: Colors.grey[200],
      child: InteractiveViewer(
        // Ограничиваем минимальный масштаб единицей, чтобы карта всегда была вписана
        minScale: 1.0,
        maxScale: 4.0,
        // Позволяет взаимодействовать с элементами внутри (тапать по местам)
        child: FittedBox(
          fit: BoxFit.contain,
          alignment: Alignment.center,
          child: SizedBox(
            width: canvasWidth,
            height: canvasHeight,
            child: Stack(
              children: [
                /// =====================
                /// WALLS
                /// =====================
                ...layout.layout.walls.map((w) {
                  return Positioned(
                    left: w.x,
                    top: w.y,
                    child: Container(
                      width: w.width,
                      height: w.height,
                      color: Colors.brown[300],
                    ),
                  );
                }),

                /// =====================
                /// PLACES
                /// =====================
                ...layout.layout.places.map((p) {
                  final place = placesMap[p.id];
                  if (place == null) return const SizedBox();

                  final isSelected = selectedPlaceIds.contains(place.id);
                  final isUnavailable =
                  unavailablePlaceIds.contains(place.id);

                  final color = _resolveColor(
                    place: place,
                    isSelected: isSelected,
                    isUnavailable: isUnavailable,
                  );

                  final iconSize = min(p.width, p.height) * 0.5;

                  return Positioned(
                    left: p.x,
                    top: p.y,
                    child: GestureDetector(
                      onTap: onPlaceTap != null
                          ? () => onPlaceTap!(place)
                          : null,
                      child: Stack(
                        alignment: Alignment.center,
                        children: [
                          Transform.rotate(
                            angle: p.rotation * pi / 180,
                            child: Container(
                              width: p.width,
                              height: p.height,
                              decoration: BoxDecoration(
                                color: color,
                                border: Border(
                                  left: BorderSide(color: Colors.blue[900]!, width: 4),
                                  right: BorderSide(color: Colors.blue[900]!, width: 4),
                                  bottom: BorderSide(color: Colors.blue[900]!, width: 4),
                                ),
                                borderRadius: BorderRadius.circular(10),
                              ),
                            ),
                          ),

                          /// контент внутри

                          Column(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              Image.asset(
                                'assets/images/open_desk.png',
                                width: iconSize,
                                height: iconSize,
                              ),
                              Text(
                                place.label,
                                style: const TextStyle(
                                  fontWeight: FontWeight.bold,
                                ),
                              ),
                            ],
                          ),
                        ],
                      ),
                    ),
                  );
                }),
              ],
            ),
          ),
        ),
      ),
    );
  }

  /// =====================
  /// COLOR LOGIC
  /// =====================
  Color _resolveColor({
    required Place place,
    required bool isSelected,
    required bool isUnavailable,
  }) {
    if (!place.isActive) {
      return Colors.grey;
    }

    if (isUnavailable) {
      return Colors.red;
    }

    if (isSelected) {
      return Colors.green;
    }

    return Colors.blue.shade300;
  }
}
