package user_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/4udiwe/avito-pvz/pkg/postgres"
	"github.com/4udiwe/coworking/auth-service/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *UserRepository {
	return &UserRepository{pg}
}

func (r *UserRepository) Create(
	ctx context.Context,
	user entity.User,
) (entity.User, error) {

	logrus.Infof("Creating user: %s", user.Email)

	query, args, _ := r.Builder.
		Insert("users").
		Columns("email", "password_hash").
		Values(user.Email, user.PasswordHash).
		Suffix("RETURNING id, email, password_hash, is_active, created_at, updated_at").
		ToSql()

	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			logrus.WithField("email", user.Email).Warn("User already exists")
			return entity.User{}, ErrUserAlreadyExists
		}
		logrus.WithError(err).WithField("email", user.Email).Error("Failed to create user")
		return entity.User{}, fmt.Errorf("create user: %w", err)
	}

	logrus.Infof("User created: %s (%s)", user.ID, user.Email)

	return user, nil
}

func (r *UserRepository) AttachRole(
	ctx context.Context,
	userID uuid.UUID,
	roleCode string,
) error {

	logrus.Infof("Attaching role %s to user %s", roleCode, userID)

	query, args, _ := r.Builder.
		Insert("user_roles").
		Columns("user_id", "role_id").
		Select(
			r.Builder.
				Select("?", "id").
				From("roles").
				Where("code = ?", roleCode),
		).
		ToSql()

	cmd, err := r.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{"user_id": userID, "role": roleCode}).Error("AttachRole failed")
		return fmt.Errorf("attach role: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		// roleCode не найден
		logrus.WithField("role", roleCode).Warn("AttachRole: role not found for user")
		return ErrRoleNotFound
	}

	logrus.Infof("Role %s attached to user %s", roleCode, userID)
	return nil
}

func (r *UserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (entity.User, error) {

	logrus.Infof("Fetching user by email: %s", email)

	query, args, _ := r.Builder.
		Select(
			"u.id",
			"u.email",
			"u.password_hash",
			"u.is_active",
			"u.created_at",
			"u.updated_at",
			"r.id",
			"r.code",
			"r.name",
		).
		From("users u").
		Join("user_roles ur ON ur.user_id = u.id").
		Join("roles r ON r.id = ur.role_id").
		Where("u.email = ?", email).
		ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, query, args...)
	if err != nil {
		logrus.WithError(err).WithField("email", email).Error("GetByEmail: query failed")
		return entity.User{}, err
	}
	defer rows.Close()

	var user entity.User
	var roles []entity.Role

	for rows.Next() {
		var role entity.Role

		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
			&role.ID,
			&role.Code,
			&role.Name,
		); err != nil {
			logrus.WithError(err).WithField("email", email).Error("GetByEmail: row scan failed")
			return entity.User{}, err
		}

		roles = append(roles, role)
	}

	if user.ID == uuid.Nil {
		logrus.WithField("email", email).Warn("GetByEmail: user not found")
		return entity.User{}, ErrUserNotFound
	}

	user.Roles = roles
	logrus.WithFields(logrus.Fields{"user_id": user.ID, "email": user.Email}).Info("User fetched by email")
	return user, nil
}

func (r *UserRepository) GetByID(
	ctx context.Context,
	userID uuid.UUID,
) (entity.User, error) {

	logrus.Infof("Fetching user by ID: %s", userID)
	query, args, _ := r.Builder.
		Select(
			"id",
			"email",
			"password_hash",
			"is_active",
			"created_at",
			"updated_at",
		).
		From("users").
		Where("id = ?", userID).
		ToSql()
	var user entity.User
	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logrus.WithField("user_id", userID).Warn("GetByID: user not found")
			return entity.User{}, ErrUserNotFound
		}
		logrus.WithError(err).WithField("user_id", userID).Error("GetByID: query failed")
		return entity.User{}, err
	}
	logrus.WithFields(logrus.Fields{"user_id": user.ID, "email": user.Email}).Info("User fetched by ID")
	return user, nil
}
