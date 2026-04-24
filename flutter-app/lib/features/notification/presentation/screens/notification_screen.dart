import 'package:coworking_app/core/l10n/localization_helper.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../../core/utils/bloc_load_state.dart';
import '../../bloc/notification_bloc.dart';
import '../../bloc/notification_event.dart';
import '../../bloc/notification_state.dart';
import '../widgets/notification_tile.dart';

class NotificationScreen extends StatefulWidget {
  const NotificationScreen({super.key});

  @override
  State<NotificationScreen> createState() => _NotificationScreenState();
}

class _NotificationScreenState extends State<NotificationScreen> {
  final ScrollController _scrollController = ScrollController();
  static const int _pageSize = 20;

  @override
  void initState() {
    super.initState();

    // initial fetch
    context.read<NotificationBloc>().add(
      FetchNotifications(limit: _pageSize, offset: 0),
    );

    // infinite scroll
    _scrollController.addListener(() {
      final bloc = context.read<NotificationBloc>();
      final state = bloc.state.notifications;

      if (_scrollController.position.pixels >=
              _scrollController.position.maxScrollExtent - 100 &&
          state.status != LoadStatus.loading &&
          state.data != null &&
          state.data!.length >= _pageSize) {
        bloc.add(
          FetchNotifications(limit: _pageSize, offset: state.data!.length),
        );
      }
    });
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }

  Future<void> _onRefresh() async {
    context.read<NotificationBloc>().add(
      FetchNotifications(limit: _pageSize, offset: 0, refresh: true),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(context.l10n.notifications),
        actions: [
          IconButton(
            icon: const Icon(Icons.done_all),
            tooltip: context.l10n.markAllAsRead,
            onPressed: () {
              context.read<NotificationBloc>().add(MarkAllNotificationsRead());
            },
          ),
        ],
      ),
      body: BlocBuilder<NotificationBloc, NotificationState>(
        builder: (context, state) {
          final notifications = state.notifications.data ?? [];

          if (state.notifications.status == LoadStatus.error) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Icon(Icons.error_outline, size: 48, color: Colors.red),
                  const SizedBox(height: 16),
                  Text(
                    state.notifications.error ?? 'Error loading notifications',
                  ),
                  ElevatedButton(
                    onPressed: _onRefresh,
                    child: const Text('Retry'),
                  ),
                ],
              ),
            );
          }

          if (state.notifications.status == LoadStatus.loading &&
              notifications.isEmpty) {
            return const Center(child: CircularProgressIndicator());
          }

          if (notifications.isEmpty) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(
                    Icons.notifications_none,
                    size: 64,
                    color: Theme.of(context).disabledColor,
                  ),
                  const SizedBox(height: 16),
                  const Text('No notifications yet'),
                ],
              ),
            );
          }

          return RefreshIndicator(
            onRefresh: _onRefresh,
            child: ListView.builder(
              controller: _scrollController,
              itemCount: notifications.length,
              itemBuilder: (context, index) {
                final item = notifications[index];
                return NotificationTile(
                  notification: item,
                  onTap: () {
                    if (!item.isRead) {
                      context.read<NotificationBloc>().add(
                        MarkNotificationRead(item.id),
                      );
                    }
                    if (item.actionUrl != null) {
                      // context.go(item.actionUrl!);
                    }
                  },
                );
              },
            ),
          );
        },
      ),
    );
  }
}
