import 'dart:typed_data';

abstract class CreateCoworkingEvent {}

class AddMediaFilesEvent extends CreateCoworkingEvent {
  final List<({Uint8List bytes, String name})> files;
  AddMediaFilesEvent(this.files);
}

class RemoveMediaItemEvent extends CreateCoworkingEvent {
  final String localId;
  RemoveMediaItemEvent(this.localId);
}

class SubmitCreateCoworkingEvent extends CreateCoworkingEvent {
  final String name;
  final String address;
  SubmitCreateCoworkingEvent(this.name, this.address);
}
