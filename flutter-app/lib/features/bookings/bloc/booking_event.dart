abstract class BookingEvent {}

class FetchActiveBookings extends BookingEvent {
  final bool refresh;
  FetchActiveBookings({this.refresh = false});
}

class FetchHistoryBookings extends BookingEvent {
  final bool refresh;
  FetchHistoryBookings({this.refresh = false});
}

class LoadMoreActive extends BookingEvent {}

class LoadMoreHistory extends BookingEvent {}

class CancelBookingEvent extends BookingEvent {
  final String id;
  final String? reason;
  CancelBookingEvent(this.id, {this.reason});
}

class HighlightBooking extends BookingEvent {
  final String bookingId;
  HighlightBooking(this.bookingId);
}
