package user_service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/4udiwe/coworking/auth-service/internal/auth"
	"github.com/4udiwe/coworking/auth-service/internal/entity"
	mock_transactor "github.com/4udiwe/coworking/auth-service/internal/mocks"
	auth_repository "github.com/4udiwe/coworking/auth-service/internal/repository/auth"
	user_repository "github.com/4udiwe/coworking/auth-service/internal/repository/user"
	service "github.com/4udiwe/coworking/auth-service/internal/service"
	"github.com/4udiwe/coworking/auth-service/internal/service/mocks"
)

func TestService_Register(t *testing.T) {
	var (
		ctx       = context.Background()
		email     = "test@mail.com"
		password  = "password123"
		roleCode  = "USER"
		userID    = uuid.New()
		hash      = "hashed_password"
		tokenHash = "hashed_refresh"
		tokens    = &auth.Tokens{
			AccessToken:  "access",
			RefreshToken: "refresh",
			ExpiresIn:    900,
		}
		arbitraryErr = errors.New("arbitrary error")
	)

	type MockBehavior func(
		ur *mocks.MockUserRepository,
		ar *mocks.MockAuthRepository,
		a *mocks.MockAuth,
		h *mocks.MockHasher,
		tx *mock_transactor.MockTransactor,
	)

	for _, tc := range []struct {
		name         string
		email        string
		password     string
		roleCode     string
		mockBehavior MockBehavior
		want         *auth.Tokens
		wantError    error
	}{
		{
			name:     "success",
			email:    email,
			password: password,
			roleCode: roleCode,
			mockBehavior: func(ur *mocks.MockUserRepository, ar *mocks.MockAuthRepository, a *mocks.MockAuth, h *mocks.MockHasher, tx *mock_transactor.MockTransactor) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				h.EXPECT().HashPassword(password).Return(hash, nil)

				ur.EXPECT().Create(ctx, gomock.Any()).Return(entity.User{ID: userID}, nil)
				ur.EXPECT().AttachRole(ctx, userID, roleCode).Return(nil)

				a.EXPECT().GenerateTokens(gomock.Any()).Return(tokens, nil)
				a.EXPECT().HashToken(tokens.RefreshToken).Return(tokenHash)

				ar.EXPECT().SaveRefreshToken(ctx, userID, tokenHash, gomock.Any()).Return(nil)
			},
			want:      tokens,
			wantError: nil,
		},
		{
			name:     "empty email",
			email:    "",
			password: password,
			roleCode: roleCode,
			mockBehavior: func(*mocks.MockUserRepository, *mocks.MockAuthRepository, *mocks.MockAuth, *mocks.MockHasher, *mock_transactor.MockTransactor) {
				// No expectations
			},
			want:      nil,
			wantError: service.ErrEmptyEmail,
		},
		{
			name:     "empty password",
			email:    email,
			password: "",
			roleCode: roleCode,
			mockBehavior: func(*mocks.MockUserRepository, *mocks.MockAuthRepository, *mocks.MockAuth, *mocks.MockHasher, *mock_transactor.MockTransactor) {
				// No expectations
			},
			want:      nil,
			wantError: service.ErrEmptyPassword,
		},
		{
			name:     "empty role code",
			email:    email,
			password: password,
			roleCode: "",
			mockBehavior: func(*mocks.MockUserRepository, *mocks.MockAuthRepository, *mocks.MockAuth, *mocks.MockHasher, *mock_transactor.MockTransactor) {
				// No expectations
			},
			want:      nil,
			wantError: service.ErrEmptyRoleCode,
		},
		{
			name:     "hash password error",
			email:    email,
			password: password,
			roleCode: roleCode,
			mockBehavior: func(_ *mocks.MockUserRepository, _ *mocks.MockAuthRepository, _ *mocks.MockAuth, h *mocks.MockHasher, tx *mock_transactor.MockTransactor) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				h.EXPECT().HashPassword(password).Return("", arbitraryErr)
			},
			want:      nil,
			wantError: service.ErrCannotRegisterUser,
		},
		{
			name:     "user already exists",
			email:    email,
			password: password,
			roleCode: roleCode,
			mockBehavior: func(ur *mocks.MockUserRepository, _ *mocks.MockAuthRepository, _ *mocks.MockAuth, h *mocks.MockHasher, tx *mock_transactor.MockTransactor) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				h.EXPECT().HashPassword(password).Return(hash, nil)

				ur.EXPECT().Create(ctx, gomock.Any()).
					Return(entity.User{}, user_repository.ErrUserAlreadyExists)
			},
			want:      nil,
			wantError: service.ErrUserAlreadyExists,
		},
		{
			name:     "role not found",
			email:    email,
			password: password,
			roleCode: roleCode,
			mockBehavior: func(ur *mocks.MockUserRepository, _ *mocks.MockAuthRepository, _ *mocks.MockAuth, h *mocks.MockHasher, tx *mock_transactor.MockTransactor) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				h.EXPECT().HashPassword(password).Return(hash, nil)

				ur.EXPECT().Create(ctx, gomock.Any()).Return(entity.User{ID: userID}, nil)
				ur.EXPECT().AttachRole(ctx, userID, roleCode).
					Return(user_repository.ErrRoleNotFound)
			},
			want:      nil,
			wantError: service.ErrRoleNotFound,
		},
		{
			name:     "attach role generic error",
			email:    email,
			password: password,
			roleCode: roleCode,
			mockBehavior: func(ur *mocks.MockUserRepository, _ *mocks.MockAuthRepository, _ *mocks.MockAuth, h *mocks.MockHasher, tx *mock_transactor.MockTransactor) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				h.EXPECT().HashPassword(password).Return(hash, nil)

				ur.EXPECT().Create(ctx, gomock.Any()).Return(entity.User{ID: userID}, nil)
				ur.EXPECT().AttachRole(ctx, userID, roleCode).Return(arbitraryErr)
			},
			want:      nil,
			wantError: service.ErrCannotRegisterUser,
		},
		{
			name:     "generate tokens error",
			email:    email,
			password: password,
			roleCode: roleCode,
			mockBehavior: func(ur *mocks.MockUserRepository, _ *mocks.MockAuthRepository, a *mocks.MockAuth, h *mocks.MockHasher, tx *mock_transactor.MockTransactor) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				h.EXPECT().HashPassword(password).Return(hash, nil)

				ur.EXPECT().Create(ctx, gomock.Any()).Return(entity.User{ID: userID}, nil)
				ur.EXPECT().AttachRole(ctx, userID, roleCode).Return(nil)

				a.EXPECT().GenerateTokens(gomock.Any()).Return(nil, arbitraryErr)
			},
			want:      nil,
			wantError: service.ErrCannotRegisterUser,
		},
		{
			name:     "save refresh token error",
			email:    email,
			password: password,
			roleCode: roleCode,
			mockBehavior: func(ur *mocks.MockUserRepository, ar *mocks.MockAuthRepository, a *mocks.MockAuth, h *mocks.MockHasher, tx *mock_transactor.MockTransactor) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				h.EXPECT().HashPassword(password).Return(hash, nil)

				ur.EXPECT().Create(ctx, gomock.Any()).Return(entity.User{ID: userID}, nil)
				ur.EXPECT().AttachRole(ctx, userID, roleCode).Return(nil)

				a.EXPECT().GenerateTokens(gomock.Any()).Return(tokens, nil)
				a.EXPECT().HashToken(tokens.RefreshToken).Return(tokenHash)

				ar.EXPECT().SaveRefreshToken(ctx, userID, tokenHash, gomock.Any()).
					Return(arbitraryErr)
			},
			want:      nil,
			wantError: service.ErrCannotRegisterUser,
		},
		{
			name:     "create user generic error",
			email:    email,
			password: password,
			roleCode: roleCode,
			mockBehavior: func(ur *mocks.MockUserRepository, _ *mocks.MockAuthRepository, _ *mocks.MockAuth, h *mocks.MockHasher, tx *mock_transactor.MockTransactor) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				h.EXPECT().HashPassword(password).Return(hash, nil)

				ur.EXPECT().Create(ctx, gomock.Any()).Return(entity.User{}, arbitraryErr)
			},
			want:      nil,
			wantError: service.ErrCannotRegisterUser,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ur := mocks.NewMockUserRepository(ctrl)
			ar := mocks.NewMockAuthRepository(ctrl)
			a := mocks.NewMockAuth(ctrl)
			h := mocks.NewMockHasher(ctrl)
			tx := mock_transactor.NewMockTransactor(ctrl)

			tc.mockBehavior(ur, ar, a, h, tx)

			s := service.New(ur, ar, tx, a, h)

			out, err := s.Register(ctx, tc.email, tc.password, tc.roleCode)
			assert.ErrorIs(t, err, tc.wantError)
			assert.Equal(t, tc.want, out)
		})
	}
}

