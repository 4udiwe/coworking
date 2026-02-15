package user_service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/4udiwe/avito-pvz/pkg/transactor"
	"github.com/4udiwe/coworking/auth-service/internal/auth"
	"github.com/4udiwe/coworking/auth-service/internal/entity"
	auth_repository "github.com/4udiwe/coworking/auth-service/internal/repository/auth"
	user_repository "github.com/4udiwe/coworking/auth-service/internal/repository/user"
	"github.com/sirupsen/logrus"
)

type Service struct {
	userRepo UserRepository
	authRepo AuthRepository
	tx       transactor.Transactor
	auth     Auth
	hasher   Hasher
}

func New(userRepo UserRepository, authRepo AuthRepository, tx transactor.Transactor, auth Auth, hasher Hasher) *Service {
	return &Service{
		userRepo: userRepo,
		authRepo: authRepo,
		tx:       tx,
		auth:     auth,
		hasher:   hasher,
	}
}

func (s *Service) Register(
	ctx context.Context,
	email string,
	password string,
	roleCode string,
) (*auth.Tokens, error) {
	// Input validation
	if email == "" {
		logrus.WithField("field", "email").Warn("Register attempt with empty email")
		return nil, ErrEmptyEmail
	}
	if password == "" {
		logrus.WithField("field", "password").Warn("Register attempt with empty password")
		return nil, ErrEmptyPassword
	}
	if roleCode == "" {
		logrus.WithField("field", "roleCode").Warn("Register attempt with empty role code")
		return nil, ErrEmptyRoleCode
	}

	logrus.WithFields(logrus.Fields{
		"email": email,
		"role":  roleCode,
	}).Info("Register started")

	var tokens *auth.Tokens

	err := s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		// Hash password
		passwordHash, err := s.hasher.HashPassword(password)
		if err != nil {
			logrus.WithError(err).Error("Password hashing failed")
			return fmt.Errorf("%w: %v", ErrPasswordHashingFailed, err)
		}

		// Create user
		user := entity.User{
			Email:        email,
			PasswordHash: passwordHash,
			IsActive:     true,
		}

		user, err = s.userRepo.Create(ctx, user)
		if err != nil {
			if errors.Is(err, user_repository.ErrUserAlreadyExists) {
				logrus.WithField("email", email).Warn("User already exists")
				return ErrUserAlreadyExists
			}
			logrus.WithError(err).WithField("email", email).Error("Failed to create user")
			return fmt.Errorf("create user: %w", err)
		}

		// Attach role
		if err := s.userRepo.AttachRole(ctx, user.ID, roleCode); err != nil {
			if errors.Is(err, user_repository.ErrRoleNotFound) {
				logrus.WithField("roleCode", roleCode).Error("Role not found")
				return ErrRoleNotFound
			}
			logrus.WithError(err).WithFields(logrus.Fields{
				"userID":   user.ID,
				"roleCode": roleCode,
			}).Error("Failed to attach role")
			return fmt.Errorf("attach role: %w", err)
		}

		// Generate tokens
		tokens, err = s.auth.GenerateTokens(user)
		if err != nil {
			logrus.WithError(err).WithField("userID", user.ID).Error("Token generation failed")
			return fmt.Errorf("%w: %v", ErrTokenGenerationFailed, err)
		}

		// Save refresh token
		if err := s.authRepo.SaveRefreshToken(
			ctx,
			user.ID,
			s.auth.HashToken(tokens.RefreshToken),
			time.Now().Add(7*24*time.Hour),
		); err != nil {
			logrus.WithError(err).WithField("userID", user.ID).Error("Failed to save refresh token")
			return fmt.Errorf("%w: %v", ErrCannotSaveRefreshToken, err)
		}

		return nil
	})

	if err != nil {
		// Don't wrap errors that are already service errors
		if errors.Is(err, ErrUserAlreadyExists) ||
			errors.Is(err, ErrRoleNotFound) ||
			errors.Is(err, ErrEmptyEmail) ||
			errors.Is(err, ErrEmptyPassword) ||
			errors.Is(err, ErrEmptyRoleCode) {
			return nil, err
		}
		logrus.WithError(err).WithField("email", email).Error("Registration failed")
		return nil, ErrCannotRegisterUser
	}

	logrus.WithField("email", email).Info("Register completed")
	return tokens, nil
}

