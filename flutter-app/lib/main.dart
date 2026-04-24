import 'dart:io';

import 'package:coworking_app/features/auth/bloc/auth_event.dart';
import 'package:coworking_app/features/auth/presentation/screens/auth_gate.dart';
import 'package:coworking_app/features/notification/bloc/notification_bloc.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:intl/date_symbol_data_local.dart';
import 'package:provider/provider.dart';

import 'core/di/service_locator.dart';
import 'core/l10n/locale_provider.dart';
import 'core/navigation/app_router.dart';
import 'features/auth/bloc/auth_bloc.dart';
import 'generated_l10n/app_localizations.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  await init();

  await initializeDateFormatting('ru', null);

  final localeProvider = LocaleProvider();
  await localeProvider.initialize();

  runApp(
    ChangeNotifierProvider(
      create: (_) => localeProvider,
      child: MultiBlocProvider(
        providers: [
          BlocProvider<AuthBloc>(
            create: (_) => sl<AuthBloc>()..add(AuthCheckSession()),
          ),
          BlocProvider<NotificationBloc>(
            create: (_) {
              final bloc = sl<NotificationBloc>();
              // Only start polling on desktop platforms (Windows, Linux, macOS)
              // Web and mobile (Android, iOS) don't use polling
              if (!kIsWeb && !Platform.isAndroid && !Platform.isIOS) {
                return bloc..startPolling();
              }
              return bloc;
            },
          ),
        ],
        child: const MyApp(),
      ),
    ),
  );
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return Consumer<LocaleProvider>(
      builder: (context, localeProvider, child) {
        return MaterialApp(
          title: 'Coworking App',
          locale: localeProvider.currentLocale,
          localizationsDelegates: AppLocalizations.localizationsDelegates,
          supportedLocales: AppLocalizations.supportedLocales,
          onGenerateRoute: AppRouter.generateRoute,
          home: const AuthGate(),
          theme: ThemeData(
            colorScheme: ColorScheme.fromSeed(seedColor: Colors.blue),
          ),
        );
      },
    );
  }
}