func TestService_Login(t *testing.T) {
	var (
		ctx          = context.Background()
		email        = "test@mail.com"
		password     = "password123"
		userID       = uuid.New()
		hash         = "hashed_password"
		tokenHash    = "hashed_refresh"
		arbitraryErr = errors.New("arbitrary error")
		user         = entity.User{ID: userID, Email: email, PasswordHash: hash, IsActive: true}
		inactiveUser = entity.User{ID: userID, Email: email, PasswordHash: hash, IsActive: false}
		tokens       = &auth.Tokens{
			AccessToken:  "access",
			RefreshToken: "refresh",
			ExpiresIn:    900,
		}
	)

	type MockBehavior func(
		ur *mocks.MockUserRepository,
		ar *mocks.MockAuthRepository,
		a *mocks.MockAuth,
		h *mocks.MockHasher,
		tx *mock_transactor.MockTransactor,
	)

	for _, tc := range []struct {
		name         string
		email        string
		password     string
		mockBehavior MockBehavior
		want         *auth.Tokens
		wantErr      error
	}{
		{
			name:     "success",
			email:    email,
			password: password,
			mockBehavior: func(ur *mocks.MockUserRepository, ar *mocks.MockAuthRepository, a *mocks.MockAuth, h *mocks.MockHasher, tx *mock_transactor.MockTransactor) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) })
				ur.EXPECT().GetByEmail(ctx, email).Return(user, nil)
				h.EXPECT().CheckPasswordHash(password, hash).Return(true)
				a.EXPECT().GenerateTokens(user).Return(tokens, nil)
				a.EXPECT().HashToken(tokens.RefreshToken).Return(tokenHash)
				ar.EXPECT().SaveRefreshToken(ctx, userID, tokenHash, gomock.Any()).Return(nil)
			},
			want:    tokens,
			wantErr: nil,
		},
		{
			name:     "empty email",
			email:    "",
			password: password,
			mockBehavior: func(*mocks.MockUserRepository, *mocks.MockAuthRepository, *mocks.MockAuth, *mocks.MockHasher, *mock_transactor.MockTransactor) {
				// No expectations
			},
			want:    nil,
			wantErr: service.ErrEmptyEmail,
		},
		{
			name:     "empty password",
			email:    email,
			password: "",
			mockBehavior: func(*mocks.MockUserRepository, *mocks.MockAuthRepository, *mocks.MockAuth, *mocks.MockHasher, *mock_transactor.MockTransactor) {
				// No expectations
			},
			want:    nil,
			wantErr: service.ErrEmptyPassword,
		},
		{
			name:     "user not found",
			email:    email,
			password: password,
			mockBehavior: func(ur *mocks.MockUserRepository, _ *mocks.MockAuthRepository, _ *mocks.MockAuth, _ *mocks.MockHasher, tx *mock_transactor.MockTransactor) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) })
				ur.EXPECT().GetByEmail(ctx, email).Return(entity.User{}, user_repository.ErrUserNotFound)
			},
			want:    nil,
			wantErr: service.ErrUserNotFound,
		},
		{
			name:     "inactive user",
			email:    email,
			password: password,
			mockBehavior: func(ur *mocks.MockUserRepository, _ *mocks.MockAuthRepository, _ *mocks.MockAuth, _ *mocks.MockHasher, tx *mock_transactor.MockTransactor) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) })
				ur.EXPECT().GetByEmail(ctx, email).Return(inactiveUser, nil)
			},
			want:    nil,
			wantErr: service.ErrInvalidCredentials,
		},
		{
			name:     "invalid password",
			email:    email,
			password: password,
			mockBehavior: func(ur *mocks.MockUserRepository, _ *mocks.MockAuthRepository, _ *mocks.MockAuth, h *mocks.MockHasher, tx *mock_transactor.MockTransactor) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) })
				ur.EXPECT().GetByEmail(ctx, email).Return(user, nil)
				h.EXPECT().CheckPasswordHash(password, hash).Return(false)
			},
			want:    nil,
			wantErr: service.ErrInvalidCredentials,
		},
		{
			name:     "generate tokens error",
			email:    email,
			password: password,
			mockBehavior: func(ur *mocks.MockUserRepository, _ *mocks.MockAuthRepository, a *mocks.MockAuth, h *mocks.MockHasher, tx *mock_transactor.MockTransactor) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) })
				ur.EXPECT().GetByEmail(ctx, email).Return(user, nil)
				h.EXPECT().CheckPasswordHash(password, hash).Return(true)
				a.EXPECT().GenerateTokens(user).Return(nil, arbitraryErr)
			},
			want:    nil,
			wantErr: service.ErrInvalidCredentials,
		},
		{
			name:     "save refresh token error",
			email:    email,
			password: password,
			mockBehavior: func(ur *mocks.MockUserRepository, ar *mocks.MockAuthRepository, a *mocks.MockAuth, h *mocks.MockHasher, tx *mock_transactor.MockTransactor) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) })
				ur.EXPECT().GetByEmail(ctx, email).Return(user, nil)
				h.EXPECT().CheckPasswordHash(password, hash).Return(true)
				a.EXPECT().GenerateTokens(user).Return(tokens, nil)
				a.EXPECT().HashToken(tokens.RefreshToken).Return(tokenHash)
				ar.EXPECT().SaveRefreshToken(ctx, userID, tokenHash, gomock.Any()).Return(arbitraryErr)
			},
			want:    nil,
			wantErr: service.ErrInvalidCredentials,
		},
		{
			name:     "get user by email generic error",
			email:    email,
			password: password,
			mockBehavior: func(ur *mocks.MockUserRepository, _ *mocks.MockAuthRepository, _ *mocks.MockAuth, _ *mocks.MockHasher, tx *mock_transactor.MockTransactor) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) })
				ur.EXPECT().GetByEmail(ctx, email).Return(entity.User{}, arbitraryErr)
			},
			want:    nil,
			wantErr: service.ErrInvalidCredentials,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ur := mocks.NewMockUserRepository(ctrl)
			ar := mocks.NewMockAuthRepository(ctrl)
			a := mocks.NewMockAuth(ctrl)
			h := mocks.NewMockHasher(ctrl)
			tx := mock_transactor.NewMockTransactor(ctrl)

			tc.mockBehavior(ur, ar, a, h, tx)

			s := service.New(ur, ar, tx, a, h)
			out, err := s.Login(ctx, tc.email, tc.password)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, out)
		})
	}
}

