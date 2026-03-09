package worker

import (
	"encoding/json"

	"github.com/4udiwe/cowoking/scheduler-service/internal/entity"
)

// TimerToEventMapper конвертирует entity.Timer в entity.OutboxEvent
type TimerToEventMapper struct{}

func (m *TimerToEventMapper) Map(timer entity.Timer) entity.OutboxEvent {
	switch timer.Type.ID {
	case entity.TimerTypeBookingReminderID:
		if timer.UserID == nil {
			panic("booking reminder timer has nil userID")
		}
		payload := ReminderPayload{
			BookingID: timer.BookingID,
			UserID:    *timer.UserID,
		}
		data, _ := json.Marshal(payload)
		var payloadMap map[string]any
		json.Unmarshal(data, &payloadMap)
		return entity.OutboxEvent{
			AggregateType: "reminder",
			AggregateID:   timer.ID,
			EventType:     "triggered",
			Payload:       payloadMap,
		}
	case entity.TimerTypeBookingExpireID:
		payload := ExpirePayload{
			BookingID: timer.BookingID,
		}
		data, _ := json.Marshal(payload)
		var payloadMap map[string]any
		json.Unmarshal(data, &payloadMap)
		return entity.OutboxEvent{
			AggregateType: "booking",
			AggregateID:   timer.ID,
			EventType:     "expire",
			Payload:       payloadMap,
		}
	default:
		// Любой новый таймер сразу будет заметен в логах
		panic("unknown timer type: " + string(timer.Type.Name))
	}
}
