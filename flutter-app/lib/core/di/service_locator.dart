import 'dart:convert';
import 'dart:io';

import 'package:coworking_app/core/models/auth_tokens.dart';
import 'package:coworking_app/core/services/fcm_service.dart';
import 'package:coworking_app/features/admin/data/media_reposotory.dart';
import 'package:coworking_app/features/auth/bloc/auth_event.dart';
import 'package:coworking_app/features/notification/bloc/notification_bloc.dart';
import 'package:coworking_app/features/user/data/user_repository.dart';
import 'package:flutter/foundation.dart';
import 'package:get_it/get_it.dart';
import 'package:http/http.dart' as http;

import '../../features/admin/bloc/admin_bloc.dart';
import '../../features/admin/data/admin_repository.dart';
import '../../features/notification/data/notification_repository.dart';
import '../../features/user/bloc/user_bloc.dart';
import '../api/api_client.dart';
import '../storage/token_storage.dart';

import '../../features/auth/data/auth_repository.dart';
import '../../features/coworking/data/coworking_repository.dart';
import '../../features/bookings/data/booking_repository.dart';
import '../../features/analytics/data/analytics_repository.dart';

import '../../features/auth/bloc/auth_bloc.dart';
import '../../features/coworking/bloc/coworking_bloc.dart';
import '../../features/bookings/bloc/booking_bloc.dart';
import '../../features/analytics/bloc/analytics_bloc.dart';

final sl = GetIt.instance;

Future<void> init() async {
  // TokenStorage
  sl.registerLazySingleton<TokenStorage>(() => TokenStorage());

  // ApiClient
  sl.registerLazySingleton<ApiClient>(() {
    final baseUrl = kIsWeb ? "http://localhost:8080" : "http://10.0.2.2:8080";

    return ApiClient(
      baseUrl: baseUrl,
      tokenStorage: sl(),
      onRefresh: (refreshToken) async {
        final response = await http.post(
          Uri.parse("$baseUrl/auth/refresh"),
          headers: {"Content-Type": "application/json"},
          body: jsonEncode({"refreshToken": refreshToken}),
        );
        if (response.statusCode == 200) {
          return AuthTokens.fromJson(jsonDecode(response.body));
        }
        throw Exception("refresh failed: ${response.statusCode}");
      },
      onSessionExpired: () {
        // AuthBloc уже создан к этому моменту (lazySingleton)
        sl<AuthBloc>().add(AuthSessionExpired());
      },
    );
  });

  // -----------------------------
  // Repositories
  // -----------------------------
  sl.registerLazySingleton<AuthRepository>(
    () => AuthRepository(apiClient: sl(), tokenStorage: sl()),
  );

  sl.registerLazySingleton<AdminRepository>(
    () => AdminRepository(apiClient: sl()),
  );

  sl.registerLazySingleton<CoworkingRepository>(
    () => CoworkingRepository(apiClient: sl()),
  );

  sl.registerLazySingleton<BookingRepository>(
    () => BookingRepository(apiClient: sl()),
  );

  sl.registerLazySingleton<AnalyticsRepository>(
    () => AnalyticsRepository(apiClient: sl()),
  );

  sl.registerLazySingleton<UserRepository>(
    () => UserRepository(apiClient: sl()),
  );

  sl.registerLazySingleton<NotificationRepository>(
    () => NotificationRepository(apiClient: sl()),
  );

  sl.registerLazySingleton<MediaRepository>(
    () => MediaRepository(apiClient: sl()),
  );

  // -----------------------------
  // BLoC
  // -----------------------------
  sl.registerLazySingleton<AuthBloc>(
    () => AuthBloc(
      authRepository: sl(),
      tokenStorage: sl(),
      fcmService: (!kIsWeb && (Platform.isAndroid || Platform.isIOS))
          ? sl<FCMService>()
          : null,
    ),
  );

  sl.registerFactory<AdminBloc>(() => AdminBloc(repository: sl()));

  sl.registerFactory<CoworkingBloc>(() => CoworkingBloc(repository: sl()));

  sl.registerFactory<BookingBloc>(() => BookingBloc(repository: sl()));

  sl.registerFactory<AnalyticsBloc>(() => AnalyticsBloc(repository: sl()));

  sl.registerFactory<UserBloc>(() => UserBloc(repository: sl()));

  // FCM Service
  sl.registerLazySingleton<FCMService>(
    () => FCMService(repository: sl<NotificationRepository>()),
  );

  sl.registerLazySingleton<NotificationBloc>(
    () => NotificationBloc(repository: sl<NotificationRepository>()),
  );
}
