package notification_builder

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/4udiwe/coworking/notification-service/internal/entity"
)

var ErrUnsupportedEvent = errors.New("unsupported event")

type DefaultBuilder struct{}

func New() *DefaultBuilder {
	return &DefaultBuilder{}
}

func (b *DefaultBuilder) Build(event Event) (entity.Notification, error) {

	switch event.Type {

	case entity.BookingCreatedNotificationType:
		return b.buildBookingCreated(event)

	case entity.BookingCancelledNotificationType:
		return b.buildBookingCancelled(event)

	case entity.BookingReminderNotificationType:
		return b.buildBookingReminder(event)

	case entity.BookingExpiredNotificationType:
		return b.buildBookingExpired(event)

	default:
		return entity.Notification{}, ErrUnsupportedEvent
	}
}

func (b *DefaultBuilder) buildBookingCreated(event Event) (entity.Notification, error) {

	place := fmt.Sprintf("%v", event.Payload["placeId"])
	start := fmt.Sprintf("%v", event.Payload["startTime"])

	title := "Бронирование создано"
	body := fmt.Sprintf("Рабочее место %s забронировано на %s", place, start)

	payload, _ := json.Marshal(event.Payload)

	return entity.Notification{
		UserID: event.UserID,

		Type: entity.BookingCreatedNotificationType,

		Title: title,
		Body:  body,

		Payload: payload,
	}, nil
}

func (b *DefaultBuilder) buildBookingCancelled(event Event) (entity.Notification, error) {

	place := fmt.Sprintf("%v", event.Payload["placeId"])

	title := "Бронирование отменено"
	body := fmt.Sprintf("Бронирование рабочего места %s отменено", place)

	payload, _ := json.Marshal(event.Payload)

	return entity.Notification{
		UserID: event.UserID,

		Type: entity.BookingCancelledNotificationType,

		Title: title,
		Body:  body,

		Payload: payload,
	}, nil
}

func (b *DefaultBuilder) buildBookingReminder(event Event) (entity.Notification, error) {

	place := fmt.Sprintf("%v", event.Payload["placeId"])
	start := fmt.Sprintf("%v", event.Payload["startTime"])

	title := "Напоминание о бронировании"
	body := fmt.Sprintf("Через 10 минут начинается бронирование места %s (%s)", place, start)

	payload, _ := json.Marshal(event.Payload)

	return entity.Notification{
		UserID: event.UserID,

		Type: entity.BookingReminderNotificationType,

		Title: title,
		Body:  body,

		Payload: payload,
	}, nil
}

func (b *DefaultBuilder) buildBookingExpired(event Event) (entity.Notification, error) {

	place := fmt.Sprintf("%v", event.Payload["placeId"])
	endTime := fmt.Sprintf("%v", event.Payload["endTime"])

	title := "Время бронирования истекло"
	body := fmt.Sprintf("Бронирование места %s закончилось в %s", place, endTime)

	payload, _ := json.Marshal(event.Payload)

	return entity.Notification{
		UserID: event.UserID,

		Type: entity.BookingExpiredNotificationType,

		Title: title,
		Body:  body,

		Payload: payload,
	}, nil
}
