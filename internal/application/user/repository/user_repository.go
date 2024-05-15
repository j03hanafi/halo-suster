package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/j03hanafi/halo-suster/common/id"
	"github.com/j03hanafi/halo-suster/common/logger"
	"github.com/j03hanafi/halo-suster/internal/domain"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r UserRepository) Register(ctx context.Context, user *domain.User) error {
	callerInfo := "[UserRepository.Register]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	user.ID = id.New()
	user.CreatedAt = time.Now()

	insertQuery := `INSERT INTO users (id, nip, name, password, is_it, img_url, created_at) VALUES (@id, @nip, @name, @password, @is_it, @img_url, @created_at)`
	args := pgx.NamedArgs{
		"id":         user.ID,
		"nip":        user.NIP,
		"name":       user.Name,
		"password":   user.Password,
		"is_it":      user.Role == domain.RoleIT,
		"img_url":    user.ImgURL,
		"created_at": user.CreatedAt,
	}

	_, err := r.db.Exec(ctx, insertQuery, args)
	if err != nil {
		l.Error("failed to register user", zap.Error(err))

		pgErr := &pgconn.PgError{}
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return new(domain.ErrDuplicateNIP)
		}

		return err
	}

	return nil
}

func (r UserRepository) GetByNIP(ctx context.Context, dUser *domain.User) (*domain.User, error) {
	callerInfo := "[UserRepository.GetByNIP]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	mUser := userAcquire()
	defer userRelease(mUser)

	selectQuery := `SELECT id, nip, name, password, is_it FROM users WHERE nip = @nip`
	args := pgx.NamedArgs{"nip": dUser.NIP}
	query, err := r.db.Query(ctx, selectQuery, args)
	if err != nil {
		l.Error("failed to login user", zap.Error(err))
		return dUser, err
	}

	*mUser, err = pgx.CollectOneRow(query, pgx.RowToStructByNameLax[user])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dUser, new(domain.ErrUserNotFound)
		}

		l.Error("failed to login user", zap.Error(err))
		return dUser, err
	}

	dUser.ID = mUser.ID
	dUser.Name = mUser.Name
	dUser.Password = mUser.Password

	if mUser.IsIT {
		dUser.Role = domain.RoleIT
	} else {
		dUser.Role = domain.RoleNurse
	}

	return dUser, nil
}

func (r UserRepository) UpdateNurse(ctx context.Context, user *domain.User) error {
	callerInfo := "[UserRepository.UpdateNurse]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	updateQuery := `UPDATE users SET nip = @nip, name = @name WHERE id = @id`
	args := pgx.NamedArgs{
		"id":   user.ID,
		"nip":  user.NIP,
		"name": user.Name,
	}

	result, err := r.db.Exec(ctx, updateQuery, args)
	if err != nil {
		l.Error("failed to update user", zap.Error(err))

		pgErr := &pgconn.PgError{}
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return new(domain.ErrDuplicateNIP)
		}

		return err
	}

	if result.RowsAffected() == 0 {
		l.Error("user not found")
		return new(domain.ErrUserNotFound)
	}

	return nil
}

func (r UserRepository) DeleteNurse(ctx context.Context, user *domain.User) error {
	callerInfo := "[UserRepository.DeleteNurse]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	deleteQuery := `DELETE FROM users WHERE id = @id AND nip LIKE '303%'`
	args := pgx.NamedArgs{"id": user.ID}

	result, err := r.db.Exec(ctx, deleteQuery, args)
	if err != nil {
		l.Error("failed to delete user", zap.Error(err))
		return err
	}

	if result.RowsAffected() == 0 {
		l.Error("user not found / is not a nurse")
		return new(domain.ErrNotFoundOrNotNurse)
	}

	return nil
}

func (r UserRepository) UpdateAccess(ctx context.Context, user *domain.User) error {
	callerInfo := "[UserRepository.UpdateAccess]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	updateQuery := `UPDATE users SET password = @password WHERE id = @id AND nip LIKE '303%'`
	args := pgx.NamedArgs{
		"id":       user.ID,
		"password": user.Password,
	}

	result, err := r.db.Exec(ctx, updateQuery, args)
	if err != nil {
		l.Error("failed to update user access", zap.Error(err))
		return err
	}

	if result.RowsAffected() == 0 {
		l.Error("user not found / is not a nurse")
		return new(domain.ErrNotFoundOrNotNurse)
	}

	return nil
}

var _ UserRepositoryContract = (*UserRepository)(nil)
