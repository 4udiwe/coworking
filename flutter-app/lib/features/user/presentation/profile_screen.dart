import 'package:coworking_app/features/user/presentation/sessions_page.dart';
import 'package:coworking_app/core/l10n/localization_helper.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../core/di/service_locator.dart';
import '../../../core/utils/bloc_load_state.dart';
import '../../auth/bloc/auth_bloc.dart';
import '../../auth/bloc/auth_event.dart';

import '../bloc/user_bloc.dart';
import '../bloc/user_event.dart';
import '../bloc/user_state.dart';

class ProfileScreen extends StatelessWidget {
  const ProfileScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) => sl<UserBloc>()..add(RefreshUserData()),
      child: BlocListener<UserBloc, UserState>(
        listenWhen: (prev, curr) =>
            prev.messageId != curr.messageId, // любое новое сообщение
        listener: (context, state) {
          if (state.actionMessage != null) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                duration: const Duration(seconds: 2),
                content: Text(state.actionMessage!),
                backgroundColor: state.isError ? Colors.red : Colors.green,
              ),
            );
          }
        },
        child: LayoutBuilder(
          builder: (context, constraints) {
            final isWide = constraints.maxWidth > 800;

            if (!isWide) {
              /// 📱 Mobile (как сейчас)
              return _ProfileView(isWide: isWide);
            }

            /// 💻 Tablet / Desktop
            return _ProfileWithSessions(isWide: isWide);
          },
        ),
      ),
    );
  }
}

class _ProfileWithSessions extends StatelessWidget {
  final bool isWide;
  const _ProfileWithSessions({required this.isWide});

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        /// 🔹 Левая часть — профиль
        Expanded(flex: 3, child: _ProfileViewContent(isWide: isWide)),

        /// 🔹 Правая часть — сессии
        Expanded(
          flex: 2,
          child: Center(
            child: Container(
              constraints: BoxConstraints(
                maxHeight: MediaQuery.of(context).size.height * 0.8,
                maxWidth: 500,
              ),
              margin: const EdgeInsets.symmetric(horizontal: 40),
              child: Card(
                elevation: 8,
                shadowColor: Colors.black26,
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(24),
                ),
                clipBehavior:
                    Clip.antiAlias, // Чтобы контент не вылезал за скругления
                child: const _EmbeddedSessionsView(),
              ),
            ),
          ),
        ),
      ],
    );
  }
}

class _ProfileView extends StatelessWidget {
  final bool isWide;
  const _ProfileView({required this.isWide});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(context.l10n.profileScreenTitle),
        centerTitle: true,
      ),
      body: _ProfileViewContent(isWide: isWide),
    );
  }
}

class _ProfileViewContent extends StatelessWidget {
  final bool isWide;
  const _ProfileViewContent({required this.isWide});

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<UserBloc, UserState>(
      builder: (context, state) {
        /// 🔹 Loading профиля
        if (state.profile.status == LoadStatus.loading) {
          return const _LoadingView();
        }

        /// 🔹 Ошибка профиля
        if (state.profile.status == LoadStatus.error) {
          return _ErrorView(message: state.profile.error!);
        }

        /// 🔹 Успех
        if (state.profile.status == LoadStatus.success) {
          final user = state.profile.data!;

          return RefreshIndicator(
            onRefresh: () async {
              context.read<UserBloc>().add(RefreshUserData());
            },
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(16),
              physics: const AlwaysScrollableScrollPhysics(),
              child: Column(
                children: [
                  const SizedBox(height: 16),

                  _AvatarSection(name: '${user.firstName} ${user.lastName}'),

                  const SizedBox(height: 24),

                  _InfoCard(
                    children: [
                      _InfoTile(
                        icon: Icons.person,
                        title: 'Имя',
                        value: user.firstName,
                      ),
                      _InfoTile(
                        icon: Icons.person_outline,
                        title: 'Фамилия',
                        value: user.lastName,
                      ),
                      _InfoTile(
                        icon: Icons.email,
                        title: 'Email',
                        value: user.email,
                      ),
                      _InfoTile(
                        icon: Icons.badge,
                        title: 'Роли',
                        value: state.profile.data!.roles
                            .map((e) => e.name)
                            .join(", "),
                      ),
                    ],
                  ),

                  const SizedBox(height: 24),

                  if (!isWide) ...[
                    /// 🔥 Сессии
                    SizedBox(
                      width: double.infinity,
                      child: ElevatedButton.icon(
                        style: ElevatedButton.styleFrom(
                          backgroundColor: Colors.blue,
                          foregroundColor: Colors.white,
                          padding: const EdgeInsets.symmetric(vertical: 14),
                        ),
                        onPressed: () {
                          Navigator.pushNamed(context, '/sessions');
                        },
                        icon: const Icon(Icons.devices),
                        label: const Text('Управление сессиями'),
                      ),
                    ),

                    const SizedBox(height: 24),
                  ],

                  _LogoutButton(),
                ],
              ),
            ),
          );
        }

        return const SizedBox();
      },
    );
  }
}

