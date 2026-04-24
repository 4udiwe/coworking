import 'package:coworking_app/features/notification/bloc/notification_event.dart';
import 'package:coworking_app/features/notification/presentation/screens/notification_screen.dart';
import 'package:coworking_app/features/user/presentation/profile_screen.dart';
import 'package:coworking_app/core/l10n/localization_helper.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../features/admin/presentation/admin_screen.dart';
import '../../features/coworking/presentation/screens/coworking_list_screen.dart';

import '../../features/auth/bloc/auth_bloc.dart';
import '../../features/auth/bloc/auth_state.dart';
import '../../features/bookings/presentation/screens/bookings_screen.dart';

// 🔔 NEW
import '../../features/notification/bloc/notification_bloc.dart';
import '../../features/notification/bloc/notification_state.dart';
import '../../features/notification/presentation/widgets/app_toast.dart';
import 'features/notification/service/toast_service.dart';

class MainScreen extends StatefulWidget {
  const MainScreen({super.key});

  @override
  State<MainScreen> createState() => _MainScreenState();
}

class _MainScreenState extends State<MainScreen> {
  int _index = 0;

  /// ❗ анти-дубли
  final Set<String> _shownNotifications = {};

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<AuthBloc, AuthState>(
      builder: (context, authState) {
        if (authState.status != AuthStatus.authenticated) {
          return const Scaffold(body: Center(child: Text('Not authenticated')));
        }

        final pages = <Widget>[
          const CoworkingPage(),
          const BookingsPage(),
          const NotificationScreen(),
          const ProfileScreen(),
          if (authState.isAdmin) const AdminScreen(),
        ];

        final items = <BottomNavigationBarItem>[
          BottomNavigationBarItem(
            icon: const Icon(Icons.work),
            label: context.l10n.navCoworkings,
          ),
          BottomNavigationBarItem(
            icon: const Icon(Icons.calendar_today),
            label: context.l10n.navBookings,
          ),
          BottomNavigationBarItem(
            icon: BlocBuilder<NotificationBloc, NotificationState>(
              builder: (context, state) {
                final count = state.unreadCount;

                return Stack(
                  clipBehavior: Clip.none,
                  children: [
                    const Icon(Icons.notifications),
                    if (count > 0)
                      Positioned(
                        right: -6,
                        top: -4,
                        child: Container(
                          padding: const EdgeInsets.all(4),
                          decoration: BoxDecoration(
                            color: Colors.red,
                            borderRadius: BorderRadius.circular(10),
                          ),
                          constraints: const BoxConstraints(
                            minWidth: 16,
                            minHeight: 16,
                          ),
                          child: Text(
                            count > 99 ? '99+' : '$count',
                            style: const TextStyle(
                              color: Colors.white,
                              fontSize: 10,
                            ),
                            textAlign: TextAlign.center,
                          ),
                        ),
                      ),
                  ],
                );
              },
            ),
            label: context.l10n.navNotifications,
          ),
          BottomNavigationBarItem(
            icon: const Icon(Icons.person),
            label: context.l10n.navProfile,
          ),
          if (authState.isAdmin && kIsWeb)
            BottomNavigationBarItem(
              icon: const Icon(Icons.admin_panel_settings),
              label: context.l10n.navAdmin,
            ),
        ];

        if (_index >= pages.length) {
          _index = 0;
        }

        return BlocListener<NotificationBloc, NotificationState>(
          listenWhen: (prev, curr) => prev.messageId != curr.messageId,
          listener: (context, state) {
            final notifications = state.notifications.data ?? [];
            print("notifications = ${notifications.length}");
            for (final n in notifications) {
              print("notification = ${n.title}");
            }

            if (notifications.isEmpty) return;

            final latest = notifications.first;

            if (_shownNotifications.contains(latest.id)) return;
            _shownNotifications.add(latest.id);

            if (!_shouldShowToast(latest.type)) return;

            ToastService.show(
              context,
              builder: (remove, animation) {
                return AppToast(
                  title: latest.title,
                  body: latest.displayBody,
                  onMarkRead: () {
                    context.read<NotificationBloc>().add(
                      MarkNotificationRead(latest.id),
                    );
                    remove();
                  },

                  /// 🔥 NAVIGATION
                  onTap: () {
                    if (latest.actionUrl != null) {
                      Navigator.of(context).pushNamed(latest.actionUrl!);
                    }
                  },
                );
              },
            );
          },
          child: Scaffold(
            body: pages[_index],
            bottomNavigationBar: BottomNavigationBar(
              type: BottomNavigationBarType.shifting,
              selectedItemColor: Colors.blue,
              unselectedItemColor: Colors.grey,
              currentIndex: _index,
              onTap: (value) {
                setState(() => _index = value);
              },
              items: items,
            ),
          ),
        );
      },
    );
  }

  bool _shouldShowToast(String type) {
    return type == 'booking_created' ||
        type == 'booking_cancelled' ||
        type == 'booking_reminder';
  }
}
