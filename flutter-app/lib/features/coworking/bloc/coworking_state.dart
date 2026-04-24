import '../../../core/models/coworking.dart';
import '../../../core/models/layout.dart';
import '../../../core/models/place.dart';
import '../../../core/utils/bloc_load_state.dart';

enum CoworkingView { list, details }

class CoworkingState {
  final LoadState<List<Coworking>> coworkings;

  final LoadState<Coworking> selectedCoworking;
  final Place? selectedPlace;

  final LoadState<Layout> layout;
  final LoadState<List<Place>> places;
  final LoadState<List<Place>> availablePlaces;

  final LoadState bookingResult;

  final DateTime? selectedStart;
  final DateTime? selectedEnd;

  final CoworkingView currentView;

  final String? actionMessage;
  final String? messageId;
  final bool isError;

  const CoworkingState({
    this.coworkings = const LoadState(),
    this.selectedCoworking = const LoadState(),
    this.selectedPlace,
    this.layout = const LoadState(),
    this.places = const LoadState(),
    this.availablePlaces = const LoadState(),
    this.bookingResult = const LoadState(),
    this.currentView = CoworkingView.list,
    this.selectedStart,
    this.selectedEnd,
    this.actionMessage,
    this.messageId,
    this.isError = false,
  });

  CoworkingState copyWith({
    LoadState<List<Coworking>>? coworkings,
    LoadState<Coworking>? selectedCoworking,
    Place? Function()?
    selectedPlace, // Используем функцию для возможности передачи null
    LoadState<Layout>? layout,
    LoadState<List<Place>>? places,
    LoadState<List<Place>>? availablePlaces,
    LoadState? bookingResult,
    CoworkingView? currentView,
    DateTime? Function()? selectedStart,
    DateTime? Function()? selectedEnd,
    String? Function()? actionMessage,
    String? messageId,
    bool? isError,
  }) {
    return CoworkingState(
      coworkings: coworkings ?? this.coworkings,
      selectedCoworking: selectedCoworking ?? this.selectedCoworking,
      selectedPlace: selectedPlace != null
          ? selectedPlace()
          : this.selectedPlace,
      layout: layout ?? this.layout,
      places: places ?? this.places,
      availablePlaces: availablePlaces ?? this.availablePlaces,
      bookingResult: bookingResult ?? this.bookingResult,
      currentView: currentView ?? this.currentView,
      selectedStart: selectedStart != null
          ? selectedStart()
          : this.selectedStart,
      selectedEnd: selectedEnd != null ? selectedEnd() : this.selectedEnd,
      actionMessage: actionMessage != null
          ? actionMessage()
          : this.actionMessage,
      messageId: messageId ?? this.messageId,
      isError: isError ?? this.isError,
    );
  }
}
