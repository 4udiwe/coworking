import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:file_picker/file_picker.dart';
import 'package:flutter/services.dart';
import 'package:intl/intl.dart';

import '../../../../core/di/service_locator.dart';
import '../../../../core/models/layout.dart';
import '../../../../core/models/place.dart';
import '../../../../core/utils/bloc_load_state.dart';
import '../../../../core/widgets/layout_preview.dart';
import '../../../analytics/bloc/analytics_bloc.dart';
import '../../../analytics/bloc/analytics_event.dart';
import '../../../analytics/bloc/analytics_state.dart';
import '../../../analytics/widgets/heatmap_view.dart';
import '../../bloc/admin_bloc.dart';
import '../../bloc/admin_event.dart';
import '../../bloc/admin_state.dart';
import 'bookings_list_panel.dart';

class CoworkingDetailsPanel extends StatefulWidget {
  final int tabIndex;
  final Function(int) onTabChanged;

  const CoworkingDetailsPanel({
    super.key,
    required this.tabIndex,
    required this.onTabChanged,
  });

  /// 🔹 Static layout preview widget
  static Widget layoutPreview(Layout layout, List<Place> places) {
    return LayoutPreview(
      layout: layout,
      places: places,
      unavailablePlaceIds: places
          .where((p) => !p.isActive)
          .map((e) => e.id)
          .toSet(),
    );
  }

  @override
  State<CoworkingDetailsPanel> createState() => _CoworkingDetailsPanelState();
}

class _CoworkingDetailsPanelState extends State<CoworkingDetailsPanel> {
  final ScrollController _tabsScrollController = ScrollController();

  @override
  void dispose() {
    _tabsScrollController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<AdminBloc, AdminState>(
      builder: (context, state) {
        if (state.selectedCoworking.status == LoadStatus.loading) {
          return const Center(child: CircularProgressIndicator());
        }
        if (state.selectedCoworking.status == LoadStatus.error) {
          return Center(
            child: Text(
              'Failed to load coworking. ${state.selectedCoworking.error!}',
            ),
          );
        }
        if (state.selectedCoworking.status == LoadStatus.initial) {
          return const Center(child: Text('No coworking selected'));
        }

        final coworking = state.selectedCoworking.data!;

        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            /// 🔹 Header
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  coworking.name,
                  style: const TextStyle(
                    fontSize: 22,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                ElevatedButton(
                  onPressed: () {
                    _confirmSetActive(
                      context,
                      coworking.id,
                      coworking.isActive,
                    );
                  },
                  style: ButtonStyle(
                    backgroundColor: MaterialStateProperty.all(
                      coworking.isActive ? Colors.red : Colors.green,
                    ),
                  ),
                  child: Text(
                    state.selectedCoworking.data!.isActive
                        ? 'Deactivate'
                        : 'Activate',
                    style: const TextStyle(color: Colors.white),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),

            /// 🔹 Tabs with Scrollbar and Visual Indicators
            SizedBox(
              height: 54, // Высота под табы и скроллбар
              child: Scrollbar(
                controller: _tabsScrollController,
                thumbVisibility: true,
                thickness: 4,
                radius: const Radius.circular(8),
                child: Padding(
                  padding: const EdgeInsets.only(
                    bottom: 8,
                  ), // Отступ для самого скроллбара
                  child: SingleChildScrollView(
                    controller: _tabsScrollController,
                    scrollDirection: Axis.horizontal,
                    physics: const BouncingScrollPhysics(),
                    child: Row(
                      children: [
                        _tabButton('Info', 0),
                        _tabButton('Places', 1),
                        _tabButton('Layout', 2),
                        _tabButton('Active bookings', 3),
                        _tabButton('Analytics', 4),
                      ],
                    ),
                  ),
                ),
              ),
            ),
            const Divider(height: 1),

            /// 🔹 Content
            Expanded(child: _buildTabContent(state)),
          ],
        );
      },
    );
  }

