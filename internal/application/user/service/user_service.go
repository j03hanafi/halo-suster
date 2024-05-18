package service

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/j03hanafi/halo-suster/common/logger"
	"github.com/j03hanafi/halo-suster/common/security"
	"github.com/j03hanafi/halo-suster/internal/application/user/repository"
	"github.com/j03hanafi/halo-suster/internal/domain"
)

type UserService struct {
	userRepository repository.UserRepositoryContract
	contextTimeout time.Duration
}

func NewUserService(timeout time.Duration, userRepository repository.UserRepositoryContract) *UserService {
	return &UserService{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (s UserService) RegisterIT(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	callerInfo := "[UserService.RegisterIT]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	password, err := security.HashPassword(user.Password)
	if err != nil {
		l.Error("failed to hash password", zap.Error(err))
		return err
	}

	user.Password = password
	user.Role = domain.RoleIT
	err = s.userRepository.Register(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s UserService) GenerateToken(ctx context.Context, user *domain.User) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	callerInfo := "[UserService.GenerateToken]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	token, err := security.GenerateAccessToken(user)
	if err != nil {
		l.Error("failed to generate token", zap.Error(err))
		return "", err
	}

	s.userRepository.SaveJWTCache(ctx, token, user)

	return token, nil
}

func (s UserService) LoginIT(ctx context.Context, user *domain.User) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	callerInfo := "[UserService.LoginIT]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	suppliedPassword := user.Password

	user, err := s.userRepository.GetByNIP(ctx, user)
	if err != nil {
		l.Error("failed to get user by NIP", zap.Error(err))
		return user, err
	}

	if err = security.ComparePassword(user.Password, suppliedPassword); err != nil {
		l.Error("failed to compare password", zap.Error(err))
		return user, new(domain.ErrInvalidPassword)
	}

	return user, nil
}

func (s UserService) LoginNurse(ctx context.Context, user *domain.User) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	callerInfo := "[UserService.LoginNurse]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	suppliedPassword := user.Password

	user, err := s.userRepository.GetByNIP(ctx, user)
	if err != nil {
		l.Error("failed to get user by NIP", zap.Error(err))
		return user, err
	}

	if user.Password == "" {
		return user, new(domain.ErrAccessNotAllowed)
	}

	if err = security.ComparePassword(user.Password, suppliedPassword); err != nil {
		l.Error("failed to compare password", zap.Error(err))
		return user, new(domain.ErrInvalidPassword)
	}

	return user, nil
}

func (s UserService) RegisterNurse(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	callerInfo := "[UserService.RegisterNurse]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	user.Role = domain.RoleNurse
	err := s.userRepository.Register(ctx, user)
	if err != nil {
		l.Error("failed to register nurse", zap.Error(err))
		return err
	}

	return nil
}

func (s UserService) UpdateNurse(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	callerInfo := "[UserService.UpdateNurse]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	err := s.userRepository.UpdateNurse(ctx, user)
	if err != nil {
		l.Error("failed to update nurse", zap.Error(err))
		return err
	}

	return nil
}

func (s UserService) DeleteNurse(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	callerInfo := "[UserService.DeleteNurse]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	err := s.userRepository.DeleteNurse(ctx, user)
	if err != nil {
		l.Error("failed to delete nurse", zap.Error(err))
		return err
	}

	return nil
}

func (s UserService) UpdateAccess(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	callerInfo := "[UserService.UpdateAccess]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	password, err := security.HashPassword(user.Password)
	if err != nil {
		l.Error("failed to hash password", zap.Error(err))
		return err
	}

	user.Password = password
	err = s.userRepository.UpdateAccess(ctx, user)
	if err != nil {
		l.Error("failed to update access", zap.Error(err))
		return err
	}

	return nil
}

func (s UserService) GetUsers(
	ctx context.Context,
	filter *domain.FilterUser,
	users domain.Users,
) (domain.Users, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	callerInfo := "[UserService.GetUsers]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	users, err := s.userRepository.GetUsers(ctx, filter, users)
	if err != nil {
		l.Error("failed to get users", zap.Error(err))
		return users, err
	}

	return users, nil
}

var _ UserServiceContract = (*UserService)(nil)
