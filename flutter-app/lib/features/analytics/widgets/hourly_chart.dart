import 'package:fl_chart/fl_chart.dart';
import 'package:flutter/material.dart';

class HourlyHistogram extends StatelessWidget {
  final Map<int, int> data;
  final RangeValues? selectedRange;

  const HourlyHistogram({super.key, required this.data, this.selectedRange});

  @override
  Widget build(BuildContext context) {
    if (data.isEmpty) return const SizedBox.shrink();

    // 1. Преобразуем UTC часы от сервера в локальное время пользователя
    final Map<int, int> localData = data.map((utcHour, count) {
      final now = DateTime.now();
      // Создаем объект в UTC и переводим в local, затем забираем час
      final localHour = DateTime.utc(
        now.year,
        now.month,
        now.day,
        utcHour,
      ).toLocal().hour;
      return MapEntry(localHour, count);
    });

    // 2. Сортируем уже локальные данные
    final sortedKeys = localData.keys.toList()..sort();
    final List<FlSpot> spots = sortedKeys.map((hour) {
      return FlSpot(hour.toDouble(), localData[hour]!.toDouble());
    }).toList();

    final maxVal = data.values.reduce((a, b) => a > b ? a : b);
    final double maxY = maxVal == 0 ? 5.0 : maxVal.toDouble() * 1.3;

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16),
      child: SizedBox(
        height: 140,
        child: LineChart(
          LineChartData(
            minX: 9,
            maxX: 17,
            maxY: maxY,
            minY: 0,
            gridData: FlGridData(
              show: true,
              drawHorizontalLine: false, // Горизонтальные не нужны
              drawVerticalLine: true,
              verticalInterval: 1, // Каждая единица по X (каждый час)
              getDrawingVerticalLine: (value) {
                return FlLine(
                  color: Colors.grey.withOpacity(0.7),
                  strokeWidth: 2,
                  dashArray: [2, 8], // Делает линию пунктирной [длина штриха, длина пропуска]
                );
              },
            ),
            borderData: FlBorderData(show: false),
            lineTouchData: const LineTouchData(enabled: false),

            titlesData: FlTitlesData(
              leftTitles: const AxisTitles(
                sideTitles: SideTitles(showTitles: false),
              ),
              rightTitles: const AxisTitles(
                sideTitles: SideTitles(showTitles: false),
              ),
              topTitles: AxisTitles(
                sideTitleAlignment: SideTitleAlignment.outside,
                sideTitles: SideTitles(
                  showTitles: true,
                  interval: 1, // Показывать метку каждый час
                  getTitlesWidget: (value, meta) {
                    return Text(
                      '${value.toInt()}:00',
                      style: const TextStyle(fontSize: 12, color: Colors.grey),
                    );
                  },
                ),
              ),
              bottomTitles: const AxisTitles(
                sideTitles: SideTitles(showTitles: false),
              ),
            ),

            lineBarsData: [
              LineChartBarData(
                spots: spots,
                isCurved: true, // Сглаживание линии
                curveSmoothness: 0.35,
                barWidth: 3,
                color: Colors.blue,
                dotData: const FlDotData(
                  show: false,
                ), // Скрываем точки для чистоты
                belowBarData: BarAreaData(
                  show: true,
                  gradient: LinearGradient(
                    colors: [
                      Colors.blue.withOpacity(0.5),
                      Colors.blue.withOpacity(0.1),
                    ],
                    begin: Alignment.topCenter,
                    end: Alignment.bottomCenter,
                  ),
                ),
              ),
            ],
          ),
          duration: const Duration(milliseconds: 400),
          curve: Curves.easeInOutCubic,
        ),
      ),
    );
  }
}
