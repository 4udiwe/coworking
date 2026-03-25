package user_repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/4udiwe/avito-pvz/pkg/postgres"
	"github.com/4udiwe/coworking/auth-service/internal/entity"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
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
		Columns("email", "first_name", "last_name", "password_hash").
		Values(user.Email, user.FirstName, user.LastName, user.PasswordHash).
		Suffix("RETURNING id, first_name, last_name, email, password_hash, is_active, created_at, updated_at").
		ToSql()

	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
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

	const query = `
		INSERT INTO user_roles (user_id, role_id)
		SELECT $1, r.id
		FROM roles r
		WHERE r.code = $2
		ON CONFLICT DO NOTHING
	`

	cmd, err := r.GetTxManager(ctx).Exec(ctx, query, userID, roleCode)
	if err != nil {
		logrus.WithError(err).
			WithFields(logrus.Fields{
				"user_id": userID,
				"role":    roleCode,
			}).
			Error("AttachRole failed")
		return fmt.Errorf("attach role: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		logrus.WithField("role", roleCode).
			Warn("AttachRole: role not found")
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
			"u.first_name",
			"u.last_name",
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
			&user.FirstName,
			&user.LastName,
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

	logrus.Infof("Fetching user with roles by ID: %s", userID)

	const query = `
		SELECT
			u.id,
			u.first_name,
			u.last_name,
			u.email,
			u.password_hash,
			u.is_active,
			u.created_at,
			u.updated_at,
			r.id         AS role_id,
			r.code       AS role_code,
			r.name       AS role_name,
			r.created_at AS role_created_at
		FROM users u
		LEFT JOIN user_roles ur ON ur.user_id = u.id
		LEFT JOIN roles r ON r.id = ur.role_id
		WHERE u.id = $1
	`

	rows, err := r.GetTxManager(ctx).Query(ctx, query, userID)
	if err != nil {
		logrus.WithError(err).WithField("user_id", userID).Error("GetByID: query failed")
		return entity.User{}, err
	}
	defer rows.Close()

	var (
		user  entity.User
		found bool
	)

	for rows.Next() {
		var raw rawUserRole

		err := rows.Scan(
			&raw.ID,
			&raw.FirstName,
			&raw.LastName,
			&raw.Email,
			&raw.PasswordHash,
			&raw.IsActive,
			&raw.CreatedAt,
			&raw.UpdatedAt,
			&raw.RoleID,
			&raw.RoleCode,
			&raw.RoleName,
			&raw.RoleCreatedAt,
		)
		if err != nil {
			return entity.User{}, err
		}

		if !found {
			user = entity.User{
				ID:           raw.ID,
				FirstName:    raw.FirstName,
				LastName:     raw.LastName,
				Email:        raw.Email,
				PasswordHash: raw.PasswordHash,
				IsActive:     raw.IsActive,
				CreatedAt:    raw.CreatedAt,
				UpdatedAt:    raw.UpdatedAt,
			}
			found = true
		}

		// если роль есть (LEFT JOIN может вернуть NULL)
		if raw.RoleID != uuid.Nil {
			user.Roles = append(user.Roles, entity.Role{
				ID:        raw.RoleID,
				Code:      entity.RoleCode(raw.RoleCode),
				Name:      raw.RoleName,
				CreatedAt: raw.RoleCreatedAt,
			})
		}
	}

	if !found {
		return entity.User{}, ErrUserNotFound
	}

	return user, nil
}

func (r *UserRepository) GetUsers(
	ctx context.Context,
	page, pageSize int,
	searchQuery, filterRole, sortField *string,
	filterIsActive *bool,
) ([]entity.User, int64, error) {

	logrus.WithFields(logrus.Fields{
		"page":           page,
		"pageSize":       pageSize,
		"searchQuery":    searchQuery,
		"filterRole":     filterRole,
		"sortField":      sortField,
		"filterIsActive": filterIsActive,
	}).Info("Fetching users")

	builder := r.Builder.
		Select(
			"u.id",
			"u.first_name",
			"u.last_name",
			"u.email",
			"u.password_hash",
			"u.is_active",
			"u.created_at",
			"u.updated_at",

			// roles aggregation
			`COALESCE(
				json_agg(
					DISTINCT jsonb_build_object(
						'id', r.id,
						'code', r.code,
						'name', r.name,
						'created_at', r.created_at
					)
				) FILTER (WHERE r.id IS NOT NULL),
				'[]'
			) AS roles`,

			// total count
			"COUNT(*) OVER() AS total",
		).
		From("users u").
		LeftJoin("user_roles ur ON ur.user_id = u.id").
		LeftJoin("roles r ON r.id = ur.role_id")

	// --- filters ---
	if searchQuery != nil && *searchQuery != "" {
		builder = builder.Where(squirrel.Or{
			squirrel.ILike{"u.email": "%" + *searchQuery + "%"},
			squirrel.ILike{"u.first_name": "%" + *searchQuery + "%"},
			squirrel.ILike{"u.last_name": "%" + *searchQuery + "%"},
		})
	}

	if filterRole != nil && *filterRole != "" {
		builder = builder.Where(squirrel.Eq{"r.code": *filterRole})
	}

	if filterIsActive != nil {
		builder = builder.Where(squirrel.Eq{"u.is_active": *filterIsActive})
	}

	// --- grouping ---
	builder = builder.GroupBy("u.id")

	// --- sorting ---
	allowedSort := map[string]string{
		"created_at": "u.created_at",
		"email":      "u.email",
		"first_name": "u.first_name",
		"last_name":  "u.last_name",
	}

	orderBy := "u.created_at DESC"
	if sortField != nil {
		if col, ok := allowedSort[*sortField]; ok {
			orderBy = col + " DESC"
		}
	}

	builder = builder.OrderBy(orderBy)

	// --- pagination ---
	offset := (page - 1) * pageSize
	builder = builder.Offset(uint64(offset)).Limit(uint64(pageSize))

	// --- build SQL ---
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build query: %w", err)
	}

	rows, err := r.GetTxManager(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	var users []entity.User
	var total int64

	var firstName sql.NullString
	var lastName sql.NullString

	for rows.Next() {
		var u entity.User
		var rolesJSON []byte

		if err := rows.Scan(
			&u.ID,
			&firstName,
			&lastName,
			&u.Email,
			&u.PasswordHash,
			&u.IsActive,
			&u.CreatedAt,
			&u.UpdatedAt,
			&rolesJSON,
			&total,
		); err != nil {
			return nil, 0, err
		}

		if err := json.Unmarshal(rolesJSON, &u.Roles); err != nil {
			return nil, 0, err
		}

		if firstName.Valid {
			u.FirstName = firstName.String
		}

		if lastName.Valid {
			u.LastName = lastName.String
		}

		users = append(users, u)
	}

	return users, total, nil
}

func (r *UserRepository) SetActive(ctx context.Context, userID uuid.UUID, active bool) error {
	query, args, _ := r.Builder.
		Update("users").
		Set("is_active", active).
		Set("updated_at", time.Now()).
		Where("id = ?", userID).
		ToSql()
	cmd, err := r.GetTxManager(ctx).Exec(ctx, query, args...,
	)
	if err != nil {
		logrus.WithError(err).WithField("user_id", userID).Error("SetActive: query failed")
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) ClearRoles(ctx context.Context, userID uuid.UUID) error {
	const sql = `DELETE FROM user_roles WHERE user_id = $1`
	query, args, _ := r.Builder.
		Delete("user_roles").
		Where("user_id = ?", userID).
		ToSql()

	cmd, err := r.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		logrus.WithError(err).WithField("user_id", userID).Error(": query failed")
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}
