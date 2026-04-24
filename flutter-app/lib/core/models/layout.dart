class Layout {
  final String id;
  final String coworkingId;
  final int version;
  final LayoutSchema layout;
  final DateTime createdAt;

  Layout({
    this.id = '',
    this.coworkingId = '',
    this.version = 0,
    required this.layout,
    required this.createdAt,
  });

  factory Layout.fromJson(Map<String, dynamic> json) => Layout(
    id: json['id'],
    coworkingId: json['coworkingId'],
    version: json['version'],
    layout: LayoutSchema.fromJson(Map<String, dynamic>.from(json['layout'])),
    createdAt: DateTime.parse(json['createdAt']),
  );

  Map<String, dynamic> toJson() => {
    'id': id,
    'coworkingId': coworkingId,
    'version': version,
    'layout': layout.toJson(),
    'createdAt': createdAt.toIso8601String(),
  };
}

class LayoutSchema {
  final int formatVersion;
  final CanvasSize canvas;
  final List<WallEntity> walls;
  final List<PlaceEntity> places;

  LayoutSchema({
    required this.formatVersion,
    required this.canvas,
    required this.walls,
    required this.places,
  });

  factory LayoutSchema.fromJson(Map<String, dynamic> json) => LayoutSchema(
    formatVersion: json['formatVersion'],
    canvas: CanvasSize.fromJson(json['canvas']),
    walls: (json['walls'] as List).map((e) => WallEntity.fromJson(e)).toList(),
    places: (json['places'] as List).map((e) => PlaceEntity.fromJson(e)).toList(),
  );

  Map<String, dynamic> toJson() => {
    'formatVersion': formatVersion,
    'canvas': canvas.toJson(),
    'walls': walls.map((e) => e.toJson()).toList(),
    'places': places.map((e) => e.toJson()).toList(),
  };
}

class CanvasSize {
  final double width;
  final double height;

  CanvasSize({required this.width, required this.height});

  factory CanvasSize.fromJson(Map<String, dynamic> json) => CanvasSize(
    width: (json['width'] as num).toDouble(),
    height: (json['height'] as num).toDouble(),
  );

  Map<String, dynamic> toJson() => {'width': width, 'height': height};
}

class WallEntity {
  final String id;
  final double x;
  final double y;
  final double width;
  final double height;
  final double rotation;

  WallEntity({
    required this.id,
    required this.x,
    required this.y,
    required this.width,
    required this.height,
    required this.rotation,
  });

  factory WallEntity.fromJson(Map<String, dynamic> json) => WallEntity(
    id: json['id'],
    x: (json['x'] as num).toDouble(),
    y: (json['y'] as num).toDouble(),
    width: (json['width'] as num).toDouble(),
    height: (json['height'] as num).toDouble(),
    rotation: (json['rotation'] as num).toDouble(),
  );

  Map<String, dynamic> toJson() => {
    'id': id,
    'x': x,
    'y': y,
    'width': width,
    'height': height,
    'rotation': rotation,
  };
}

class PlaceEntity {
  final String id;
  final String type;
  final double x;
  final double y;
  final double rotation;

  final double width;
  final double height;

  PlaceEntity({
    required this.id,
    required this.type,
    required this.x,
    required this.y,
    required this.rotation,
    this.width = 80,
    this.height = 80,
  });

  factory PlaceEntity.fromJson(Map<String, dynamic> json) => PlaceEntity(
    id: json['id'],
    type: json['type'],
    x: (json['x'] as num).toDouble(),
    y: (json['y'] as num).toDouble(),
    rotation: (json['rotation'] as num).toDouble(),
  );

  Map<String, dynamic> toJson() => {
    'id': id,
    'type': type,
    'x': x,
    'y': y,
    'rotation': rotation,
  };
}

class LayoutVersion {
  final int version;
  final DateTime createdAt;

  LayoutVersion({required this.version, required this.createdAt});

  factory LayoutVersion.fromJson(Map<String, dynamic> json) => LayoutVersion(
    version: json['version'],
    createdAt: DateTime.parse(json['createdAt']),
  );
}
