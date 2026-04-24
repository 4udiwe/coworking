import 'package:coworking_app/features/auth/presentation/screens/login_screen.dart';
import 'package:coworking_app/features/user/presentation/profile_screen.dart';
import 'package:coworking_app/features/auth/presentation/screens/register_screen.dart';
import 'package:coworking_app/features/bookings/presentation/screens/bookings_screen.dart';
import 'package:coworking_app/features/user/presentation/sessions_page.dart';
import 'package:coworking_app/main_screen.dart';
import 'package:flutter/material.dart';

import '../../features/coworking/presentation/screens/coworking_list_screen.dart';

class AppRouter {
  static Route<dynamic> generateRoute(RouteSettings settings) {
    switch (settings.name) {

      case '/login':
        return MaterialPageRoute(builder: (_) => const LoginScreen());
      case '/register':
        return MaterialPageRoute(builder: (_) => const RegisterScreen());

      case '/main':
        return MaterialPageRoute(
          builder: (_) => const MainScreen(),
        );

      case '/coworkings':
        return MaterialPageRoute(
          builder: (_) => const CoworkingPage(),
        );

      case '/bookings':
        final uri = Uri.parse(settings.name!);

        final tab = uri.queryParameters['tab'];
        final bookingId = uri.queryParameters['bookingId'];

        return MaterialPageRoute(
          builder: (_) => BookingsPage(
            initialTab: tab,
            highlightBookingId: bookingId,
          ),
        );

      case '/profile':
        return MaterialPageRoute(
          builder: (_) => const ProfileScreen(),
        );

      case '/sessions':
        return MaterialPageRoute(
          builder: (_) => const SessionsPage(),
        );

      default:
        return MaterialPageRoute(
          builder: (_) => const MainScreen(),
        );
    }
  }
}