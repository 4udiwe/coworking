package worker

import (
	"encoding/json"

	"github.com/4udiwe/cowoking/scheduler-service/internal/entity"
	"github.com/sirupsen/logrus"
)

// TimerToEventMapper конвертирует entity.Timer в entity.OutboxEvent
type TimerToEventMapper struct{}

func (m *TimerToEventMapper) Map(timer entity.Timer) entity.OutboxEvent {
	switch timer.Type.ID {
	case entity.TimerTypeBookingReminderID:
		if timer.UserID == nil {
			logrus.Error("booking reminder timer has nil userID")
		}
		payload := ReminderPayload{
			BookingID:  timer.BookingID,
			UserID:     *timer.UserID,
			PlaceID:    timer.PlaceID,
			PlaceLabel: timer.PlaceLabel,
			StartTime:  timer.StartTime,
			EndTime:    timer.EndTime,
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
		if timer.UserID == nil {
			logrus.Error("booking expire timer has nil userID")
		}
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
		logrus.Error("unknown timer type to map in scheduler worker: " + string(timer.Type.Name))
		return entity.OutboxEvent{}
	}
}
