import 'package:flutter_bloc/flutter_bloc.dart';
import '../../../core/utils/bloc_load_state.dart';
import '../data/admin_repository.dart';
import 'admin_event.dart';
import 'admin_state.dart';

class AdminBloc extends Bloc<AdminEvent, AdminState> {
  final AdminRepository repository;

  AdminBloc({required this.repository}) : super(const AdminState()) {
    on<SetAdminViewEvent>((event, emit) {
      emit(state.copyWith(currentView: event.view));
    });

    on<FetchUsersEvent>(_onFetchUsers);
    on<SelectUserEvent>(_onSelectUser);
    on<UpdateUserRolesEvent>(_onUpdateUserRoles);
    on<SetUserActiveEvent>(_onSetUserActive);

    /// =======================
    /// QUERY EVENTS
    /// =======================

    on<FetchCoworkingsEvent>(_onFetchCoworkings);
    on<FetchCoworkingByIdEvent>(_onFetchCoworkingById);
    on<FetchPlacesEvent>(_onFetchPlaces);
    on<FetchLayoutVersionsEvent>(_onFetchLayoutVersions);
    on<FetchLayoutByVersionEvent>(_onFetchLayoutByVersion);
    on<FetchLatestLayoutEvent>(_onFetchLatestLayout);
    on<FetchActiveBookingsEvent>(_onFetchActiveBookings);
    on<RefreshCoworkingEvent>(_onRefreshCoworking);

    /// =======================
    /// COMMAND EVENTS
    /// =======================

    on<CreateCoworkingEvent>(_onCreateCoworking);
    on<UpdateCoworkingEvent>(_onUpdateCoworking);
    on<DeactivateCoworkingEvent>(_onDeactivateCoworking);
    on<ActivateCoworkingEvent>(_onActivateCoworking);
    on<AddPlacesEvent>(_onAddPlaces);
    on<SetPlaceActiveEvent>(_onSetPlaceActive);
    on<AdminCancelBookingEvent>(_onAdminCancelBooking);
    on<CreateLayoutEvent>(_onCreateLayout);
    on<DeleteLayoutEvent>(_onDeleteLayout);
    on<SetActiveLayout>(_onSetActiveLayout);

    on<UpdateBookingsFilterEvent>(_onUpdateBookingsFilter);
    on<ChangeBookingsPageEvent>(_onChangeBookingsPage);
  }

  /// =======================
  /// QUERY HANDLERS
  /// =======================

  Future<void> _onRefreshCoworking(
    RefreshCoworkingEvent event,
    Emitter<AdminState> emit,
  ) async {
    try {
      final data = await repository.getCoworkingById(event.id);

      emit(
        state.copyWith(
          selectedCoworking: state.selectedCoworking.copyWith(
            data: data,
            status: LoadStatus.success,
          ),
        ),
      );
    } catch (e) {
      emit(state.copyWith(actionMessage: e.toString(), isError: true));
    }
  }

  Future<void> _onFetchCoworkings(
    FetchCoworkingsEvent event,
    Emitter<AdminState> emit,
  ) async {
    emit(
      state.copyWith(coworkings: const LoadState(status: LoadStatus.loading)),
    );

    try {
      final data = await repository.getCoworkings();

      emit(
        state.copyWith(
          coworkings: LoadState(status: LoadStatus.success, data: data),
        ),
      );
    } catch (e) {
      emit(
        state.copyWith(
          coworkings: LoadState(status: LoadStatus.error, error: e.toString()),
        ),
      );
    }
  }

  Future<void> _onFetchCoworkingById(
    FetchCoworkingByIdEvent event,
    Emitter<AdminState> emit,
  ) async {
    emit(
      state.copyWith(
        selectedCoworking: const LoadState(status: LoadStatus.loading),
        places: const LoadState(),
        layout: const LoadState(),
        layoutVersions: const LoadState(),
        bookings: const LoadState(),
        bookingsFilter: const BookingsFilter(),
      ),
    );

    try {
      final data = await repository.getCoworkingById(event.id);

      emit(
        state.copyWith(
          selectedCoworking: LoadState(status: LoadStatus.success, data: data),
        ),
      );
    } catch (e) {
      emit(
        state.copyWith(
          selectedCoworking: LoadState(
            status: LoadStatus.error,
            error: e.toString(),
          ),
        ),
      );
    }
  }

