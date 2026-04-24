import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../../core/models/user.dart';
import '../../../../core/utils/bloc_load_state.dart';
import '../../../auth/bloc/auth_bloc.dart';
import '../../bloc/admin_bloc.dart';
import '../../bloc/admin_event.dart';
import '../../bloc/admin_state.dart';

class UsersPanel extends StatelessWidget {
  const UsersPanel({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<AdminBloc, AdminState>(
      builder: (context, state) {
        final usersState = state.users;

        if (usersState.status == LoadStatus.loading) {
          return const Center(child: CircularProgressIndicator());
        }

        if (usersState.status == LoadStatus.error) {
          return Center(child: Text(usersState.error ?? 'Error'));
        }

        final users = usersState.data?.items ?? [];

        // Get current user ID from AuthBloc
        final currentUserId = context.select<AuthBloc, String?>((bloc) {
          final state = bloc.state;
          if (state.isAuthenticated) {
            return state.userClaims!.userId;
          }
          return null;
        });

        return Row(
          children: [
            /// LEFT — LIST
            Expanded(
              flex: 3,
              child: ListView.builder(
                itemCount: users.length,
                itemBuilder: (context, index) {
                  final u = users[index];
                  final current = u.id == currentUserId;

                  return ListTile(
                    title: Text(u.fullName),
                    subtitle: Text(u.email),
                    trailing: Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        ?current
                            ? Row(
                                mainAxisSize: MainAxisSize.min,
                                children: [
                                  Text('You', style: TextStyle(fontSize: 16.0)),
                                  SizedBox(width: 6.0),
                                  Icon(Icons.person),
                                ],
                              )
                            : null,
                        SizedBox(width: 20.0),
                        Text(
                          u.isActive ? 'Active' : 'Inactive',
                          style: TextStyle(
                            color: u.isActive ? Colors.green : Colors.red,
                          ),
                        ),
                      ],
                    ),
                    onTap: () {
                      context.read<AdminBloc>().add(SelectUserEvent(u.id));
                    },
                  );
                },
              ),
            ),

            const VerticalDivider(),

            /// RIGHT — DETAILS
            Expanded(flex: 3, child: _UserDetails()),
          ],
        );
      },
    );
  }
}

class _UserDetails extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return BlocBuilder<AdminBloc, AdminState>(
      builder: (context, state) {
        final userState = state.selectedUser;

        if (userState.status == LoadStatus.loading) {
          return const Center(child: CircularProgressIndicator());
        }

        if (userState.data == null) {
          return const Center(child: Text('Select user'));
        }

        final user = userState.data!;

        final roles = user.roles.map((e) => e.roleCode).toList();

        return Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(user.fullName, style: const TextStyle(fontSize: 20)),
              Text(user.email),

              const SizedBox(height: 16),

              /// ACTIVE SWITCH
              SwitchListTile(
                title: const Text('Active'),
                value: user.isActive,
                onChanged: (value) {
                  context.read<AdminBloc>().add(
                    SetUserActiveEvent(user.id, value),
                  );
                },
              ),

              const SizedBox(height: 16),

              /// ROLES
              const Text('Roles'),

              CheckboxListTile(
                title: const Text('Admin'),
                value: roles.contains(UserRole.admin),
                onChanged: (value) {
                  final newRoles = [...roles];
                  value == true
                      ? newRoles.add(UserRole.admin)
                      : newRoles.remove(UserRole.admin);

                  context.read<AdminBloc>().add(
                    UpdateUserRolesEvent(user.id, newRoles),
                  );
                },
              ),

              CheckboxListTile(
                title: const Text('Teacher'),
                value: roles.contains(UserRole.teacher),
                onChanged: (value) {
                  final newRoles = [...roles];
                  value == true
                      ? newRoles.add(UserRole.teacher)
                      : newRoles.remove(UserRole.teacher);

                  context.read<AdminBloc>().add(
                    UpdateUserRolesEvent(user.id, newRoles),
                  );
                },
              ),

              CheckboxListTile(
                title: const Text('Student'),
                value: roles.contains(UserRole.student),
                onChanged: (value) {
                  final newRoles = [...roles];
                  value == true
                      ? newRoles.add(UserRole.student)
                      : newRoles.remove(UserRole.student);

                  context.read<AdminBloc>().add(
                    UpdateUserRolesEvent(user.id, newRoles),
                  );
                },
              ),
            ],
          ),
        );
      },
    );
  }
}
