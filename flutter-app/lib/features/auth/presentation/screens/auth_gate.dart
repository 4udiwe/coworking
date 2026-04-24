import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../../main_screen.dart';
import '../../bloc/auth_bloc.dart';
import '../../bloc/auth_state.dart';
import 'login_screen.dart';

class AuthGate extends StatelessWidget {
  const AuthGate({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocListener<AuthBloc, AuthState>(
      // BlocListener для показа snackbar/диалога
      listenWhen: (prev, curr) => curr.status == AuthStatus.sessionExpired,
      listener: (context, state) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Сессия завершена на другом устройстве'),
          ),
        );
      },
      child: BlocBuilder<AuthBloc, AuthState>(
        builder: (context, state) {
          if (state.status == AuthStatus.loading ||
              state.status == AuthStatus.initial) {
            return const Scaffold(
              body: Center(child: CircularProgressIndicator()),
            );
          }

          if (state.status == AuthStatus.authenticated) {
            return const MainScreen();
          }

          return const LoginScreen(); // покрывает unauthenticated и sessionExpired
        },
      ),
    );
  }
}