class _EmbeddedSessionsView extends StatelessWidget {
  const _EmbeddedSessionsView();

  @override
  Widget build(BuildContext context) {
    return BlocProvider.value(
      value: context.read<UserBloc>()..add(LoadAllSessions()),
      child: const SessionsScreen(),
    );
  }
}

class _LoadingView extends StatelessWidget {
  const _LoadingView();
  @override
  Widget build(BuildContext context) {
    return const Center(child: CircularProgressIndicator());
  }
}

class _ErrorView extends StatelessWidget {
  final String message;
  const _ErrorView({required this.message});
  @override
  Widget build(BuildContext context) {
    return Center(
      child: Text(
        message,
        style: const TextStyle(color: Colors.red),
        textAlign: TextAlign.center,
      ),
    );
  }
}

class _InfoTile extends StatelessWidget {
  final IconData icon;
  final String title;
  final String value;
  const _InfoTile({
    required this.icon,
    required this.title,
    required this.value,
  });
  @override
  Widget build(BuildContext context) {
    return ListTile(
      leading: Icon(icon, color: Colors.blue),
      title: Text(title),
      subtitle: Text(value),
    );
  }
}

class _LogoutButton extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: double.infinity,
      child: ElevatedButton.icon(
        style: ElevatedButton.styleFrom(
          backgroundColor: Colors.red,
          foregroundColor: Colors.white,
          padding: const EdgeInsets.symmetric(vertical: 14),
        ),
        onPressed: () => _confirmLogout(context),
        icon: const Icon(Icons.logout),
        label: const Text('Выйти'),
      ),
    );
  }

  void _confirmLogout(BuildContext context) {
    showDialog(
      context: context,
      builder: (_) => AlertDialog(
        title: const Text('Выход'),
        content: const Text('Вы уверены, что хотите выйти?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Отмена'),
          ),
          TextButton(
            onPressed: () {
              Navigator.pop(context);
              context.read<AuthBloc>().add(AuthLogout());
            },
            child: const Text('Выйти'),
          ),
        ],
      ),
    );
  }
}

class _AvatarSection extends StatelessWidget {
  final String name;
  const _AvatarSection({required this.name});
  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        CircleAvatar(
          radius: 48,
          backgroundColor: Colors.grey.shade200,
          child: const Icon(Icons.person, size: 48, color: Colors.grey),
        ),
        const SizedBox(height: 12),
        Text(name, style: Theme.of(context).textTheme.titleLarge),
      ],
    );
  }
}

class _InfoCard extends StatelessWidget {
  final List<Widget> children;
  const _InfoCard({required this.children});
  @override
  Widget build(BuildContext context) {
    return Card(
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      elevation: 2,
      child: Column(children: children),
    );
  }
}
