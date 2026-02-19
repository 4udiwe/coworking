package booking_service

import "errors"

var (
	ErrCoworkingAlreadyExists     = errors.New("coworking already exists")
	ErrCoworkingNotFound          = errors.New("coworking not found")
	ErrCoworkingHasActiveBookings = errors.New("cannot deactivate coworking with active bookings")
	ErrInvalidLayoutSchema        = errors.New("invalid layout schema")
	ErrInvalidLayoutSchemaVersion = errors.New("invalid layout schema version")

	ErrCannotCreateCoworking = errors.New("cannot create coworking")
	ErrCannotUpdateCoworking = errors.New("cannot update coworking")
	ErrCannotFetchCoworking  = errors.New("cannot fetch coworking")
	ErrCannotCreateLayout    = errors.New("cannot create layout version")
	ErrCannotFetchLayout     = errors.New("cannot fetch layout version")

	ErrPlaceAlreadyExists     = errors.New("place already exists")
	ErrPlaceNotFound          = errors.New("place not found")
	ErrPlaceHasActiveBookings = errors.New("cannot deactivate place with active bookings")

	ErrCannotCreatePlace = errors.New("cannot create place")
	ErrCannotUpdatePlace = errors.New("cannot update place")
	ErrCannotFetchPlace  = errors.New("cannot fetch place")

	ErrBookingStartTimeAfterEndTime      = errors.New("booking start time cannot be after end time")
	ErrBookingStartTimeInPast            = errors.New("booking start time cannot be in the past")
	ErrBookingStartTimeEqualEndTime      = errors.New("booking start time cannot be equal to end time")
	ErrBookingTimeNotMultipleOfHour      = errors.New("booking start and end times must be multiples of an hour")
	ErrBookingDurationLessThanOneHour    = errors.New("booking duration must be at least one hour")
	ErrBookingDurationMoreThanThreeHours = errors.New("booking duration cannot exceed three hours")
	ErrPlaceInactive                     = errors.New("cannot book an inactive place")
	ErrCoworkingInactive                 = errors.New("cannot book a place in an inactive coworking")
	ErrBookingTimeConflict               = errors.New("booking time conflicts with an existing booking")
	ErrBookingNotFound                   = errors.New("booking not found")
	ErrBookingAlreadyCancelled           = errors.New("booking is already cancelled")
	ErrBookingAlreadyCompleted           = errors.New("booking is already completed")

	ErrCannotCreateBooking   = errors.New("cannot create booking")
	ErrCannotCancelBooking   = errors.New("cannot cancel booking")
	ErrCannotCompleteBooking = errors.New("cannot complete booking")
	ErrCannotFetchBooking    = errors.New("cannot fetch booking")
)
