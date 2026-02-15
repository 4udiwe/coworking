package auth_repository

import "errors"

var (
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)