  Future<void> _onFetchPlaces(
    FetchPlacesEvent event,
    Emitter<AdminState> emit,
  ) async {
    emit(state.copyWith(places: const LoadState(status: LoadStatus.loading)));

    try {
      final places = await repository.getPlaces(event.coworkingId);

      emit(
        state.copyWith(
          places: LoadState(status: LoadStatus.success, data: places),
        ),
      );
    } catch (e) {
      emit(
        state.copyWith(
          places: LoadState(status: LoadStatus.error, error: e.toString()),
        ),
      );
    }
  }

  Future<void> _onFetchLayoutVersions(
    FetchLayoutVersionsEvent event,
    Emitter<AdminState> emit,
  ) async {
    emit(
      state.copyWith(
        layoutVersions: const LoadState(status: LoadStatus.loading),
      ),
    );

    try {
      final versions = await repository.getLayoutVersions(event.coworkingId);

      emit(
        state.copyWith(
          layoutVersions: LoadState(status: LoadStatus.success, data: versions),
        ),
      );
    } catch (e) {
      emit(
        state.copyWith(
          layoutVersions: LoadState(
            status: LoadStatus.error,
            error: e.toString(),
          ),
        ),
      );
    }
  }

  Future<void> _onFetchLayoutByVersion(
    FetchLayoutByVersionEvent event,
    Emitter<AdminState> emit,
  ) async {
    emit(state.copyWith(layout: const LoadState(status: LoadStatus.loading)));

    try {
      final layout = await repository.getLayoutVersion(
        event.coworkingId,
        event.version,
      );
      emit(
        state.copyWith(
          layout: LoadState(status: LoadStatus.success, data: layout),
        ),
      );
    } catch (e) {
      emit(
        state.copyWith(
          layout: LoadState(status: LoadStatus.error, error: e.toString()),
        ),
      );
    }
  }

  Future<void> _onFetchLatestLayout(
    FetchLatestLayoutEvent event,
    Emitter<AdminState> emit,
  ) async {
    emit(state.copyWith(layout: const LoadState(status: LoadStatus.loading)));

    try {
      final layout = await repository.getLatestLayout(event.coworkingId);

      emit(
        state.copyWith(
          layout: LoadState(status: LoadStatus.success, data: layout),
        ),
      );
    } catch (e) {
      emit(
        state.copyWith(
          layout: LoadState(status: LoadStatus.error, error: e.toString()),
        ),
      );
    }
  }

  Future<void> _onFetchActiveBookings(
    FetchActiveBookingsEvent event,
    Emitter<AdminState> emit,
  ) async {
    final coworkingId = event.coworkingId;

    final filter = state.bookingsFilter;

    emit(
      state.copyWith(
        bookings: state.bookings.copyWith(status: LoadStatus.loading),
      ),
    );

    try {
      final data = await repository.getActiveBookings(
        coworkingId: coworkingId,
        page: filter.page,
        pageSize: filter.pageSize,
        date: filter.date,
        placeType: filter.placeType,
        sortBy: filter.sortBy,
      );

      emit(
        state.copyWith(
          bookings: LoadState(status: LoadStatus.success, data: data),
        ),
      );
    } catch (e) {
      emit(
        state.copyWith(
          bookings: state.bookings.copyWith(
            status: LoadStatus.error,
            error: e.toString(),
          ),
        ),
      );
    }
  }

  /// =======================
  /// COMMAND HANDLERS
  /// =======================

  Future<void> _onCreateCoworking(
    CreateCoworkingEvent event,
    Emitter<AdminState> emit,
  ) async {
    try {
      await repository.createCoworking(event.name, event.address);

      emit(state.copyWith(actionMessage: 'Coworking created', isError: false));

      add(FetchCoworkingsEvent());
    } catch (e) {
      emit(
        state.copyWith(actionMessage: 'Error: ${e.toString()}', isError: true),
      );
    }
  }

