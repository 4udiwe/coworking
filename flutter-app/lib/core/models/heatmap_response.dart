import 'heatmap_cell.dart';

class CoworkingHeatmap {
  final List<HeatmapCell> heatmap;

  CoworkingHeatmap({required this.heatmap});

  factory CoworkingHeatmap.fromJson(Map<String, dynamic> json) {
    final list = json['heatmap'] as List;
    return CoworkingHeatmap(
      heatmap: list.map((e) => HeatmapCell.fromJson(e)).toList(),
    );
  }
}

class PlaceHeatmap {
  final List<HeatmapCell> heatmap;

  PlaceHeatmap({required this.heatmap});

  factory PlaceHeatmap.fromJson(Map<String, dynamic> json) {
    final list = json['heatmap'] as List;
    return PlaceHeatmap(
      heatmap: list.map((e) => HeatmapCell.fromJson(e)).toList(),
    );
  }
}