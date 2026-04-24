import '../../../core/models/place.dart';

abstract class CoworkingEvent {}

/// 📌 список
class FetchCoworkings extends CoworkingEvent {}

/// 📌 выбор коворкинга
class SelectCoworking extends CoworkingEvent {
  final String id;
  SelectCoworking(this.id);
}

/// 📌 загрузка layout
class FetchLayout extends CoworkingEvent {
  final String coworkingId;
  FetchLayout(this.coworkingId);
}

/// 📌 загрузка всех мест
class FetchPlaces extends CoworkingEvent {
  final String coworkingId;
  FetchPlaces(this.coworkingId);
}

/// 📌 доступные места по времени
class FetchAvailablePlaces extends CoworkingEvent {
  final String coworkingId;
  final DateTime start;
  final DateTime end;

  FetchAvailablePlaces(this.coworkingId, this.start, this.end);
}

class SelectPlace extends CoworkingEvent {
  final Place place;
  SelectPlace(this.place);
}

class SelectTimeRange extends CoworkingEvent {
  final DateTime start;
  final DateTime end;

  SelectTimeRange(this.start, this.end);
}

/// 📌 бронирование
class CreateBooking extends CoworkingEvent {
  final String placeId;
  final DateTime start;
  final DateTime end;

  CreateBooking(this.placeId, this.start, this.end);
}

class SelectDate extends CoworkingEvent {
  final DateTime date;

  SelectDate(this.date);
}

/// 📌 сброс сообщений (snackbar и т.д.)
class ClearAction extends CoworkingEvent {}