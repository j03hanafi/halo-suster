package medical

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/j03hanafi/halo-suster/common/configs"
	"github.com/j03hanafi/halo-suster/internal/application/medical/handler"
	"github.com/j03hanafi/halo-suster/internal/application/medical/repository"
	"github.com/j03hanafi/halo-suster/internal/application/medical/service"
)

func NewModule(router fiber.Router, db *pgxpool.Pool, jwtMiddleware fiber.Handler) {
	ctxTimeout := time.Duration(configs.Get().App.ContextTimeout) * time.Second

	medicalRepository := repository.NewMedicalRepository(db)
	medicalService := service.NewMedicalService(ctxTimeout, medicalRepository)
	handler.NewMedicalHandler(router, jwtMiddleware, medicalService)
}
