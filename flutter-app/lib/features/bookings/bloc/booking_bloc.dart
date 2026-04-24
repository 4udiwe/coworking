import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../core/models/booking.dart';
import '../data/booking_repository.dart';
import 'booking_event.dart';
import 'booking_state.dart';

class BookingBloc extends Bloc<BookingEvent, BookingState> {
  final BookingRepository repository;

  BookingBloc({required this.repository}) : super(const BookingState()) {
    on<FetchActiveBookings>(_fetchActive);
    on<FetchHistoryBookings>(_fetchHistory);
    on<LoadMoreActive>(_loadMoreActive);
    on<LoadMoreHistory>(_loadMoreHistory);
    on<CancelBookingEvent>(_cancel);

    on<HighlightBooking>((event, emit) {
      emit(state.copyWith(
        highlightedBookingId: () => event.bookingId,
      ));
    });
  }

  Future<void> _fetchActive(FetchActiveBookings e, Emitter<BookingState> emit) async {
    emit(state.copyWith(activeLoading: true, activePage: 1));

    try {
      final res = await repository.getBookings(
          page: 1, status: BookingStatus.active.name, pageSize: 10);

      emit(state.copyWith(
        activeItems: res.items,
        activeLoading: false,
        activeHasMore: 1 < res.totalPages,
        activePage: 1,
      ));
    } catch (e) {
      emit(state.copyWith(
        activeLoading: false,
        actionMessage: 'Failed to fetch active bookings',
        messageId: DateTime.now().millisecondsSinceEpoch.toDouble(),
        isError: true,
      ));
    }
  }

  Future<void> _fetchHistory(FetchHistoryBookings e, Emitter<BookingState> emit) async {
    emit(state.copyWith(historyLoading: true, historyPage: 1));

    try {
      final res = await repository.getBookings(page: 1, pageSize: 10);

      emit(state.copyWith(
        historyItems: res.items,
        historyLoading: false,
        historyHasMore: 1 < res.totalPages,
        historyPage: 1,
      ));
    } catch (e) {
      emit(state.copyWith(
        historyLoading: false,
        actionMessage: 'Failed to fetch history bookings',
        messageId: DateTime.now().millisecondsSinceEpoch.toDouble(),
        isError: true,
      ));
    }
  }

  Future<void> _loadMoreActive(LoadMoreActive e, Emitter<BookingState> emit) async {
    if (!state.activeHasMore || state.activeLoadingMore) return;

    emit(state.copyWith(activeLoadingMore: true));

    final next = state.activePage + 1;
    final res = await repository.getBookings(page: next, status: BookingStatus.active.name, pageSize: 10);

    emit(state.copyWith(
      activeItems: [...state.activeItems, ...res.items],
      activePage: next,
      activeHasMore: next < res.totalPages,
      activeLoadingMore: false,
    ));
  }

  Future<void> _loadMoreHistory(LoadMoreHistory e, Emitter<BookingState> emit) async {
    if (!state.historyHasMore || state.historyLoadingMore) return;

    emit(state.copyWith(historyLoadingMore: true));

    final next = state.historyPage + 1;
    final res = await repository.getBookings(page: next, pageSize: 10);

    emit(state.copyWith(
      historyItems: [...state.historyItems, ...res.items],
      historyPage: next,
      historyHasMore: next < res.totalPages,
      historyLoadingMore: false,
    ));
  }

  Future<void> _cancel(CancelBookingEvent e, Emitter<BookingState> emit) async {
    try {
      await repository.cancelBooking(e.id, reason: e.reason);


      add(FetchActiveBookings());
      add(FetchHistoryBookings());
      emit(state.copyWith(
        actionMessage: 'Booking canceled',
        messageId: DateTime.now().millisecondsSinceEpoch.toDouble(),
        isError: false,
      ));
    } catch (_) {
      emit(state.copyWith(
        actionMessage: 'Failed to cancel booking',
        messageId: DateTime.now().millisecondsSinceEpoch.toDouble(),
        isError: true,
      ));
    }
  }
}