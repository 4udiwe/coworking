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

    print("Initial tap = " + widget.initialTab.toString());

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

// В _BookingList добавляем Map ключей для каждого элемента
class _BookingListState extends State<_BookingList> {
  final ScrollController _controller = ScrollController();
  final Map<String, GlobalKey> _itemKeys = {};
  bool _scrolled = false;

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
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  void _scrollToHighlighted() {
    if (widget.highlightBookingId == null || _scrolled) return;
    final key = _itemKeys[widget.highlightBookingId];
    if (key?.currentContext != null) {
      _scrolled = true;
      Scrollable.ensureVisible(
        key!.currentContext!,
        duration: const Duration(milliseconds: 400),
        curve: Curves.easeInOut,
        alignment: 0.3, // элемент окажется в верхней трети экрана
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<BookingBloc, BookingState>(
      builder: (context, state) {
        final items = widget.isActive ? state.activeItems : state.historyItems;
        final loading = widget.isActive
            ? state.activeLoading
            : state.historyLoading;
        final loadingMore = widget.isActive
            ? state.activeLoadingMore
            : state.historyLoadingMore;

        // Скроллим после рендера
        if (widget.highlightBookingId != null &&
            items.isNotEmpty &&
            !_scrolled) {
          WidgetsBinding.instance.addPostFrameCallback((_) {
            _scrollToHighlighted();
          });
        }

        if (loading && items.isEmpty) {
          return const Center(child: CircularProgressIndicator());
        }

        if (items.isEmpty) {
          return Center(
            child: Text(
              widget.isActive
                  ? context.l10n.noActiveBookings
                  : context.l10n.noHistoryBookings,
              style: TextStyle(color: Colors.grey.shade600, fontSize: 16),
            ),
          );
        }

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
          // Listener сбрасывает подсветку при любом касании списка
          child: Listener(
            onPointerDown: (_) {
              // Сигнал тайлу через BLoC или просто перестраиваем —
              // здесь проще всего снять highlight через setState родителя,
              // но т.к. тайл сам управляет своей анимацией через _dismissHighlight,
              // достаточно что onTap на тайле её сбрасывает.
              // Если нужно сбросить даже без тапа на сам тайл — используй подход ниже.
            },
            child: ListView.builder(
              controller: _controller,
              itemCount: items.length + (loadingMore ? 1 : 0),
              itemBuilder: (_, i) {
                if (i >= items.length) {
                  return const Padding(
                    padding: EdgeInsets.all(16),
                    child: Center(child: CircularProgressIndicator()),
                  );
                }

                final booking = items[i];
                final isHighlighted = booking.id == widget.highlightBookingId;

                // Создаём ключ для элемента при необходимости
                final key = _itemKeys.putIfAbsent(
                  booking.id,
                  () => GlobalKey(),
                );

                return BookingTile(
                  key: key,
                  booking: booking,
                  isActive: widget.isActive,
                  isHighlighted: isHighlighted,
                );
              },
            ),
          ),
        );
      },
    );
  }
}
