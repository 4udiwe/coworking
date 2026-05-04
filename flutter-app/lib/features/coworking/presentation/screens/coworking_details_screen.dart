import 'package:coworking_app/core/l10n/localization_helper.dart';
import 'package:coworking_app/features/analytics/bloc/analytics_event.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:intl/intl.dart';

import '../../../../core/di/service_locator.dart';
import '../../../../core/models/layout.dart';
import '../../../../core/models/place.dart';
import '../../../../core/widgets/layout_preview.dart';
import '../../../../core/widgets/skeleton.dart';
import '../../../analytics/bloc/analytics_bloc.dart';
import '../../bloc/coworking_bloc.dart';
import '../../bloc/coworking_event.dart';
import '../../bloc/coworking_state.dart';
import '../widgets/time_range_picker.dart';

final _timeFormat = DateFormat('HH:mm');

class CoworkingDetailsScreen extends StatefulWidget {
  final String coworkingId;
  const CoworkingDetailsScreen({super.key, required this.coworkingId});

  @override
  State<CoworkingDetailsScreen> createState() => _CoworkingDetailsScreenState();
}

class _CoworkingDetailsScreenState extends State<CoworkingDetailsScreen> {
  final ScrollController _scrollController = ScrollController();

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }

  void _scrollToInfo(bool isWide) {
    if (!_scrollController.hasClients) return;

    if (isWide) {
      _scrollController.animateTo(
        0,
        duration: const Duration(milliseconds: 300),
        curve: Curves.easeInOut,
      );
    } else {
      final double mapHeight = MediaQuery.of(context).size.height * 0.45;
      if (_scrollController.offset < mapHeight) {
        _scrollController.animateTo(
          mapHeight,
          duration: const Duration(milliseconds: 600),
          curve: Curves.fastOutSlowIn,
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        leading: BackButton(
          onPressed: () {
            context.read<CoworkingBloc>().add(ClearAction());
            Navigator.pop(context);
          },
        ),
        title: BlocBuilder<CoworkingBloc, CoworkingState>(
          buildWhen: (prev, curr) =>
              prev.selectedCoworking.data?.name !=
              curr.selectedCoworking.data?.name,
          builder: (context, state) =>
              Text(state.selectedCoworking.data?.name ?? "Coworking"),
        ),
      ),
      body: BlocBuilder<CoworkingBloc, CoworkingState>(
        builder: (context, state) {
          final layout = state.layout.data;
          final places = state.places.data ?? [];
          final currentCoworking = state.selectedCoworking.data;

          // Проверяем, что загружен именно тот коворкинг, который мы запрашивали
          final isLoading =
              currentCoworking == null ||
              currentCoworking.id != widget.coworkingId ||
              layout == null;

          if (isLoading) {
            return const _DetailsSkeleton();
          }

          return BlocProvider(
            key: ValueKey(currentCoworking.id),
            create: (_) => sl<AnalyticsBloc>()
              ..add(LoadWeekdayEvent(currentCoworking.id))
              ..add(
                LoadHourlyEvent(
                  currentCoworking.id,
                  state.selectedStart!.weekday,
                ),
              ),
            child: MultiBlocListener(
              listeners: [
                BlocListener<CoworkingBloc, CoworkingState>(
                  listenWhen: (prev, curr) =>
                      prev.selectedStart?.weekday !=
                      curr.selectedStart?.weekday,
                  listener: (context, state) {
                    if (state.selectedStart != null) {
                      context.read<AnalyticsBloc>().add(
                        LoadHourlyEvent(
                          widget.coworkingId,
                          state.selectedStart!.weekday,
                        ),
                      );
                    }
                  },
                ),
                BlocListener<CoworkingBloc, CoworkingState>(
                  listenWhen: (prev, curr) =>
                      prev.selectedPlace?.id != curr.selectedPlace?.id &&
                      curr.selectedPlace != null,
                  listener: (context, state) {
                    final isWide = MediaQuery.of(context).size.width > 800;
                    _scrollToInfo(isWide);
                  },
                ),
                BlocListener<CoworkingBloc, CoworkingState>(
                  listenWhen: (prev, curr) =>
                      prev.messageId != curr.actionMessage &&
                      curr.actionMessage != null,
                  listener: (context, state) {
                    ScaffoldMessenger.of(context).showSnackBar(
                      state.isError
                          ? SnackBar(
                              content: Text(context.l10n.errorWhileCreatingBooking),
                              backgroundColor: Colors.red,
                            )
                          : SnackBar(
                              content: Text(context.l10n.bookingCreatedSuccess),
                              backgroundColor: Colors.green,
                            ),
                    );
                  },
                ),
              ],
              child: LayoutBuilder(
                builder: (context, constraints) {
                  final isWide = constraints.maxWidth > 800;
                  if (isWide) {
                    return _DesktopLayout(
                      state,
                      layout,
                      places,
                      _scrollController,
                    );
                  } else {
                    return _MobileLayout(
                      state,
                      layout,
                      places,
                      _scrollController,
                    );
                  }
                },
              ),
            ),
          );
        },
      ),
    );
  }
}

class _DetailsSkeleton extends StatelessWidget {
  const _DetailsSkeleton();

