package booking_service

import (
	"context"
	"testing"
	"time"

	"github.com/4udiwe/cowoking/booking-service/internal/entity"
)

// dummyTransactor is used to test validation logic that returns
// before any transaction is started.
type dummyTransactor struct{}

func (d dummyTransactor) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestCreateBooking_ValidationErrors(t *testing.T) {
	now := time.Now().UTC()

	type testCase struct {
		name      string
		start     time.Time
		end       time.Time
		wantError error
	}

	tests := []testCase{
		{
			name:      "start_after_end",
			start:     now.Add(2 * time.Hour),
			end:       now.Add(1 * time.Hour),
			wantError: ErrBookingStartTimeAfterEndTime,
		},
		{
			name:      "start_equal_end",
			start:     now.Add(1 * time.Hour),
			end:       now.Add(1 * time.Hour),
			wantError: ErrBookingStartTimeEqualEndTime,
		},
		{
			name:      "start_in_past",
			start:     now.Add(-1 * time.Hour),
			end:       now.Add(1 * time.Hour),
			wantError: ErrBookingStartTimeInPast,
		},
	}

	svc := &BookingService{
		txManager: dummyTransactor{},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := svc.CreateBooking(context.Background(), entity.Booking{
				StartTime: tc.start,
				EndTime:   tc.end,
			})
			if err != tc.wantError {
				t.Fatalf("expected error %v, got %v", tc.wantError, err)
			}
		})
	}
}

