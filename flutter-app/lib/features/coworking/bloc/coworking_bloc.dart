import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../core/utils/bloc_load_state.dart';
import '../data/coworking_repository.dart';
import 'coworking_event.dart';
import 'coworking_state.dart';

class CoworkingBloc extends Bloc<CoworkingEvent, CoworkingState> {
  final CoworkingRepository repository;

  CoworkingBloc({required this.repository}) : super(const CoworkingState()) {
    on<FetchCoworkings>(_onFetchCoworkings);
    on<SelectCoworking>(_onSelectCoworking);
    on<FetchLayout>(_onFetchLayout);
    on<FetchPlaces>(_onFetchPlaces);
    on<FetchAvailablePlaces>(_onFetchAvailablePlaces);
    on<CreateBooking>(_onCreateBooking);
    on<ClearAction>(_onClearAction);

    on<SelectPlace>((event, emit) {
      final availableIds =
          state.availablePlaces.data?.map((e) => e.id).toSet() ?? {};
      if (!availableIds.contains(event.place.id)) {
        return;
      }

      emit(state.copyWith(selectedPlace: () => event.place));
    });

    on<SelectTimeRange>(_onSelectTimeRange);
  }

  Future<void> _onSelectTimeRange(
    SelectTimeRange event,
    Emitter<CoworkingState> emit,
  ) async {
    emit(
      state.copyWith(
        selectedStart: () => event.start,
        selectedEnd: () => event.end,
        selectedPlace: () => null,
      ),
    );

    if (state.selectedCoworking.data != null) {
      add(
        FetchAvailablePlaces(
          state.selectedCoworking.data!.id,
          event.start,
          event.end,
        ),
      );
    }
  }

  /// =======================
  /// COWORKINGS LIST
  /// =======================
  Future<void> _onFetchCoworkings(
    FetchCoworkings event,
    Emitter<CoworkingState> emit,
  ) async {
    emit(
      state.copyWith(
        coworkings: state.coworkings.copyWith(status: LoadStatus.loading),
      ),
    );

    try {
      final data = await repository.getCoworkings();
      print("Fetched coworkings: ${data.length}");
      print(
        "Coworkings: ${data.map((c) => c.name + "images count ${c.imageIDs.length}").join(', ')}",
      );
      emit(
        state.copyWith(
          coworkings: LoadState(data: data, status: LoadStatus.success),
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

  /// =======================
  /// SELECT COWORKING
  /// =======================
  Future<void> _onSelectCoworking(
    SelectCoworking event,
    Emitter<CoworkingState> emit,
  ) async {
    // При выборе нового коворкинга обязательно очищаем старые данные деталей,
    // чтобы UI не отображал данные предыдущего коворкинга во время загрузки.
    emit(
      state.copyWith(
        selectedCoworking: const LoadState(status: LoadStatus.loading),
        layout: const LoadState(),
        places: const LoadState(),
        availablePlaces: const LoadState(),
        selectedPlace: () => null,
        currentView: CoworkingView.details,
      ),
    );

    try {
      final coworking = await repository.getCoworkingById(event.id);

      // Устанавливаем дефолтное время
      var now = DateTime.now();
      if (now.hour > 17) {
        now = now.add(Duration(days: 1));
      }
      final defaultStart = DateTime(now.year, now.month, now.day, now.hour + 1);
      final defaultEnd = DateTime(now.year, now.month, now.day, now.hour + 2);

      emit(
        state.copyWith(
          selectedCoworking: LoadState(
            data: coworking,
            status: LoadStatus.success,
          ),
          selectedStart: () => defaultStart,
          selectedEnd: () => defaultEnd,
        ),
      );

      add(FetchLayout(event.id));
      add(FetchPlaces(event.id));
      add(FetchAvailablePlaces(event.id, defaultStart, defaultEnd));
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

  /// =======================
  /// LAYOUT
  /// =======================
  Future<void> _onFetchLayout(
    FetchLayout event,
    Emitter<CoworkingState> emit,
  ) async {
    emit(state.copyWith(layout: const LoadState(status: LoadStatus.loading)));

    try {
      final layout = await repository.getLayout(event.coworkingId);
      emit(
        state.copyWith(
          layout: LoadState(data: layout, status: LoadStatus.success),
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

  /// =======================
  /// PLACES
  /// =======================
  Future<void> _onFetchPlaces(
    FetchPlaces event,
    Emitter<CoworkingState> emit,
  ) async {
    emit(state.copyWith(places: const LoadState(status: LoadStatus.loading)));

    try {
      final places = await repository.getPlaces(event.coworkingId);
      emit(
        state.copyWith(
          places: LoadState(data: places, status: LoadStatus.success),
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

  /// =======================
  /// AVAILABLE PLACES
  /// =======================
  Future<void> _onFetchAvailablePlaces(
    FetchAvailablePlaces event,
    Emitter<CoworkingState> emit,
  ) async {
    emit(
      state.copyWith(
        availablePlaces: LoadState(
          status: LoadStatus.loading,
          data: state.availablePlaces.data,
        ),
      ),
    );

    try {
      final result = await repository.getAvailablePlaces(
        event.coworkingId,
        event.start,
        event.end,
      );

      emit(
        state.copyWith(
          availablePlaces: LoadState(
            data: result.places,
            status: LoadStatus.success,
          ),
        ),
      );
    } catch (e) {
      emit(
        state.copyWith(
          availablePlaces: LoadState(
            status: LoadStatus.error,
            error: e.toString(),
          ),
        ),
      );
    }
  }

  /// =======================
  /// BOOKING
  /// =======================
  Future<void> _onCreateBooking(
    CreateBooking event,
    Emitter<CoworkingState> emit,
  ) async {
    emit(
      state.copyWith(
        bookingResult: state.bookingResult.copyWith(status: LoadStatus.loading),
      ),
    );

    try {
      await repository.createBooking(event.placeId, event.start, event.end);

      add(
        FetchAvailablePlaces(
          state.selectedCoworking.data!.id,
          event.start,
          event.end,
        ),
      );

      emit(
        state.copyWith(
          selectedPlace: () => null,
          bookingResult: LoadState(status: LoadStatus.success),
          actionMessage: () => "Booking created successfully",
          messageId: DateTime.now().millisecondsSinceEpoch.toString(),
          isError: false,
        ),
      );
    } catch (e) {
      emit(
        state.copyWith(
          bookingResult: LoadState(
            status: LoadStatus.error,
            error: e.toString(),
          ),
          actionMessage: () => e.toString(),
          messageId: DateTime.now().millisecondsSinceEpoch.toString(),
          isError: true,
        ),
      );
    }
  }

  void _onClearAction(ClearAction event, Emitter<CoworkingState> emit) {
    emit(
      state.copyWith(
        selectedPlace: () => null,
        selectedCoworking: const LoadState(),
        layout: const LoadState(),
        places: const LoadState(),
        availablePlaces: const LoadState(),
        bookingResult: const LoadState(),
        actionMessage: () => null,
        isError: false,
      ),
    );
  }
}
