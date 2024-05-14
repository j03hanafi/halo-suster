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
	versionInfo := versionAcquire()
	defer versionRelease(versionInfo)

	versionInfo.Version = configs.Get().App.Version

	res := baseResponseAcquire()
	defer baseResponseRelease(res)

	res.Message = "API version"
	res.Data = versionInfo

	return c.JSON(res)
}

func (h infoHandler) Health(c *fiber.Ctx) error {
	healthInfo := healthAcquire()
	defer healthRelease(healthInfo)

	healthInfo.Status = "connected"
	healthInfo.IdleConns = h.db.Stat().IdleConns()
	healthInfo.TotalConns = h.db.Stat().TotalConns()
	healthInfo.MaxConns = h.db.Stat().MaxConns()

	res := baseResponseAcquire()
	defer baseResponseRelease(res)

	res.Message = "API up and running"
	res.Data = healthInfo

	return c.JSON(res)
}
