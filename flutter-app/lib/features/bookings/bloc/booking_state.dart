import '../../../core/models/booking.dart';

class BookingState {
  final List<Booking> activeItems;
  final List<Booking> historyItems;

  final bool activeLoading;
  final bool historyLoading;

  final bool activeLoadingMore;
  final bool historyLoadingMore;

  final int activePage;
  final int historyPage;

  final bool activeHasMore;
  final bool historyHasMore;

  final String? actionMessage;
  final double? messageId;
  final bool isError;

  final String? highlightedBookingId;

  const BookingState({
    this.activeItems = const [],
    this.historyItems = const [],
    this.activeLoading = false,
    this.historyLoading = false,
    this.activeLoadingMore = false,
    this.historyLoadingMore = false,
    this.activePage = 1,
    this.historyPage = 1,
    this.activeHasMore = true,
    this.historyHasMore = true,
    this.actionMessage,
    this.messageId,
    this.isError = false,
    this.highlightedBookingId,
  });

  BookingState copyWith({
    List<Booking>? activeItems,
    List<Booking>? historyItems,
    bool? activeLoading,
    bool? historyLoading,
    bool? activeLoadingMore,
    bool? historyLoadingMore,
    int? activePage,
    int? historyPage,
    bool? activeHasMore,
    bool? historyHasMore,
    String? actionMessage,
    double? messageId,
    bool? isError,
    String? Function()? highlightedBookingId,
  }) {
    return BookingState(
      activeItems: activeItems ?? this.activeItems,
      historyItems: historyItems ?? this.historyItems,
      activeLoading: activeLoading ?? this.activeLoading,
      historyLoading: historyLoading ?? this.historyLoading,
      activeLoadingMore: activeLoadingMore ?? this.activeLoadingMore,
      historyLoadingMore: historyLoadingMore ?? this.historyLoadingMore,
      activePage: activePage ?? this.activePage,
      historyPage: historyPage ?? this.historyPage,
      activeHasMore: activeHasMore ?? this.activeHasMore,
      historyHasMore: historyHasMore ?? this.historyHasMore,
      actionMessage: actionMessage ?? this.actionMessage,
      messageId: messageId ?? this.messageId,
      isError: isError ?? this.isError,
      highlightedBookingId: highlightedBookingId != null ? highlightedBookingId() : this.highlightedBookingId,
    );
  }
}
