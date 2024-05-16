package repository

import (
	"context"

	"github.com/j03hanafi/halo-suster/internal/domain"
)

type UserRepositoryContract interface {
	Register(ctx context.Context, user *domain.User) error
	GetByNIP(ctx context.Context, user *domain.User) (*domain.User, error)
	UpdateNurse(ctx context.Context, user *domain.User) error
	DeleteNurse(ctx context.Context, user *domain.User) error
	UpdateAccess(ctx context.Context, user *domain.User) error
	GetUsers(ctx context.Context, filter *domain.FilterUser, users domain.Users) (domain.Users, error)
}
