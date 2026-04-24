import 'package:coworking_app/core/models/user_session.dart';
import 'package:coworking_app/features/user/bloc/user_event.dart';
import 'package:coworking_app/features/user/bloc/user_state.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../core/utils/bloc_load_state.dart';
import '../data/user_repository.dart';

class UserBloc extends Bloc<UserEvent, UserState> {
  final UserRepository repository;

  UserBloc({required this.repository}) : super(const UserState()) {
    on<LoadProfile>(_onLoadProfile);
    on<LoadActiveSessions>(_onLoadActiveSessions);
    on<LoadAllSessions>(_onLoadAllSessions);
    on<RevokeSessionEvent>(_onRevokeSession);
    on<RefreshUserData>(_onRefresh);
  }

  Future<void> _onLoadProfile(
    LoadProfile event,
    Emitter<UserState> emit,
  ) async {
    emit(state.copyWith(profile: const LoadState(status: LoadStatus.loading)));

    try {
      final user = await repository.getProfile();

      emit(
        state.copyWith(
          profile: LoadState(status: LoadStatus.success, data: user),
        ),
      );
    } catch (e) {
      emit(
        state.copyWith(
          profile: LoadState(status: LoadStatus.error, error: e.toString()),
          actionMessage: "Failed to load profile: ${e.toString()}",
          messageId: DateTime.now().toString(),
          isError: true,
        ),
      );
    }
  }

  Future<void> _onLoadActiveSessions(
    LoadActiveSessions event,
    Emitter<UserState> emit,
  ) async {
    emit(
      state.copyWith(
        activeSessions: const LoadState(status: LoadStatus.loading),
      ),
    );

    try {
      final sessions = await repository.getActiveSessions();

      emit(
        state.copyWith(
          activeSessions: LoadState(status: LoadStatus.success, data: sessions),
        ),
      );
    } catch (e) {
      emit(
        state.copyWith(
          activeSessions: LoadState(
            status: LoadStatus.error,
            error: e.toString(),
          ),
          actionMessage: "Failed to load active sessions: ${e.toString()}",
          messageId: DateTime.now().toString(),
          isError: true,
        ),
      );
    }
  }

  Future<void> _onLoadAllSessions(
    LoadAllSessions event,
    Emitter<UserState> emit,
  ) async {
    emit(
      state.copyWith(allSessions: const LoadState(status: LoadStatus.loading)),
    );

    try {
      final sessions = await repository.getAllSessions();

      /* сортировка сессий
      cначала current
      затем active
      затем revoked

      также упорядочивание по lastUsedAt внутри группы
      */
      sessions.sort((a, b) {
        int priority(UserSession s) {
          if (s.current) return 0;
          if (!s.revoked) return 1;
          return 2;
        }

        final pCompare = priority(a).compareTo(priority(b));
        if (pCompare != 0) return pCompare;

        return b.lastUsedAt.compareTo(a.lastUsedAt);
      });

      emit(
        state.copyWith(
          allSessions: LoadState(status: LoadStatus.success, data: sessions),
        ),
      );
    } catch (e) {
      emit(
        state.copyWith(
          allSessions: LoadState(status: LoadStatus.error, error: e.toString()),
          actionMessage: e.toString(),
          messageId: DateTime.now().toString(),
          isError: true,
        ),
      );
    }
  }

  Future<void> _onRevokeSession(
    RevokeSessionEvent event,
    Emitter<UserState> emit,
  ) async {
    try {
      await repository.revokeSession(event.sessionId);
      emit(
        state.copyWith(
          actionMessage: "Session revoked",
          messageId: DateTime.now().toString(),
          isError: false,
        ),
      );

      /// 🔥 обновляем только активные сессии
      add(LoadActiveSessions());
      add(LoadAllSessions());
    } catch (e) {
      emit(
        state.copyWith(
          activeSessions: LoadState(
            status: LoadStatus.error,
            error: e.toString(),
          ),
          actionMessage: e.toString(),
          messageId: DateTime.now().toString(),
          isError: true,
        ),
      );
    }
  }

  Future<void> _onRefresh(
    RefreshUserData event,
    Emitter<UserState> emit,
  ) async {
    add(LoadProfile());
    add(LoadActiveSessions());
  }
}
