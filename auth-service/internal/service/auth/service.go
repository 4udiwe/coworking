package auth_service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/4udiwe/avito-pvz/pkg/transactor"
	"github.com/4udiwe/coworking/auth-service/internal/auth"
	"github.com/4udiwe/coworking/auth-service/internal/entity"
	user_repository "github.com/4udiwe/coworking/auth-service/internal/repository/user"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// MAX_ACTIVE_SESSIONS_PER_USER defines the maximum number of concurrent active sessions
// for a single user at any given time.
//
// WHY THIS LIMIT:
// - Prevents unbounded session accumulation
// - Typical user has 1-2 devices (mobile, web)
// - Value of 5 accommodates power users (phone, tablet, laptop, desktop, smart device)
// - Without limit: 1000 users × 365 days = 365,000 sessions in DB
// - With limit: 1000 users × max(5 active) = ~5,000 concurrent sessions
//
// WHAT HAPPENS:
// When user logs in on 6th device while already having 5 sessions:
// 1. System detects limit exceeded
// 2. Automatically revokes the oldest session
// 3. New login succeeds, user has 5 sessions again
// 4. Event logged for monitoring
const MAX_ACTIVE_SESSIONS_PER_USER = 5

type Service struct {
	userRepo UserRepository
	authRepo AuthRepository
	tx       transactor.Transactor
	auth     Auth
	hasher   Hasher

	refreshTokenTTL time.Duration
}

