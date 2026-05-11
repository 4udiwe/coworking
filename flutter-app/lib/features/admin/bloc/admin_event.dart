import 'package:coworking_app/core/models/layout.dart';
import 'package:coworking_app/core/models/place.dart';

import 'admin_state.dart';

abstract class AdminEvent {}

class CreateCoworkingEvent extends AdminEvent {
  final String name;
  final String address;
  final List<String> mediaIDs;
  CreateCoworkingEvent(this.name, this.address, this.mediaIDs);
}

class UpdateCoworkingEvent extends AdminEvent {
  final String id;
  final String name;
  final String address;
  UpdateCoworkingEvent(this.id, this.name, this.address);
}

class DeactivateCoworkingEvent extends AdminEvent {
  final String id;

  DeactivateCoworkingEvent(this.id);
}

class ActivateCoworkingEvent extends AdminEvent {
  final String id;
  ActivateCoworkingEvent(this.id);
}

class AddPlacesEvent extends AdminEvent {
  final String coworkingId;
  final List<Place> places;
  AddPlacesEvent(this.coworkingId, this.places);
}

class SetPlaceActiveEvent extends AdminEvent {
  final String placeId;
  final bool active;
  SetPlaceActiveEvent(this.placeId, this.active);
}

class AdminCancelBookingEvent extends AdminEvent {
  final String bookingId;
  final String? reason;
  AdminCancelBookingEvent(this.bookingId, {this.reason});
}

class CreateLayoutEvent extends AdminEvent {
  final String coworkingId;
  final LayoutSchema layout;
  final int version;

  CreateLayoutEvent(this.coworkingId, this.layout, this.version);
}

class DeleteLayoutEvent extends AdminEvent {
  final String coworkingId;
  final int version;

  DeleteLayoutEvent(this.coworkingId, this.version);
}

class FetchLayoutVersions extends AdminEvent {
  final String coworkingId;
  FetchLayoutVersions(this.coworkingId);
}

class FetchLayoutByVersion extends AdminEvent {
  final String coworkingId;
  final int version;
  FetchLayoutByVersion(this.coworkingId, this.version);
}

class SetActiveLayout extends AdminEvent {
  final String coworkingId;
  final int version;
  SetActiveLayout(this.coworkingId, this.version);
}

class FetchPlaces extends AdminEvent {
  final String coworkingId;
  FetchPlaces(this.coworkingId);
}

class FetchCoworkingsEvent extends AdminEvent {
  FetchCoworkingsEvent();
}

class FetchCoworkingByIdEvent extends AdminEvent {
  final String id;
  FetchCoworkingByIdEvent(this.id);
}

class FetchPlacesEvent extends AdminEvent {
  final String coworkingId;
  FetchPlacesEvent(this.coworkingId);
}

class FetchLayoutVersionsEvent extends AdminEvent {
  final String coworkingId;
  FetchLayoutVersionsEvent(this.coworkingId);
}

class FetchLayoutByVersionEvent extends AdminEvent {
  final String coworkingId;
  final int version;

  FetchLayoutByVersionEvent(this.coworkingId, this.version);
}

class FetchLatestLayoutEvent extends AdminEvent {
  final String coworkingId;
  FetchLatestLayoutEvent(this.coworkingId);
}

class FetchActiveBookingsEvent extends AdminEvent {
  final String coworkingId;
  // final int? page;
  // final int? pageSize;
  // final DateTime? date;
  // final String? placeType;
  // final String? sortBy;

  FetchActiveBookingsEvent({
    required this.coworkingId,
    // this.page ,
    // this.pageSize,
    // this.date,
    // this.placeType,
    // this.sortBy = 'desc',
  });
}

class LoadLayoutFromJsonEvent extends AdminEvent {
  final String jsonString;
  LoadLayoutFromJsonEvent(this.jsonString);
}

class RefreshCoworkingEvent extends AdminEvent {
  final String id;
  RefreshCoworkingEvent(this.id);
}

class SetAdminViewEvent extends AdminEvent {
  final AdminView view;

  SetAdminViewEvent(this.view);
}

class FetchUsersEvent extends AdminEvent {
  final String? search;
  final int page;
  final int size;
  final String? role;
  final bool? isActive;
  final String? sort;

  FetchUsersEvent({
    this.search,
    this.page = 1,
    this.size = 20,
    this.role,
    this.isActive,
    this.sort,
  });
}

class SelectUserEvent extends AdminEvent {
  final String userId;
  SelectUserEvent(this.userId);
}

class UpdateUserRolesEvent extends AdminEvent {
  final String userId;
  final List<String> roles;

  UpdateUserRolesEvent(this.userId, this.roles);
}

class SetUserActiveEvent extends AdminEvent {
  final String userId;
  final bool active;

  SetUserActiveEvent(this.userId, this.active);
}

class UpdateBookingsFilterEvent extends AdminEvent {
  final DateTime? date;
  final bool clearDate;
  final String? placeType;
  final bool clearPlaceType;
  final String? sortBy;

  UpdateBookingsFilterEvent({
    this.date,
    this.clearDate = false,
    this.placeType,
    this.clearPlaceType = false,
    this.sortBy,
  });
}

class ChangeBookingsPageEvent extends AdminEvent {
  final int page;

  ChangeBookingsPageEvent(this.page);
}