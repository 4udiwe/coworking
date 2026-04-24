class HourlyLoad {
  final Map<int, int> load;

  HourlyLoad({required this.load});

  factory HourlyLoad.fromJson(Map<String, dynamic> json) {
    final raw = json['load'] as Map<String, dynamic>;
    return HourlyLoad(
      load: raw.map((k, v) => MapEntry(int.parse(k), v as int)),
    );
  }
}

class WeekdayLoad {
  final Map<int, int> load;

  WeekdayLoad({required this.load});

  factory WeekdayLoad.fromJson(Map<String, dynamic> json) {
    final raw = json['load'] as Map<String, dynamic>;
    return WeekdayLoad(
      load: raw.map((k, v) => MapEntry(int.parse(k), v as int)),
    );
  }
}