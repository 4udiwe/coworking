package get_users

import (
	"net/http"
	"strings"

	"github.com/4udiwe/coworking/auth-service/internal/api"
	"github.com/4udiwe/coworking/auth-service/internal/api/dto"
	"github.com/4udiwe/coworking/auth-service/internal/entity"
	"github.com/4udiwe/coworking/auth-service/pkg/decorator"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

const PAGE_NUMBER int = 1
const PAGE_SIZE int = 10

type handler struct {
	s UserService
}

func New(userService UserService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{s: userService})
}

type Request = dto.GetUsersRequest

func (h *handler) Handle(ctx echo.Context, in Request) error {
	in.Search = strings.ToLower(in.Search)
	var pageSize, pageNumber int
	if in.Size != nil {
		pageSize = *in.Size
	} else {
		pageSize = PAGE_SIZE
	}
	if in.Page != nil {
		pageNumber = *in.Page
	} else {
		pageNumber = PAGE_NUMBER
	}
	users, total, err := h.s.GetUsers(ctx.Request().Context(), pageNumber, pageSize, &in.Search, in.Role, in.IsActive, &in.Sort)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, dto.PaginatedUsers{
		Page:  pageNumber,
		Size:  pageSize,
		Total: total,
		Users: lo.Map(users, func(user entity.User, _ int) dto.User {
			return dto.User{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
				IsActive:  user.IsActive,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				Roles: lo.Map(user.Roles, func(r entity.Role, _ int) dto.Role {
					return dto.Role{
						ID:       r.ID,
						RoleCode: string(r.Code),
						Name:     r.Name,
					}
				}),
			}
		}),
	})
}
