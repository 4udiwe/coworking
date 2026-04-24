class HeatmapCell {
  final int weekday; // 0-6 или 1-7 (уточни backend)
  final int hour;    // 0-23
  final int count;

  HeatmapCell({
    required this.weekday,
    required this.hour,
    required this.count,
  });

  factory HeatmapCell.fromJson(Map<String, dynamic> json) => HeatmapCell(
    weekday: json['weekday'],
    hour: json['hour'],
    count: json['count'],
  );
}