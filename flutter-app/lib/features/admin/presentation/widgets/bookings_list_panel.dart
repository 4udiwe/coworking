import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../../core/utils/bloc_load_state.dart';
import '../../bloc/admin_bloc.dart';
import '../../bloc/admin_event.dart';
import '../../bloc/admin_state.dart';
import '../../../../core/models/booking.dart';



String _formatDate(DateTime dt) {
  return '${dt.day.toString().padLeft(2, '0')}.'
      '${dt.month.toString().padLeft(2, '0')}.'
      '${dt.year}';
}

String _formatTime(DateTime dt) {
  return '${dt.hour.toString().padLeft(2, '0')}:'
      '${dt.minute.toString().padLeft(2, '0')}';
}

Widget _stats(List<Booking> pageBookings, PaginatedBookings paginated) {
  return Row(
    children: [
      Text(
        'Items on page: ${pageBookings.length} / Total items: ${paginated.totalItems}',
        style: const TextStyle(fontWeight: FontWeight.bold),
      ),
    ],
  );
}

class BookingListPanel extends StatelessWidget {
  const BookingListPanel({super.key});

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        const Text(
          'Active bookings',
          style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
        ),
        const SizedBox(height: 12),

        _Filters(),

        const SizedBox(height: 12),

        Expanded(
          child: BlocBuilder<AdminBloc, AdminState>(
            builder: (context, state) {
              final bookingsState = state.bookings;
              final paginated = bookingsState.data;

              final bookings = paginated?.items ?? [];

              return Column(
                children: [
                  _stats(bookings, paginated!),

                  const SizedBox(height: 12),

                  Expanded(
                    child: Stack(
                      children: [
                        AnimatedOpacity(
                          duration: const Duration(milliseconds: 250),
                          opacity: bookingsState.status == LoadStatus.loading ? 0.4 : 1,
                          child: _Table(bookings),
                        ),

                        if (bookingsState.status == LoadStatus.loading)
                          Positioned.fill(
                            child: Container(
                              color: Colors.black.withOpacity(0.05),
                              child: const Center(child: CircularProgressIndicator()),
                            ),
                          ),
                      ],
                    ),
                  ),

                  _Pagination(paginated),
                ],
              );
            },
          ),
        ),
      ],
    );
  }
}

class _Filters extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    final bloc = context.read<AdminBloc>();
    final filter = bloc.state.bookingsFilter;
    print("FILTER VALUE: ${filter.placeType}");


    return Row(
      children: [
        if (filter.date != null)
          IconButton(
            icon: const Icon(Icons.close, size: 18),
            onPressed: () {
              bloc.add(UpdateBookingsFilterEvent(clearDate: true));
            },
          ),

        /// DATE
        TextButton(
          onPressed: () async {
            final picked = await showDatePicker(
              context: context,
              initialDate: filter.date ?? DateTime.now(),
              firstDate: DateTime.now().subtract(const Duration(days: 365)),
              lastDate: DateTime.now().add(const Duration(days: 365)),
            );

            if (picked != null) {
              bloc.add(UpdateBookingsFilterEvent(date: picked));
            }
          },
          child: Text(
            filter.date == null
                ? 'Select date'
                : 'Date: ${_formatDate(filter.date!)}',
            style: TextStyle(color: Colors.blue.shade800, fontSize: 16),
          ),
        ),

        const SizedBox(width: 16),

        /// TYPE
        DropdownButton<String?>(
          style: TextStyle(color: Colors.blue.shade800, fontSize: 16),
          hint: const Text('Place type'),
          focusColor: Colors.transparent,
          value: filter.placeType,
          items: [
            DropdownMenuItem(value: 'open_desk', child: Text('Open desk')),
            DropdownMenuItem(value: 'meeting_room', child: Text('Meeting room')),
            DropdownMenuItem(value: 'private_office', child: Text('Private office')),
          ],
          onChanged: (v) {
            bloc.add(UpdateBookingsFilterEvent(placeType: v));
          },
        ),

        const Spacer(),

        /// RESET
        TextButton(
          onPressed: () {
            bloc.add(UpdateBookingsFilterEvent(
              date: null,
              clearDate: true,
              placeType: null,
              clearPlaceType: true,
              sortBy: 'desc',
            ));
          },
          child: Text('Reset', style: TextStyle(color: Colors.blue.shade800, fontSize: 16)),
        ),
      ],
    );
  }
}

class _Table extends StatelessWidget {
  final List<Booking> bookings;

  const _Table(this.bookings);

  @override
  Widget build(BuildContext context) {
    final sortBy = context.watch<AdminBloc>().state.bookingsFilter.sortBy;

    return SingleChildScrollView(
      child: DataTable(
        sortColumnIndex: 2,
        sortAscending: sortBy == 'asc',
        columns: [
          const DataColumn(label: Text('User')),
          const DataColumn(label: Text('Place')),
          DataColumn(
            label: const Text('Start'),
            onSort: (_, asc) {
              context.read<AdminBloc>().add(
                UpdateBookingsFilterEvent(
                  sortBy: asc ? 'asc' : 'desc',
                ),
              );
            },
          ),
          const DataColumn(label: Text('End')),
          const DataColumn(label: Text('Date')),
          const DataColumn(label: Text('')),
        ],
        rows: bookings.isEmpty
            ? [
          const DataRow(cells: [
            DataCell(Text('No data')),
            DataCell(Text('')),
            DataCell(Text('')),
            DataCell(Text('')),
            DataCell(Text('')),
            DataCell(Text('')),
          ])
        ]
            : bookings.map((b) {
          return DataRow(cells: [
            DataCell(Text(b.userName.isNotEmpty ? b.userName : b.userId)),
            DataCell(Text('${b.place.label} (${b.place.placeType.toTitleCase(context)})')),
            DataCell(Text(_formatTime(b.startTime.toLocal()))),
            DataCell(Text(_formatTime(b.endTime.toLocal()))),
            DataCell(Text(_formatDate(b.startTime.toLocal()))),
            DataCell(
              IconButton(
                icon: const Icon(Icons.cancel, color: Colors.red),
                onPressed: () {
                  context.read<AdminBloc>().add(
                    AdminCancelBookingEvent(b.id),
                  );
                },
              ),
            ),
          ]);
        }).toList(),
      ),
    );
  }
}

class _Pagination extends StatelessWidget {
  final PaginatedBookings paginated;

  const _Pagination(this.paginated);

  @override
  Widget build(BuildContext context) {
    final bloc = context.read<AdminBloc>();

    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        IconButton(
          icon: const Icon(Icons.chevron_left),
          onPressed: paginated.page > 1
              ? () => bloc.add(ChangeBookingsPageEvent(paginated.page - 1))
              : null,
        ),
        Text('Page ${paginated.page} of ${paginated.totalPages}'),
        IconButton(
          icon: const Icon(Icons.chevron_right),
          onPressed: paginated.page < paginated.totalPages
              ? () => bloc.add(ChangeBookingsPageEvent(paginated.page + 1))
              : null,
        ),
      ],
    );
  }
}