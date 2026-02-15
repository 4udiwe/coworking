package user_service

import "errors"

var (
	// Authentication errors
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")

	// User errors
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	// Input validation errors
	ErrEmptyEmail    = errors.New("email cannot be empty")
	ErrEmptyPassword = errors.New("password cannot be empty")
	ErrEmptyRoleCode = errors.New("role code cannot be empty")
	ErrEmptyToken    = errors.New("token cannot be empty")

	// Service errors
	ErrCannotRegisterUser        = errors.New("cannot register user")
	ErrPasswordHashingFailed     = errors.New("password hashing failed")
	ErrTokenGenerationFailed     = errors.New("token generation failed")
	ErrRoleNotFound             = errors.New("role not found")
	ErrCannotSaveRefreshToken   = errors.New("cannot save refresh token")
	ErrCannotRevokeRefreshToken = errors.New("cannot revoke refresh token")
)
