import 'package:fl_chart/fl_chart.dart';
import 'package:flutter/material.dart';

class WeekdayChart extends StatelessWidget {
  final Map<int, int> data;
  final int selectedIndex;
  final int currentWeekday;

  const WeekdayChart({
    super.key,
    required this.data,
    required this.selectedIndex,
    required this.currentWeekday,
  });

  static const days = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс'];

  @override
  Widget build(BuildContext context) {
    // Вычисляем максимальное значение в данных
    final maxVal = data.values.isEmpty ? 0 : data.values.reduce((a, b) => a > b ? a : b);

    // Динамические значения для линий
    final double topValue = maxVal.toDouble();
    final double midValue = (maxVal / 2).roundToDouble();

    // Вычисляем maxY с запасом, чтобы верхняя линия и её метка поместились
    // Если maxVal = 0, ставим дефолтный 5 для красоты пустого графика
    final double chartMaxY = maxVal == 0 ? 5.0 : maxVal.toDouble() * 1.3;

    return LayoutBuilder(
      builder: (context, constraints) {
        // Рассчитываем ширину одного столбика исходя из ширины экрана.
        final barWidth = constraints.maxWidth / 8.5;

        return SizedBox(
          height: double.infinity, // Высота для размещения столбиков и подписей
          child: BarChart(
            BarChartData(
              extraLinesData: ExtraLinesData(
                horizontalLines: maxVal > 0
                  ? [
                      HorizontalLine(
                        y: midValue,
                        color: Colors.grey.withOpacity(0.3),
                        strokeWidth: 2,
                        dashArray: [4, 4],
                        label: HorizontalLineLabel(
                          show: true,
                          alignment: Alignment.topLeft,
                          labelResolver: (line) => midValue.toInt().toString(),
                          style: TextStyle(color: Colors.grey[500], fontSize: 10, fontWeight: FontWeight.bold),
                        ),
                      ),
                      HorizontalLine(
                        y: topValue,
                        color: Colors.grey.withOpacity(0.3),
                        strokeWidth: 2,
                        dashArray: [4, 4],
                        label: HorizontalLineLabel(
                          show: true,
                          alignment: Alignment.topLeft,
                          labelResolver: (line) => topValue.toInt().toString(),
                          style: TextStyle(color: Colors.grey[500], fontSize: 10, fontWeight: FontWeight.bold),
                        ),
                      ),
                    ]
                  : [],
              ),
              maxY: chartMaxY,
              alignment: BarChartAlignment.spaceAround,
              gridData: const FlGridData(show: false), // Отключаем стандартную сетку
              borderData: FlBorderData(show: false),
              titlesData: const FlTitlesData(
                leftTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
                rightTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
                topTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
                bottomTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
              ),
              barGroups: List.generate(7, (i) {
                // Маппим i (0..6) в ключ дня недели (1..7), начиная с currentWeekday
                final weekdayKey = (currentWeekday + i - 1) % 7 + 1;
                final value = data[weekdayKey] ?? 0;

                return BarChartGroupData(
                  x: i,
                  barRods: [
                    BarChartRodData(
                      toY: value.toDouble(),
                      width: barWidth, // Адаптивная ширина
                      borderRadius: const BorderRadius.vertical(top: Radius.circular(6)),
                      color: i == selectedIndex ? Colors.blue : Colors.blue.withOpacity(0.2),
                      backDrawRodData: BackgroundBarChartRodData(
                        show: true,
                        toY: chartMaxY,
                        color: Colors.grey.withOpacity(0.05),
                      ),
                    ),
                  ],
                );
              }),
            ),
            swapAnimationDuration: const Duration(milliseconds: 400),
            swapAnimationCurve: Curves.easeInOutCubic,
          ),
        );
      },
    );
  }
}
