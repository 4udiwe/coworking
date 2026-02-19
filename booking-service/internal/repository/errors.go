package repository

import (
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

const (
	StatusActive    = 1
	StatusCancelled = 2
	StatusCompleted = 3
)

var (
	ErrAlreadyExists       = errors.New("booking already exists")
	ErrBookingNotFound     = errors.New("booking not found")
	ErrBookingTimeConflict = errors.New("booking time conflict")
	ErrInvalidBookingTime  = errors.New("invalid booking time")
	ErrInvalidDuration     = errors.New("invalid booking duration")
	ErrInvalidStatus       = errors.New("invalid booking status")

	ErrPlaceNotFound = errors.New("place not found")

	ErrCoworkingNotFound = errors.New("coworking not found")
	ErrLayoutNotFound    = errors.New("layout version not found")
)

func MapPgError(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	switch pgErr.Code {

	// ============================
	// UNIQUE
	// ============================
	case pgerrcode.UniqueViolation:
		return ErrAlreadyExists

	// ============================
	// FK
	// ============================
	case pgerrcode.ForeignKeyViolation:
		switch pgErr.ConstraintName {
		case "booking_place_id_fkey":
			return ErrPlaceNotFound
		case "booking_status_id_fkey":
			return ErrInvalidStatus
		default:
			return err
		}

	// ============================
	// CHECK
	// ============================
	case pgerrcode.CheckViolation:
		switch pgErr.ConstraintName {
		case "chk_time_order":
			return ErrInvalidBookingTime
		case "chk_duration_hours":
			return ErrInvalidDuration
		default:
			return err
		}

	// ============================
	// EXCLUSION
	// ============================
	case pgerrcode.ExclusionViolation:
		return ErrBookingTimeConflict

	default:
		return err
	}
}
