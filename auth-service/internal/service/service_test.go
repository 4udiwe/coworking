package user_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/4udiwe/coworking/auth-service/internal/auth"
	"github.com/4udiwe/coworking/auth-service/internal/entity"
	service "github.com/4udiwe/coworking/auth-service/internal/service"

	mock_tx "github.com/4udiwe/coworking/auth-service/internal/mocks"
	m "github.com/4udiwe/coworking/auth-service/internal/service/mocks"
)

func TestService_Register(t *testing.T) {
	type mocks struct {
		ur *m.MockUserRepository
		ar *m.MockAuthRepository
		tx *mock_tx.MockTransactor
		a  *m.MockAuth
		h  *m.MockHasher
	}

	tests := []struct {
		name         string
		mockBehavior func(m mocks)
		expectedErr  error
	}{
		{
			name: "success",
			mockBehavior: func(m mocks) {
				userID := uuid.New()

				m.tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				m.h.EXPECT().HashPassword("pass").Return("hash", nil)

				m.ur.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(entity.User{ID: userID, Email: "mail"}, nil)

				m.ur.EXPECT().AttachRole(gomock.Any(), userID, "student").Return(nil)

				m.a.EXPECT().
					GenerateTokens(gomock.Any(), gomock.Any()).
					Return(&auth.Tokens{RefreshToken: "rt"}, nil)

				m.a.EXPECT().HashToken("rt").Return("hashRT")

				m.ar.EXPECT().
					CreateSession(gomock.Any(), gomock.Any(), "hashRT").
					Return(nil)
			},
		},
		{
			name: "hash password fail",
			mockBehavior: func(m mocks) {
				m.tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})
				m.h.EXPECT().HashPassword("pass").Return("", errors.New("fail"))
			},
			expectedErr: errors.New("fail"),
		},
		{
			name: "create user fail",
			mockBehavior: func(m mocks) {
				m.tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})
				m.h.EXPECT().HashPassword("pass").Return("hash", nil)
				m.ur.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(entity.User{}, errors.New("fail"))
			},
			expectedErr: errors.New("fail"),
		},
		{
			name: "attach role fail",
			mockBehavior: func(m mocks) {
				userID := uuid.New()
				m.tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})
				m.h.EXPECT().HashPassword("pass").Return("hash", nil)
				m.ur.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(entity.User{ID: userID}, nil)
				m.ur.EXPECT().AttachRole(gomock.Any(), userID, "student").Return(errors.New("fail"))
			},
			expectedErr: errors.New("fail"),
		},
		{
			name: "transaction fail",
			mockBehavior: func(m mocks) {
				m.tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).
					Return(errors.New("tx fail"))
			},
			expectedErr: errors.New("tx fail"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks{
				ur: m.NewMockUserRepository(ctrl),
				ar: m.NewMockAuthRepository(ctrl),
				tx: mock_tx.NewMockTransactor(ctrl),
				a:  m.NewMockAuth(ctrl),
				h:  m.NewMockHasher(ctrl),
			}

			s := service.New(m.ur, m.ar, m.tx, m.a, m.h, 7*24*time.Hour)

			if tt.mockBehavior != nil {
				tt.mockBehavior(m)
			}

			_, err := s.Register(context.Background(), "mail", "pass", "student", "ua", "ip")

			if tt.expectedErr != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_Login(t *testing.T) {
	type mocks struct {
		ur *m.MockUserRepository
		ar *m.MockAuthRepository
		tx *mock_tx.MockTransactor
		a  *m.MockAuth
		h  *m.MockHasher
	}

	userID := uuid.New()

	user := entity.User{
		ID:           userID,
		Email:        "mail",
		PasswordHash: "hash",
	}

	tests := []struct {
		name         string
		mockBehavior func(m mocks)
		expectedErr  error
	}{
		{
			name: "success",
			mockBehavior: func(m mocks) {

				m.tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				m.ur.EXPECT().
					GetByEmail(gomock.Any(), "mail").
					Return(user, nil)

				m.h.EXPECT().
					CheckPasswordHash("pass", "hash").
					Return(true)

				m.a.EXPECT().
					GenerateTokens(user, gomock.Any()).
					Return(&auth.Tokens{RefreshToken: "rt"}, nil)

				m.a.EXPECT().
					HashToken("rt").
					Return("hashRT")

				m.ar.EXPECT().
					CreateSession(gomock.Any(), gomock.Any(), "hashRT").
					Return(nil)
			},
		},
		{
			name: "get user fail",
			mockBehavior: func(m mocks) {

				m.tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				m.ur.EXPECT().
					GetByEmail(gomock.Any(), "mail").
					Return(entity.User{}, errors.New("fail"))
			},
			expectedErr: service.ErrInvalidCredentials,
		},
		{
			name: "wrong password",
			mockBehavior: func(m mocks) {

				m.tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				m.ur.EXPECT().
					GetByEmail(gomock.Any(), "mail").
					Return(user, nil)

				m.h.EXPECT().
					CheckPasswordHash("pass", "hash").
					Return(false)
			},
			expectedErr: service.ErrInvalidCredentials,
		},
		{
			name: "generate tokens fail",
			mockBehavior: func(m mocks) {

				m.tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				m.ur.EXPECT().
					GetByEmail(gomock.Any(), "mail").
					Return(user, nil)

				m.h.EXPECT().
					CheckPasswordHash("pass", "hash").
					Return(true)

				m.a.EXPECT().
					GenerateTokens(user, gomock.Any()).
					Return(nil, errors.New("fail"))
			},
			expectedErr: service.ErrInvalidCredentials,
		},
		{
			name: "create session fail",
			mockBehavior: func(m mocks) {

				m.tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				m.ur.EXPECT().
					GetByEmail(gomock.Any(), "mail").
					Return(user, nil)

				m.h.EXPECT().
					CheckPasswordHash("pass", "hash").
					Return(true)

				m.a.EXPECT().
					GenerateTokens(user, gomock.Any()).
					Return(&auth.Tokens{RefreshToken: "rt"}, nil)

				m.a.EXPECT().
					HashToken("rt").
					Return("hashRT")

				m.ar.EXPECT().
					CreateSession(gomock.Any(), gomock.Any(), "hashRT").
					Return(errors.New("fail"))
			},
			expectedErr: service.ErrInvalidCredentials,
		},
		{
			name: "transaction fail",
			mockBehavior: func(m mocks) {
				m.tx.EXPECT().
					WithinTransaction(gomock.Any(), gomock.Any()).
					Return(errors.New("tx fail"))
			},
			expectedErr: service.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks{
				ur: m.NewMockUserRepository(ctrl),
				ar: m.NewMockAuthRepository(ctrl),
				tx: mock_tx.NewMockTransactor(ctrl),
				a:  m.NewMockAuth(ctrl),
				h:  m.NewMockHasher(ctrl),
			}

			s := service.New(m.ur, m.ar, m.tx, m.a, m.h, 7*24*time.Hour)

			if tt.mockBehavior != nil {
				tt.mockBehavior(m)
			}

			_, err := s.Login(context.Background(), "mail", "pass", "ua", "ip")

			if tt.expectedErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_Refresh(t *testing.T) {
	type mocks struct {
		ur *m.MockUserRepository
		ar *m.MockAuthRepository
		tx *mock_tx.MockTransactor
		a  *m.MockAuth
		h  *m.MockHasher
	}

	tests := []struct {
		name         string
		mockBehavior func(m mocks, sessionID, userID uuid.UUID, user entity.User, validSession entity.Session)
		expectedErr  error
	}{
		{
			name: "success",
			mockBehavior: func(m mocks, sessionID, userID uuid.UUID, user entity.User, validSession entity.Session) {
				m.a.EXPECT().ParseRefreshToken("rt").Return(&auth.RefreshClaims{SessionID: sessionID}, nil)

				m.tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					},
				)

				m.ar.EXPECT().GetSessionByID(gomock.Any(), sessionID).Return(validSession, nil)
				m.ur.EXPECT().GetByID(gomock.Any(), userID).Return(user, nil)
				m.ar.EXPECT().UpdateLastUsedAt(gomock.Any(), sessionID).Return(nil)
				m.ar.EXPECT().RevokeSession(gomock.Any(), sessionID).Return(nil)

				m.a.EXPECT().GenerateTokens(user, gomock.Any()).Return(&auth.Tokens{RefreshToken: "newRT"}, nil)
				m.a.EXPECT().HashToken("newRT").Return("hashNew")
				m.ar.EXPECT().CreateSession(gomock.Any(), gomock.Any(), "hashNew").Return(nil)
			},
		},
		{
			name: "parse token fail",
			mockBehavior: func(m mocks, sessionID, userID uuid.UUID, user entity.User, validSession entity.Session) {
				m.a.EXPECT().ParseRefreshToken("rt").Return(nil, errors.New("fail"))
			},
			expectedErr: service.ErrInvalidRefreshTokenFormat,
		},
		{
			name: "session not found",
			mockBehavior: func(m mocks, sessionID, userID uuid.UUID, user entity.User, validSession entity.Session) {
				m.a.EXPECT().ParseRefreshToken("rt").Return(&auth.RefreshClaims{SessionID: sessionID}, nil)
				m.tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) },
				)
				m.ar.EXPECT().GetSessionByID(gomock.Any(), sessionID).Return(entity.Session{}, errors.New("not found"))
			},
			expectedErr: service.ErrSessionNotFound,
		},
		{
			name: "session revoked",
			mockBehavior: func(m mocks, sessionID, userID uuid.UUID, user entity.User, validSession entity.Session) {
				m.a.EXPECT().ParseRefreshToken("rt").Return(&auth.RefreshClaims{SessionID: sessionID}, nil)
				m.tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) },
				)
				m.ar.EXPECT().GetSessionByID(gomock.Any(), sessionID).Return(entity.Session{
					ID:        sessionID,
					UserID:    userID,
					ExpiresAt: time.Now().Add(time.Hour),
					Revoked:   true,
				}, nil)
			},
			expectedErr: service.ErrSessionExpired,
		},
		{
			name: "session expired",
			mockBehavior: func(m mocks, sessionID, userID uuid.UUID, user entity.User, validSession entity.Session) {
				m.a.EXPECT().ParseRefreshToken("rt").Return(&auth.RefreshClaims{SessionID: sessionID}, nil)
				m.tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) },
				)
				m.ar.EXPECT().GetSessionByID(gomock.Any(), sessionID).Return(entity.Session{
					ID:        sessionID,
					UserID:    userID,
					ExpiresAt: time.Now().Add(-time.Hour),
					Revoked:   false,
				}, nil)
			},
			expectedErr: service.ErrSessionExpired,
		},
		{
			name: "user not found",
			mockBehavior: func(m mocks, sessionID, userID uuid.UUID, user entity.User, validSession entity.Session) {
				m.a.EXPECT().ParseRefreshToken("rt").Return(&auth.RefreshClaims{SessionID: sessionID}, nil)
				m.tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) },
				)
				m.ar.EXPECT().GetSessionByID(gomock.Any(), sessionID).Return(validSession, nil)
				m.ur.EXPECT().GetByID(gomock.Any(), userID).Return(entity.User{}, errors.New("not found"))
			},
			expectedErr: service.ErrUserNotFound,
		},
		{
			name: "user inactive",
			mockBehavior: func(m mocks, sessionID, userID uuid.UUID, user entity.User, validSession entity.Session) {
				m.a.EXPECT().ParseRefreshToken("rt").Return(&auth.RefreshClaims{SessionID: sessionID}, nil)
				m.tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) },
				)
				m.ar.EXPECT().GetSessionByID(gomock.Any(), sessionID).Return(validSession, nil)
				m.ur.EXPECT().GetByID(gomock.Any(), userID).Return(entity.User{ID: userID, IsActive: false}, nil)
			},
			expectedErr: service.ErrUserInactive,
		},
		{
			name: "update last used fail",
			mockBehavior: func(m mocks, sessionID, userID uuid.UUID, user entity.User, validSession entity.Session) {
				m.a.EXPECT().ParseRefreshToken("rt").Return(&auth.RefreshClaims{SessionID: sessionID}, nil)
				m.tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) },
				)
				m.ar.EXPECT().GetSessionByID(gomock.Any(), sessionID).Return(validSession, nil)
				m.ur.EXPECT().GetByID(gomock.Any(), userID).Return(user, nil)
				m.ar.EXPECT().UpdateLastUsedAt(gomock.Any(), sessionID).Return(errors.New("fail"))
			},
			expectedErr: service.ErrCannotUpdateSession,
		},
		{
			name: "revoke session fail",
			mockBehavior: func(m mocks, sessionID, userID uuid.UUID, user entity.User, validSession entity.Session) {
				m.a.EXPECT().ParseRefreshToken("rt").Return(&auth.RefreshClaims{SessionID: sessionID}, nil)
				m.tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) },
				)
				m.ar.EXPECT().GetSessionByID(gomock.Any(), sessionID).Return(validSession, nil)
				m.ur.EXPECT().GetByID(gomock.Any(), userID).Return(user, nil)
				m.ar.EXPECT().UpdateLastUsedAt(gomock.Any(), sessionID).Return(nil)
				m.ar.EXPECT().RevokeSession(gomock.Any(), sessionID).Return(errors.New("fail"))
			},
			expectedErr: service.ErrCannotRevokeSession,
		},
		{
			name: "generate tokens fail",
			mockBehavior: func(m mocks, sessionID, userID uuid.UUID, user entity.User, validSession entity.Session) {
				m.a.EXPECT().ParseRefreshToken("rt").Return(&auth.RefreshClaims{SessionID: sessionID}, nil)
				m.tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) },
				)
				m.ar.EXPECT().GetSessionByID(gomock.Any(), sessionID).Return(validSession, nil)
				m.ur.EXPECT().GetByID(gomock.Any(), userID).Return(user, nil)
				m.ar.EXPECT().UpdateLastUsedAt(gomock.Any(), sessionID).Return(nil)
				m.ar.EXPECT().RevokeSession(gomock.Any(), sessionID).Return(nil)
				m.a.EXPECT().GenerateTokens(user, gomock.Any()).Return(nil, errors.New("fail"))
			},
			expectedErr: service.ErrCannotGenerateTokens,
		},
		{
			name: "transaction fail",
			mockBehavior: func(m mocks, sessionID, userID uuid.UUID, user entity.User, validSession entity.Session) {
				m.a.EXPECT().ParseRefreshToken("rt").Return(&auth.RefreshClaims{SessionID: sessionID}, nil)
				m.tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).Return(errors.New("tx fail"))
			},
			expectedErr: service.ErrInvalidRefreshToken,
		},
		{
			name: "create session fail",
			mockBehavior: func(m mocks, sessionID, userID uuid.UUID, user entity.User, validSession entity.Session) {
				m.a.EXPECT().ParseRefreshToken("rt").Return(&auth.RefreshClaims{SessionID: sessionID}, nil)
				m.tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					},
				)
				m.ar.EXPECT().GetSessionByID(gomock.Any(), sessionID).Return(validSession, nil)
				m.ur.EXPECT().GetByID(gomock.Any(), userID).Return(user, nil)
				m.ar.EXPECT().UpdateLastUsedAt(gomock.Any(), sessionID).Return(nil)
				m.ar.EXPECT().RevokeSession(gomock.Any(), sessionID).Return(nil)
				m.a.EXPECT().GenerateTokens(user, gomock.Any()).Return(&auth.Tokens{RefreshToken: "newRT"}, nil)
				m.a.EXPECT().HashToken("newRT").Return("hashNew")
				m.ar.EXPECT().CreateSession(gomock.Any(), gomock.Any(), "hashNew").Return(errors.New("fail"))
			},
			expectedErr: service.ErrInvalidRefreshToken,
		},
	}	

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks{
				ur: m.NewMockUserRepository(ctrl),
				ar: m.NewMockAuthRepository(ctrl),
				tx: mock_tx.NewMockTransactor(ctrl),
				a:  m.NewMockAuth(ctrl),
				h:  m.NewMockHasher(ctrl),
			}

			sessionID := uuid.New()
			userID := uuid.New()
			user := entity.User{ID: userID, IsActive: true}
			validSession := entity.Session{
				ID:        sessionID,
				UserID:    userID,
				ExpiresAt: time.Now().Add(time.Hour),
				Revoked:   false,
			}

			if tt.mockBehavior != nil {
				tt.mockBehavior(m, sessionID, userID, user, validSession)
			}

			s := service.New(m.ur, m.ar, m.tx, m.a, m.h, 7*24*time.Hour)
			_, err := s.Refresh(context.Background(), "rt", "ua", "ip")

			if tt.expectedErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}


