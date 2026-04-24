package entity

import (
	"time"

	"github.com/google/uuid"
)

type RoleCode string

const (
	RoleStudent RoleCode = "student"
	RoleTeacher RoleCode = "teacher"
	RoleAdmin   RoleCode = "admin"
)

type Role struct {
	ID        uuid.UUID
	Code      RoleCode
	Name      string
	CreatedAt time.Time
}
