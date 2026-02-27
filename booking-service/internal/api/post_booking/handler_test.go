package post_booking

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/4udiwe/cowoking/booking-service/internal/api/dto"
	"github.com/4udiwe/cowoking/booking-service/internal/entity"
	booking_service "github.com/4udiwe/cowoking/booking-service/internal/service/booking"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type mockBookingService struct {
	err error
}

func (m *mockBookingService) CreateBooking(_ context.Context, _ entity.Booking) error {
	return m.err
}

func TestPostBookingHandler_ValidationError(t *testing.T) {
	e := echo.New()

	h := &handler{s: &mockBookingService{err: booking_service.ErrBookingStartTimeAfterEndTime}}

	req := httptest.NewRequest(http.MethodPost, "/bookings", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	now := time.Now().UTC()
	in := dto.CreateBookingRequest{
		UserID:      uuid.New(),
		PlaceID:     uuid.New(),
		CoworkingID: uuid.New(),
		StartTime:   now.Add(2 * time.Hour),
		EndTime:     now.Add(1 * time.Hour),
	}

	err := h.Handle(c, in)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	httpErr, ok := err.(*echo.HTTPError)
	if !ok {
		t.Fatalf("expected *echo.HTTPError, got %T", err)
	}
	if httpErr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, httpErr.Code)
	}
	if !errors.Is(httpErr, err) {
		// just to make sure error is propagated
		t.Fatalf("expected same error to be returned")
	}
}