  @override
  Widget build(BuildContext context) {
    final isWide = MediaQuery.of(context).size.width > 800;

    if (isWide) {
      return Row(
        children: [
          const Expanded(
            flex: 3,
            child: Padding(
              padding: EdgeInsets.all(16.0),
              child: Skeleton(height: double.infinity),
            ),
          ),
          Expanded(
            flex: 2,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: const [
                Padding(
                  padding: EdgeInsets.all(16),
                  child: Skeleton(height: 30, width: 200),
                ),
                Padding(
                  padding: EdgeInsets.symmetric(horizontal: 16),
                  child: Skeleton(height: 120),
                ),
                SizedBox(height: 16),
                Padding(
                  padding: EdgeInsets.symmetric(horizontal: 16),
                  child: Skeleton(height: 60),
                ),
                SizedBox(height: 16),
                Padding(
                  padding: EdgeInsets.symmetric(horizontal: 16),
                  child: Skeleton(height: 200),
                ),
              ],
            ),
          ),
        ],
      );
    }

    return SingleChildScrollView(
      child: Column(
        children: [
          Skeleton(height: MediaQuery.of(context).size.height * 0.45),
          const SizedBox(height: 16),
          const Padding(
            padding: EdgeInsets.symmetric(horizontal: 16),
            child: Skeleton(height: 120),
          ),
          const SizedBox(height: 16),
          const Padding(
            padding: EdgeInsets.symmetric(horizontal: 16),
            child: Skeleton(height: 60),
          ),
          const SizedBox(height: 16),
          const Padding(
            padding: EdgeInsets.symmetric(horizontal: 16),
            child: Skeleton(height: 200),
          ),
        ],
      ),
    );
  }
}

class _MobileLayout extends StatelessWidget {
  final CoworkingState state;
  final Layout layout;
  final List<Place> places;
  final ScrollController scrollController;

  const _MobileLayout(
    this.state,
    this.layout,
    this.places,
    this.scrollController,
  );

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      controller: scrollController,
      child: Column(
        children: [
          SizedBox(
            height: MediaQuery.of(context).size.height * 0.45,
            child: LayoutPreview(
              layout: layout,
              places: places,
              selectedPlaceIds: state.selectedPlace != null
                  ? {state.selectedPlace!.id}
                  : {},
              unavailablePlaceIds: _getUnavailable(state),
              onPlaceTap: (place) {
                context.read<CoworkingBloc>().add(SelectPlace(place));
              },
            ),
          ),
          const TimeRangePicker(),
          _BottomPanel(state: state),
        ],
      ),
    );
  }
}

Set<String> _getUnavailable(CoworkingState state) {
  final all = state.places.data ?? [];
  final available = state.availablePlaces.data ?? [];
  final availableIds = available.map((e) => e.id).toSet();

  return all
      .where((p) => !availableIds.contains(p.id))
      .map((e) => e.id)
      .toSet();
}

class _DesktopLayout extends StatelessWidget {
  final CoworkingState state;
  final Layout layout;
  final List<Place> places;
  final ScrollController scrollController;

  const _DesktopLayout(
    this.state,
    this.layout,
    this.places,
    this.scrollController,
  );

  @override
  Widget build(BuildContext context) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Expanded(
          flex: 3,
          child: LayoutPreview(
            layout: layout,
            places: places,
            selectedPlaceIds: state.selectedPlace != null
                ? {state.selectedPlace!.id}
                : {},
            unavailablePlaceIds: _getUnavailable(state),
            onPlaceTap: (place) {
              context.read<CoworkingBloc>().add(SelectPlace(place));
            },
          ),
        ),
        Expanded(
          flex: 2,
          child: SingleChildScrollView(
            controller: scrollController,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisAlignment: MainAxisAlignment.start,
              children: [
                Padding(
                  padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
                  child: Text(
                    state.selectedCoworking.data?.name ?? 'coworking name',
                    style: const TextStyle(
                      fontWeight: FontWeight.bold,
                      fontSize: 22,
                    ),
                    textAlign: TextAlign.start,
                  ),
                ),
                const TimeRangePicker(),
                _BottomPanel(state: state),
              ],
            ),
          ),
        ),
      ],
    );
  }
}

class _BottomPanel extends StatelessWidget {
  final CoworkingState state;

  const _BottomPanel({required this.state});

  @override
  Widget build(BuildContext context) {
    // Чтобы UI не "прыгал", мы резервируем место под панель,
    // либо используем AnimatedSwitcher
    return AnimatedSize(
      duration: const Duration(milliseconds: 300),
      child: Container(
        constraints: const BoxConstraints(minHeight: 100),
        width: double.infinity,
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
        child: state.selectedPlace == null
            ? const Center(
                child: Text(
                  "Выберите место на карте",
                  style: TextStyle(color: Colors.grey, fontSize: 16),
                ),
              )
            : Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisSize: MainAxisSize.min,
                children: [
                  Text(
                    "Место ${state.selectedPlace!.label}",
                    style: const TextStyle(
                      fontSize: 24,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    "Тип места: ${state.selectedPlace!.placeType.toTitleCase(context)}",
                  ),
                  const SizedBox(height: 8),
                  Text(
                    "Доступно на выбранное время ${_timeFormat.format(state.selectedStart!)} - ${_timeFormat.format(state.selectedEnd!)}",
                  ),
                  const SizedBox(height: 24),
                  SizedBox(
                    width: double.infinity,
                    child: ElevatedButton(
                      style: ElevatedButton.styleFrom(
                        backgroundColor: Colors.blue,
                        foregroundColor: Colors.white,
                        padding: const EdgeInsets.symmetric(vertical: 16),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(8),
                        ),
                      ),
                      onPressed: () {
                        context.read<CoworkingBloc>().add(
                          CreateBooking(
                            state.selectedPlace!.id,
                            state.selectedStart!,
                            state.selectedEnd!,
                          ),
                        );
                      },
                      child: const Text(
                        "Забронировать",
                        style: TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
      ),
    );
  }
}

extension StringExtension on String {
  String toTitleCase() {
    if (isEmpty) return this;
    return "${this[0].toUpperCase()}${substring(1).toLowerCase()}";
  }
}