func TestService_Refresh(t *testing.T) {
	var (
		ctx          = context.Background()
		userID       = uuid.New()
		email        = "test@mail.com"
		tokenHash    = "hashed_refresh"
		arbitraryErr = errors.New("arbitrary error")
		user         = entity.User{ID: userID, Email: email, IsActive: true}
		inactiveUser = entity.User{ID: userID, Email: email, IsActive: false}
		tokens       = &auth.Tokens{
			AccessToken:  "access",
			RefreshToken: "refresh",
			ExpiresIn:    900,
		}
	)

	type MockBehavior func(
		ur *mocks.MockUserRepository,
		ar *mocks.MockAuthRepository,
		a *mocks.MockAuth,
		tx *mock_transactor.MockTransactor,
	)

	for _, tc := range []struct {
		name         string
		token        string
		mockBehavior MockBehavior
		want         *auth.Tokens
		wantErr      error
	}{
		{
			name:  "success",
			token: "refresh",
			mockBehavior: func(ur *mocks.MockUserRepository, ar *mocks.MockAuthRepository, a *mocks.MockAuth, tx *mock_transactor.MockTransactor) {
				a.EXPECT().ValidateRefreshToken("refresh").Return(email, nil)

				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) })

				a.EXPECT().HashToken("refresh").Return(tokenHash)
				ar.EXPECT().GetUserByRefreshToken(ctx, tokenHash).Return(userID, nil)
				ur.EXPECT().GetByID(ctx, userID).Return(user, nil)
				ar.EXPECT().RevokeRefreshToken(ctx, tokenHash).Return(nil)
				a.EXPECT().GenerateTokens(user).Return(tokens, nil)
				a.EXPECT().HashToken(tokens.RefreshToken).Return(tokenHash)
				ar.EXPECT().SaveRefreshToken(ctx, userID, tokenHash, gomock.Any()).Return(nil)
			},
			want:    tokens,
			wantErr: nil,
		},
		{
			name:  "empty token",
			token: "",
			mockBehavior: func(*mocks.MockUserRepository, *mocks.MockAuthRepository, *mocks.MockAuth, *mock_transactor.MockTransactor) {
				// No expectations
			},
			want:    nil,
			wantErr: service.ErrEmptyToken,
		},
		{
			name:  "invalid refresh token format",
			token: "refresh",
			mockBehavior: func(_ *mocks.MockUserRepository, _ *mocks.MockAuthRepository, a *mocks.MockAuth, _ *mock_transactor.MockTransactor) {
				a.EXPECT().ValidateRefreshToken("refresh").Return("", arbitraryErr)
			},
			want:    nil,
			wantErr: service.ErrInvalidRefreshToken,
		},
		{
			name:  "token not found or expired",
			token: "refresh",
			mockBehavior: func(_ *mocks.MockUserRepository, ar *mocks.MockAuthRepository, a *mocks.MockAuth, tx *mock_transactor.MockTransactor) {
				a.EXPECT().ValidateRefreshToken("refresh").Return(email, nil)

				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) })

				a.EXPECT().HashToken("refresh").Return(tokenHash)
				ar.EXPECT().GetUserByRefreshToken(ctx, tokenHash).Return(uuid.UUID{}, auth_repository.ErrInvalidRefreshToken)
			},
			want:    nil,
			wantErr: service.ErrInvalidRefreshToken,
		},
		{
			name:  "failed to get user by refresh token",
			token: "refresh",
			mockBehavior: func(_ *mocks.MockUserRepository, ar *mocks.MockAuthRepository, a *mocks.MockAuth, tx *mock_transactor.MockTransactor) {
				a.EXPECT().ValidateRefreshToken("refresh").Return(email, nil)

				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) })

				a.EXPECT().HashToken("refresh").Return(tokenHash)
				ar.EXPECT().GetUserByRefreshToken(ctx, tokenHash).Return(uuid.UUID{}, arbitraryErr)
			},
			want:    nil,
			wantErr: service.ErrInvalidRefreshToken,
		},
		{
			name:  "user not found",
			token: "refresh",
			mockBehavior: func(ur *mocks.MockUserRepository, ar *mocks.MockAuthRepository, a *mocks.MockAuth, tx *mock_transactor.MockTransactor) {
				a.EXPECT().ValidateRefreshToken("refresh").Return(email, nil)

				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) })

				a.EXPECT().HashToken("refresh").Return(tokenHash)
				ar.EXPECT().GetUserByRefreshToken(ctx, tokenHash).Return(userID, nil)
				ur.EXPECT().GetByID(ctx, userID).Return(entity.User{}, user_repository.ErrUserNotFound)
			},
			want:    nil,
			wantErr: service.ErrInvalidRefreshToken,
		},
		{
			name:  "get user by id generic error",
			token: "refresh",
			mockBehavior: func(ur *mocks.MockUserRepository, ar *mocks.MockAuthRepository, a *mocks.MockAuth, tx *mock_transactor.MockTransactor) {
				a.EXPECT().ValidateRefreshToken("refresh").Return(email, nil)

				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) })

				a.EXPECT().HashToken("refresh").Return(tokenHash)
				ar.EXPECT().GetUserByRefreshToken(ctx, tokenHash).Return(userID, nil)
				ur.EXPECT().GetByID(ctx, userID).Return(entity.User{}, arbitraryErr)
			},
			want:    nil,
			wantErr: service.ErrInvalidRefreshToken,
		},
		{
			name:  "inactive user",
			token: "refresh",
			mockBehavior: func(ur *mocks.MockUserRepository, ar *mocks.MockAuthRepository, a *mocks.MockAuth, tx *mock_transactor.MockTransactor) {
				a.EXPECT().ValidateRefreshToken("refresh").Return(email, nil)

				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) })

				a.EXPECT().HashToken("refresh").Return(tokenHash)
				ar.EXPECT().GetUserByRefreshToken(ctx, tokenHash).Return(userID, nil)
				ur.EXPECT().GetByID(ctx, userID).Return(inactiveUser, nil)
			},
			want:    nil,
			wantErr: service.ErrInvalidRefreshToken,
		},
		{
			name:  "revoke token error",
			token: "refresh",
			mockBehavior: func(ur *mocks.MockUserRepository, ar *mocks.MockAuthRepository, a *mocks.MockAuth, tx *mock_transactor.MockTransactor) {
				a.EXPECT().ValidateRefreshToken("refresh").Return(email, nil)

				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) })

				a.EXPECT().HashToken("refresh").Return(tokenHash)
				ar.EXPECT().GetUserByRefreshToken(ctx, tokenHash).Return(userID, nil)
				ur.EXPECT().GetByID(ctx, userID).Return(user, nil)
				ar.EXPECT().RevokeRefreshToken(ctx, tokenHash).Return(arbitraryErr)
			},
			want:    nil,
			wantErr: service.ErrInvalidRefreshToken,
		},
		{
			name:  "generate tokens error",
			token: "refresh",
			mockBehavior: func(ur *mocks.MockUserRepository, ar *mocks.MockAuthRepository, a *mocks.MockAuth, tx *mock_transactor.MockTransactor) {
				a.EXPECT().ValidateRefreshToken("refresh").Return(email, nil)

				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) })

				a.EXPECT().HashToken("refresh").Return(tokenHash)
				ar.EXPECT().GetUserByRefreshToken(ctx, tokenHash).Return(userID, nil)
				ur.EXPECT().GetByID(ctx, userID).Return(user, nil)
				ar.EXPECT().RevokeRefreshToken(ctx, tokenHash).Return(nil)
				a.EXPECT().GenerateTokens(user).Return(nil, arbitraryErr)
			},
			want:    nil,
			wantErr: service.ErrInvalidRefreshToken,
		},
		{
			name:  "save refresh token error",
			token: "refresh",
			mockBehavior: func(ur *mocks.MockUserRepository, ar *mocks.MockAuthRepository, a *mocks.MockAuth, tx *mock_transactor.MockTransactor) {
				a.EXPECT().ValidateRefreshToken("refresh").Return(email, nil)

				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) })

				a.EXPECT().HashToken("refresh").Return(tokenHash)
				ar.EXPECT().GetUserByRefreshToken(ctx, tokenHash).Return(userID, nil)
				ur.EXPECT().GetByID(ctx, userID).Return(user, nil)
				ar.EXPECT().RevokeRefreshToken(ctx, tokenHash).Return(nil)
				a.EXPECT().GenerateTokens(user).Return(tokens, nil)
				a.EXPECT().HashToken(tokens.RefreshToken).Return(tokenHash)
				ar.EXPECT().SaveRefreshToken(ctx, userID, tokenHash, gomock.Any()).Return(arbitraryErr)
			},
			want:    nil,
			wantErr: service.ErrInvalidRefreshToken,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ur := mocks.NewMockUserRepository(ctrl)
			ar := mocks.NewMockAuthRepository(ctrl)
			a := mocks.NewMockAuth(ctrl)
			tx := mock_transactor.NewMockTransactor(ctrl)

			tc.mockBehavior(ur, ar, a, tx)

			s := service.New(ur, ar, tx, a, nil)
			out, err := s.Refresh(ctx, tc.token)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, out)
		})
	}
}

