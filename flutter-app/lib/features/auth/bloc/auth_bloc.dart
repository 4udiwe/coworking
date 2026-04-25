import 'dart:io';

import 'package:bloc/bloc.dart';
import 'package:coworking_app/core/models/auth_requests.dart';
import 'package:coworking_app/core/services/fcm_service.dart';
import 'package:flutter/foundation.dart';
import '../../../core/storage/token_storage.dart';
import '../data/auth_repository.dart';
import '../data/jwt_parser.dart';
import 'auth_event.dart';
import 'auth_state.dart';

class AuthBloc extends Bloc<AuthEvent, AuthState> {
  final AuthRepository authRepository;
  final TokenStorage tokenStorage;
  final FCMService? fcmService;

  AuthBloc({
    required this.authRepository,
    required this.tokenStorage,
    this.fcmService,
  }) : super(AuthState()) {
    on<AuthRegister>(_onRegister);
    on<AuthLogin>(_onLogin);
    on<AuthLogout>(_onLogout);
    on<AuthCheckSession>(_onCheckSession);
    on<AuthRefresh>(_onRefresh);
    on<AuthSessionExpired>(_onSessionExpired);
  }

  // ----------------- Event Handlers -----------------

  Future<void> _onRegister(AuthRegister event, Emitter<AuthState> emit) async {
    emit(state.copyWith(status: AuthStatus.loading));

    try {
      final tokens = await authRepository.register(
        RegisterRequest(
          name: event.name,
          lastName: event.lastName,
          email: event.email,
          password: event.password,
          roleCode: event.role,
        ),
      );

      await tokenStorage.saveTokens(tokens);

      final claims = JwtClaimsParser.parse(tokens.accessToken);

      emit(
        state.copyWith(status: AuthStatus.authenticated, userClaims: claims),
      );
      _registerFCMToken();
    } catch (e) {
      emit(state.copyWith(status: AuthStatus.failure, error: e.toString()));
    }
  }

  Future<void> _onLogin(AuthLogin event, Emitter<AuthState> emit) async {
    emit(state.copyWith(status: AuthStatus.loading));

    try {
      final tokens = await authRepository.login(
        LoginRequest(email: event.email, password: event.password),
      );

      await tokenStorage.saveTokens(tokens);

      final claims = JwtClaimsParser.parse(tokens.accessToken);

      emit(
        state.copyWith(status: AuthStatus.authenticated, userClaims: claims),
      );
      _registerFCMToken();
    } catch (e) {
      emit(state.copyWith(status: AuthStatus.failure, error: e.toString()));
    }
  }

  Future<void> _onLogout(AuthLogout event, Emitter<AuthState> emit) async {
    emit(state.copyWith(status: AuthStatus.loading));
    try {
      final tokens = await tokenStorage.readTokens();
      if (tokens != null) {
        await authRepository.logout(tokens.refreshToken);
      }
    } catch (_) {
      // игнорируем ошибку сети — всё равно чистим локально
    } finally {
      await tokenStorage.clearTokens();
      emit(const AuthState(status: AuthStatus.unauthenticated));
    }
  }

  // Сессия убита извне — только чистим локальное состояние, без HTTP
  Future<void> _onSessionExpired(
    AuthSessionExpired event,
    Emitter<AuthState> emit,
  ) async {
    await tokenStorage.clearTokens();
    emit(const AuthState(status: AuthStatus.unauthenticated));
  }

  Future<void> _onCheckSession(
    AuthCheckSession event,
    Emitter<AuthState> emit,
  ) async {
    emit(state.copyWith(status: AuthStatus.loading));

    try {
      final tokens = await tokenStorage.readTokens();

      if (tokens == null) {
        emit(const AuthState(status: AuthStatus.unauthenticated));
        return;
      }

      /// 🔹 access живой
      if (!JwtClaimsParser.isExpired(tokens.accessToken)) {
        final claims = JwtClaimsParser.parse(tokens.accessToken);

        emit(
          state.copyWith(status: AuthStatus.authenticated, userClaims: claims),
        );
        return;
      }

      /// 🔹 access истёк → пробуем refresh
      try {
        final newTokens = await authRepository.refresh();
        await tokenStorage.saveTokens(newTokens);

        final claims = JwtClaimsParser.parse(newTokens.accessToken);

        emit(
          state.copyWith(status: AuthStatus.authenticated, userClaims: claims),
        );
      } catch (_) {
        /// ❌ refresh тоже умер
        await tokenStorage.clearTokens();
        emit(const AuthState(status: AuthStatus.unauthenticated));
      }
    } catch (_) {
      emit(const AuthState(status: AuthStatus.unauthenticated));
    }
  }

  Future<void> _onRefresh(AuthRefresh event, Emitter<AuthState> emit) async {
    try {
      final tokens = await tokenStorage.readTokens();

      if (tokens == null) {
        emit(const AuthState(status: AuthStatus.unauthenticated));
        return;
      }

      final newTokens = await authRepository.refresh();
      await tokenStorage.saveTokens(newTokens);

      final claims = JwtClaimsParser.parse(newTokens.accessToken);

      emit(
        state.copyWith(status: AuthStatus.authenticated, userClaims: claims),
      );
    } catch (_) {
      emit(const AuthState(status: AuthStatus.unauthenticated));
    }
  }

  Future<void> _registerFCMToken() async {
    if (kIsWeb) {
      return;
    }
    if (!Platform.isAndroid && !Platform.isIOS) {
      return;
    }
    if (fcmService == null) {
      return;
    }

    try {
      await fcmService!.registerToken();
    } catch (e) {
      print('FCM token register failed: $e');
    }
  }
}
