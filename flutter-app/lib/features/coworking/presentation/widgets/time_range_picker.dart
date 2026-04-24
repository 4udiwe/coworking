import 'package:coworking_app/core/utils/bloc_load_state.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:intl/intl.dart';
import 'package:syncfusion_flutter_sliders/sliders.dart';
import 'package:syncfusion_flutter_core/theme.dart';
import '../../../analytics/bloc/analytics_bloc.dart';
import '../../../analytics/bloc/analytics_event.dart';
import '../../../analytics/bloc/analytics_state.dart';
import '../../../analytics/widgets/hourly_chart.dart';
import '../../../analytics/widgets/weekday_chart.dart';
import '../../bloc/coworking_bloc.dart';
import '../../bloc/coworking_event.dart';
import '../../bloc/coworking_state.dart';

class TimeRangePicker extends StatefulWidget {
  const TimeRangePicker({super.key});

  @override
  State<TimeRangePicker> createState() => _TimeRangePickerState();
}

class _TimeRangePickerState extends State<TimeRangePicker> {
  static const int startHour = 9;
  static const int endHour = 18;

  late SfRangeValues _values;
  late DateTime _selectedDate;
  late int currentWeekday;
  int _selectedIndex = 0;
  bool _isHourlyExpanded = false;

  @override
  void initState() {
    super.initState();
    final now = DateTime.now();
    // Если сейчас 17:00 или позже, бронирование на сегодня недоступно
    final bool isLate = now.hour >= 17;

    // Инициализируем дату: сегодня (если до 17:00) или завтра
    _selectedDate = DateTime(
      now.year,
      now.month,
      now.day,
    ).add(Duration(days: isLate ? 1 : 0));

    currentWeekday = _selectedDate.weekday;

    // Инициализируем диапазон времени
    final minH = _getMinHour(_selectedDate);
    // Инициализируем диапазон ближайшим доступным часом
    _values = SfRangeValues(minH, (minH + 1).clamp(minH, endHour.toDouble()));
  }

  double _getMinHour(DateTime date) {
    final now = DateTime.now();
    if (_isSameDate(date, now)) {
      // Бронирование доступно только со следующего целого часа
      final nextHour = now.hour + 1;
      return nextHour.toDouble().clamp(
        startHour.toDouble(),
        endHour.toDouble(),
      );
    }
    return startHour.toDouble();
  }

