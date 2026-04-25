import 'dart:async';

import 'package:coworking_app/core/services/fcm_service.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../core/utils/bloc_load_state.dart';
import '../data/notification_repository.dart';
import 'notification_event.dart';
import 'notification_state.dart';

class NotificationBloc extends Bloc<NotificationEvent, NotificationState> {
  final NotificationRepository repository;
  final FCMService? fcmService; // nullable — не используется на десктопе/вебе

  Timer? _pollTimer;
  StreamSubscription? _fcmSubscription;

  NotificationBloc({required this.repository, this.fcmService})
    : super(const NotificationState()) {
    on<FetchNotifications>(_onFetchNotifications);
    on<PollNotifications>(_onPollNotifications);
    on<FetchUnreadCount>(_onFetchUnreadCount);
    on<MarkNotificationRead>(_onMarkRead);
    on<MarkAllNotificationsRead>(_onMarkAllRead);
    on<ClearAction>(_onClearAction);
    on<FCMMessageReceived>(_onFCMMessageReceived);

    _initFCMListener();
  }

  void _initFCMListener() {
    if (fcmService == null) return;

    // Foreground: приложение открыто
    _fcmSubscription = fcmService!.onForegroundMessage.listen((message) {
      add(FCMMessageReceived());
    });

    // Background tap: юзер тапнул на уведомление
    fcmService!.onMessageOpenedApp.listen((message) {
      add(FCMMessageReceived());
    });
  }

  void startPolling() {
    _pollTimer?.cancel();

    _pollTimer = Timer.periodic(const Duration(seconds: 10), (_) {
      add(PollNotifications());
    });
  }

  void stopPolling() {
    _pollTimer?.cancel();
  }

  @override
  Future<void> close() {
    _pollTimer?.cancel();
    _fcmSubscription?.cancel();
    return super.close();
  }

  Future<void> _onFCMMessageReceived(
    FCMMessageReceived event,
    Emitter<NotificationState> emit,
  ) async {
    // Просто рефрешим список и счётчик — бекенд уже сохранил уведомление
    add(FetchNotifications(refresh: true));
    add(FetchUnreadCount());

    emit(
      state.copyWith(
        actionMessage: () => "new_notifications",
        messageId: DateTime.now().millisecondsSinceEpoch.toString(),
      ),
    );
  }

  Future<void> _onFetchNotifications(
    FetchNotifications event,
    Emitter<NotificationState> emit,
  ) async {
    final current = state.notifications.data ?? [];
    if (!event.refresh) {
      emit(
        state.copyWith(
          notifications: state.notifications.copyWith(
            status: LoadStatus.loading,
          ),
        ),
      );
    }

    try {
      final fetched = await repository.getNotifications(
        limit: event.limit,
        offset: event.offset,
      );
      final combined = event.refresh
          ? fetched.notifications
          : [...current, ...fetched.notifications];

      emit(
        state.copyWith(
          notifications: LoadState(data: combined, status: LoadStatus.success),
        ),
      );
    } catch (e) {
      emit(
        state.copyWith(
          notifications: LoadState(
            data: current,
            status: LoadStatus.error,
            error: e.toString(),
          ),
        ),
      );
    }
  }

  /// =======================
  /// POLLING
  /// =======================
  Future<void> _onPollNotifications(
    PollNotifications event,
    Emitter<NotificationState> emit,
  ) async {
    try {
      final response = await repository.getNotifications(
        since: state.lastFetchedAt,
        limit: 10,
      );

      final current = state.notifications.data ?? [];

      // merge (dedupe)
      final map = {for (var n in current) n.id: n};

      for (final n in response.notifications) {
        map[n.id] = n;
      }

      final merged = map.values.toList()
        ..sort((a, b) => b.createdAt.compareTo(a.createdAt));

      emit(
        state.copyWith(
          notifications: LoadState(data: merged, status: LoadStatus.success),
          lastFetchedAt: () => response.notifications.isNotEmpty
              ? response.notifications.first.createdAt.add(
                  Duration(milliseconds: 1),
                )
              : state.lastFetchedAt,
        ),
      );

      // 👉 можно триггерить UI (toast)
      if (response.notifications.isNotEmpty) {
        emit(
          state.copyWith(
            actionMessage: () => "new_notifications",
            messageId: DateTime.now().millisecondsSinceEpoch.toString(),
          ),
        );
      }
      add(FetchUnreadCount());
    } catch (_) {}
  }

  /// =======================
  /// UNREAD COUNT
  /// =======================
  Future<void> _onFetchUnreadCount(
    FetchUnreadCount event,
    Emitter<NotificationState> emit,
  ) async {
    try {
      final count = await repository.getUnreadCount();

      emit(state.copyWith(unreadCount: count));
    } catch (_) {}
  }

  /// =======================
  /// MARK READ
  /// =======================
  Future<void> _onMarkRead(
    MarkNotificationRead event,
    Emitter<NotificationState> emit,
  ) async {
    try {
      await repository.markRead(event.id);

      final updated = state.notifications.data?.map((n) {
        if (n.id == event.id) {
          return n.copyWith(isRead: true);
        }
        return n;
      }).toList();

      emit(
        state.copyWith(
          notifications: state.notifications.copyWith(data: updated),
        ),
      );

      add(FetchUnreadCount());
    } catch (_) {}
  }

  /// =======================
  /// MARK ALL READ
  /// =======================
  Future<void> _onMarkAllRead(
    MarkAllNotificationsRead event,
    Emitter<NotificationState> emit,
  ) async {
    try {
      await repository.readAll();

      final updated = state.notifications.data
          ?.map((n) => n.copyWith(isRead: true))
          .toList();

      emit(
        state.copyWith(
          notifications: state.notifications.copyWith(data: updated),
          unreadCount: 0,
        ),
      );
    } catch (_) {}
  }

  void _onClearAction(ClearAction event, Emitter<NotificationState> emit) {
    emit(state.copyWith(actionMessage: () => null, isError: false));
  }
}
