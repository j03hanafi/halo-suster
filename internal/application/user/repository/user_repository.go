package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"

	"github.com/j03hanafi/halo-suster/common/id"
	"github.com/j03hanafi/halo-suster/common/logger"
	"github.com/j03hanafi/halo-suster/internal/domain"
)

type UserRepository struct {
	db       *pgxpool.Pool
	jwtCache *cache.Cache
}

func NewUserRepository(db *pgxpool.Pool, jwtCache *cache.Cache) *UserRepository {
	return &UserRepository{db: db, jwtCache: jwtCache}
}

func (r UserRepository) Register(ctx context.Context, user *domain.User) error {
	callerInfo := "[UserRepository.Register]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	user.ID = id.New()
	user.CreatedAt = time.Now()

	insertQuery := `INSERT INTO users (id, nip, name, password, is_it, img_url, created_at) 
		VALUES (@id, @nip, @name, @password, @is_it, @img_url, @created_at)`
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
	rows, err := r.db.Query(ctx, selectQuery, args)
	if err != nil {
		l.Error("failed to login user", zap.Error(err))
		return dUser, err
	}

	*mUser, err = pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[user])
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

	updateQuery := `UPDATE users SET nip = @nip, name = @name WHERE id = @id AND nip LIKE '303%'`
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

func (r UserRepository) GetUsers(
	ctx context.Context,
	filter *domain.FilterUser,
	users domain.Users,
) (domain.Users, error) {
	callerInfo := "[UserRepository.GetUsers]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	conditions, params := r.filterUser(filter)
	getQuery := `SELECT id, nip, name, created_at FROM users` + conditions

	rows, err := r.db.Query(ctx, getQuery, params)
	if err != nil {
		l.Error("failed to get users", zap.Error(err))
		return users, err
	}

	dUser := domain.UserAcquire()
	defer domain.UserRelease(dUser)

	_, err = pgx.ForEachRow(rows, []any{&dUser.ID, &dUser.NIP, &dUser.Name, &dUser.CreatedAt}, func() error {
		users = append(users, *dUser)
		return nil
	})
	if err != nil {
		l.Error("failed to get users", zap.Error(err))
		return users, err
	}

	return users, nil
}

func (r UserRepository) filterUser(filter *domain.FilterUser) (string, pgx.NamedArgs) {
	const totalConditions = 4
	conditions, params := make([]string, 0, totalConditions), pgx.NamedArgs{}

	if !id.IsZero(filter.UserID) {
		conditions = append(conditions, "id = @id")
		params["id"] = filter.UserID
	}

	if filter.Name != "" {
		conditions = append(conditions, "name ILIKE @name")
		params["name"] = "%" + filter.Name + "%"
	}

	if filter.NIP != "" {
		conditions = append(conditions, "nip LIKE @nip")
		params["nip"] = filter.NIP + "%"
	}

	if filter.Role != "" {
		conditions = append(conditions, "is_it = @is_it")
		params["is_it"] = filter.Role == domain.RoleIT
	}

	order := " ORDER BY created_at DESC"
	if filter.CreatedAt != "" && (filter.CreatedAt == "asc" || filter.CreatedAt == "desc") {
		order = " ORDER BY created_at " + filter.CreatedAt
	}

	const totalLimitOffset = 2
	limitOffset := make([]string, 0, totalLimitOffset)

	limitOffset = append(limitOffset, "LIMIT @limit")
	params["limit"] = 5
	if filter.Limit != 0 {
		params["limit"] = filter.Limit
	}

	if filter.Offset != 0 {
		limitOffset = append(limitOffset, "OFFSET @offset")
		params["offset"] = filter.Offset
	}

	queryConditions := ""
	if len(conditions) > 0 {
		queryConditions = " WHERE " + strings.Join(conditions, " AND ")
	}

	queryConditions += order

	if len(limitOffset) > 0 {
		queryConditions += " " + strings.Join(limitOffset, " ")
	}

	return queryConditions, params
}

func (r UserRepository) SaveJWTCache(_ context.Context, token string, user *domain.User) {
	userClaim := domain.User{
		ID:   user.ID,
		NIP:  user.NIP,
		Name: user.Name,
		Role: user.Role,
	}
	r.jwtCache.Set(token, userClaim, 0)
}

var _ UserRepositoryContract = (*UserRepository)(nil)
