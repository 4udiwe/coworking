package notification_builder

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/4udiwe/coworking/notification-service/internal/entity"
)

var ErrUnsupportedEvent = errors.New("unsupported event")

// StandardPayload defines the standardized notification payload structure
type StandardPayload struct {
	Type       string                 `json:"type"`
	BookingID  string                 `json:"bookingId,omitempty"`
	PlaceID    string                 `json:"placeId,omitempty"`
	PlaceLabel string                 `json:"placeLabel,omitempty"`
	StartTime  string                 `json:"startTime,omitempty"`
	EndTime    string                 `json:"endTime,omitempty"`
	Extra      map[string]interface{} `json:"extra,omitempty"`
}

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
	placeLabel := fmt.Sprintf("%v", event.Payload["placeLabel"])
	start := fmt.Sprintf("%v", event.Payload["startTime"])
	bookingID := fmt.Sprintf("%v", event.Payload["bookingId"])

	title := "Бронирование создано"
	body := fmt.Sprintf("Рабочее место %s забронировано на %s", place, start)

	// Create standardized payload
	payload := StandardPayload{
		Type:       "booking",
		BookingID:  bookingID,
		PlaceID:    place,
		PlaceLabel: placeLabel,
		StartTime:  start,
	}

	// Add any extra fields from original payload
	extraFields := make(map[string]interface{})
	for k, v := range event.Payload {
		if k != "bookingId" && k != "placeId" && k != "placeLabel" && k != "startTime" {
			extraFields[k] = v
		}
	}
	if len(extraFields) > 0 {
		payload.Extra = extraFields
	}

	payloadBytes, _ := json.Marshal(payload)

	// Construct action URL to open booking details
	actionURL := fmt.Sprintf("/bookings?tab=active&bookingId=%s", bookingID)

	return entity.Notification{
		UserID: event.UserID,

		Type: entity.BookingCreatedNotificationType,

		Title: title,
		Body:  body,

		Payload:   payloadBytes,
		ActionURL: &actionURL,
	}, nil
}

func (b *DefaultBuilder) buildBookingCancelled(event Event) (entity.Notification, error) {

	place := fmt.Sprintf("%v", event.Payload["placeId"])
	placeLabel := fmt.Sprintf("%v", event.Payload["placeLabel"])
	bookingID := fmt.Sprintf("%v", event.Payload["bookingId"])

	title := "Бронирование отменено"
	body := fmt.Sprintf("Бронирование рабочего места %s отменено", place)

	// Create standardized payload
	payload := StandardPayload{
		Type:       "booking",
		BookingID:  bookingID,
		PlaceID:    place,
		PlaceLabel: placeLabel,
	}

	payloadBytes, _ := json.Marshal(payload)

	actionURL := fmt.Sprintf("/bookings?tab=all&bookingId=%s", bookingID)

	return entity.Notification{
		UserID: event.UserID,

		Type: entity.BookingCancelledNotificationType,

		Title: title,
		Body:  body,

		Payload:   payloadBytes,
		ActionURL: &actionURL,
	}, nil
}

func (b *DefaultBuilder) buildBookingReminder(event Event) (entity.Notification, error) {

	place := fmt.Sprintf("%v", event.Payload["placeId"])
	placeLabel := fmt.Sprintf("%v", event.Payload["placeLabel"])
	start := fmt.Sprintf("%v", event.Payload["startTime"])
	bookingID := fmt.Sprintf("%v", event.Payload["bookingId"])

	title := "Напоминание о бронировании"
	body := fmt.Sprintf("Через 10 минут начинается бронирование места %s (%s)", place, start)

	// Create standardized payload
	payload := StandardPayload{
		Type:       "booking",
		BookingID:  bookingID,
		PlaceID:    place,
		PlaceLabel: placeLabel,
		StartTime:  start,
	}

	payloadBytes, _ := json.Marshal(payload)

	// Construct action URL to open booking details
	actionURL := fmt.Sprintf("/bookings?tab=active&bookingId=%s", bookingID)

	return entity.Notification{
		UserID: event.UserID,

		Type: entity.BookingReminderNotificationType,

		Title: title,
		Body:  body,

		Payload:   payloadBytes,
		ActionURL: &actionURL,
	}, nil
}

func (b *DefaultBuilder) buildBookingExpired(event Event) (entity.Notification, error) {

	place := fmt.Sprintf("%v", event.Payload["placeId"])
	placeLabel := fmt.Sprintf("%v", event.Payload["placeLabel"])
	endTime := fmt.Sprintf("%v", event.Payload["endTime"])
	bookingID := fmt.Sprintf("%v", event.Payload["bookingId"])

	title := "Время бронирования истекло"
	body := fmt.Sprintf("Бронирование места %s закончилось в %s", place, endTime)

	// Create standardized payload
	payload := StandardPayload{
		Type:       "booking",
		BookingID:  bookingID,
		PlaceID:    place,
		PlaceLabel: placeLabel,
		EndTime:    endTime,
	}

	payloadBytes, _ := json.Marshal(payload)

	actionURL := fmt.Sprintf("/bookings?tab=all&bookingId=%s", bookingID)

	return entity.Notification{
		UserID: event.UserID,

		Type: entity.BookingExpiredNotificationType,

		Title: title,
		Body:  body,

		Payload:   payloadBytes,
		ActionURL: &actionURL,
	}, nil
}
