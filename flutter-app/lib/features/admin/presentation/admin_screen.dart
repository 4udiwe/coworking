import 'package:coworking_app/features/admin/presentation/widgets/admin_sidebar.dart';
import 'package:coworking_app/features/admin/presentation/widgets/coworking_details_panel.dart';
import 'package:coworking_app/features/admin/presentation/widgets/users_panel.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../core/utils/bloc_load_state.dart';
import '../bloc/admin_bloc.dart';
import '../bloc/admin_event.dart';
import '../bloc/admin_state.dart';

import '../../../../core/di/service_locator.dart';

class AdminScreen extends StatelessWidget {
  const AdminScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) => sl<AdminBloc>()..add(FetchCoworkingsEvent()),
      child: BlocListener<AdminBloc, AdminState>(
        listenWhen: (prev, curr) => prev.actionMessage != curr.actionMessage,
        listener: (context, state) {
          if (state.actionMessage != null) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(state.actionMessage!),
                backgroundColor: state.isError ? Colors.red : Colors.green,
              ),
            );
          }
        },
        child: const Scaffold(body: AdminViewWrapper()),
      ),
    );
  }
}

/// =====================
/// ADMIN VIEW WRAPPER
/// =====================
class AdminViewWrapper extends StatelessWidget {
  const AdminViewWrapper({super.key});

  @override
  Widget build(BuildContext context) {
    return Row(
      children: const [
        /// 🔹 SIDEBAR
        Expanded(flex: 1, child: AdminSidebar()),

        /// 🔹 MAIN CONTENT
        Expanded(flex: 5, child: AdminMainContent()),
      ],
    );
  }
}

/// =====================
/// MAIN CONTENT PANEL
/// =====================
class AdminMainContent extends StatefulWidget {
  const AdminMainContent({super.key});

  @override
  State<AdminMainContent> createState() => _AdminMainContentState();
}

class _AdminMainContentState extends State<AdminMainContent> {
  int tabIndex = 0;

  void onTabChanged(int index) {
    setState(() {
      tabIndex = index;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(16),
      child: BlocBuilder<AdminBloc, AdminState>(
        builder: (context, state) {
          final isWide = tabIndex == 3 || tabIndex == 4;

          switch (state.currentView) {
            case AdminView.coworkingDetails:
              // Выносим контент панелей в переменные, чтобы не дублировать логику
              final leftChild = switch (state.layout.status) {
                LoadStatus.loading => const Center(
                  child: CircularProgressIndicator(),
                ),
                LoadStatus.error => Center(
                  child: Text('Failed: ${state.layout.error}'),
                ),
                LoadStatus.success => CoworkingDetailsPanel.layoutPreview(
                  state.layout.data!,
                  state.places.data ?? [],
                ),
                LoadStatus.initial => const Center(
                  child: Text('Select a coworking'),
                ),
              };

              final rightChild = CoworkingDetailsPanel(
                tabIndex: tabIndex,
                onTabChanged: onTabChanged,
              );

              return TweenAnimationBuilder<double>(
                duration: const Duration(
                  milliseconds: 500,
                ), // 500мс обычно достаточно для комфортной плавности
                curve: Curves.easeInOutCubic,
                // Анимируем коэффициент от 0.0 до 1.0
                tween: Tween<double>(end: isWide ? 1.0 : 0.0),
                builder: (context, t, _) {
                  // Интерполируем flex с большим множителем (например, 10000) для идеальной плавности
                  // t=0 (not wide): left=4000, right=2000
                  // t=1 (wide):     left=3000, right=5000
                  final leftFlex = (4000 + (3000 - 4000) * t).toInt();
                  final rightFlex = (2000 + (5000 - 2000) * t).toInt();

                  return Row(
                    children: [
                      Expanded(flex: leftFlex, child: leftChild),
                      const SizedBox(width: 24),
                      Expanded(flex: rightFlex, child: rightChild),
                    ],
                  );
                },
              );

            case AdminView.users:
              return const UsersPanel();
          }
        },
      ),
    );
  }
}