  Widget _tabButton(String text, int index) {
    final isSelected = widget.tabIndex == index;

    return Padding(
      padding: const EdgeInsets.only(right: 8),
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 250),
        decoration: BoxDecoration(
          border: Border(
            bottom: BorderSide(
              color: isSelected ? Colors.blue : Colors.transparent,
              width: 3,
            ),
          ),
        ),
        child: TextButton(
          onPressed: () {
            widget.onTabChanged(index);
            if (index == 4) {
              final coworkingId = context
                  .read<AdminBloc>()
                  .state
                  .selectedCoworking
                  .data
                  ?.id;
              if (coworkingId != null) {
                BlocProvider(
                  create: (_) =>
                      sl<AnalyticsBloc>()
                        ..add(LoadCoworkingHeatmapEvent(coworkingId)),
                );
              }
            }
          },
          style: TextButton.styleFrom(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(4),
            ),
          ),
          child: Text(
            text,
            style: TextStyle(
              color: isSelected ? Colors.blue : Colors.black87,
              fontWeight: isSelected ? FontWeight.bold : FontWeight.w500,
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildTabContent(AdminState state) {
    switch (widget.tabIndex) {
      case 0:
        return _infoTab(state);
      case 1:
        return _placesTab(state);
      case 2:
        return _layoutTab(state);
      case 3:
        return _bookingsTab(state);
      case 4:
        return _analyticsTab(state);
      default:
        return const SizedBox();
    }
  }

  /// =====================
  /// ANALYTICS TAB
  /// =====================
  Widget _analyticsTab(AdminState state) {
    return BlocProvider(
      create: (_) =>
          sl<AnalyticsBloc>()
            ..add(LoadCoworkingHeatmapEvent(state.selectedCoworking.data!.id)),
      child: BlocBuilder<AnalyticsBloc, AnalyticsState>(
        builder: (context, analyticsState) {
          if (analyticsState.coworkingHeatmapState.status ==
              LoadStatus.loading) {
            return const Center(child: CircularProgressIndicator());
          }

          if (analyticsState.coworkingHeatmapState.status == LoadStatus.error) {
            return Center(
              child: Text(
                'Error loading analytics: ${analyticsState.coworkingHeatmapState.error}',
              ),
            );
          }

          final heatmap = analyticsState.coworkingHeatmapState.data?.heatmap;

          return SingleChildScrollView(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.center,
              children: [
                const Text(
                  'Карта загрузки мест коворкинга',
                  style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
                ),
                const SizedBox(height: 16),
                if (heatmap != null && heatmap.isNotEmpty)
                  HeatmapView(data: heatmap)
                else
                  const Center(child: Text('No data available for heatmap')),
              ],
            ),
          );
        },
      ),
    );
  }

  /// =====================
  /// INFO TAB
  /// =====================
  Widget _infoTab(AdminState state) {
    final coworking = state.selectedCoworking.data!;
    final dateFormat = DateFormat('dd MMM yy, HH:mm', 'RU');
    return Padding(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _infoRow(Icons.business, 'Name', coworking.name),
          _infoRow(Icons.location_on, 'Address', coworking.address),
          _infoRow(
            Icons.calendar_today,
            'Created',
            dateFormat.format(coworking.createdAt),
          ),
          _infoRow(
            Icons.edit_calendar,
            'Updated',
            dateFormat.format(coworking.updatedAt),
          ),
        ],
      ),
    );
  }

  // Вспомогательный виджет для красоты строк
  Widget _infoRow(IconData icon, String label, String value) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: Row(
        children: [
          Icon(icon, size: 20, color: Colors.blueGrey),
          const SizedBox(width: 12),
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                label,
                style: const TextStyle(fontSize: 12, color: Colors.grey),
              ),
              Text(
                value,
                style: const TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.w500,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  /// =====================
  /// PLACES TAB
  /// =====================
  Widget _placesTab(AdminState state) {
    if (state.places.status == LoadStatus.loading) {
      return const Center(child: CircularProgressIndicator());
    }
    if (state.places.status == LoadStatus.error) {
      return const Center(child: Text('Failed to load places'));
    }
    if (state.places.status == LoadStatus.success &&
        state.places.data == null) {
      return const Center(child: Text('No places'));
    }
    final places = state.places.data!;

    return Column(
      children: [
        /// Header
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            const Text(
              'Places',
              style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
            ),
            ElevatedButton(
              onPressed: () => _showAddPlacesDialog(
                context,
                state.selectedCoworking.data!.id,
                state.places.data!.length,
              ),
              child: const Text('Add Places'),
            ),
          ],
        ),
        const SizedBox(height: 16),

        /// Table
        Expanded(
          child: SingleChildScrollView(
            child: DataTable(
              columns: const [
                DataColumn(label: Text('Label')),
                DataColumn(label: Text('Type')),
                DataColumn(label: Text('Active')),
              ],
              rows: places.map((p) {
                return DataRow(
                  cells: [
                    DataCell(Text(p.label)),
                    DataCell(Text(p.placeType.toTitleCase(context))),
                    DataCell(
                      Switch(
                        value: p.isActive,
                        onChanged: (value) {
                          context.read<AdminBloc>().add(
                            SetPlaceActiveEvent(p.id, value),
                          );
                        },
                      ),
                    ),
                  ],
                );
              }).toList(),
            ),
          ),
        ),
      ],
    );
  }

