import 'package:bloc/bloc.dart';
import 'package:coworking_app/core/utils/bloc_load_state.dart';
import '../data/analytics_repository.dart';
import 'analytics_event.dart';
import 'analytics_state.dart';

class AnalyticsBloc extends Bloc<AnalyticsEvent, AnalyticsState> {
  final AnalyticsRepository repository;

  AnalyticsBloc({required this.repository}) : super(const AnalyticsState()) {
    on<LoadHourlyEvent>(_onLoadHourly);
    on<LoadWeekdayEvent>(_onLoadWeekday);
    on<LoadCoworkingHeatmapEvent>(_onLoadCoworkingHeatmap);
    on<LoadPlaceHeatmapEvent>(_onLoadPlaceHeatmap);
  }

  Future<void> _onLoadHourly(
      LoadHourlyEvent event,
      Emitter<AnalyticsState> emit,
      ) async {
    try {
      emit(state.copyWith(hourlyState: LoadState(status: LoadStatus.loading)));

      final data = await repository.getHourlyLoad(event.coworkingId, event.weekday);

      emit(state.copyWith(
        hourlyState: LoadState(
          data: data,
          status: LoadStatus.success,
      )));
    } catch (e) {
      emit(state.copyWith(
        hourlyState: LoadState(
          status: LoadStatus.error,
          error: e.toString(),
        ),
      ));
    }
  }

  Future<void> _onLoadWeekday(
      LoadWeekdayEvent event,
      Emitter<AnalyticsState> emit,
      ) async {
    try {
      emit(state.copyWith(weekdayState: LoadState(status: LoadStatus.loading)));

      final data = await repository.getWeekdayLoad(event.coworkingId);

      emit(state.copyWith(
          weekdayState: LoadState(
            data: data,
            status: LoadStatus.success,
          )));
    } catch (e) {
      emit(state.copyWith(
        weekdayState: LoadState(
          status: LoadStatus.error,
          error: e.toString(),
        ),
      ));
    }
  }

  Future<void> _onLoadCoworkingHeatmap(
      LoadCoworkingHeatmapEvent event,
      Emitter<AnalyticsState> emit,
      ) async {
    try {
      emit(state.copyWith(
          coworkingHeatmapState: LoadState(status: LoadStatus.loading)));

      final data = await repository.getCoworkingHeatmap(event.coworkingId);

      emit(state.copyWith(
          coworkingHeatmapState: LoadState(
            data: data,
            status: LoadStatus.success,
          )));
    } catch (e) {
      emit(state.copyWith(
        coworkingHeatmapState: LoadState(
          status: LoadStatus.error,
          error: e.toString(),
        ),
      ));
    }
  }

  Future<void> _onLoadPlaceHeatmap(
      LoadPlaceHeatmapEvent event,
      Emitter<AnalyticsState> emit,
      ) async {
    try {
      emit(state.copyWith(
          placeHeatmapState: LoadState(status: LoadStatus.loading)));

      final data = await repository.getPlaceHeatmap(event.placeId);

      emit(state.copyWith(
          placeHeatmapState: LoadState(
            data: data,
            status: LoadStatus.success,
          )));
    } catch (e) {
      emit(state.copyWith(
        placeHeatmapState: LoadState(
          status: LoadStatus.error,
          error: e.toString(),
        ),
      ));
    }
  }
}