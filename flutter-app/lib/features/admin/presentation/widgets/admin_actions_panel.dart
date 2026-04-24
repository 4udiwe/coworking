import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../bloc/admin_bloc.dart';
import '../../bloc/admin_event.dart';
import '../../bloc/admin_state.dart';

class AdminActionsPanel extends StatefulWidget {
  const AdminActionsPanel({super.key});

  @override
  State<AdminActionsPanel> createState() => _AdminActionsPanelState();
}

class _AdminActionsPanelState extends State<AdminActionsPanel> {

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        ListTile(
          leading: const Icon(Icons.person),
          title: const Text('Users'),
          onTap: () {
            context.read<AdminBloc>().add(
              SetAdminViewEvent(AdminView.users),
            );

            /// сразу грузим пользователей
            context.read<AdminBloc>().add(
              FetchUsersEvent(),
            );
          },
        ),
      ],
    );
  }
}