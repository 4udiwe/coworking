package user_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/4udiwe/avito-pvz/pkg/transactor"
	"github.com/4udiwe/coworking/auth-service/internal/entity"
	user_repository "github.com/4udiwe/coworking/auth-service/internal/repository/user"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Service struct {
	userRepo UserRepository
	tx       transactor.Transactor
}

func New(
	userRepo UserRepository,
	tx transactor.Transactor,
) *Service {
	return &Service{
		userRepo: userRepo,
		tx:       tx,
	}
}

func (s *Service) GetUsers(
	ctx context.Context,
	page, pageSize int,
	searchQuery *string,
	filterRole *string,
	filterIsActive *bool,
	sortField *string,
) (users []entity.User, total int64, err error) {

	logrus.WithFields(logrus.Fields{
		"page":  page,
		"limit": pageSize,
		"query": searchQuery,
	}).Info("GetUsers called")

	users, total, err = s.userRepo.GetUsers(ctx, page, pageSize, searchQuery, filterRole, sortField, filterIsActive)

	if err != nil {
		logrus.WithError(err).Error("failed to get users")
		return nil, 0, ErrCannotFetchUsers
	}

	return users, total, nil
}

func (s *Service) SetUserActive(
	ctx context.Context,
	userID uuid.UUID,
	active bool,
) error {

	logrus.WithFields(logrus.Fields{
		"userID": userID,
		"active": active,
	}).Info("SetUserActive called")

	err := s.userRepo.SetActive(ctx, userID, active)
	if err != nil {
		if errors.Is(err, user_repository.ErrUserNotFound) {
			logrus.WithField("userID", userID).Warn("user not found")
			return ErrUserNotFound
		}

		logrus.WithError(err).Error("failed to update user active status")
		return fmt.Errorf("set active: %w", err)
	}

	return nil
}

func (s *Service) UpdateUserRoles(
	ctx context.Context,
	userID uuid.UUID,
	roles []string,
) error {

	logrus.WithFields(logrus.Fields{
		"userID": userID,
		"roles":  roles,
	}).Info("UpdateUserRoles called")

	if len(roles) == 0 {
		return ErrEmptyRoles
	}

	return s.tx.WithinTransaction(ctx, func(ctx context.Context) error {

		// Проверяем пользователя
		_, err := s.userRepo.GetByID(ctx, userID)
		if err != nil {
			if errors.Is(err, user_repository.ErrUserNotFound) {
				return ErrUserNotFound
			}
			return fmt.Errorf("get user: %w", err)
		}

		// 🔥 Удаляем старые роли
		if err := s.userRepo.ClearRoles(ctx, userID); err != nil {
			logrus.WithError(err).Error("failed to clear roles")
			return fmt.Errorf("clear roles: %w", err)
		}

		// 🔥 Добавляем новые роли
		for _, roleCode := range roles {
			if err := s.userRepo.AttachRole(ctx, userID, roleCode); err != nil {
				if errors.Is(err, user_repository.ErrRoleNotFound) {
					logrus.WithField("role", roleCode).Warn("role not found")
					return ErrRoleNotFound
				}
				return fmt.Errorf("attach role: %w", err)
			}
		}

		return nil
	})
}

func (s *Service) GetUserInfo(
	ctx context.Context,
	userID uuid.UUID,
) (entity.User, error) {
	logrus.Info("GetUserInfo attempt")

	if userID == uuid.Nil {
		return entity.User{}, ErrEmptyUserID
	}

	logrus.WithField("user_id", userID).Info("Fetching user by ID")
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return entity.User{}, ErrUserNotFound
	}

	return user, nil
}