func TestService_Logout(t *testing.T) {
	var (
		ctx          = context.Background()
		refreshToken = "refresh_token"
		tokenHash    = "hashed_refresh"
		arbitraryErr = errors.New("arbitrary error")
	)

	type MockBehavior func(
		ar *mocks.MockAuthRepository,
		a *mocks.MockAuth,
	)

	for _, tc := range []struct {
		name         string
		refreshToken string
		mockBehavior MockBehavior
		wantErr      error
	}{
		{
			name:         "success",
			refreshToken: refreshToken,
			mockBehavior: func(ar *mocks.MockAuthRepository, a *mocks.MockAuth) {
				a.EXPECT().
					HashToken(refreshToken).
					Return(tokenHash)

				ar.EXPECT().
					RevokeRefreshToken(ctx, tokenHash).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name:         "empty token",
			refreshToken: "",
			mockBehavior: func(*mocks.MockAuthRepository, *mocks.MockAuth) {
				// No expectations
			},
			wantErr: service.ErrEmptyToken,
		},
		{
			name:         "invalid token",
			refreshToken: refreshToken,
			mockBehavior: func(ar *mocks.MockAuthRepository, a *mocks.MockAuth) {
				a.EXPECT().
					HashToken(refreshToken).
					Return(tokenHash)

				ar.EXPECT().
					RevokeRefreshToken(ctx, tokenHash).
					Return(auth_repository.ErrInvalidRefreshToken)
			},
			wantErr: service.ErrInvalidRefreshToken,
		},
		{
			name:         "failed to revoke token",
			refreshToken: refreshToken,
			mockBehavior: func(ar *mocks.MockAuthRepository, a *mocks.MockAuth) {
				a.EXPECT().
					HashToken(refreshToken).
					Return(tokenHash)

				ar.EXPECT().
					RevokeRefreshToken(ctx, tokenHash).
					Return(arbitraryErr)
			},
			wantErr: service.ErrCannotRevokeRefreshToken,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ur := mocks.NewMockUserRepository(ctrl)
			ar := mocks.NewMockAuthRepository(ctrl)
			a := mocks.NewMockAuth(ctrl)
			h := mocks.NewMockHasher(ctrl)
			tx := mock_transactor.NewMockTransactor(ctrl)

			tc.mockBehavior(ar, a)

			s := service.New(ur, ar, tx, a, h)

			err := s.Logout(ctx, tc.refreshToken)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}
