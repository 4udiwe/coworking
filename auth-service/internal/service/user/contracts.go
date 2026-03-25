package user_service

import (
	"context"

	"github.com/4udiwe/coworking/auth-service/internal/entity"
	"github.com/google/uuid"
)

type UserRepository interface {
	AttachRole(ctx context.Context, userID uuid.UUID, roleCode string) error
	GetByID(ctx context.Context, userID uuid.UUID) (entity.User, error)
	GetUsers(
		ctx context.Context,
		page, pageSize int,
		searchQuery, filterRole, sortField *string,
		filterIsActive *bool,
	) ([]entity.User, int64, error)
	SetActive(ctx context.Context, userID uuid.UUID, active bool) error
	ClearRoles(ctx context.Context, userID uuid.UUID) error
}