  Future<void> _onUpdateCoworking(
    UpdateCoworkingEvent event,
    Emitter<AdminState> emit,
  ) async {
    try {
      await repository.updateCoworking(event.id, event.name, event.address);

      emit(state.copyWith(actionMessage: 'Coworking updated', isError: false));

      add(FetchCoworkingsEvent());
      add(RefreshCoworkingEvent(event.id));
    } catch (e) {
      emit(
        state.copyWith(actionMessage: 'Error: ${e.toString()}', isError: true),
      );
    }
  }

  Future<void> _onDeactivateCoworking(
    DeactivateCoworkingEvent event,
    Emitter<AdminState> emit,
  ) async {
    try {
      await repository.setCoworkingActive(event.id, false);

      emit(
        state.copyWith(actionMessage: 'Coworking deactivated', isError: false),
      );

      add(FetchCoworkingsEvent());
      add(RefreshCoworkingEvent(event.id));
    } catch (e) {
      emit(
        state.copyWith(actionMessage: 'Error: ${e.toString()}', isError: true),
      );
    }
  }

  Future<void> _onActivateCoworking(
    ActivateCoworkingEvent event,
    Emitter<AdminState> emit,
  ) async {
    try {
      await repository.setCoworkingActive(event.id, true);

      emit(
        state.copyWith(actionMessage: 'Coworking activated', isError: false),
      );

      add(FetchCoworkingsEvent());
      add(RefreshCoworkingEvent(event.id));
    } catch (e) {
      emit(
        state.copyWith(actionMessage: 'Error: ${e.toString()}', isError: true),
      );
    }
  }

  Future<void> _onAddPlaces(
    AddPlacesEvent event,
    Emitter<AdminState> emit,
  ) async {
    try {
      await repository.addPlaces(event.coworkingId, event.places);

      emit(state.copyWith(actionMessage: 'Places added', isError: false));

      add(FetchPlacesEvent(event.coworkingId));
    } catch (e) {
      emit(
        state.copyWith(actionMessage: 'Error: ${e.toString()}', isError: true),
      );
    }
  }

  Future<void> _onSetPlaceActive(
    SetPlaceActiveEvent event,
    Emitter<AdminState> emit,
  ) async {
    try {
      await repository.setPlaceActive(event.placeId, event.active);

      emit(state.copyWith(actionMessage: 'Place updated', isError: false));

      if (state.selectedCoworking.data != null) {
        add(FetchPlacesEvent(state.selectedCoworking.data!.id));
      }
    } catch (e) {
      emit(
        state.copyWith(actionMessage: 'Error: ${e.toString()}', isError: true),
      );
    }
  }

  Future<void> _onAdminCancelBooking(
    AdminCancelBookingEvent event,
    Emitter<AdminState> emit,
  ) async {
    try {
      await repository.adminCancelBooking(
        event.bookingId,
        reason: event.reason,
      );

      emit(state.copyWith(actionMessage: 'Booking cancelled', isError: false));

      final coworkingId = state.selectedCoworking.data?.id;
      if (coworkingId != null) {
        add(FetchActiveBookingsEvent(coworkingId: coworkingId));
      }
    } catch (e) {
      emit(
        state.copyWith(actionMessage: 'Error: ${e.toString()}', isError: true),
      );
    }
  }

  Future<void> _onCreateLayout(
    CreateLayoutEvent event,
    Emitter<AdminState> emit,
  ) async {
    try {
      await repository.createLayout(
        event.coworkingId,
        event.layout,
        event.version,
      );

      emit(state.copyWith(actionMessage: 'Layout created', isError: false));

      add(FetchLayoutVersionsEvent(event.coworkingId));
      add(FetchLatestLayoutEvent(event.coworkingId));
    } catch (e) {
      emit(
        state.copyWith(actionMessage: 'Error: ${e.toString()}', isError: true),
      );
    }
  }

  Future<void> _onDeleteLayout(
    DeleteLayoutEvent event,
    Emitter<AdminState> emit,
  ) async {
    try {
      await repository.deleteLayout(event.coworkingId, event.version);

      emit(state.copyWith(actionMessage: 'Layout deleted', isError: false));

      add(FetchLayoutVersionsEvent(event.coworkingId));
    } catch (e) {
      emit(
        state.copyWith(actionMessage: 'Error: ${e.toString()}', isError: true),
      );
    }
  }

