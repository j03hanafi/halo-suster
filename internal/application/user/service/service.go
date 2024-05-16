package service

import (
	"context"

	"github.com/j03hanafi/halo-suster/internal/domain"
)

type UserServiceContract interface {
	RegisterIT(ctx context.Context, user *domain.User) error
	GenerateToken(ctx context.Context, user *domain.User) (string, error)
	LoginIT(ctx context.Context, user *domain.User) (*domain.User, error)
	LoginNurse(ctx context.Context, user *domain.User) (*domain.User, error)

	RegisterNurse(ctx context.Context, user *domain.User) error
	UpdateNurse(ctx context.Context, user *domain.User) error
	DeleteNurse(ctx context.Context, user *domain.User) error
	UpdateAccess(ctx context.Context, user *domain.User) error
	GetUsers(ctx context.Context, filter *domain.FilterUser, users domain.Users) (domain.Users, error)
}
