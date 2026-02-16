package user_service

import "errors"

var (
	// Authentication errors
	ErrInvalidCredentials        = errors.New("invalid email or password")
	ErrInvalidRefreshToken       = errors.New("invalid refresh token")
	ErrInvalidRefreshTokenFormat = errors.New("invalid refresh token format")

	// User errors
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserInactive      = errors.New("user is inactive")

	// Session errors
	ErrSessionNotFound     = errors.New("session not found")
	ErrSessionExpired      = errors.New("session is expired")
	ErrCannotUpdateSession = errors.New("cannot update session")
	ErrCannotRevokeSession = errors.New("cannot revoke session")

	// Input validation errors
	ErrEmptyEmail    = errors.New("email cannot be empty")
	ErrEmptyPassword = errors.New("password cannot be empty")
	ErrEmptyRoleCode = errors.New("role code cannot be empty")
	ErrEmptyToken    = errors.New("token cannot be empty")

	// Service errors
	ErrCannotRegisterUser       = errors.New("cannot register user")
	ErrPasswordHashingFailed    = errors.New("password hashing failed")
	ErrTokenGenerationFailed    = errors.New("token generation failed")
	ErrRoleNotFound             = errors.New("role not found")
	ErrCannotSaveRefreshToken   = errors.New("cannot save refresh token")
	ErrCannotRevokeRefreshToken = errors.New("cannot revoke refresh token")
	ErrCannotGenerateTokens     = errors.New("cannot generate tokens")
	ErrCannotFetchSessions      = errors.New("cannot fetch sessions")
)
