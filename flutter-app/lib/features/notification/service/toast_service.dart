import 'package:flutter/material.dart';

class ToastService {
  static final List<_ToastItem> _items = [];
  static OverlayEntry? _overlay;


  static void show(
      BuildContext context, {
        required Widget Function(VoidCallback remove, Animation<double> animation) builder,
      }) {
    final item = _ToastItem();

    final controller = AnimationController(
      vsync: Navigator.of(context),
      duration: const Duration(milliseconds: 300),
    );

    final animation = CurvedAnimation(
      parent: controller,
      curve: Curves.easeOut,
      reverseCurve: Curves.easeIn,
    );

    item.controller = controller;
    item.animation = animation;

    item.builder = (remove, animation) => builder(remove, animation);

    _items.insert(0, item);

    _showOverlay(context);

    controller.forward().then((_) {
      Future.delayed(const Duration(seconds: 5), () {
        if (_items.contains(item)) {
          item.controller.reverse().then((_) {
            _items.remove(item);
            _overlay?.markNeedsBuild();
            if (_items.isEmpty) {
              _overlay?.remove();
              _overlay = null;
            }
          });
        }
      });
    });
  }

  static void _showOverlay(BuildContext context) {
    if (_overlay != null) {
      _overlay!.markNeedsBuild();
      return;
    }

    final overlay = Overlay.of(context);

    _overlay = OverlayEntry(
      builder: (_) {
        return Positioned(
          top: 50,
          left: 0,
          right: 0,
          child: SafeArea(
            child: Column(
              children: _items.map((item) {
                return AnimatedBuilder(
                  animation: item.animation,
                  builder: (_, child) {
                    return FadeTransition(
                      opacity: item.animation,
                      child: SlideTransition(
                        position: Tween<Offset>(
                          begin: const Offset(0, -0.3),
                          end: Offset.zero,
                        ).animate(item.animation),
                        child: child,
                      ),
                    );
                  },
                  child: item.builder(
                        () {
                      item.controller.reverse().then((_) {
                        _items.remove(item);
                        _overlay?.markNeedsBuild();

                        if (_items.isEmpty) {
                          _overlay?.remove();
                          _overlay = null;
                        }
                      });
                    },
                    item.animation,
                  ),
                );
              }).toList(),
            ),
          ),
        );
      },
    );

    overlay.insert(_overlay!);
  }
}

class _ToastItem {
  late Widget Function(VoidCallback remove, Animation<double> animation) builder;

  late AnimationController controller;
  late Animation<double> animation;
}