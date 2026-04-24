import 'package:coworking_app/core/models/coworking.dart';
import '../../../core/models/booking.dart';
import '../../../core/models/layout.dart';
import '../../../core/models/place.dart';
import '../../../core/models/user.dart';
import '../../../core/utils/bloc_load_state.dart';


enum AdminView {
  coworkingDetails,
  users,
}

class AdminState {
  final LoadState<List<Coworking>> coworkings;

  final LoadState<Coworking> selectedCoworking;
  final LoadState<List<Place>> places;

  final LoadState<List<LayoutVersion>> layoutVersions;
  final LoadState<Layout> layout;

  final LoadState<PaginatedBookings> bookings;
  final BookingsFilter bookingsFilter;

  final LoadState<PaginatedUsers> users;
  final LoadState<User> selectedUser;

  final String? actionMessage;
  final bool isError;

  final AdminView currentView;

  const AdminState({
    this.coworkings = const LoadState(),
    this.selectedCoworking = const LoadState(),
    this.places = const LoadState(),
    this.layoutVersions = const LoadState(),
    this.layout = const LoadState(),
    this.bookings = const LoadState(),
    this.bookingsFilter = const BookingsFilter(),
    this.users = const LoadState(),
    this.selectedUser = const LoadState(),
    this.actionMessage = '',
    this.isError = false,
    this.currentView = AdminView.coworkingDetails,
  });

  AdminState copyWith({
    LoadState<List<Coworking>>? coworkings,
    LoadState<Coworking>? selectedCoworking,
    LoadState<List<Place>>? places,
    LoadState<List<LayoutVersion>>? layoutVersions,
    LoadState<Layout>? layout,
    LoadState<PaginatedBookings>? bookings,
    BookingsFilter? bookingsFilter,
    LoadState<PaginatedUsers>? users,
    LoadState<User>? selectedUser,
    String? actionMessage,
    bool? isError,
    AdminView? currentView,
  }) {
    return AdminState(
      coworkings: coworkings ?? this.coworkings,
      selectedCoworking: selectedCoworking ?? this.selectedCoworking,
      places: places ?? this.places,
      layoutVersions: layoutVersions ?? this.layoutVersions,
      layout: layout ?? this.layout,
      bookings: bookings ?? this.bookings,
      bookingsFilter: bookingsFilter ?? this.bookingsFilter,
      users: users ?? this.users,
      selectedUser: selectedUser ?? this.selectedUser,
      actionMessage: actionMessage,
      isError: isError ?? this.isError,
      currentView: currentView ?? this.currentView,
    );
  }
}

class BookingsFilter {
  final DateTime? date;
  final String? placeType;
  final String sortBy;
  final int page;
  final int pageSize;

  const BookingsFilter({
    this.date,
    this.placeType,
    this.sortBy = 'desc',
    this.page = 1,
    this.pageSize = 6,
  });

  BookingsFilter copyWith({
    DateTime? date,
    bool clearDate = false,
    String? placeType,
    bool clearPlaceType = false,
    String? sortBy,
    int? page,
    int? pageSize,
    bool resetPage = false,
  }) {
    return BookingsFilter(
      date: clearDate ? null : (date ?? this.date),
      placeType: clearPlaceType ? null : (placeType ?? this.placeType),
      sortBy: sortBy ?? this.sortBy,
      page: resetPage ? 1 : (page ?? this.page),
      pageSize: pageSize ?? this.pageSize,
    );
  }
}