  /// =====================
  /// LAYOUT TAB
  /// =====================
  Widget _layoutTab(AdminState state) {
    final coworking = state.selectedCoworking.data!;

    final layoutState = state.layout;
    final versionsState = state.layoutVersions;

    final layout = layoutState.data;
    final versions = versionsState.data ?? [];

    final isLayoutLoading = layoutState.status == LoadStatus.loading;
    final isLayoutError = layoutState.status == LoadStatus.error;

    versions.sort((a, b) => b.version.compareTo(a.version));

    final nextVersion = versions.isEmpty ? 1 : versions.first.version + 1;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        /// 🔹 ACTIONS
        Row(
          children: [
            ElevatedButton(
              onPressed: () =>
                  _loadLayoutFromFile(context, coworking.id, nextVersion),
              child: const Text('Upload JSON'),
            ),
            const SizedBox(width: 12),
            ElevatedButton(
              onPressed: layout == null
                  ? null
                  : () => _exportLayoutToFile(layout),
              child: const Text('Download JSON'),
            ),
          ],
        ),

        const SizedBox(height: 16),

        /// 🔹 CURRENT VERSION / STATUS
        if (isLayoutLoading)
          const Center(child: CircularProgressIndicator())
        else if (isLayoutError)
          const Text('Failed to load layout')
        else if (layout != null)
          Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: Colors.blue.withOpacity(0.05),
              borderRadius: BorderRadius.circular(8),
            ),
            child: Row(
              children: [
                const Icon(Icons.layers),
                const SizedBox(width: 8),
                Text('Current version: v${layout.version}'),
              ],
            ),
          ),

        const SizedBox(height: 16),

        const Text(
          'Versions',
          style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
        ),

        const SizedBox(height: 8),

        /// 🔹 VERSIONS LIST
        Expanded(
          child: Builder(
            builder: (_) {
              if (state.layoutVersions.status == LoadStatus.loading) {
                return const Center(child: CircularProgressIndicator());
              }

              if (state.layoutVersions.status == LoadStatus.error) {
                return Center(
                  child: Text(
                    'Failed to load versions: ${state.layoutVersions.error}',
                  ),
                );
              }

              if (versions.isEmpty) {
                return const Center(child: Text('No versions'));
              }

              return ListView.builder(
                itemCount: versions.length,
                itemBuilder: (context, index) {
                  final v = versions[index];
                  final isCurrent = layout?.version == v.version;
                  final dateFormat = DateFormat('dd.MM.yy, HH:mm', 'RU');

                  return Card(
                    child: ListTile(
                      leading: Icon(
                        isCurrent ? Icons.check_circle : Icons.history,
                        color: isCurrent ? Colors.green : Colors.grey,
                      ),
                      title: Text('Version ${v.version}'),
                      subtitle: Row(
                        children: [
                          Icon(
                            Icons.calendar_today,
                            size: 16,
                            color: Colors.blueGrey,
                          ),
                          const SizedBox(width: 4),
                          Text(
                            'Created',
                            style: const TextStyle(
                              fontSize: 12,
                              color: Colors.grey,
                            ),
                          ),
                          const SizedBox(width: 4),
                          Text(
                            dateFormat.format(v.createdAt),
                            style: const TextStyle(
                              fontSize: 12,
                              fontWeight: FontWeight.w500,
                            ),
                          ),
                        ],
                      ),
                      trailing: isCurrent
                          ? const Text('Current')
                          : Row(
                              mainAxisSize: MainAxisSize.min,
                              children: [
                                TextButton(
                                  onPressed: () {
                                    context.read<AdminBloc>().add(
                                      SetActiveLayout(coworking.id, v.version),
                                    );
                                  },
                                  child: const Text('Activate'),
                                ),
                                IconButton(
                                  icon: const Icon(Icons.delete),
                                  onPressed: () {
                                    context.read<AdminBloc>().add(
                                      DeleteLayoutEvent(
                                        coworking.id,
                                        v.version,
                                      ),
                                    );
                                  },
                                ),
                              ],
                            ),
                    ),
                  );
                },
              );
            },
          ),
        ),
      ],
    );
  }

  /// =====================
  /// HELPERS
  /// =====================
  void _confirmSetActive(
    BuildContext context,
    String coworkingId,
    bool active,
  ) {
    showDialog(
      context: context,
      builder: (_) => AlertDialog(
        title: Text(active ? 'Confirm Deactivation' : 'Confirm Activation'),
        content: Text(
          'Do you really want to ${active ? 'deactivate' : 'activate'} this coworking?',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              context.read<AdminBloc>().add(
                active
                    ? DeactivateCoworkingEvent(coworkingId)
                    : ActivateCoworkingEvent(coworkingId),
              );
              Navigator.pop(context);
            },
            child: Text(active ? 'Deactivate' : 'Activate'),
          ),
        ],
      ),
    );
  }

  void _showAddPlacesDialog(
    BuildContext context,
    String coworkingId,
    int placesNumber,
  ) {
    final countController = TextEditingController();
    String selectedType = 'open_desk';
    final adminBloc = context.read<AdminBloc>();

    showDialog(
      context: context,
      builder: (dialogContext) {
        return StatefulBuilder(
          builder: (context, setState) {
            return AlertDialog(
              title: const Text('Add Places'),
              content: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  TextField(
                    controller: countController,
                    keyboardType: TextInputType.number,
                    decoration: const InputDecoration(labelText: 'Count'),
                  ),
                  const SizedBox(height: 12),
                  DropdownButton<String>(
                    value: selectedType,
                    items: const [
                      DropdownMenuItem(
                        value: 'open_desk',
                        child: Text('Open desk'),
                      ),
                      DropdownMenuItem(value: 'room', child: Text('Room')),
                    ],
                    onChanged: (value) {
                      setState(() => selectedType = value!);
                    },
                  ),
                ],
              ),
              actions: [
                TextButton(
                  onPressed: () => Navigator.pop(dialogContext),
                  child: const Text('Cancel'),
                ),
                ElevatedButton(
                  onPressed: () {
                    final count = int.tryParse(countController.text) ?? 0;
                    if (count <= 0) return;
                    final places = List.generate(
                      count,
                      (index) => Place(
                        label: 'p${index + placesNumber}',
                        placeType: PlaceType.fromJson(selectedType),
                      ),
                    );
                    adminBloc.add(AddPlacesEvent(coworkingId, places));
                    Navigator.pop(dialogContext);
                  },
                  child: const Text('Add'),
                ),
              ],
            );
          },
        );
      },
    );
  }

  /// =====================
  /// LAYOUT JSON LOAD / EXPORT
  /// =====================
  Future<void> _loadLayoutFromFile(
    BuildContext context,
    String coworkingId,
    int newVersion,
  ) async {
    final result = await FilePicker.platform.pickFiles(
      type: FileType.custom,
      allowedExtensions: ['json'],
    );
    if (result == null) return;
    try {
      final fileBytes = result.files.single.bytes;
      final jsonStr = utf8.decode(fileBytes!);
      final jsonMap = jsonDecode(jsonStr);
      final layoutSchema = LayoutSchema.fromJson(jsonMap);

      context.read<AdminBloc>().add(
        CreateLayoutEvent(coworkingId, layoutSchema, newVersion),
      );
    } catch (e) {
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(const SnackBar(content: Text('Invalid JSON file')));
    }
  }

  void _exportLayoutToFile(Layout? layout) {
    if (layout == null) return;
    final jsonStr = jsonEncode(layout.toJson());
    // TODO: implement actual file save (Desktop or Web)
    // For now, just copy to clipboard
    Clipboard.setData(ClipboardData(text: jsonStr));
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Layout JSON copied to clipboard')),
    );
  }

  /// =====================
  /// ACTIVE BOOKINGS TAB
  /// ===================== //
  Widget _bookingsTab(AdminState state) {
    return BookingListPanel();
  }
}
