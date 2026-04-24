import 'package:flutter/material.dart';
import '../../../core/models/heatmap_cell.dart';

class HeatmapView extends StatelessWidget {
  final List<HeatmapCell> data;

  // Настройки сетки
  static const int startHour = 9;
  static const int endHour = 17;
  static const List<String> weekDays = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс'];

  const HeatmapView({super.key, required this.data});

  @override
  Widget build(BuildContext context) {
    // 1. Маппинг данных: переводим UTC -> Local и группируем по локальным "день_час"
    final Map<String, int> localDataMap = {};
    
    // Берем любой понедельник в UTC (2024-01-01 был понедельником) для базы расчета
    final baseMonday = DateTime.utc(2024, 1, 1);

    for (var cell in data) {
      // Создаем дату в UTC: Понедельник + смещение дней + час сервера
      final utcDateTime = baseMonday.add(
        Duration(days: cell.weekday - 1, hours: cell.hour),
      );
      
      // Переводим в локальное время пользователя
      final localDateTime = utcDateTime.toLocal();
      
      final int localWd = localDateTime.weekday; // 1 (Mon) - 7 (Sun)
      final int localH = localDateTime.hour;
      
      // Нам интересны только часы в выбранном диапазоне (9-17)
      if (localH >= startHour && localH <= endHour) {
        final key = "${localWd}_$localH";
        // Суммируем, так как несколько UTC часов могут попасть в один локальный час (хотя обычно 1 к 1)
        localDataMap[key] = (localDataMap[key] ?? 0) + cell.count;
      }
    }

    // 2. Находим максимум только среди тех данных, что попали в диапазон, для расчета яркости цвета
    final maxVal = localDataMap.values.isEmpty 
        ? 1 
        : localDataMap.values.reduce((a, b) => a > b ? a : b);

    final int hoursCount = endHour - startHour + 1; // Количество столбцов (9)

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        // Шапка с часами (09:00, 10:00...)
        Padding(
          padding: const EdgeInsets.only(left: 45, bottom: 8),
          child: Row(
            children: List.generate(hoursCount, (i) {
              return Expanded(
                child: Text(
                  '${startHour + i}',
                  textAlign: TextAlign.center,
                  style: const TextStyle(fontSize: 10, color: Colors.grey, fontWeight: FontWeight.bold),
                ),
              );
            }),
          ),
        ),
        
        // Сетка: Дни недели + Квадраты данных (Адаптивная отрисовка по строкам)
        ...List.generate(7, (dayIdx) {
          final int weekday = dayIdx + 1;
          return Padding(
            padding: const EdgeInsets.only(bottom: 3),
            child: Row(
              children: [
                // Название дня недели (фиксированная ширина для выравнивания с шапкой)
                SizedBox(
                  width: 45,
                  child: Text(
                    weekDays[dayIdx],
                    style: const TextStyle(
                      fontSize: 12,
                      fontWeight: FontWeight.bold,
                      color: Colors.blueGrey,
                    ),
                  ),
                ),
                // Строка с часами
                ...List.generate(hoursCount, (hourIdx) {
                  final int hour = startHour + hourIdx;
                  final count = localDataMap["${weekday}_$hour"] ?? 0;

                  return Expanded(
                    child: Padding(
                      padding: EdgeInsets.only(
                        right: hourIdx == hoursCount - 1 ? 0 : 3,
                      ),
                      child: AspectRatio(
                        aspectRatio: 1.4,
                        child: Container(
                          decoration: BoxDecoration(
                            color: count == 0 ? Colors.grey[100] : _color(count, maxVal),
                            borderRadius: BorderRadius.circular(4),
                          ),
                          child: count > 0
                              ? Center(
                                  child: Text(
                                    '$count',
                                    style: const TextStyle(
                                      fontSize: 10,
                                      color: Colors.white,
                                      fontWeight: FontWeight.bold,
                                    ),
                                  ),
                                )
                              : null,
                        ),
                      ),
                    ),
                  );
                }),
              ],
            ),
          );
        }),
      ],
    );
  }

  Color _color(int value, int max) {
    final ratio = value / max;
    // Градиент от голубого (мало) до красного (пик)
    if (ratio < 0.25) return Colors.blue.shade200;
    if (ratio < 0.5) return Colors.blue.shade400;
    if (ratio < 0.75) return Colors.orange.shade400;
    return Colors.red.shade400;
  }
}
