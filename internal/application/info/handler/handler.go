package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/j03hanafi/halo-suster/common/configs"
)

type infoHandler struct {
	db *pgxpool.Pool
}

func NewInfoHandler(router fiber.Router, db *pgxpool.Pool) {
	handler := infoHandler{
		db: db,
	}

	infoRouter := router.Group("/info")
	infoRouter.Get("/version", handler.Version)
	infoRouter.Get("/health", handler.Health)
}

func (h infoHandler) Version(c *fiber.Ctx) error {
	versionInfo := version{
		Version: configs.Get().App.Version,
	}

	res := baseResponse{
		Message: "API version",
		Data:    versionInfo,
	}

	return c.JSON(res)
}

func (h infoHandler) Health(c *fiber.Ctx) error {
	res := baseResponse{
		Message: "API is up and running",
	}

	dbData := health{
		Status:     "connected",
		IdleConns:  h.db.Stat().IdleConns(),
		TotalConns: h.db.Stat().TotalConns(),
		MaxConns:   h.db.Stat().MaxConns(),
	}

	res.Data = dbData

	return c.JSON(res)
}