func (s *Service) Login(
	ctx context.Context,
	email string,
	password string,
) (*auth.Tokens, error) {
	// Input validation
	if email == "" {
		logrus.WithField("field", "email").Warn("Login attempt with empty email")
		return nil, ErrEmptyEmail
	}
	if password == "" {
		logrus.WithField("field", "password").Warn("Login attempt with empty password")
		return nil, ErrEmptyPassword
	}

	logrus.WithField("email", email).Info("Login attempt")

	var tokens *auth.Tokens

	err := s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		// Get user by email
		user, err := s.userRepo.GetByEmail(ctx, email)
		if err != nil {
			if errors.Is(err, user_repository.ErrUserNotFound) {
				logrus.WithField("email", email).Debug("User not found")
				return ErrUserNotFound
			}
			logrus.WithError(err).WithField("email", email).Error("Failed to retrieve user")
			return fmt.Errorf("get user by email: %w", err)
		}

		// Check if user is active
		if !user.IsActive {
			logrus.WithField("email", email).Warn("Login attempt with inactive user")
			return ErrInvalidCredentials
		}

		// Verify password
		if !s.hasher.CheckPasswordHash(password, user.PasswordHash) {
			logrus.WithField("email", email).Debug("Invalid password")
			return ErrInvalidCredentials
		}

		// Generate tokens
		tokens, err = s.auth.GenerateTokens(user)
		if err != nil {
			logrus.WithError(err).WithField("userID", user.ID).Error("Token generation failed")
			return fmt.Errorf("%w: %v", ErrTokenGenerationFailed, err)
		}

		// Save refresh token
		if err := s.authRepo.SaveRefreshToken(
			ctx,
			user.ID,
			s.auth.HashToken(tokens.RefreshToken),
			time.Now().Add(7*24*time.Hour),
		); err != nil {
			logrus.WithError(err).WithField("userID", user.ID).Error("Failed to save refresh token")
			return fmt.Errorf("%w: %v", ErrCannotSaveRefreshToken, err)
		}

		return nil
	})

	if err != nil {
		// Don't wrap validation and auth errors
		if errors.Is(err, ErrInvalidCredentials) ||
			errors.Is(err, ErrUserNotFound) ||
			errors.Is(err, ErrEmptyEmail) ||
			errors.Is(err, ErrEmptyPassword) {
			return nil, err
		}
		logrus.WithError(err).WithField("email", email).Error("Login failed")
		return nil, ErrInvalidCredentials
	}

	logrus.WithField("email", email).Info("Login successful")
	return tokens, nil
}

func (s *Service) Refresh(
	ctx context.Context,
	refreshToken string,
) (*auth.Tokens, error) {
	// Input validation
	if refreshToken == "" {
		logrus.WithField("field", "refreshToken").Warn("Refresh attempt with empty token")
		return nil, ErrEmptyToken
	}

	logrus.Info("Refresh started")

	// Validate token format
	email, err := s.auth.ValidateRefreshToken(refreshToken)
	if err != nil {
		logrus.WithError(err).Debug("Invalid refresh token format")
		return nil, ErrInvalidRefreshToken
	}

	var tokens *auth.Tokens

	err = s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		tokenHash := s.auth.HashToken(refreshToken)

		// Get user ID by refresh token
		userID, err := s.authRepo.GetUserByRefreshToken(ctx, tokenHash)
		if err != nil {
			if errors.Is(err, auth_repository.ErrInvalidRefreshToken) {
				logrus.WithField("email", email).Debug("Refresh token not found or expired")
				return ErrInvalidRefreshToken
			}
			logrus.WithError(err).Error("Failed to get user by refresh token")
			return fmt.Errorf("get user by refresh token: %w", err)
		}

		// Get user by ID
		user, err := s.userRepo.GetByID(ctx, userID)
		if err != nil {
			if errors.Is(err, user_repository.ErrUserNotFound) {
				logrus.WithField("userID", userID).Error("User not found during refresh")
				return ErrInvalidRefreshToken
			}
			logrus.WithError(err).WithField("userID", userID).Error("Failed to retrieve user")
			return fmt.Errorf("get user by id: %w", err)
		}

		// Check if user is active
		if !user.IsActive {
			logrus.WithField("userID", userID).Warn("Refresh attempt with inactive user")
			return ErrInvalidRefreshToken
		}

		// Revoke old refresh token
		if err := s.authRepo.RevokeRefreshToken(ctx, tokenHash); err != nil {
			logrus.WithError(err).WithField("userID", userID).Error("Failed to revoke old refresh token")
			return fmt.Errorf("%w: %v", ErrCannotRevokeRefreshToken, err)
		}

		// Generate new tokens
		tokens, err = s.auth.GenerateTokens(user)
		if err != nil {
			logrus.WithError(err).WithField("userID", userID).Error("Token generation failed")
			return fmt.Errorf("%w: %v", ErrTokenGenerationFailed, err)
		}

		// Save new refresh token
		if err := s.authRepo.SaveRefreshToken(
			ctx,
			user.ID,
			s.auth.HashToken(tokens.RefreshToken),
			time.Now().Add(7*24*time.Hour),
		); err != nil {
			logrus.WithError(err).WithField("userID", userID).Error("Failed to save refresh token")
			return fmt.Errorf("%w: %v", ErrCannotSaveRefreshToken, err)
		}

		return nil
	})

	if err != nil {
		// Don't wrap validation errors
		if errors.Is(err, ErrInvalidRefreshToken) ||
			errors.Is(err, ErrEmptyToken) {
			return nil, err
		}
		logrus.WithError(err).Error("Refresh failed")
		return nil, ErrInvalidRefreshToken
	}

	logrus.WithField("email", email).Info("Refresh successful")
	return tokens, nil
}

func (s *Service) Logout(
	ctx context.Context,
	refreshToken string,
) error {
	// Input validation
	if refreshToken == "" {
		logrus.WithField("field", "refreshToken").Warn("Logout attempt with empty token")
		return ErrEmptyToken
	}

	logrus.Info("Logout started")

	tokenHash := s.auth.HashToken(refreshToken)

	// Revoke refresh token
	if err := s.authRepo.RevokeRefreshToken(ctx, tokenHash); err != nil {
		if errors.Is(err, auth_repository.ErrInvalidRefreshToken) {
			logrus.WithField("token_hash", tokenHash).Debug("Token already revoked or not found")
			return ErrInvalidRefreshToken
		}
		logrus.WithError(err).WithField("token_hash", tokenHash).Error("Failed to revoke refresh token")
		return fmt.Errorf("%w: %v", ErrCannotRevokeRefreshToken, err)
	}

	logrus.WithField("token_hash", tokenHash).Info("Refresh token revoked successfully")
	return nil
}