  Future<void> _onSetActiveLayout(
    SetActiveLayout event,
    Emitter<AdminState> emit,
  ) async {
    try {
      await repository.setActiveLayout(event.coworkingId, event.version);

      emit(
        state.copyWith(actionMessage: 'Active layout updated', isError: false),
      );

      add(FetchLayoutVersionsEvent(event.coworkingId));
      add(FetchLatestLayoutEvent(event.coworkingId));
    } catch (e) {
      emit(
        state.copyWith(actionMessage: 'Error: ${e.toString()}', isError: true),
      );
    }
  }

  /// =======================
  /// USERS
  /// =======================

  Future<void> _onFetchUsers(
    FetchUsersEvent event,
    Emitter<AdminState> emit,
  ) async {
    emit(state.copyWith(users: const LoadState(status: LoadStatus.loading)));

    try {
      final data = await repository.getUsers(
        search: event.search,
        page: event.page,
        size: event.size,
        role: event.role,
        isActive: event.isActive,
        sort: event.sort,
      );

      emit(
        state.copyWith(
          users: LoadState(status: LoadStatus.success, data: data),
        ),
      );
    } catch (e) {
      emit(
        state.copyWith(
          users: LoadState(status: LoadStatus.error, error: e.toString()),
        ),
      );
    }
  }

  Future<void> _onSelectUser(
    SelectUserEvent event,
    Emitter<AdminState> emit,
  ) async {
    emit(
      state.copyWith(selectedUser: const LoadState(status: LoadStatus.loading)),
    );

    try {
      final user = await repository.getUserById(event.userId);

      emit(
        state.copyWith(
          selectedUser: LoadState(status: LoadStatus.success, data: user),
        ),
      );
    } catch (e) {
      emit(
        state.copyWith(
          selectedUser: LoadState(
            status: LoadStatus.error,
            error: e.toString(),
          ),
        ),
      );
    }
  }

  Future<void> _onUpdateUserRoles(
    UpdateUserRolesEvent event,
    Emitter<AdminState> emit,
  ) async {
    try {
      await repository.updateUserRoles(event.userId, event.roles);

      emit(state.copyWith(actionMessage: 'Roles updated', isError: false));

      add(SelectUserEvent(event.userId));
      add(FetchUsersEvent());
    } catch (e) {
      emit(state.copyWith(actionMessage: e.toString(), isError: true));
    }
  }

  Future<void> _onSetUserActive(
    SetUserActiveEvent event,
    Emitter<AdminState> emit,
  ) async {
    try {
      await repository.setUserActive(event.userId, event.active);

      emit(state.copyWith(actionMessage: 'User updated', isError: false));

      add(SelectUserEvent(event.userId));
      add(FetchUsersEvent());
    } catch (e) {
      emit(state.copyWith(actionMessage: e.toString(), isError: true));
    }
  }

  Future<void> _onUpdateBookingsFilter(
    UpdateBookingsFilterEvent event,
    Emitter<AdminState> emit,
  ) async {
    final newFilter = state.bookingsFilter.copyWith(
      date: event.date,
      clearDate: event.clearDate,
      placeType: event.placeType,
      clearPlaceType: event.clearPlaceType,
      sortBy: event.sortBy,
      resetPage: true,
    );

    emit(state.copyWith(bookingsFilter: newFilter));

    final coworkingId = state.selectedCoworking.data?.id;
    if (coworkingId != null) {
      add(FetchActiveBookingsEvent(coworkingId: coworkingId));
    }
  }

  Future<void> _onChangeBookingsPage(
    ChangeBookingsPageEvent event,
    Emitter<AdminState> emit,
  ) async {
    final newFilter = state.bookingsFilter.copyWith(page: event.page);

    emit(state.copyWith(bookingsFilter: newFilter));

    final coworkingId = state.selectedCoworking.data?.id;
    if (coworkingId != null) {
      add(FetchActiveBookingsEvent(coworkingId: coworkingId));
    }
  }
}
