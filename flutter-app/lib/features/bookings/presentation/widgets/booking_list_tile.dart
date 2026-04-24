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

class _BookingTileState extends State<BookingTile> {
  bool expanded = false;

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
  void initState() {
    super.initState();

    if (widget.isHighlighted) {
      expanded = true;

      // авто-убрать подсветку через время (опционально)
      Future.delayed(const Duration(seconds: 5), () {
        if (mounted) {
          setState(() {});
        }
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final b = widget.booking;

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      child: Material(
        borderRadius: BorderRadius.circular(24),
        color: widget.isHighlighted
            ? Colors.yellow.shade300
            : Theme.of(context).cardColor,
        child: InkWell(
          borderRadius: BorderRadius.circular(24),
          onTap: () => {setState(() => expanded = !expanded)},
          child: AnimatedContainer(
            duration: const Duration(milliseconds: 200),
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // =====================
                // HEADER (COMPACT VIEW)
                // =====================
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    // DATE + TIME
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

                    // STATUS + ICON
                    Row(
                      children: [
                        Container(
                          padding: const EdgeInsets.symmetric(
                            horizontal: 10,
                            vertical: 6,
                          ),
                          decoration: BoxDecoration(
                            color: _statusColor(
                              b.status.name,
                            ).withOpacity(0.15),
                            borderRadius: BorderRadius.circular(12),
                          ),
                          child: Text(
                            switch(b.status){
                              BookingStatus.active => context.l10n.bookingStatusActive,
                              BookingStatus.cancelled => context.l10n.bookingStatusCancelled,
                              BookingStatus.completed => context.l10n.bookingStatusCompleted,
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

                // =====================
                // EXPANDED CONTENT
                // =====================
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
                      '${context.l10n.bookingLabelCancelReason} ${b.cancelReason}',
                    ),
                  ],

                  if (b.status == BookingStatus.active) ...[
                    const SizedBox(height: 16),
                    SizedBox(
                      width: double.infinity,
                      child: ElevatedButton(
                        onPressed: () {
                          context.read<BookingBloc>().add(
                            CancelBookingEvent(
                              b.id,
                              reason: 'cancelled_by_user',
                            ),
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
        ),
      ),
    );
  }
}
