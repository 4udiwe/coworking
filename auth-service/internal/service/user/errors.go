package user_service

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrRoleNotFound     = errors.New("role not found")
	ErrEmptyRoles       = errors.New("roles cannot be empty")
	ErrEmptyUserID      = errors.New("empty user id")
	ErrCannotFetchUsers = errors.New("cannot fetch users")
)