  @override
  Widget build(BuildContext context) {
    final now = DateTime.now();
    final bool isLate = now.hour >= 17;

    // Генерируем список дат. Если уже поздно, начинаем с завтрашнего дня.
    final dates = List.generate(7, (i) {
      final date = now.add(Duration(days: i + (isLate ? 1 : 0)));
      return DateTime(date.year, date.month, date.day);
    });

    return BlocBuilder<CoworkingBloc, CoworkingState>(
      builder: (context, state) {
        final minAllowed = _getMinHour(_selectedDate);

        // Настройки стиля
        const double trackHeight = 16.0; // Делаем полоску шире
        final Color unavailableColor =
            Colors.grey.shade400; // Темно-серый для недоступного
        final Color availableColor =
            Colors.grey.shade200; // Светло-серый для доступного
        final Color activeColor = Colors.blue;

        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            if (isLate)
              const Padding(
                padding: EdgeInsets.fromLTRB(16, 0, 16, 8),
                child: Text(
                  "Бронирование на сегодня уже недоступно (после 17:00)",
                  style: TextStyle(
                    color: Colors.orange,
                    fontSize: 13,
                    fontWeight: FontWeight.w500,
                  ),
                ),
              ),

            /// Загрузка по дням недели
            BlocBuilder<AnalyticsBloc, AnalyticsState>(
              builder: (context, state) {
                final weekdayState = state.weekdayState;

                return AnimatedSwitcher(
                  duration: const Duration(milliseconds: 300),
                  child: _buildWeekdayChartContent(weekdayState),
                );
              },
            ),

            /// 🔷 Выбор даты
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 8),
              child: Row(
                children: dates.map((date) {
                  final isSelected = _isSameDate(date, _selectedDate);

                  return Expanded(
                    child: GestureDetector(
                      onTap: () {
                        setState(() {
                          _selectedDate = date;
                          _selectedIndex = dates.indexOf(date);

                          // При смене даты сбрасываем интервал на "ближайший доступный"
                          final minH = _getMinHour(date);
                          _values = SfRangeValues(
                            minH,
                            (minH + 1).clamp(minH, endHour.toDouble()),
                          );
                        });

                        final start = _toDateTime(_values.start, date);
                        final end = _toDateTime(_values.end, date);

                        context.read<CoworkingBloc>().add(
                          SelectTimeRange(start, end),
                        );
                        context.read<AnalyticsBloc>().add(
                          LoadHourlyEvent(
                            state.selectedCoworking.data!.id,
                            date.weekday,
                          ),
                        );
                      },
                      child: AnimatedContainer(
                        duration: const Duration(milliseconds: 250),
                        margin: const EdgeInsets.symmetric(
                          horizontal: 4,
                          vertical: 6,
                        ),
                        padding: const EdgeInsets.symmetric(vertical: 8),
                        decoration: BoxDecoration(
                          color: isSelected ? Colors.blue : Colors.grey[200],
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Text(
                              _weekdayShort(date),
                              style: TextStyle(
                                color: isSelected ? Colors.white : Colors.black,
                                fontWeight: FontWeight.bold,
                                fontSize: 12,
                              ),
                            ),
                            Text(
                              DateFormat('d MMM', 'ru').format(date),
                              style: TextStyle(
                                color: isSelected ? Colors.white : Colors.black,
                                fontSize: 14,
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
                  );
                }).toList(),
              ),
            ),

            const SizedBox(height: 16),

            /// Загрузка по часам
            BlocBuilder<AnalyticsBloc, AnalyticsState>(
              builder: (context, state) {
                final hourlyState = state.hourlyState;

                if (hourlyState.status == LoadStatus.loading) {
                  return const SizedBox(
                    height: 40,
                    child: Center(child: CircularProgressIndicator()),
                  );
                }

                if (hourlyState.status == LoadStatus.success &&
                    hourlyState.data != null &&
                    hourlyState.data!.load.isNotEmpty) {
                  return Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 8),
                    child: Column(
                      children: [
                        Center(
                          child: TextButton(
                            onPressed: () {
                              setState(() {
                                _isHourlyExpanded = !_isHourlyExpanded;
                              });
                            },
                            child: Text(
                              _isHourlyExpanded ? "Скрыть график" : "Показать почасовую загрузку",
                              style: const TextStyle(
                                color: Colors.blue,
                                fontSize: 12,
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                          ),
                        ),
                        AnimatedCrossFade(
                          firstChild: const SizedBox(width: double.infinity),
                          secondChild: Padding(
                            padding: const EdgeInsets.symmetric(horizontal: 8),
                            child: HourlyHistogram(data: hourlyState.data!.load),
                          ),
                          crossFadeState: _isHourlyExpanded
                              ? CrossFadeState.showSecond
                              : CrossFadeState.showFirst,
                          duration: const Duration(milliseconds: 300),
                          sizeCurve: Curves.easeInOut,
                        ),
                        const SizedBox(height: 8),
                      ],
                    ),
                  );
                }

                return const SizedBox.shrink();
              },
            ),

            /// 🔷 Текущий диапазон времени
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: Row(
                children: [
                  Text(
                    "${_format(_values.start)} — ${_format(_values.end)}",
                    style: const TextStyle(
                      fontSize: 14,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  if (minAllowed > startHour)
                    Row(
                      children: [
                        const SizedBox(width: 16),
                        Container(
                          width: 30,
                          height: 16,
                          decoration: BoxDecoration(
                            color: unavailableColor,
                            borderRadius: BorderRadius.circular(10),
                          ),
                        ),
                        const SizedBox(width: 4),
                        const Text(
                          "прошедшее время\n(недоступно для бронирования)",
                          style: TextStyle(fontSize: 12, color: Colors.grey),
                          maxLines: 2,
                        ),
                      ],
                    ),
                ],
              ),
            ),

            /// 🔷 Time range slider
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: SizedBox(
                height: 90.0,
                child: LayoutBuilder(
                  builder: (context, constraints) {
                    final double totalHours = (endHour - startHour).toDouble();
                    final double passedHours = (minAllowed - startHour).toDouble();

                    const double horizontalPadding = 24.0;
                    final double trackWidth = constraints.maxWidth - (horizontalPadding * 2);
                    final double passedWidth = (passedHours / totalHours) * trackWidth;

                    return Stack(
                      alignment: Alignment.center,
                      children: [
                        Positioned(
                          left: horizontalPadding,
                          top: 28,
                          child: Container(
                            width: trackWidth,
                            height: trackHeight,
                            decoration: BoxDecoration(
                              color: availableColor,
                              borderRadius: BorderRadius.circular(trackHeight / 2),
                            ),
                            child: Row(
                              children: [
                                if (passedWidth > 0)
                                  Container(
                                    width: passedWidth,
                                    height: trackHeight,
                                    decoration: BoxDecoration(
                                      color: unavailableColor,
                                      borderRadius: BorderRadius.circular(trackHeight / 2),
                                    ),
                                  ),
                              ],
                            ),
                          ),
                        ),
                        SfRangeSliderTheme(
                          data: SfRangeSliderThemeData(
                            activeLabelStyle: const TextStyle(
                              color: Colors.blue,
                              fontSize: 14,
                              fontWeight: FontWeight.bold,
                            ),
                            inactiveLabelStyle: const TextStyle(
                              color: Colors.black87,
                              fontSize: 14,
                              fontWeight: FontWeight.w500,
                            ),
                            inactiveTrackColor: Colors.transparent,
                            activeTrackColor: activeColor,
                            activeTrackHeight: trackHeight,
                            inactiveTrackHeight: trackHeight,
                            activeTickColor: Colors.blue,
                            inactiveTickColor: Colors.black26,
                            tickOffset: const Offset(0, 4),
                          ),
                          child: SfRangeSlider(
                            min: startHour.toDouble(),
                            max: endHour.toDouble(),
                            values: _values.copyWith(
                              start: _values.start - 0.001,
                              end: _values.end + 0.001,
                            ),
                            interval: 1,
                            stepSize: 1,
                            showTicks: true,
                            showLabels: true,
                            activeColor: activeColor,
                            edgeLabelPlacement: EdgeLabelPlacement.auto,
                            labelPlacement: LabelPlacement.onTicks,
                            onChanged: (SfRangeValues values) {
                              setState(() => _values = _normalize(values));
                            },
                            onChangeEnd: (SfRangeValues values) {
                              final normalized = _normalize(values);
                              context.read<CoworkingBloc>().add(
                                SelectTimeRange(
                                  _toDateTime(normalized.start, _selectedDate),
                                  _toDateTime(normalized.end, _selectedDate),
                                ),
                              );
                            },
                          ),
                        ),
                      ],
                    );
                  },
                ),
              ),
            ),
          ],
        );
      },
    );
  }

  Widget _buildWeekdayChartContent(dynamic weekdayState) {
    if (weekdayState.status == LoadStatus.loading) {
      return const SizedBox(
        key: ValueKey('loading'),
        height: 120,
        child: Center(child: CircularProgressIndicator()),
      );
    }

    if (weekdayState.status == LoadStatus.success &&
        weekdayState.data != null &&
        weekdayState.data!.load.isNotEmpty) {
      return Padding(
        key: const ValueKey('chart'),
        padding: const EdgeInsets.symmetric(horizontal: 8),
        child: AnimatedContainer(
          duration: const Duration(milliseconds: 300),
          curve: Curves.easeInOut,
          height: _isHourlyExpanded ? 60 : 120,
          child: WeekdayChart(
            data: weekdayState.data!.load,
            selectedIndex: _selectedIndex,
            currentWeekday: currentWeekday,
          ),
        ),
      );
    }

    return const SizedBox.shrink(key: ValueKey('empty'));
  }

  SfRangeValues _normalize(SfRangeValues values) {
    final minAllowed = _getMinHour(_selectedDate);

    double start = values.start.round().toDouble();
    double end = values.end.round().toDouble();

    if (start < minAllowed) start = minAllowed;
    if (end < start + 1) end = start + 1;

    if (end - start > 3) end = start + 3;

    if (end > endHour) {
      end = endHour.toDouble();
      start = (end - 1).clamp(minAllowed, endHour.toDouble());
    }

    if (end - start > 3) start = end - 3;
    if (start < minAllowed) start = minAllowed;

    return SfRangeValues(start, end);
  }

  DateTime _toDateTime(double hour, DateTime date) {
    return date.copyWith(
      hour: hour.toInt(),
      minute: 0,
      second: 0,
      millisecond: 0,
    );
  }

  String _format(double hour) {
    final h = hour.toInt().toString().padLeft(2, '0');
    return "$h:00";
  }

  String _weekdayShort(DateTime date) {
    const weekdays = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс'];
    return weekdays[date.weekday - 1];
  }

  bool _isSameDate(DateTime a, DateTime b) {
    return a.year == b.year && a.month == b.month && a.day == b.day;
  }
}
