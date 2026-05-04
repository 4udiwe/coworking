import 'package:coworking_app/core/di/service_locator.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:coworking_app/core/l10n/localization_helper.dart';
import 'package:intl/intl.dart';

import '../../../../core/models/booking.dart';
import '../../bloc/booking_bloc.dart';
import '../../bloc/booking_event.dart';

final _dateOnlyFormat = DateFormat('dd MMMM', 'ru');
final _timeFormat = DateFormat('HH:mm');

class BookingTile extends StatefulWidget {
  final Booking booking;
  final bool isActive;
  final bool isHighlighted;

  const BookingTile({
    super.key,
    required this.booking,
    required this.isActive,
    this.isHighlighted = false,
  });

  @override
  State<BookingTile> createState() => _BookingTileState();
}

class _BookingTileState extends State<BookingTile>
    with SingleTickerProviderStateMixin {
  bool expanded = false;
  bool _highlightActive = false;

  AnimationController? _animController;
  Animation<Color?>? _colorAnim;

  // Цвета подсветки
  static const _highlightColor = Colors.blue;
  static const _blinkCount = 3;

  @override
  void initState() {
    super.initState();

    if (widget.isHighlighted) {
      expanded = true;
      _highlightActive = true;
      _startBlinkAnimation();
    }
  }

  void _startBlinkAnimation() {
    _animController = AnimationController(
      vsync: this,
      // Один цикл мигания: fade in + fade out
      duration: const Duration(milliseconds: 400),
    );

    _colorAnim = ColorTween(
      begin: Colors.transparent,
      end: _highlightColor.withOpacity(0.5),
    ).animate(_animController!);

    // Запускаем N миганий, затем оставляем подсветку
    int blinksLeft = _blinkCount;

    _animController!.addStatusListener((status) {
      if (!mounted) return;

      if (status == AnimationStatus.completed) {
        if (blinksLeft > 1) {
          blinksLeft--;
          _animController!.reverse();
        }
        // После последнего forward — остаёмся на completed (подсветка держится)
      } else if (status == AnimationStatus.dismissed) {
        _animController!.forward();
      }
    });

    _animController!.forward();
  }

  void _dismissHighlight() {
    if (!_highlightActive) return;
    sl<BookingBloc>().add(DisableBookingHighlight());
    setState(() => _highlightActive = false);
    _animController?.stop();
    _animController?.dispose();
    _animController = null;
  }

  @override
  void dispose() {
    _animController?.dispose();
    super.dispose();
  }

  Color _statusColor(String status) {
    switch (status) {
      case 'active':
        return Colors.green;
      case 'cancelled':
        return Colors.red;
      default:
        return Colors.grey;
    }
  }

  @override
  Widget build(BuildContext context) {
    final b = widget.booking;

    // Базовый цвет карточки
    final baseColor = Theme.of(context).cardColor;

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      child: AnimatedBuilder(
        animation: _animController ?? const AlwaysStoppedAnimation(0),
        builder: (context, child) {
          Color cardColor;

          if (!_highlightActive) {
            cardColor = baseColor;
          } else if (_animController != null) {
            // Во время анимации — интерполируем
            cardColor = Color.lerp(
                  baseColor,
                  _highlightColor.withOpacity(0.4),
                  _colorAnim!.value?.opacity ?? 0,
                ) ??
                baseColor;
          } else {
            cardColor = baseColor;
          }

          return Material(
            borderRadius: BorderRadius.circular(24),
            color: cardColor,
            child: InkWell(
              borderRadius: BorderRadius.circular(24),
              onTap: () {
                _dismissHighlight();
                setState(() => expanded = !expanded);
              },
              child: AnimatedContainer(
                duration: const Duration(milliseconds: 200),
                padding: const EdgeInsets.all(16),
                child: child,
              ),
            ),
          );
        },
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // HEADER
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      _dateOnlyFormat.format(b.startTime.toLocal()),
                      style: const TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      '${_timeFormat.format(b.startTime.toLocal())} - ${_timeFormat.format(b.endTime.toLocal())}',
                      style: TextStyle(color: Colors.grey[600]),
                    ),
                  ],
                ),
                Row(
                  children: [
                    Container(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 10,
                        vertical: 6,
                      ),
                      decoration: BoxDecoration(
                        color: _statusColor(b.status.name).withOpacity(0.15),
                        borderRadius: BorderRadius.circular(12),
                      ),
                      child: Text(
                        switch (b.status) {
                          BookingStatus.active =>
                            context.l10n.bookingStatusActive,
                          BookingStatus.cancelled =>
                            context.l10n.bookingStatusCancelled,
                          BookingStatus.completed =>
                            context.l10n.bookingStatusCompleted,
                        },
                        style: TextStyle(
                          color: _statusColor(b.status.name),
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                    ),
                    const SizedBox(width: 8),
                    Icon(expanded ? Icons.expand_less : Icons.expand_more),
                  ],
                ),
              ],
            ),

            // EXPANDED
            if (expanded) ...[
              const SizedBox(height: 16),
              const Divider(),
              const SizedBox(height: 12),
              Text(
                '${context.l10n.bookingLabelCoworking} ${b.place.coworkingName}',
                style: const TextStyle(fontWeight: FontWeight.w500),
              ),
              Text(
                '${context.l10n.bookingLabelPlace} ${b.place.label}',
                style: const TextStyle(fontWeight: FontWeight.w500),
              ),
              const SizedBox(height: 8),
              if (b.cancelReason != null && b.cancelReason!.isNotEmpty) ...[
                const SizedBox(height: 8),
                Text(
                    '${context.l10n.bookingLabelCancelReason} ${b.cancelReason}'),
              ],
              if (b.status == BookingStatus.active) ...[
                const SizedBox(height: 16),
                SizedBox(
                  width: double.infinity,
                  child: ElevatedButton(
                    onPressed: () {
                      context.read<BookingBloc>().add(
                            CancelBookingEvent(b.id,
                                reason: 'cancelled_by_user'),
                          );
                      setState(() => expanded = false);
                    },
                    child: Text(context.l10n.bookingButtonCancel),
                  ),
                ),
              ],
            ],
          ],
        ),
      ),
    );
  }
}