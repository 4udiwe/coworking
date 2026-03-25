package get_users

import (
	"context"

	"github.com/4udiwe/coworking/auth-service/internal/entity"
)

type UserService interface {
	GetUsers(
		ctx context.Context,
		page, pageSize int,
		searchQuery *string,
		filterRole *string,
		filterIsActive *bool,
		sortField *string,
	) (users []entity.User, total int64, err error)
}
