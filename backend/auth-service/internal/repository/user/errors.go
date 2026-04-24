package user_repository

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrRoleNotFound      = errors.New("role not found")
	ErrUserNotFound      = errors.New("user not found")
)
