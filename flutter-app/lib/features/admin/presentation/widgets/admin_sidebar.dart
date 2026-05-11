import 'package:coworking_app/features/admin/presentation/widgets/create_coworking/create_coworking_dialog.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:coworking_app/core/l10n/localization_helper.dart';

import '../../../../core/utils/bloc_load_state.dart';
import '../../bloc/admin_bloc.dart';
import '../../bloc/admin_event.dart';
import '../../bloc/admin_state.dart';
import 'admin_actions_panel.dart';

class AdminSidebar extends StatelessWidget {
  const AdminSidebar({super.key});

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: 280,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: const EdgeInsets.all(14),
            child: Text(
              context.l10n.adminPanelTitle,
              style: const TextStyle(fontSize: 20),
            ),
          ),
          const Divider(),

          // ── Add coworking button ──
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
            child: SizedBox(
              width: double.infinity,
              child: OutlinedButton.icon(
                icon: const Icon(Icons.add, size: 18),
                label: const Text('Add Coworking'),
                style: OutlinedButton.styleFrom(
                  alignment: Alignment.centerLeft,
                  padding: const EdgeInsets.symmetric(
                    horizontal: 12,
                    vertical: 10,
                  ),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(8),
                  ),
                ),
                onPressed: () => CreateCoworkingDialog.show(context),
              ),
            ),
          ),

          /// 🔹 Upper part: coworkings list
          Expanded(
            child: BlocBuilder<AdminBloc, AdminState>(
              builder: (context, state) {
                if (state.coworkings.status == LoadStatus.loading) {
                  return const Center(child: CircularProgressIndicator());
                }
                if (state.coworkings.status == LoadStatus.error) {
                  return Center(
                    child: Text(
                      '${context.l10n.adminFailedLoadCoworkings}. ${state.coworkings.error!}',
                    ),
                  );
                }
                if (state.coworkings.status == LoadStatus.success &&
                    state.coworkings.data == null) {
                  return Center(child: Text(context.l10n.adminNoCoworkings));
                }

                final coworkings = state.coworkings.data ?? [];
                if (coworkings.isEmpty) {
                  return Center(child: Text(context.l10n.adminNoCoworkings));
                }

                return ListView.builder(
                  itemCount: coworkings.length,
                  itemBuilder: (context, index) {
                    final c = coworkings[index];
                    return _CoworkingListItem(
                      coworkingId: c.id,
                      name: c.name,
                      address: c.address,
                      isActive: c.isActive,
                    );
                  },
                );
              },
            ),
          ),

          const Divider(),

          /// 🔹 Lower part: admin actions
          Padding(
            padding: const EdgeInsets.all(12.0),
            child: const AdminActionsPanel(),
          ),
        ],
      ),
    );
  }
}

/// =====================
/// COWORKING LIST ITEM
/// =====================
class _CoworkingListItem extends StatefulWidget {
  final String coworkingId;
  final String name;
  final String address;
  final bool isActive;

  const _CoworkingListItem({
    required this.coworkingId,
    required this.name,
    required this.address,
    required this.isActive,
  });

  @override
  State<_CoworkingListItem> createState() => _CoworkingListItemState();
}

class _CoworkingListItemState extends State<_CoworkingListItem> {
  bool isHovered = false;

  @override
  Widget build(BuildContext context) {
    // Получаем состояние из блока, чтобы узнать, выбран ли этот коворкинг
    return BlocBuilder<AdminBloc, AdminState>(
      buildWhen: (previous, current) =>
          previous.selectedCoworking.data?.id !=
          current.selectedCoworking.data?.id,
      builder: (context, state) {
        final isSelected =
            state.selectedCoworking.data?.id == widget.coworkingId;

        return MouseRegion(
          onEnter: (_) => setState(() => isHovered = true),
          onExit: (_) => setState(() => isHovered = false),
          child: Container(
            // Плавное изменение фона при наведении или выборе
            decoration: BoxDecoration(
              color: isSelected
                  ? Colors.blue.withOpacity(0.1) // Цвет при выборе
                  : (isHovered
                        ? Colors.black.withOpacity(0.04)
                        : Colors.transparent),
              borderRadius: BorderRadius.circular(8),
            ),
            margin: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
            child: ListTile(
              // Скругляем углы для эффекта нажатия (InkWell внутри ListTile)
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(8),
              ),
              title: Text(
                widget.name,
                style: TextStyle(
                  fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
                  color: isSelected ? Colors.blue[700] : null,
                ),
              ),
              subtitle: Text(
                widget.isActive ? 'Active' : 'Inactive',
                style: TextStyle(
                  color: widget.isActive ? Colors.green : Colors.red,
                  fontSize: 12,
                ),
              ),
              trailing: AnimatedOpacity(
                duration: const Duration(milliseconds: 200),
                opacity: isHovered ? 1.0 : 0.0,
                child: IconButton(
                  icon: const Icon(Icons.edit, size: 20),
                  onPressed: () => _showEditDialog(context),
                  tooltip: 'Edit coworking',
                ),
              ),
              onTap: () {
                final bloc = context.read<AdminBloc>();
                bloc.add(SetAdminViewEvent(AdminView.coworkingDetails));
                bloc.add(FetchCoworkingByIdEvent(widget.coworkingId));
                bloc.add(FetchPlacesEvent(widget.coworkingId));
                bloc.add(FetchLayoutVersionsEvent(widget.coworkingId));
                bloc.add(FetchLatestLayoutEvent(widget.coworkingId));
                bloc.add(
                  FetchActiveBookingsEvent(coworkingId: widget.coworkingId),
                );
              },
            ),
          ),
        );
      },
    );
  }

  void _showEditDialog(BuildContext context) {
    final nameController = TextEditingController(text: widget.name);
    final addressController = TextEditingController(text: widget.address);

    showDialog(
      context: context,
      builder: (_) {
        return AlertDialog(
          title: const Text('Edit Coworking'),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              TextField(
                controller: nameController,
                decoration: const InputDecoration(labelText: 'Name'),
              ),
              TextField(
                controller: addressController,
                decoration: const InputDecoration(labelText: 'Address'),
              ),
            ],
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text('Cancel'),
            ),
            ElevatedButton(
              onPressed: () {
                context.read<AdminBloc>().add(
                  UpdateCoworkingEvent(
                    widget.coworkingId,
                    nameController.text,
                    addressController.text,
                  ),
                );
                Navigator.pop(context);
              },
              child: const Text('Save'),
            ),
          ],
        );
      },
    );
  }

  // void _confirmDeactivate(BuildContext context) {
  //   showDialog(
  //     context: context,
  //     builder: (_) => AlertDialog(
  //       title: const Text('Confirm Deactivation'),
  //       content: Text('Do you really want to deactivate "${widget.name}"?'),
  //       actions: [
  //         TextButton(
  //           onPressed: () => Navigator.pop(context),
  //           child: const Text('Cancel'),
  //         ),
  //         ElevatedButton(
  //           onPressed: () {
  //             context.read<AdminBloc>().add(
  //               DeactivateCoworkingEvent(widget.coworkingId),
  //             );
  //             Navigator.pop(context);
  //           },
  //           child: const Text('Deactivate'),
  //         ),
  //       ],
  //     ),
  //   );
  // }
}
