package user

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/j03hanafi/halo-suster/common/configs"
	"github.com/j03hanafi/halo-suster/internal/application/user/handler"
	"github.com/j03hanafi/halo-suster/internal/application/user/repository"
	"github.com/j03hanafi/halo-suster/internal/application/user/service"
)

func NewModule(router fiber.Router, db *pgxpool.Pool, jwtMiddleware fiber.Handler) {
	ctxTimeout := time.Duration(configs.Get().App.ContextTimeout) * time.Second

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(ctxTimeout, userRepository)
	handler.NewUserHandler(router, jwtMiddleware, userService)
}
