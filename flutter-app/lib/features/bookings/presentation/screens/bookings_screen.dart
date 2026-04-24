import 'package:coworking_app/features/bookings/presentation/widgets/booking_list_tile.dart';
import 'package:coworking_app/core/l10n/localization_helper.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../bloc/booking_bloc.dart';
import '../../bloc/booking_event.dart';
import '../../bloc/booking_state.dart';

import '../../../../core/di/service_locator.dart';

class BookingsPage extends StatelessWidget {
  final String? initialTab;
  final String? highlightBookingId;
  const BookingsPage({super.key, this.initialTab, this.highlightBookingId});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) => sl<BookingBloc>()
        ..add(FetchActiveBookings())
        ..add(FetchHistoryBookings()),
      child: Builder(
        builder: (context) {
          return BlocListener<BookingBloc, BookingState>(
            listenWhen: (prev, curr) => prev.messageId != curr.messageId,
            listener: (context, state) {
              if (state.actionMessage != null) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    duration: Duration(seconds: 1),
                    content: Text(state.actionMessage!),
                    backgroundColor: state.isError ? Colors.red : Colors.green,
                  ),
                );
              }
            },
            child: BookingsScreen(
              initialTab: initialTab,
              highlightBookingId: highlightBookingId,
            ),
          );
        },
      ),
    );
  }
}

// =====================
// MAIN SCREEN
// =====================

class BookingsScreen extends StatefulWidget {
  final String? initialTab;
  final String? highlightBookingId;

  const BookingsScreen({super.key, this.initialTab, this.highlightBookingId});

  @override
  State<BookingsScreen> createState() => _BookingsScreenState();
}

class _BookingsScreenState extends State<BookingsScreen>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();

    _tabController = TabController(length: 2, vsync: this);

    // 👉 tab
    if (widget.initialTab == 'history') {
      _tabController.index = 1;
    } else {
      _tabController.index = 0;
    }

    // 👉 highlight
    if (widget.highlightBookingId != null) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        context.read<BookingBloc>().add(
          HighlightBooking(widget.highlightBookingId!),
        );
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(context.l10n.bookingsScreenTitle),
        automaticallyImplyLeading: false,
        bottom: TabBar(
          controller: _tabController,
          tabs: [
            Tab(text: context.l10n.bookingsTabActive),
            Tab(text: context.l10n.bookingsTabHistory),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          _BookingList(
            isActive: true,
            highlightBookingId: widget.highlightBookingId,
          ),
          _BookingList(
            isActive: false,
            highlightBookingId: widget.highlightBookingId,
          ),
        ],
      ),
    );
  }
}

// =====================
// LIST WITH INFINITE SCROLL
// =====================

class _BookingList extends StatefulWidget {
  final bool isActive;
  final String? highlightBookingId;

  const _BookingList({required this.isActive, this.highlightBookingId});

  @override
  State<_BookingList> createState() => _BookingListState();
}

class _BookingListState extends State<_BookingList> {
  bool _scrolled = false;

  final ScrollController _controller = ScrollController();

  @override
  void initState() {
    super.initState();
    _controller.addListener(_onScroll);
  }

  void _onScroll() {
    if (!mounted) return;

    if (_controller.position.pixels >=
        _controller.position.maxScrollExtent - 200) {
      final bloc = context.read<BookingBloc>();

      if (widget.isActive) {
        bloc.add(LoadMoreActive());
      } else {
        bloc.add(LoadMoreHistory());
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<BookingBloc, BookingState>(
      builder: (context, state) {
        final items = widget.isActive ? state.activeItems : state.historyItems;

        if (!_scrolled &&
            widget.highlightBookingId != null &&
            items.isNotEmpty) {
          WidgetsBinding.instance.addPostFrameCallback((_) {
            final index = items.indexWhere(
              (b) => b.id == widget.highlightBookingId,
            );

            if (index != -1) {
              _controller.animateTo(
                index * 100, // примерная высота
                duration: const Duration(milliseconds: 400),
                curve: Curves.easeInOut,
              );
              _scrolled = true;
            }
          });
        }

        final loading = widget.isActive
            ? state.activeLoading
            : state.historyLoading;

        final loadingMore = widget.isActive
            ? state.activeLoadingMore
            : state.historyLoadingMore;

        return RefreshIndicator(
          onRefresh: () async {
            if (widget.isActive) {
              context.read<BookingBloc>().add(
                FetchActiveBookings(refresh: true),
              );
            } else {
              context.read<BookingBloc>().add(
                FetchHistoryBookings(refresh: true),
              );
            }
          },
          child: loading && items.isEmpty
              ? const Center(child: CircularProgressIndicator())
              : items.isEmpty
              ? Center(
                  child: Text(
                    widget.isActive
                        ? context.l10n.noActiveBookings
                        : context.l10n.noHistoryBookings,
                    style: TextStyle(color: Colors.grey.shade600, fontSize: 16),
                  ),
                )
              : ListView.builder(
                  controller: _controller,
                  itemCount: items.length + (loadingMore ? 1 : 0),
                  itemBuilder: (_, i) {
                    if (i >= items.length) {
                      return const Padding(
                        padding: EdgeInsets.all(16),
                        child: Center(child: CircularProgressIndicator()),
                      );
                    }

                    return BookingTile(
                      booking: items[i],
                      isActive: widget.isActive,
                      isHighlighted: items[i].id == widget.highlightBookingId,
                    );
                  },
                ),
        );
      },
    );
  }
}
