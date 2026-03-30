# Каталог событий и топиков Kafka  
ИС «Коворкинг»

Документ является **единым контрактом** событийного взаимодействия микросервисов.

# Event Envelope (единый формат всех сообщений Kafka)
```json
{
  "eventId": "UUID",
  "eventType": "string",
  "occurredAt": "RFC3339",
  "data": {}
}
```
- event_id — уникальный UUID события
- event_type — тип события (например: "booking.created")
- occurred_at — точное время возникновения события в домене
- data — конкретный payload (структура описана ниже для каждого события)

# Topics
Ниже перечислены **все топики Kafka**, используемые в системе.

| Topic	              | Описание	                   | Публикует            |
|---------------------|------------------------------|----------------------|
| booking.events	    | Жизненный цикл бронирований  | booking-service      |
| auth.events	        | События аутентификации	     | scheduler-service    |
| notification.events | Уведомления пользователям    | notification-service |
| scheduler.events	  | Таймеры и отложенные события | scheduler-service    |


# TOPIC: booking.events
## booking.booking.created
- Описание: Создано новое бронирование
- Публикует: booking-service
- Слушают: notification, analytics, scheduler

```json
{
  "bookingId": "UUID",
  "userId": "UUID",
  "placeId": "UUID",
  "startTime": "RFC3339",
  "endTime": "RFC3339"
}
```

## booking.booking.cancelled
- Описание: Бронирование отменено
- Публикует: booking-service
- Слушают: notification, analytics, scheduler

```json
{
  "bookingId": "UUID",
  "reason": "string"
}
```

## booking.booking.completed
- Описание: Бронирование завершено по времени
- Публикует: booking-service
- Слушают: notification, analytics

```json
{
  "bookingId": "UUID"
}
```

# TOPIC: scheduler.events
## scheduler.reminder.triggered
- Описание: Сработало напоминание о начале бронирования
- Публикует: scheduler-service
- Слушают: notification-service

```json
{
  "bookingId": "UUID",
  "userId": "UUID"
}
```

## scheduler.booking.expire
- Описание: Запрос на изменение статуса бронирования на "Завершено"
- Публикует: scheduler-service
- Слушают: booking-service

```json
{
  "bookingId": "UUID"
}
```


# TOPIC: notification.events
## notification.sent
- Описание: Уведомление отправлено пользователю
- Публикует: notification-service
- Слушают: notification, analytics-service

```json
{
  "userId": "UUID",
  "notificationType": "booking_reminder",
  "notificationId": "UUID"
}
```

# TOPIC: auth.events
## auth.sessions.cleanup
- Описание: Запрос на очистку старых revoked сессий
- Публикует: scheduler-service
- Слушают: auth-service
- Частота: 2 раза в день (каждые 12 часов)

```json
{
  "retentionDays": 10
}
```

**Описание параметров:**
- `retentionDays` — удалять revoked сессии, которым больше этого количества дней

**Логика в auth-service:**
- Получает событие с указанным retentionDays
- Выполняет: `DELETE FROM refresh_tokens WHERE revoked = true AND created_at < now() - INTERVAL 'retentionDays days'`
- Логирует количество удалённых записей
