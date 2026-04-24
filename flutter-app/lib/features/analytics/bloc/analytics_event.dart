abstract class AnalyticsEvent {}

class LoadHourlyEvent extends AnalyticsEvent {
  final String coworkingId;
  final int weekday;
  LoadHourlyEvent(this.coworkingId, this.weekday);
}

class LoadWeekdayEvent extends AnalyticsEvent {
  final String coworkingId;
  LoadWeekdayEvent(this.coworkingId);
}

class LoadCoworkingHeatmapEvent extends AnalyticsEvent {
  final String coworkingId;
  LoadCoworkingHeatmapEvent(this.coworkingId);
}

class LoadPlaceHeatmapEvent extends AnalyticsEvent {
  final String placeId;
  LoadPlaceHeatmapEvent(this.placeId);
}