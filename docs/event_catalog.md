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
| notification.events | Уведомления пользователям    | notification-service |
| scheduler.events	  | Таймеры и отложенные события | scheduler-service    |
| analytics.events	  | Сырые события для аналитики	 | booking-service      |


# TOPIC: booking.events
## booking.created
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

## booking.cancelled
- Описание: Бронирование отменено
- Публикует: booking-service
- Слушают: notification, analytics

```json
{
  "bookingId": "UUID",
  "reason": "string"
}
```

## booking.completed
- Описание: Бронирование завершено по времени
- Публикует: scheduler-service
- Слушают: booking-service, analytics

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

# TOPIC: notification.events
## notification.sent
- Описание: Уведомление отправлено пользователю
- Публикует: notification-service
- Слушают: analytics-service

```json
{
  "userId": "UUID",
  "channel": "push | email",
  "type": "booking_reminder"
}
```

# TOPIC: analytics.events
## analytics.booking.recorded
- Описание: Событие бронирования записано в аналитику
- Публикует: analytics-service
- Слушают: (нет)

```json
{
  "bookingId": "UUID",
  "coworkingId": "UUID",
  "hour": 14
}
```