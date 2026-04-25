import 'dart:io';
import 'package:coworking_app/features/notification/data/notification_repository.dart';
import 'package:firebase_messaging/firebase_messaging.dart';

@pragma('vm:entry-point')
Future<void> firebaseMessagingBackgroundHandler(RemoteMessage message) async {
  // Вызывается когда приложение в фоне/закрыто
  // Здесь нельзя обращаться к BLoC — только легковесная логика
  print('Background message: ${message.messageId}');
}

class FCMService {
  final FirebaseMessaging _messaging = FirebaseMessaging.instance;
  final NotificationRepository repository;

  FCMService({required this.repository});

  Future<void> initialize() async {
    // Запрашиваем разрешение (iOS обязательно, Android 13+)
    final settings = await _messaging.requestPermission(
      alert: true,
      badge: true,
      sound: true,
    );

    if (settings.authorizationStatus == AuthorizationStatus.denied) return;

    // Получаем и регистрируем токен
    await registerToken();

    // Слушаем обновления токена (токен может меняться)
    _messaging.onTokenRefresh.listen(_onTokenRefresh);
  }

  Future<void> registerToken() async {
    try {
      final token = await _messaging.getToken();
      if (token == null) return;

      final platform = Platform.isIOS ? 'ios' : 'android';
      await repository.registerDeviceToken(token: token, platform: platform);
    } catch (e) {
      print('FCM token registration failed: $e');
    }
  }

  Future<void> _onTokenRefresh(String token) async {
    try {
      final platform = Platform.isIOS ? 'ios' : 'android';
      await repository.registerDeviceToken(token: token, platform: platform);
    } catch (e) {
      print('FCM token refresh failed: $e');
    }
  }

  /// Вызвать из UI-слоя — передаёт сообщения в BLoC
  Stream<RemoteMessage> get onForegroundMessage => FirebaseMessaging.onMessage;

  /// Когда юзер тапнул на уведомление (приложение было в фоне)
  Stream<RemoteMessage> get onMessageOpenedApp =>
      FirebaseMessaging.onMessageOpenedApp;

  /// Получить сообщение которое открыло приложение из закрытого состояния
  Future<RemoteMessage?> getInitialMessage() => _messaging.getInitialMessage();
}
