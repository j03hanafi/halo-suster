package application

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/j03hanafi/halo-suster/common/configs"
	"github.com/j03hanafi/halo-suster/internal/application/info"
)

func New(server *fiber.App, db *pgxpool.Pool, jwtMiddleware fiber.Handler) {
	router := server.Group(configs.Get().API.BaseURL)

	info.NewModule(router, db)
}
