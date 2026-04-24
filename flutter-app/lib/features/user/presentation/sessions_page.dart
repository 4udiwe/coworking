import 'package:coworking_app/core/l10n/localization_helper.dart';
import 'package:coworking_app/features/user/presentation/widgets/session_tile.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../core/di/service_locator.dart';
import '../../../core/utils/bloc_load_state.dart';
import '../bloc/user_bloc.dart';
import '../bloc/user_event.dart';
import '../bloc/user_state.dart';

class SessionsPage extends StatelessWidget {
  const SessionsPage({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) => sl<UserBloc>()..add(LoadAllSessions()),
      child: const SessionsScreen(),
    );
  }
}

class SessionsScreen extends StatefulWidget {
  final bool showAppBar;
  const SessionsScreen({super.key, this.showAppBar = true});

  @override
  State<SessionsScreen> createState() => _SessionsScreenState();
}

class _SessionsScreenState extends State<SessionsScreen>
    with SingleTickerProviderStateMixin {
  @override
  Widget build(BuildContext context) {
    final content = BlocListener<UserBloc, UserState>(
      listenWhen: (prev, curr) => prev.messageId != curr.messageId,
      listener: (context, state) {
        if (state.actionMessage != null) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(state.actionMessage!),
              backgroundColor: state.isError ? Colors.red : Colors.green,
              duration: const Duration(seconds: 2),
            ),
          );
        }
      },
      child: const SessionList(),
    );

    if (!widget.showAppBar) {
      return content;
    }

    return Scaffold(
      appBar: AppBar(title: Text(context.l10n.sessions)),
      body: content,
    );
  }
}

class SessionList extends StatelessWidget {
  const SessionList({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<UserBloc, UserState>(
      builder: (context, state) {
        final items = state.allSessions.data ?? [];
        final loading = state.allSessions.status == LoadStatus.loading;

        if (loading && items.isEmpty) {
          return const Center(child: CircularProgressIndicator());
        }

        if (items.isEmpty) {
          return const Center(child: Text('Нет активных сессий'));
        }

        return RefreshIndicator(
          onRefresh: () async {
            context.read<UserBloc>().add(LoadAllSessions());
          },
          child: ListView.builder(
            padding: const EdgeInsets.only(top: 8, bottom: 40),
            itemCount: items.length,
            itemBuilder: (_, i) {
              return Padding(
                padding: const EdgeInsets.symmetric(
                  horizontal: 12,
                  vertical: 4,
                ),
                child: SessionTile(session: items[i]),
              );
            },
          ),
        );
      },
    );
  }
}