func New(
	userRepo UserRepository,
	authRepo AuthRepository,
	tx transactor.Transactor,
	auth Auth,
	hasher Hasher,
	refreshTokenTTL time.Duration,
) *Service {
	return &Service{
		userRepo:        userRepo,
		authRepo:        authRepo,
		tx:              tx,
		auth:            auth,
		hasher:          hasher,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *Service) Register(
	ctx context.Context,
	email string,
	password string,
	firstName string,
	lastName string,
	roleCode string,
	userAgent string,
	deviceInfo string,
	ip string,
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
			FirstName:    firstName,
			LastName:     lastName,
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

		sessionID := uuid.New()

		// Generate tokens
		tokens, err = s.auth.GenerateTokens(user, sessionID)
		if err != nil {
			logrus.WithError(err).WithField("userID", user.ID).Error("Token generation failed")
			return fmt.Errorf("%w: %v", ErrTokenGenerationFailed, err)
		}

		// Save refresh token
		if err := s.authRepo.CreateSession(
			ctx,
			entity.Session{
				ID:         sessionID,
				UserID:     user.ID,
				UserAgent:  userAgent,
				IPAddress:  ip,
				DeviceName: &deviceInfo,
				ExpiresAt:  time.Now().Add(s.refreshTokenTTL),
			},
			s.auth.HashToken(tokens.RefreshToken),
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
	userAgent string,
	deviceInfo string,
	ip string,
) (*auth.Tokens, error) {

	if email == "" {
		logrus.WithField("field", "email").Warn("Login attempt with empty email")
		return nil, ErrEmptyEmail
	}
	if password == "" {
		logrus.WithField("field", "password").Warn("Login attempt with empty password")
		return nil, ErrEmptyPassword
	}

	var tokens *auth.Tokens

	err := s.tx.WithinTransaction(ctx, func(ctx context.Context) error {

		user, err := s.userRepo.GetByEmail(ctx, email)
		if err != nil {
			return ErrUserNotFound
		}

		if !s.hasher.CheckPasswordHash(password, user.PasswordHash) {
			return ErrInvalidCredentials
		}

		// ============================================
		// NEW: Enforce session limit per user (MAX_ACTIVE_SESSIONS_PER_USER)
		// ============================================
		// WHAT THIS DOES:
		// 1. Get all ACTIVE (non-revoked) sessions for this user
		// 2. If user already has 5+ sessions, revoke the oldest one
		// 3. This allows new login to proceed
		//
		// WHY:
		// - Prevents unlimited session accumulation
		// - User with 8 devices gets limited to 5 active sessions
		// - When login on 6th device → oldest (1st) auto-revoked
		// - User maintains their 5 most recent/active sessions
		// ============================================
		activeSessions, err := s.authRepo.GetUserSessions(ctx, user.ID, true)
		if err != nil {
			// If we can't fetch sessions, log but allow login anyway
			// (better to have login succeed than fail due to DB issue)
			logrus.WithError(err).WithField("user_id", user.ID).
				Warn("Failed to check session limit during login, proceeding anyway")
		} else if len(activeSessions) >= MAX_ACTIVE_SESSIONS_PER_USER {
			// User has reached session limit
			// Revoke the oldest session to make room for new one
			if err := s.authRepo.DeleteOldestSessionByUser(ctx, user.ID); err != nil {
				logrus.WithError(err).WithField("user_id", user.ID).
					Warn("Failed to revoke oldest session")
				// Don't fail login - new session creation may still succeed
			} else {
				// Log successful auto-revocation for monitoring and debugging
				logrus.WithFields(logrus.Fields{
					"user_id":              user.ID,
					"email":                email,
					"active_sessions":      len(activeSessions),
					"max_sessions":         MAX_ACTIVE_SESSIONS_PER_USER,
					"user_agent":           userAgent,
					"device_info":          deviceInfo,
					"action":               "auto_revoke_oldest_session",
				}).Info("User reached session limit, auto-revoked oldest session")
			}
		}

		sessionID := uuid.New()

		tokens, err = s.auth.GenerateTokens(user, sessionID)
		if err != nil {
			return ErrCannotGenerateTokens
		}

		return s.authRepo.CreateSession(
			ctx,
			entity.Session{
				ID:         sessionID,
				UserID:     user.ID,
				UserAgent:  userAgent,
				IPAddress:  ip,
				DeviceName: &deviceInfo,
				ExpiresAt:  time.Now().Add(s.refreshTokenTTL),
			},
			s.auth.HashToken(tokens.RefreshToken),
		)
	})

	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return tokens, nil
}

func (s *Service) Refresh(
	ctx context.Context,
	refreshToken string,
	userAgent string,
	deviceInfo string,
	ip string,
) (*auth.Tokens, error) {
	logrus.WithField("userAgent", userAgent).Info("Refresh token attempt")

	if refreshToken == "" {
		return nil, ErrEmptyToken
	}

	claims, err := s.auth.ParseRefreshToken(refreshToken)
	if err != nil {
		logrus.WithError(err).WithField("userAgent", userAgent).Warn("Invalid refresh token")
		return nil, ErrInvalidRefreshTokenFormat
	}

	var tokens *auth.Tokens

	err = s.tx.WithinTransaction(ctx, func(ctx context.Context) error {

		session, err := s.authRepo.GetSessionByID(ctx, claims.SessionID)
		if err != nil {
			return ErrSessionNotFound
		}

		if session.Revoked || session.ExpiresAt.Before(time.Now()) {
			return ErrSessionExpired
		}

		user, err := s.userRepo.GetByID(ctx, session.UserID)
		if err != nil {
			return ErrUserNotFound
		}
		if !user.IsActive {
			return ErrUserInactive
		}

		// фиксируем использование
		if err := s.authRepo.UpdateLastUsedAt(ctx, session.ID); err != nil {
			return ErrCannotUpdateSession
		}

		// ROTATION
		if err := s.authRepo.RevokeSession(ctx, session.ID); err != nil {
			return ErrCannotRevokeSession
		}

		newSessionID := uuid.New()

		tokens, err = s.auth.GenerateTokens(user, newSessionID)
		if err != nil {
			return ErrCannotGenerateTokens
		}

		return s.authRepo.CreateSession(
			ctx,
			entity.Session{
				ID:         newSessionID,
				UserID:     user.ID,
				UserAgent:  userAgent,
				IPAddress:  ip,
				DeviceName: &deviceInfo,
				ExpiresAt:  time.Now().Add(s.refreshTokenTTL),
			},
			s.auth.HashToken(tokens.RefreshToken),
		)
	})

	if err != nil {
		// Don't wrap service errors
		if errors.Is(err, ErrSessionNotFound) ||
			errors.Is(err, ErrSessionExpired) ||
			errors.Is(err, ErrUserNotFound) ||
			errors.Is(err, ErrCannotUpdateSession) ||
			errors.Is(err, ErrCannotRevokeSession) ||
			errors.Is(err, ErrCannotGenerateTokens) ||
			errors.Is(err, ErrUserInactive) {
			return nil, err
		}
		logrus.WithError(err).WithField("userAgent", userAgent).Warn("Refresh token failed")
		return nil, ErrInvalidRefreshToken
	}

	return tokens, nil
}

func (s *Service) Logout(
	ctx context.Context,
	refreshToken string,
) error {

	claims, err := s.auth.ParseRefreshToken(refreshToken)
	if err != nil {
		return ErrInvalidRefreshToken
	}

	return s.authRepo.RevokeSession(ctx, claims.SessionID)
}

func (s *Service) GetUserSessions(
	ctx context.Context,
	userID uuid.UUID,
	onlyActive bool,
) ([]entity.Session, error) {
	logrus.Info("GetUserSessions attempt")

	if userID == uuid.Nil {
		return nil, ErrEmptyUserID
	}

	logrus.WithField("user_id", userID).Info("Fetching user sessions")

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	sessions, err := s.authRepo.GetUserSessions(ctx, userID, onlyActive)
	if err != nil {
		return nil, ErrCannotFetchSessions
	}

	return sessions, nil
}

func (s *Service) RevokeSession(
	ctx context.Context,
	sessionID uuid.UUID,
) error {
	logrus.WithField("session_id", sessionID).Info("RevokeSession attempt")

	err := s.authRepo.RevokeSession(ctx, sessionID)
	if err != nil {
		logrus.WithError(err).WithField("session_id", sessionID).Error("Failed to revoke session")
		return ErrCannotRevokeSession
	}

	return nil
}

