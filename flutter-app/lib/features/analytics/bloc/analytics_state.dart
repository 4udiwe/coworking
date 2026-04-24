import 'package:coworking_app/core/models/heatmap_response.dart';
import 'package:coworking_app/core/models/load_response.dart';

import '../../../core/utils/bloc_load_state.dart';

class AnalyticsState {
  final LoadState<HourlyLoad> hourlyState;
  final LoadState<WeekdayLoad> weekdayState;
  final LoadState<CoworkingHeatmap> coworkingHeatmapState;
  final LoadState<PlaceHeatmap> placeHeatmapState;

  const AnalyticsState({
    this.hourlyState = const LoadState(),
    this.weekdayState = const LoadState(),
    this.coworkingHeatmapState = const LoadState(),
    this.placeHeatmapState = const LoadState(),
  });

  AnalyticsState copyWith({
    LoadState<HourlyLoad>? hourlyState,
    LoadState<WeekdayLoad>? weekdayState,
    LoadState<CoworkingHeatmap>? coworkingHeatmapState,
    LoadState<PlaceHeatmap>? placeHeatmapState,
  }) {
    return AnalyticsState(
      hourlyState: hourlyState ?? this.hourlyState,
      weekdayState: weekdayState ?? this.weekdayState,
      coworkingHeatmapState:
          coworkingHeatmapState ?? this.coworkingHeatmapState,
      placeHeatmapState: placeHeatmapState ?? this.placeHeatmapState,
    );
  }
}
