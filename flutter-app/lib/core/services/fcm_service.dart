import 'dart:convert';
import 'dart:io';
import 'dart:ui';

import 'package:coworking_app/features/notification/data/notification_repository.dart';
import 'package:firebase_core/firebase_core.dart';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:http/http.dart' as http;

// ===========================================================================
// TOP-LEVEL — работают в отдельном изоляте, нет доступа к GetIt / BLoC
// ===========================================================================

@pragma('vm:entry-point')
Future<void> firebaseMessagingBackgroundHandler(RemoteMessage message) async {
  await Firebase.initializeApp();

  final localNotifications = FlutterLocalNotificationsPlugin();
  await localNotifications.initialize(
    const InitializationSettings(
      android: AndroidInitializationSettings('@mipmap/ic_launcher'),
      iOS: DarwinInitializationSettings(),
    ),
    // обработчики тапов нужны и здесь — приложение может быть закрыто
    onDidReceiveNotificationResponse: onNotificationTap,
    onDidReceiveBackgroundNotificationResponse: onBackgroundNotificationTap,
  );

  final data = message.data;

  await localNotifications.show(
    data['notificationId']?.hashCode ?? 0,
    data['title'],
    data['body'],
    const NotificationDetails(
      android: AndroidNotificationDetails(
        'default_channel',
        'Уведомления',
        importance: Importance.high,
        priority: Priority.high,
        actions: [
          AndroidNotificationAction(
            'mark_read',
            'Прочитано',
            cancelNotification: true,
          ),
        ],
      ),
    ),
    payload: jsonEncode({
      'notificationId': data['notificationId'],
      'actionUrl': data['actionUrl'],
    }),
  );
}

/// Вызывается когда приложение ОТКРЫТО и юзер тапнул на уведомление/кнопку.
/// Может обращаться к FCMService через статический коллбэк.
@pragma('vm:entry-point')
void onNotificationTap(NotificationResponse response) {
  FCMService.handleNotificationResponse(response);
}

/// Вызывается когда приложение В ФОНЕ/ЗАКРЫТО и юзер нажал action кнопку.
/// Отдельный изолят — только прямой HTTP, без GetIt/BLoC.
@pragma('vm:entry-point')
Future<void> onBackgroundNotificationTap(NotificationResponse response) async {
  if (response.actionId != 'mark_read') return;

  final payload = response.payload != null
      ? jsonDecode(response.payload!) as Map<String, dynamic>
      : null;

  final id = payload?['notificationId'] as String?;
  if (id == null) return;

  const storage = FlutterSecureStorage();
  final token = await storage.read(
    key: 'access_token',
  ); // твой ключ в TokenStorage
  if (token == null) return;

  try {
    await http.patch(
      Uri.parse('http://10.0.2.2:8080/notifications/$id'),
      headers: {'Authorization': 'Bearer $token'},
    );
  } catch (e) {
    // нет возможности логировать в UI — просто игнорируем
  }
}

// ===========================================================================
// FCMService
// ===========================================================================

class FCMService {
  final FlutterLocalNotificationsPlugin _localNotifications =
      FlutterLocalNotificationsPlugin();
  final FirebaseMessaging _messaging = FirebaseMessaging.instance;
  final NotificationRepository repository;

  // Статические коллбэки — доступны из top-level onNotificationTap
  static void Function(String actionUrl)? onNavigate;
  static void Function(String notificationId)? onMarkRead;
  static VoidCallback? onNewMessage;

  FCMService({required this.repository});

  Future<void> initialize() async {
    final settings = await _messaging.requestPermission(
      alert: true,
      badge: true,
      sound: true,
    );

    if (settings.authorizationStatus == AuthorizationStatus.denied) return;

    await _initLocalNotifications();

    // Foreground — приложение открыто, FCM не показывает уведомление сам
    FirebaseMessaging.onMessage.listen(_onForegroundMessage);

    // Background — юзер тапнул на системное уведомление
    FirebaseMessaging.onMessageOpenedApp.listen((msg) {
      _handleTap(msg.data['notificationId'], msg.data['actionUrl']);
    });

    // Closed — приложение открылось по тапу на уведомление
    final initial = await _messaging.getInitialMessage();
    if (initial != null) {
      _handleTap(initial.data['notificationId'], initial.data['actionUrl']);
    }

    await registerToken();
    _messaging.onTokenRefresh.listen(_onTokenRefresh);
  }

  Future<void> _initLocalNotifications() async {
    await _localNotifications.initialize(
      const InitializationSettings(
        android: AndroidInitializationSettings('@mipmap/ic_launcher'),
        iOS: DarwinInitializationSettings(),
      ),
      onDidReceiveNotificationResponse: onNotificationTap,
      onDidReceiveBackgroundNotificationResponse: onBackgroundNotificationTap,
    );

    await _localNotifications
        .resolvePlatformSpecificImplementation<
          AndroidFlutterLocalNotificationsPlugin
        >()
        ?.createNotificationChannel(
          const AndroidNotificationChannel(
            'default_channel',
            'Уведомления',
            importance: Importance.high,
          ),
        );
  }

  Future<void> _onForegroundMessage(RemoteMessage message) async {
    // final data = message.data;

    // await _localNotifications.show(
    //   data['notificationId']?.hashCode ?? 0,
    //   data['title'],
    //   data['body'],
    //   const NotificationDetails(
    //     android: AndroidNotificationDetails(
    //       'default_channel',
    //       'Уведомления',
    //       importance: Importance.high,
    //       priority: Priority.high,
    //       actions: [
    //         AndroidNotificationAction(
    //           'mark_read',
    //           'Прочитано',
    //           cancelNotification: true,
    //         ),
    //       ],
    //     ),
    //   ),
    //   payload: jsonEncode({
    //     'notificationId': data['notificationId'],
    //     'actionUrl': data['actionUrl'],
    //   }),
    // );

    onNewMessage?.call();
  }

  /// Вызывается из top-level onNotificationTap — имеет доступ к статическим коллбэкам
  static void handleNotificationResponse(NotificationResponse response) {
    final payload = response.payload != null
        ? jsonDecode(response.payload!) as Map<String, dynamic>
        : null;

    if (response.actionId == 'mark_read') {
      final id = payload?['notificationId'] as String?;
      if (id != null) onMarkRead?.call(id);
      return;
    }

    // Обычный тап — навигация
    final notificationId = payload?['notificationId'] as String?;
    final actionUrl = payload?['actionUrl'] as String?;

    if (notificationId != null) onMarkRead?.call(notificationId);
    if (actionUrl != null && actionUrl.isNotEmpty) onNavigate?.call(actionUrl);
  }

  void _handleTap(String? notificationId, String? actionUrl) {
    if (notificationId != null) FCMService.onMarkRead?.call(notificationId);
    if (actionUrl != null && actionUrl.isNotEmpty) {
      FCMService.onNavigate?.call(actionUrl);
    }
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
}
