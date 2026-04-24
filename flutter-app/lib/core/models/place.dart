import 'package:coworking_app/core/l10n/localization_helper.dart';
import 'package:flutter/material.dart';

enum PlaceType {
  openDesk,
  meetingRoom;


  static PlaceType fromJson(String value) {
    switch (value) {
      case 'open_desk':
        return PlaceType.openDesk;
      case 'meeting_room':
        return PlaceType.meetingRoom;
      default:
        return PlaceType.openDesk;
    }
  }

  static String toJson(PlaceType value) {
    switch (value) {
      case PlaceType.openDesk:
        return 'open_desk';
      case PlaceType.meetingRoom:
        return 'meeting_room';
    }
  }

  String toTitleCase(BuildContext context) {
    switch (this) {
      case PlaceType.openDesk:
        return context.l10n.placeTypeOpenDesk;
      case PlaceType.meetingRoom:
        return 'Meeting room';
    }
  }
}

class Place {
  String id;
  String coworkingId;
  String coworkingName;
  String label;
  PlaceType placeType;
  bool isActive;
  DateTime createdAt;
  DateTime updatedAt;

  Place({
    this.id = '',
    this.coworkingId = '',
    this.coworkingName = '',
    this.label = '',
    this.placeType = PlaceType.openDesk,
    this.isActive = true,
    DateTime? createdAt,
    DateTime? updatedAt,
  })  : createdAt = createdAt ?? DateTime.now(),
        updatedAt = updatedAt ?? DateTime.now();

  factory Place.fromJson(Map<String, dynamic> json) => Place(
    id: json['id'] ?? '',
    coworkingId: json['coworkingId'] ?? '',
    coworkingName: json['coworkingName'] ?? '',
    label: json['label'] ?? '',
    placeType: PlaceType.fromJson(json['placeType']),
    isActive: json['isActive'] ?? true,
    createdAt: json['createdAt'] != null
        ? DateTime.parse(json['createdAt'])
        : DateTime.now(),
    updatedAt: json['updatedAt'] != null
        ? DateTime.parse(json['updatedAt'])
        : DateTime.now(),
  );

  Map<String, dynamic> toJson() => {
    'id': id,
    'coworkingId': coworkingId,
    'coworkingName': coworkingName,
    'label': label,
    'placeType': PlaceType.toJson(placeType),
    'isActive': isActive,
    'createdAt': createdAt.toIso8601String(),
    'updatedAt': updatedAt.toIso8601String(),
  };
}