func TestService_Logout(t *testing.T) {
	type mocks struct {
		ar *m.MockAuthRepository
		a  *m.MockAuth
		tx *mock_tx.MockTransactor
	}

	tests := []struct {
		name         string
		mockBehavior func(m mocks)
		expectedErr  error
	}{
		{
			name: "success",
			mockBehavior: func(m mocks) {
				m.a.EXPECT().ParseRefreshToken("rt").Return(&auth.RefreshClaims{SessionID: uuid.New()}, nil)
				m.ar.EXPECT().RevokeSession(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "parse token fail",
			mockBehavior: func(m mocks) {
				m.a.EXPECT().ParseRefreshToken("rt").Return(nil, errors.New("fail"))
			},
			expectedErr: service.ErrInvalidRefreshToken,
		},
		{
			name: "revoke session fail",
			mockBehavior: func(m mocks) {
				m.a.EXPECT().ParseRefreshToken("rt").Return(&auth.RefreshClaims{SessionID: uuid.New()}, nil)
				m.ar.EXPECT().RevokeSession(gomock.Any(), gomock.Any()).Return(errors.New("fail"))
			},
			expectedErr: errors.New("fail"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks{
				ar: m.NewMockAuthRepository(ctrl),
				a:  m.NewMockAuth(ctrl),
				tx: mock_tx.NewMockTransactor(ctrl),
			}

			s := service.New(nil, m.ar, m.tx, m.a, nil, 7*24*time.Hour)

			if tt.mockBehavior != nil {
				tt.mockBehavior(m)
			}

			err := s.Logout(context.Background(), "rt")

			if tt.expectedErr != nil {
				require.Error(t, err)
				if err != nil && err.Error() == tt.expectedErr.Error() {
					// ok
				} else {
					require.ErrorIs(t, err, tt.expectedErr)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_Register_Validation(t *testing.T) {
	s := service.New(nil, nil, nil, nil, nil, 7*24*time.Hour)
	_, err := s.Register(context.Background(), "", "pass", "student", "ua", "ip")
	require.ErrorIs(t, err, service.ErrEmptyEmail)
	_, err = s.Register(context.Background(), "mail", "", "student", "ua", "ip")
	require.ErrorIs(t, err, service.ErrEmptyPassword)
	_, err = s.Register(context.Background(), "mail", "pass", "", "ua", "ip")
	require.ErrorIs(t, err, service.ErrEmptyRoleCode)
}

func TestService_Refresh_EmptyToken(t *testing.T) {
	s := service.New(nil, nil, nil, nil, nil, 7*24*time.Hour)
	_, err := s.Refresh(context.Background(), "", "ua", "ip")
	require.ErrorIs(t, err, service.ErrEmptyToken)
}
