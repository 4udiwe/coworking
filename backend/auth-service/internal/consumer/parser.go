package consumer

import (
	"encoding/json"
	"fmt"

	"github.com/4udiwe/big-bob-pizza/order-service/pkg/kafka"
)

func ParseCleanupEvent(data []byte) (*IncomingEvent, error) {
	var env kafka.Envelope

	if err := json.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("invalid envelope: %w", err)
	}

	var p Payload

	if env.Data == nil {
		return nil, fmt.Errorf("invalid payload for %s: empty data", env.EventType)
	}

	if err := json.Unmarshal(env.Data, &p); err != nil {
		return nil, fmt.Errorf("invalid payload for %s: %w", env.EventType, err)
	}

	return &IncomingEvent{
		Type:       EventType(env.EventType),
		OccurredAt: env.OccurredAt,
		Payload:    p,
	}, nil
}
