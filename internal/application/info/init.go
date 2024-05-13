package info

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/j03hanafi/halo-suster/internal/application/info/handler"
)

func NewModule(router fiber.Router, db *pgxpool.Pool) {
	handler.NewInfoHandler(router, db)
}
